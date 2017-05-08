#!/bin/bash

#
# Test the cli utilities that demonstrate features of dataset package.
#
if [ -d "testdata/characters-index.bleve" ]; then
    rm -fR "testdata/characters-index.bleve" 
fi
if [ -d "testdata/characters" ]; then
    rm -fR testdata/characters
fi
if [ -f "testdata/characters.csv" ]; then
    rm testdata/characters.csv
fi
if [ -f "testdata/characters-index.json" ]; then
    rm testdata/characters-index.json
fi

#
# Pick version of cli to test with
if [ -f "bin/dataset" ] && [ -f "bin/dsfind" ] && [ -f "bin/dsindexer" ]; then
    export PATH="./bin":$PATH
fi

# Generate CSV test data
cat<<FILE1 > testdata/characters.csv
last_name,first_name,email
sam,mojo,mojo.sam@zbs.example.org
frieda,little,little.frieda@zbs.example.org
flanders,Jack,captain.jack@zbs@example.org
art,far seeing,old.far.seeing.art@zbs.example.org
shoe,ruby,ruby2@zbs.example.org
turu,t.j.,t.j.turu@zbs.example.org
the andover,rhodes,arhodes@another.example.org
kapur,rodant,rodant.kapur@zbs.example.org
FILE1

# Generate the index mapping
cat<<FILE2 > testdata/characters-index.json
{
    "last_name":{
        "object_path":".last_name"
    },
    "first_name":{
        "object_path":".first_name"
    },
    "email":{
        "object_path":".email"
    }
}
FILE2

# Initialize an empty repository
$(dataset init testdata/characters)
# Load the data
dataset -uuid import testdata/characters.csv
dsindexer testdata/characters-index.json
echo "Sorting records by descending last_name"
dsfind -sort="-last_name" -fields="last_name" -indexes="testdata/characters-index.bleve" "*"
echo "Sorting records by ascending last_name"
dsfind -sort="last_name" -fields="last_name" -indexes="testdata/characters-index.bleve" "*"
