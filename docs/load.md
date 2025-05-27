load [OPTION]
============

This will read a JSONL stream (JSONL, see https://jsonlines.org) 
and store the objects in a dataset collection. The objects must have
two attributes, __key__ and __object__.  The key value is used as
the key assigned in the collection and object value is the JSON stored.
Load reads from standard input.

The object structure matches the schema used by dataset collections
that are using SQLite3 for their object store. The dataset collection
loading the objects can use any dataset collection storage format
supported by the version of dataset featuring load and dump verbs.

# OPTION

-o, -overwrite
: If an object exists in the collection with the same key replace it.

-m, -max-capacity INTEGER
: Objects can be large in JSONL so you have the option of setting the
maximum buffer size for a single object. The integer value should be
greater than zero. The unit value is measured in mega bytes. "1" is
one meta byte, "10" would be ten mega bytes.

# EXAMPLE

Load a JSONL file, duplicate objects will not be overwritten.

~~~shell
    dataset load mycollection.ds <mycollection.jsonl
~~~

Load a JSONL file, duplicate objects will be overwritten.

~~~shell
    dataset load -overwrite mycollection.ds <mycollection.jsonl
~~~

