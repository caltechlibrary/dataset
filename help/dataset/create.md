
# create

## Syntax

```
    cat JSON_DOCNAME | dataset create JSON_RECORD_ID
    dataset -i JSON_DOCNAME create JSON_RECORD_ID
    dataset create JSON_RECORD_ID JSON_VALUE
```

## Description

create adds or replaces a JSON document to a collection. The JSON document can be read from a 
standard in, a named file or expressed on the command line.

## Usage

In the following three examples *blob.json* is a file on the local file system
contains JSON data containing the JSON_VALUE of `{"name":"Jane Doe"}`.  The JSON_RECORD_ID we will 
create is _r1_. 

```shell
    cat blob.json | dataset create r1
    dataset -i blob.json create r1
    dataset create r1 '{"name":"Jane Doe"}'
```

