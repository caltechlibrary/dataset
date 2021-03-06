
# haskey

## Syntax

```
    dataset [OPTIONS] haskey COLLECTION_NAME KEY_TO_CHECK_FOR
```

## Description

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

## Usage

```
    dataset haskey people.ds '0000-0003-0900-6903'
    dataset haskey people.ds r1
```

In python

```
    dataset.has_key('people.ds', '0000-0003-0900-6903')
    dataset.has_key('people.ds', 'r1')
```

Related topics: [keys](keys.html)

