
## _dataset_ Operations

The basic operations support by *dataset* are listed below organized by collection and JSON document level.

### Collection Level

+ [init](dataset/init.html) creates a collection
+ [import](dataset/import.html) JSON documents from rows of a CSV file
+ [import-gsheet](dataset/import.html) JSON documents from rows of a Google Sheet
+ [export](dataset/export.html) JSON documents from a collection into a CSV file
+ [export-gsheet](dataset/export-gsheet.html) JSON documents from a collection into a Google Sheet
+ [keys](dataset/keys.html) list keys of JSON documents in a collection, supports filtering and sorting
+ [haskey](dataset/haskey.html) returns true if key is found in collection, false otherwise
+ [count](dataset/count.html) returns the number of documents in a collection, supports filtering for subsets
+ [extract](dataset/extract.html) unique JSON attribute values from a collection

### JSON Document level

+ [create](dataset/create.html) a JSON document in a collection
+ [read](dataset/read.html) back a JSON document in a collection
+ [update](dataset/update.html) a JSON document in a collection
+ [delete](dataset/delete.html) a JSON document in a collection
+ [join](dataset/join.html) a JSON document with a document in a collection
+ [list](dataset/list.html) the lists JSON records as an array for the supplied keys
+ [path](dataset/path.html) list the file path for a JSON document in a collection

### JSON Document Attachments

+ [attach](dataset/attach.html) a file to a JSON document in a collection
+ [attachments](dataset/attachments.html) lists the files attached to a JSON document in a collection
+ [detach](dataset/detach.html) retrieve an attached file associated with a JSON document in a collection
+ [prune](dataset/prune.html) delete one or more attached files of a JSON document in a collection

### Search

+ [indexer](dataset/indexer.html) indexes JSON documents in a collection for searching with _find_
+ [deindexer](dataset/deindexer.html) de-indexes (removes) JSON documents from an index
+ [find](dataset/find.html) provides a search indexed full text interface into a collection


## Example

Common operations using the *dataset* command line tool

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document
+ import a CSV file as JSON documents
+ how to remove a *dataset* collection

```shell
    # Create a collection "mystuff" inside the directory called demo
    dataset init demo/mystuff
    # if successful an expression to export the collection name is show
    export DATASET=demo/mystuff

    # Create a JSON document 
    dataset create freda.json '{"name":"freda","email":"freda@inverness.example.org"}'
    # If successful then you should see an OK or an error message

    # Read a JSON document
    dataset read freda.json

    # Path to JSON document
    dataset path freda.json

    # Update a JSON document
    dataset update freda.json '{"name":"freda","email":"freda@zbs.example.org"}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys

    # Delete a JSON document
    dataset delete freda.json

    # Import CSV file as JSON documents using column 1 as JSON document name
    # (if no column given the row number will be used for the JSON document name)
    dataset import my-data.csv 1

    # To remove the collection just use the Unix shell command
    # /bin/rm -fR demo/mystuff
```

