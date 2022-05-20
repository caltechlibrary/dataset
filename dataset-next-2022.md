---
title: Dataset next
draft: true
---

Dataset next
============

Ideas
-----

I've been using dataset for several years now. It has proven helpful. The strength of dataset appears to be in severl areas. A simple clear API of commands on the command line make it easy to pickup and make useful. Storing JSON Object documents in a pairtree makes it easy to integrated into a text friendly environment like that found on Unix systems.  The added ability to store non-JSON documents along side the JSON document as attachments he proven useful but could be refined to be more seemless (e.g. rather than passing an explicit version number you could increas the semver by passing a keyword like patch, minor, major which would then increment the semver appropriately).

Dataset has a deliberate limitiations. Those limitations should be lossened with caution.  The metadata of a collcetion is described in a JSON document that needs to fit into memory. This includes a map of keys. Likewise various actions like frames are done by loading JSON documents into memory making choices and updating the in memory frame metadata before writing out all the frame in done JSON document. This limits it's performance. Evaluation involves opening a JSON document in the pairtree as well as updating a potentially large JSON frame document.  While working with data in a batch mode this generally isn't a problem beyond being slow. It does limit the total number of documents you can work with. That limitation doesn't apply in database systems. SQL engines with JSON column support offer a way to scale.

Taking advantage of mature databases platforms like SQLite3 (for frame storage) and MySQL 8 for web service hosted dataset implementation could be a way to improve performance as well as allow improved concurrent usage. This is an expansion of the origin limitation of single process/user interacting with a dataset collections. If pursued it needs to be easily to import/export between filesystem/pairtree collections and SQL stored collections. There is a question of the cli supporting both types of storage or if there should be a clear split between pairtree storage and SQL storage with the cli handling pairtrees and the daemon handling SQL storage.

In the v1 series of dataset frames were introduced but calculating the frame is generally slow for larger datasets. This is partiailly the overhead of storing records as individual documents. A faster implementation of frames could leverage SQL storage supporting JSON columns. This would be true for both storage and for updates as the frame elements could be manimulated individually rather than having to read the whole frame set into memory, make the updates and write them out again. If the frame set is a table filtering could be done as a join on the primary object table with the frame rows storing the key and extracted values.

There are generally two practices in using dataset in Caltech Library. The command line is used interactively or Python is used to programatically interact with collections (e.g. like in reporting or the feeds project). Python has been support via a C shared library called libdataset.  While this has worked well it also has been a challenge to maintain requiring acccess to each platform we support.  I don't think this is sustatinable. Since the introduction of datasetd (the web service implementation of dataset) py_dataset could be rewritten to use the web service implementation of dataset (i.e. datasetd). This would allow dataset/datasetd to support any platform where Go canbe cross compile and would help if dataset is distributed as a snap. For that to be practable the web implementation would need to support all the frame interactions and may need to support attachments as well.

Another area of cruft is the integrated help system. It makes more sense to focus that on GitHub, godoc and possible publish to a site like readthedocs.io from the GitHub repository.

Goals
-----

1. Extended dataset's usefulness
2. Improve performance
3. Simplify features (e.g. prune the crufty bits)
4. Expand usage beyond Tom and Robert


Proposals
---------

1. (braking change) datasetd should should store data in a SQL engine that support JSON columns, e.g. MySQL 8
  a. should improve performance and allow for better concurrent usage
  b. improve frames support
  c. allow for fulltext search indexes to be defined on specific attributes
2. Frames for a pairtree based dataset could be implemented using a SQLite3 db rather than a simple collection of JSON documents representing the frame. This should extend performance and make frames more flexible, this would also allow indexing attributes and search
3. Versioning of attachments needs to be more automatic. A set of four version keywords could make it easier.
  a. __set__ would set the initial version number (defaults to 0.0.0)
  b. __patch__ would increment the patch number in the semver, if versioning is enabled in the collection then update will assume patch increment
  c. __minor__ would increment the minor number and set patch to zero
  d. __major__ would increment the major number and set minor and patch to zer
4. Version JSON documents as well as attachments
5. Versioning of JSON documents and attachments should be global to the collection
6. Dot notation needs to be brought inline with JSON dot notation practices used in the SQL engines with JSON column support, see [SQLite3](https://www.sqlite.org/json1.html), [MySQL 8](https://dev.mysql.com/doc/refman/8.0/en/json.html) and [Postgres 9](https://www.postgresql.org/docs/9.3/functions-json.html)
7. Easy import/export to/from pairtree based dataset collections
8. Drop libdataset


Leveraging SQL with JSON column support
---------------------------------------

When initializing a new SQL based collection a directory will get created and a collections.json document will also be create.  This will help in supporting import/export of JSON collections to/from pairtree and SQL engines.

The v1 structure of a collection is defined by a directory name (e.g. mydataset.ds) containing a collection.json file (e.g. mydata.ds/collection.json).

When supporting SQL storage the collections.json should identify that the storage type is a SQL storage engine (e.g. `"storage_type": "mysql"`) and a porter to how to access that storage (e.g. `"storage_access": "..."`). The collection.json document SHOULD NOT have any secrets. The access could be passed via the environment or via a seperate file containing a DSN.  

If the "storage_type" attribute is not present it is assumed that storage is local disk in a pairtree. Storage type is set a collection creation. E.g.

- init, intialize dataset as a pairtree
- init-mysql, intialize dataset using MySQL 8 for JSON document storage

Additional verbs for converting collection could be

- import FROM_COLLECTION
- export TO_COLLECTION

A SQL based dataset collections could be stored in a single SQL database as tables. This would allow for easier collection migration and replication.

The column structure of a SQL based collection could be

- `_Key VARCHAR(255) NOT NULL PRIMARY KEY`
- `Object JSON`
- `Version VARCHAR DEFAULT 0.0.0`
- `Created DATETIME DEFAULT CURRENT_TIMESTAMP`
- `Updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`

The column structure for a SQL base frame set could be

- `_Key VARCHAR(255) NOT NULL PRIMARY KEY`
- `Extract JSON` (the extracted attributes exposed by the frame)
- `Updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`

A search index could be defined as a frame with a full text index of the atrtibutes.

Stored procedures or triggers could implement the JSON document versioning via a copy to a history table. The latest version would be saved in the primary table, versions would be stored in the history table where `_Key` becomes `_Key` combined with `Version`

In a pairtree implementation JSON documents would be in semver directory like attachments. They could share the same versioning mechanism but be in serpate directories, e.g. `_attachments`, `_objects`.

Attachments can be large. A decision needs to be made if attachments make sense in a MySQL based collection or if datasetd will need to handle concurrency and atomicity of attachment actions. E.g. a pairtree could continue to be used but we need a means of preventing two actions (e.g. attach, replace, removing) from colliding. Ideas might be to require versioning on attachments and lock the versioned directory on attach, replace and remove.

Code organization
-----------------

The v1 series of dataset source code is rather organic. It needs to be structured so that it is easier to read, understand and curate.

- semver.go models semver behaviors
- dotpath.go models dotpaths and JSON object behaviors
- collection.go should hold the collection level actions and datastructure
- objects.go should hold the object level actions of dataset
- pairtree.go should hold pairtree structure and methods
- cli.go should hold the outer methods for implementing the dataset CLI
- daemon.go should hold the wrapper that implements the datasetd daemon
- ptstore holds the code for the pairtree local disk implementation
  - ptstore/storage.go handle mapping objects and attachments to disk in the pairtree
  - ptstore/frames.go should handling implementing SQLite3 frames for pairtree implementation
  - ptstore/versioning.go should handle the version mapping on disk
  - ptstore/attachments.go should hold the attachment implementation
- sqlstore holds the code hanlding a SQL engine storage using JSON columns
  - sqlstore/sql.go - SQL primatives for mapping actions to the SQL store
  - sqlstore/frames.go should hold the SQL implementation of frames
  - sqlstore/storage.go should handle mapping objects into MySQL storage
  - sqlstore/versioning.go should handle the version mapping in MySQL tables
- cmd/dataset/dataset.go is a light wrapper envoking run methods in cli
- cmd/datasetd/datasetd.go is a light wrapper envoking the run methods in daemon.go

Questions
---------

- Should datasetd resources be managed through its own client (e.g. datasetctl) or use the dataset cli?
