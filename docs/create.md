
# create

## Syntax

```
    cat JSON_DOCNAME | dataset COLLECTION_NAME create KEY
    dataset -i JSON_DOCNAME COLLECTION_NAME create KEY
    dataset COLLECTION_NAME create KEY JSON_VALUE
    dataset COLLECTION_NAME create KEY JSON_FILENAME
```

## Description

create adds or replaces a JSON document to a collection. The JSON document can be read from a 
standard in, a named file (with a ".json" file extension) or expressed literally on the command line.

## Usage

In the following four examples *jane-doe.json* is a file on the local file system
contains JSON data containing the JSON_VALUE of `{"name":"Jane Doe"}`.  The KEY we will 
create is _r1_. Collection is "people.ds".  The following are equivalent in resulting record.

```shell
    cat jane-doe.json | dataset people.ds create r1
    dataset -i blob.json people.ds create r1
    dataset people.ds create r1 jane-doe.json
    dataset people.ds create r1 '{"name":"Jane Doe"}'
```

Related topics: [update](update.html), [read](read.html), and [delete](delete.html)

