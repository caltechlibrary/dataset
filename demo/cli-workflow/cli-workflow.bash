#!/bin/bash

cd $(dirname $0)
#
# Test the cli utilities that demonstrate features of dataset package.
#
if [ -d "characters-index.bleve" ]; then
    rm -fR "characters-index.bleve" 
fi
if [ -d "characters" ]; then
    rm -fR characters
fi
if [ -f "characters.csv" ]; then
    rm characters.csv
fi
if [ -f "characters-index.json" ]; then
    rm characters-index.json
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
Or,Mark,mark.or@scientists.example.org
Or,Ann,ann.or@scientists.example.org
Or,And,andor@digital-circus.example.org
FILE1

# Generate the index mapping
cat<<FILE2 > characters-index.json
{
    "last_name":{
        "object_path":".last_name",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":"true"
    },
    "first_name":{
        "object_path":".first_name",
        "field_mapping":"text",
        "analyzer":"keyword",
        "store":"true"
    },
    "email":{
        "object_path":".email",
        "analyzer":"simple",
        "store":"true"
    }
}
FILE2

# Initialize an empty repository
$(dataset init characters)
# Load the data
dataset -uuid import characters.csv
dsindexer characters-index.json
echo "Sorting records by descending last_name"
dsfind -sort="-last_name" -fields="last_name" -indexes="characters-index.bleve" "*"
echo "Sorting records by ascending last_name"
dsfind -sort="last_name" -fields="last_name" -indexes="characters-index.bleve" "*"
