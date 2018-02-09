
# USAGE

	dsindexer [OPTIONS] INDEX_DEF_JSON INDEX_NAME

## SYNOPSIS


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
	   "name":{
		   "object_path": ".name"
	   },
	   "occupation": {
		   "object_path":".occupation"
	   }
   }
```

Based on this definition the "id" and "dob" fields would not be included in the index.



## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    DATASET   # Set the dataset collection you're working with
```

## OPTIONS

Options will override any corresponding environment settings.

```
    -batch                    Set the size index batch, default is 100
    -c, -collection           sets the collection to be used
    -e, -examples             display examples
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -i, -input                input file name
    -id-file                  Create/Update an index for the ids in file
    -l, -license              display license
    -max-procs                Change the maximum number of CPUs that can executing simultaneously
    -nl, -newline             if set to false to suppress a trailing newline
    -o, -output               output file name
    -p, -pretty               pretty print output
    -quiet                    suppress error messages
    -t                        the label of the type of document you are indexing, e.g. accession, agent/person
    -update                   updating is slow, use this app if you want to update an exists
    -v, -version              display version
```


dsindexer v0.0.18-dev
