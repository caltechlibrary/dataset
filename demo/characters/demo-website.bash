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
if [ -d emails.bleve ]; then
    rm -fR emails.bleve
fi
$(dataset init characters)
dataset -uuid import characters.csv
dsindexer characters.json
#dsindexer names.json
#dsindexer emails.json
echo "Open your web browser and go to http://localhost:8011"
dsws -t search.tmpl htdocs characters.bleve # names.bleve emails.bleve

