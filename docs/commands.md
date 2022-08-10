commands
========

Documentation on individual commands can be see with

```shell
    dataset -help COMMAND_NAME
```

where "COMMAND_NAME" is 

replaced with one of the commands below --

- [init](init.html) - initialize a new collection if none exists, requires a
  path to collection
- [keys](keys.html) - returns the keys to stdout, one key per line
- [has-key](haskey.html) - returns true is key is in collection, false otherwise
- [create](create.html) - creates a new JSON document or replace an existing
  one in collection 
- [read](read.html) - displays a JSON document to stdout
- [update](update.html) - updates a JSON document in collection
- [delete](delete.html) - removes a JSON document from collection
- Data frames
    - [frames](frames.html) - list the data frames defined for a collection
    - [frame](frame.html) - defines a new data frame 
    - [frame-def](frame-def.html) - returns a frame's object definition
    - [frame-keys](frame-keys.html) - returns a frame's key list
    - [frame-objects](frame-objects.html) - returns a frame's object list
    - [reframe](reframe.html) - uses the existing frame definition replacing all objects using a new key list
    - [refresh](refresh.html) - updates the objects in a data frame based on the current status of the collection.
    - [delete-frame](delete-frame.html) - remove a frame from a collection
- [count](count.html) - returns a count of keys in a collection
- Attachments
    - [attachments](attachments.html) - lists any attached content for JSON document
    - [attach](attach.html) - attaches a non-JSON content to a JSON record
    - [retrieve](retrieve.html) - returns attachments for a JSON document
    - [prune](prune.html) - remove attachments to a JSON document
- [check](check.html) - will check a collection against current version (for pairtree storage collections)
- [repair](repair.html) - will attempt to repair/upgrade a collection (for pairtree storage collections)

NOTE: The options create, update can read JSON documents piped 
from standard in if you use the '-i -' or '-include -' option. 
Likewise keys can be read from standard input with the '-i -' 
or '-include -' options for read, list, keys and count.

