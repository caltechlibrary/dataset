#!/bin/bash

#
# Filtering Keys using data frames
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
    dataset -pretty read friends.ds "${KEY}"
done

#
# Frames, filter for given name "Mojo"
#

# Step 1. create an unfiltered frame of all records
dataset frame friends.ds "unfiltered" "._Key=id" ".given=given" ".family=family">/dev/null

# Step 2. Filter the unfiltered frame creating a "filtered" data frame
for OBJ in $(dataset frame-objects friends.ds unfiltered | jsonrange -values ); do
    GIVEN=$(echo "${OBJ}" | jsoncols -i - .given | sed -E 's/"//g')
    if [ "${GIVEN}" = "Mojo" ]; then
        echo "${OBJ}" | jsoncols -i - .id | sed -E 's/"//g'
    fi
done |\
# Step 3. create our filtered frame
 dataset -i - frame friends.ds filtered "._Key=id" ".given=given" ".family=family">/dev/null

echo "We now have a frame with only Mojo."
dataset -pretty frame friends.ds filtered
