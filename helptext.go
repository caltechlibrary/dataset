/**
 * helptext.go holds the common help docuemntation shared between the dataeset and datasetd commands.
 */
package dataset

const (
  DatasetHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}

# NAME

{app_name} 

# SYNOPSIS

{app_name} [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

{app_name} command line interface supports creating JSON object
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


A word about "keys". {app_name} uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensistive.

A word about "GLOBAL_OPTIONS" in v2 of {app_name}.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
{app_name}).  There are three "GLOBAL_OPTIONS" that are exception
and they are ` + "`" + `-version` + "`" + `, ` + "`" + `-help` + "`" + `
and ` + "`" + `-license` + "`" + `. All other options come
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
   {app_name} help init

   {app_name} init my_objects.ds 

   {app_name} model my_objects.ds

   {app_name} help create

   {app_name} create my_objects.ds "123" '{"one": 1}'

   {app_name} create my_objects.ds "234" mydata.json 
   
   cat <<EOT | {app_name} create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   {app_name} update my_objects.ds "123" '{"one": 1, "two": 2}'

   {app_name} delete my_objects.ds "345"

   {app_name} keys my_objects.ds
~~~

This is an example of initializing a Pairtree JSON documentation
collection using the environment.

~~~
{app_name} init '${C_NAME}' pairtree
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

To specify Postgres as the storage for your JSON document collection.
You'd use something like --

~~~
{app_name} init '${C_NAME}' \\
  'postgres://${USER}@localhost/${DB_NAME}?sslmode=disable'
~~~


In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME. USER is used
for the Postgres username and DB_NAME is used for the Postgres
database name.  The sslmode option was specified because Postgres
in this example was restricted to localhost on a single user machine.


{app_name} {version}

`

  DatasetdHelpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SETTINGS_FILE

# DESCRIPTION

{app_name} provides a web service for one or more dataset collections. Requires the
collections to exist (e.g. created previously with the dataset cli). It requires a
settings JSON or YAML file that decribes the web service configuration and
permissions per collection that are available via the web service.

# OPTIONS

-help
: display detailed help

-license
: display license

-version
: display version

-debug
: log debug information

# SETTINGS_FILE

The settings files provides {app_name} with the configuration
of the service web service and associated dataset collection(s).

It can be writen as either a JSON or YAML file. If it is a YAML file
you should use the ".yaml" extension so that {app_name} will correctly
parse the YAML.

The top level YAML attributes are

host
: (required) hostname a port for the web service to listen on, e.g. localhost:8485

htdocs
: (optional) if set static content will be serviced based on this path. This is a
good place to implement a browser side UI in HTML, CSS and JavaScript.

collections
: (required) A list of dataset collections that will be supported with this
web service. The dataset collections can be pairtrees or SQL stored. The
latter is preferred for web access to avoid problems of write collisions.

The collections object is a list of configuration objects. The configuration
attributes you should supply are as follows.

dataset
: (required) The path to the dataset collection you are providing a web API to.

query
: (optional) is map of query name to SQL statement. A POST is used to access
the query (i.e. a GET or POST To the path "`+"`"+`/api/<COLLECTION_NAME>/query/<QUERY_NAME>/<FIELD_NAMES>`+"`"+`")
The parameters submitted in the post are passed to the SQL statement.
NOTE: Only dataset collections using a SQL store are supported. The SQL
needs to conform the SQL dialect of the store being used (e.g. MySQL, Postgres,
SQLite3). The SQL statement functions with the same contraints of dsquery SQL
statements. The SQL statement is defined as a YAML text blog.

## API Permissions

The following are permissioning attributes for the collection. These are
global to the collection and by default are set to false. A read only API 
would normally only include "keys" and "read" attributes set to true.

keys
: (optional, default false) allow object keys to be listed

create
: (optional, default false) allow object creation through a POST to the web API

read
: (optional, default false) allow object to be read through a GET from the web API

update
: (optional, default false) allow object updates through a PUT to the web API.

delete
: (optional, default false) allow object deletion through a DELETE to the web API.

attachments
: (optional, default false) list object attachments through a GET to the web API.

attach
: (optional, default false) Allow adding attachments through a POST to the web API.

retrieve
: (optional, default false) Allow retrieving attachments through a GET to the web API.

prune
: (optional, default false) Allow removing attachments through a DELETE to the web API.

versions
: (optional, default false) Allow setting versioning of attachments via POST to the web API.


# EXAMPLES

Starting up the web service

~~~
   {app_name} settings.yaml
~~~

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the {app_name} 
service to interact with the collection via cURL. 

To create the dataset collection we use the "dataset" command and the
"vi" text edit (use can use your favorite text editor instead of vi).

~~~
    createdb t1
    dataset init t1.ds \
	   "postgres://$PGUSER:$PGPASSWORD@/t1?sslmode=disable"
	vi settings.yaml
~~~

You can create the "settings.yaml" with this Bash script.
I've created an htdocs directory to hold the static content
to interact with the dataset web service.

~~~
mkdir htdocs
cat <<EOT >settings.yaml
host: localhost:8485
htdocs: htdocs
collections:
  # Each collection is an object. The path prefix is
  # /api/<dataset_name>/...
  - dataset: t1.ds
    # What follows are object level permissions
	keys: true
    create: true
    read: true
	update: true
	delete: true
    # These are attachment related permissions
	attachments: true
	attach: true
	retrieve: true
	prune: true
    # This sets versioning behavior
	versions: true
EOT
~~~

Now we can run {app_name} and make the dataset collection available
via HTTP.

~~~
    {app_name} settings.yaml
~~~

You should now see the start up message and any log information display
to the console. You should open a new shell sessions and try the following.

We can now use cURL to post the document to the "api//t1.ds/object/one" end
point. 

~~~
    curl -X POST http://localhost:8485/api/t1.ds/object/one \
	    -d '{"one": 1}'
~~~

Now we can list the keys available in our collection.

~~~
    curl http://localhost:8485/api/t1.ds/keys
~~~

We should see "one" in the response. If so we can try reading it.

~~~
    curl http://localhost:8485/api/t1.ds/read/one
~~~

That should display our JSON document. Let's try updating (replacing)
it. 

~~~
    curl -X POST http://localhost:8485/api/t1.ds/object/one \
	    -d '{"one": 1, "two": 2}'
~~~

If you read it back you should see the updated record. Now lets try
deleting it.

~~~
	curl http://localhost:8485/api/t1.ds/object/one
~~~

List the keys and you should see that "one" is not longer there.

~~~
    curl http://localhost:8485/api/t1.ds/keys
~~~

You can run a query named 'browse' that is defined in the YAML configuration like this.

~~~
	curl http://localhost:8485/api/t1.ds/query/browse
~~~

or 

~~~
	curl -X POST -H 'Content-type:application/json' -d '{}' http://localhost:8485/api/t1.ds/query/browse
~~~

In the shell session where {app_name} is running press "ctr-C"
to terminate the service.


{app_name} {version}
`

	DatasetdApiText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}


# {app_name} REST API

{app_name} provides a RESTful JSON API for working with a dataset collection. This document describes the path expressions and to interact with the API.  Note some of the methods and paths require permissions as set in the {app_name} YAML or JSON [settings file]({app_name}_yaml.5.md).

## basic path expressions

There are three basic forms of the URL paths supported by the API.

- `+"`"+`/api/<COLLECTION_NAME>/keys`+"`"+`, get a list of all keys in the the collection
- `+"`"+`/api/<COLLECTION_NAME>/object/<OPTIONS>`+"`"+`, interact with an object in the collection (e.g. create, read, update, delete)
- `+"`"+`/api/<COLLECTION_NAME>/query/<QUERY_NAME>/<FIELDS>`+"`"+`, query the collection and receive a list of objects in response

The "`+"`"+`<COLLECTION_NAME>`+"`"+`" would be the name of the dataset collection, e.g. "mydata.ds".

The "`+"`"+`<OPTIONS>`+"`"+`" holds any additional parameters related to the verb. Options are separated by the path delimiter (i.e. "/"). The options are optional. They do not require a trailing slash.

The "`+"`"+`<QUERY_NAME>`+"`"+`" is the query name defined in the YAML configuration for the specific collection.

The "`+"`"+`<FIELDS>`+"`"+`" holds the set of fields being passed into the query. These are delimited with the path separator like with options (i.e. "/"). Fields are optional and they do not require a trailing slash.

## HTTP Methods

The {app_name} REST API follows the rest practices. Good examples are POST creates, GET reads, PUT updates, and DELETE removes. It is important to remember that the HTTP method and path expression need to match form the actions you'd take using the command line version of dataset. For example to create a new object you'd use the object path without any options and a POST expression. You can do a read of an object using the GET method along withe object style path.

## Content Type and the API

The REST API works with JSON data. The service does not support multipart urlencoded content. You MUST use the content type of `+"`"+`application/json`+"`"+` when performing a POST, or PUT. This means if you are building a user interface for a collections {app_name} service you need to appropriately use JavaScript to send content into the API and set the content type to `+"`"+`application/json`+"`"+`.

## Examples

Here's an example of a list, in YAML, of people in a collection called "people.ds". There are some fields for the name, sorted name, display name and orcid. The pid is the "key" used to store the objects in our collection.

~~~yaml
people:
  - pid: doe-jane
    family: Doe
    lived: Jane
    orcid: 9999-9999-9999-9999
~~~

In JSON this would look like

~~~json
{
  "people": [
    {
      "pid": "doe-jane",
      "family": "Doe",
      "lived": "Jane",
      "orcid": "9999-9999-9999-9999"
    }
  ]
}
~~~

### create

The create action is formed with the object URL path, the POST http method and the content type of "application/json". It POST data is expressed as a JSON object.

The object path includes the dataset key you'll assign in the collection. The key must be unique and not currently exist in the collection.

If we're adding an object with the key of "doe-jane" to our collection called "people.ds" then the object URL path would be  `+"`"+`/api/people.ds/object/doe-jane`+"`"+`. NOTE: the object key is included as a single parameter after "object" path element.

Adding an object to our collection using curl looks like the following.

~~~shell
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"pid": "doe-jane", "family": "Doe", "lived": "Jane", "orcid": "9999-9999-9999-9999" }' \
  http://localhost:8485/api/people.ds/object/doe-jane  
~~~

### read

The read action is formed with the object URL path, the GET http method and the content type of "application/json".  There is no data
aside from the URL to request the object. Here's what it would look like using curl to access the API.

~~~shell
curl http://localhost:8485/api/people.ds/object/doe-jane  
~~~

### update

Like create update is formed from the object URL path, content type of "application/json" the data is expressed as a JSON object.
Onlike create we use the PUT http method.

Here's how you would use curl to get the JSON expression of the object called "doe-jane" in your collection.

~~~shell
curl -X PUT \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"pid": "doe-jane", "family": "Doe", "lived": "Jane", "orcid": "9999-9999-9999-9999" }' \
  http://localhost:8485/api/people.ds/object/doe-jane  
~~~

This will overwrite the existing "doe-jane". NOTE the record must exist or you will get an error.

### delete

If you want to delete the "doe-jane" record in "people.ds" you perform an http DELETE method and form the url like a read.

~~~shell
curl -X DELETE http://localhost:8485/api/people.ds/object/doe-jane  
~~~

## query

The query path lets you run a predefined query from your settings YAML file. The http method used is a POST. This is becaue we need to send data inorder to receive a response. The resulting data is expressed as a JSON array of object. Like with create, read, update and delete you use the content type of "application/json".

In the settings file the queries are named. The query names are unique. One or many queries may be defined. The SQL expression associated with the name run as a prepared statement and parameters are mapped into based on the URL path provided. This allows you use many fields in forming your query.

Let's say we have a query called "full_name". It is defined to run the following SQL.

~~~sql
select src
from people
where src->>'family' like ?
  and src->>'lived' like ?
order by family, lived
~~~

NOTE: The SQL is has to retain the constraint of a single object per row, normally this will be "src" for dataset collections.

When you form a query path we need to indicate that the parameter for family and lived names need to get mapped to their respect positional references in the SQL. This is done as following url path. In this example "full_name" is the name of the query while "family" and "lived" are the values mapped into the parameters.

~~~
/api/people.ds/query/full_name/family/lived
~~~

The web form could look like this.  

~~~
<form id="query_name">
   <label for="family">Family</label> <input id="family" name="family" ><br/>
   <label for="lived">Lived</label> <input id="lived" name="lived" ><br/>
   <button type="submit">Search</button>
</form>
~~~

REMEMBER: the JSON API only supports the content type of "application/json" so you can use the browser's action and method in the form.

You would include JavaScript in the your HTML to pull the values out of the form and create a JSON object. If I searched
for someone who had the family name "Doe" and he lived name of "Jane" the object submitted to query might look like the following. 

~~~json
{
    "family": "Doe"
    "lived": "Jane"
}
~~~

The curl expression would look like the following simulating the form submission would look like the following.


~~~shell
curl -X POST \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{"family": "Doe", "lived": "Jane" }' \
  http://localhost:8485/api/people.ds/query/full_name/family/lived
~~~


`

	DatasetdServiceText = `%{app_name}(5) user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}


# {app_name} Service

The {app_name} based application can be configured to be managed by
systemd. You need to create a an appropriate service file with
Unit, Service and Install described.

## Example

Below is a generic {app_name} systemd style service file for a project
called citesearch implemented as citesearch.yaml using {app_name} to provide
the web service.

~~~
[Unit]

Description=A Citation search engine for CaltechTHESIS search

[Service]
Type=simple
ExecStart=/usr/local/bin/{app_name} /Sites/citesearch/citesearch.yaml

[Install]
WantedBy=multi-user.target
~~~

`

	DatasetdYAMLText = `%{app_name}(5) user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}


# {app_name} YAML configuration

The dataset RESTful JSON API is configured using either a YAML or JSON file. YAML is preferred as it is more readable but JSON remains supported for backward compatibility. What follows is the description of the YAML configuration. Note option elements are optional and for booleans will default to false if missing.

## Top level

host
: (required) this is the hostname and port for the web service, e.g. localhost:8485

htdocs
: (optional) if this is a non-empty it will be used as the path to static resouce provided with the web service.
These are useful for prototyping user interfaces with HTML, CSS and JavaScript interacting the RESTful JSON API.


collections
: (required), a list of datasets to be manage via the web service.

Each collection object has the following properties. Notes if you are trying to provide a read-only API
then you will want to include permissions for keys, read and probably query (to provide a search feature).

dataset
: (required) this is a path to your dataset collection.

query
: (optional) Is a map of query name to SQL statements. Each name will trigger a the execution of a SQL statement.
The query expects a POST. Fields are mapped to the SQL statement parameters. If a pairtree store is used a
indexing will be needed before this will work as it would use the SQLite 3 database to execute the SQL statement against.
Otherwise the SQL statement would conform to the SQL dialect of the SQL storage used (e.g. Postgres, MySQL or SQLite3).
The SQL statements need to conform to the same constraints as dsquery's implementation of SQL statements.

## API Permissions

API permissions are global. They are controlled with the following attributes. If the attributes are set to true
then they enable that permission. If you want to create a read only API then set keys, read to true. Query
support can be added via the query parameter. These are indepent so if you didn't want to allow keys or full
objects to be retrieve you could just provide access via defined queries.

keys
: (optional, default false) If true allow keys for the collection to be retrieved with a GET to `+"`"+`/api/<COLLECTION_NAME>/keys`+"`"+`

read
: (optional, default false) If true allow objects to be read via a GET to `+"`"+`/api/<COLLLECTION_NAME>/object/<KEY>`+"`"+`

create
: (optional, default false) If true allow object to be created via a POST to `+"`"+`/api/<COLLLECTION_NAME>/object`+"`"+`

update
: (optional, default false) If true allow object to be updated via a PUT  to `+"`"+`/api/<COLLECTION_NAME>/object/<KEY>`+"`"+`

delete
: (optional, default false) If true allow obejct to be deleted via a DELETE to `+"`"+`/api/<COLLECTION_NAME>/object/<KEY>`+"`"+`

attachments
: (optional, default false) list object attachments through a GET to the web API.

attach
: (optional, default false) Allow adding attachments through a POST to the web API.

retrieve
: (optional, default false) Allow retrieving attachments through a GET to the web API.

prune
: (optional, default false) Allow removing attachments through a DELETE to the web API.

versions
: (optional, default false) Allow setting versioning of attachments via POST to the web API.


`

	DSQueryHelpText = `%{app_name}(1) dataset user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] C_NAME SQL_STATEMENT [PARAMS]

# DESCRIPTION

__{app_name}__ is a tool to support SQL queries of dataset collections. 
Pairtree based collections should be index before trying to query them
(see '-index' option below). Pairtree collections use the SQLite 3
dialect of SQL for querying.  For collections using a SQL storage
engine (e.g. SQLite3, Postgres and MySQL), the SQL dialect reflects
the SQL of the storage engine.

The schema is the same for all storage engines.  The scheme for the JSON
stored documents have a four column scheme.  The columns are "_key", 
"created", "updated" and "src". "_key" is a string (aka VARCHAR),
"created" and "updated" are timestamps while "src" is a JSON column holding
the JSON document. The table name reflects the collection
name without the ".ds" extension (e.g. data.ds is stored in a database called
data having a table also called data).

The output of __{app_name}__ is a JSON array of objects. The order of the
objects is determined by the your SQL statement and SQL engine. There
is an option to generate a 2D grid of values in JSON, CSV or YAML formats.
See OPTIONS for details.

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

SQL_STATEMENT
: The SQL statement should conform to the SQL dialect used for the
JSON store for the JSON store (e.g.  Postgres, MySQL and SQLite 3).
The SELECT clause should return a single JSON object type per row.
__{app_name}__ returns an JSON array of JSON objects returned
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
representation with the first row as the attribute names.

-yaml STRING_OF_ATTRIBUTE_NAMES
: Like -grid this takes our list of dataset objects and a list of attribute
names but rather than create a 2D JSON of values it creates YAML 
representation.

-index
: This will create a SQLite3 index for a collection. This enables {app_name}
to query pairtree collections using SQLite3 SQL dialect just as it would for
SQL storage collections (i.e. don't use with postgres, mysql or sqlite based
dataset collections. It is not needed for them). Note the index is always
built before executing the SQL statement.

# EXAMPLES

Generate a list of JSON objects with the ` + "`" + `_key` + "`" + ` value
merged with the object stored as the ` + "`" + `._Key` + "`" + ` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (note: in Postgres the ` + "`" + `||` + "`" + ` is very helpful).

~~~
{app_name} data.ds "SELECT jsonb_build_object('_Key', _key)::jsonb || src::jsonb FROM data"
~~~

In this example we're returning the "src" in our collection by querying
for a "id" attribute in the "src" column. The id is passed in as an attribute
using the Postgres positional notatation in the statement.

~~~
{app_name} data.ds "SELECT src FROM data WHERE src->>'id' = $1 LIMIT 1" "xx103-3stt9"
~~~

`

	DSImporterHelpText = `%{app_name}(1) dataset user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] C_NAME CSV_FILENAME KEY_COLUMN

# DESCRIPTION

__{app_name}__ is a tool to import CSV content into a dataset collection
where the column headings become the attribute names and the row values
become the attribute values.

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

CSV_FILENAME
: The name of the CSV file to import

KEY_COLUMN
: The column name to use the they object key. If none is provided then
the first column is used as the object key. Keys values must be unique.


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-comma
: Set column delimiter

-comment
: Set row comment delimiter

-overwrite
: Overwrite objects on key collision

# EXAMPLES

Import a file with three columns

- item_code
- title
- location

The "item_code" is unique for each row. The data is stored
in a file called "books.csv". We are importing the CSV file
into a collections called. "shelves.ds"

~~~
{app_name} shelves.ds books.csv "item_code"
~~~

`

)
