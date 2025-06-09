%C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe(1) user manual | version 2.2.7 9c44ac2
% R. S. Doiel and Tom Morrell
% 2025-06-02

# NAME

C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe 

# SYNOPSIS

C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe command line interface supports creating JSON object
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

frame-names
: lists the frames defined in a collection

frame
: will add a data frame to a collection 

frame-def
: will return the definition of a frame

frame-keys
: will retrieve the object keys in a frame

frame-objects
: will retrieve the object list in a frame

reframe
: will recreate a frame using its existing definition but
  replacing objects based on a new set of keys provided

refresh
: will update all objects currently in the frame based on the
  current state of the collection. Any keys deleted in the collection
  will be delete from the frame.

delete-frame
: will remove a frame from the collection

has-frame
: will return true (exit 0) if frame exists, false (exit 1)
  if not

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


A word about "keys". C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensistive.

A word about "GLOBAL_OPTIONS" in v2 of C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe).  There are three "GLOBAL_OPTIONS" that are exception
and they are `-version`, `-help`
and `-license`. All other options come
after the verb and apply to the specific action the verb
implements.

# STORAGE TYPE

There are currently three support storage options for JSON documents in a dataset collection.

- SQLite3 database (default),
- Pairtree (pre-2.1 default)
- Postgres

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
   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe help init

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe init my_objects.ds 

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe model my_objects.ds

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe help create

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe create my_objects.ds "123" '{"one": 1}'

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe create my_objects.ds "234" mydata.json 
   
   cat <<EOT | C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe update my_objects.ds "123" '{"one": 1, "two": 2}'

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe delete my_objects.ds "345"

   C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe keys my_objects.ds
~~~

This is an example of initializing a Pairtree JSON documentation
collection using the environment.

~~~
C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe init '${C_NAME}' pairtree
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

To specify Postgres as the storage for your JSON document collection.
You'd use something like --

~~~
C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe init '${C_NAME}' \\
  'postgres://${USER}@localhost/${DB_NAME}?sslmode=disable'
~~~


In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME. USER is used
for the Postgres username and DB_NAME is used for the Postgres
database name.  The sslmode option was specified because Postgres
in this example was restricted to localhost on a single user machine.


C:\Users\rsdoi\Source\GitHub\CaltechLibrary\dataset\bin\dataset.exe 2.2.7


