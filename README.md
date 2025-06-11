Dataset Project
===============
[![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

The Dataset Project provides tools for working with collections of JSON documents easily. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It is suitable for temporary storage of JSON objects in data processing pipelines as well as a persistent storage mechanism for collections of JSON objects.

The Dataset Project provides a command line program and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science projects as well as building web applications that work with metadata.

dataset, a command line tool
----------------------------

[dataset](doc/dataset.md) is a command line tool for working with collections of [JSON](https://en.wikipedia.org/wiki/JSON) documents. Collections can be stored on the file system in a [pairtree](https://datatracker.ietf.org/doc/html/draft-kunze-pairtree-01) or stored in a SQL database that supports JSON columns like SQLite3, PostgreSQL or MySQL.

The __dataset__ command line tool supports common data management operations as

- initialization of a collection
- dump and load JSON lines files into collection
- CRUD operations on a collection
- Query a collection using SQL

See [Getting started with dataset](how-to/getting-started-with-dataset.md) for a tour and tutorial.

datasetd is dataset implemented as a web service
------------------------------------------------

[datasetd](docs/datasetd.md) is a JSON REST web service and static file host. It provides a JSON API supporting the main operations found in the __dataset__ command line program. This allows dataset collections to be integrated safely into web applications or be used concurrently by multiple processes.

The Dataset Web Service can host multiple collections each with their own custom query API defined in a simple YAML configuration file.

Design choices
--------------

__dataset__ and __datasetd__ are intended to be simple tools for managing collections JSON object documents in a predictable structured way. The dataset web service allows multi process or multi user access to a dataset collection via HTTP.

__dataset__ is guided by the idea that you should be able to work with JSON documents as easily as you can any plain text document on the Unix command line. __dataset__ is intended to be simple to use with minimal setup (e.g.  `dataset init mycollection.ds` creates a new collection called 'mycollection.ds').

- __dataset__ and __datasetd__ store JSON object documents in collections
  - Storage of the JSON documents may be either in a pairtree on disk or in a SQL database using JSON columns (e.g. SQLite3 or MySQL 8)
  - dataset collections are made up of a directory containing a collection.json and codemeta.json files.
  - collection.json metadata file describing the collection, e.g. storage type, name, description, if versioning is enabled
  - codemeta.json is a [codemeta](https://codemeta.github.io) file describing the nature of the collection, e.g. authors, description, funding
  - collection objects are accessed by their key, a unique identifier, made up of lower case alpha numeric characters
  - collection names are usually lowered case and usually have a `.ds` extension for easy identification

__dataset__ collection storage options
  - SQL store stores JSON documents in a JSON column
    - SQLite3 (default), PostgreSQL >= 12 and MySQL 8 are the current SQL databases support
    - A "DSN URI" is used to identify and gain access to the SQL database
    - The DSN URI maybe passed through the environment
  - [pairtree](https://datatracker.ietf.org/doc/html/draft-kunze-pairtree-01) (depricated, will be removed in v3)
    - the pairtree path is always lowercase
    - non-JSON attachments can be associated with a JSON document and found in a directories organized by semver (semantic version number)
    - versioned JSON documents are created along side the current JSON document but are named using both their key and semver

__datasetd__ is a web service
  - it is intended as a back end web service run on localhost
    - it runs on localhost and a designated port (port 8485 is the default)
    - supports multiple collections each can have their own configuration for global object permissions and supported SQL queries

The choice of plain UTF-8 is intended to help future proof reading dataset collections.  Care has been taken to keep _dataset_ simple enough and light weight enough that it will run on a machine as small as a Raspberry Pi Zero while being equally  comfortable on a more resource rich server or desktop environment. _dataset_ can be re-implement in any programming language supporting file input and output, common string operations and along with JSON encoding and decoding functions. The current  implementation is in the Go language.

Features
--------

[dataset](docs/dataset.md) supports
- Collection level
  - [Initialize](docs/init.md) a new dataset collection
  - Codemeta file support for describing the collection contents
  - [Dump](docs/load.md) a collection to a JSON lines document
  - [Load](docs/load.md) a collection from a JSON lines document
  - Listing [Keys](docs/keys.md) in a collection
- Object level actions
  - [create](docs/create.md)
  - [read](docs/read.md)
  - [update](docs/update.md)
  - [delete](docs/delete.md)
  - [keys](docs/keys.md)
  - [has-key](docs/has-key.md)
  - Documents as attachments
    - [attachments](docs/attachments.md) (list)
    - [attach](docs/attach.md) (create/update)
    - [retrieve](docs/retrieve.md) (read)
    - [prune](docs/prune.md) (delete)

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


Both __dataset__  and __datasetd__ maybe useful for general data science applications needing JSON object management or in implementing repository systems in research libraries and archives.


Limitations of __dataset__ and __datasetd__
-------------------------------------------

__dataset__ has many limitations, some are listed below

- the pairtree implementation it is not a multi-process, multi-user data store
- it is not a general purpose database system
- it stores all keys in lower case in order to deal with file systems 
- it stores collection names as lower case to deal with file systems that
  are not case sensitive
- **it should NOT be used for sensitive, confidential or secret information** because it lacks access controls and data encryption

__datasetd__ is a simple web service intended to run on "localhost:8485".

- it does not include support for authentication
- it does not support access control for users or roles
- it does not encrypt the data it stores
- it does not support HTTPS
- it does not provide auto key generation
- it limits the size of JSON documents stored to the size supported by
  with host SQL JSON columns
- it limits the size of attached files to less than 250 MiB
- it does not support partial JSON record updates or retrieval
- it does not provide an interactive Web UI for working with dataset
  collections
- **it should NOT be used for sensitive, confidential or secret information** because it lacks access controls and data encryption


Read next ...
-------------

- About the [dataset](docs/dataset.md) command
- About [datasetd](docs/datasetd.md) web service
- [Installation](INSTALL.md)
- [License](LICENSE)
- [Contributing](CONTRIBUTING.md)
- [Code of conduct](CODE_OF_CONDUCT.md)
- Explore __dataset__ and __datasetd__
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

Compiled versions are provided for Linux (x86, aarch64), Mac OS X (x86 and M1), Windows 11 (x86, aarch64) and Raspberry Pi OS (ARM7).

[github.com/caltechlibrary/dataset/releases](https://github.com/caltechlibrary/dataset/releases)

Related projects
----------------

You can use __dataset__ from Python via the [py_dataset](https://github.com/caltechlibrary/py_dataset) package. 

You can use __dataset__ from Deno+TypeScript by running datasetd and access it with [ts_dataset](https://github.com/caltechlibraray/ts_dataset).

