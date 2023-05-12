init
====

Syntax
------

```shell
    dataset init COLLECTION_NAME [DSN_URI]
```

Description
-----------

_init_ creates a collection. Collections are created on local 
disc.

Usage
-----

The following example command create a dataset collection 
named "data.ds".

```shell
    dataset init data.ds
```

NOTE: After each evocation of `dataset init` if all went well 
you will be shown an `OK` if everything went OK, otherwise
an error message. 

By default dataset cli creates pairtree collections. You can now optionally 
store your documents in a SQL database (e.g. SQLite3, MySQL 8). This can
improve performance for large collections as well as support multi-user or
multi-process concurrent use of a collection. To use a SQL storage engine
you need to provide a "DSN_URI". The DSN_URI is formed by setting the "protocl" of the URL to either "sqlite://" or "mysql://" followed by a DSN
(data source name) as described by the database/sql package in Go.

This examples shows using SQLite3 storage for the JSON documents in
a "collection.db" stored inside the "data.ds" collection.

```shell
    dataset init data.ds "sqlite://collection.db"
```

Here's a variation using MySQL 8 as the storage engine storing the
collection in the "collections" database.

```shell
    dataset init data.ds "mysql://DB_USER:DB_PASSWD@/collections"
```


