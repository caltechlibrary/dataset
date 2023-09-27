// dsquery.go is a command line program for working an dataset collections using the
// dataset v2 SQL store for the JSON documents (e.g. Postgres, MySQL, SQLite 3).
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
// @author Tom Morrell, <tmorrell@caltech.edu>
//
// Copyright (c) 2023, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/v2"
)

var (
	helpText = `%{app_name}(1) dataset user manual | version {version} {release_hash}
% R. S. Doiel and Tom Morrell
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] C_NAME SQL_STATEMENT [PARAMS]

# DESCRIPTION

__{app_name}__ is a tool to support SQL queries of dataset collections that
use SQL storage for the collection's JSON documents.  It takes a dataset
collection name and a sql statement returning the results. This will allow us
to improve our feeds building process by taking advantage of SQL and the
collection's SQL database engine.

The scheme for the JSON stored documents have a four column scheme. 
The columns are "_key", "created", "updated" and "src". The are stored
in a table with the same name as the database which is formed from the
C_NAME without extension (e.g. data.ds is stored in a database called
data having a table also called data).

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

SQL_STATEMENT
: The SQL statement should conform to the SQL dialect used for the
JSON store for the JSON store (e.g.  Postgres, MySQL and SQLite 3).
The SELECT clause should return a single JSON object type per row.
__{app_name}__ returns an JSON array of JSON objects returned
by the SQL query.

PARAMS
: Is optional, it is any values you want to pass to the SQL_STATEMENT.

# SQL Store Scheme

_key
: The key or id used to identify the JSON documented stored.

src
: This is a JSON column holding the JSON document

created
: The date the JSON document was created in the table

updated
: The date the JSON document was updated


# OPTIONS

help
: display help

license
: display license

version
: display version


# EXAMPLES

Generate a list of JSON objects with the `+"`"+`_key`+"`"+` value
merged with the object stored as the `+"`"+`._Key`+"`"+` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (note: in Postgres the `+"`"+`||`+"`"+` is very helpful).

~~~
{app_name} data.ds "SELECT jsonb_build_object('_Key', _key)::jsonb || src::jsonb FROM data"
~~~

In this example we're returning the "src" in our collection by querying
for a "id" attribute in the "src" column. The id is passed in as an attribute
using the Postgres positional notatation in the statement.

~~~
{app_name} data.ds "SELECT src FROM data WHERE src->>'id' = $1 LIMIT 1" "xx103-3stt9"
~~~

`
)

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following will be set when version.go is generated
	version := dataset.Version
	releaseDate := dataset.ReleaseDate
	releaseHash := dataset.ReleaseHash
	fmtHelp := dataset.FmtHelp
	pretty := false

	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&pretty, "pretty", false, "pretty JSON output")
	flag.Parse()
	args := flag.Args()

	if showHelp {
		fmt.Fprintf(os.Stdout, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(os.Stdout, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(os.Stdout, "%s\n", dataset.LicenseText)
		os.Exit(0)
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing C_NAME and SQL_STATEMENT")
		os.Exit(10)
	}
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "missing SQL_STATEMENT")
		os.Exit(10)
	}
	// Create a Ep3Util object
	app := new(dataset.DSQuery)
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing action, don't know what to do\n")
		os.Exit(10)
	}
	// To start we assume the first parameter is an action
	cName, stmt, params := args[0], args[1], []string{}
	if len(args) > 2 {
		params = args[2:]
	}
	app.Pretty = pretty
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr, cName, stmt, params); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(10)
	}
}
