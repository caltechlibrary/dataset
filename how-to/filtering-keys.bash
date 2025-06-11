#!/bin/bash

#
# Filtering Keys using dsquery
#


#
# create a collections with init
#
if [ -d "friends.ds" ]; then
    rm -fR friends.ds
fi
dataset init friends.ds

#
# create some records
#
dataset create friends.ds Frieda  \
  '{"display_name":"Little Frieda","given": "Frieda", "family": "Little"}'
dataset create friends.ds Mojo \
  '{"display_name": "Mojo Sam, the Yudoo Man", "given": "Mojo", "family": "Sam"}'
dataset create friends.ds Jack \
  '{"display_name": "Jack Flanders", "given": "Jack", "family": "Flanders"}'
dataset create friends.ds Mozulla \
  '{"display_name": "Professor Mozulla", "given": "Marvin", "family": "Mozulla"}'

## Display the records we created.
for KEY in Frieda Mojo Jack Mozulla; do
    echo "Reading ${KEY} object"
    dataset read friends.ds "${KEY}"
done

#
# Filter using dsquery for given name "Mojo"
#

# Step 1. Show some records so I can figure out what part of the JSON object I want.
echo "Look at the Mojo record and see what the fields are I need."
dataset dump friends.ds Mojo | jq .
# Reviewing the records I see I'm iterested in `_Key`, `src->>'given'` and  `src->>'family'`


# Step 2. do our filtering iterating over the unfiltered frame (piping the results)
# This SQL statement I'll want should looke something like this.
cat <<SQL >mojo-filter.sql
SELECT '"' || _Key || '"'
FROM friends
WHERE src->>'given' LIKE 'Mojo' 
   OR src->>'family' LIKE 'Mojo'
SQL

# Step 3. Run the SQL query using dsquery, pretty print the output with jq.
echo "Keys for given or family names of 'Mojo'"
dsquery -sql mojo-filter.sql friends.ds | jq -r .[0]

