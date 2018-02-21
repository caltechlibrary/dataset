
# dataset   [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

_dataset_ is a small collection of command line tools for working with JSON (object) documents stored as 
collections.  [This](docs/dataset/) include basic storage actions (e.g. CRUD operations, filtering
and extraction) as well as [indexing](docs/dsindexer/), [searching](docs/dsfind/) and even 
[web hosting](docs/dsws/).  A project goal of _dataset_ is to "play nice" with shell scripts and other 
Unix tools (e.g. it respects standard in, out and error with minimal side effects). This means it is 
easily scriptable via Bash, Posix shell or interpretted languages like Python.

_dataset_ is also a golang package for managing JSON documents and their attachments on disc or in cloud storage
(e.g. Amazon S3, Google Cloud Storage). The command line utilities excersize this package extensively.

The inspiration for creating _dataset_ was the desire to process metadata as JSON document collections using
Unix shell utilities and pipe lines. While it has grown in capabilities that remains a core use case.

_dataset_ organanizes JSON documents by unique names in collections. Collections are represented
as an index into a series of buckets. The buckets are subdirectories (or paths under cloud storage services) 
holding individual JSON documents and their attachments. The JSON documents in a collection as assigned to a
bucket (and the bucket generated if necessary) automatically when the document is added to the collection.
The assigment to the buckets is round robin determined by the order of addition. This avoids having too
many documents assigned to a single path (e.g. on some Unix there is a limit to how many documents are held
in a single directory). This means you can list and manipulate the JSON documents directly with common
Unix commands like ls, find, grep or their cloud counter parts.


### Limitations of _dataset_

_dataset_ has many limitations, some are listed below

+ it is not a real-time data store
+ it is not a repository management system
+ it is not a general purpose multiuser database system


## Operations

The basic operations support by *dataset* are listed below organized by collection and JSON document level.

### Collection Level

+ Create a collection
+ List the JSON document ids in a collection
+ Create named lists of JSON document ids (aka select lists)
+ Read back a named list of JSON document ids
+ Delete a named list of JSON document ids
+ Import JSON documents from rows of a CSV file or Google Sheets
+ Filter JSON documents and return a list of matching ids
+ Extract Unique JSON attribute values from a collection

### JSON Document level

+ Create a JSON document in a collection
+ Update a JSON document in a collection
+ Read back a JSON document in a collection
+ Delete a JSON document in a collection
+ Join a JSON document with a document in a collection

Additionally

+ Attach a file to a JSON document in a collection
+ List the files attached to a JSON document in a collection
+ Update a file attached to a JSON document in a collection
+ Delete one or more attached files of a JSON document in a collection

## Examples

Common operations using the *dataset* command line tool

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```shell
    # Create a collection "mystuff" inside the directory called demo
    dataset init demo/mystuff
    # if successful an expression to export the collection name is show
    export DATASET=demo/mystuff

    # Create a JSON document 
    dataset create freda.json '{"name":"freda","email":"freda@inverness.example.org"}'
    # If successful then you should see an OK or an error message

    # Read a JSON document
    dataset read freda.json

    # Path to JSON document
    dataset path freda.json

    # Update a JSON document
    dataset update freda.json '{"name":"freda","email":"freda@zbs.example.org"}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys

    # Filter for the name "freda"
    dataset filter '(eq .name "freda")'

    # Join freda-profile.json with "freda" adding unique key/value pairs
    dataset join update freda freda-profile.json

    # Join freda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from freda-profile.json
    dataset join overwrite freda freda-profile.json

    # Delete a JSON document
    dataset delete freda.json

    # To remove the collection just use the Unix shell command
    # /bin/rm -fR demo/mystuff
```

## Releases

Compiled versions are provided for Linux (amd64), Mac OS X (amd64), Windows 10 (amd64) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

