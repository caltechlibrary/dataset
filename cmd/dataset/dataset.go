// dataset is a command line tool, Go package, shared library and Python package for working with JSON objects as collections on local disc.
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

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v2"
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

init [STORAGE_TYPE]
: Initialize a new dataset collection

model
: provides an experimental interactive data model generator creating
the "model.yaml" file in the data set collection's root directory.

create
: creates a new JSON document in the collection

read
: retrieves the "current" version of a JSON document from 
  the collection writing it standard out

update
: updates a JSON document in the collection

delete
: removes all versions of a JSON document from the collection

keys
: returns a list of keys in the collection

codemeta:
: copies metadata a codemeta file and updates the 
  collections metadata

attach
: attaches a document to a JSON object record

attachments
: lists the attachments associated with a JSON object record

retrieve
: creates a copy local of an attachement in a JSON record

detach
: will copy out the attachment to a JSON document 
  into the current directory 

prune
: removes an attachment (including all versions) from a JSON record

set-versioning
: will set the versioning of a collection. The versioning
  value can be "", "none", "major", "minor", or "patch"

get-versioning
: will display the versioning setting for a collection

dump
: This will write out all dataset collection records in a JSONL document.
JSONL shows on JSON object per line, see https://jsonlines.org for details.
The object rendered will have two attributes, "key" and "object". The
key corresponds to the dataset collection key and the object is the JSON
value retrieved from the collection.

load
: This will read JSON objects one per line from standard input. This
format is often called JSONL, see https://jsonlines.org. The object
has two attributes, key and object. 

join [OPTIONS] c_name, key, JSON_SRC
: This will join a new object provided on the command line with an
existing object in the collection.


A word about "keys". {app_name} uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensistive.

A word about "GLOBAL_OPTIONS" in v2 of {app_name}.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
{app_name}).  There are three "GLOBAL_OPTIONS" that are exception
and they are ` + "`" + `-version` + "`" + `, ` + "`" + `-help` + "`" + `
and ` + "`" + `-license` + "`" + `. All other options come
after the verb and apply to the specific action the verb
implements.

# STORAGE TYPE

There are currently three support storage options for JSON documents in a dataset collection.

- SQLite3 database >= 3.40 (default)
- Postgres >= 12
- MySQL 8
- Pairtree (pre-2.1 default)

STORAGE TYPE are specified as a DSN URI except for pairtree which is just "pairtree".


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

   {app_name} model my_objects.ds

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

This is an example of initializing a Pairtree JSON documentation
collection using the environment.

~~~
{app_name} init '${C_NAME}' pairtree
~~~

In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME.

To specify Postgres as the storage for your JSON document collection.
You'd use something like --

~~~
{app_name} init '${C_NAME}' \\
  'postgres://${USER}@localhost/${DB_NAME}?sslmode=disable'
~~~


In this case '${C_NAME}' is the name of your JSON document
read from the environment varaible C_NAME. USER is used
for the Postgres username and DB_NAME is used for the Postgres
database name.  The sslmode option was specified because Postgres
in this example was restricted to localhost on a single user machine.


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
