Datasetd
========

Overview
--------

__datasetd__ is a minimal web service typically run on localhost port 8485 that exposes a dataset collection as a web service. It features a subset of functionallity available with the dataset command line program. __datasetd__ does support multi-process/asynchronous update to a dataset collection. 

__datasetd__ is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a stand alone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL style solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). __datasetd__ is a simple, miminal service.

NOTE: You could run __datasetd__ with access control based on a predictable set of URL paths by a web server such as Apache2 or NginX proxying to __datasetd__. That would require a robust understanding of the front end web server, it's access control mechanisms and how to defend a proxies service. That is beyond the skope of this project.

Configuration
-------------

__datasetd__ can make one or more dataset collections visible over HTTP/HTTPS. The dataset collections hosted need to be avialable on the same file system as where __datasetd__ is running. __datasetd__ is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line.  

The "settings.json" file has the following structure

```
    {
        "host": "localhost:8483",
        "collections": {
            "<COLLECTION_ID>": {
                "dataset": "<PATH_TO_DATASET_COLLECTION>",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": false
            }
        }
    }
```

In the "collections" object the "<COLLECTION_ID>" is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at "<COLLECTION_ID>". For each collection you can allow the following sub-paths, "create", "read", "update", "delete" and "keys". These sub-paths correspond to their counter parts in the dataset command line tool. By varying the settings of these you can support read only collections, drop off collections and function as a object store behind a web application.

Running datasetd
----------------

__datasetd__ runs as a HTTP service and as such can be exploited in the same manner as other services using HTTP.  You should only run __datasetd__ on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by __datasetd__ regardless of the file permissions they may in their account.
E.g. If all dataset collections are in a directory only allowed access to be the "web-data" user but another user on the system can run cURL then they can access the dataset collections based on the rights of the "web-data" user.  This is a typical situation for most web services and you need to be aware of it if you choose to run __datasetd__.

Supported Features
------------------

__datasetd__ provide a limitted subset of actions support by the standard datset command line tool. It only supports the following verbs

1. keys (return a list of all keys in the collection)
2. create (create a new JSON document in the collection)
3. read (read a JSON document from a collection)
4. update (update a JSON document in the collection)
5. delete (delete a JSON document in the collection)

Each of theses "actions" can be restricted in the configuration (
i.e. "settings.json" file) by setting the value to "false". If the
attribute for the action is not specified in the JSON settings file
then it is assumed to be "false".

Example
-------

E.g. if I have a settings file for "recipes" based on the collection
"recipes.ds" and want to make it read only I would make the attribute
"read" set to true and if I want the option of listing the keys in the collection I would set that true also.

```
{
    "host": "localhost:8485",
    "collections": {
        "recipes": {
            "dataset": "recipes.ds",
            "keys": true,
            "read": true
        }
    }
}
```

I would start __datasetd__ with the following command line.

```shell
    datasetd settings.json
```

This would display the start up message and log output of the service.

In another shell session I could then use cURL to list the keys and read
a record. In this example I assume that "waffles" is a JSON document
in dataset collection "recipes.ds".

```shell
    curl http://localhost:8485/recipies/read/waffles
```

This would return the "waffles" JSON document or a 404 error if the 
document was not found.

Listing the keys for "recipes.ds" could be done with this cURL command.

```shell
    curl http://localhost:8485/recipies/keys
```

This would return a list of keys, one per line. You could show
all JSON documents in the collection be retrieving a list of keys
and iterating over them using cURL. Here's a simple example in Bash.

```shell
    for KEY in $(curl http://localhost:8485/recipes/keys); do
       curl "http://localhost/8485/recipe/read/${KEY}"
    done
```

Documentation
-------------

__datasetd__ provide documentation as plain text output via request
to the service end points without parameters. Continuing with our
"recipes" example. Try the following URLs with cURL.

```
    curl http://localhost:8485
    curl http://localhost:8485/recipes
    curl http://localhost:8485/recipes/read
```





__datasetd__ is intended to be combined with other services like Solr 8.9.
__datasetd__ only implements the simplest of object storage.
