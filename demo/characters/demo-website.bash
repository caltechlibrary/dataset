#!/bin/bash
if [ -d characters ]; then
    rm -fR characters
fi
if [ -d characters.bleve ]; then
    rm -fR characters.bleve
fi
if [ -d names.bleve ]; then
    rm -fR names.bleve
fi
if [ -d stories.bleve ]; then
    rm -fR stories.bleve
fi
$(dataset init characters)
dataset -uuid import characters.csv
dsindexer characters.json
dsindexer stories.json
dsindexer names.json
echo "Open your web browser and go to http://localhost:8011"
dsws -t search.tmpl htdocs characters.bleve # names.bleve stories.bleve

