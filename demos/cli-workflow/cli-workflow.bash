#!/bin/bash

cd $(dirname $0)

#
# Test the cli utilities that demonstrate features of dataset package.
#

for NAME in characters plays datatypes; do
    # Remove stale imported data
    if [ -f "${NAME}.csv" ]; then
        echo "Removing stale ${NAME}.csv"
        rm "${NAME}.csv"
    fi
    # Remove stale index definition
    if [ -d "${NAME}.json" ]; then
        echo "Removing stale ${NAME}.json"
        rm "${NAME}.json"
    fi
    # Remove stale dataset
    if [ -d "${NAME}" ]; then
        echo "Removing stale ${NAME}"
        rm -fR "${NAME}"
    fi
    # Remove stale indexe
    if [ -d "${NAME}.bleve" ]; then
        echo "Removing stale ${NAME}.bleve"
        rm -fR "${NAME}.bleve"
    fi
done

#
# Pick version of cli to test with
#
if [ -f "bin/dataset" ] && [ -f "bin/dsfind" ] && [ -f "bin/dsindexer" ]; then
    export PATH="./bin":$PATH
fi

# Generate CSV test data
cat<<CSV1 > characters.csv
name,email
Jack Flanders,captain.jack@zbs.example.org
Sam Mojo,mojo.sam@zbs.example.org
Frieda Little,little.frieda@zbs.example.org
Old Far-Seeing Art,old.far.seeing.art@zbs.example.org
Doctor Mazoola,dr.mazzola@secret-labs.zbs.example.org
Chief Wampum,chief.wampum@zbs.example.org
Lord Henry Jowls,lord.henry@inverness.zbs.example.org
CSV1

cat<<CSV2 > plays.csv
title,year,characters
The Fourth Tower of Inverness,1972,"Jack Flanders, Little Fredia, Narrator, Dr. Mazoola, Chief Wampum, Old Far-Seeing Art, Lord Henry Jowls, Meanie Eenie, Lady Sarah Jowls, Whirlitzer"
Moon Over Morocco,1974,"Jack Flanders, Kasbah Kelly, Mojo Sam, Little Flossic, Sunny Skies, Layla Oolupi, Queen Azora, Narrator, Storyteller Mustafa, Comtese Zazeenia, Abu, Taxi Driver, Marmaduke"
The Ah-Ha Phenomenon,1977,"Jack Flanders, Sir Seymour Jowls, Cynthia, Hostess, Archivist, Narrator, Troll, Chief Wampum, Wizard"
The Incredible Adventures of Jack Flanders,1978,"Jack Flanders, Little Frieda, Doctor Mazoola, Narrator, Captian Swallow, Marquis of Carambas, Mojo Sam, The Pirate Queen, Old Far-Seeing Art, Chief Wampum, Owl Eyes, Sorcerer, Waitress"
CSV2

cat<<CSV3 > datatypes.csv
string,year,int,float,date,geo
Hello world,1999,12,12.4,1999-12-12,"50.7226968,-4.3453813"
Goodbye world,1834,1,1.5,1834-02-06,"52.3337236,-6.4907425"
Muddling around the world,1901,0,0,1901-12-31,"10.0344444,139.7700281"
Give me liberty and give chow fun noodles,1775,1,1.0,1775-03-23,"34.1376576,-118.127463"
"Caltech Pasadena, California",1891,3,3.0,1891-09-01,"34.138,-118.125"
CSV3

# Generate the index mapping (we're calling it characters.json)
cat<<DEF1 > characters.json
{
    "uuid":{
        "object_path":".uuid",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":true
    },
    "name":{
        "object_path":".name",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":true
    },
    "email":{
        "object_path":".email",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":true
    }
}
DEF1

cat<<DEF2 > plays.json
{
    "uuid":{
        "object_path":".uuid",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":true
    },
    "title": {
        "object_path":".title",
        "field_mapping":"text",
        "analyzer":"standard",
        "store":true
    },
    "year": {
        "object_path":".year",
        "field_mapping":"numeric",
        "store":true
    },
    "characters": {
        "object_path":".characters",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":true
    }
}
DEF2

cat<<DEF3 > datatypes.json
{
    "uuid":{
        "object_path":".uuid",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":true
    },
    "string": {
        "object_path":".string",
        "field_mapping":"text",
        "analyzer":"standard",
        "store":true
    },
    "year": {
        "object_path":".year",
        "field_mapping":"numeric",
        "store":true
    },
    "int": {
        "object_path":".int",
        "field_mapping":"numeric",
        "analyzer":"simple",
        "store":true
    },
    "float": {
        "object_path":".float",
        "field_mapping":"numeric",
        "analyzer":"simple",
        "store":true
    },
    "date": {
        "object_path":".date",
        "field_mapping":"datetime",
        "store":true
    },
    "geo": {
        "object_path":".geo",
        "field_mapping":"geopoint",
        "store":true
    }
}
DEF3

# Initialize an empty repository
$(dataset init characters)
# Load the data
dataset -uuid import characters.csv
dsindexer characters.json
$(dataset init plays)
dataset -uuid import plays.csv
dsindexer plays.json
$(dataset init datatypes)
dataset -uuid import datatypes.csv
dsindexer datatypes.json
unset DATASET
echo

echo "Sorting records by ascending name (simple analyzer)"
echo
dsfind -indexes="characters.bleve" -size 25 -sort="+name" -csv -fields="name,email" "*"
echo

echo "Combining indexes query for \"Jack ZBS\" in characters.bleve and plays.bleve (sort by title,name)"
echo
dsfind -indexes="plays.bleve,characters.bleve" -size 100 -csv -fields="uuid,title,year,characters,name,email" -sort="title,name" "+Jack +ZBS"
echo

#
# Demonstrate formats and interactions with query strings for datatype collection
#
export DATASET=datatypes
echo "Listing datatypes indexes, sort by descending year"
echo
dsfind -csv -fields="uuid,string,year,int,float,date" -sort="-year" "*"
echo

echo "Listing datatypes indexes, sort by descending year, for range 1700 to 1902"
echo
dsfind -csv -fields="uuid,string,year,int,float,date" -sort="-year" "+year:>=1700 +year:<=1902"
echo

echo "Listing in date range 1700 to 1910, ascending dates"
echo
dsfind -csv -fields="uuid,string,year,int,float,date" -sort="+date" '+date:>="1700-01-01" +date:<="1910-12-31"'
echo

echo "Looking for Caltech, Pasadena, CA, USA"
echo
dsfind -indexes="datatypes.bleve" -csv -fields="uuid,string,date,geo.lat,geo.lng" -sort="+geo.lat,+geo.lng" '+geo.lat:34.138 +geo.lng:-118.125'
echo 
