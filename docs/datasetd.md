
Datasetd
========

Overview
--------

_datasetd_ is a minimal web service intended to run on localhost port 8485. It presents one or more dataset collections as a web service. It features a subset of functionallity available with the dataset command line program. _datasetd_ does support multi-process/asynchronous update to a dataset collection.

_datasetd_ is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a standalone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL data management solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). _datasetd_ is a simple, miminal service.

NOTE: You could run _datasetd_ could be combined with a front end web service like Apache 2 or NginX and through them provide access control based on _datasetd_'s predictable URL paths. That would require a robust understanding of the front end web server, it's access control mechanisms and how to defend a proxied service. That is beyond the skope of this project.

Configuration
-------------

_datasetd_ can make one or more dataset collections visible over HTTP. The dataset collections hosted need to be avialable on the same file system as where _datasetd_ is running. _datasetd_ is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line to a appropriate JSON settings.

The "settings.json" file has the following structure

```json
    {
        "host": "localhost:8485",
        "dsn_url": "mysql://DB_USER:DB_PASSWORD\@DB_NAME",
        "collections": [
            {
                "dataset": "<PATH_TO_DATASET_COLLECTION>",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": false,
                "attach": false,
                "retrieve": false,
                "prune": false,
                "frame-read": true,
                "frame-write": false
           }
        ]
    }
```

In the "collections" object the "<COLLECTION_ID>" is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at "<PATH_TO_DATASET_COLLECTION>". For each collection you can allow the following sub-paths for JSON object interaction "keys", "create", "read", "update" and "delete". JSON document attachments are supported by "attach", "retrieve", "prune". If any of these attributes are missing from the settings they are assumed to be set to false.

The sub-paths correspond to their counter parts in the dataset command line tool. By varying the settings of these you can support read only collections, drop off collections or function as a object store running behind a web application.

Running datasetd
----------------

_datasetd_ runs as a HTTP service and as such can be exploited in the same manner as other services using HTTP.  You should only run _datasetd_ on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by _datasetd_ regardless of the file permissions they may in their account.

Example: If all dataset collections are in a directory only allowed access to be the "web-data" user but another users on the machine have access to curl they can access the dataset collections based on the rights of the "web-data" user by access the HTTP service.  This is a typical situation for most localhost based web services and you need to be aware of it if you choose to run _datasetd_.

_datasetd_ should NOT be used to store confidential, sensitive or secret information.


Supported Features
------------------

_datasetd_ provides a limitted subset of actions supportted by the standard datset command line tool. It only supports the following actions

- collections (return a list of collections available)
- collection (return the codemeta for a collection)
- keys (return the list of keys in a collection)
- has-keys (return the list of keys in a collection)
- object (CRUD operations on a JSON document via REST calls)
- frames (return a list of frames available in a collection)
- has-frame (return a true if frame exists, false otherwise)
- frame (CRUD operations on a frame via REST calls)
- frame-objects (get a frame's list of objects)
- frame-keys (get a frame's list of keys)
- attachments (list attachments for a JSON document)
- attachment (CRUD operations on attachment via REST calls)

Each of theses "actions" can be restricted in the configuration (
i.e. "settings.json" file) by setting the value to "false". If the
attribute for the action is not specified in the JSON settings file
then it is assumed to be "false".

Use case
--------

In this use case a dataset collection called "recipes.ds" has been previously created and populated using the command line tool.

If I have a settings file for "recipes" based on the collection
"recipes.ds" and want to make it read only I would make the attribute
"read" set to true and if I want the option of listing the keys in the collection I would set that true also.

```json
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

I would start _datasetd_ with the following command line.

```shell
    datasetd settings.json
```

This would display the start up message and log output of the service.

In another shell session I could then use curl to list the keys and read
a record. In this example I assume that "waffles" is a JSON document
in dataset collection "recipes.ds".

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
     -d '{"ingredients":["banana","ice cream","chocalate syrup"]}'
```

End points
----------

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
- `/<COLLECTION_ID>/has-key/<KEY>` returns true if a key is found for a JSON document or false otherwise
- `/<COLLECTION_ID>/object/<KEY>` performs CRUD operations on a JSON document, a GET retrieves the JSON document, a POST creates it, PUT updates it and DELETE removes it.
- `/<COLLECTION_ID>/attachments/<KEY>` returns a list of attachments assocated with the JSON document
- `/<COLLECTION_ID>/attachment/<KEY>/<FILENAME>` allows you to perform CRUD operations on an attachment. Create is done with a POST, read (retrieval) is done wiht a GET, replacement is done with a PUT and deleting an attachment (pruning) is done with a DELETE http method.
- `/<COLLECTION_ID>/frames` list the frames defined for a collection
- `/<COLLECTION_ID>/has-frame/<FRAME_NAME>` returns true if frame is defined otherwise false
- `/<COLLECTION_ID>/frame/<FRAME_NAME>` a GET will return the frame definition, a POST will create a frame, a DELETE will remove a frame, a PUT without a body will cause the frame to be refreshed and a PUT with an array of keys will cause the frame to be reframed
- `/<COLLECTION_ID>/frame-objects/<FRAME_NAME>` will return a list of the frame's objects
- `/<COLLECTION_ID>/frame-keys/<FRAME_NAME>` will return a list of keys in the frame

