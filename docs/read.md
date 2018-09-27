
# read

## Syntax

```
    dataset read COLLECTION_NAME KEY
```

## Description

The writes the JSON document to standard out (unless you've 
specific an alternative location with the "-output" option)
for the given KEY.

## Usage

An example we're assuming there is a JSON document with a KEY 
of "r1". Our collection name is "data.ds"

```shell
    dataset read data.ds r1
```

Related topics: [keys](keys.html), [create](create.html), [update](update.html), [delete](delete.html)

