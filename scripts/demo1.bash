#!/bin/bash

echo "create collections demo/friends"
datatset init demo/friends
echo "Save the collection name in the environment"
export DATASET_COLLECTION=demo/friends
echo "Add some records ..."
for Name in "Captain Jack Flanders" "Little Frieda" "Mojo Sam" "Dominique" "Kamela" "Ruby" "Angel Sisters"; do
  dataset create $(slugify $Name) '{"name":"'$Name'"}'
done
