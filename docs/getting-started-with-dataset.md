
# Getting started with dataset

_dataset_ is a set of tools to manage JSON objects as a collection of key/value pairs stored on either your
local file system, AWS S3 or Google Cloud Storage. These documents can be interated over or retrieved individually.
There is also a full text indexer for supporting fielded or full text searches based on the index definitions.
One final feature of _dataset_ is the ability to add attachments to your JSON objects. These attachments are stored
in a simple archive format called tar. Basic metadata can be retrieved, and the attachments can be retreive as a group
or individually.

## Getting dataset onto your computer

The command line tools that form dataset are available for installation from https://github.com/caltechlibrary/dataset/releases/latest.
Find the zip file associated with your computer type and operating systems and download it. Once downloaded you can unzip the zip
file and copy the programs into are directory in a local "bin" directory of your comptuer. For full instructions on installation see
[INSTALL.md](../install.html).

## Basic workflow with dataset

_dataset_'s focus is in storing JSON documents in collections. The documents are stored in a bucketed directory structure and
named for the "key" provided but themselves just remain plain JSON on disc. When you first start working with a dataset you
will need to initialize the collection. This creates the bucket directories and associated metadata so you can easily
retrieve your documents. If you were to initialize a dataset collection called "FavoriteThings" it would look like --

```shell
    dataset init FavoriteThings
```

If the command is successful you'll see output that looks like

```shell
    export DATASET=FavoriteThings
```

This is a suggested command to run your shell session. It sets the default DATASET to operate on. With out it
you'd have to explicit indicite which collection to use with the `-c` or `-collection` option. To save your
self some typing cut and paste the export statement now into your terminal session.

Next you'll want to add some records to the collection of "FavoriteThings".  The records we're going to add need
to be expressed as JSON objects. You need to decide on a key (the thing you'll used to retrieve the record later)
and document to store.  For this example I'm going to use the key, "beverage" and a document that looks like
`{"thing": "coffe"}`.  If you've set the DATASET environment variable you can run the following command --

```shell
    dataset create beverage '{"thing":"coffee"}'
```

If all goes well you'll get a response of "OK".  If you forgot to set the environment variable you can use the 
`-c` option to achieve the same thing.

```shell
    dataset -c FavoriteThings create beverage '{"thing":"coffee"}'
```

Later if your have forgotten what your favorite beverage was you can read it back with

```shell
    dataset read beverage
```

Or using the `-c` option

```shell
    dataset -c FavoriteThings read beverage
```

To list all your favorite things keys try

```shell
    dataset keys
```

or 

```shell
    dataset -c FavoriteThings keys
```





