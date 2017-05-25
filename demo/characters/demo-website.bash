#!/bin/bash
if [ -d characters ]; then
    rm -fR characters
fi
if [ -d characters.bleve ]; then
    rm -fR characters.bleve
fi
$(dataset init characters)
dataset import characters.csv
dsindexer characters.json
echo "Open your web browser and go to http://localhost:8011"
dsws -t search.tmpl htdocs characters.bleve
