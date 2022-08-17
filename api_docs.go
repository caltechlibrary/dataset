// api_docs.go is a part of dataset
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package dataset

var (
	//
	// documentation for running Daemon
	//
	WebDescription = `
USAGE

   {app_name} SETTINGS_FILE

This runs a {app_name} daemon providing a web service for interacting
with collections of JSON documents. Most paths take the form of 
The end points provided by the web service are formed in the following
steps.

  1. collection name
  2. "rest"
  3. "keys", "docs", "frames", "attachments"
  4. "key", "frame_name", "key/attachment_name"

The CRUD operations are mapped to the HTTP method per the REST conventions.

`

	WebExamples = `
EXAMPLES

Assuming the service is running on localhost:8485 and you have cURL
installed you could do the following to access the web service
implemented by {app_name}.

	# list the collections avialable
    curl -X GET http://localhost:8485/_collections  

	# list the keys from my-stuff collection
	curl -X GET http://localhost:8485/my-stuff/rest/keys 

    # See if a key "1" exists in the my-stuff collection
	curl -X GET http://localhost:8485/my-stuff/rest/keys/1

	# read record "1" from my-stuff collection
	curl -X GET http://localhost:8485/my-stuff/rest/docs/1 

    # Add a codemeta.json to set the collection's metadata
	curl -X POST http://locahost:8485/my-stuff/codemeta \\
	   -F 'codemeta=@codemeta.json'           # 

    # Get the metadata for a collection
	curl -X GET http://localhost:8485/my-stuff/codemeta
 
`

	//
	// documentation for end points support by dataset Daemon
	//

	// EndPointREADME a copy of docs/datasetd.md
	EndPointREADME = `
{app_name}
==========

Overview
--------

{app_name} is a minimal web service intended to run on localhost port 8485. It presents one or more dataset collections as a web service. It features a subset of functionallity available with the dataset command line program. {app_name} does support multi-process/asynchronous update to a dataset collection.

{app_name} is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a standalone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL data management solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). {app_name} is a simple, miminal service.

NOTE: You could run {app_name} could be combined with a front end web service like Apache 2 or NginX and through them provide access control based on {app_name}'s predictable URL paths. That would require a robust understanding of the front end web server, it's access control mechanisms and how to defend a proxied service. That is beyond the skope of this project.

Configuration
-------------

{app_name} can make one or more dataset collections visible over HTTP. The dataset collections hosted need to be avialable on the same file system as where {app_name} is running. {app_name} is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line to a appropriate JSON settings.

The "settings.json" file has the following structure

` + "```" + `
    {
        "host": "localhost:8485",
		"sql_type": "mysql",
        "dsn": "<DSN_STRING>"
    }
` + "```" + `


Running {app_name}
----------------

{app_name} runs as a HTTP service and as such can be exploited in the same manner as other services using HTTP.  You should only run {app_name} on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by {app_name} regardless of the file permissions they may in their account.

Example: If all dataset collections are in a directory only allowed access to be the "web-data" user but another users on the machine have access to curl they can access the dataset collections based on the rights of the "web-data" user by access the HTTP service.  This is a typical situation for most localhost based web services and you need to be aware of it if you choose to run {app_name}.

{app_name} should NOT be used to store confidential, sensitive or secret information.


Supported Features
------------------


{app_name} provides a RESTful web service for accessing a collection's
metdata, keys, documents, frames and attachments. The form of the path
is generally '/rest/<COLLECTION_NAME>/<DOC_TYPE>/<ID>[/<NAME>]'. REST
maps the CRUD operations to POST (create), GET (read), PUT (update),
and DELETE (delete). There are four general types of objects in
a dataset collection
  
  1. keys (point to a JSON document, these are unique identifiers)
  2. docs are the JSON documents
  3. frames hold data frames (aggregation's of an collection's content)
  4. attachments hold files attached to JSON documents

Additionally you can list all the collections available in the web service
as well as collection level metadata (as a codemeta.json document).

Collections can have their CRUD operations turned on or off based on
the columns set in the "_collections" table of the database hosting
the web service.


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

I would start {app_name} with the following command line.

` + "```" + `shell
    {app_name} settings.json
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

{app_name} provide documentation as plain text output via request
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

The following end points are supported by {app_name}

- '/' returns documentation for {app_name}
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

	EndPointCollections = `
Collections (end point)
=======================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

This provides a JSON list of collections available from the running _{app_name}_ service.

Example
=======

The assumption is that we have _{app_name}_ running on port "8485" of "localhost" and a set of collections, "t1" and "t2", defined in the "settings.json" used at launch.

` + "```" + `{.json}
    [
      "t1",
      "t2"
    ]
` + "```" + `

`

	EndPointCollection = `
Collection (end point)
=======================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

This provides a metadata as JSON for a specific collection. This may including attributes like authorship, funding and contributions.

If this end point is request with a GET method then the data is returned, if requested with a POST method the date is updated the updated metadata returned. The POST must submit JSON encoded object with the mime type of "application/json".

The metadata fields are

- "dataset" (string, semver, version of dataset managing collection)
- "name" (string) name of dataset collection
- "contact" (string) free format contact info
- "description" (string) 
- "doi" (string) a DOI assigned to the collection
- "created" (string) a date string in RFC1123 format
- "version" (string) the collection's version as a semver
- "author" (array of PersonOrOrg) a list of authors of the collection
- "contributor" (array of PersonOrOrg) a list of contributors to a collection
- "funder" (array of PersonOrOrg) a list of funders of the collection
- "annotations" (an object) this is a map to any ad-hoc fields for the collection's metadata

The PersonOrOrg structure holds the metadata for either a person or
organization. This is inspired by codemeta's peron or organization object
scheme. For a person you'd have a structure like

- "@type" (the string "Person")
- "@id" (string) the person's ORCID
- "givenName" (string) person's given name
- "familyName" (string) person's family name
- "affiliation" (array of PersonOrOrg) an list of affiliated organizations

For an organization structure like

- "@type" (the string "Organization")
- "@id" (string) the orgnization's ROR
- "name" (string) name of organization

Example
=======

The assumption is that we have _{app_name}_ running on port "8485" of "localhost" and a collection named characters is defined in the "settings.json" used at launch.

Retrieving metatadata

` + "```" + `{.shell}
    curl -X GET https://localhost:8485/collection/characters
` + "```" + `

This would return the metadata found for our characters' collection.

` + "```" + `
    {
        "dataset_version": "v0.1.10",
        "name": "characters.ds",
        "created": "2021-07-28T11:32:36-07:00",
        "version": "v0.0.0",
        "author": [
            {
                "@type": "Person",
                "@id": "https://orcid.org/0000-0000-0000-0000",
                "givenName": "Jane",
                "familyName": "Doe",
                "affiliation": [
                    {
                        "@type": "Organization",
                        "@id": "https://ror.org/05dxps055",
                        "name": "California Institute of Technology"
                    }
                ]
            }
        ],
        "contributor": [
            {
                "@type": "Person",
                "givenName": "Martha",
                "familyName": "Doe",
                "affiliation": [
                    {
                        "@type": "Organization",
                        "@id": "https://ror.org/05dxps055",
                        "name": "California Institute of Technology"
                    }
                ]
            }
        ],
        "funder": [
            {
                "@type": "Organization",
                "name": "Caltech Library"
            }
        ],
        "annotation": {
            "award": "00000000000000001-2021"
        }
    }
` + "```" + `

Update metadata requires a POST with content type "application/json". In
this example the dataset collection is named "t1" only the "name" and
"dataset_version" set.

` + "```" + `{.shell}
    curl -X POST -H 'Content-Type: application/json' \
    http://localhost:8485/collection/t1 \
    -d '{"author":[{"@type":"Person","givenName":"Jane","familyName":"Doe"}]}'
` + "```" + `

The curl calls returns

` + "```" + `{.json}
    {
        "dataset_version": "1.0.1",
        "name": "T1.ds",
        "author": [
            {
                "@type": "Person",
                "givenName": "Robert",
                "familyName": "Doiel"
            }
        ]
    }
` + "```" + `

`

	EndPointKeys = `
Keys (end point)
================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointDocument = `
Create (end point)
==================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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
	EndPointRead = `
Read (end point)
================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointUpdate = `
Update (end point)
==================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointDelete = `
Delete (end point)
==================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointAttach = `
Attach (end point)
==================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointRetrieve = `
Retrieve (end point)
====================

Interacting with the _{app_name}_ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

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

	EndPointPrune = `
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

	//
	// datasetd provide collection access/management via a simple
	// HTTP/HTTPS API
	//

	WEBDescription = `
USAGE
=====

	{app_name} SETTINGS_FILENAME

SYNPOSIS
--------

{app_name} is a web service for serving dataset collections
via HTTP/HTTPS.

DETAIL
------

{app_name} is a minimal web service typically run on localhost port 8485
that exposes a dataset collection as a web service. It features a subset of functionality available with the dataset command line program. {app_name} does support multi-process/asynchronous update to a dataset collection. 

{app_name} is notable in what it does not provide. It does not provide user/role access restrictions to a collection. It is not intended to be a stand alone web service on the public internet or local area network. It does not provide support for search or complex querying. If you need these features I suggest looking at existing mature NoSQL style solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects). {app_name} is a simple, miminal service.

NOTE: You could run {app_name} with access control based on a set of set of URL paths by running {app_name} behind a full feature web server like Apache 2 or NginX but that is beyond the skope of this project.

Configuration
-------------

{app_name} can make one or more dataset collections visible over HTTP/HTTPS. The dataset collections hosted need to be avialable on the same file system as where {app_name} is running. {app_name} is configured by reading a "settings.json" file in either the local directory where it is launch or by a specified directory on the command line.  

The "settings.json" file has the following structure

` + "```" + `
    {
        "host": "localhost:8483",
		"sql_type": "mysql",
		"dsn": "DB_USER:DB_PASSWORD3@/DB_NAME"
    }
` + "```" + `

The "host" is the URL listened to by the dataset daemon, the "sql_type" is
usually "mysql" though could be "sqlite", the "dsn" is the data source name
used to initialized the connection to the SQL engine. It is SQL engine
specific. E.g. if "sql_type" is "sqlite" then the "dsn" might be "file:DB_NAME?cache=shared".

Running {app_name}
------------------

{app_name} runs as a HTTP/HTTPS service and as such can be exploit as other network based services can be.  It is recommend you only run {app_name} on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by {app_name} regardless of the file permissions they may in their account.  E.g. If all dataset collections are in a directory only allowed access to be the "web-data" user but another user on the system can run cURL then they can access the dataset collections based on the rights of the "web-data" user.  This is a typical situation for most web services and you need to be aware of it if you choose to run {app_name}. A way to minimize the problem would be to run {app_name} in a container restricted to the specific user.

Supported Features
------------------

{app_name} provide a limitted subset of actions support by the standard datset command line tool. It only supports the following verbs

1. init (create a new collection SQL based collection)
2. keys (return a list of all keys in the collection)
3. has-key (return true if has key false otherwise)
4. create (create a new JSON document in the collection)
5. read (read a JSON document from a collection)
6. update (update a JSON document in the collection)
7. delete (delete a JSON document in the collection)
8. frames (list frames available)
9. frame (define a frame)
10. frame-def (show frame definition)
11. frame-objects (return list of framed objects)
12. refresh (refresh all the objects in a frame)
13. reframe (update the definition and reload the frame)
14. delete-frame (remove the frame)
15. has-frame (returns true if frame exists, false otherwise)
16. codemeta (imports a codemeta JSON file providing collection metadata)

Each of theses "actions" can be restricted in the _collections table (
) by setting the value to "false". If the
attribute for the action is not specified in the JSON settings file
then it is assumed to be "false".

Working with {app_name}
---------------------

E.g. if I have a settings file for "recipes" based on the collection
"recipes.ds" and want to make it read only I would make the attribute
"read" set to true and if I want the option of listing the keys in the collection I would set that true also.

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

I would start {app_name} with the following command line.

    {app_name} settings.json

This would display the start up message and log output of the service.

In another shell session I could then use cURL to list the keys and read
a record. In this example I assume that "waffles" is a JSON document
in dataset collection "recipes.ds".

    curl http://localhost:8485/recipies/read/waffles

This would return the "waffles" JSON document or a 404 error if the 
document was not found.

Listing the keys for "recipes.ds" could be done with this cURL command.

    curl http://localhost:8485/recipies/keys

This would return a list of keys, one per line. You could show
all JSON documents in the collection be retrieving a list of keys
and iterating over them using cURL. Here's a simple example in Bash.

    for KEY in $(curl http://localhost:8485/recipes/keys); do
       curl "http://localhost/8485/recipe/read/${KEY}"
    done

Access Documentation
--------------------

{app_name} provide documentation as plain text output via request
to the service end points without parameters. Continuing with our
"recipes" example. Try the following URLs with cURL.

    curl http://localhost:8485
    curl http://localhost:8485/recipes
    curl http://localhost:8485/recipes/read


{app_name} is intended to be combined with other services like Solr 8.9.
{app_name} only implements the simplest of object storage.

`

	examples = `

EXAMPLE

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the {app_name} 
service to interact with the collection via cURL. 

To create the dataset collection we use the "dataset" command and the
"vi" text edit (use can use your favorite text editor instead of vi).

    dataset init t1.ds
	vi settings.json

In the "setttings.json" file the JSON should look like.

    {
		"host": "localhost:8485",
		"sql_type": "mysql",
		"dsn": "DB_USER:DB_PASSWORD@/DB_NAME"
	}

Now we can run {app_name} and make the dataset collection available
via HTTP.

    {app_name} settings.json

You should now see the start up message and any log information display
to the console. You should open a new shell sessions and try the following.

We can now use cURL to post the document to the "/t1/create/one" end
point. 

    curl -X POST http://localhost:8485/t1/create/one \
	    -d '{"one": 1}'

Now we can list the keys available in our collection.

    curl http://localhost:8485/t1/keys

We should see "one" in the response. If so we can try reading it.

    curl http://localhost:8485/t1/read/one

That should display our JSON document. Let's try updating (replacing)
it. 

    curl -X POST http://localhost:8485/t1/update/one \
	    -d '{"one": 1, "two": 2}'

If you read it back you should see the updated record. Now lets try
deleting it.

	curl http://localhost:8485/t1/delete/one

List the keys and you should see that "one" is not longer there.

    curl http://localhost:8485/t1/keys

In the shell session where {app_name} is running press "ctr-C"
to terminate the service.

`
)
