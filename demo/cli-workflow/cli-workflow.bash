#!/bin/bash

cd $(dirname $0)

#
# Test the cli utilities that demonstrate features of dataset package.
#

# Remove stale imported data
if [ -f "characters-to-import.csv" ]; then
    rm characters-to-import.csv
fi
# Remove stale dataset
if [ -d "characters" ]; then
    rm -fR characters
fi
# Remove stale indexe
if [ -d "characters.bleve" ]; then
    rm -fR "characters.bleve" 
fi
# Remove stale index definition
if [ -f "characters.json" ]; then
    rm characters.json
fi

#
# Pick version of cli to test with
#
if [ -f "bin/dataset" ] && [ -f "bin/dsfind" ] && [ -f "bin/dsindexer" ]; then
    export PATH="./bin":$PATH
fi

# Generate CSV test data
cat<<FILE1 > characters.csv
last_name,first_name,email
sam,mojo,mojo.sam@zbs.example.org
frieda,little,little.frieda@zbs.example.org
flanders,Jack,captain.jack@zbs.example.org
art,far seeing,old.far.seeing.art@zbs.example.org
shoe,ruby,ruby2@zbs.example.org
turu,t.j.,t.j.turu@zbs.example.org
andover,rhodes,arhodes@another.example.org
kapur,rodant,rodant.kapur@zbs.example.org
Li,Ho,ho.li@scientists.example.org
Lee,Ho,ho.lee@scientists.example.org
Oh,Mark,mark.oh@scientists.example.org
Oh,Ann,ann.oh@scientists.example.org
or,and,andor@digital-circus.example.org
an,and,andan@digital-circus.example.org
Oi,And,andoi@digital-circus.example.org
On,And,andon@digital-circus.example.org
Of,And,andof@digital-circus.example.org
the,off,offthe@digital-circus.example.org
FILE1

# Generate the index mapping (we're calling it characters.json)
cat<<FILE2 > characters.json
{
    "last_name":{
        "object_path":".last_name",
        "field_mapping":"text",
        "analyzer":"simple",
        "store":"true"
    },
    "first_name":{
        "object_path":".first_name",
        "field_mapping":"text",
        "analyzer":"standard",
        "store":"true"
    },
    "email":{
        "object_path":".email",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":"true"
    }
}
FILE2

# Initialize an empty repository
$(dataset init characters)
# Load the data
dataset -uuid import characters.csv
dsindexer characters.json
#echo "Sorting records by descending last_name"
#dsfind -size 25 -sort="-last_name" -fields="last_name" "*"
#echo "Sorting records by ascending last_name"
#dsfind -size 25 -sort="last_name" -fields="last_name" "*"
echo "Revsere sort by last name output as CSV file"
dsfind -size 25 -csv -fields="last_name:first_name:email" -sort="last_name:first_name" "*"
