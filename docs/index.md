dataset
=======

The documentation is organized around the command line options 
and as a series of "how to" style examples.

- [getting started with dataset](../how-to/getting-started-with-dataset.html) (covers both Bash and Python)
- Explore additional other [tutorials](../how-to/)
- [Overview](description.md)
- [Examples](examples.md)

Command line program documentation
----------------------------------

- [dataset](dataset.md) - usage page for managing collections with _dataset_

Internal project concepts
-------------------------

- [upgrading a collection](../how-to/upgrading-a-collection.md) - Describes how to upgrade a collection from a previous version of dataset to a new one

__dataset__ Operations
----------------------

The basic operations support by *dataset* are listed below organized 
by collection and JSON document level.

A word about keys
-----------------

__dataset__ is based around the concept of key/value pairs where
the key is the unique identifier for an object stored (i.e. the 
value) in the collection. Each storage option supported by dataset
and its own issues around what things can be called. **Keys should be
lower case alpha numeric or underscore only.** E.g. the pairtree storage
relies on the file system to store the JSON objects. Some file
systems are not case sensitive, others face challenges with
non-alpha numeric filenames.


Collection Level
----------------

- [init](init.md) creates a collection
- [history](history.md), how history is implemented
- [keys](keys.md) list keys of JSON documents in a collection, supports filtering and sorting
- [query](query.md), query lets to create a list of objects via SQL statements
- [dump](dump.md) dumps an entire collection to standard in [JSON line format](https://jsonlines.org/)
- [load](load.md) loads a [JSON line formatted stream](https://jsonlines.org/) into a collection
- [codemeta](codemeta.md)

JSON Document level
-------------------

- [create](create.md) a JSON document in a collection
- [read](read.md) back a JSON document in a collection
- [update](update.md) a JSON document in a collection
- [delete](delete.md) a JSON document in a collection
- [haskey](haskey.md) returns true if key is found in collection, false otherwise

datasetd as a web service
=========================

Since version v2.x is a web service providing access to dataset
collections. This is described in the [datasetd](datasetd.md) 
documentation page.

[datasetd](datasetd.md) supports the following end points.

Storage engines
===============

In v3 dataset JSON is stored in an SQL database with JSON column support. Currently three SQL databases can be used to store the JSON documents, SQLite 3 (default), PostgreSQL and (with manual configuration), MySQL 8. PostgreSQL.  See [storage engines](storage-engines.md) for more details.

Compatibity
===========

v3 is not compatible with v2 or earlier. You can use the v2 dump and 
use load in v3 to migrate the collection.
