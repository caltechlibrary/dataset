#!/bin/bash

if [ ! -d "data.ds" ]; then
    echo "missing data collection for demo"
    exit 1
fi

# Index collection
if [ -f "index.bleve" ]; then
    rm -fR "index.bleve"
fi

echo "Indexing creating family_name, given_name, display_name fields via templates"
dataset indexer data.ds idxdefn.json index.bleve

# Show CSV output for indexes records
dataset find -csv -size 100 -sort "orcid" -fields "orcid,family_name,given_name,display_name" "index.bleve" "*"


