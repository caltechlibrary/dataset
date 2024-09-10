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
	"io"
	"os"
	"path"
	"strings"

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

__{app_name}__ is a tool to support SQL queries of dataset collections. 
Pairtree based collections should be index before trying to query them
(see '-index' option below). Pairtree collections use the SQLite 3
dialect of SQL for querying.  For collections using a SQL storage
engine (e.g. SQLite3, Postgres and MySQL), the SQL dialect reflects
the SQL of the storage engine.

The schema is the same for all storage engines.  The scheme for the JSON
stored documents have a four column scheme.  The columns are "_key", 
"created", "updated" and "src". "_key" is a string (aka VARCHAR),
"created" and "updated" are timestamps while "src" is a JSON column holding
the JSON document. The table name reflects the collection
name without the ".ds" extension (e.g. data.ds is stored in a database called
data having a table also called data).

The output of __{app_name}__ is a JSON array of objects. The order of the
objects is determined by the your SQL statement and SQL engine. There
is an option to generate a 2D grid of values in JSON, CSV or YAML formats.
See OPTIONS for details.

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

-help
: display help

-license
: display license

-version
: display version

-pretty
: pretty print the resulting JSON array

-sql SQL_FILENAME
: read SQL from a file. If filename is "-" then read SQL from standard input.

-grid STRING_OF_ATTRIBUTE_NAMES
: Returns list as a 2D grid of values. This options requires a comma delimited
string of attribute names for the outer object to include in grid output. It
can be combined with -pretty options.

-csv STRING_OF_ATTRIBUTE_NAMES
: Like -grid this takes our list of dataset objects and a list of attribute
names but rather than create a 2D JSON array of values it creates CSV 
representation with the first row as the attribute names.

-yaml STRING_OF_ATTRIBUTE_NAMES
: Like -grid this takes our list of dataset objects and a list of attribute
names but rather than create a 2D JSON of values it creates YAML 
representation.

-index
: This will create a SQLite3 index for a collection. This enables {app_name}
to query pairtree collections using SQLite3 SQL dialect just as it would for
SQL storage collections (i.e. don't use with postgres, mysql or sqlite based
dataset collections. It is not needed for them). Note the index is always
built before executing the SQL statement.

# EXAMPLES

Generate a list of JSON objects with the ` + "`" + `_key` + "`" + ` value
merged with the object stored as the ` + "`" + `._Key` + "`" + ` attribute.
The colllection name "data.ds" which is implemented using Postgres
as the JSON store. (note: in Postgres the ` + "`" + `||` + "`" + ` is very helpful).

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
	pretty, ptIndex, grid, csv, asYaml := false, false, "", "", ""
	sqlFName := ""

	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&pretty, "pretty", false, "pretty JSON output")
	flag.StringVar(&grid, "grid", grid, "return JSON grid of values, requires a comma delimited string of attribute names")
	flag.StringVar(&csv, "csv", csv, "return csv file using the attribute names from list of objects")
	flag.StringVar(&asYaml, "yaml", asYaml, "return YAML file using the attribute names from list of objects")
	flag.StringVar(&sqlFName, "sql", sqlFName, "read SQL statement from a file")
	flag.BoolVar(&ptIndex, "index", ptIndex, "create a SQLite 3 'index' for a collection.")
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

	// Create a DSQuery object and evaluate the command line options
	app := new(dataset.DSQuery)
	cName, stmt, params := "", "", []string{}
	if sqlFName != "" {
		in := os.Stdin
		if sqlFName != "-" {
			var err error
			in, err = os.Open(sqlFName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(10)
			}
			defer in.Close()
		}
		src, err := io.ReadAll(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(10)
		}
		stmt = fmt.Sprintf("%s", src)
	}
	for _, arg := range args {
		switch {
		case cName == "":
			cName = arg
		case stmt == "":
			stmt = arg
		default:
			params = append(params, arg)
		}
	}
	if cName == "" {
		fmt.Fprintf(os.Stderr, "missing C_NAME\n")
		os.Exit(10)
	}
	if stmt == "" {
		fmt.Fprintf(os.Stderr, "missing SQL_STATEMENT\n")
		os.Exit(10)
	}
	app.Pretty = pretty
	app.PTIndex = ptIndex
	attributes := []string{}
	if grid != "" {
		app.AsGrid = true
		app.AsCSV = false
		app.AsYAML = false
		attributes = strings.Split(grid, ",")
	}
	if csv != "" {
		app.AsGrid = false
		app.AsCSV = true
		app.AsYAML = false
		attributes = strings.Split(csv, ",")
	}
	if asYaml != "" {
		app.AsGrid = false
		app.AsCSV = false
		app.AsYAML = true
		attributes = strings.Split(asYaml, ",")
	}
	if app.AsGrid || app.AsCSV || app.AsYAML { 
		app.Attributes = []string{}
		for _, attr := range attributes {
			app.Attributes = append(app.Attributes, strings.TrimSpace(attr))
		}
	}
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr, cName, stmt, params); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(10)
	}
}
