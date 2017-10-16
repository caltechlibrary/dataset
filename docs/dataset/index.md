
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
+ [keys](keys.html) - returns the keys to stdout, one key per line
+ [haskey](haskeys.html) - returns true is key is in collection, false otherwise
+ [path](path.html) - given a document name return the full path to document
+ [attach](attach.html) - attaches a non-JSON content to a JSON record 
+ [attachments](attachments.html) - lists any attached content for JSON document
+ [attached](attached.html) - returns attachments for a JSON document 
+ [detach](detach.html) - remove attachments to a JSON document
+ [import](import.html) - import a CSV file's rows as JSON documents
    + [import-gsheet](import-gsheet.html) - import a Google Sheets sheet rows as JSON documents
+ [export](export.html) - export a CSV file based on filtered results of collection records rendering dotpaths associated with column names
+ [extract](extract.html) - will return a unique list of unique values based on the associated dot path described in the JSON docs
+ [commands](commands.html) - return a list of command names

