
# extract

This is an experimental feature of dataset and may be removed in the future.

## Syntax

```
    dataset COLLECTION_NAME extract FILTER DOTPATH
```

## Description

extract returns a list of unique values across documents in a collection based on the FILTER and
DOTPATH provided (for DOTPATH see `dataset -help dotpath` and FITLER see `dataset -help filter`).

## Usage

In this example we're turning a list of unique author orcid ids across the collection. The filter
we use is "true". The author field is an array so our dotpath uses that notation. Collection name
is "publications.ds"

```shell
    dataset publications.ds extract true .authors[:].orcid
```

Related topics: dotpath, export-csv and import-csv

