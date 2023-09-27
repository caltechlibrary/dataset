%dsquery(1) dataset user manual | version 2.1.4 978b141
% R. S. Doiel and Tom Morrell
% 2023-09-27

# NAME

dsquery

# SYNOPSIS

dsquery [OPTIONS] C_NAME SQL_STATEMENT [PARAMS]

# DESCRIPTION

__dsquery__ is a tool to support SQL queries of dataset collections that
use SQL storage for the collection's JSON documents.  It takes a dataset
collection name and a sql statement returning the results. This will allow us
to improve our feeds building process by taking advantage of SQL and the
collection's SQL database engine.

The scheme for the JSON stored documents have a four column scheme. 
The columns are "_key", "created", "updated" and "src". The are stored
in a table with the same name as the database which is formed from the
C_NAME without extension (e.g. data.ds is stored in a database called
data having a table also called data).

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

help
: display help

license
: display license

version
: display version


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


