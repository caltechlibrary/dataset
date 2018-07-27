
# deindexer

This is an experimental feature of dataset and maybe removed 
in the future.

## Syntax

```
    dataset deindexer INDEX_NAME KEY_FILENAME
```

## Description

Deindexer removes records from an index. It requires a list of 
keys in a key file, one key per line.

## Usage

```
    dataset deindexer author-title.bleve titles-to-delete.txt
```

Related topics: [indexer](indexer.html), [find](find.html)

