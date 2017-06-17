#!/bin/bash
if [ -d characters ]; then
    rm -fR characters
fi
if [ -d characters.bleve ]; then
    rm -fR characters.bleve
fi
$(dataset init characters)
dataset import htdocs/characters.csv
dsindexer htdocs/characters.json characters.bleve
echo "Open your web browser and go to http://localhost:8011"
dsws -dev-mode=true -t htdocs/search.tmpl htdocs characters.bleve
