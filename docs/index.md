dataset3
========

The documentation is organized around the command line options and as a series of "how to" style examples.

- [overview](description.md)
- [examples](examples.md)
- [getting started with dataset](getting_started_with_dataset.html) (covers both Bash and Python)
- [A shell example using dataset3 and dataset3d](a_shell_example.bash)
- Special topcis
  - [Using Dataset from Windows' command prompt and powershell](Windows-cmd-prompt.md)
 

Command line program documentation
----------------------------------

- [dataset3](dataset3.md) - usage page for managing collections with __dataset3__

__dataset3__ Operations
-----------------------

The basic operations support by __dataset3__ are listed below organized by collection and JSON document level.

A word about keys
-----------------

__dataset3__ is based around the concept of key/value pairs where the key is the unique identifier for an object stored (i.e. the value) in the collection. Each storage option supported by dataset and its own issues around what things can be called. **Keys should be lower case alpha numeric or underscore only.** E.g. the pairtree storage relies on the file system to store the JSON objects. some file systems are not case sensitive, others face challenges with non-alpha numeric filenames.


Collection Level
----------------

- [init](init.md) creates a collection
- [keys](keys.md) list keys of JSON documents in a collection, supports filtering and sorting
- [dump](dump.md) dumps an entire collection to standard in [JSON line format](https://jsonlines.org/)
- [load](load.md) loads a [JSON line formatted stream](https://jsonlines.org/) into a collection
- [query](query.md), query lets to create a list of objects via SQL statements
- [codemeta](codemeta.md)
- [history](history.md), how history is implemented

JSON Document level
-------------------

- [create](create.md) a JSON document in a collection
- [read](read.md) back a JSON document in a collection
- [update](update.md) a JSON document in a collection
- [delete](delete.md) a JSON document in a collection
- [haskey](haskey.md) returns true if key is found in collection, false otherwise

dataset3d as a web service
==========================

Since version v2.x is a web service providing access to dataset
collections. This is described in the [dataset3d](dataset3d.md) 
documentation page.

[dataset3d](dataset3d.md) supports the following end points.

### End points

- [Collections Endpoint](collections-endpoint.md)
- [Collection Endpoint](collection-endpoint.md)
- [Create Endpoint](create-endpoint.md)
- [Read Endpoint](read-endpoint.md)
- [Update Endpoint](update-endpoint.md)
- [Delete Endpoint](delete-endpoint.md)
- [Keys Endpoint](keys-endpoint.md)
- [Query Endpoint](query-endpoint.md)

Storage engines
===============

In v3 dataset JSON is stored in an SQL database with JSON column support. Currently three SQL databases can be used to store the JSON documents, SQLite 3 (default), PostgreSQL and (with manual configuration), MySQL 8. PostgreSQL.  See [storage engines](storage-engines.md) for more details.

- [Storage Engines](storage-engines.md)

Compatibity
===========

v3 is not compatible with v2 or earlier. You can use the last dataset v2 to dump a collections objects and then 
load them in with __dataset3__ load.

~~~shell
# Dump the v2 collection as a JSON lines document
dataset dump old_data.ds >data.jsonl
# Initialize the v3 collection and then load the JSON line document
dataset3 init new_data.ds
dataset3 load new_data.ds <data.jsonl
~~~

