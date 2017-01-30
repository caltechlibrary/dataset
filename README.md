
# dataset

A go package for managing JSON documents stored on disc. *dataset* is also a
command line tool. It stores one of more collections of JSON documents. Typically
you'd have a directory that holds collections, each collection holds buckets and 
each bucket holds some JSON documents. Both the package and command line tool 
allow you to interact with that logical structure on disc.

## layout

+ dataset (directory on file system)
    + collection (directory on file system)
        + collection.json - metadata about collection
            + maps the filename of the JSON blob stored to a bucket in the collection
            + e.g. file "mydocs.jons" stored in bucket "aa" would have a map of {"mydocs.json": "aa"}
        + keys.json - a list of keys in the collection (it is the default select list)
        + BUCKETS - a sequence of alphabet names for buckets holding JSON documents
            + Buckets let supporting common commands like ls, tree, etc. when the doc count is high
        + SELECT_LIST.json - a JSON document holding an array of keys
            + the default select list is "keys", it is not mutable by Push, Pop, Shift and Unshift
            + select lists cannot be named "keys" or "collection"

BUCKETS are names without meaning normally using Alphabetic characters. A dataset defined with four buckets
might looks like aa, ab, ba, bb.



## operations

+ Collection level 
    + Create (collection) - creates or opens collection structure on disc, creates collection.json and keys.json if new
    + Open (collection) - opens an existing collections and reads collection.json into memory
    + Close (collection) - writes changes to collection.json to disc if dirty
    + Delete (collection) - removes a collection from disc
    + Keys (collection) - list of keys in the collection
    + Select (collection) - returns the request select list, will create the list if not exist, append keys if provided
    + Clear (collection) - Removes a select list from a collection and disc
    + Lists (collection) - returns the names of the available select lists
+ JSON document level
    + Create (JSON document) - saves a new JSON blob or overwrites and existing one on  disc with given blob name, updates keys.json if needed
    + Read (JSON document)) - finds the JSON document in the buckets and returns the JSON document contents
    + Update (JSON document) - updates an existing blob on disc (record must already exist)
    + Delete (JSON document) - removes a JSON blob from its disc
    + Path (JSON document) - returns the path to the JSON document
+ Select list level
    + First (select list) - returns the value of the first key in the select list (non-distructively)
    + Last (select list) - returns the value of the last key in the select list (non-distructively)
    + Rest (select list) - returns values of all keys in the select list except the first (non-destructively)
    + List (select list) - returns values of all keys in the select list (non-destructively)
    + Length (select list) - returns the number of keys in a select list
    + Push (select list) - appends one or more keys to an existing select list
    + Pop (select list) - returns the last key in select list and removes it
    + Unshift (select list) - inserts one or more new keys at the beginning of the select list
    + Shift (select list) - returns the first key in a select list and removes it
    + Sort (select list) - orders the select lists' keys in ascending or descending alphabetical order
    + Reverse (select list) - flips the order of the keys in the select list

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



