
# Getting started with dataset

_dataset_ is a set of tools for managing JSON (object) documents as a collection of key/value pairs stored on either your
local file system, AWS S3 or Google Cloud Storage. These documents can be interated over or retrieved individually.
There is also a full text indexer for supporting fielded or full text searches based on the index definitions.
One final feature of _dataset_ is the ability to add attachments to your JSON objects. These attachments are stored
in a simple archive format called tar. Basic metadata can be retrieved, and the attachments can be retreive as a group
or individually. Attachments can be removed.

## Getting dataset onto your computer

The command line _dataset_ is available for installation from https://github.com/caltechlibrary/dataset/releases/latest.
Find the zip file associated with your computer type and operating system then download it. Once downloaded you can unzip the zip
file and copy the programs into a local directory called "bin" on your comptuer. For full instructions on installation see
[INSTALL.md](../install.html).

## Basic workflow with dataset

_dataset_'s focus is in storing JSON (object) documents in collections. The documents are stored in a bucketed directory structure and
named for the "key" provided. The documents remain plain text JSON on disc. When you first start working with a dataset you
will need to initialize the collection. This creates the bucket directories and associated metadata so you can easily
retrieve your documents. If you were to initialize a dataset collection called "FavoriteThings.ds" it would look like --

```shell
    dataset init FavoriteThings.ds
```

If the command is successful you'll see output that looks like

```shell
    export DATASET="FavoriteThings.ds"
```

This is a suggested command to run your shell session. It sets the default DATASET to operate on. With out it
you need to explicit indicite which collection using the command line option `-c` or `-collection`. To save your
self some typing cut and paste the export statement now into your terminal session.

Next you'll want to add some records to the collection of "FavoriteThings.ds".  The records we're going to add need
to be expressed as JSON objects. You need to decide on a key (the thing you'll used to retrieve the record later)
of the document to store.  For this example I'm going to use the key, "beverage" and a document that looks like
`{"thing": "coffee"}`.  If you've set the DATASET environment variable you can run the following command --

```shell
    dataset create beverage '{"thing":"coffee"}'
```

If all goes well you'll get a response of "OK".  If you forgot to set the environment variable you can use the 
`-c` option to achieve the same thing.

```shell
    dataset -c FavoriteThings.ds create beverage '{"thing":"coffee"}'
```

Later if your have forgotten what your favorite beverage was you can read it back with

```shell
    dataset read beverage
```

Or using the `-c` option

```shell
    dataset -c FavoriteThings.ds read beverage
```

To list all your favorite things keys try

```shell
    dataset keys
```

or 
k
```shell
    dataset -c FavoriteThings.ds keys
```

## Adding an existing JSON document to a collection

One of my favorite things is music. I happen to have a JSON document that I started currating with
song and performers names. The document is called `music.json`. I can add this to my collection too.

```json
    {
       "songs": ["Blue Rondo al la Turk", "Larks Tongues in Aspic", "Bernie's Tune", "Perdido"],
       "performers": ["Dave Brubeck Quartet", "King Crimson", "Dirk Fischer", "L.A. Guitar Quartet"]
    }
```

I can add this to my collection of *FavoriteThings.ds* this way using the key "songs-performers". 

```shell
    dataset -c FavoriteThings.ds create "songs-performers" music.json
```

Notice that the organization of the JSON documents do not impose a common structure (though that is
often useful). We can list the documents using our key command.

```shell
    dataset -c FavoriteThings.ds keys
```

Would return something like

```
    beverage
    songs-performers
```

The should list out "beverage" and "songs-performers". 

I can create a JSON list of the objects stored using the "list" command.

```shell
    dataset -c FavoriteThings.ds list beverage songs-performers
```

Would return something like

```json
    [
        {
            "_Key": "beverage",
            "thing": "coffee"
        },
        {
            "_Key": "songs-performers",
            "performers": [
                "Dave Brubeck Quartet",
                "King Crimson",
                "Dirk Fischer",
                "L.A. Guitar Quartet"
            ],
            "songs": [
                "Blue Rondo al la Turk",
                "Larks Tongues in Aspic",
                "Bernie's Tune",
                "Perdido"
            ]
        }
    ]
```


