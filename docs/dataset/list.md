
# list

## Syntax

```
    dataset list JSON_RECORD_IDS
```

## Description

Write a JSON array to standard out (unless you've specific an alternative 
location with the "-output" option) for the provided JSON_RECORD_IDS.
If no ids are provided (or none are found) then an empty JSON array is return.

## Usage

An example we're assuming there is a JSON document with a JSON_RECORD_ID of "r1".

```shell
    dataset list r1
```

If "r1" was '{"one":1}' then the output would be

```json
    [{"one":1}]
```

