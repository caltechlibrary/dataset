
# delete

## Syntax

```
    dataset COLLECTION_NAME delete KEY
```

## Description

+ delete - removes a JSON document from collection
  + requires JSON document name

## Usage

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

```shell
    dataset publications.ds delete r1
```

Related topics: create, read, and update

