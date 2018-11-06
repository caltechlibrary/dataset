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
dataset import characters.ds "characters.csv" 1

echo ""
echo "TESTING filtering with keys and grid"
dataset keys characters.ds '(eq .name "Mojo Sam")' > mojo-sam.keys

echo "Now we use a grid, and two datatools utilities (jsonrange and jsoncols) to display I selected results"
dataset grid characters.ds mojo-sam.keys "._Key" ".name" ".title" ."year" | jsonrange -i - --values | while read LINE; do
   echo "${LINE}" | jsoncols -nl=true  -csv -i - .[0] .[1] .[3]
done
 
#FIXME: export needs a frame ...
echo ""
echo "TESTING export csv (eq .name \"Mojo Sam\")"
if [ -f "mojo.csv" ]; then
    rm "mojo.csv"
fi
dataset export characters.ds "mojo.csv" '(eq .name "Mojo Sam")' "._Key,.name,.title,.year" "ID,Name,Title,Year"
cat mojo.csv
echo ""
echo "TESTING export csv true"
if [ -f "all.csv" ]; then
    rm "all.csv"
fi
dataset export characters.ds "all.csv" '(true)' "._Key,.name,.title,.year" "ID,Name,Title,Year"
cat all.csv

echo ""
echo "TESTING finding Mojo Sam and years"
dataset keys characters.ds '(eq .name "Mojo Sam")' > mojo-sam.keys
dataset grid characters.ds mojo-sam.keys .year | jsonrange -i - --values | while read LINE; do
    echo $LINE | jsoncols -i - -nl=true  .[0]
done | sort 

echo ""
echo "TESTING dataset indexer"
dataset indexer characters.ds characters.bmap characters.bleve 
if [[ "$?" = "1" ]]; then
    echo "Can't create characters.bleve, aborting"
    exit 1
fi
dataset indexer characters.ds names.json names.bleve 
dataset indexer characters.ds titles.json titles.bleve 
echo ""
echo "TESTING dataset find"
echo "        Find Mojo or Sam"
dataset find characters.bleve "Mojo Sam"
echo "        Find Mojo or Sam"
dataset find -fields="name,title,year" characters.bleve "Mojo Sam"
echo "        Find Mojo"
dataset find -fields="name,title,year" titles.bleve "Mojo"
echo "        Find Mojo or Sam"
dataset find -fields="name,title,year" names.bleve "Mojo Sam"
echo "        Find year is 2002"
dataset find -size=1000 -csv -fields="name,title,year" characters.bleve 'year:2002'
echo "        Find year is 2002"
dataset find -size=1000 -csv -fields="name,title,year" names.bleve 'year:2002'
echo "        Find Mojo or Sam"
dataset find -size=1000 -indexes="names.bleve,titles.bleve" -csv -fields="name,title,year" "Mojo Sam" 
echo "        Find Jack and Flanders"
dataset find -size=1000 -sort='year' -csv -fields="title,year" \
    characters.bleve '+name,"Jack" +name,"Flanders"'
echo "        Find Jack and Flanders"
dataset find -size=1000 -sort='title' -csv -fields="title,year" \
    characters.bleve '+name,"Jack" +name,"Flanders"'
echo "        Find Mojo or Frieda"
dataset find -size=1000 -sort='name' -csv -fields="name,title,year" \
    characters.bleve "Mojo Frieda"
echo "Tests completed"
