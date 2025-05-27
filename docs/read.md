
read
====

Syntax
------

~~~shell
    dataset read COLLECTION_NAME KEY
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
    dataset read data.ds r1
~~~

Options
-------

Normally dataset adds two values when it stores an object, `._Key`
and possibly `._Attachments`. You can get the object without these
added attributes by using the `-c` or `-clean` option.


~~~shell
    dataset read -clean data.ds r1
~~~

