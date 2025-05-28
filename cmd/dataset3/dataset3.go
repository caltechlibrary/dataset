// dataset is a command line tool, Go package, shared library and Python package for working with JSON objects as collections on local disc.
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

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v3"
)

const (
	helpText = `%{app_name}(1) user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}

# NAME

{app_name} 

# SYNOPSIS

{app_name} [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

{app_name} command line interface supports creating JSON object
collections and managing the JSON object documents in a collection.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

# SUPPORTED VERBS

help
: will give documentation of help on a verb, e.g. "help create"

init C_NAME
: Initialize a new dataset collection named with C_NAME.

create [OPTION] C_NAME KEY
: creates a new JSON document in the collection

read [OPTION] C_NAME KEY
: retrieves the "current" version of a JSON document from 
  the collection writing it standard out

update [OPTION] C_NAME KEY
: updates a JSON document in the collection

delete C_NAME KEY
: removes all versions of a JSON document from the collection

keys [OPTION] C_NAME
: returns a list of keys in the collection

codemeta C_NAME [PATH_TO_NEW_CODEMETA_JSON]
: displays an existing codemetada.json file for the collections or
if an optional path to a new codemeta.json it copies the file and updates the 
collections metadata.

dump C_NAME
: This will write out all dataset collection records in a JSONL document.
JSONL shows on JSON object per line, see https://jsonlines.org for details.
The object rendered will have two attributes, "key" and "object". The
key corresponds to the dataset collection key and the object is the JSON
value retrieved from the collection.

load [OPTION] C_NAME
: This will read JSON objects one per line from standard input. This
format is often called JSONL, see https://jsonlines.org. The object
has two attributes, key and object. 

A word about "keys". {app_name} uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys are composed as lower case 
alpha numeric characters but may include period, dash and underscore.
While keys maybe provided in upper and lower case they are always
converted to lowercase internally.

There are three "GLOBAL_OPTIONS" in v3 of {app_name}.T hey are 
` + "`" + `-version` + "`" + `, ` + "`" + `-help` + "`" + `
and ` + "`" + `-license` + "`" + `. All other options come
after the verb and apply to the specific action the verb
implements.

# STORAGE TYPE

There are currently three support storage options for JSON documents in a dataset collection.

- SQLite3 database (default),

The following storage engines were removed in v3 -- pairtree, MySQL and PostgreSQL. If you need
to migrate data from a v2 dataset instance use the dump verb. Then you can import the data
into the v3 dataset collection using the load verb.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
   {app_name} help init

   {app_name} init my_objects.ds 

   {app_name} help create

   {app_name} create my_objects.ds "123" '{"one": 1}'

   {app_name} create my_objects.ds "234" mydata.json 
   
   cat <<EOT | {app_name} create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   {app_name} update my_objects.ds "123" '{"one": 1, "two": 2}'

   {app_name} delete my_objects.ds "345"

   {app_name} keys my_objects.ds
~~~

This is an example of initializing a JSON documentation
collection.

~~~
{app_name} init '${C_NAME}'
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

{app_name} {version}

`
)

var (
	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
)

func fmtTxt(src string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(src, "{app_name}", appName), "{version}", version)
}

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", dataset.FmtHelp(helpText, appName, dataset.Version, dataset.ReleaseDate, dataset.ReleaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", dataset.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, dataset.Version, dataset.ReleaseHash)
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Fprintf(eout, "%s\n", fmtTxt(helpText, appName, dataset.Version))
		os.Exit(1)
	}

	// Application Logic
	err := dataset.RunCLI(in, out, eout, args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
