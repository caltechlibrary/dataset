
# Documentation for dataset

The documentation is organized around the command line options 
and as a series of "how to" style examples.

+ [getting started with dataset](../how-to/getting-started-with-dataset.html) (covers both Bash and Python)
+ Explore additional other [tutorials](../how-to/)

## Command line program documentation

+ [dataset](dataset.html) - usage page for managing collections with _dataset_

## Internal project concepts

+ [upgrading a collection](../how-to/upgrading-a-collection.html) - Describes how to upgrade a collection from a previous version of dataset to a new one
+ [how attachments work](../how-to/how-attachments-work.html) - Detailed descriptiojn of attachments and their metadata

## _dataset_ Operations

The basic operations support by *dataset* are listed below organized 
by collection and JSON document level.

### Collection Level

+ [init](init.html) creates a collection
+ [import](import-csv.html) (csv) JSON documents from rows of a CSV file
+ [export](export-csv.html) (csv) JSON documents from a collection into a CSV file
+ [keys](keys.html) list keys of JSON documents in a collection, supports filtering and sorting
+ [haskey](haskey.html) returns true if key is found in collection, false otherwise
+ [count](count.html) returns the number of documents in a collection, supports filtering for subsets
+ [data frame support](../how-to/collections-and-data-frames.html) provides a persistant metadata associated with the collection as data frames
    + [frame](frame.html)
    + [frame-objects](frame-objects.html)
    + [frame-grid](frame-grid.html)
    + [frames](frames.html)
    + [refresh](refresh.html)
    + [reframe](reframe.html)
    + [delete-frame](delete-frame.html)
    + [hasframe](hasframe.html)

### JSON Document level

+ [create](create.html) a JSON document in a collection
+ [read](read.html) back a JSON document in a collection
+ [update](update.html) a JSON document in a collection
+ [delete](delete.html) a JSON document in a collection
+ [join](join.html) a JSON document with a document in a collection
+ [list](list.html) the lists JSON records as an array for the supplied keys
+ [path](path.html) list the file path for a JSON document in a collection

### JSON Document Attachments

+ [attach](attach.html) a file to a JSON document in a collection
+ [attachments](attachments.html) lists the files attached to a JSON document in a collection
+ [detach](detach.html) retrieve an attached file associated with a JSON document in a collection
+ [prune](prune.html) delete one or more attached files of a JSON document in a collection

### Samples and cloning

+ [sample](sample.html) - getting a random sample of keys
+ [clone](clone.html) - clone a repository
+ [clone-sample](clone-sample.html) - cloning a respository into training and test collections

### Collection health

+ [check](check.html) - checks a collection against the current version of tools
+ [repair](repair.html) - repairs/upgrades a collection based on the current verison of the tool

