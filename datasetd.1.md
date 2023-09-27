%datasetd (1) user manual | verion 2.1.4 {release_hash}"
% R. S. Doiel
% {release_date}

# NAME

datasetd

# SYNOPSIS

datasetd [OPTIONS] SETTINGS_JSON_FILE

# DESCRIPTION

Runs a web service for one or more dataset collections. Requires
the collections to exist (e.g. created previously with the dataset
cli). It requires a settings JSON file that decribes the web service
configuration and permissions per collection that are available via
the web service.

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
   datasetd settings.json
~~~

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the datasetd 
service to interact with the collection via cURL. 

To create the dataset collection we use the "dataset" command and the
"vi" text edit (use can use your favorite text editor instead of vi).

~~~
    dataset init t1.ds
	vi settings.json
~~~

In the "setttings.json" file the JSON should look like.

~~~
    {
		"host": "localhost:8485",
		"sql_type": "mysql",
		"dsn": "DB_USER:DB_PASSWORD@/DB_NAME"
	}
~~~

Now we can run datasetd and make the dataset collection available
via HTTP.

~~~
    datasetd settings.json
~~~

You should now see the start up message and any log information display
to the console. You should open a new shell sessions and try the following.

We can now use cURL to post the document to the "/t1/create/one" end
point. 

~~~
    curl -X POST http://localhost:8485/t1/create/one \
	    -d '{"one": 1}'
~~~

Now we can list the keys available in our collection.

~~~
    curl http://localhost:8485/t1/keys
~~~

We should see "one" in the response. If so we can try reading it.

~~~
    curl http://localhost:8485/t1/read/one
~~~

That should display our JSON document. Let's try updating (replacing)
it. 

~~~
    curl -X POST http://localhost:8485/t1/update/one \
	    -d '{"one": 1, "two": 2}'
~~~

If you read it back you should see the updated record. Now lets try
deleting it.

~~~
	curl http://localhost:8485/t1/delete/one
~~~

List the keys and you should see that "one" is not longer there.

~~~
    curl http://localhost:8485/t1/keys
~~~

In the shell session where datasetd is running press "ctr-C"
to terminate the service.


datasetd 2.1.4


