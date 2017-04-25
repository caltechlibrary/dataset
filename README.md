
# dataset

_dataset_ is a go package for managing JSON documents stored on disc along with attached documents. 
*dataset* is also a command line tool exercising the features of the _dataset_ package.
It organanizes JSON documents by unique name in collections distributing the content across 
subdirectories (buckets) that allow easy process with common Unix text utilities. It also is
friendly for scripting languages like Bash (available on many operating systems such as Linux, Mac OS X,
Windows 10).

## Operations

A project goal of _dataset_ is to "play nice" with shell scripts and other Unix tools (e.g. it 
respects standard in, out and error with minimal side effects).  The operations support by *dataset* 
command line tools are listed below organized at the collection level, JSON document level.

### Collection Level

+ Create a collection
+ List the JSON document ids in a collection
+ Create named lists of JSON document ids (aka select lists)
+ Read back a named list of JSON document ids
+ Delete a named list of JSON document ids
+ Import JSON documents from rows of a CSV file

### JSON Document level

+ Create a JSON document in a collection
+ Update a JSON document in a collection
+ Read back a JSON document in a collection
+ Delete a JSON document in a collection

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

    # Delete a JSON document
    dataset delete freda.json

    # To remove the collection just use the Unix shell command
    # /bin/rm -fR demo/mystuff
```

Common operations shown in Golang

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```go
    // Create a collection "mystuff" inside the directory called demo
    collection, err := dataset.Create("demo/mystuff", dataset.GenerateBucketNames("ab", 2))
    if err != nil {
        log.Fatalf("%s", err)
    }
    defer collection.Close()
    // Create a JSON document 
    docName := "freda.json"
    document := map[string]string{"name":"freda","email":"freda@inverness.example.org"}
    if err := collection.Create(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Read a JSON document
    if err := collection.Read(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Update a JSON document
    document["email"] = "freda@zbs.example.org"
    if err := collection.Update(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Delete a JSON document
    if err := collection.Delete(docName); err != nil {
        log.Fatalf("%s", err)
    }
```

## Releases

Compiled versions are provided for Linux (amd64), Mac OS X (amd64), Windows 10 (amd64) and Raspbian (ARM7). 
See https://github.com/caltechlibrary/dataset/releases.

