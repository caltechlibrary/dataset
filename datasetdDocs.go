package dataset

import ()

//
// End point documentation for datasetd.go
//

// readmeDocument() provides a copy of docs/datasetd.md
func readmeDocument() string {
	return `
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

` + "```" + `
    {
        "host": "localhost:8485",
        "collections": {
            '<COLLECTION_ID>': {
                "dataset": "<PATH_TO_DATASET_COLLECTION>",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": false,
                "attach": false,
                "retrieve": false,
                "prune": false
            }
        }
    }
` + "```" + `

In the "collections" object the '<COLLECTION_ID>' is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at '<PATH_TO_DATASET_COLLECTION>'. For each collection you can allow the following sub-paths for JSON object interaction "keys", "create", "read", "update" and "delete". JSON document attachments are supported by "attach", "retrieve", "prune". If any of these attributes are missing from the settings they are assumed to be set to false.

The sub-paths correspond to their counter parts in the dataset command line tool. By varying the settings of these you can support read only collections, drop off collections or function as a object store running behind a web application.

Running datasetd
----------------

_datasetd_ runs as a HTTP service and as such can be exploited in the same manner as other services using HTTP.  You should only run _datasetd_ on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by _datasetd_ regardless of the file permissions they may in their account.

Example: If all dataset collections are in a directory only allowed access to be the "web-data" user but another users on the machine have access to curl they can access the dataset collections based on the rights of the "web-data" user by access the HTTP service.  This is a typical situation for most localhost based web services and you need to be aware of it if you choose to run _datasetd_.

_datasetd_ should NOT be used to store confidential, sensitive or secret information.


Supported Features
------------------

_datasetd_ provides a limitted subset of actions supportted by the standard datset command line tool. It only supports the following verbs

1. keys (return a list of all keys in the collection)
    - must be a GET request
2. create (create a new JSON document in the collection)
    - must be a POST request ended as JSON with a content type of "application/json"
3. read (read a JSON document from a collection)
    - must be a GET request
4. update (update a JSON document in the collection)
    - must be a POST request ended as JSON with a content type of "application/json"
5. delete (delete a JSON document in the collection)
    - must be a GET request
6. collections (list as a JSON array of objects the collections avialable)
    - must be a GET request
7. attach allows you to upload via a POST (not JSON encoded) an attachment to a JSON document. The attachment is limited in size to 250 MiB. The POST must be a multi-part encoded web form where the upload name is identified as "filename" in the form and the URL path identifies the name to use for the saved attachment.
8. retrieve allows you to download an versioned attachment from a JSON document
9. prune removes versioned attachments from a JSON document

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

` + "```" + `
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
` + "```" + `

I would start _datasetd_ with the following command line.

` + "```" + `shell
    datasetd settings.json
` + "```" + `

This would display the start up message and log output of the service.

In another shell session I could then use curl to list the keys and read
a record. In this example I assume that "waffles" is a JSON document
in dataset collection "recipes.ds".

` + "```" + `shell
    curl http://localhost:8485/recipies/read/waffles
` + "```" + `

This would return the "waffles" JSON document or a 404 error if the
document was not found.

Listing the keys for "recipes.ds" could be done with this curl command.

` + "```" + `shell
    curl http://localhost:8485/recipies/keys
` + "```" + `

This would return a list of keys, one per line. You could show
all JSON documents in the collection be retrieving a list of keys
and iterating over them using curl. Here's a simple example in Bash.

` + "```" + `shell
    for KEY in $(curl http://localhost:8485/recipes/keys); do
       curl "http://localhost/8485/recipe/read/${KEY}"
    done
` + "```" + `

Add a new JSON object to a collection.

` + "```" + `shell
    KEY="sunday"
    curl -X POST -H 'Content-Type:application/json' \
        "http://localhost/8485/recipe/create/${KEY}" \
     -d '{"ingredients":["banana","ice cream","chocalate syrup"]}'
` + "```" + `

Online Documentation
--------------------

_datasetd_ provide documentation as plain text output via request
to the service end points without parameters. Continuing with our
"recipes" example. Try the following URLs with curl.

` + "```" + `
    curl http://localhost:8485
    curl http://localhost:8485/recipes
    curl http://localhost:8485/recipes/create
    curl http://localhost:8485/recipes/read
    curl http://localhost:8485/recipes/update
    curl http://localhost:8485/recipes/delete
    curl http://localhost:8485/recipes/attach
    curl http://localhost:8485/recipes/retrieve
    curl http://localhost:8485/recipes/prune
` + "```" + `

End points
----------

The following end points are supported by _datasetd_

- '/' returns documentation for _datasetd_
- '/collections' returns a list of available collections.

The following end points are per colelction. They are available
for each collection where the settings are set to true. Some end points require POST HTTP method and specific content types.

The terms '<COLLECTION_ID>', '<KEY>' and '<SEMVER>' refer to
the collection path, the string representing the "key" to a JSON document and semantic version number for attachment. Unless specified
end points support the GET method exclusively.

- '/<COLLECTION_ID>' returns general dataset documentation with some tailoring to the collection.
- '/<COLLECTION_ID>/keys' returns a list of keys available in the collection
- '/<COLLECTION_ID>/create' returns documentation on the 'create' end point
- '/<COLLECTION_IO>/create/<KEY>' requires the POST method with content type header of 'application/json'. It can accept JSON document up to 1 MiB in size. It will create a new JSON document in the collection or return an HTTP error if that fails
- '/<COLLECTION_ID>/read' returns documentation on the 'read' end point
- '/<COLLECTION_ID>/read/<KEY>' returns a JSON object for key or a HTTP error
- '/<COLLECTION_ID>/update' returns documentation on the 'update' end point
- '/COLLECTION_ID>/update/<KEY>' requires the POST method with content type header of 'application/json'. It can accept JSON document up to 1 MiB is size. It will replace an existing document in the collection or return an HTTP error if that fails
- '/<COLLECTION_ID>/delete' returns documentation on the 'delete' end point
- '/COLLECTION_ID>/delete/<KEY>' requires the GET method. It will delete a JSON document for the key provided or return an HTTP error
- '/<COLLECTION_ID>/attach' returns documentation on attaching a file to a JSON document in the collection.
- '/COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>' requires a POST method and expects a multi-part web form providing the filename in the 'filename' field. The <FILENAME> in the URL is used in storing the file. The document will be written the JSON document directory by '<KEY>' in sub directory indicated by '<SEMVER>'. See https://semver.org/ for more information on semantic version numbers.
- '/<COLLECTION_ID>/retrieve' returns documentation on how to retrieve a versioned attachment from a JSON document.
- '/<COLLECTION_ID>/retrieve/<KEY>/<SEMVER>/<FILENAME>' returns the versioned attachment from a JSON document or an HTTP error if that fails
- '/<COLLECTION_ID>/prune' removes a versioned attachment from a JSON document or returns an HTTP error if that fails.
- '/<COLLECTION_ID>/prune/<KEY>/<SEMVER>/<FILENAME>' removes a versioned attachment from a JSON document.

`
}

func keysDocument(collectionID string) string {
	return `
Keys (end point)
================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

This end point lists keys available in a collection.

    'http://localhost:8485/<COLLECTION_ID>/keys'

Requires a "GET" method.

The keys are turned as a JSON array or http error if not found.

Example
-------

In this example '<COLLECTION_ID>' is "t1".

` + "```" + `{.shell}
    curl http://localhost:8485/t1/keys
` + "```" + `

The document return looks some like

` + "```" + `
    [
        "one",
        "two",
        "three"
    ]
` + "```" + `

For a "t1" containing the keys of "one", "two" and "three".

`
}

func createDocument(collectionID string) string {
	return `
Create (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Create a JSON document in the collection. Requires a unique key in the URL and the content most be JSON less than 1 MiB in size.

    'http://localhost:8485/<COLLECTION_ID>/created/<KEY>' 

Requires a "POST" HTTP method with.

Creates a JSON document for the '<KEY>' in collection '<COLLECTION_ID>'. On success it returns HTTP 201 OK. Otherwise an HTTP error if creation fails.

The "POST" needs to be JSON encoded and using a Content-Type of "application/json" in the request header.

Example
-------

The '<COLLECTION_ID>' is "t1", the '<KEY>' is "one" The content posted is

` + "```" + `{.json}
    {
       "one": 1
    }
` + "```" + `

Posting using CURL is done like

` + "```" + `shell
    curl -X POST -H 'Content-Type: application.json' \
      -d '{"one": 1}' \
      http://locahost:8485/t1/create/one 
` + "```" + `

`
}

func readDocument(collectionID string) string {
	return `
Read (end point)
================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Retrieve a JSON document from a collection.

    'http://localhost:8485/<COLLECTION_ID>/read/<KEY>'

Requires a "GET" HTTP method.

Returns the JSON document for given '<KEY>' found in '<COLLECTION_ID>' or a HTTP error if not found.

Example
-------

Curl accessing "t1" with a key of "one"

` + "```" + `{.shell}
    curl http://localhost:8485/t1/read/one
` + "```" + `

An example JSON document (this example happens to have an attachment) returned.

` + "```" + `
{
   "_Attachments": [
      {
         "checksums": {
            "0.0.1": "bb327f7bcca0f88649f1c6acfdc0920f"
         },
         "created": "2021-10-11T11:09:51-07:00",
         "href": "T1.ds/pairtree/on/e/0.0.1/a1.png",
         "modified": "2021-10-11T11:09:51-07:00",
         "name": "a1.png",
         "size": 32511,
         "sizes": {
            "0.0.1": 32511
         },
         "version": "0.0.1",
         "version_hrefs": {
            "0.0.1": "T1.ds/pairtree/on/e/0.0.1/a1.png"
         }
      }
   ],
   "_Key": "one",
   "four": "four",
   "one": 1,
   "three": 3,
   "two": 2
}
` + "```" + `

`
}

func updateDocument(collectionID string) string {
	return `
Update (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Update a JSON document in the collection. Requires a key to an existing
JSON record in the URL and the content most be JSON less than 1 MiB in size.

    'http://localhost:8485/<COLLECTION_ID>/update/<KEY>'

Requires a "POST" HTTP method.

Update a JSON document for the '<KEY>' in collection '<COLLECTION_ID>'. On success it returns HTTP 200 OK. Otherwise an HTTP error if creation fails.

The "POST" needs to be JSON encoded and using a Content-Type of "application/json" in the request header.

Example
-------

The '<COLLECTION_ID>' is "t1", the '<KEY>' is "one" The revised content posted is

` + "```" + `{.json}
    {
       "one": 1,
       "two": 2,
       "three": 3,
       "four": 4
    }
` + "```" + `

Posting using CURL is done like

` + "```" + `shell
    curl -X POST -H 'Content-Type: application.json' \
      -d '{"one": 1, "two": 2, "three": 3, "four": 4}' \
      http://locahost:8485/t1/update/one 
` + "```" + `

`
}

func deleteDocument(collectionID string) string {
	return `
Delete (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Delete a JSON document in the collection. Requires the document key and collection name.

    'http://localhost:8485/<COLLECTION_ID>/delete/<KEY>'

Requires a 'GET' HTTP method.

Deletes a JSON document for the '<KEY>' in collection '<COLLECTION_ID>'. On success it returns HTTP 200 OK. Otherwise an HTTP error if creation fails.

Example
-------

The '<COLLECTION_ID>' is "t1", the '<KEY>' is "one" The content posted is

Posting using CURL is done like

` + "```" + `shell
    curl -X GET -H 'Content-Type: application.json' \
      http://locahost:8485/t1/delete/one 
` + "```" + `

`
}

func attachDocument(collectionID string) string {
	return `
Attach (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Attaches a document to a JSON Document using '<KEY>', '<SEMVER>' and '<FILENAME>'.

    'http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>'

Requires a "POST" method. The "POST" is expected to be a multi-part web form providing the source filename in the field "filename".  The document will be written the JSON document directory by '<KEY>' in sub directory indicated by '<SEMVER>'.

See https://semver.org/ for more information on semantic version numbers.

Example
=======

In this example the '<COLLECTION_ID>' is "t1", the '<KEY>' is "one" and
the content upload is "a1.png" in the home directory "/home/jane.doe".
The '<SEMVER>' is "0.0.1".

` + "```" + `shell
    curl -X POST -H 'Content-Type: multipart/form-data' \
       -F 'filename=@/home/jane.doe/a1.png' \
       http://localhost:8485/t1/attach/one/0.0.1/a1.png
` + "```" + `

NOTE: The URL contains the filename used in the saved attachment. If
I didn't want to call it "a1.png" I could have provided a different
name in the URL path.

`
}

func retrieveDocument(collectionID string) string {
	return `
Retrieve (end point)
====================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Retrieves an s attached document from a JSON record using '<KEY>', '<SEMVER>' and '<FILENAME>'.

    'http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>'

Requires a POST method and expects a multi-part web form providing the filename. The document will be written the JSON document directory by '<KEY>' in sub directory indicated by '<SEMVER>'. 

See https://semver.org/ for more information on semantic version numbers.

Example
-------

In this example we're retieving the '<FILENAME>' of "a1.png", with the '<SEMVER>' of "0.0.1" from the '<COLLECTION_ID>' of "t1" and '<KEY>'
of "one" using curl.

` + "```" + `{.shell}
    curl http://localhost:8485/t1/retrieve/one/0.0.1/a1.png
` + "```" + `

This should trigger a download of the "a1.png" image file in the
collection for document "one".

`
}

func pruneDocument(collectionID string) string {
	return `
Prune (end point)
=================

Removes an attached document from a JSON record using '<KEY>', '<SEMVER>' and '<FILENAME>'.

    'http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>'

Requires a GET method. Returns an HTTP 200 OK on success or an HTTP error code if not.

See https://semver.org/ for more information on semantic version numbers.

Example
-------

In this example '<COLLECTION_ID>' is "t1", '<KEY>' is "one", '<SEMVER>' is "0.0.1" and '<FILENAME>' is "a1.png". Once again our example uses curl.

` + "```" + `
    curl http://localhost:8485/t1/prune/one/0.0.1/a1.png
` + "```" + `

This will cause the attached file to be removed from the record
and collection.

`
}
