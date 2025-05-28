query
============

This will run a SQL query against the SQLite3 JSON store. It returns
a list of JSON objects or an error. The SQL query must only return
one object per row.

# SYNOPSIS

dataset3 query [OPTIONS] C_NAME [SQL_STATEMENT] [PARAMS]

# DESCRIPTION

__dataset3 query__ is a tool to support SQL queries of dataset collections. 
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

The output of __dataset3 query__ is a JSON array of objects. The order of the
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
a parameter __dataset3__ will expect to read SQL from standard
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
    dataset3 query mycollection.ds \\
      "select src from mycollection order by created desc"
~~~

Read the SQL statement from a file called "report.sql".

~~~shell
    dataset3 query -sql report.sql mycollection.ds
~~~

Generate a list of JSON objects with the `_Key` value
merged with the object stored as the `._Key` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (NOTE: in PostgreSQL the `||` is very helpful).

~~~
dataset3 query data.ds "SELECT json_object('key', _Key) FROM data"
~~~

In this example we're returning the "src" in our collection by querying
for a "id" attribute in the "src" column. The id is passed in as an attribute
using the Postgres positional notatation in the statement.

~~~
dataset3 query data.ds "SELECT src FROM data WHERE src->>'id' = $1 LIMIT 1" "xx103-3stt9"
~~~

This is an example of sending a formated query to return a list of objects with version info.

~~~
cat <<SQL | dataset3 query data.ds
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

