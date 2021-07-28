#!/bin/bash

#
# Getting Started using dataset and Stephen Dolan's jq.
#


#
# create a collections with init
#
if [ -d "friends.ds" ]; then
    rm -fR friends.ds
fi
dataset init friends.ds


#
# create, read, update and delete
#

## create
dataset create friends.ds  Frieda  '{"name":"Little Frieda","email":"frieda@inverness.example.org"}'
dataset create friends.ds  Mojo '{"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"}'
dataset create friends.ds Jack '{"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"}'
dataset create friends.ds Mazulla '{"name": "Professor Mazulla", "email": "mm@alchemist.example.org"}'

## read
for KEY in Frieda Mojo Jack; do
    echo "Reading ${KEY} profile"
    dataset read friends.ds "${KEY}" | jq .
done

## Add a "catch_phrase" to existing records.
function add_field() {
    KEY="${1}"
    FIELD="${2}"
    VALUE="${3}"
    # form the field into a key/value pair as JSON
    SRC="{\"${FIELD}\": \"${VALUE}\"}"
    dataset join friends.ds "${KEY}" "${SRC}"
}

add_field Frieda catch_phrase "Wowee Zowee"
add_field Mojo catch_phrase "Feet Don't Fail Me Now!"
add_field Jack catch_phrase "What is coming at you is coming from you"
add_field Frieda given "Frieda"
add_field Mojo given "Mojo"
add_field Jack given "Jack"
add_fields Mazulla given "Marvin"
add_field Frieda family "Little"
add_field Mojo family "Sam"
add_field Jack family "Flanders"
add_fields Mazulla family "Mazulla"

# Display our ammeded records.
for KEY in Frieda Mojo Jack; do
    echo "Reading ${KEY} profile"
    dataset read friends.ds "${KEY}" | jq .
done

# Updating (replacing) Frieda's record with new email address using sed
dataset read friends.ds Frieda |\
jq . | \
sed -E 's/"email":"frieda@inverness.example.org"/"email":"frieda@venus.example.org"/' \
| dataset update friends.ds Frieda 

## delete example, remove Mazulla
dataset delete friends.ds Mazulla
dataset keys friends.ds


#
# Keys and count
#

# count
cnt=$(dataset count friends.ds)
echo "Total Records now: ${cnt}"

# keys
for KEY in $(dataset keys friends.ds); do
    if [ "${KEY}" != "" ]; then
        echo "${KEY}"
    fi
done

#
# Frames, filter for given name "Mojo"
#
dataset frame friends.ds "unfiltered" "._Key=id" ".given=given" ".family=family">/dev/null
#dataset frame friends.ds unfiltered | jq .
for OBJ in $(dataset frame-objects friends.ds unfiltered | jsonrange -values ); do
    #echo -n "DEBUG "; echo "${OBJ}" | jq .given 
    GIVEN=$(echo "${OBJ}" | jq .given | sed -E 's/"//g')
    #echo "DEBUG given -> ${GIVEN}"
    if [ "${GIVEN}" = "Mojo" ]; then
        echo "${OBJ}" | jq .id | sed -E 's/"//g'
    fi
done |dataset -i - frame friends.ds filtered "._Key=id" ".given=given" ".family=family" | jq .

echo "We now have a frame with only Mojo."
