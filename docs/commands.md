commands
========

Documentation on individual commands can be see with

```shell
    dataset -help COMMAND_NAME
```

where "COMMAND_NAME" is 

replaced with one of the commands below --

- [init](init.md) - initialize a new collection if none exists, requires a
  path to collection
- [keys](keys.md) - returns the keys to stdout, one key per line
- [has-key](haskey.md) - returns true is key is in collection, false otherwise
- [create](create.md) - creates a new JSON document or replace an existing
  one in collection 
- [read](read.md) - displays a JSON document to stdout
- [update](update.md) - updates a JSON document in collection
- [delete](delete.md) - removes a JSON document from collection
- [count](count.md) - returns a count of keys in a collection
- Attachments
    - [attachments](attachments.md) - lists any attached content for JSON document
    - [attach](attach.md) - attaches a non-JSON content to a JSON record
    - [retrieve](retrieve.md) - returns attachments for a JSON document
    - [prune](prune.md) - remove attachments to a JSON document
- [dump](dump.md) - export collection to a JSON lines file
- [load](load.md) - import a collection using a JSON lines file

NOTE: The options create, update can read JSON documents piped 
from standard in if you use the '-i -' or '-include -' option. 
Likewise keys can be read from standard input with the '-i -' 
or '-include -' options for read, list, keys and count.
