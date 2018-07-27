
# commands

Documentation on individual commands can be see with
`dataset -help COMMAND_NAME` where "COMMAND_NAME" is 
replaced with one of the commands below --

+ [init](init.html) - initialize a new collection if none exists, requires a
  path to collection
+ [create](create.html) - creates a new JSON document or replace an existing
  one in collection 
+ [read](read.html) - displays a JSON document to stdout
+ [update](update.html) - updates a JSON document in collection
+ [delete](delete.html) - removes a JSON document from collection
+ [join](join.html) - brings the functionality of jsonjoin to the dataset
  command.
+ [filter](filter.html) - takes a filter and returns an unordered list of keys
  that match filter expression
    + [dotpath](dotpath.html) - reach into an object to return a value(s)
+ [keys](keys.html) - returns the keys to stdout, one key per line
+ [haskey](haskey.html) - returns true is key is in collection, false otherwise
+ [count](count.html) - returns a count of keys in a collection
+ [path](path.html) - given a document name return the full path to document
+ [attach](attach.html) - attaches a non-JSON content to a JSON record
+ [attachments](attachments.html) - lists any attached content for JSON document
+ [detach](detach.html) - returns attachments for a JSON document
+ [prune](prune.html) - remove attachments to a JSON document
+ [import-csv](import-csv.html) - import a CSV file's rows as JSON documents
    + [import-gsheet](import-gsheet.html) - import a Google Sheets sheet rows
      as JSON documents
+ [export-csv](export-csv.html) - export a CSV file based on filtered results of
  collection records rendering dotpaths associated with column names
    + [export-gsheet](export-gsheet.html) - export a Collection of JSON
      documents to Google Sheets sheet rows
+ [extract](extract.html) - will return a unique list of unique values based on
  the associated dot path described in the JSON docs
    + [dotpath](dotpath.html) - reach into an object to return a value(s)
+ [check](check.html) - will check a collection against current version
+ [repair](repair.html) - will attempt to repair/upgrade a collection
+ [migrate](migrate.md) - will migrate the file layout (e.g. buckets, pairtree)

NOTE: The options create, update can read JSON documents piped 
from standard in if you use the '-i -' or '-include -' option. 
Likewise keys can be read from standard input with the '-i -' 
or '-include -' options for read, list, keys and count.

