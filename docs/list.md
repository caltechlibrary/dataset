
# list

## Syntax

```
    dataset COLLECTION_NAME list KEYS
```

## Description

Write a JSON array to standard out (unless you've specific an 
alternative location with the "-output" option) for the provided 
KEYS.  If no ids are provided (or none are found) then an empty 
JSON array is return.

## Usage

An example we're assuming there is are JSON documents with a KEYS 
of "r1", "r2", and "r3". Our collection name is "data.ds"

```shell
    dataset data.ds list r1 r2 r3
```

If "r1" was '{"one":1}', "r2" was '{"two":2}' and "r3" was '{"three":3}' 
then the output would be

```json
    [{"one":1},{"two":2},{"three":3}]
```

Related topics: [keys](keys.html), [read](read.html)

