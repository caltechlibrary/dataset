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
	"flag"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	//"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
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

	// Application Verbs
	vInit        *cli.Verb // init
	vStatus      *cli.Verb // status
	vCreate      *cli.Verb // create
	vRead        *cli.Verb // read
	vUpdate      *cli.Verb // update
	vDelete      *cli.Verb // delete
	vJoin        *cli.Verb // join
	vKeys        *cli.Verb // keys
	vHasKey      *cli.Verb // haskey
	vCount       *cli.Verb // count
	vPath        *cli.Verb // path
	vAttach      *cli.Verb // attach
	vAttachments *cli.Verb // attachments
	vDetach      *cli.Verb // detach
	vPrune       *cli.Verb // prune
	vGrid        *cli.Verb // grid
	vImport      *cli.Verb // import
	vExport      *cli.Verb // export
	vCheck       *cli.Verb // check
	vRepair      *cli.Verb // repair
	vMigrate     *cli.Verb // migrate
	vIndexer     *cli.Verb // indexer
	vDeindexer   *cli.Verb // deindexer
	vFind        *cli.Verb // find
	vCloneSample *cli.Verb // clone-sample
	vClone       *cli.Verb // clone
	vFrame       *cli.Verb // frame
	vFrames      *cli.Verb // frames
	vReframe     *cli.Verb // reframe
	vFrameLabels *cli.Verb // frame-labels
	vFrameTypes  *cli.Verb // frame-types
	vFrameDelete *cli.Verb // delete-frame
	vSyncSend    *cli.Verb // sync-send
	vSyncRecieve *cli.Verb // sync-recieve
)

func fnInit(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnInit() not implemented\n")
	return 1
}

func fnStatus(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnStatus() not implemented\n")
	return 1
}

func fnCreate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnCreate() not implemented\n")
	return 1
}

func fnRead(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnRead() not implemented\n")
	return 1
}

func fnUpdate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnUpdate() not implemented\n")
	return 1
}

func fnDelete(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnDelete() not implemented\n")
	return 1
}

func fnJoin(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnJoin() not implemented\n")
	return 1
}

func fnKeys(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnKeys() not implemented\n")
	return 1
}

func fnHasKey(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnHasKeys() not implemented\n")
	return 1
}

func fnCount(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnHasKeys() not implemented\n")
	return 1
}

func fnPath(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnPath() not implemented\n")
	return 1
}

func fnAttach(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnAttach() not implemented\n")
	return 1
}

func fnAttachments(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnAttachments() not implemented\n")
	return 1
}

func fnDetach(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnDetach() not implemented\n")
	return 1
}

func fnPrune(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnPrune() not implemented\n")
	return 1
}

func fnGrid(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnGrid() not implemented\n")
	return 1
}

func fnImport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnImport() not implemented\n")
	return 1
}

func fnExport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnExport() not implemented\n")
	return 1
}

func fnCheck(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnCheck() not implemented\n")
	return 1
}

func fnRepair(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnRepair() not implemented\n")
	return 1
}

func fnMigrate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnMigrate() not implemented\n")
	return 1
}

func fnIndexer(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnIndexer() not implemented\n")
	return 1
}

func fnDeindexer(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnDeindexer() not implemented\n")
	return 1
}

func fnFind(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFind() not implemented\n")
	return 1
}

func fnCloneSample(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnCloneSample() not implemented\n")
	return 1
}

func fnClone(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnClone() not implemented\n")
	return 1
}

func fnFrame(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFrame() not implemented\n")
	return 1
}

func fnFrames(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFrames() not implemented\n")
	return 1
}

func fnReframe(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnReframe() not implemented\n")
	return 1
}

func fnFrameLabels(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFrameLabels() not implemented\n")
	return 1
}

func fnFrameTypes(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFrameTypes() not implemented\n")
	return 1
}

func fnFrameDelete(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnFrameDelete() not implemented\n")
	return 1
}

func fnSyncSend(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnSyndSend() not implemented\n")
	return 1
}

func fnSyncRecieve(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	fmt.Fprintf(eout, "fnSyndRecieve() not implemented\n")
	return 1
}

func main() {
	//FIXME: Replace with your base package .Version attribute
	app := cli.NewCli(dataset.Version)
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
	vInit = app.NewVerb("init", "initialize a collection", fnInit)
	vStatus = app.NewVerb("status", "collection status", fnStatus)
	vCreate = app.NewVerb("create", "create a JSON object", fnCreate)
	vRead = app.NewVerb("read", "read a JSON object", fnRead)
	vUpdate = app.NewVerb("update", "update a JSON object", fnUpdate)
	vDelete = app.NewVerb("delete", "delete a JSON object", fnDelete)
	vJoin = app.NewVerb("join", "join data to a JSON object", fnJoin)
	vKeys = app.NewVerb("keys", "list keys in collection", fnKeys)
	vHasKey = app.NewVerb("haskey", "check for key in collection", fnHasKey)
	vCount = app.NewVerb("count", "count JSON objects", fnCount)
	vPath = app.NewVerb("path", "path to JSON object", fnPath)
	vAttach = app.NewVerb("attach", "attach a file to JSON object", fnAttach)
	vAttachments = app.NewVerb("attachments", "list attachments for a JSON object", fnAttachments)
	vDetach = app.NewVerb("detach", "detach a copy of the attachment from a JSON object", fnDetach)
	vPrune = app.NewVerb("prune", "prune an the attachment to a JSON object", fnPrune)
	vGrid = app.NewVerb("grid", "create a 2D JSON array from JSON objects", fnGrid)
	vImport = app.NewVerb("import", "import a table (CSV, GSheet) as JSON bject into a collection", fnImport)
	vExport = app.NewVerb("export", "export a table (CSV, GSheet) from a collection of JSON objects", fnExport)
	vCheck = app.NewVerb("check", "check a collection for errors", fnCheck)
	vRepair = app.NewVerb("repair", "repair a collection", fnRepair)
	vMigrate = app.NewVerb("migrate", "migrate a collection's layout", fnMigrate)
	vIndexer = app.NewVerb("indexer", "index a JSON object in a collection", fnIndexer)
	vDeindexer = app.NewVerb("deindex", "remove a JSON object from an index", fnDeindexer)
	vFind = app.NewVerb("find", "find a JSON object base on a dot path and value", fnFind)
	vCloneSample = app.NewVerb("clone-sample", "clone a sample from a collection", fnCloneSample)
	vClone = app.NewVerb("clone", "clone a collection", fnClone)
	vFrame = app.NewVerb("frame", "create a data frame", fnFrame)
	vFrames = app.NewVerb("frames", "list frames in a collection", fnFrames)
	vReframe = app.NewVerb("reframe", "re-create an existing frame", fnReframe)
	vFrameLabels = app.NewVerb("frame-labels", "set labels for a frame", fnFrameLabels)
	vFrameTypes = app.NewVerb("frame-types", "set the types for columns in a frame", fnFrameTypes)
	vFrameDelete = app.NewVerb("delete-frame", "delete a frame from a collection", fnFrameDelete)
	vSyncSend = app.NewVerb("sync-send", "sync from a collection to a target (e.g. CSV, GSheet)", fnSyncSend)
	vSyncRecieve = app.NewVerb("sync-recieve", "sync a collection from a source table (e.g. CSV, GSheet)", fnSyncRecieve)

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
