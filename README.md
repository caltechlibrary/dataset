Dataset Project
===============
[![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

The Dataset Project provides tools for working with collections of
JSON documents stored on the local file system in a pairtree or
in a SQL database supporting JSON columns. Two tools are provided
by the project -- a command line interface (dataset) and a
[RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer)
web service (datasetd).

dataset, a command line tool
----------------------------

[dataset](doc/dataset.md) is a command line tool for working with
collections of [JSON](https://en.wikipedia.org/wiki/JSON) documents.
Collections can be stored on the file system in a pairtree directory
structure or stored in a SQL database that supports JSON columns
(currently SQLite3 or MySQL 8 are supported).  Collections using the
file system store the JSON documents in a
[pairtree](https://datatracker.ietf.org/doc/html/draft-kunze-pairtree-01).
The JSON documents are plain UTF-8 source. This means the objects can be
accessed with common [Unix](https://en.wikipedia.org/wiki/Unix)
text processing tools as well as most programming languages.

The __dataset__ command line tool supports common data management operations
such as initialization of collections; document creation, reading,
updating and deleting; listing keys of JSON objects in the collection;
and associating non-JSON documents (attachments) with specific JSON
documents in the collection.

### enhanced features include

- aggregate objects into data [frames](docs/frame.md)
- generate sample sets of keys and objects
- clone a collection
- clone a collection into training and test samples

See [Getting started with dataset](how-to/getting-started-with-dataset.md) for a tour and tutorial.


datasetd, dataset as a web service
----------------------------------

[datasetd](doc/datasetd.md) is a RESTful web service implementation of the
_dataset_ command line program. It features a sub-set of capability found
in the command line tool. This allows dataset collections to be integrated
safely into web applications or used concurrently by multiple processes.
It achieves this by storing the dataset collection in a SQL database
using JSON columns.

Design choices
--------------

_dataset_ and _datasetd_ are intended to be simple tools for managing
collections JSON object documents in a predictable structured way.

_dataset_ is guided by the idea that you should be able to work with
JSON documents as easily as you can any plain text document on the Unix
command line. _dataset_ is intended to be simple to use with minimal
setup (e.g.  `dataset init mycollection.ds` creates a new collection
called 'mycollection.ds').

- _dataset_ and _datasetd_ store JSON object documents in collections.
  - Storage of the JSON documents may be either in a pairtree on disk
    or in a SQL database using JSON columns (e.g. SQLite3 or MySQL 8)
  - dataset collections are made up of a directory containing a
    collection.json and codemeta.json files.
  - collection.json metadata file describing the collection,
    e.g. storage type, name, description, if versioning is enabled
  - codemeta.json is a [codemeta](https://codemeta.github.io) file describing the nature of the collection, e.g. authors, description, funding
  - collection objects are accessed by their key, a unique identifier made of lower case alpha numeric characters
  - collection names are usually lowered case and usually have a `.ds`
    extension for easy identification

_datatset_ collection storage options
  - [pairtree](https://datatracker.ietf.org/doc/html/draft-kunze-pairtree-01) is the default disk organization of a dataset collection
    - the pairtree path is always lowercase
    - non-JSON attachments can be associated with a JSON document and
      found in a directories organized by semver (semantic version number)
    - versioned JSON documents are created along side the current JSON document but are named using both their key and semver
  - SQL store stores JSON documents in a JSON column
    - SQLite3 and MySQL 8 are the current SQL databases support
    - A "DSN URI" is used to identify and gain access to the SQL database
    - The DSN URI maybe passed through the environment

_datasetd_ is a web service
  - is intended as a back end web service run on localhost
    - by default it runs on localhost port 8485
    - supports collections that use the SQL storage engine
  - **should never be used as a public facing web service**
    - there are no user level access mechanisms
    - anyone with access to the web service end point has access to the dataset collection content


The choice of plain UTF-8 is intended to help future proof reading dataset
collections.  Care has been taken to keep _dataset_ simple enough and light
weight enough that it will run on a machine as small as a Raspberry Pi Zero
while being equally comfortable on a more resource rich server or desktop
environment. _dataset_ can be re-implement in any programming language
supporting file input and output, common string operations and along with
JSON encoding and decoding functions. The current implementation is in the
Go language.


Features
--------

[dataset](docs/dataset.md) supports
- Initialize a new dataset collection
  - Define metadata about the collection using a codemeta.json file
  - Define a keys file holding a list of allocated keys in the collection
  - Creates a pairtree for object storage
- Codemeta file support for describing the collection contents
- Simple JSON object versioning
- Listing [Keys](docs/keys.md) in a collection
- Object level actions
  - [create](docs/create.md)
  - [read](docs/read.md)
  - [update](docs/update.md)
  - [delete](docs/delete.md)
  - [keys](docs/keys.md)
  - [has-key](docs/has-key.md)
  - [sample](docs/sample.md)
  - [clone](docs/clone.md)
  - [clone-sample](docs/clone-sample.md)
  - Documents as attachments
    - [attachments](docs/attacments.md) (list)
    - [attach](docs/attach.md) (create/update)
    - [retrieve](docs/retrieve.md) (read)
    - [prune](docs/prune.md) (delete)
- The ability to create data [frames](docs/frame.md) from while
  collections or based on keys lists
  - frames are defined using a list of keys and a lost
    [dot paths](docs/dotpath.md) describing what is to be pulled out
    of a stored JSON objects and into the frame
  - frame level actions
    - frames, list the frame names in the collection
    - frame, define a frame, does not overwrite an existing frame with
      the same name
    - frame-def, show the frame definition (in case we need it for some
      reason)
    - frame-keys, return a list of keys in the frame
    - frame-objects, return a list of objects in the frame
    - refresh, using the current frame definition reload all the objects
      in the frame given a key list
    - reframe, replace the frame definition then reload the objects in
      the frame using the existing key list
    - has-frame, check to see if a frame exists
    - delete-frame remove the frame

[datasetd](docs/datasetd.md) supports

- List [collections](docs/collections-endpoint.md) available from the
  web service
- List a [collection](collection-endpoint.md)'s metadata
- List a collection's [Keys](docs/keys-endpoint.md)
- Object level actions
    - [create](docs/create-endpoint.md)
    - [read](docs/read-endpoint.md)
    - [update](docs/update-endpoint.md)
    - [delete](docs/delete-endpoint.md)
    - Documents as attachments
        - [attach](docs/attach-endpoint.md)
        - [retrieve](docs/retrieve-endpoint.md)
        - [prune](docs/prune-endpoint.md)
- The ability to create data [frames](docs/frame.md) from
  collections or based on keys lists and dot paths to form a new object
  - [dot paths](docs/dotpath.md) describing
    what is to be pulled out of a stored JSON objects

Both _dataset_  and _datasetd_ maybe useful for general data science
applications needing JSON object management or in implementing repository
systems in research libraries and archives.


Limitations of _dataset_ and _datasetd_
---------------------------------------

_dataset_ has many limitations, some are listed below

- the pairtree implementation it is not a multi-process, multi-user
  data store
- it is not a general purpose database system
- it stores all keys in lower case in order to deal with file systems
  that are not case sensitive, compatibility needed by a pairtree
- it stores collection names as lower case to deal with file systems that
  are not case sensitive
- it does not have a built-in query language, search or sorting
- it should NOT be used for sensitive or secret information

_datasetd_ is a simple web service intended to run on "localhost:8485".

- it does not include support for authentication
- it does not support a query language, search or sorting
- it does not support access control by users or roles
- it does not provide auto key generation
- it limits the size of JSON documents stored to the size supported by
  with host SQL JSON columns
- it limits the size of attached files to less than 250 MiB
- it does not support partial JSON record updates or retrieval
- it does not provide an interactive Web UI for working with dataset
  collections
- it does not support HTTPS or "at rest" encryption
- it should NOT be used for sensitive or secret information


Read next ...
-------------

- About the [dataset](docs/dataset.md) command
- About [datasetd](docs/datasetd.md) web service
- [Installation](INSTALL.md)
- [License](LICENSE)
- [Contributing](CONTRIBUTING.md)
- [Code of conduct](CODE_OF_CONDUCT.md)
- Explore _dataset_ and _datasetd_
    - [Getting Started with Dataset](how-to/getting-started-with-dataset.md "Python examples as well as command line")
    - [How To](how-to/) guides
    - [Reference Documentation](docs/).
    - [Topics](docs/topics.md)

Authors and history
-------------------

- R. S. Doiel
- Tommy Morrell

Releases
--------

Compiled versions are provided for Linux (x86), Mac OS X (x86 and M1),
Windows 11 (x86) and Raspberry Pi OS (ARM7).

[github.com/caltechlibrary/dataset/releases](https://github.com/caltechlibrary/dataset/releases)

Related projects
----------------

You can use _dataset_ from Python via the [py_dataset](https://github.com/caltechlibrary/py_dataset) package.
