#!/bin/bash

# Create our collection if needed
if [ ! -f testdata/fiends/collection.json ]; then
	echo "Creating testdata/friends.ds"
	dataset init testdata/friends.ds
fi
export DATASET=testdata/friends.ds
echo "Creating document 'littlefreda.json'"
dataset create littlefreda.json '{"name":"Freda","email":"little.freda@inverness.example.org"}'
for KY in $(dataset keys); do
	echo "Path: $(dataset path "$KY")"
	echo "Doc: $(dataset read "$KY")"
done
#dataset delete littlefreda.json
