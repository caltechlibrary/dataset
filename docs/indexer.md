
# indexer

This is an experimental feature of dataset and maybe removed in the future.

## Syntax

```
    dataset [OPTIONS] COLLECTION_NAME indexer INDEX_NAME INDEX_MAP_FILENAME
```

## Description

Indexer creates a [Blevesearch](https://blevesearch.com) index with used used by the [find](find.html)
command for searching a collection. The indexes support an elastic search like query language so
in additional to general full text support it also includes the ability to scope results by fields
defined in the index. The Bleve search packages supports full text search across many languages,
time ranges and geo points.

Index names should end in `.bleve` and the INDEX_MAP_FILE is a JSON file organized as described in
[defining-indexes](../defining-indexes.html) document.

## Usage

Our collection name in this example is "publications.ds"

```
    dataset publications.ds indexer author-title.bleve author-title.json
```

Related topics: [dotpath](dotpath.html), [find](find.html)

