%datasetd(1) user manual | version 2.3.1 9a3d898
% R. S. Doiel
% 2025-07-10

# NAME

datasetd

# SYNOPSIS

datasetd [OPTIONS] SETTINGS_FILE

# DESCRIPTION

datasetd provides a web service for one or more dataset collections. Requires the
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

The settings files provides datasetd with the configuration
of the service web service and associated dataset collection(s).

It can be writen as either a JSON or YAML file. If it is a YAML file
you should use the ".yaml" extension so that datasetd will correctly
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
the query (i.e. a GET or POST To the path "`/api/<COLLECTION_NAME>/query/<QUERY_NAME>/<FIELD_NAMES>`")
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
   datasetd settings.yaml
~~~

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the datasetd 
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

Now we can run datasetd and make the dataset collection available
via HTTP.

~~~
    datasetd settings.yaml
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

In the shell session where datasetd is running press "ctr-C"
to terminate the service.


datasetd 2.3.1

