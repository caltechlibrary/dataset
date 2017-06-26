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
if [ -d titles.bleve ]; then
    rm -fR titles.bleve
fi
$(dataset init characters)
echo "Testing import csv file"
dataset -uuid import htdocs/characters.csv
echo "Testing filter '(eq .name "Mojo Sam")'"
dataset filter '(eq .name "Mojo Sam")'
echo "Testing dsindexer"
dsindexer htdocs/characters.json characters.bleve
dsindexer htdocs/names.json names.bleve
dsindexer htdocs/titles.json titles.bleve
echo "Testing dsfind"
dsfind "Mojo Sam"
dsfind -indexes=characters.bleve -fields="name,title,year" "Mojo Sam"
dsfind -indexes=titles.bleve -fields="name,title,year" "Mojo"
dsfind -indexes=names.bleve -fields="name,title,year" "Mojo Sam"
dsfind -size=1000 -indexes=titles.bleve -csv -fields="name,title,year" 'year:2002'
dsfind -size=1000 -indexes=names.bleve -csv -fields="name,title,year" 'year:2002'
dsfind -size=1000 -indexes="names.bleve,titles.bleve" -csv -fields="name,title,year" "Mojo Sam" 
dsfind -size=1000 -sort='year' -indexes=characters.bleve -csv -fields="title,year" '+name,"Jack" +name,"Flanders"'
dsfind -size=1000 -sort='title' -indexes=characters.bleve -csv -fields="title,year" '+name,"Jack" +name,"Flanders"'
dsfind -size=1000 -sort='name' -indexes=characters.bleve -csv -fields="name,title,year" "Mojo Frieda"
