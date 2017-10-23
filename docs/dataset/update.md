
# update

## Syntax

```
    dataset update JSON_RECORD_ID
```

## Description

update will replace a JSON document in a dataset collection for a given JSON_RECORD_ID.
By default the JSON document is read from standard input but you can specific a spefic
file with the "-input" option. The JSON document should aready exist in the collection
when you use update.


## Usage

In this example we assume there is a JSON document on local disc named _blob.json_. It
contains `{"name":"Jane Doe"}` and the JSON_RECORD_ID is "jane.doe". In the first
one we specify the full JSON document via the command line after the JSON_RECORD_ID.
In the second example we read the data from _blob.json_. Finally in the last we
read the JSON document from standard input and save the update to "jane.doe".

```shell
    dataset update jane.doe '{"name":"Jane Doiel"}'
    dataset -i blob.json update jane.doe
    cat blob.json | dataset update jane.doe
```

