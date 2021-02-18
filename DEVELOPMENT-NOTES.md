
# Developer notes

## Requirements

Compiler requirements

+ go1.16 or better

## Recommend

+ Python v3.8 or better to use the C-Shared library libdataset

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
        collection, err := dataset.InitCollection("demo/mystuff.ds")
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

_dataset_ is built on both Golang's standard packages, Caltech Library 
packages and a few 3rd party packages.  At this has not been necessary 
to vendor any packages assuming you're building from the master branch.

## Caltech Library packages

+ [github.com/caltechlibrary/dotpath](https://github.com/caltechlibrary/dotpath)
    + provides dot path style notation to reach into JSON objects
+ [github.com/caltechlibrary/storage](github.com/caltechlibrary/storage)
    + provides a unified storage interaction supporting local disc and AWS S3 storage
+ [github.com/caltechlibrary/tmplfn](https://github.com/caltechlibrary/tmplfn)
    + provides additional template functionality used to format web search results
    + provides a filter engine leveraging the pipeline notation in Go's text templates


## 3rd party packages

+ [Markdown packages] - used to support rendering Markdown embedded in JSON objects
    + [github.com/microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday)
    + [github.com/russross/blackfriday](https://github.com/russross/blackfriday)

