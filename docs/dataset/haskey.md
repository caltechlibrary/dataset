
# haskey

## Syntax

```
    dataset [OPTIONS] COLLECTION_NAME haskey KEY_TO_CHECK_FOR
```

## Description

Checks if a given key is in the a collection. Returns "true" if found, "false" otherwise.
The collection name is "people.ds"

## Usage

```
    dataset people.ds haskey '0000-0003-0900-6903'
    dataset people.ds haskey r1
```

