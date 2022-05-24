
# Developer notes

## Requirements

(needed if compiling from source)

+ Golang 1.18.2 or better

## Recommend

(recommended if compiling from source)

- Bash
- GNU Make
- codemeta2cff (part of [datatools](https://github.com/caltechlibrary/datatools))
- Python 3.9
- Pandoc and [mkpage](https://github.com/caltechlibrary/mkpage) (a Pandoc pre-processor)
- Snapcraft if generating a snap package of dataset/datasetd

## Using the _dataset_ package

+ create/initialize collection
+ create a JSON document in a collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

```go
    package main
    
    import (
        "github.com/caltechlibrary/dataset"
        "log"
    )
    
    func main() {
        // Create a collection "mystuff" inside the directory called demo
        collection, err := dataset.InitCollection("demo/mystuff.ds", dataset.PTSTORE)
        if err != nil {
            log.Fatalf("%s", err)
        }
        defer collection.Close()
        // Create a JSON document
        docName := "freda.json"
        document := map[string]interface{}{
            "name":  "freda",
            "email": "freda@inverness.example.org",
        }
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
    }
```


## package requirements

_dataset_ is built on both Golang's standard packages and Caltech Library 
packages.

## Caltech Library packages

+ [github.com/caltechlibrary/dotpath](https://github.com/caltechlibrary/dotpath)
    + provides dot path style notation to reach into JSON objects

