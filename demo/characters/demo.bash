#!/bin/bash
if [ -d characters ]; then
    rm -fR characters
fi
$(dataset init characters)
dataset -uuid import characters.csv
dsindexer characters.json
dsindexer names.json
dsindexer emails.json
dsfind "Mojo Sam"
dsfind -indexes=characters.bleve "Mojo Sam"
dsfind -indexes=emails.bleve "Mojo Sam"
dsfind -indexes=emails.bleve "mojo.sam"
dsfind -indexes=emails.bleve "zbs.example.org"
dsfind -indexes=names.bleve "Mojo Sam"
dsfind -indexes=names.bleve:emails.bleve "Mojo Sam" 
dsfind -sort='-name'  -indexes=characters.bleve:names.bleve:emails.bleve "email:zbs.example.org"
