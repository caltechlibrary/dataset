
# count

## Syntax

```
    dataset count [FILTER EXPRESSION]
```

## Description

This returns a count of the keys in the collection. It is reasonable quick as only the
collection metadata is read in. *count* also can accept a filter expression. This is slower
as it iterates over all the records and counts those which evaluate to true based on the
filter expression provided.

## Usage

Count all records

```shell
    dataset count
```

Count records where the `.published` field is true.

```shell
    dataset count '(eq .published true)'
```

