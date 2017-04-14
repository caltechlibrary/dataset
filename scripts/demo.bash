#!/bin/bash

if [ "$DATASET" != "" ]; then
	OLD_DATASET=$DATASET
fi

# Create a collection "mystuff" inside the directory called demo
dataset init demo/mystuff
# if successful an expression to export the collection name is show
export DATASET=demo/mystuff

# Create a JSON document 
dataset create freda.json '{"name":"freda","email":"freda@inverness.example.org"}'
# If successful then you should see an OK or an error message

# Read a JSON document
dataset read freda.json

# Path to JSON document
dataset path freda.json

# Update a JSON document
dataset update freda.json '{"name":"freda","email":"freda@zbs.example.org"}'
# If successful then you should see an OK or an error message

# List the keys in the collection
dataset keys

# Delete a JSON document
dataset delete freda.json

# To remove the collection just use the Unix shell command
# /bin/rm -fR demo/mystuff

if [ "$OLD_DATASET" != "" ]; then
	export DATASET=$OLD_DATASET
else
	unset DATASET
fi
