
# USAGE

	dataset [OPTIONS] COLLECTION_NAME ACTION [ACTION PARAMETERS...]

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
    attach         Attach a document (file) to a JSON record in a collection
    attachments    List of attachments associated with a JSON record in a collection
    check          Check the health of a dataset collection
    clone          Clone a collection from a list of keys into a new collection
    clone-sample   Clone a collection into a sample size based training collection and test collection
    count          Counts the number of records in a collection, accepts a filter for sub-counts
    create         Create a JSON record in a collection
    deindexer      Remove record(s) from a Bleve index for a collection
    delete         Delete a JSON record (and attachments) from a collection
    detach         Copy an attach out of an associated JSON record in a collection
    export         Export a JSON records from a collection to a CSV file
    export-gsheet  Export a collection's JSON records to a GSheet
    extract        Extract unique values from JSON records in a collection based on a dot path expression
    find           Query a bleve index(es) associated with a collection
    haskey         Returns true if key is in collection, false otherwise
    import         Import a CSV file's rows as JSON records into a collection
    import-gsheet  Import a GSheet rows as JSON records into a collection
    indexer        Create/Update a Bleve index of a collection
    init           Initialize a dataset collection
    join           Join a JSON record with a new JSON object in a collection
    keys           List the keys in a collection, support filtering and sorting
    list           List the JSON records as an array for provided record ids
    path           Show the file system path to a JSON record in a collection
    prune          Remove attachments from a JSON record in a collection
    read           Read back a JSON record from a collection
    repair         Try to repair a damaged dataset collection
    status         Checks to see if a collection name contains a 'collection.json' file
    update         Update a JSON record in a collection
```


Related: [attach](attach.html), [attachments](attachments.html), [check](check.html), [clone](clone.html), [clone-sample](clone-sample.html), [count](count.html), [create](create.html), [deindexer](deindexer.html), [delete](delete.html), [detach](detach.html), [export](export.html), [export-gsheet](export-gsheet.html), [extract](extract.html), [find](find.html), [haskey](haskey.html), [import](import.html), [import-gsheet](import-gsheet.html), [indexer](indexer.html), [init](init.html), [join](join.html), [keys](keys.html), [list](list.html), [path](path.html), [prune](prune.html), [read](read.html), [repair](repair.html), [status](status.html), [update](update.html)

dataset v0.0.39
