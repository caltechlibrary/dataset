dataset
=======
[![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

__dataset__ is a command line tool, a Go package, and an experimental C shared 
library for working with [JSON](https://en.wikipedia.org/wiki/JSON) 
objects as collections. Collections are stored on your local disk.
JSON objects are stored in collections as plain UTF-8 text files.
This means the objects can be accessed with common Unix text processing
tools as well as most programming languages.  __dataset__ is also available as a
Python package, see [py_dataset](https://github.com/caltechlibrary/py_dataset)

The [dataset](docs/dataset.html) command line tool supports common data 
management operations such as initialization of collections; document creation, 
reading, updating and deleting; listing keys of JSON objects in the collection.

__datasets__'s enhanced features include

- aggregate objects into data [frames](docs/frame.html)
- import, export and synchronize JSON objects to and from CSV files
- generate sample sets of keys and objects

See [Getting started with dataset](how-to/getting-started-with-dataset.html)
for a tour and tutorial. Both the command line and examples in Python 3 using
using [py_dataset](https://github.com/caltechlibrary/py_dataset) are included.


Design choices
--------------

__dataset__ isn't a database or a replacement for a repository system. 
It is tool to manage JSON documents in a predictable and structured way.
__dataset__ is guided by the idea that you should be able to work with JSON
documents as easily as you can any plain text document on Unix. __dataset__
is intended to be simple to use with minimal setup (e.g. 
`dataset init mycollection.ds` creates a new collection called 
'mycollection.ds'). It is built around the following abstractions

- dataset stores JSON objects in collections
- collections are folder(s) containing
    - collection.json metadata file
    - a pairtree of JSON object documents
    - support for attachments to JSON documents


The choice of plain UTF-8 is intended to help future proof reading
dataset collections.  Care has been taken to keep _dataset_ simple enough
and light weight enough that it will run on a machine as small as a
Raspberry Pi Zero while being equally comfortable on a more resource
rich server or desktop environment. __dataset__ can be re-implement in
any programming language supporting file input and output, common string
operations and a JSON encoding and decoding.


Example Workflow
----------------

A typical processing pattern is to write a "harvester" 
which then stores it results in a __dataset__ collection.
This is often followed by another program that transforms or aggregates
harvested material before rendering a prepared output, e.g. web pages
or data files. At Caltech Library the harvesters are typically written in Python
or Bash storing the results in a dataset collection. Depending on 
the performance needs transform and aggregates stages are written 
either in Python or Go and our final rendering stages are typically 
written in Python or as simple Bash scripts.


Features
--------

[dataset](docs/dataset) supports 

- Basic storage actions ([create](docs/create.html), [read](docs/read.html), [update](docs/update.html) and [delete](docs/delete.html))
- listing of collection [keys](docs/keys.html)
- import/export of [CSV](how-to/working-with-csv.html) files
- The ability to reshape data by performing simple object [joins](docs/join.html)
- The ability to create data [frames](docs/frame.html) from while collections or based on keys lists
    - frames are defined using [dot paths](docs/dotpath.html) describing what is to be pulled out of a stored JSON objects

You can work with dataset collections via the 
[command line tool](docs/dataset.html), via Go using the 
[dataset package](https://godoc.org/github.com/caltechlibrary/dataset) 
or in Python 3.8 using the 
[py_dataset](https://github.com/caltechlibrary/py_dataset) python package.  _dataset_ is useful for general data science applications 
which need intermediate JSON object management but not 
a full blown database.


Limitations of __dataset__
--------------------------

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store (it stores files on disk without locking)
- it is not a replacement for a repository management system
- it is not a general purpose database system
- it does not supply automatic version control on collections, objects or attachments
- it stores all keys to lower case in order to deal with file systems that are not case sensitive
- it does not have a built-in query language for filtering or sorting


Read next ...
-------------

Explore _dataset_ through 
[A Shell Example](how-to/a-shell-example.html "command line example"),
[Getting Started with Dataset](how-to/getting-started-with-dataset.html "Python examples as well as command line"),
[How To](how-to/) guides,
[topics](docs/topics.html) and [Documentation](docs/).

Releases
--------

Compiled versions are provided for Linux (x86), Mac OS X (x86 and M1), 
Windows 10 (x86) and Raspberry Pi OS (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

You can use __dataset__ from Python via the [py_dataset](https://github.com/caltechlibrary/py_dataset) package.
