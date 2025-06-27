%dataset(1) user manual | version 2.3.0 c421d1c
% R. S. Doiel and Tom Morrell
% 2025-06-26

# NAME

dataset 

# SYNOPSIS

dataset [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

dataset command line interface supports creating JSON object
collections and managing the JSON object documents in a collection.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

# SUPPORTED VERBS

help
: will give documentation of help on a verb, e.g. "help create"

init [STORAGE_TYPE]
: Initialize a new dataset collection

model
: provides an experimental interactive data model generator creating
the "model.yaml" file in the data set collection's root directory.

create
: creates a new JSON document in the collection

read
: retrieves the "current" version of a JSON document from 
  the collection writing it standard out

update
: updates a JSON document in the collection

delete
: removes all versions of a JSON document from the collection

keys
: returns a list of keys in the collection

codemeta:
: copies metadata a codemeta file and updates the 
  collections metadata

attach
: attaches a document to a JSON object record

attachments
: lists the attachments associated with a JSON object record

retrieve
: creates a copy local of an attachement in a JSON record

detach
: will copy out the attachment to a JSON document 
  into the current directory 

prune
: removes an attachment (including all versions) from a JSON record

set-versioning
: will set the versioning of a collection. The versioning
  value can be "", "none", "major", "minor", or "patch"

get-versioning
: will display the versioning setting for a collection

dump
: This will write out all dataset collection records in a JSONL document.
JSONL shows on JSON object per line, see https://jsonlines.org for details.
The object rendered will have two attributes, "key" and "object". The
key corresponds to the dataset collection key and the object is the JSON
value retrieved from the collection.

load
: This will read JSON objects one per line from standard input. This
format is often called JSONL, see https://jsonlines.org. The object
has two attributes, key and object. 

join [OPTIONS] c_name, key, JSON_SRC
: This will join a new object provided on the command line with an
existing object in the collection.


A word about "keys". dataset uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensistive.

A word about "GLOBAL_OPTIONS" in v2 of dataset.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
dataset).  There are three "GLOBAL_OPTIONS" that are exception
and they are `-version`, `-help`
and `-license`. All other options come
after the verb and apply to the specific action the verb
implements.

# STORAGE TYPE

There are currently three support storage options for JSON documents in a dataset collection.

- SQLite3 database >= 3.40 (default)
- Postgres >= 12
- MySQL 8
- Pairtree (pre-2.1 default)

STORAGE TYPE are specified as a DSN URI except for pairtree which is just "pairtree".


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
   dataset help init

   dataset init my_objects.ds 

   dataset model my_objects.ds

   dataset help create

   dataset create my_objects.ds "123" '{"one": 1}'

   dataset create my_objects.ds "234" mydata.json 
   
   cat <<EOT | dataset create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   dataset update my_objects.ds "123" '{"one": 1, "two": 2}'

   dataset delete my_objects.ds "345"

   dataset keys my_objects.ds
~~~

This is an example of initializing a Pairtree JSON documentation
collection using the environment.

~~~
dataset init '${C_NAME}' pairtree
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

To specify Postgres as the storage for your JSON document collection.
You'd use something like --

~~~
dataset init '${C_NAME}' \\
  'postgres://${USER}@localhost/${DB_NAME}?sslmode=disable'
~~~


In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME. USER is used
for the Postgres username and DB_NAME is used for the Postgres
database name.  The sslmode option was specified because Postgres
in this example was restricted to localhost on a single user machine.


dataset 2.3.0


