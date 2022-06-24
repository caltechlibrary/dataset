
Filters and sorting
===================

__dataset__ does not support internally sorting or filtering of keys.
It does support data frames which can be used to do that via external
tools.

Example
-------

1. Create a data frame using `frame` verb containing the field ".given" and the record ".id"
2. Iterate over the frame objects and in the frame using `frame-objects` verb.
3. For desired keys output the key and send to a new "filtered" frame using `frame` verb.

NOTE: In the example below I've used __jsonrange__ and __jsoncols__ for iterating
and filtering our objects. These are provided by [datatools](https://github.com/caltechlibrary/datatools/releases). See [filtering-keys.bash](filtering-keys.bash).

```shell
#
# Frames, filter for given name "Mojo"
#

# Step 1.
dataset frame friends.ds "unfiltered" "._Key=id" ".given=given" ".family=family">/dev/null

# Step 2. do our filtering iterating over the unfiltered frame (piping the results)
for OBJ in $(dataset frame-objects friends.ds unfiltered | jsonrange -values ); do
    GIVEN=$(echo "${OBJ}" | jsoncols -i - .given | sed -E 's/"//g')
    # This is the filter, we're checking if the record is about Mojo.
    if [ "${GIVEN}" = "Mojo" ]; then
        echo "${OBJ}" | jsoncols -i - .id | sed -E 's/"//g'
    fi
done |\
# Step 3. create a filtered frame
  dataset -i - frame friends.ds filtered "._Key=id" ".given=given" ".family=family">/dev/null

echo "We now have a frame with only Mojo."
dataset -pretty frame friends.ds filtered
```


