
USAGE

 dataset [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PARAMETER ...]

SYNOPSIS

dataset command line interface supports creating JSON object
collections and managing the JSON object documents in a collection. As of v3
SQLite3 is the storage option supported.  Prior versions of dataset supported
storing objects in a pairtree, in MySQL and PostgreSQL.  You can use dump
from a dataset v2 then use load with dataset v3 to migrate your collection to
v3.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

SUPPORTED VERBS

help VERB
: will give this this documentation of help on a verb

init C_NAME
: creates a new dataset collection. By default the collection uses
  SQLite3 for a JSON store.

create [OPTION] C_NAME KEY [DATA]
: creates a new document in the collection

read [OPTION] C_NAME KEY
: retrieves a document from the collection writing it standard out

update [OPTION] C_NAME KEY [DATA]
: updates a document in the collection

delete C_NAME KEY
: removes a document from the collection

keys [OPTION] C_NAME
: returns a list of keys, one per line held in the collection

codemeta C_NAME [PATH_TO_NEW_CODEMETA_JSON]
: display existing collection codemeta.json file or 
if a path to a codemeta.json file is provided copies it into
the collection's root directory.

load [OPTION] C_NAME
: imports a json lines document into the current collection

dump C_NAME
: exports the collection aas a json lines document

query [OPTION] C_NAME SQL_STMT
: allows you to pass a SQL query that returns a single object
  or value to the SQLite 3 engine. This provides a flexible way to work
  with the objects in the collection as well as create lists of objects.

history is implemented as a SQLite3 history table where the verison number
is a single ascending integer value. The key and version number form a
complex primary key to ensure uniqueness in the history table.

A word about "keys". dataset uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys are formed from lowered case
alphanumeric characters but may also include period, dash and underscore
characters. If you enter an upper case key it will be stored as lower.
Case sensitivity does not provide uniquenes for keys.

There are three "GLOBAL_OPTIONS". They are `-version`,`-help` and `-license`. All other
options come after the verb and apply to the specific action the verb
implements.

