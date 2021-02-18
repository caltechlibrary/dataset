
# count

## Syntax

```
    dataset count COLLECTION_NAME [FILTER EXPRESSION]
```

## Description

This returns a count of the keys in the collection. It is reasonable 
quick as only the collection metadata is read in. *count* also can 
accept a filter expression. This is slower as it iterates over all 
the records and counts those which evaluate to true based on the
filter expression provided.

## Usage

Count all records in collection "publications.ds"

```shell
    dataset count "publications.ds"
```

Count records where the ".published" field is true.

```shell
    dataset count "publications.ds" '(eq .published true)'
```

Related topic: [keys](keys.html)

