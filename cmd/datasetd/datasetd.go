// datasetd implements a web service for working with dataset collections.
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
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
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/v2"
)

const (
	helpText = `%{app_name} (1) user manual | verion {version} {release_hash}"
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] SETTINGS_JSON_FILE

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
   {app_name} settings.json
~~~

In this example we cover a short life cycle of a collection
called "t1.ds". We need to create a "settings.json" file and
an empty dataset collection. Once ready you can run the {app_name} 
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

Now we can run {app_name} and make the dataset collection available
via HTTP.

~~~
    {app_name} settings.json
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

func fmtTxt(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display detailed help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")

	flag.Parse()
	args := flag.Args()

	out := os.Stdout
	//eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, dataset.Version))
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintf(out, "%s\n", dataset.LicenseText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(out, "%s %s\n", appName, dataset.Version)
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
	if err := dataset.RunAPI(appName, settings); err != nil {
		fmt.Fprintf(os.Stderr, "RunWebAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
}
