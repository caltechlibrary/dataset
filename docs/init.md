
init
====

Syntax
------

~~~shell
    dataset init COLLECTION_NAME
~~~

Description
-----------

_init_ creates a collection. Collections are created on local 
disc. By default it uses a SQLite3 database called "collection.db"
in the dataset directory for storing JSON Objects. As of v3 only
SQLite3 is supported.

Usage
-----

The following example command create a dataset collection 
named "data.ds".

~~~shell
    dataset init data.ds
~~~

NOTE: After each evocation of `dataset init` if all went well 
you will be shown an `OK` if everything went OK, otherwise
an error message. 

~~~shell
    dataset init data.ds
~~~

