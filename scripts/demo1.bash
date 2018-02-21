#!/bin/bash

if [ ! -f demo/friends ]; then
	echo "create collections demo/friends.ds"
	dataset init demo/friends.ds
fi
echo "Save the collection name in the environment"
export DATASET="demo/friends.ds"
echo "Add some records ..."
for NAME in "Captain Jack Flanders" "Little Frieda" "Mojo Sam" "Dominique" "Kamela" "Ruby" "Angel Sisters"; do
	echo "Add '$NAME'"
	dataset create "${NAME// /-}" "{\"name\":\"$NAME\"}"
done

TMP=$(dataset list keys)
dataset select names "$TMP"
