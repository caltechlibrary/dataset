
# dataset   [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

_dataset_ is a command line tool, Go package, shared library and Python package for working with [JSON](https://en.wikipedia.org/wiki/JSON) objects as collection. Collections can be stored on disc, in S3 or Google Cloud Storage.
JSON objects are stored in collections as plain UTF-8 text. This means the objects can be accessed with common Unix text processing tools as well as
most programming languages with text processing support. 

The [dataset](docs/dataset.html) command line tool supports common data manage operations such as initialization of collections, creation, reading, updating and deleting JSON objects in the collection. Some of its enhanced features include the ability to generate data [grids](docs/grid.html) and [frames](docs/frame.html), the ability to import and export JSON object to and from CSV files and Google Sheets.  It even includes an experimental search feature by the integrating [Blevesearch](http://www.blevesearch.com) indexing and search engine library developed for [CouchDB](http://couchdb.apache.org/).

In addition to the command line tool dataset includes a C shared library called libdataset which is used for integration in a Python module of the same name.  _dataset_ itself is written in a [Go](https://golang.org) package which can also in other Go based projects.  _libdataset_ could be used as a bases for integration with other languages that support a C API (e.g. [Julia](https://julialang.org/)).

See [getting-started-with-datataset.md](how-to/getting-started-with-dataset.html) for a tour and tutorial.


## Origin story

The inspiration for creating _dataset_ was the desire to process metadata as JSON object collections using simple Unix shell utilities and data pipelines. The core use case evolved at [Caltech Library](https://library.caltech.edu) working with various repository systems' API (e.g. [EPrints](https://en.wikipedia.org/wiki/EPrints) and and [Invenio](https://en.wikipedia.org/wiki/Invenio)). It has allowed the library to build an aggregated view of heterogeneous content (see https://feeds.library.caltech.edu) as well as facilitate ad-hoc analysis and data enhancement for a number of internal library projects.


## Design choices

_dataset_ isn't a database or repository system. It is intended to be simple and easier to use with minimal setup (e.g. `dataset init mycollection.ds` would create a new collection called 'mycollection.ds').  It is built around a few abstractions (e.g. dataset stores JSON objects in collections, collections are a folder containing a JSON file called collections.json and buckets, buckets containing the JSON objects and any attachments, the collections.json file describes the mapping of keys to buckets).  It takes minimal system resources
and keeps all content, except JSON object attachments, in plain UTF-8 text (attachments are kept in tar files).

A the typical library processing pattern is to write a "harvester" which stores it results in a _dataset_ collection, the use either a shell script or Python program to transform the collections content and finally redeploy the augmented results.

Care has been taken to keep _dataset_ simple enough and light weight enough that it will run on a machine as small as a Raspberry Pi while being equally comfortable on a more resource rich server or desktop environment.


## Features

[dataset](docs/dataset) supports 

- Basic storage actions ([create](docs/create.html), [read](docs/read.html), [update](docs/update.html) and [delete](docs/delete.html))
- listing of collection [keys](docs/keys.html) (including filtering and sorting)
- import/export  of [CSV](how-to/import-csv-rows-as-json-documents.html) files and [Google Sheets](how-to/gsheet-integration.html)
- An experimental full text [search](how-to/searchable-datasets.html) interface based on [Blevesearch](https://blevesearch.com)
- The ability to reshape data by performing simple object [joins](docs/join.html)
- The ability to create data [grids](docs/grid.html) and [frames](docs/frame.html) from collections based 
  on keys lists and [dot paths](docs/dotpath.html) into the JSON objects stored

You can work with dataset collections via the [command line tool](docs/dataset.html), via Go using the 
[dataset package](https://godoc.org/github.com/caltechlibrary/dataset) or in
Python 3.6 using a python package.  _dataset_ is useful for general data science applications which 
need intermediate JSON object storage but not a full blown database.


### Limitations of _dataset_

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store (it's just files on disc without any locking)
- it is not a repository management system
- it is not a general purpose multi-user database system
- it does not supply version control on collections or objects (though integrating it with git, mercurial or subversion would be trivial)


## Example

Below is a simple example of shell based interaction with dataset collations using the command line dataset tool.

```shell
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    bin/dataset friends.ds init
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    bin/dataset friends.ds create frieda '{"name":"frieda","email":"frieda@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    bin/dataset friends.ds read frieda
    
    # Path to JSON document
    bin/dataset friends.ds path frieda

    # Update a JSON document
    bin/dataset friends.ds update frieda '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    bin/dataset friends.ds keys

    # Get keys filtered for the name "frieda"
    bin/dataset friends.ds keys '(eq .name "frieda")'

    # Join frieda-profile.json with "frieda" adding unique key/value pairs
    bin/dataset friends.ds join append frieda frieda-profile.json

    # Join frieda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from frieda-profile.json
    bin/dataset friends.ds join overwrite frieda frieda-profile.json

    # Delete a JSON document
    bin/dataset friends.ds delete frieda

    # Import data from a CSV file using column 1 as key
    bin/dataset -quiet -nl=false friends.ds import-csv my-data.csv 1

    # To remove the collection just use the Unix shell command
    rm -fR friends.ds
```

## Releases

Compiled versions are provided for Linux (amd64), Mac OS X (amd64), Windows 10 (amd64) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

