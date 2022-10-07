dataset
=======

The documentation is organized around the command line options 
and as a series of "how to" style examples.

- [getting started with dataset](../how-to/getting-started-with-dataset.html) (covers both Bash and Python)
- Explore additional other [tutorials](../how-to/)

Command line program documentation
----------------------------------

- [dataset](dataset.html) - usage page for managing collections with _dataset_

Internal project concepts
-------------------------

- [upgrading a collection](../how-to/upgrading-a-collection.html) - Describes how to upgrade a collection from a previous version of dataset to a new one
- [how attachments work](../how-to/how-attachments-work.html) - Detailed description of attachments and their metadata

__dataset__ Operations
----------------------

The basic operations support by *dataset* are listed below organized 
by collection and JSON document level.

A word about keys
-----------------

__dataset__ is based around the concept of key/value pairs where
the key is the unique identifier for an object stored (i.e. the 
value) in the collection. Each storage option supported by dataset
and its own issues around what things can be called. **Keys should be
lower case alpha numeric or underscore only.** E.g. the pairtree storage
relies on the file system to store the JSON objects. Some file
systems are not case sensitive, others face challenges with
non-alpha numeric filenames.


Collection Level
----------------

- [init](init.html) creates a collection
- [keys](keys.html) list keys of JSON documents in a collection, supports filtering and sorting
- [has-key](haskey.html) returns true if key is found in collection, false otherwise
- [count](count.html) returns the number of documents in a collection, supports filtering for subsets
- [data frame support](../how-to/collections-and-data-frames.html) provides a persistent metadata associated with the collection as data frames
    - [frame](frame.html)
    - [frame-objects](frame-objects.html)
    - [frames](frames.html)
    - [refresh](refresh.html)
    - [reframe](reframe.html)
    - [delete-frame](delete-frame.html)
    - [has-frame](hasframe.html)

JSON Document level
-------------------

- [create](create.html) a JSON document in a collection
- [read](read.html) back a JSON document in a collection
- [update](update.html) a JSON document in a collection
- [delete](delete.html) a JSON document in a collection

JSON Document Attachments
-------------------------

- [attach](attach.html) a file to a JSON document in a collection
- [attachments](attachments.html) lists the files attached to a JSON document in a collection
- [retrieve](retrieve.html) retrieve an attached file associated with a JSON document in a collection
- [prune](prune.html) delete one or more attached files of a JSON document in a collection

Samples and cloning
-------------------

- [sample](sample.html) - getting a random sample of keys
- [clone](clone.html) - clone a repository
- [clone-sample](clone-sample.html) - cloning a repository into training and test collections

Collection health
-----------------

The following commands are provided to support pairtrees and a limitted amount of backward campatiblity with v1 dataset collections.

- [check](check.html) - checks a collection against the current version of tools
- [repair](repair.html) - repairs/upgrades a collection based on the current version of the tool

datasetd as a web service
=========================

New as of version v2.x is a web service providing access to dataset
collections. This is described in the [datasetd](datasetd.html) 
documentation page.

[datasetd](datasetd.html) supports the following end points.

Storage engines
===============

In v2 dataset is starting to suport storing your JSON document in a SQL database. This is an experimental feature and likely to contain some surprises in the 2.0 series implementation. Currently three SQL databases can be used to store the JSON documents, SQLite 3 (used in dataset's test suites), MySQL 8 (used by some related projects but under tested), Postgres 14 (not tested yet).  See [storage engines](storage-engines.html) for more details.
