package main

//
// datasetd provide collection access/management via a simple
// HTTP/HTTPS API
//

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	description = `
USAGE
=====

	{app_name} [SETTINGS_FILENAME]

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

    {
        "host": "localhost:8483",
        "collections": {
            "COLLECTION_ID": {
                "dataset": "PATH_TO_DATASET_COLLECTION",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": false
            }
        }
    }

In the "collections" object the "COLLECTION_ID" is a string which will be used as the start of the path in the URL. The "dataset" attribute sets the path to the dataset collection made available at "COLLECTION_ID". For each
collection you can allow the following sub-paths, "create", "read", "update", "delete" and "keys". These sub-paths correspond to their counter parts in the dataset command line tool. In this way would can have a
dataset collection function as a drop box, a read only list or a simple JSON
object storage service.

Running datasetd
----------------

{app_name} runs as a HTTP/HTTPS service and as such can be exploit as other network based services can be.  It is recommend you only run {app_name} on localhost on a trusted machine. If the machine is a multi-user machine all users can have access to the collections exposed by {app_name} regardless of the file permissions they may in their account.
E.g. If all dataset collections are in a directory only allowed access to be the "web-data" user but another user on the system can run cURL then they can access the dataset collections based on the rights of the "web-data" user.  This is a typical situation for most web services and you need to be aware of it if you choose to run {app_name}.

Supported Features
------------------

{app_name} provide a limitted subset of actions support by the standard datset command line tool. It only supports the following verbs

1. keys (return a list of all keys in the collection)
2. create (create a new JSON document in the collection)
3. read (read a JSON document from a collection)
4. update (update a JSON document in the collection)
5. delete (delete a JSON document in the collection)

Each of theses "actions" can be restricted in the configuration (
i.e. "settings.json" file) by setting the value to "false". If the
attribute for the action is not specified in the JSON settings file
then it is assumed to be "false".

Working with datasetd
---------------------

E.g. if I have a settings file for "recipes" based on the collection
"recipes.ds" and want to make it read only I would make the attribute
"read" set to true and if I want the option of listing the keys in the collection I would set that true also.

    {
        "host": "localhost:8485",
        "recipes": {
            "dataset": "recipes.ds",
            "keys": true,
            "read": true
        }
    }

I would start {app_name} with the following command line.

    datasetd settings.json

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
		"t1": {
			"dataset": "t1.ds",
			"keys": true,
			"create": true,
			"read": true,
			"update": true,
			"delete": true
		}
	}

Now we can run {app_name} and make the dataset collection available
via HTTP.

    datasetd settings.json

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

	license = `
{app_name} {version}

Copyright (c) 2021, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

func main() {
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		dataset.DisplayUsage(os.Stdout, appName, flagSet, description, examples, license)
		os.Exit(0)
	}

	if showLicense {
		dataset.DisplayLicense(os.Stdout, appName, license)
		os.Exit(0)
	}

	if showVersion {
		dataset.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}

	/* Looking for settings.json */
	settings := "settings.json"
	if len(args) > 0 {
		settings = args[0]
	}

	/* Iniitialize API */
	if err := dataset.InitDatasetAPI(settings); err != nil {
		fmt.Fprintf(os.Stderr, "InitDatasetAPI(%q) failed, %s", settings, err)
		os.Exit(1)
	}

	/* Run API */
	if err := dataset.RunDatasetAPI(appName); err != nil {
		fmt.Fprintf(os.Stderr, "RunDatasetAPI(%q) failed, %s", appName, err)
	}
}
