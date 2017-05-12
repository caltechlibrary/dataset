
# dsindexer

## USAGE

    dsindexer [OPTIONS] INDEX_DEFINITION [INDEX_NAME]

## SYNOPSIS

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
	   "name":{
		   "object_path": ".name"
	   },
	   "occupation": {
		   "object_path":".occupation"
	   }
   }
```

Based on this definition the "id" and "dob" fields would not be included in the index.
OPTIONS


```
	-c	sets the collection to be used
	-collection	sets the collection to be used
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-v	display version
	-version	display version
```

## EXAMPLES

In the example the index will be created for a collection called "characters".

```shell
    dsindexer -c characters email-mapping.json email-index
```

This will build a Bleve index called "email-index" based on the index defined
in "email-mapping.json".


dsindexer v0.0.2