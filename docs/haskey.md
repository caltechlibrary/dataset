
haskey
======

Syntax
------

~~~shell
    dataset3 [OPTIONS] haskey COLLECTION_NAME KEY_TO_CHECK_FOR
~~~

Description
-----------

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

Usage
-----

~~~shell
    dataset3 haskey people.ds '0000-0003-0900-6903'
    dataset3 haskey people.ds r1
~~~

