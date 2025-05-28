
Dataset Project, v3
===================
 
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

[![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

Dataset Project provides tools for working with JSON documents as a collection hosted in a SQL JSON store (e.g. SQLite3). It is an exploration of the use and usefulness of JSON documents in the setting of libraries, archives and museums. It may have application further afield as metadata management benefits a wide range of application domains.

Dataset v3 stores JSON documents using a SQL engine with JSON column support. These days that most popular SQL implementations support JSON columns (e.g. SQLite3, PostgreSQL, MySQL). V3 uses a extremely simple table structure for both the current object state and object history as well as integration with the SQL's growing support for working with JSON objects generally.

Two tools are provided by the Dataset Project v3

[dataset3](dataset3.1.md)
: is a command line interface for managing JSON documents. 

[dataset3d](dataset3d.1.md)
: JSON Web Service. This service supports file attachments.

Design choices
--------------

__dataset3__ and __dataset3d__ are intended to be simple tools for managing collections JSON object documents in a predictable structured way.

__dataset3__ is guided by the idea that you should be able to work with JSON documents as easily as you can any plain text document on the Unix command line. _dataset_ is intended to be simple to use with minimal setup (e.g.  `dataset3 init mycollection.ds` creates a new collection called 'mycollection.ds').

__dataset3__ and __dataset3d__ store JSON object documents in collections.

- Storage in a SQL database using JSON columns (example SQLite3, or PostgreSQL and MySQL 8 with extra configuration)
- dataset collections are made up of a directory containing a collection.json and codemeta.json files.
- collection.json metadata file describing the collection, e.g. storage type, name, description, if versioning is enabled
- codemeta.json is a [codemeta](https://codemeta.github.io) file describing the nature of the collection, e.g. authors, description, funding
- collection objects are accessed by their key, a unique identifier made of lower case alpha numeric characters
- collection names are usually lowered case and usually have a `.ds` extension for easy identification


The __dataset3d__ provides a web service has the following constraints.

- is intended as a back end web service run on localhost
  - by default it runs on localhost port 8485
  - supports collections that use the SQL storage engine
- **should never be used as a public facing web service**
  - there are no user level access mechanisms
  - anyone with access to the web service end point has access to the dataset collection content


The choice of plain UTF-8 is intended to help future proof reading dataset collections.  Care has been taken to keep _dataset_ simple enough and light weight enough that it will run on a machine as small as a Raspberry Pi Zero while being equally comfortable on a more resource rich server or desktop environment. _dataset_ can be re-implement in any programming language supporting file input and output, common string operations and along with JSON encoding and decoding functions. The current implementation is in the Go language.


Features
--------

-[dataset3](docs/dataset3.md) supports
- Initialize a new dataset collection
  - Define metadata about the collection using a codemeta.json file
  - Define a keys file holding a list of allocated keys in the collection
  - Creates a pairtree for object storage
- Codemeta file support for describing the collection contents
- Simple JSON object versioning
- Listing [Keys](docs/keys.md) in a collection
- SQL [query](docs/query.md) of collection
- Object level actions
  - [create](docs/create.md)
  - [read](docs/read.md)
  - [update](docs/update.md)
  - [delete](docs/delete.md)
  - [keys](docs/keys.md)
  - [haskey](docs/has-key.md)

If attachments are configured via a `dataset3d` YAML configuration then the following additional verbs are `attachments` attribute. The dataset tools will need the appropriate file permissions to support attachment operations. The YAML configuration should have the same name as the dataset collection but end with the `.yaml` instead of `.ds`. E.g. if my collection is called "mydata.ds" the attachment configuration is "mydata.yaml" in the same directory.

[dataset3d](docs/dataset3d.md) supports

- List [collections](docs/collections-endpoint.md) available from the web service
- List a [collection](collection-endpoint.md)'s metadata
- List a collection's [Keys](docs/keys-endpoint.md)
- Object level actions
  - [create](docs/create-endpoint.md)
  - [read](docs/read-endpoint.md)
  - [update](docs/update-endpoint.md)
  - [delete](docs/delete-endpoint.md)
  - [query](docs/query-endpoint.md) (for queries defined in the configuration YAML file)

Both __dataset3__  and __dataset3d__ maybe useful for general data science applications needing JSON object management or in implementing repository systems in research libraries and archives.


Limitations of __dataset3__ and __dataset3d__
---------------------------------------------

__dataset3__ has many limitations, some are listed below

- it is not a general purpose database system
- it stores all keys in lower case in order to deal with file systems
  that are not case sensitive, compatibility needed by a pairtree
- it stores collection names as lower case to deal with file systems that
  are not case sensitive
- it does not have a built-in query language, search or sorting. Instead
  it relies on the SQL engine's SQL dialect for manipulating the JSON objects
- it should NOT be used for sensitive or secret information

__dataset3__ is a simple web service intended to run on local host, e.g. "localhost:8485".

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

### Authors

- Doiel, R. S.
- Morrell, Thomas E

### Maintainers

- Doiel, R. S.
- Morrell, Thomas E

Software Requirements
---------------------

- [Golang](https://golang.org) &gt;&#x3D; 1.24.2
- [CMTools](https://caltechlibrary.github.io/CMTools) &gt;&#x3D; 0.0.29
- [Pandoc](https://pandoc.org) &gt;&#x3D; 3.1
- GNU Make &gt;&#x3D; 3.8

Software Suggestions
--------------------

- [jq](https://jqlang.org) &gt;&#x3D; 1.7

Read next ...
-------------

- Explore _dataset_ and __dataset3d__
  - [Getting Started with Dataset](how-to/getting-started-with-dataset.md "Python examples as well as command line")
  - [User Manual](user_manual.md)
  - [Reference Documentation](docs/).
  - [Topics](docs/topics.md)
- [Getting Help, Reporting bugs](https://github.com/caltechlibrary/dataset/issues)
- [LICENSE](https://caltechlibrary.github.io/dataset/LICENSE)
- [Installation](INSTALL.md)
- [Contributing](CONTRIBUTING.md)
- [Code of conduct](CODE_OF_CONDUCT.md)
- [About](about.md)

