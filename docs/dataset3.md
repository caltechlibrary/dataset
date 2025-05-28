dataset3
========

USAGE
-----

	dataset3 [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [ACTION PARAMETERS...]

SYNOPSIS
--------

dataset is a command line tool demonstrating dataset package for 
managing JSON documents stored on disc. A dataset is organized 
around collections, collections contain a pairtree holding specific 
JSON documents and related content.  In addition to the JSON 
documents dataset maintains metadata for management of the 
documents, their attachments as well as a ability to generate 
select lists based JSON document keys (aka JSON document names).

GLOBAL OPTIONS
--------------

Options can be general (e.g. `--help`) or specific to a verb.
General options are

~~~
    -e, -examples             display examples
    -h, -help                 display help
    -l, -license              display license
    -v, -version              display version
    -verbose                  output rows processed on importing from CSV
~~~


VERBS
-----

~~~
    create         Create a JSON record in a collection
    delete         Delete a JSON record (and attachments) from a collection
    dump           export a dataset collection as JSON lines stream
    haskey        Returns true if key is in collection, false otherwise
    init           Initialize a dataset collection
    keys           List the keys in a collection
    load           import a dataset collection using JSON lines stream
    read           Read back a JSON record from a collection
    update         Update a JSON record in a collection
    query          query a collection using SQL
~~~


