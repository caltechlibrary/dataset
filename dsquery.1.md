%dsquery(1) dataset user manual | version 2.1.8 030d0c1
% R. S. Doiel and Tom Morrell
% 2023-10-11

# NAME

dsquery

# SYNOPSIS

dsquery [OPTIONS] C_NAME SQL_STATEMENT [PARAMS]

# DESCRIPTION

__dsquery__ is a tool to support SQL queries of dataset collections. 
Pairtree based collections should be index before trying to query them
(see '-index' option below). Pairtree collections use the SQLite 3
dialect of SQL for querying them.  For collections using a SQL storage
engine (e.g. SQLite3, Postgres and MySQL), the SQL dialect used is that
of the SQL storage engine chosen.

The schema is the same for all storage engines.  The scheme for the JSON
stored documents have a four column scheme.  The columns are "_key", 
"created", "updated" and "src". "_key" is a string (aka VARCHAR),
"created" and "updated" are timestamps while "src" is a JSON column holding
the JSON document. The table name reflects the collection
name without the ".ds" extension (e.g. data.ds is stored in a database called
data having a table also called data).

The output of __dsquery__ is a JSON arrary of objects. The order of the
objects is determined by the your SQL statement and SQL engine. There
is an option to generate a 2D grid of values and CSV format are also
supported as options (see '-grid' and '-csv' below).

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

SQL_STATEMENT
: The SQL statement should conform to the SQL dialect used for the
JSON store for the JSON store (e.g.  Postgres, MySQL and SQLite 3).
The SELECT clause should return a single JSON object type per row.
__dsquery__ returns an JSON array of JSON objects returned
by the SQL query.

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


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-pretty
: pretty print the resulting JSON array

-sql SQL_FILENAME
: read SQL from a file. If filename is "-" then read SQL from standard input.

-grid STRING_OF_ATTRIBUTE_NAMES
: Returns list as a 2D grid of values. This options requires a comma delimited
string of attribute names for the outer object to include in grid output. It
can be combined with -pretty options.

-csv STRING_OF_ATTRIBUTE_NAMES
: Like -grid this takes our list of dataset objects and a list of attribute
names but rather than create a 2D JSON array of values it creates CSV 
represnetation with the first row as the attribute names.

-index
: This will create a SQLite3 index for a collection. This enables dsquery
to query pairtree collections using SQLite3 SQL dialect just as it would for
SQL storage collections (i.e. don't use with postgres, mysql or sqlite based
dataset collections. It is not needed for them). Note the index is always
built before executing the SQL statement.

# EXAMPLES

Generate a list of JSON objects with the `_key` value
merged with the object stored as the `._Key` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (note: in Postgres the `||` is very helpful).

~~~
dsquery data.ds "SELECT jsonb_build_object('_Key', _key)::jsonb || src::jsonb FROM data"
~~~

In this example we're returning the "src" in our collection by querying
for a "id" attribute in the "src" column. The id is passed in as an attribute
using the Postgres positional notatation in the statement.

~~~
dsquery data.ds "SELECT src FROM data WHERE src->>'id' = $1 LIMIT 1" "xx103-3stt9"
~~~


