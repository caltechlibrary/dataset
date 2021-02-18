[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)


# dataset   [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

_dataset_ is a command line tool, Go package, and an experimental C shared 
library for working with [JSON](https://en.wikipedia.org/wiki/JSON) 
objects as collections. Collections can be stored on local disc.
JSON objects are stored in collections as 
plain UTF-8 text. This means the objects can be accessed with common 
Unix text processing tools as well as most programming languages. 
_dataset_ is also available as a Python package, see 
[py_dataset](https://github.com/caltechlibrary/py_dataset)

The [dataset](docs/dataset.html) command line tool supports common data 
manage operations such as initialization of collections, creation, 
reading, updating and deleting JSON objects in the collection. Some of 
its enhanced features include the ability to generate data 
[frames](docs/frame.html) as well as the ability to 
import, export and synchronize JSON objects to and from CSV files. 

_dataset_ is written in the [Go](https://golang.org) programming language.
It can be used as a Go package by other Go based software. Go supports
generating C shared libraries. By compiling the Go source you can
create a _libdataset_ C shared library. The C shared library is currently
being used by the Digital Library Development Group in Caltech Library from
Python 3.8 (see [py_dataset](https://github.com/caltecehlibrary/py_dataset "link to github repo for py_dataset")).
This approach looks promising if you need support from other programming
languages (e.g. [Julia](https://julialang.org/) can call shared libraries
easily with a ccall function). 


See [getting-started-with-datataset.md](how-to/getting-started-with-dataset.html) for a tour and tutorial. Include are both the command line as well
as examples in Python using [py_dataset](https://github.com/caltechlibrary/py_dataset).


## Design choices

_dataset_ isn't a database or a replacement for repository systems. 
It is guided by the idea that you should be able to work with text 
files, the JSON objects documents, with standard Unix text utilities.
It is intended to be simple to use with minimal setup (e.g. 
`dataset init mycollection.ds` creates a new collection called 
'mycollection.ds'). It is built around a few abstractions --
dataset stores JSON objects in collections, collections are folder(s) 
containing a pairtree of JSON object documents and any attachments, a 
collections.json file describing the mapping of keys to folder locations).
_dataset_ takes minimal system resources and keeps all content, 
except JSON object attachments, in plain UTF-8 text. 

The choice of plain UTF-8 and future proof reading dataset collections.  
Care has been taken to keep _dataset_ simple enough and light weight 
enough that it will run on a machine as small as a Raspberry Pi while 
being equally comfortable on a more resource rich server or desktop 
environment. It should be easy to do alternative implementations
in any language having a good string library, JSON support and memory
management.


## Workflows

A typical library processing pattern is to write a "harvester" 
which then stores it results in a _dataset_ collection. Write something
that transforms or aggregates harvested options and then write
a final rendering program to prepare the data for the web. The
the hearvesters are typically written in Python or as a simple Bash
scripts storing the results in a dataset collection. Depending on 
the performance needs transform and aggregates stages are written 
either in Python or Go and our final rendering stages are typically 
written in Python or as simple Bash scripts.


## Features

[dataset](docs/dataset) supports 

- Basic storage actions ([create](docs/create.html), [read](docs/read.html), [update](docs/update.html) and [delete](docs/delete.html))
- listing of collection [keys](docs/keys.html) (including filtering and sorting)
- import/export of [CSV](how-to/working-with-csv.html) files
- The ability to reshape data by performing simple object [joins](docs/join.html)
- The ability to create data [frames](docs/frame.html) from collections based on keys lists and [dot paths](docs/dotpath.html) into stored JSON objects

You can work with dataset collections via the 
[command line tool](docs/dataset.html), via Go using the 
[dataset package](https://godoc.org/github.com/caltechlibrary/dataset) 
or in Python 3.8 using the 
[py_dataset](https://github.com/caltechlibrary/py_dataset) python package.  _dataset_ is useful for general data science applications 
which need intermediate JSON object management but not 
a full blown database.


### Limitations of _dataset_

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store (it's files on "disc" without locking)
- it is not a replacement for a repository management system
- it is not a general purpose database system
- it does not supply automatic version control on collections or objects

## Read next ...

Explore _dataset_ through 
[A Shell Example](how-to/a-shell-example.html "command line example"),
[Getting Started with Dataset](how-to/getting-started-with-dataset.html "pyton examples as well as command line"),
[How To](how-to/) guides,
[topics](docs/topics.html) and [Documentation](docs/).

## Releases

Compiled versions are provided for Linux (x86), Mac OS X (x86 and M1), 
Windows 10 (x86) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

You can use _dataset_ from Python via the [py_dataset](https://github.com/caltechlibrary/py_dataset) package.
