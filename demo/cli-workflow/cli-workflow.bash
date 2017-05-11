#!/bin/bash

cd $(dirname $0)

#
# Test the cli utilities that demonstrate features of dataset package.
#

# Remove stale imported data
if [ -f "characters.csv" ]; then
    rm characters.csv
fi
if [ -f "plays.csv" ]; then
    rm plays.csv
fi
# Remove stale dataset
if [ -d "characters" ]; then
    rm -fR characters
fi
if [ -d "plays" ]; then
    rm -fR plays
fi
# Remove stale indexe
if [ -d "characters.bleve" ]; then
    rm -fR "characters.bleve" 
fi
if [ -d "plays.bleve" ]; then
    rm -fR "plays.bleve"
fi
# Remove stale index definition
if [ -f "characters.json" ]; then
    rm characters.json
fi
if [ -d "plays.json" ]; then
    rm plays.json
fi

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

# Generate the index mapping (we're calling it characters.json)
cat<<DEF1 > characters.json
{
    "name":{
        "object_path":".last_name",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":"true"
    },
    "email":{
        "object_path":".email",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":"true"
    }
}
DEF1

cat<<DEF2 > plays.json
{
    "title": {
        "object_path":".title",
        "field_mapping":"text",
        "analyzer":"standard",
        "store":"true"
    },
    "year": {
        "object_path":".year",
        "field_mapping":"numeric",
        "store":"true"
    },
    "characters": {
        "object_path":".characters",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":"true"
    }
}
DEF2

# Initialize an empty repository
$(dataset init characters)
# Load the data
dataset -uuid import characters.csv
dsindexer characters.json
$(dataset init plays)
dataset -uuid import plays.csv
dsindexer plays.json
unset DATASET

#echo "Sorting records by descending last_name"
#dsfind -size 25 -sort="-last_name" -fields="last_name" "*"
#echo "Sorting records by ascending last_name"
#dsfind -size 25 -sort="last_name" -fields="last_name" "*"
echo "Revsere sort by last name output as CSV file"
dsfind -indexes="plays.bleve:characters.bleve" -size 100 -csv -fields="title:year:characters:name:email" -sort="title:name" "Jack ZBS"
