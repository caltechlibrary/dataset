#!/bin/bash
if [[ -d "characters.ds" ]]; then
    rm -fR "characters.ds"
fi
if [[ -d "characters.bleve" ]]; then
    rm -fR "characters.bleve"
fi
if [[ -d "names.bleve" ]]; then
    rm -fR "names.bleve"
fi
if [[ -d "titles.bleve" ]]; then
    rm -fR "titles.bleve"
fi
if [[ -f "mojo-sam.keys" ]]; then
    rm mojo-sam.keys
fi
dataset init "characters.ds"
echo "TESTING import csv file"
dataset characters.ds import-csv "characters.csv" 1

echo ""
echo "TESTING filtering with keys and grid"
dataset characters.ds keys '(eq .name "Mojo Sam")' > mojo-sam.keys

echo "Now we use a grid, and two datatools utilities (jsonrange and jsoncols) to display I selected results"
dataset characters.ds grid mojo-sam.keys "._Key" ".name" ".title" ."year" | jsonrange -i - --values | while read LINE; do
   echo "${LINE}" | jsoncols -nl=true  -csv -i - .[0] .[1] .[3]
done
 
echo ""
echo "TESTING export-csv (eq .name \"Mojo Sam\")"
if [ -f "mojo.csv" ]; then
    rm "mojo.csv"
fi
dataset characters.ds export-csv "mojo.csv" '(eq .name "Mojo Sam")' "._Key,.name,.title,.year" "ID,Name,Title,Year"
cat mojo.csv
echo ""
echo "TESTING export-csv true"
if [ -f "all.csv" ]; then
    rm "all.csv"
fi
dataset characters.ds export-csv "all.csv" '(true)' "._Key,.name,.title,.year" "ID,Name,Title,Year"
cat all.csv

echo ""
echo "TESTING finding Mojo Sam and years"
dataset characters.ds keys '(eq .name "Mojo Sam")' > mojo-sam.keys
dataset characters.ds grid mojo-sam.keys .year | jsonrange -i - --values | while read LINE; do
    echo $LINE | jsoncols -i - -nl=true  .[0]
done | sort 

echo ""
echo "TESTING dataset indexer"
dataset characters.ds indexer characters.bleve characters.json
if [[ "$?" = "1" ]]; then
    echo "Can't create characters.bleve, aborting"
    exit 1
fi
dataset characters.ds indexer names.bleve names.json 
dataset characters.ds indexer titles.bleve titles.json
echo ""
echo "TESTING dataset find"
echo "        Find Mojo or Sam"
dataset find characters.bleve "Mojo Sam"
echo "        Find Mojo or Sam"
dataset -fields="name,title,year" find characters.bleve "Mojo Sam"
echo "        Find Mojo"
dataset -fields="name,title,year" find titles.bleve "Mojo"
echo "        Find Mojo or Sam"
dataset -fields="name,title,year" find names.bleve "Mojo Sam"
echo "        Find year is 2002"
dataset -size=1000 -csv -fields="name,title,year" find characters.bleve 'year:2002'
echo "        Find year is 2002"
dataset -size=1000 -csv -fields="name,title,year" find names.bleve 'year:2002'
echo "        Find Mojo or Sam"
dataset -size=1000 -indexes="names.bleve,titles.bleve" -csv -fields="name,title,year" find "Mojo Sam" 
echo "        Find Jack and Flanders"
dataset -size=1000 -sort='year' -csv -fields="title,year" \
    find characters.bleve '+name,"Jack" +name,"Flanders"'
echo "        Find Jack and Flanders"
dataset -size=1000 -sort='title' -csv -fields="title,year" \
    find characters.bleve '+name,"Jack" +name,"Flanders"'
echo "        Find Mojo or Frieda"
dataset -size=1000 -sort='name' -csv -fields="name,title,year" \
    find characters.bleve "Mojo Frieda"
echo "Tests completed"
