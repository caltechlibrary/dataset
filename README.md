
# dataset

A go package for managing JSON documents stored on disc. A dataset stores
one of more collections, collections store the a buckted distribution of documents
as well as minimal metadata about the collection.

## disc layout

+ dataset
    + collection
        + collection.json - metadata about collection
            + maps the filename of the JSON blob stored to a bucket in the collection
            + e.g. file "mydocs.jons" stored in bucket "aa" would have a map of {"mydocs.json": "aa"}
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
+ Record level
    + Create (record) - saves a new JSON blob to disc with given blob name (sets dirty flag on collection)
    + Read (record) - finds the bucket the record is in and returns the JSON blob
    + Update (record) - updates an existing blob on disc (sets dirty flag on collection)
    + Delete (record) - removes a JSON blob from its disc (sets the dirty flag on collection)

## Example

Common operations shown in Golang

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```go
    // Create a collection "mystuff" inside the directory called dataset
    collection, err := dataset.Create("dataset/mystuff", dataset.GenerateBucketNames("ab", 2))
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
