COLLECTIONS AND FRAMES
======================

__dataset__ stores JSON objects and can store also data frames similar to
that used in Python, R and Julia. This document outlines the ideas
behind __dataset__\'s implementation of data frames.

COLLECTIONS
-----------

Collections are at the core of the __dataset__ tool. A collection is either a
pairtree directory structure storing JSON objects in plain text
or a SQL table with a JSON column for the object. Both
support optional attachments. The root folder should contain two files
*collection.json* and *codemeta.json*. 
*collection.json* file contains operational metadata, metadata that dataset
needs to work with the collection of JSON documents. The *codemeta.json*
file is general metadata describing the collection. See [codemeta](https://codemete.github.io) for details on the *codemeta.json* file structure.

For the pairtree based collections one of the guiding ideas
behind has been to keep everything in plain text (i.e. UTF-8)
whenever reasonable. The affords the opportunity to easily interact
with the data via standard Unix tools or easily from any language
that supports working with JSON documents. A pairtree to a large
degree is future proof as it is very likely the ability to read the
file system and text files will persist long after dataset development
stops.

Over the course of the last several years the limitations of storing
JSON documents directly on the file system have been more compelling.
Additionally common, mature SQL database systems have pickup the
ability to store JSON as a column type (e.g. current systems
include SQLite3, MySQL 8, Postgres 14). With version 2 of dataset
this ability has been embraced more directly.  Collections of 
the JSON documents can be stored in a table in a supported database.
While this is less future proof than plain text it does increase
the flexibility of working with large collections, likewise a
dataset collection running with a SQL storage engine can be cloned
to a pairtree storage collection for long term preservation.

The dataset project provides Go package for working
with dataset collections. This means you can use dataset as a storage
engine directly in your own Go based projects. As of version 2
this is documented along with command line interactions.

Dataset isn\'t a database (there are plenty of JSON oriented databases
out there, e.g. CouchDB, MongoDB and No SQL storage systems for MySQL
and Postgresql). __dataset__\'s focus is on providing a mechanism to
manage JSON documents (objects). It supports the ability to attach non-JSON documents to a JSON document record as well as for working with
JSON collections as data frames. Dataset collections are like a mini
repository system avoiding the complexity of more mature repostiory
back ends like [Fedora Repository](https://duraspace.org/fedora/).

By working with JSON documents dataset can be used to feed full text
search engines like [Solr](https://solr.apache.org/),
[OpenSearch](https://www.opensearch.org/), and small engines like
[LunrJS](https://lunrjs.com/). Likewise when combined with Go's struct
types can be used to support building structured data repositories
customized to specific needs.


DATA FRAMES
-----------

Working with subsets of data in a collection is useful, particularly
ordered subsets. Implementing this started me thinking about the
similarity to data frames in Python, Julia and Octave. A *frame* is an
ordered list of objects.  Frames can be retrieved as a list of objects 
or as a list of keys.  Frames contain a additional metadata to help them
persist.  Frames include enough metadata to efficiently refresh framed
objects or even replace all objects in the list based on a new set of
keys. 

__dataset__ stores frames with the collection so they are is available for
later processing. The objects in a frame reflect the objects as they
existed when the frame was generated. Frames can be refreshed
to match the current state of the collection.

Frames become handy when moving data from JSON documents (tree like) to
other formats like spreadsheets (table like). This is because the
data frame's structure is defined based on paths into objects in the
collections. These pathes are mapped to "labels" structure the framed
objects. Frames can be used to simplify a complex record into a
simpler model for indexing in your favorite search engine.

Frames are stored in the collection's `_frames` sub directory. One
JSON document per frame combining both the definition and frame content.


FRAME OPERATIONS
----------------

-   frame (define a frame)
-   frame-def (read a frame definition, i.e. name, dot paths and labels)
-   frame-objects (read a frame's object list)
-   frame-keys (read a frame's key list)
-   frames (return a list of frame names)
-   reframe (replace all frame objects with objects indicated by a new list of keys)
-   refresh (update objects in a frame while pruning objects no longer
    in the collection)
-   has-frame (check to see if a frame exists in the collection)
-   delete-frame


### Create a frame

Example creating a frame named \"dois-and-titles\"

```shell
    dataset keys Pubs.ds >pubs.keys
    dataset frame-create -i pubs.keys Pubs.ds dois-and-titles \
        ".doi=DOI" \
        ".title=Title"
```

### Retrieve an existing frame's objects

An example of getting the frame\'s object list only.

```shell
    dataset frame-objects Pubs.ds dois-and-titles
```

### Regenerating a frame

Regenerating \"dois-and-titles\".

```shell
    dataset refresh Pubs.ds dois-and-titles
```

### Updating keys associated with the frame

In this example we want to "reframe" our "titles-and-dois"
data frame. We get the current list of keys in the collection
and regenerate the objects in the data frame using the new
list of keys.

```shell
    dataset Pubs.ds keys >updated.keys
    dataset reframe -i updated.keys Pubs.ds titles-and-dios
```

### Removing a frame

```shell
    dataset delete-frame Pubs.ds titles-and-dios
```

Listing available frames
------------------------

```shell
    dataset frames Pubs.ds
```

