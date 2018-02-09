
## Description

dsindexer is a command line tool for creating a Bleve index based on records in a dataset 
collection. dsindexer reads a JSON document for the index definition and uses that to
configure the Bleve index built based on the dataset collection. If an index
name is not provided then the index name will be the same as the definition file
with the .json replaced by .bleve.

A index definition is JSON document where the indexable record is defined
along with dot paths into the JSON collection record being indexed.

If your collection has records that look like

```json
    {"name":"Frieda Kahlo","occupation":"artist","id":"Frida_Kahlo","dob":"1907-07-06"}
```

and your wanted an index of names and occupation then your index definition file could
look like

```json
    {
        "types": {
            "default": {
                "enabled": true,
                "dynamic": true,
                "fields": [
                    {
                        "name": "name",
                        "type": "text",
                        "analyzer": "standard",
                        "store": true,
                        "index": true
                    },
                    {
                        "name": "occupation",
                        "type": "text",
                        "analyzer": "standard",
                        "store": true,
                        "index": true
                    }
                ]
            }
        }
    }
```

Based on this definition the "id" and "dob" fields would not be included in the index.

