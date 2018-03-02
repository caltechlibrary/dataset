
# USAGE

	dataset [OPTIONS] [ACTION] [ACTION PARAMETERS...]

## SYNOPSIS


dataset is a command line tool demonstrating dataset package for managing 
JSON documents stored on disc. A dataset is organized around collections,
collections contain buckets holding specific JSON documents and related content.
In addition to the JSON documents dataset maintains metadata for management
of the documents, their attachments as well as a ability to generate select lists
based JSON document keys (aka JSON document names).



## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    DATASET   # Set the working path to your dataset collection
```

## OPTIONS

Options will override any corresponding environment settings. Options are shared between all actions and must precede the action on the command line.

```
    -batch, -size             (indexer, deindexer, find) set the number of records per response
    -c, -collection           sets the collection to be used
    -client-secret            (import-gsheet, export-gsheet) set the client secret path and filename for GSheet access
    -csv                      (find) format results as a CSV document, used with fields option
    -csv-skip-header          (find) don't output a header row, only values for csv output
    -e, -examples             display examples
    -explain                  (find) explain results in a verbose JSON document
    -fields                   (find) comma delimited list of fields to display in the results
    -from                     (find) return the result starting with this result number
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -highlight                (find) display highlight in search results
    -highlighter              (find) set the highlighter (ansi,html) for search results
    -i, -input                input file name
    -ids, -ids-only           (find) output only a list of ids from results
    -json                     (find) format results as a JSON document
    -key-file                 operate on the record keys contained in file, one key per line
    -l, -license              display license
    -nl, -newline             if set to false suppress the trailing newline
    -o, -output               output file name
    -overwrite                overwrite will treat a create as update if the record exists
    -p, -pretty               pretty print output
    -quiet                    suppress error messages
    -sample                   set the sample size when listing keys
    -sort                     (find) a comma delimited list of field names to sort by
    -use-header-row           (import) use the header row as attribute names in the JSON document
    -uuid                     (import) generate a UUID for a new JSON document name
    -v, -version              display version
    -verbose                  output rows processed on importing from CSV
```


## ACTIONS

```
    attach          Attach a document (file) to a JSON record in a collection
    attachments     List of attachments associated with a JSON record in a collection
    check           Check the health of a dataset collection
    count           Counts the number of records in a collection, accepts a filter for sub-counts
    create          Create a JSON record in a collection
    deindexer       Remove record(s) from a Bleve index for a collection
    delete          Delete a JSON record (and attachments) from a collection
    detach          Copy an attach out of an associated JSON record in a collection
    export          Export a JSON records from a collection to a CSV file
    export-gsheet   Export a collection's JSON records to a GSheet
    extract         Extract unique values from JSON records in a collection based on a dot path expression
    find            Query a bleve index(es) associated with a collection
    haskey          Returns true if key is in collection, false otherwise
    import          Import a CSV file's rows as JSON records into a collection
    import-gsheet   Import a GSheet rows as JSON records into a collection
    indexer         Create/Update a Bleve index of a collection
    init            Initialize a dataset collection
    join            Join a JSON record with a new JSON object in a collection
    keys            List the keys in a collection, support filtering and sorting
    list            List the JSON records as an array for provided record ids
    path            Show the file system path to a JSON record in a collection
    prune           Remove attachments from a JSON record in a collection
    read            Read back a JSON record from a collection
    repair          Try to repair a damaged dataset collection
    status          Checks to see if a collection name contains a 'collection.json' file
    update          Update a JSON record in a collection
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
   dataset export titles2016.csv '(eq 2016 (year .pubDate))' \
           '.id,.title,.pubDate' 'id,title,publication date'
```

If wanted to extract a unqie list of all ORCIDs from a collection 

```shell
   dataset extract true .authors[:].orcid
```

If you wanted to extract a list of ORCIDs from publications in 2016.

```shell
   dataset extract '(eq 2016 (year .pubDate))' .authors[:].orcid
```


You can augement JSON key/value pairs for a JSON document in your collection
using the join operation. This works similar to the datatools cli called jsonjoin.

Let's assume you have a record in your collection with a key 'jane.doe'. It has
three fields - name, email, age.

```json
    {"name":"Doe, Jane", "email": "jd@example.org", age: 42}
```

You also have an external JSON document called profile.json. It looks like

```json
    {"name": "Doe, Jane", "email": "jane.doe@example.edu", "bio": "world renowned geophysist"}
```

You can merge the unique fields in profile.json with your existing jane.doe record

```shell
    dataset join append jane.doe profile.json
```

The result would look like

```json
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42, "bio": "renowned geophysist"}
```

If you wanted to overwrite the common fields you would use 'join overwrite'

```shell
    dataset join overwrite jane.doe profile.json
```

Which would result in a record like

```json
    {"name":"Doe, Jane", "email": "jane.doe@example.edu", "age": 42, "bio": "renowned geophysist"}
```



Related: [attach](attach.html), [attachments](attachments.html), [check](check.html), [count](count.html), [create](create.html), [deindexer](deindexer.html), [delete](delete.html), [detach](detach.html), [export](export.html), [export-gsheet](export-gsheet.html), [extract](extract.html), [find](find.html), [haskey](haskey.html), [import](import.html), [import-gsheet](import-gsheet.html), [indexer](indexer.html), [init](init.html), [join](join.html), [keys](keys.html), [list](list.html), [path](path.html), [prune](prune.html), [read](read.html), [repair](repair.html), [status](status.html), [update](update.html)

dataset v0.0.28-dev
