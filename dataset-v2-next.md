---
title: Dataset next
draft: true
---

Dataset next
============

Ideas
-----

We've been using dataset for several years now. It has proven helpful. The
strength of dataset appears to be in severl areas. A simple clear API of
commands on the command line make it easy to pickup and use. Storing
JSON Object documents in a pairtree makes it easy to integrated into a
text friendly environment like that found on Unix systems.  The added
ability to store non-JSON documents along side the JSON document as
attachments he proven useful but could be refined to be more seemless
(e.g. you pass a semver when you attach a document).

Dataset was concieved with deliberate limitiations. This in part 
due to the options available at the time (e.g. MySQL, Postgres, 
MongDB, CouchDB, Redis) all imposed a high level of complexity to do
conceptually simple things. While many of the limitations were
deliberate it is time to consider loosing some. This should be done with
a degree of caution.

In the intervening years since starting the dataset project the NoSQL
and SQL database engines have started to converge in capabilities. This
is particularly true on the SQL engine side. SQLite 3, MySQL 8, and
Postgres 14 all have mature support for storing JSON objects in a
column. This provides an opportunity for dataset itself. It can use
those engines for storing hetrogenious collections of JSON objects. The
use case where this is particularly helpful is when running multi-user,
multi-proccess support for interacting with a dataset collection.
If dataset provides a web service the SQL engines can be used to store
the objects reliably. This allows for larger dataset collections as well as
concurrent interactions. The SQL engines provide the necessary record
locking to avoid curruption on concurrent writes.

In developing a version 2 of dataset an eye needs to be kept to several
areas --

1. reduction of complexity.
   a. simplification of codebase
   b. simplification of operation
2. learn from other systems. 
   a. align with good data practices
   b. adopt standards, e.g. codemeta for general metadata
3. **elimination** of underused "features".

With the introduction of Go 1.18 some of this can be achieved through
a better organization of code, some by applying lessons learned over the
last several years and some by reorganizing the underlying persistenent
structure of the collections themselves (e.g. take advantage of SQL
storage options as well as the tried and true pairtree model)

Proposed updates
----------------

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

Dataset has supported a form of versioning attachments for some time.
It's has not supported versioning of JSON objects, that is desirable.
Likewise the JSON support for attachments has been achieved by explicitly
passing [semver](https://semver.org/) strings when attaching a document.
This is not ideal.  The versioning process should be automatic. Retaining
a semver raises a question, what is the increment value to change?
Should you increment by major version, minor version or patch level?
There is no "right" answer for the likely use cases for dataset. The
incremented level could be set collection wide, e.g. "my_collection.ds"
might increment patch level with each update, "your_collection.ds" might
increment the major level. That needs to be explored through using the
tool. Versioning should be across the collection meaning both the JSON
documents and attachments should be versioning consistently or not
versioned at all.

Dataset frames have proved very helpful. Where possible code should be
simplified and frames should be available regardless of JSON document
storage type. As we continue to use frames in growing collections
performance will need to be improved. In practice the object
list or keys associated with a frame are used not the direct
representation of the frame in memory. This is an area suited to
refinement. The internals could be changed to improve performance as
long as the access to the keys/objects in the frame remains consistent. 
E.g. Frames don't have to be stored as plain text on disk. The code
for frames needs to be reviewed and positioned for possible evolution
as needs evolve with frame usage.

Before frames were implemented data grids were tried. In practical
usage frames replaced grids. The data grids code can be removed from
dataset. The few places where they are used in our feeds processing
are scheduled to be rewritten to use regular frames. It is a good time
to prune this "feature".

Importing, syncing and exporting to CSV is a canidate for a rethink.
While it makes it easy to get started with dataset maintaining the
ability to syncronization between a CSV representation and a dataset
collection is overly complex. While CSV support in Go is very good but
so are the Python libraries for working with CSV files. Processing
objects in a collection is more commonly done in a script (e.g. Python
using py_dataset) then directly in Go. It may make more sense to either
simplify or drop support for CSV for the version 1 level integration.
How much does Go bring to the table beyond Python? Does this need to be
"built-in" to dataset or should it be left to scripting a dataset service
or resource? Does import/export support of CSV files make dataset easier
to use beyond the library? If so does that extend to SQL tables in general?


There are generally two practices in using dataset in Caltech Library. The
command line is used interactively or Python is used programatically to
process collections (e.g. like in reporting or the feeds project).
Python has been support via a C shared library called libdataset.  While
this has worked well it also has been a challenge to maintain requiring
acccess to each os/hardware platform the cli supports.  I don't think
this is sustainable.  Since the introduction of datasetd (the web service
implementation of dataset) py_dataset could be rewritten to use the web
service implementation of dataset (i.e. datasetd) and this would fit
most of our use cases now and planned in the near future. It would avoid
some hard edge cases we've run across where the Go run time and Python
run need to be kept in sync.

Dropping libdataset support would allow dataset/datasetd to be cross
compiled for all supported platforms using only the Go tool chain.
It would make fully supporting snap installs possible.

A large area of cruft is the integrated help system. It makes more sense
to focus that on GitHub, godoc and possibly publish to a site like
readthedocs.io from the GitHub repository than to sustain a high level
of direct help integration with the cli or web service.

Goals
-----

1. Extended dataset's usefulness
2. Improve performance
3. Simplify features (e.g. prune the crufty bits)
4. Expand usage beyond Tom and Robert


Proposals
---------

In moving to version 2 there will be breaking changes.

1. (braking change) datasetd should should store data in a SQL engine
that support JSON columns, e.g. MySQL 8
   a. should improve performance and allow for better concurrent usage
   b. improve frames support
   c. facilitate integration with fulltext search engines, e.g. Lunr,
     Solr, Elasticsearch
2. Cleanup frames and clarify their behavior, position the code for
persisting frames efficiently. (e.g. explore frames implemented
using SQLite 3 database and tables)
3. Versioning of attachments needs to be automatic. A set of four
version keywords could make it easier.
   a. __set__ would set the initial version number (defaults to 0.0.0)
   b. __patch__ would increment the patch number in the semver, if
      versioning is enabled in the collection then update will assume
      patch increment
   c. __minor__ would increment the minor number and set patch to zero
   d. __major__ would increment the major number and set minor and patch
to zer
4. JSON objects should be versioned if the collection is versioned.
5. Versioning of JSON documents and attachments should be global to the
collection, i.e. everything is versioned or nothing is versioned
6. Dot notation needs reviewed. Look at how SQL databases are interacting with JSON columns. Is there a convergence in notation?
   a. [SQLite3](https://www.sqlite.org/json1.html), 
   b. [MySQL 8](https://dev.mysql.com/doc/refman/8.0/en/json.html) and
   c. [Postgres 9](https://www.postgresql.org/docs/9.3/functions-json.html)
7. Easily clone to/from pairtree and SQL stored dataset collections
8. Drop libdataset, it has been a time sync and constrainged dataset's
evolution
9. Automated migration from version 1 to version 2 databases
(via check/repair) for primary JSON documents

Leveraging SQL with JSON column support
---------------------------------------

When initializing a new SQL based collection a directory will get created
and a collections.json document will also be create.  This will help in
supporting import/export (aka cloning) of JSON collections to/from
pairtree and SQL engines.

The v1 structure of a collection is defined by a directory name (e.g.
mydataset.ds) containing a collection.json file (e.g. 
mydata.ds/collection.json).

When supporting SQL storage the collections.json should identify that the
storage type is a SQL storage engine targetted (e.g. `"sqlite3", "mysql"`)
a URI like string could be used to define the SQL stored based on Go's DNS (data source name). The storage engine could be indentified as the "protocal" in the URI.  The collection.json document SHOULD NOT require storing any secrets. Secrets can be passed via the environment. Loading a configuration should automatically check for this situation (e.g. you're running a datasetd cprocess in a container and the settings.json file needs to be stored in the project's GitHub repo)

If the "storage type" is not present it is assumed that storage
is local disk in a pairtree. Storage type is set at collection creation.
E.g.

- `init COLLECTION_NAME`, intialize dataset as a pairtree
- `init COLLECTION_NAME DSN_URI`, intialize dataset using SQLite3 or MySQL 8 for JSON document storage depending on the values in DSN_URI

A SQL based dataset collections could be stored in a single SQL database
as tables. This would allow for easier collection migration and replication.

The desired column structure of a SQL based collection could be

- `Key VARCHAR(255) NOT NULL PRIMARY KEY`
- `Object JSON`
- `Created DATETIME DEFAULT CURRENT_TIMESTAMP`
- `Updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`

NOTE: The problem is specifying automatic update timestamps isn't
standard across SQL implementations. It may make sense to only have one
or other other. This needs to be explored further.

The column structure for a SQL base frame set could be

- `Key VARCHAR(255) NOT NULL PRIMARY KEY`
- `Extract JSON` (the extracted attributes exposed by the frame)
- `Updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`

A search index could be defined as a frame with a full text index of the
atrtibutes.

Stored procedures or triggers could implement the JSON document versioning
via a copy to a history table. The latest version would be saved in the
primary table, versions would be stored in the history table where `_Key`
becomes `Key` combined with `Version`

In a pairtree implementation JSON documents could use the same 
semver settings as attachment. Need to think about how this is
organized on disk. Also attachments should not be stored in a SQL 
engine (we have big attachments). The could be stored in their own
pairtree. Using versioning on JSON documents and attachments should
function the same way but the implementation may need to very.


Code organization
-----------------

The v1 series of dataset source code is rather organic. It needs to be
structured so that it is easier to read, understand and curate. In
Go version 1.18 we can keep all the packages in the same repository.
This means code for pairtree, semver, etc. can be maintained in the
same repository easily now. This beings us an opportunity to refine
things.

- collection.go should hold the general collection level actions and collection level data structures
- frames.go should hold the frames implementation indepent of the JSON store being used
- attachments.go should hold the attachments implementation indepent of the JSON store being used
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
- semver/semver.go models semver behaviors
- dotpath/dotpath.go models dotpaths and JSON object behaviors
- pairtree/pairtree.go should hold pairtree structure and methods
- cli/cli.go should hold the outer methods for implementing the dataset CLI
  - base assumption, single user, single process
- api/api.go should hold the wrapper that implements the datasetd daemon
  - base assumption, multi user, multi process
- cmd/dataset/dataset.go is a light wrapper envoking run methods in cli
- cmd/datasetd/datasetd.go is a light wrapper envoking the run methods in ebapi.go

