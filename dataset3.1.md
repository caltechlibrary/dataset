%dataset3(1) user manual | version 3.0.0-alpha ee097b6
% R. S. Doiel and Tom Morrell
% 2025-05-28

# NAME

dataset3 

# SYNOPSIS

dataset3 [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

dataset3 command line interface supports creating JSON object
collections and managing the JSON object documents in a collection.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

# SUPPORTED VERBS

help
: will give documentation of help on a verb, e.g. "help create"

init C_NAME
: Initialize a new dataset collection named with C_NAME.

create [OPTION] C_NAME KEY
: creates a new JSON document in the collection

read [OPTION] C_NAME KEY
: retrieves the "current" version of a JSON document from 
  the collection writing it standard out

update [OPTION] C_NAME KEY
: updates a JSON document in the collection

delete C_NAME KEY
: removes all versions of a JSON document from the collection

keys [OPTION] C_NAME
: returns a list of keys in the collection

codemeta C_NAME [PATH_TO_NEW_CODEMETA_JSON]
: displays an existing codemetada.json file for the collections or
if an optional path to a new codemeta.json it copies the file and updates the 
collections metadata.

dump C_NAME
: This will write out all dataset collection records in a JSONL document.
JSONL shows on JSON object per line, see https://jsonlines.org for details.
The object rendered will have two attributes, "key" and "object". The
key corresponds to the dataset collection key and the object is the JSON
value retrieved from the collection.

load [OPTION] C_NAME
: This will read JSON objects one per line from standard input. This
format is often called JSONL, see https://jsonlines.org. The object
has two attributes, key and object. 

A word about "keys". dataset3 uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys are composed as lower case 
alpha numeric characters but may include period, dash and underscore.
While keys maybe provided in upper and lower case they are always
converted to lowercase internally.

There are three "GLOBAL_OPTIONS" in v3 of dataset3.T hey are 
`-version`, `-help`
and `-license`. All other options come
after the verb and apply to the specific action the verb
implements.

# STORAGE TYPE

There are currently three support storage options for JSON documents in a dataset collection.

- SQLite3 database (default),

The following storage engines were removed in v3 -- pairtree, MySQL and PostgreSQL. If you need
to migrate data from a v2 dataset instance use the dump verb. Then you can import the data
into the v3 dataset collection using the load verb.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
   dataset3 help init

   dataset3 init my_objects.ds 

   dataset3 help create

   dataset3 create my_objects.ds "123" '{"one": 1}'

   dataset3 create my_objects.ds "234" mydata.json 
   
   cat <<EOT | dataset3 create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   dataset3 update my_objects.ds "123" '{"one": 1, "two": 2}'

   dataset3 delete my_objects.ds "345"

   dataset3 keys my_objects.ds
~~~

This is an example of initializing a JSON documentation
collection.

~~~
dataset3 init '${C_NAME}'
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

dataset3 3.0.0-alpha


