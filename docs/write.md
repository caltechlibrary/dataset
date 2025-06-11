
create
======

Syntax
------

~~~shell
    cat JSON_DOCNAME | dataset3 write COLLECTION_NAME KEY
    dataset3 write -i JSON_DOCNAME COLLECTION_NAME KEY
    dataset3 write COLLECTION_NAME KEY JSON_VALUE
    dataset3 write COLLECTION_NAME KEY JSON_FILENAME
~~~

Description
-----------

write adds or replaces a JSON document to a collection. The JSON 
document can be read from a standard in, a named file (with a 
".json" file extension) or expressed literally on the command line.

Usage
-----

In the following four examples *jane-doe.json* is a file on the 
local file system contains JSON data containing the JSON_VALUE 
of `{"name":"Jane Doe"}`.  The KEY we will write is _r1_. 
Collection is "people.ds".  The following are equivalent in 
resulting record.

~~~shell
    cat jane-doe.json | dataset3 write people.ds r1
    dataset3 write -i blob.json people.ds r1
    dataset3 write people.ds r1 '{"name":"Jane Doe"}'
    dataset3 write people.ds r1 jane-doe.json
~~~

