#!/bin/bash
if [ -d "characters.ds" ]; then
    rm -fR "characters.ds"
fi
if [ -d "characters.bleve" ]; then
    rm -fR "characters.bleve"
fi
$(dataset init "characters.ds")
dataset import "characters.csv"
dsindexer "characters.json" "characters.bleve"
echo "Open your web browser and go to http://localhost:8011"
dsws -dev-mode=true -t "search.tmpl" "." "characters.bleve"
