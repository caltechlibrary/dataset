
# dataset   [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

_dataset_ is a command line tool, Go package, and an experimental C shared 
library for working with [JSON](https://en.wikipedia.org/wiki/JSON) 
objects as collections. Collections can be stored on disc or in 
Cloud Storage.  JSON objects are stored in collections as 
plain UTF-8 text. This means the objects can be accessed with common 
Unix text processing tools as well as most programming languages.

The [dataset](docs/dataset.html) command line tool supports common data 
manage operations such as initialization of collections, creation, 
reading, updating and deleting JSON objects in the collection. Some of 
its enhanced features include the ability to generate data 
[frames](docs/frame.html) as well as the ability to 
import and export JSON objects to and from CSV files and Google Sheets.

_dataset_ is written in the [Go](https://golang.org) programming language.
It can be used as a Go package by other Go based software. Go supports
generating C shared libraries. By compiling the Go source you can
create a _libdataset_ C shared library. The C shared library is currently
being used by the DLD Group in Caltech Library experimentally from
Python 3.  This approach looks promising to support other languages 
(e.g. [Julia](https://julialang.org/) can easily use dataset via its 
ccall function, while R, Octave and NodeJS would probably need some 
C++ wrapping code).


See [getting-started-with-datataset.md](how-to/getting-started-with-dataset.html) for a tour and tutorial.

## Design choices

_dataset_ isn't a database or a replacement for repository systems. 
It is guided by the idea that you should be able to work with text 
files, the JSON objects documents, with standard Unix text utilities.
It is intended to be simple to use with minimal setup (e.g. 
`dataset init mycollection.ds` would create a new collection called 
'mycollection.ds'). It is built around a few abstractions --
dataset stores JSON objects in collections, collections are a folder(s) 
containing the JSON object documents and any attachments, a 
collections.json file describes the mapping of keys to folder locations).
_dataset_ takes minimal system resources and keeps all content, 
except JSON object attachments, in plain UTF-8 text. Attachments
are stored using the venerable "tar" archive format. 

The choice of plain UTF-8 and tar balls is intended to help future 
proof reading dataset collections.  Care has been taken to keep 
_dataset_ simple enough and light weight enough that it will run 
on a machine as small as a Raspberry Pi while being equally 
comfortable on a more resource rich server or desktop 
environment. It should be easy to do alternative implementations
in any language that has good string, JSON support and memory
management.


## Workflows

A typical library processing pattern is to write a "harvester" 
which then stores it results in a _dataset_ collection. Write something
that transforms or aggregates harvested options and then write
a final rendering program to prepare the data for the web. The
the hearvesters are typically written in Python or as a simple Bash
script storing the results in dataset. Depending on the performance
needs our transform and aggregates stage are written either in Python
or Go and our final rendering stages are typically written in Python or
as simple Bash scripts.


## Features

[dataset](docs/dataset) supports 

- Basic storage actions ([create](docs/create.html), [read](docs/read.html), [update](docs/update.html) and [delete](docs/delete.html))
- listing of collection [keys](docs/keys.html) (including filtering and sorting)
- import/export  of [CSV](how-to/working-with-csv.html) files and [Google Sheets](how-to/working-with-gsheets.html)
- An experimental full text [search](how-to/indexing-and-search.html) interface based on [Blevesearch](https://blevesearch.com)
- The ability to reshape data by performing simple object [joins](docs/join.html)
- The ability to create data [grids](docs/grid.html) and [frames](docs/frame.html) from collections based 
  on keys lists and [dot paths](docs/dotpath.html) into the JSON objects stored

You can work with dataset collections via the 
[command line tool](docs/dataset.html), via Go using the 
[dataset package](https://godoc.org/github.com/caltechlibrary/dataset) 
or in Python 3.7 using a python package.  _dataset_ is useful for general 
data science applications which need intermediate JSON object management 
but not a full blown database.


### Limitations of _dataset_

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store (it's files on "disc" without locking)
- it is not a replacement for a repository management system
- it is not a general purpose database system
- it does not supply version control on collections or objects

## Read next ...

Explore _dataset_ through 
[A Shell Example](how-to/a-shell-example.html),
[Getting Started with Dataset](how-to/getting-started-with-dataset.html),
[How To](how-to/) guides,
[topics](docs/topics.html) and [Documentation](docs/).

## Releases

Compiled versions are provided for Linux (amd64), Mac OS X (amd64), 
Windows 10 (amd64) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

