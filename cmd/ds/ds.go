//
// dataset is a command line tool, Go package, shared library and Python package for working with JSON objects as collections on disc, in an S3 bucket or in Cloud Storage
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
//
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
)

var (
	synopsis = `
_dataset_ is a command line tool for working with JSON objects as 
collections on disc, in an S3 bucket or in Cloud Storage.
`

	description = `
The [dataset](docs/dataset.html) command line tool supports 
common data management operations for JSON objects stored
as collections.  

Features:

- Basic storage actions (*create*, *read*, *update* and *delete*)
- Listing of collection *keys* (including filtering and sorting)
- Import/Export to/from CSV and Google Sheets
- An experimental full text search interface
- The ability to reshape data by performing simple object *joins*
- The ability to create data *grids* and *frames* from 
  keys lists and "dot paths" using a collections' JSON objects

Limitations:

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store 
  (it's files on "disc" without any locking)
`

	examples = `
Below is a simple example of shell based interaction with dataset 
a collection using the command line dataset tool.

` + "```shell" + `
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    dataset init friends.ds
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    dataset friends.ds create frieda '{"name":"frieda","email":"frieda@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    dataset friends.ds read frieda
    
    # Path to JSON document
    dataset friends.ds path frieda

    # Update a JSON document
    dataset friends.ds update frieda '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset friends.ds keys

    # Get keys filtered for the name "frieda"
    dataset friends.ds keys '(eq .name "frieda")'

    # Join frieda-profile.json with "frieda" adding unique key/value pairs
    dataset friends.ds join append frieda frieda-profile.json

    # Join frieda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from frieda-profile.json
    dataset friends.ds join overwrite frieda frieda-profile.json

    # Delete a JSON document
    dataset friends.ds delete frieda

    # Import data from a CSV file using column 1 as key
    dataset -quiet -nl=false friends.ds import-csv my-data.csv 1

    # To remove the collection just use the Unix shell command
    rm -fR friends.ds
` + "```" + `

`

	bugs = `
_dataset_ is NOT multi-user and doesn't have file locking abilities.
This means if you have multiple processing running _dataset_ on
the same collection doing writes you'll probably have corruption
before too long.
`

	license = `
Copyright (c) 2018, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	newLine          bool
	quiet            bool
	prettyPrint      bool
	generateMarkdown bool
	generateManPage  bool

	// Application Options
)

func main() {
	//FIXME: Replace with your base package .Version attribute
	app := cli.NewCli("v0.0.0")
	//FIXME: if you need the app name then...
	//appName := app.AppName()

	// Add Help Docs
	app.SectionNo = 1 // The manual page section number
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))
	app.AddHelp("bugs", []byte(bugs))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "output manpage markup")

	// Application Options
	//FIXME: Add any application specific options

	// Application Verbs
	//FIXME: If the application is verb based add your verbs here
	//(e.g. app.NewVerb(STRING_VERB, STRING_DESCRIPTION, FUNC_POINTER)

	// We're ready to process args
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	// Application Logic
	exitCode := app.Run(args)
	if exitCode != 0 {
		os.Exit(exitCode)
	}

	if newLine {
		fmt.Fprintln(app.Out, "")
	}
}
