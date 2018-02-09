#!/bin/bash
if [ -d "characters.ds" ]; then
    rm -fR "characters.ds"
fi
if [ -d "characters.bleve" ]; then
    rm -fR "characters.bleve"
fi
if [ -d "names.bleve" ]; then
    rm -fR "names.bleve"
fi
if [ -d "titles.bleve" ]; then
    rm -fR "titles.bleve"
fi
$(dataset init "characters.ds")
echo "TESTING import csv file"
dataset -uuid import "characters.csv"
echo ""
echo "TESTING filter '(eq .name \"Mojo Sam\")'"
dataset filter '(eq .name "Mojo Sam")'
echo ""
echo "TESTING export (eq .name \"Mojo Sam\")"
if [ -f "mojo.csv" ]; then
    rm "mojo.csv"
fi
dataset "export" "mojo.csv" '(eq .name "Mojo Sam")' ".name,.title,.year" "Name,Title,Year"
cat mojo.csv
echo ""
echo "TESTING export true"
if [ -f "all.csv" ]; then
    rm "all.csv"
fi
dataset "export" "all.csv" '(true)' ".name,.title,.year" "Name,Title,Year"
cat all.csv

echo ""
echo "TESTING extract '(eq .name \"Mojo Same\") .year"
dataset extract '(eq .name "Mojo Sam")' ".year"

echo ""
echo "TESTING dsindexer"
dsindexer "characters.json" "characters.bleve"
dsindexer "names.json" "names.bleve"
dsindexer "titles.json" "titles.bleve"
echo ""
echo "TESTING dsfind"
echo "        Find Mojo or Sam"
dsfind "Mojo Sam"
echo "        Find Mojo or Sam"
dsfind -indexes="characters.bleve" -fields="name,title,year" "Mojo Sam"
echo "        Find Mojo"
dsfind -indexes="titles.bleve" -fields="name,title,year" "Mojo"
echo "        Find Mojo or Sam"
dsfind -indexes="names.bleve" -fields="name,title,year" "Mojo Sam"
echo "        Find year is 2002"
dsfind -size=1000 -indexes="titles.bleve" -csv -fields="name,title,year" 'year:2002'
echo "        Find year is 2002"
dsfind -size=1000 -indexes="names.bleve" -csv -fields="name,title,year" 'year:2002'
echo "        Find Mojo or Sam"
dsfind -size=1000 -indexes="names.bleve,titles.bleve" -csv -fields="name,title,year" "Mojo Sam" 
echo "        Find Jack and Flanders"
dsfind -size=1000 -sort='year' -indexes="characters.bleve" -csv -fields="title,year" '+name,"Jack" +name,"Flanders"'
echo "        Find Jack and Flanders"
dsfind -size=1000 -sort='title' -indexes="characters.bleve" -csv -fields="title,year" '+name,"Jack" +name,"Flanders"'
echo "        Find Mojo or Frieda"
dsfind -size=1000 -sort='name' -indexes="characters.bleve" -csv -fields="name,title,year" "Mojo Frieda"
echo "Tests completed"
