dump
============

This will dump all the JSON objects in a collection, one
object per line (see https://jsonlines.org) (aka JSONL). 

The objects are written to standard output. Dump is the complement of
load verb. The objects dumped reflect the structured using when storing
objects in an SQLite3 database regardless of the store of the specific
collection. Like clone it provides a means of easily moving your data out
of a dataset collection.

Example
-------

~~~shell
    dataset3 dump mycollection.ds >mycollection.jsonl
~~~

