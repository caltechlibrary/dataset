#!/bin/bash

# Setup our collection
dataset init JoinDemo.ds
export DATASET="JoinDemo.ds"

cat <<EOF

join will allow you to merge (updating or overwriting) an existing JSON document stored in
a collection. Two JOIN_TYPES are available -- "update" and "overwrite".  With "update" 
only new fields will be added to the record. If you specify "overwrite" new fields will be 
added and existing fields in common will be overwritten.

join us helpful in building up an aggregated record where you have a common JSON_RECORD_ID.

EOF

# Create a Jane Doe profile record from profile.json
echo "Creating a record called Jane.Doe from person.json"
dataset -i person.json create Jane.Doe 
echo "Reading it back..."
dataset read Jane.Doe

# read it back
echo "Join profile.json using 'join update'..."
dataset join update Jane.Doe profile.json
echo "Reading it back..."
dataset read Jane.Doe
echo "Join profile.json using 'join overwrite'..."
dataset join overwrite Jane.Doe profile.json
dataset read Jane.Doe
