
## Indexing a collection

In the example the index will be created for a collection called *characters.ds*.

```shell
    dataset characters.ds indexer email-mapping.json email-index
```

This will build a Bleve index called "email-index" based on the index defined
in "email-mapping.json".


## Searching an index

In this example we have already indexes a collection called "characters.ds". The
index name in *characters.bleve* which we will use for searching.

```shell
    dataset find characters.bleve "Jack Flanders"
```

This would search the Bleve index named *characters.bleve* for the string "Jack Flanders" 
returning records that matched based on how the index was defined.

