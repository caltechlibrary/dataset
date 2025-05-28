
# Developer notes

## Requirements

(when compiling from source)

- Golang 1.24.3 or better

## Recommend

(recommended if compiling from source and for development)

- Bash
- GNU Make
- Pandoc (is used as a pre-processor to generate version.go and other files)

## Using the _dataset_ package

- create/initialize collection
- create a JSON document in a collection
- read a JSON document
- update a JSON document
- delete a JSON document

~~~golang
    package main
    
    import (
        "github.com/caltechlibrary/dataset/v3"
        "log"
    )
    
    func main() {
        // Create a collection "mystuff" inside the directory called demo
        collection, err := dataset.InitCollection("demo/mystuff.ds", "")
        if err != nil {
            log.Fatalf("%s", err)
        }
        defer collection.Close()
        // Create a JSON document
        key := "freda"
        document := map[string]interface{}{
            "name":  "freda",
            "email": "freda@inverness.example.org",
        }
        if err := collection.Create(key, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Read a JSON document
        if err := collection.Read(key, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Update a JSON document
        document["email"] = "freda@zbs.example.org"
        if err := collection.Update(key, document); err != nil {
            log.Fatalf("%s", err)
        }
        // Delete a JSON document
        if err := collection.Delete(key); err != nil {
            log.Fatalf("%s", err)
        }
    }
~~~


## package requirements

_dataset3_ is built on both Golang's standard packages and Caltech Library packages.

## Caltech Library packages

- [github.com/caltechlibrary/dotpath](https://github.com/caltechlibrary/dotpath)
  - provides dot path style notation to reach into JSON objects

