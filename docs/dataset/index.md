
# USAGE

```
    dataset [OPTIONS] COMMAND_AND_PARAMETERS
```


# Description

dataset is a command line tool demonstrating dataset package for managing
JSON documents stored on disc. A dataset is organized around collections,
collections contain buckets holding specific JSON documents and related content.
In addition to the JSON documents dataset maintains metadata for management
of the documents, their attachments as well as a ability to generate select lists
based JSON document keys (aka JSON document names).

## OPTIONS

```
	-c	sets the collection to be used
	-collection	sets the collection to be used
	-example	display example(s)
	-h	display help
	-help	display help
	-i	input filename
	-input	input filename
	-l	display license
	-license	display license
	-no-newline	suppress a trailing newline on output
	-o	output filename
	-output	output filename
	-quiet	suppress error and status output
	-use-header-row	use the header row as attribute names in the JSON document
	-uuid	generate a UUID for a new JSON document name
	-v	display version
	-verbose	output rows processed on importing from CSV
	-version	display version
```

# Help

The following topics are available by specifying the name after the "-help" option.
E.g. to get help on the "init" command.

```
    dataset -help init
```

## Topics

+ [init](init.html) - initialize a new collection if none exists, requires a path to collection
+ [create](create.html) - creates a new JSON document or replace an existing one in collection
+ [read](read.html) - displays a JSON document to stdout
+ [update](update.html) - updates a JSON document in collection
+ [delete](delete.html) - removes a JSON document from collection
+ [join](join.html) - brings the functionality of jsonjoin to the dataset command.
+ [filter](filter.html) - takes a filter and returns an unordered list of keys that match filter expression
    + [dotpath](dotpath.html) - reach into an object to return a value(s)
+ [keys](keys.html) - returns the keys to stdout, one key per line
+ [haskey](haskeys.html) - returns true is key is in collection, false otherwise
+ [count](count.html) - returns a count of keys in a collection
+ [path](path.html) - given a document name return the full path to document
+ [attach](attach.html) - attaches a non-JSON content to a JSON record
+ [attachments](attachments.html) - lists any attached content for JSON document
+ [attached](attached.html) - returns attachments for a JSON document
+ [detach](detach.html) - remove attachments to a JSON document
+ [import](import.html) - import a CSV file's rows as JSON documents
    + [import-gsheet](import-gsheet.html) - import a Google Sheets sheet rows as JSON documents
+ [export](export.html) - export a CSV file based on filtered results of collection records rendering dotpaths associated with column names
    + [export-gsheet](export-gsheet.html) - export a Collection of JSON documents to Google Sheets sheet rows
+ [extract](extract.html) - will return a unique list of unique values based on the associated dot path described in the JSON docs
    + [dotpath](dotpath.html) - reach into an object to return a value(s)
