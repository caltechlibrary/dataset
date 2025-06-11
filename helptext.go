// This is part of the dataset package.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrell, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2025, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package dataset

const (
    // List help subjects
    cliHelp = `
Help Topics
-----------

Overview
========

- Description
- History
- Examples

Collections
===========

- Codemeta
- Init
- Load
- Dump

Objects
=======

- Write
- Read
- Delete
- Keys
- HasKey
- Query

`
	// cliDescription describes how to use the cli
	cliDescription = `
SYNOPSIS

 {app_name} [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PARAMETER ...]

DESCRIPTION

{app_name} command line interface supports creating JSON object
collections and managing the JSON object documents in a collection. As of v3
SQLite3 is the storage option supported.  Prior versions of dataset supported
storing objects in a pairtree in additiojnal to SQlite3, MySQL and PostgreSQL.
You can use dump from a dataset v2 then use load with dataset v3 to migrate
your collection to v3.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

SUPPORTED VERBS

help VERB
: will give this this documentation of help on a verb

init C_NAME
: creates a new dataset collection. By default the collection uses
  SQLite3 for a JSON store.

write [OPTION] C_NAME KEY [DATA]
: creates or replaces a JSON document in the collection

read [OPTION] C_NAME KEY
: retrieves a document from the collection writing it standard out

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

history is implemented as a SQL history table where the verison number
is a single ascending integer value. The key and version number form a
complex primary key to ensure uniqueness in the history table.

A word about "keys". {app_name} uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys are formed from lowered case
alphanumeric characters but may also include period, dash and underscore
characters. If you enter an upper case key it will be stored as lower.
Case sensitivity does not provide uniquenes for keys.

There are three "GLOBAL_OPTIONS". They are ` + "`" + `-version` + "`," + "`" + 
`-help` + "`" + ` and ` + "`" + `-license` + "`" + `. All other
options come after the verb and apply to the specific action the verb
implements.

`


    // cliExamples lists some examples of using the cli
	cliExamples = `
EXAMPLES

~~~
   {app_name} help init

   {app_name} init my_objects.ds 

   {app_name} help write

   {app_name} write my_objects.ds "123" '{"one": 1}'

   cat <<EOT | {app_name} write my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   {app_name} write my_objects.ds "123" '{"one": 1, "two": 2}'

   {app_name} delete my_objects.ds "345"

   {app_name} keys my_objects.ds

   {app_name} hasKey my_objects.ds "345"

   {app_name} dump my_objects.ds >objects.jsonl

   {app_name} load my_objects.ds <objects.jsonl

   cat <<SQL | {app_name} query my_objects.ds 
   select json_object('key', _Key, 'version', version) as obj
   from my_objects_history
   where _Key "345"
   order by version desc
   SQL
~~~

`

	//
	// cli specific help, not exported
	//

	// Taken from docs/write.md
	cliWrite = `
write
======

Syntax
------

~~~shell
    cat JSON_DOCNAME | {app_name} write COLLECTION_NAME KEY
    {app_name} write -i JSON_DOCNAME COLLECTION_NAME KEY
    {app_name} write COLLECTION_NAME KEY JSON_VALUE
    {app_name} write COLLECTION_NAME KEY JSON_FILENAME
~~~

Description
-----------

write adds or replaces a JSON document to a collection. The JSON 
document can be read from a standard in, a named file (with a 
".json" file extension) or expressed literally on the command line.

Usage
-----

In the following four examples *jane-doe.json* is a file on the 
local file system contains JSON data containing the JSON_VALUE 
of ` + "`" + `{"name":"Jane Doe"}` + "`" + `.  The KEY we will create is _r1_. 
Collection is "people.ds".  The following are equivalent in 
resulting record.

~~~shell
    cat jane-doe.json | {app_name} write people.ds r1
    {app_name} write -i blob.json people.ds r1
    {app_name} write people.ds r1 '{"name":"Jane Doe"}'
    {app_name} write people.ds r1 jane-doe.json
    cat jane-doe.json | {app_name} write people.ds r1
    {app_name} write people.ds r1 <jane-doe.json
~~~

`

	cliRead = `
read
====

Syntax
------

~~~shell
    {app_name} read [OPTION] COLLECTION_NAME KEY
~~~

Description
-----------

The writes the JSON document to standard out (unless you've 
specific an alternative location with the "-output" option)
for the given KEY.

Usage
-----

An example we're assuming there is a JSON document with a KEY 
of "r1". Our collection name is "data.ds"

~~~shell
    {app_name} read data.ds r1
~~~

Options
-------

Normally {app_name} outputs the JSON object as presented by the storage engine.
Use the `+"`"+`-jsonl`+"`"+` to force it to a single line (JSON line format).


~~~shell
    {app_name} read -jsonl data.ds r1
~~~

`

	cliDelete = `
delete
======

Syntax
------

~~~shell
    {app_name} delete COLLECTION_NAME KEY
~~~

Description
-----------

- delete - removes a JSON document from collection
  - requires JSON document name

Usage
-----

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

~~~shell
    {app_name} delete publications.ds r1
~~~

`

	cliKeys = `
keys
====

Syntax
------

~~~shell
    {app_name} keys COLLECTION_NAME
~~~

Description
-----------

List the JSON_DOCUMENT_ID available in a collection. Key order is not
guaranteed. Keys are forced to lower case when the record is created
in the {app_name} (as of version 1.0.2). Note combining "keys" with
a pipe and POSIX commands like "sort" can given a rich pallet of
ways to work with your {app_name} collection's keys.

Examples
--------

Here are three examples usage. Notice the sorting is handled by
the POSIX sort command which lets you sort ascending or descending
including sorting number strings.

~~~shell
    {app_name} keys COLLECTION_NAME
    {app_name} keys COLLECTION_NAME | sort
    {app_name} keys COLLECTION_NAME | sort -n
~~~


`

	cliHasKey = `
haskey
======

Syntax
------

~~~shell
    {app_name} [OPTIONS] haskey COLLECTION_NAME KEY_TO_CHECK_FOR
~~~

Description
-----------

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

Usage
-----

~~~shell
    {app_name} haskey people.ds '0000-0003-0900-6903'
    {app_name} haskey people.ds r1
~~~

`

	cliInit = `
init
====

Syntax
------

~~~shell
    dataset init COLLECTION_NAME
~~~

Description
-----------

_init_ creates a collection. Collections are created on local 
disc. By default it uses a SQLite3 database called "collection.db"
in the dataset directory for storing JSON Objects. As of v3 only
SQLite3 is supported.

Usage
-----

The following example command create a dataset collection 
named "data.ds".

~~~shell
    dataset init data.ds
~~~

NOTE: After each evocation of ` + "`" + `dataset init` + "`" + ` if all went well 
you will be shown an ` + "`" + `OK` + "`" + ` if everything went OK, otherwise
an error message. 

~~~shell
    dataset init data.ds
~~~

`

	cliCodemeta = `
codemeta
========

The command imports a codemeta.json file into the collection replacing
it's existing metadata.

~~~shell
   {app_name} codemeta data.ds ./codemeta.json
~~~

Without the codemeta filename it returns the existing codemeta values.

`

    cliLoad = `load [OPTION]
============

This will read a JSONL stream (JSONL, see https://jsonlines.org) 
and store the objects in a dataset collection. The objects must have
two attributes, __key__ and __object__.  The key value is used as
the key assigned in the collection and object value is the JSON stored.
Load reads from standard input.

The object structure matches the schema used by dataset collections
that are using SQLite3 for their object store. The dataset collection
loading the objects can use any dataset collection storage format
supported by the version of dataset featuring load and dump verbs.

# OPTION

-o, -overwrite
: If an object exists in the collection with the same key replace it.

-m, -max-capacity INTEGER
: Objects can be large in JSONL so you have the option of setting the
maximum buffer size for a single object. The integer value should be
greater than zero. The unit value is measured in mega bytes. "1" is
one meta byte, "10" would be ten mega bytes.

# EXAMPLE

Load a JSONL file, duplicate objects will not be overwritten.

~~~shell
    {app_name} load mycollection.ds <mycollection.jsonl
~~~

Load a JSONL file, duplicate objects will be overwritten.

~~~shell
    {app_name} load -overwrite mycollection.ds <mycollection.jsonl
~~~

`

cliDump = `dump
============

This will dump all the JSON objects in a collection, one
object per line (see https://jsonlines.org) (aka JSONL). 

The objects are written to standard output. Dump is the complement of
load verb. The objects dumped reflect the structured using when storing
objects in an SQLite3 database regardless of the store of the specific
collection. Like clone it provides a means of easily moving your data out
of a dataset collection.

Example
-------

~~~shell
    {app_name} dump mycollection.ds >mycollection.jsonl
~~~

`

cliQuery = `query
============

This will run a SQL query against the SQLite3 JSON store. It returns
a list of JSON objects or an error. The SQL query must only return
one object per row.

# SYNOPSIS

{app_name} query [OPTIONS] C_NAME [SQL_STATEMENT] [PARAMS]

# DESCRIPTION

__{app_name} query__ is a tool to support SQL queries of dataset collections. 
Pairtree based collections should be index before trying to query them
(see '-index' option below). Pairtree collections use the SQLite 3
dialect of SQL for querying.  For collections using a SQL storage
engine (e.g. SQLite3, Postgres and MySQL), the SQL dialect reflects
the SQL of the storage engine.

The schema is the same for all storage engines.  The scheme for the JSON
stored documents have a four column scheme.  The columns are "_Key", 
"created", "updated", "version" and "src". "_Key" is a string (aka VARCHAR),
"created" and "updated" are timestamps while "src" is a JSON column holding
the JSON document. The table name reflects the collection
name without the ".ds" extension (e.g. data.ds is stored in a database called
data having a table also called data).

The output of __{app_name} query__ is a JSON array of objects. The order of the
objects is determined by the your SQL statement and SQL engine. There
is an option to generate a 2D grid of values in JSON, CSV or YAML formats.
See OPTIONS for details.

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

SQL_STATEMENT
: The SQL statement should conform to the SQL dialect used for the
JSON store for the JSON store (e.g. SQLite3, Postgres or MySQL 8).
The SELECT clause should return a single JSON object type per row.
__query__ returns an JSON array of JSON objects returned
by the SQL query. NOTE: If you do not provide a SQL statement as
a parameter __{app_name}__ will expect to read SQL from standard
input.

PARAMS
: Is optional, it is any values you want to pass to the SQL_STATEMENT.

# SQL Store Scheme

_key
: The key or id used to identify the JSON documented stored.

src
: This is a JSON column holding the JSON document

created
: The date the JSON document was created in the table

updated
: The date the JSON document was updated

version
: The version of the object stored (zero indexed)

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-sql SQL_FILENAME
: read SQL from a file. If filename is "-" then read SQL from standard input.

-jsonl
: Output the query result using [JSON lines](https://jsonlines.org) format.

Example
-------

Return a JSON array of all objects by descending created date.

~~~shell
    {app_name} query mycollection.ds \\
      "select src from mycollection order by created desc"
~~~

You can also redirect a SQL statement via standard out like this.

~~~shell
cat <<SQL | {app_name} query mycollection.ds
  select src
  from mycollection
  order by created desc"
SQL
~~~

Read the SQL statement from a file called "report.sql".

~~~shell
    {app_name} query -sql report.sql mycollection.ds
~~~

Generate a list of JSON objects with the `+"`"+`_Key`+"`"+` value
merged with the object stored as the `+"`"+`._Key`+"`"+` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (NOTE: in PostgreSQL the `+"`"+`||`+"`"+` is very helpful).

~~~
{app_name} query data.ds "SELECT json_object('key', _Key) FROM data"
~~~

In this example we're returning the "src" in our collection by querying
for a "id" attribute in the "src" column. The id is passed in as an attribute
using the Postgres positional notatation in the statement.

~~~
{app_name} query data.ds "SELECT src FROM data WHERE src->>'id' = $1 LIMIT 1" "xx103-3stt9"
~~~

This is an example of sending a formated query to return a list of objects with version info.

~~~
cat <<SQL | {app_name} query data.ds
select
  json_object(
    "key": _Key,
    "src": src,
    "version": version
    "created": created,
    "updated: updated
  ) as obj
from data
order by _key;
SQL
~~~

`

)
