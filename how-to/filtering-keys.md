
Filters and sorting
===================

__dataset__ supports querying a collection using SQL. In this example the datastore is assumed to be
the default v2.2 SQLite3. The tool use to list of keys filtered by an SQL statement is `dsquery`.


Example
-------

1. Decide what elements you are filter one by looking at an example record (`jq` can be used to pretty the result)
2. Write a `SELECT` SQL statement that can return a single column (i.e. `_Key`), the JSON object fields are expressed in SQL using arrow notation, e.g. `src->>'given' like 'Mojo' or src->>'family' like 'Mojo'`. The column of results needs to be an array of JSON elements, in this case a "string" element hodling the key. We get a quoted string in SQLite3 using `'"' || _Key || '"'`.
3. Using `dsquery` to execute the SQL statement and get back an array of JSON, this can then be processed using `jq` to return a single key one per line.

NOTE: In the example below I've used __jsonrange__ and __jsoncols__ for iterating
and filtering our objects. These are provided by [datatools](https://github.com/caltechlibrary/datatools/releases). See [filtering-keys.bash](filtering-keys.bash).

```shell
#
# dsquery, filter for given name "Mojo"
#

# Step 1. Show some records so I can figure out what part of the JSON object I want.
echo "Look at the Mojo record and see what the fields are I need."
dataset dump friends.ds Mojo | jq .
# Reviewing the records I see I'm iterested in `_Key`, `src->>'given'` and  `src->>'family'`


# Step 2. do our filtering iterating over the unfiltered frame (piping the results)
# This SQL statement I'll want should looke something like this.
cat <<SQL | tee mojo-filter.sql
SELECT '"' || _Key || '"'
FROM friends
WHERE src->>'given' LIKE 'Mojo' 
   OR src->>'family' LIKE 'Mojo'
SQL

# Step 3. Run the SQL query using dsquery, pretty print the output with jq.
echo "Keys for given or family names of 'Mojo'"
dsquery -sql mojo-filter.sql friends.ds | jq -r .[0]
```


