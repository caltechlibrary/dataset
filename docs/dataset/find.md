
# find

## Syntax

```
    dataset [OPTIONS] INDEX_NAMES QUERY_STRING
```

## Description

_find_ adds support for full text searching of a collection based on [Bleve](https://blevesearch.com) indexes.
It supports a [query string language]() similar to elastic search. Additionally _find_ can render the results
in various formats include plain text, JSON, and CSV.

_find_ supports using multiple indexes. List the index names separated by colons.

## Usage

Single index examples

```
    dataset find authors-title.bleve 'Robert Doiel'
    dataset find authors-title.bleve '+family:"Doiel" given:"R"'
    dataset find authors-title.bleve '+orcid:"0000-0003-0900-6903"'
```

Multi index examples (using authors-title.bleve and abstracts.bleve indexes)

```
    dataset find 'authors-title.bleve:abstracts.bleve' 'Robert Doiel'
    dataset find 'authors-title.bleve:abstracts.bleve' '+family:"Doiel" given:"R"'
    dataset find 'authors-title.bleve:abstracts.bleve' '+orcid:"0000-0003-0900-6903"'
```
