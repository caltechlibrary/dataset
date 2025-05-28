// datasetd implements a web service for working with dataset collections.
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2025, Caltech
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
	"github.com/caltechlibrary/dataset/v3"
)

const (
	helpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SETTINGS_FILE

# DESCRIPTION

{app_name} provides a web service for one or more dataset collections. Requires the
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

The settings files provides {app_name} with the configuration
of the service web service and associated dataset collection(s).

It can be writen as either a JSON or YAML file. If it is a YAML file
you should use the ".yaml" extension so that {app_name} will correctly
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
the query (i.e. a POST To the path "`+"`"+`/api/<COLLECTION_NAME>/query/<QUERY_NAME>`+"`"+`")
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

query
: (optional, default is false) allow defined queries in the JSON API

# EXAMPLES

Starting up the web service

~~~
   {app_name} settings.yaml
~~~

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the {app_name} 
service to interact with the collection via cURL. 

To create the dataset collection we use the "dataset" command and the
"vi" text edit (use can use your favorite text editor instead of vi).

~~~
    createdb t1
    dataset3 init t1.ds \
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
    # Define a query
    query:
	  list_objects: |
	    select src
		from t1
		order by _Key
    # What follows are object level permissions
	keys: true
    create: true
    read: true
	update: true
	delete: true
EOT
~~~

Now we can run {app_name} and make the dataset collection available
via HTTP.

~~~
    {app_name} start settings.yaml
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

In the shell session where {app_name} is running press "ctr-C"
to terminate the service.


{app_name} {version}

`
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	description = dataset.WebDescription
	examples    = dataset.WebExamples
	license     = dataset.License
)

func main() {
	appName := path.Base(os.Args[0])
	version, releaseDate, releaseHash, licenseText  := dataset.Version, dataset.ReleaseDate, dataset.ReleaseHash, dataset.LicenseText
	fmtHelp := dataset.FmtHelp

	// Standard Options
	debug := false
	flag.BoolVar(&showHelp, "help", false, "display detailed help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&debug, "debug", debug, "log debugging information")

	flag.Parse()
	args := flag.Args()

	out := os.Stdout
	//eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintf(out, "%s\n", licenseText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(out, "%s %s (%s %s)\n", appName, version, releaseDate, releaseHash)
		os.Exit(0)
	}

	/* Looking for settings.json */
	settings := ""
	if len(args) > 0 {
		settings = args[0]
	}
	if _, err := os.Stat(settings); err != nil {
		fmt.Fprintf(os.Stderr, `Could not find %s

Try %s -help for usage details
`, settings, appName)
		os.Exit(1)
	}

	/* Run API */
	if err := dataset.RunAPI(appName, settings, debug); err != nil {
		fmt.Fprintf(os.Stderr, "RunAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
}
