// dsimport.go is a command line program for working an dataset collections. It
// is using to import data from a CSV file into a collection.
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

{app_name} [OPTIONS] C_NAME CSV_FILENAME KEY_COLUMN

# DESCRIPTION

__{app_name}__ is a tool to import CSV content into a dataset collection
where the column headings become the attribute names and the row values
become the attribute values.

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

CSV_FILENAME
: The name of the CSV file to import

KEY_COLUMN
: The column name to use the they object key. If none is provided then
the first column is used as the object key. Keys values must be unique.


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-comma
: Set column delimiter

-comment
: Set row comment delimiter

-overwrite
: Overwrite objects on key collision

# EXAMPLES

Import a file with three columns

- item_code
- title
- location

The "item_code" is unique for each row. The data is stored
in a file called "books.csv". We are importing the CSV file
into a collections called. "shelves.ds"

~~~
{app_name} shelves.ds books.csv "item_code"
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
	overwrite, comma, comment := false, "", ""
	lazyQuotes, trimSpace := false, false

	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&overwrite, "overwrite", overwrite, "overwrite object on key collision")
	flag.StringVar(&comma, "comma", comma, "set column delimiter")
	flag.StringVar(&comment, "comment", comment, "set row comment delimiter")
	flag.BoolVar(&lazyQuotes, "lazy", lazyQuotes, "use lazy quotes")
	flag.BoolVar(&trimSpace, "trim", trimSpace, "trim leading space")
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
	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s C_NAME CSV_FILENAME KEY_COLUMN", appName)
		os.Exit(10)
	}

	// Setup our POSIX IO options
	var err error
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	// Create a DSQuery object and evaluate the command line options
	app := new(dataset.DSImport)
	app.Overwrite = overwrite
	app.LazyQuotes = lazyQuotes
	app.TrimLeadingSpace = trimSpace
	if comma != "" {
		app.Comma = comma
	}
	if comment != "" {
		app.Comment = comment
	}
	cName, csvName, keyColumn := args[0], args[1], args[2]

	// Handle arranging our input data
	if csvName != "" && csvName != "-" {
		in, err = os.Open(csvName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(64)
		}
		defer in.Close()
	}
	if err := app.Run(in, out, eout, cName, keyColumn); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(10)
	}
}
