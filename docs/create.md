create
======

Syntax
------

```shell
    cat JSON_DOCNAME | dataset create COLLECTION_NAME KEY
    dataset create COLLECTION_NAME KEY JSON_VALUE
    dataset create -i JSON_DOCNAME COLLECTION_NAME KEY
```

Description
-----------

create adds a JSON document to a collection. The JSON 
document can be read from a standard in, a named file (with a 
".json" file extension) or expressed literally on the command line.

Usage
-----

In the following four examples *jane-doe.json* is a file on the 
local file system contains JSON data containing the JSON_VALUE 
of `{"name":"Jane Doe"}`.  The KEY we will create is _r1_. 
Collection is "people.ds".  The following are equivalent in 
resulting record.

```shell
    cat jane-doe.json | dataset create people.ds r1
    dataset create -i jane-doe.json people.ds r1
    dataset create people.ds r1 '{"name":"Jane Doe"}'
```

