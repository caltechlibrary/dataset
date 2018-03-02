
# extract

## Syntax

```
    dataset extract FILTER DOTPATH
```

## Description

extract returns a list of unique values across documents in a collection based on the FILTER and
DOTPATH provided (for DOTPATH see `dataset -help dotpath` and FITLER see `dataset -help filter).

## Usage

In this example we're turning a list of unique author orcid ids across the collection. The filter
we use is "true". The author field is an array so our dotpath uses that notation.

```shell
    dataset extract true .authors[:].orcid
```

Related topics: dotpath, export and import

