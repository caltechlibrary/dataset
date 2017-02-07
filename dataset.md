
# USAGE

    dataset [OPTIONS] COMMAND_AND_PARAMETERS

## SYNOPSIS

dataset is a command line tool demonstrating dataset package for managing 
JSON documents stored on disc. A dataset stores one of more collections, 
collections store the a bucketted distribution of JSON documents
as well as metadata about the collection (e.g. collection info,
select lists).

## COMMANDS

Collection and JSON Documant related--

+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable dataset_COLLECTION
    to collection name
+ create - creates a new JSON doc or replace an existing one in collection
  + requires JSON doc name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON doc to stdout
  + requires JSON doc name
+ update - updates a JSON doc in collection
  + requires JSON doc name, followed by replacement JSON blob or 
    JSON blob read from stdin
  + JSON document must already exist
+ delete - removes a JSON doc from collection
  + requires JSON doc name
+ keys - returns the keys to stdout, one key per line
+ path - given a document name return the full path to document

Select list related--

+ select - is the command for working with lists of collection keys
    + "dataset select mylist k1 k2 k3" would create/update a select list 
      mylist adding keys k1, k2, k3
+ lists - returns the select list names associated with a collection
    + "dataset lists"
+ clear - removes a select list from the collection
    + "dataset clear mylist"
+ first - writes the first key to stdout
    + "dataset first mylist"
+ last would display the last key in the list
    + "dataset last mylist"
+ rest displays all but the first key in the list
    + "dataset rest mylist"
+ list displays a list of keys from the select list to stdout
    + "dataet list mylist" 
+ shift writes the first key to stdout and remove it from list
    + "dataset shift mylist" 
+ unshift would insert at the beginning 
    + "dataset unshift mylist k4"
+ push would append the list
    + "dataset push mylist k4"
+ pop removes last key form list and displays it
    + "dataset pop mylist" 
+ sort orders the keys alphabetically in the list
    + "dataset sort mylist asc" - sorts in ascending order
    + "dataset sort mylist desc" - sorts in descending order
+ reverse flips the order of the list
    + "dataset reverse mylists"

## OPTIONS

```
    -c          sets the collection to be used
    -collection sets the collection to be used
    -h          display help
    -help       display help
    -i          input filename
    -input      input filename
    -l          display license
    -license    display license
    -v          display version
    -version    display version
```

## EXAMPLE

This is an example of creating a dataset called testdata/friends, saving
a record called "littlefreda.json" and reading it back.

```
   dataset init testdata/friends
   export DATASET_COLLECTION=testdata/friends
   dataset create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of testdata/friends

```
   dataset -i mojosam.json create mojosam
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

```
   cat capt-jack.json | dataset create capt-jack
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

