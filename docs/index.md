
# Documentation for dataset

The documentation generallyis organized by command line programs
but a more explority approach can be taken by the list below

+ [getting started with dataset](getting-started-with-dataset.html)
+ [how to](../how-to/) -- task oriented

## Command line program documentation

+ [dataset](dataset/) - the command line tool for managing _dataset_ collections
+ [dsws](dsws/) - A web server and web service based on one or more Bleve indexes created with _dataset indexer_


## Internal project concepts

+ [file system layout](file-system-layout.html) - Describes how collections are organized
+ [defining indexes](defining-indexes.html) - Describes the index definition JSON document format
+ [cloud storage](cloud-storage.html) - Describes using Cloud Storage (e.g. Amazon S3, Google Cloud Storage)
+ [Google Spreadsheet integration](gsheet-integration.html) - describes how to setup import/export access to a Google Spreadsheet 

## _dataset_ Operations

The basic operations support by *dataset* are listed below organized by collection and JSON document level.

### Collection Level

+ [init](dataset/init.html) creates a collection
+ [import-csv](dataset/import-csv.html) JSON documents from rows of a CSV file
+ [import-gsheet](dataset/import.html) JSON documents from rows of a Google Sheet
+ [export-csv](dataset/export-csv.html) JSON documents from a collection into a CSV file
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













