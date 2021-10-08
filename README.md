Dataset Project
===============
[![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

The Dataset Project provides tools for working with collections of JSON Object documents stored on the local file system.  Two tools are provided.

dataset command line tool
-------------------------

[dataset](doc/dataset.html) is a command line tool for working with collections of [JSON](https://en.wikipedia.org/wiki/JSON) objects. Collections are stored on the file system.  JSON objects are stored in collections as plain UTF-8 text files.  This means the objects can be accessed with common [Unix](https://en.wikipedia.org/wiki/Unix) text processing tools as well as most programming languages.

The _dataset_ command line tool supports common data management operations such as initialization of collections; document creation, reading, updating and deleting; listing keys of JSON objects in the collection; and associating non-JSON documents (attachments) with specific JSON documents in the collection.

### enhanced features include

- aggregate objects into data [frames](docs/frame.html)
- import, export and synchronize JSON objects to and from CSV files
- generate sample sets of keys and objects

See [Getting started with dataset](how-to/getting-started-with-dataset.html) for a tour and tutorial.

dataset as a web service
------------------------

[datasetd](doc/datasetd) is a web service implementation of the _dataset_ command line program. It features a sub-set of capability found in the command line tool. This allows dataset collections to be integrated safely into other web applications or used by multiple processes.

Design choices
--------------

_dataset_ and _datasetd_ are intended to be simple tools for managing collections JSON object documents in a predictable structured way. 

_dataset_ and _datasetd_ are guided by the idea that you should be able to work with JSON documents as easily as you can any plain text document on Unix command line. _dataset_ is intended to be simple to use with minimal setup (e.g.  `dataset init mycollection.ds` creates a new collection called 'mycollection.ds'). 
- _dataset_ and _datasetd_ store JSON object documents in collections
    - collections are folder(s) containing
        - collection.json metadata file describing the collection and keys
        - a pairtree of JSON object documents
        - non-JSON attachments can be associated with a JSON document and found in a semver (semantic version number) named sub directory


The choice of plain UTF-8 is intended to help future proof reading dataset collections.  Care has been taken to keep _dataset_ simple enough and light weight enough that it will run on a machine as small as a Raspberry Pi Zero while being equally comfortable on a more resource rich server or desktop environment. _dataset_ can be re-implement in any programming language supporting file input and output, common string
operations and a JSON encoding and decoding. The current implementation is in the Go language.


Features
--------

[dataset](docs/dataset) supports 

- Listing [Keys](docs/keys.html) in a collection
- Object level actions
    - [create](docs/create.html)
    - [read](docs/read.html)
    - [update](docs/update.html)
    - [delete](docs/delete.html)
    - Documents attachments
        - [attach](docs/attach.html)
        - [retrieve](docs/retrieve.html)
        - [prune](docs/prune.html)
- Import and export of [CSV](how-to/working-with-csv.html) files
- The ability to reshape data by performing simple object [joins](docs/join.html)
- The ability to create data [frames](docs/frame.html) from while collections or based on keys lists
    - frames are defined using [dot paths](docs/dotpath.html) describing what is to be pulled out of a stored JSON objects

[datasetd](docs/datasetd) supports

- List collections available from the web service
- List collection [Keys](docs/keys.html)
- Object level actions
    - [create](docs/create.html)
    - [read](docs/read.html)
    - [update](docs/update.html)
    - [delete](docs/delete.html)
    - Documents attachments
        - [attach](docs/attach.html)
        - [retrieve](docs/retrieve.html)
        - [prune](docs/prune.html)

Both  _dataset_  and _datasetd_ are useful for general data science applications which need intermediate JSON object management but not a full blown database or repository system.


Limitations of _dataset_ and _datasetd_
-------------------------------------------

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store
- it is not a general purpose database system
- it does not supply automatic version control on collections, objects or attachments
- it stores all keys to lower case in order to deal with file systems that are not case sensitive
- it does not have a built-in query language, search or sorting
- it should NOT be used for sensitive or secret information

_datasetd_ is a simple web service intended to run on "localhost:8485".

- it is not a RESTful service
- it does not include support for authentication
- it does not support a query language, search or sorting
- it does not support data frames
- it does not support access control by users or roles
- it does not provide auto key generation or versioning
- it limits the size of JSON documents stored to less than 1 MiB
- it limits the size of attachment files to less than 250 MiB
- it does not support partial JSON record updates or retrieval
- it does not provide an interactive Web UI for working with dataset collections
- it does not support HTTPS or "at rest" encryption
- it should NOT be used for sensitive or secret information


Read next ...
-------------

- About the [dataset](doc/dataset.html) command
- About [datasetd](doc/datasetd.html) web service
- [Installation](install.html)
- [License](license.html)
- [Contributing](contributing.html)
- [Code of conduct](code_of_conduct.html)
- Explore _dataset_ and _datasetd_
    - [A Shell Example](how-to/a-shell-example.html "command line example")
    - [Getting Started with Dataset](how-to/getting-started-with-dataset.html "Python examples as well as command line")
    - [How To](how-to/) guides
    - [Reference Documentation](docs/).
    - [Topics](docs/topics.html)

Authors and history
-------------------

- R. S. Doiel
- Tommy Morrell

Releases
--------

Compiled versions are provided for Linux (x86), Mac OS X (x86 and M1), Windows 10 (x86) and Raspberry Pi OS (ARM7).  See https://github.com/caltechlibrary/dataset/releases.

Related projects
----------------

You can use _dataset_ from Python via the [py_dataset](https://github.com/caltechlibrary/py_dataset) package.
