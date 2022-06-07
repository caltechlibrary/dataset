---
title: Dataset next
draft: true
---

Dataset next
============

Ideas
-----

I've been using dataset for several years now. It has proven helpful. The
strength of dataset appears to be in severl areas. A simple clear API of
commands on the command line make it easy to pickup and use. Storing
JSON Object documents in a pairtree makes it easy to integrated into a
text friendly environment like that found on Unix systems.  The added
ability to store non-JSON documents along side the JSON document as
attachments he proven useful but could be refined to be more seemless
(e.g. you pass a semver when you attach a document).

Dataset has a deliberate limitiations. While most of the limitations were
deliberate it is time to consider loosing some. This should be done with
a degree of caution. An eye needs to be kept to several areas,
simplification of code and operation, reduction of complexity, elimination
of unused "features". With the introduction of Go 1.18 some of this can
be achieved through a better organization of code, some by applying
lessons learned of the last several years and some by reorganizing the
underlying persistenent structure of the collections themselves (e.g.
simply of augment the JSON documents about the collections, use
alternative means of storing JSON documents like a SQL database supporting
JSON columns). 

The metadata of a collection can be described by two JSON document.
Operational metadata (e.g. type of collection storage) is held
in a document named "collection.json". General metadata about a
collection's purpose is held in a document named "codemeta.json". 
The "codemeta.json" document should reflect the codemeta's project
for describing software and data. This has been adopted in the data
science community. 

Looking at storage options. While a pairtree is well suited for
integration into the text processing environment of Unix it is not
performant when dealing with large numbers of objects and concurrent
access. To meet the needs of scaling out a collection other options
can easily be explored. First SQL databases often support JSON columns.
This includes two commonly used in Caltech Library, i.e. MySQL 8 and
SQLite 3.  If a dataset collection is to be accessed via a web service
then using a SQL store gives us an implementation that solves concurrent
access and updates to JSON metadata. This is desirable. 

Dataset have supported a form of versioning attachments for some time.
It's has not supported versioning of JSON objects, that is desirable.
Likewise the JSON support for attachments has been achieved by explicitly
passing a semver string when attaching a document. This is not ideal.
The versioning process should be automatic but retaining a semver style
version string raises a question, what is the increment value to change?
Should you increment by major version, minor version or patch level?
There is no "right" answer for the likely use cases for dataset. The
incremented level could be set collection wide, e.g. "my_collection.ds"
might increment patch level with each update, "your_collection.ds" might
increment the major level. That needs to be explored. Also versioning
should be across the collection meaning both the JSON documents and attachments should be versioning consistently or not versioned at all.

Dataset frames has proved very helpful. Where possible code should be
simplified and frames should be available regardless of JSON document
storage type. As we continue to use frames in growing collections
performance will need to be improved. In practice the group object
list or keys associated with a frame or the primary data used from
the frame. The internals could be changed to improve performance.
They don't need necessarily be stored as plain text on disk. The code
for frames needs to be reviewed and positioned for possible evolution
as demands evolve on frames.

Before frames were implemented data grids were tried. For practical
usage frames replaced grids. The data grids code can be removed from
dataset. The few places where they are used in our feeds processing
are scheduled to be rewritten to use frames. It is a good time to
prune this "feature".

Importing and exporting to CSV is a canidate for removal. On the one
hand CSV support in Go is very good but also somewhat strict. Most
of the time when we use CSV import or export we're doing so from a 
Python program. Python also support CSV files reasonably well. It
is easy to implementing a table to object conversion in Python. How
much does Go bring to the table beyond Python? Does this need to be
"built-in" to dataset or should it be left to scripting a dataset
service or resource? 


There are generally two practices in using dataset in Caltech Library. The
command line is used interactively or Python is used to programatically
interact with collections (e.g. like in reporting or the feeds project).
Python has been support via a C shared library called libdataset.  While
this has worked well it also has been a challenge to maintain requiring
acccess to each platform we support.  I don't think this is sustatinable.
Since the introduction of datasetd (the web service implementation of
dataset) py_dataset could be rewritten to use the web service
implementation of dataset (i.e. datasetd) and this would fit most of our
use cases now and planned in the near future.

Dropping libdataset support would allow dataset/datasetd to support all
platforms where Go can be cross compile without having access to that
specific system. It would make snap installs easier.

A large area of cruft is the integrated help system. It makes more sense
to focus that on GitHub, godoc and possible publish to a site like
readthedocs.io from the GitHub repository.

Goals
-----

1. Extended dataset's usefulness
2. Improve performance
3. Simplify features (e.g. prune the crufty bits)
4. Expand usage beyond Tom and Robert


Proposals
---------

1. (braking change) datasetd should should store data in a SQL engine
  that support JSON columns, e.g. MySQL 8
  a. should improve performance and allow for better concurrent usage
  b. improve frames support
  c. facilitate integration with fulltext search engines, e.g. Lunr,
     Solr, Elasticsearch
2. Frames for a pairtree based dataset could be implemented using a
   SQLite3 db rather than a simple collection of JSON documents
   representing the frame. This should extend performance and make
   frames more flexible, this would also allow indexing attributes and
   potentionally search and sorting
3. Versioning of attachments needs to be more automatic. A set of four
   version keywords could make it easier.
  a. __set__ would set the initial version number (defaults to 0.0.0)
  b. __patch__ would increment the patch number in the semver, if
     versioning is enabled in the collection then update will assume
     patch increment
  c. __minor__ would increment the minor number and set patch to zero
  d. __major__ would increment the major number and set minor and patch
     to zer
4. v2 should support versioning JSON documents in a manner like versioned
   attachments
5. Versioning of JSON documents and attachments should be global to the
   collection, i.e. everything is versioned or nothing is versioned
6. Dot notation needs to be brought inline with JSON dot notation practices 
   used in the SQL engines with JSON column support, see
   [SQLite3](https://www.sqlite.org/json1.html), 
   [MySQL 8](https://dev.mysql.com/doc/refman/8.0/en/json.html) and
   [Postgres 9](https://www.postgresql.org/docs/9.3/functions-json.html)
7. Easy import/export to/from pairtree based dataset collections
8. Drop libdataset, it has been a time sync and constrainged dataset's
   evolution


Leveraging SQL with JSON column support
---------------------------------------

When initializing a new SQL based collection a directory will get created
and a collections.json document will also be create.  This will help in
supporting import/export of JSON collections to/from pairtree and SQL
engines.

The v1 structure of a collection is defined by a directory name (e.g.
mydataset.ds) containing a collection.json file (e.g. 
mydata.ds/collection.json).

When supporting SQL storage the collections.json should identify that the
storage type is a SQL storage engine (e.g. `"storage_type": "mysql"`) and
a porter to how to access that storage (e.g. `"storage_access": "..."`).
The collection.json document SHOULD NOT have any secrets. The access could
be passed via the environment or via a seperate file containing a DSN.  

If the "storage_type" attribute is not present it is assumed that storage
is local disk in a pairtree. Storage type is set a collection creation.
E.g.

- init, intialize dataset as a pairtree
- init-mysql, intialize dataset using MySQL 8 for JSON document storage

Additional verbs for converting collection could be

- import FROM_COLLECTION
- export TO_COLLECTION

A SQL based dataset collections could be stored in a single SQL database
as tables. This would allow for easier collection migration and replication.

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

A search index could be defined as a frame with a full text index of the
atrtibutes.

Stored procedures or triggers could implement the JSON document versioning
via a copy to a history table. The latest version would be saved in the
primary table, versions would be stored in the history table where `_Key`
becomes `_Key` combined with `Version`

In a pairtree implementation JSON documents would be in semver directory
like attachments. They could share the same versioning mechanism but be
in serpate directories, e.g. `_attachments`, `_objects`.

Attachments can be large. A decision needs to be made if attachments make
sense in a MySQL based collection or if datasetd will need to handle
concurrency and atomicity of attachment actions. E.g. a pairtree could
continue to be used but we need a means of preventing two actions (e.g.
attach, replace, removing) from colliding. Ideas might be to require
versioning on attachments and lock the versioned directory on attach,
replace and remove.

Code organization
-----------------

The v1 series of dataset source code is rather organic. It needs to be
structured so that it is easier to read, understand and curate.

- semver.go models semver behaviors
- dotpath.go models dotpaths and JSON object behaviors
- collection.go should hold the collection level actions and datastructure
- objects.go should hold the object level actions of dataset
- pairtree.go should hold pairtree structure and methods
- cli.go should hold the outer methods for implementing the dataset CLI
- webapi.go should hold the wrapper that implements the datasetd daemon
- ptstore holds the code for the pairtree local disk implementation
  - ptstore/storage.go handle mapping objects and attachments to disk in the pairtree
  - ptstore/frames.go should handling implementing frames for pairtree implementation
  - ptstore/versioning.go should handle the version mapping on disk
  - ptstore/attachments.go should hold the attachment implementation
- sqlstore holds the code hanlding a SQL engine storage using JSON columns
  - sqlstore/sql.go - SQL primatives for mapping actions to the SQL store
  - sqlstore/frames.go should hold the SQL implementation of frames
  - sqlstore/storage.go should handle mapping objects into MySQL storage
  - sqlstore/versioning.go should handle the version mapping in MySQL tables
- cmd/dataset/dataset.go is a light wrapper envoking run methods in cli
- cmd/datasetd/datasetd.go is a light wrapper envoking the run methods in ebapi.go

Questions
---------

- Should datasetd resources be managed through its own client (e.g.
  datasetctl) or use the dataset cli? Yes.
- Do all collections need a directory containing collection.json and
  codemeta.json? Yes.



