#!/bin/bash
if [ -d authors ]; then
    rm -fR authors
fi
if [ -d authors.bleve ]; then
    rm -fR authors.bleve
fi
$(dataset init authors)
for ITEM in $(ls data); do
    ID=$(jsoncols -i "data/${ITEM}" .id)
    dataset -i "data/${ITEM}" create "${ID}";
    echo "ID: ${ID}, Item: ${ITEM}"
done

echo "Indexing repository"
dsindexer authors.json

echo "Example CSV output"
dsfind -c authors -fields="id,title,authors_id,orcid" -csv 'Singh-C'


echo "Open your web browser and go to http://localhost:8011"
dsws -dev-mode -t search.tmpl htdocs authors.bleve


