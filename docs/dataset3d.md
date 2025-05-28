
Dataset3d
=========

Overview
--------

__dataset3d__ is a minimal web service intended to run on localhost port 8485 (by default). It presents one or more dataset collections as a web service. It features a subset of functionallity available with the dataset command line program. __dataset3d__ does support multi-process/asynchronous update to a dataset collection.

__dataset3d__ is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a standalone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL data management solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). __dataset3d__ is a simple, miminal service.

NOTE: You could run __dataset3d__ could be combined with a front end web service like Apache 2 or NginX and through them provide access control based on _datasetd_'s predictable URL paths. That would require a robust understanding of the front end web server, it's access control mechanisms and how to defend a proxied service. That is beyond the skope of this project.

Configuration
-------------

__dataset3d__ can make one or more dataset collections visible over HTTP. The dataset collections hosted need to be avialable on the same file system as where __dataset3d__ is running. __dataset3d__ is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line to a appropriate JSON settings.

Configuration to The "settings.yaml" file has the following structure

```yaml
host: localhost:8485
dsn_url: mysql://DB_USER:DB_PASSWORD\@DB_NAME
collections
  - dataset: <PATH_TO_DATASET_COLLECTION>
    keys: true
    create: true
    read: true
    update: true
    delete: false
```

In the "collections" object the "<COLLECTION_ID>" is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at "<PATH_TO_DATASET_COLLECTION>". For each collection you can allow the following sub-paths for JSON object interaction "keys", "create", "read", "update" and "delete". JSON document attachments are supported by "attach", "retrieve", "prune". If any of these attributes are missing from the settings they are assumed to be set to false.

The sub-paths correspond to their counter parts in the dataset command line tool. By varying the settings of these you can support read only collections, drop off collections or function as a object store running behind a web application.

Running datasetd
----------------

__dataset3d__ runs as a HTTP service and as such can be exploited in the same manner as other services using HTTP.  You should only run __dataset3d__ on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by __dataset3d__ regardless of the file permissions they may in their account.

Example: If all dataset collections are in a directory only allowed access to be the "web-data" user but another users on the machine have access to curl they can access the dataset collections based on the rights of the "web-data" user by access the HTTP service.  This is a typical situation for most localhost based web services and you need to be aware of it if you choose to run __dataset3d__.

__dataset3d__ should NOT be used to store confidential, sensitive or secret information.


Supported Features
------------------

__dataset3d__ provides a limitted subset of actions supportted by the standard datset command line tool. It only supports the following actions

- collections (return a list of collections available)
- collection (return the codemeta for a collection)
- keys (return the list of keys in a collection)
- haskeys (return the list of keys in a collection)
- object (CRUD operations on a JSON document via REST calls)
- query

Each of theses "actions" can be restricted in the configuration (i.e. "settings.yaml" file) by setting the value to "false". If the attribute for the action is not specified in the JSON settings  file then it is assumed to be "false".

Exploring the End points
------------------------

The following end points are planned for _datasetd_ in version 2.

- `/collections` returns a list of available collections.
- `/collection/<COLLECTION_ID>` with an HTTP GET returns the codemeta document describing the collection.

The following end points are per collection. They are available for each
collection where the settings are set to true. The end points 
are generally RESTful so one end point will often map to a CRUD style
operations via http methods POST to create an object, GET to "read" or retrieve an object, a PUT to update an object and DELETE to remove it.

The terms "<COLLECTION_ID>" and "<KEY>" refer to the collection path, the
string representing the "key" to a JSON document. For attachment then a
base filename is used to identify the attachment associate with a "key"
in a collection.

- `/<COLLECTION_ID>/keys` returns a list of keys available in the collection
- `/<COLLECTION_ID>/haskey/<KEY>` returns true if a key is found for a JSON document or false otherwise
- `/<COLLECTION_ID>/object/<KEY>` performs CRUD operations on a JSON document, a GET retrieves the JSON document, a POST creates it, PUT updates it and DELETE removes it.
- `/<COLLECTION_ID>/query/<QUERY_NAME>` this will let you retrieve a query results


"Recipes", a use case
---------------------

In this use case a dataset collection called "recipes.ds" has been previously created and populated using the command line tool.

If I have a settings file for "recipes" based on the collection "recipes.ds" and want to make it read only I would make the attribute "read" set to true and if I want the option of listing the keys in the collection I would set that true also.

1. create our collection called "recipes.ds"
2. load test data from [recipes.jsonl](recipes.jsonl) dump file
3. create the [recipes_api.yaml](recipes_api.yaml) YAML file
4. test our web service

### Create our "recipes.ds" collection

~~~shell
dataset3 init recipes.ds
~~~

### Load test data

~~~shell
dataset3 load recipes.ds <recipes.jsonl
~~~

### Create "recipies_api.yaml"

Here is the contents of  [recipes_api.yaml](recipes_api.yaml)

~~~yaml
#
# recipes_api.yaml example configuration
#
host: localhost:8485
# This htdocs directory is provided by cold_ui so we don't enable it.
#htdocs: htdocs
collections:
  # Each collection is an object. The path prefix is
  # /api/<dataset_name>/...
  - dataset: recipes.ds
    query:
      recipes_by_namet: |
        select json_object(
            'name', src->'name',
            'ingredients', src->'ingredients',
            'preparations', src->'preparaions'
        ) as src
        from recipies
        order by src->'name'
    keys: true
    create: true
    read: true
    update: true
~~~

### Test run our collection using __dataset3d__


I would start __dataset3d__ with the following command line.

```shell
    dataset3d start recipes_api.yaml
```

This would display the start up message and log output of the service.

In another shell session I could then use curl to list the keys and read a record. In this example I assume that "waffles" is a JSON document in dataset collection "recipes.ds".

```shell
    curl http://localhost:8485/recipies/read/waffles
```

This would return the "waffles" JSON document or a 404 error if the
document was not found.

Listing the keys for "recipes.ds" could be done with this curl command.

```shell
    curl http://localhost:8485/recipies/keys
```

This would return a list of keys, one per line. You could show
all JSON documents in the collection be retrieving a list of keys
and iterating over them using curl. Here's a simple example in Bash.

```shell
    for KEY in $(curl http://localhost:8485/recipes/keys); do
       curl "http://localhost/8485/recipe/read/${KEY}"
    done
```

Add a new JSON object to a collection.

```shell
    KEY="sunday"
    curl -X POST -H 'Content-Type:application/json' \
        "http://localhost/8485/recipe/create/${KEY}" \
     -d '{"name": "Sunday", "ingredients":{"banana": 1,"ice cream": "1/4 qt", "chocalate syrup": "2 Tbps"}}'
```

