dataset
=======

USAGE
-----

	dataset [OPTIONS] VERB [OPTIONS] COLLECTION_NAME [ACTION PARAMETERS...]

SYNOPSIS
--------


dataset is a command line tool demonstrating dataset package for 
managing JSON documents stored on disc. A dataset is organized 
around collections, collections contain a pairtree holding specific 
JSON documents and related content.  In addition to the JSON 
documents dataset maintains metadata for management of the 
documents, their attachments as well as a ability to generate 
select lists based JSON document keys (aka JSON document names).

OPTIONS
-------

Options can be general (e.g. `--help`) or specific to a verb.
General options are

```
    -e, -examples             display examples
    -h, -help                 display help
    -l, -license              display license
    -v, -version              display version
    -verbose                  output rows processed on importing from CSV
```


VERBS
-----

```
    attach         Attach a document (file) to a JSON record in a collection
    attachments    List of attachments associated with a JSON record in
                   a collection
    check          Check the health of a dataset collection (for collections
                   using pairtree storage model)
    clone          Clone a collection from a list of keys into a new
                   collection
    clone-sample   Clone a collection into a sample size based training
                   collection and test collection
    count          Counts the number of records in a collection, accepts
                   a filter for sub-counts
    create         Create a JSON record in a collection
    delete         Delete a JSON record (and attachments) from a collection
    delete-frame   remove a frame from a collection
    retrieve       Copy an attachment out of an associated JSON record in
                   a collection
    frame          define a frame in a collection
    frame-def      retrieve a frame's definition
    frame-objects  retrieve a frame's object list
    frame-keys     retrieve a frame's key list
    frames         list the available frames in a collection
    has-key        Returns true if key is in collection, false otherwise
    init           Initialize a dataset collection
    join           Join a JSON record with a new JSON object in a collection
    keys           List the keys in a collection
    prune          Remove attachments from a JSON record in a collection
    read           Read back a JSON record from a collection
    reframe        re-generate a frame with provided key list
    repair         Try to repair a damaged dataset collection (for
                   collections using a pairtree storage model)
    sample         return a random sample of keys in a collection
    update         Update a JSON record in a collection
```

