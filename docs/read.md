
read
====

Syntax
------

~~~shell
    dataset3 read COLLECTION_NAME KEY
~~~

Description
-----------

The writes the JSON document to standard out (unless you've 
specific an alternative location with the "-output" option)
for the given KEY.

Usage
-----

An example we're assuming there is a JSON document with a KEY 
of "r1". Our collection name is "data.ds"

~~~shell
    dataset3 read data.ds r1
~~~

Options
-------

Normally dataset3 outputs the JSON object as presented by the storage engine.
Use the `-jsonl` to force it to a single line (JSON line format).


~~~shell
    dataset3 read -jsonl data.ds r1
~~~

