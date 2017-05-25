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
dsindexer names.json
dsindexer stories.json
dsfind "Mojo Sam"
dsfind -indexes=characters.bleve -fields="name:story:year" "Mojo Sam"
dsfind -indexes=stories.bleve -fields="name:story:year" "Mojo"
dsfind -indexes=names.bleve -fields="name:story:year" "Mojo Sam"
dsfind -size=1000 -indexes=stories.bleve -csv -fields="name:story:year" 'year:2002'
dsfind -size=1000 -indexes=names.bleve -csv -fields="name:story:year" 'year:2002'
dsfind -size=1000 -indexes="names.bleve:stories.bleve" -csv -fields="name:story:year" "Mojo Sam" 
dsfind -size=1000 -sort='year' -indexes=characters.bleve -csv -fields="story:year" '+name:"Jack" +name:"Flanders"'
dsfind -size=1000 -sort='story' -indexes=characters.bleve -csv -fields="story:year" '+name:"Jack" +name:"Flanders"'
dsfind -size=1000 -sort='name' -indexes=characters.bleve -csv -fields="name:story:year" "Mojo Frieda"
