
# Searchable Datasets

The _dataset_ tool provides _dataset indexer_ and _dataset find_. Together
they implement indexing and search for a dataset collection. The index and search features
are based built on [Bleve](https://www.blevesearch.com) search engine. Depending on how your define your
index(es) search can provide a effective means of exploring and aggregating your
collection (or collections).


## How to build an index

In the example the index will be created for a collection called *characters.ds*.

```shell
    dataset characters.ds indexer email-mapping.json email-index
```

This will build a Bleve index called "email-index" based on the index defined
in "email-mapping.json" (more on mapping indexes at [docs/defining-indexes.md](../docs/defining-indexes.html)).

You can build multiple indexes by having multiple index definitions. For large
JSON documents with lots of text this may let you more efficiently create the indexes.
Indexes and be aggregated together using _find_.


## Searching an index

In this example we have already indexes a collection called "characters.ds". The
index name in *characters.bleve* which we will use for searching.

```shell
    dataset find characters.bleve "Jack Flanders"
```

This would search the Bleve index named *characters.bleve* for the string "Jack Flanders" 
returning records that matched based on how the index was defined.

## How to search across multiple indexes

Let's say you have created an index called *audiodramas.bleve*. That index also includes
information about characters, scenes, etc.  If you want to search both *characters.bleve*
and *audiodramas.bleve* include both with your _find_ command

```shell
    dataset find characters.bleve audiodramas.bleve "Jack Flanders"
```

