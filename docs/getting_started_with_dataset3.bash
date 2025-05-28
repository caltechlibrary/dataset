#!/bin/bash

#
# Getting Started using dataset and Stephen Dolan's jq and dataset3 query function.
#


#
# create a collections with init
#
if [ -d "friends.ds" ]; then
    rm -fR friends.ds
fi
dataset3 init friends.ds


#
# create, read, update and delete
#

## create
dataset3 create friends.ds Frieda  '{"name":"Little Frieda","email":"frieda@inverness.example.org"}'
dataset3 create friends.ds Mojo '{"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"}'
dataset3 create friends.ds Jack '{"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"}'
dataset3 create friends.ds Mazulla '{"name": "Professor Mazulla", "email": "mm@alchemist.example.org"}'

## read
for KEY in Frieda Mojo Jack; do
    echo "Reading ${KEY} profile"
    dataset read friends.ds "${KEY}" | jq .
done

## Add a "catch_phrase", "given" and "family" to existing records.
function add_field() {
    KEY="${1}"
    FIELD="${2}"
    VALUE="${3}"
    # Get original object as one line
    OBJ="$(dataset3 read friends.ds "${KEY}")"
    # form the field into a key/value pair as JSON
    cat <<JSON_SRC | jq --slurp '.[0] * .[1]' | dataset3 update friends.ds "${KEY}"
${OBJ}
{"${FIELD}": "${VALUE}"}
JSON_SRC

}

add_field Frieda catch_phrase "Wowee Zowee"
add_field Mojo catch_phrase "Feet Don't Fail Me Now!"
add_field Jack catch_phrase "What is coming at you is coming from you"
add_field Frieda given "Frieda"
add_field Mojo given "Mojo"
add_field Jack given "Jack"
add_field Mazulla given "Marvin"
add_field Frieda family "Little"
add_field Mojo family "Sam"
add_field Jack family "Flanders"
add_field Mazulla family "Mazulla"

# Display our ammeded records.
for KEY in Frieda Mojo Jack; do
    echo "Reading ${KEY} profile"
    dataset3 read friends.ds "${KEY}" | jq .
done

# Updating (replacing) Frieda's record with new email address using sed
dataset3 read friends.ds Frieda |\
jq . | \
sed -E 's/"email":"frieda@inverness.example.org"/"email":"frieda@venus.example.org"/' \
| dataset3 update friends.ds Frieda 

## delete example, remove Mazulla
dataset3 delete friends.ds Mazulla
dataset3 keys friends.ds


#
# Keys and counting
#

# List keys
dataset3 keys friends.ds

# count can be done using dataset3 query function combined with jq.
cnt=$(dataset3 query friends.ds 'select count(*) from friends' | jq -r '.[0]')
echo "Total Records now: ${cnt}"


#
# Filter fiends.ds for the name "Mojo", save a Mojo.json.
#
cat <<SQL | dataset3 query friends.ds >Mojo.json
select src
from friends
where src->>'given' = 'Mojo'
order by _Key;
SQL

echo "We now have a JSON file called Mojo.json"
