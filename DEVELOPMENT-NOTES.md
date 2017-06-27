
# Developer notes

## Using the _dataset_ package

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


## package requirements

_dataset_ is built on both Golang's standard packages, Caltech Library packages and a few 3rd party packages.
At this has not been necessary to vendor any packages assuming you're building from the master branch.

## Caltech Library packages

+ [github.com/caltechlibrary/dotpath](https://github.com/caltechlibrary/dotpath)
    + provides dot path style notation to reach into JSON objects
+ [github.com/caltechlibrary/storage](github.com/caltechlibrary/storage)
    + provides a unified storage interaction supporting local disc and AWS S3 storage
+ [github.com/caltechlibrary/tmplfn](https://github.com/caltechlibrary/tmplfn)
    + provides additional template functionality used to format web search results
    + provides a filter engine leveraging the pipeline notation in Go's text templates


## 3rd party packages

+ [bleve](https://blevesearch.com) - for indexing and search capabilities (e.g. _dsfind_ and _dsws_)
    + github.com/blevesearch/bleve
    + github.com/blevesearch/bleve
    + github.com/blevesearch/bleve/analysis/analyzer/keyword
    + github.com/blevesearch/bleve/analysis/analyzer/simple
    + github.com/blevesearch/bleve/analysis/analyzer/standard
    + github.com/blevesearch/bleve/analysis/analyzer/web
    + github.com/blevesearch/bleve/analysis/lang/ar
+ [aws sdk go](https://github.com/aws/aws-sdk-go) - supporting AWS S3 storage (used by all the cli)
    + github.com/aws/aws-sdk-go/aws
    + github.com/aws/aws-sdk-go/aws/session
    + github.com/aws/aws-sdk-go/service/s3
    + github.com/aws/aws-sdk-go/service/s3/s3iface
    + github.com/aws/aws-sdk-go/service/s3/s3manager
+ [Google UUID]() - for generating UUID when importing from CSV
    + [github.com/google/uuid](github.com/google/uuid)
+ [Markdown packages] - used to support rendering Markdown embedded in JSON objects
    + [github.com/microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday)
    + [github.com/russross/blackfriday](https://github.com/russross/blackfriday)
