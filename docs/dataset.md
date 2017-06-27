
# dataset

## USAGE

    dataset [OPTIONS] COMMAND_AND_PARAMETERS

## SYNOPSIS

dataset is a command line tool demonstrating dataset package for managing 
JSON documents stored on disc. A dataset is organized around collections,
collections contain buckets holding specific JSON documents and related content.
In addition to the JSON documents dataset maintains metadata for management
of the documents, their attachments as well as a ability to generate select lists
based JSON document keys (aka JSON document names).

## COMMANDS

Collection and JSON Documant related--

+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable DATASET
    to collection name
  + if you're using S3 for storing your dataset prefix your path with 's3://'
    'dataset init s3://mybucket/mydataset-collections'
+ create - creates a new JSON document or replace an existing one in collection
  + requires JSON document name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON document to stdout
  + requires JSON document name
+ update - updates a JSON document in collection
  + requires JSON document name, followed by replacement JSON document name or 
    JSON document read from stdin
  + JSON document must already exist
+ delete - removes a JSON document from collection
  + requires JSON document name
+ filter - takes a filter and returns an unordered list of keys that match filter expression
  + if filter expression not provided as a command line parameter then it is read from stdin
+ keys - returns the keys to stdout, one key per line
+ haskey - returns true is key is in collection, false otherwise
+ path - given a document name return the full path to document
+ attach - attaches a non-JSON content to a JSON record 
    + "dataset attach k1 stats.xlsx" would attach the stats.xlsx file to JSON document named k1
    + (stores content in a related tar file)
+ attachments - lists any attached content for JSON document
    + "dataset attachments k1" would list all the attachments for k1
+ attached - returns attachments for a JSON document 
    + "dataset attached k1" would write out all the attached files for k1
    + "dataset attached k1 stats.xlsx" would write out only the stats.xlsx file attached to k1
+ detach - remove attachments to a JSON document
    + "dataset detach k1 stats.xlsx" would rewrite the attachments tar file without including stats.xlsx
    + "dataset detach k1" would remove ALL attachments to k1
+ import - import a CSV file's rows as JSON documents
	+ "dataset import mydata.csv 1" would import the CSV file mydata.csv using column one's value as key
+ export - export a CSV file based on filtered results of collection records rendering dotpaths associated with column names
	+ "dataset export titles.csv 'true' '._id,.title,.pubDate' 'id,title,publication date'" 
	  this would export all the ids, titles and publication dates as a CSV fiile named titles.csv
+ extract - will return a unique list of unique values based on the associated dot path described in the JSON docs
    + "dataset extract true .authors[:].orcid" would extract a list of authors' orcid ids in collection


## OPTIONS

```
	-c	sets the collection to be used
	-collection	sets the collection to be used
	-h	display help
	-help	display help
	-i	input filename
	-input	input filename
	-l	display license
	-license	display license
	-skip-header-row	skip the header row (use as property names)
	-uuid	generate a UUID for a new JSON document name
	-v	display version
	-verbose	output rows processed on importing from CSV
	-version	display version
```


## EXAMPLES

This is an example of creating a dataset called testdata/friends, saving
a record called "littlefreda.json" and reading it back.

```shell
   dataset init testdata/friends
   export DATASET=testdata/friends
   dataset create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Now check to see if the key, littlefreda, is in the collection

```shell
   dataset haskey littlefreda
```

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of testdata/friends

```shell
   dataset -i mojosam.json create mojosam
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

```shell
   cat capt-jack.json | dataset create capt-jack
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Adding high-capt-jack.txt as an attachment to "capt-jack"

```shell
   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset attach capt-jack high-capt-jack.txt
```

List attachments for "capt-jack"

```shell
   dataset attachments capt-jack
```

Get the attachments for "capt-jack" (this will untar in your current directory)

```shell
   dataset attached capt-jack
```

Remove high-capt-jack.txt from "capt-jack"

```shell
    dataset detach capt-jack high-capt-jack.txt
```

Remove all attachments from "capt-jack"

```shell
   dataset detach capt-jack
```

Filter can be used to return only the record keys that return true for a given
expression. Here's is a simple case for match records where name is equal to
"Mojo Sam".

```shell
   dataset filter '(eq .name "Mojo Sam")'
```

If you are using a complex filter it can read a file in and apply it as a filter.

```shell
   dataset filter < myfilter.txt
```

Import can take a CSV file and store each row as a JSON document in dataset. In
this example we're generating a UUID for the key name of each row

```shell
   dataset -uuid import my-data.csv
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.

```shell
   dataset export titles.csv true '.id,.title,.pubDate' 'id,title,publication date'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset export titles2016.csv '(eq 2016 (year .pubDate))' '.id,.title,.pubDate' 'id,title,publication date'
```

If wanted to extract a unqie list of all ORCIDs from a collection 

```shell
   dataset extract true .authors[:].orcid
```

Finally if you wanted to extract a list of ORCIDs from publications in 2016.

```shell
   dataset extract '(eq 2016 (year .pubDate))' .authors[:].orcid
```



dataset v0.0.3-rc1
