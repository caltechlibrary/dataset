#!/bin/bash
if [ -d authors ]; then
    rm -fR authors
fi
if [ -d authors.bleve ]; then
    rm -fR authors.bleve
fi

echo "Run dataset command and create our collection from our CaltechAUTHORS sample"
read -p "Press any key to run command, ctrl-C to exit" NEXT
$(dataset init authors)
for ITEM in $(ls data/*.json); do
    ID=$(jsoncols -i "${ITEM}" .id)
    dataset -i "${ITEM}" create "${ID}";
    echo "ID: ${ID}, Item: ${ITEM}"
done

echo ""
echo "Run dsindexer to index our collection based on our definition in authors.json"
read -p "Press any key to run command, ctrl-C to exit" NEXT
dsindexer authors.json

echo "Run dsfind to generate a CSV table from id, title, authors_id, orcid searching for Singh-C"
read -p "Press any key to run command, ctrl-C to exit" NEXT
dsfind -c authors -fields="id,title,authors_id,orcid" -csv 'Singh-C'


echo ""
echo "Run dsws for a web searchable version of our collection"
read -p "Press any key to run command, ctrl-C to exit" NEXT
echo "Open your web browser and go to http://localhost:8011"
dsws -dev-mode=true -t templates authors.bleve


