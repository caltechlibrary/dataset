%datasetd(5) user manual | version 2.4.1-rc2 2d96591
% R. S. Doiel and Tom Morrell
% 2026-05-29


# datasetd YAML configuration

The dataset RESTful JSON API is configured using either a YAML or JSON file. YAML is preferred as it is more readable but JSON remains supported for backward compatibility. What follows is the description of the YAML configuration. Note option elements are optional and for booleans will default to false if missing.

## Top level

host
: (required) this is the hostname and port for the web service, e.g. localhost:8485

htdocs
: (optional) if this is a non-empty it will be used as the path to static resouce provided with the web service.
These are useful for prototyping user interfaces with HTML, CSS and JavaScript interacting the RESTful JSON API.

schemas
: (optional) A map of schema names to schema definitions. Schemas are used to validate JSON documents in collections.
  Each schema is defined as a model with elements, types, and validation rules. Schemas can be referenced
  by collections using the schema_name attribute. See the models package documentation for schema format.

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
Otherwise the SQL statement would conform to the SQL dialect of the SQL storage used (e.g. Postgres or SQLite3).
The SQL statements need to conform to the same constraints as dsquery's implementation of SQL statements.

schema_name
: (optional) The name of a schema defined in the top-level schemas map to use for validating documents in this collection.
  If set and validate is true, documents will be validated against this schema on create and update operations.

validate
: (optional, default false) If true and schema_name is set, enable schema validation for create and update operations.
  When validation fails, a 400 Bad Request response is returned with validation errors in the X-Validation-Errors header.

## API Permissions

API permissions are global. They are controlled with the following attributes. If the attributes are set to true
then they enable that permission. If you want to create a read only API then set keys, read to true. Query
support can be added via the query parameter. These are indepent so if you didn't want to allow keys or full
objects to be retrieve you could just provide access via defined queries.

keys
: (optional, default false) If true allow keys for the collection to be retrieved with a GET to `/api/<COLLECTION_NAME>/keys`

read
: (optional, default false) If true allow objects to be read via a GET to `/api/<COLLLECTION_NAME>/object/<KEY>`

create
: (optional, default false) If true allow object to be created via a POST to `/api/<COLLLECTION_NAME>/object`

update
: (optional, default false) If true allow object to be updated via a PUT  to `/api/<COLLECTION_NAME>/object/<KEY>`

delete
: (optional, default false) If true allow obejct to be deleted via a DELETE to `/api/<COLLECTION_NAME>/object/<KEY>`

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



