
# dataset

A go package for managing JSON documents stored on disc. *dataset* is also a
command line tool. It stores one of more collections of JSON documents. Typically
you'd have a directory that holds collections, each collection holds buckets and 
each bucket holds some JSON document. Both the package and command line tool 
allow you to interact with that logical structure on disc.

## layout

+ dataset
    + collection
        + collection.json - metadata about collection
            + maps the filename of the JSON blob stored to a bucket in the collection
            + e.g. file "mydocs.jons" stored in bucket "aa" would have a map of {"mydocs.json": "aa"}
        + keys.json - a list of keys in the collection
        + BUCKETS - a sequence of alphabet names for buckets holding JSON documents
            + Buckets let supporting common commands like ls, tree, etc. when the doc count is high

BUCKETS are names without meaning normally using Alphabetic characters. A dataset defined with four buckets
might looks like aa, ab, ba, bb.

## operations

+ Collection level 
    + Create (collection) - sets up a new disc scripture and creates $DATASET/$COLLECTION_NAME/collection.json
    + Open (collection) - opens an existing collections and reads collection.json into memory
    + Close (collection) - writes changes to collection.json to disc if dirty
    + Delete (collection) - removes a collection from disc
    + Keys (collection) - list of keys in the collection
+ Record level
    + Create (record) - saves a new JSON blob to disc with given blob name (sets dirty flag on collection)
    + Read (record) - finds the bucket the record is in and returns the JSON blob
    + Update (record) - updates an existing blob on disc (sets dirty flag on collection)
    + Delete (record) - removes a JSON blob from its disc (sets the dirty flag on collection)
    + Path (record) - returns the path to the JSON document

## Example

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
    export DATASET_COLLECTION=demo/mystuff

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
