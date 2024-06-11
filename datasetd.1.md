%datasetd (1) user manual | verion 2.1.11 {release_hash}"
% R. S. Doiel
% {release_date}

# NAME

datasetd

# SYNOPSIS

datasetd [OPTIONS] SETTINGS_FILE

# DESCRIPTION

Runs a web service for one or more dataset collections. Requires
the collections to exist (e.g. created previously with the dataset
cli). It requires a settings JSON or YAML file that decribes the
web service configuration and permissions per collection that are
available via the web service.

# OPTIONS

-help
: display detailed help

-license
: display license

-version
: display version


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
    # These are frame level permissions
	frame_read: true
	frame_write: true
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

In the shell session where datasetd is running press "ctr-C"
to terminate the service.


datasetd 2.1.11


