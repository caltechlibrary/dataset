
# dataset   [![DOI](https://data.caltech.edu/badge/79394591.svg)](https://data.caltech.edu/badge/latestdoi/79394591)

_dataset_ is a command line tool for working with JSON (object) documents stored as 
collections.  [This](docs/dataset/) supports basic storage actions (e.g. CRUD operations, filtering
and extraction) as well as [indexing](docs/dataset/indexer.html), [searching](docs/dataset/find.html).
A project goal of _dataset_ is to "play nice" with shell scripts and other 
Unix tools (e.g. it respects standard in, out and error with minimal side effects). This means it is 
easily scriptable via Bash, Posix shell or interpretted languages like R.

_dataset_ includes an implementation as a Python3 module. The same functionality as in the command line tool is 
replicated for Python3. (module requires Python 3.6 or better).

Finally _dataset_ is a golang package for managing JSON documents and their attachments on disc or in cloud storage
(e.g. Amazon S3, Google Cloud Storage). The command line utilities excersize this package extensively.

The inspiration for creating _dataset_ was the desire to process metadata as JSON document collections using
Unix shell utilities and pipe lines. While it has grown in capabilities that remains a core use case.

_dataset_ organanizes JSON documents by unique names in collections. Collections are represented
as an index into a series of buckets. The buckets are subdirectories (or paths under cloud storage services). 
Buckets hold individual JSON documents and their attachments. The JSON document is assigned automatically to a
bucket (and the bucket generated if necessary) when it is added to a collection. 
Assigning documents to buckets avoids having too many documents assigned to a single path (e.g. on some Unix 
there is a limit to how many documents are held in a single directory). In addition to using the _dataset_ 
comnad you can list and manipulate the JSON documents directly with common Unix commands like ls, find, grep or 
their cloud counter parts.

See [getting-started-with-datataset.md](docs/getting-started-with-dataset.html) for a tour of functionality.

### Limitations of _dataset_

_dataset_ has many limitations, some are listed below

+ it is not a multi-process, multi-user data store (it's just files on disc)
+ it is not a repository management system
+ it is not a general purpose multiuser database system


## Operations

The basic operations support by *dataset* are listed below organized by collection and JSON document level.

### Collection Level

+ [init](docs/dataset/init.html) creates a collection
+ [import-csv](docs/dataset/import-csv.html) JSON documents from rows of a CSV file
+ [import-gsheet](docs/dataset/import-gsheet.html) JSON documents from rows of a Google Sheet
+ [export-csv](docs/dataset/export-csv.html) JSON documents from a collection into a CSV file
+ [export-gsheet](docs/dataset/export-gsheet.html) JSON documents from a collection into a Google Sheet
+ [keys](docs/dataset/keys.html) list keys of JSON documents in a collection, supports filtering and sorting
+ [haskey](docs/dataset/haskey.html) returns true if key is found in collection, false otherwise
+ [count](docs/dataset/count.html) returns the number of documents in a collection, supports filtering for subsets
+ [extract](docs/dataset/extract.html) unique JSON attribute values from a collection


### JSON Document level

+ [create](docs/dataset/create.html) a JSON document in a collection
+ [read](docs/dataset/read.html) back a JSON document in a collection
+ [update](docs/dataset/update.html) a JSON document in a collection
+ [delete](docs/dataset/delete.html) a JSON document in a collection
+ [join](docs/dataset/join.html) a JSON document with a document in a collection
+ [list](docs/dataset/list.html) the lists JSON records as an array for the supplied keys
+ [path](docs/dataset/path.html) list the file path for a JSON document in a collection


### JSON Document Attachments

+ [attach](docs/dataset/attach.html) a file to a JSON document in a collection
+ [attachments](docs/dataset/attachments.html) lists the files attached to a JSON document in a collection
+ [detach](docs/dataset/detach.html) retrieve an attached file associated with a JSON document in a collection
+ [prune](docs/dataset/prune.html) delete one or more attached files of a JSON document in a collection

### Search

+ [indexer](docs/dataset/indexer.html) indexes JSON documents in a collection for searching with _find_
+ [deindexer](docs/dataset/deindexer.html) de-indexes (removes) JSON documents from an index
+ [find](docs/dataset/find.html) provides a index based full text search interface for collections


## Example

Common operations using the *dataset* command line tool

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```shell
    # Create a collection "mystuff.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    bin/dataset mystuff.ds init
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    bin/dataset mystuff.ds create freda '{"name":"freda","email":"freda@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    bin/dataset mystuff.ds read freda
    
    # Path to JSON document
    bin/dataset mystuff.ds path freda

    # Update a JSON document
    bin/dataset mystuff.ds update freda '{"name":"freda","email":"freda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    bin/dataset mystuff.ds keys

    # Get keys filtered for the name "freda"
    bin/dataset mystuff.ds keys '(eq .name "freda")'

    # Join freda-profile.json with "freda" adding unique key/value pairs
    bin/dataset mystuff.ds join append freda freda-profile.json

    # Join freda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from freda-profile.json
    bin/dataset mystuff.ds join overwrite freda freda-profile.json

    # Delete a JSON document
    bin/dataset mystuff.ds delete freda

    # Import data from a CSV file using column 1 as key
    bin/dataset -quiet -nl=false mystuff.ds import-csv my-data.csv 1

    # To remove the collection just use the Unix shell command
    rm -fR mystuff.ds
```

## Releases

Compiled versions are provided for Linux (amd64), Mac OS X (amd64), Windows 10 (amd64) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

