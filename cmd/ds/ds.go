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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
    dataset friends.ds join -overwrite frieda frieda-profile.json

    # Delete a JSON document
    dataset friends.ds delete frieda

    # Import data from a CSV file using column 1 as key
    dataset -quiet -nl=false friends.ds import my-data.csv 1

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
	showVerbose      bool

	// Application Options
	collectionName    string
	useHeaderRow      bool
	clientSecretFName string
	overwrite         bool
	batchSize         int
	keyFName          string
	collectionLayout  = "pairtree" // Default collection file layout

	// Search specific options, application Options
	showHighlight  bool
	setHighlighter string
	resultFields   string
	sortBy         string
	jsonFormat     bool
	csvFormat      bool
	csvSkipHeader  bool
	idsOnly        bool
	size           int
	from           int
	explain        bool // Note: will be force results to be in JSON format

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
	var (
		c   *dataset.Collection
		err error
	)
	fmt.Fprintf(out, "DEBUG args: %s\n", strings.Join(args, " "))
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprint(eout, "%s\n", err)
		return 1
	}
	if len(args) == 0 {
		fmt.Fprintf(eout, "Missing collection name\n")
		return 1
	}
	for _, collectionName := range args {
		switch strings.ToLower(collectionLayout) {
		case "pairtree":
			c, err = dataset.InitCollection(collectionName, dataset.PAIRTREE_LAYOUT)
		case "buckets":
			c, err = dataset.InitCollection(collectionName, dataset.BUCKETS_LAYOUT)
		default:
			fmt.Fprint(eout, "%s is an unknown layout\n", collectionLayout)
			return 1
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		c.Close()
	}
	fmt.Fprintf(out, "OK")
	return 0
}

func fnStatus(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		c   *dataset.Collection
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprint(eout, "%s\n", err)
		return 1
	}
	if len(args) == 0 {
		fmt.Fprintf(eout, "Missing collection name\n")
		return 1
	}
	for _, collectionName := range args {
		c, err = dataset.Open(collectionName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if showVerbose {
			switch c.Layout {
			case dataset.PAIRTREE_LAYOUT:
				fmt.Fprintf(out, "%s, layout pairtree, version %s\n", collectionName, c.Version)
			case dataset.BUCKETS_LAYOUT:
				fmt.Fprintf(out, "%s, layout buckets, version %s\n", collectionName, c.Version)
			default:
				fmt.Fprintf(eout, "%s, layout unknown, version %s\n", collectionName, c.Version)
				c.Close()
				return 1
			}
		}
		c.Close()
	}
	fmt.Fprintf(out, "OK")
	return 0
}

func fnCreate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		key            string
		src            []byte
		c              *dataset.Collection
		err            error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprint(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	switch len(args) {
	case 0:
		fmt.Fprintf(eout, "Missing collection name, key and JSON source\n")
		return 1
	case 1:
		fmt.Fprintf(eout, "Missing key and JSON source\n")
		return 1
	case 2:
		collectionName, key = args[0], args[1]
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing JSON source\n")
			return 1
		}
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	case 3:
		collectionName, key, src = args[0], args[1], []byte(args[2])
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	if strings.HasSuffix(key, ".json") {
		key = strings.TrimSuffix(key, ".json")
	}
	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()
	m := map[string]interface{}{}
	if err := json.Unmarshal(src, &m); err != nil {
		fmt.Fprintf(eout, "%s must be a valid JSON Object", key)
		return 1
	}
	if c.HasKey(key) == true && overwrite == true {
		if err := c.Update(key, m); err != nil {
			fmt.Fprintf(eout, "failed to create %s in %s, %s\n", key, collectionName, err)
			return 1
		}
		fmt.Fprint(out, "OK")
		return 0
	}

	if err := c.Create(key, m); err != nil {
		fmt.Fprintf(eout, "failed to create %s in %s, %s\n", key, collectionName, err)
		return 1
	}
	fmt.Fprint(out, "OK")
	return 0
}

func fnRead(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		keys           []string
		src            []byte
		c              *dataset.Collection
		err            error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprint(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		collectionName = args[0]
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		for _, line := range strings.Split(string(src), "\n") {
			s := strings.TrimSpace(line)
			if len(s) > 0 {
				keys = append(keys, s)
			}
		}
	case len(args) >= 2:
		collectionName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()
	if len(keys) == 1 {
		m := map[string]interface{}{}
		if err := c.Read(keys[0], m); err != nil {
			fmt.Fprintf(eout, "%s must be a valid JSON Object", keys[0])
			return 1
		}
		if prettyPrint {
			src, err = json.MarshalIndent(m, "", "    ")
		} else {
			src, err = json.Marshal(m)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fmt.Fprintf(out, "%s", src)
		return 0
	}

	recs := []map[string]interface{}{}
	for _, key := range keys {
		m := map[string]interface{}{}
		err := c.Read(key, m)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		recs = append(recs, m)
	}
	if prettyPrint {
		src, err = json.MarshalIndent(recs, "", "    ")
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fmt.Fprintf(out, "%s", src)
		return 0
	}
	src, err = json.Marshal(recs)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s", src)
	return 0
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
	app := cli.NewCli(dataset.Version)
	app.SetParams("COLLECTION", "VERB", "[VERB OPTIONS]", "[VERB PARAMETERS ...]")

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
	app.BoolVar(&newLine, "nl,newline", true, "if true add a trailing newline, false suppress it")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "output manpage markup")
	app.BoolVar(&showVerbose, "verbose", false, "output rows processed on importing from CSV")

	// Application Options
	app.StringVar(&collectionName, "c,collection", "", "sets the collection to be used")
	//app.BoolVar(&overwrite, "overwrite", false, "overwrite will treat a create as update if the record exists")
	//app.StringVar(&keyFName, "key-file", "", "operate on the record keys contained in file, one key per line")

	// Application Verbs
	app.VerbsRequired = true

	// Collection oriented functions
	vInit = app.NewVerb("init", "initialize a collection", fnInit)
	vInit.AddParams("COLLECTION")
	vInit.StringVar(&collectionLayout, "layout", "pairtree", "set file layout for a new collection (i.e. \"buckets\" or \"pairtree\")")
	vStatus = app.NewVerb("status", "collection status", fnStatus)
	vStatus.AddParams("COLLECTION")
	vCheck = app.NewVerb("check", "check a collection for errors", fnCheck)
	vCheck.AddParams("COLLECTION")
	vRepair = app.NewVerb("repair", "repair a collection", fnRepair)
	vRepair.AddParams("COLLECTION")
	vMigrate = app.NewVerb("migrate", "migrate a collection's layout", fnMigrate)
	vMigrate.AddParams("COLLECTION", "LAYOUT")
	vCloneSample = app.NewVerb("clone-sample", "clone a sample from a collection", fnCloneSample)
	vCloneSample.AddParams("SOURCE_COLLECTION", "SAMPLE_SIZE", "SAMPLE_COLLECTION", "[TEST_COLLECTION]")
	vClone = app.NewVerb("clone", "clone a collection", fnClone)
	vClone.AddParams("SOURCE_COLLECTION", "DESTINATION_COLLECTION")

	// Object oriented functions
	vCreate = app.NewVerb("create", "create a JSON object", fnCreate)
	vCreate.AddParams("COLLECTION", "KEY", "[JSON_SRC]")
	vCreate.StringVar(&inputFName, "i,input", "", "input file to read JSON object source from")
	vCreate.BoolVar(&overwrite, "overwrite", false, "overwrite treat a create an update if record already exists")

	vRead = app.NewVerb("read", "read a JSON object from key(s)", fnRead)
	vRead.AddParams("COLLECTION", "KEY", "[KEY ...]")
	vRead.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vUpdate = app.NewVerb("update", "update a JSON object", fnUpdate)
	vUpdate.AddParams("COLLECTION", "KEY", "[JSON_SRC]")
	vUpdate.StringVar(&inputFName, "i,input", "", "input file to read JSON object source from")

	vDelete = app.NewVerb("delete", "delete a JSON object", fnDelete)
	vDelete.AddParams("COLLECTION", "KEY", "[KEY ...]")
	vDelete.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vJoin = app.NewVerb("join", "join data to a JSON object", fnJoin)
	vJoin.AddParams("COLLECTION", "KEY", "[JSON_SRC]")
	vJoin.StringVar(&inputFName, "i,input", "", "read JSON source from file")
	vJoin.BoolVar(&overwrite, "overwrite", false, "overwrite will replace common attributes on join")

	vKeys = app.NewVerb("keys", "list keys in collection", fnKeys)
	vKeys.AddParams("COLLECTION", "[FILTER_EXPR]", "[SORT_EXPR]")

	vHasKey = app.NewVerb("haskey", "check for key in collection", fnHasKey)
	vHasKey.AddParams("COLLECTION", "KEY")

	vCount = app.NewVerb("count", "count JSON objects", fnCount)
	vCount.AddParams("COLLECTION", "[FILTER_EXPR]")

	vPath = app.NewVerb("path", "path to JSON object", fnPath)
	vPath.AddParams("COLLECTION", "[FILTER_EXPR]")

	vAttach = app.NewVerb("attach", "attach a file to JSON object", fnAttach)
	vAttach.AddParams("COLLECTION", "KEY", "FILENAMES")
	vAttachments = app.NewVerb("attachments", "list attachments for a JSON object", fnAttachments)
	vAttachments.AddParams("COLLECTION", "KEY")
	vDetach = app.NewVerb("detach", "detach a copy of the attachment from a JSON object", fnDetach)
	vDetach.AddParams("COLLECTION", "KEY", "[FILENAMES]")
	vPrune = app.NewVerb("prune", "prune an the attachment to a JSON object", fnPrune)
	vPrune.AddParams("COLLECTION", "KEY", "[FILENAMES]")

	// Import/export collections from/into tables
	vImport = app.NewVerb("import", "import a table (CSV, GSheet) as JSON bject into a collection", fnImport)
	vImport.AddParams("COLLECTION", "CSV_FILENAME|GSHEET_ID SHEET_NAME", "[RANGE]", "[KEY_COLUMN_NO]")
	vImport.StringVar(&clientSecretFName, "client-secret", "", "(import from GSheet) set the client secret path and filename for GSheet access")
	vImport.BoolVar(&useHeaderRow, "use-header-row", true, "use the header row as attribute names in the JSON object")
	vExport = app.NewVerb("export", "export a table (CSV, GSheet) from a collection of JSON objects", fnExport)
	vExport.StringVar(&clientSecretFName, "client-secret", "", "(export into a GSheet) set the client secret path and filename for GSheet access")

	// Sync send/receive collections from/to tables
	vSyncSend = app.NewVerb("sync-send", "sync a collection using a data frame sending data to a table (e.g. CSV, GSheet)", fnSyncSend)
	vSyncSend.AddParams("COLLECTION", "FRAME_NAME", "CSV_FILENAME|GSHEET_ID SHEET_NAME")
	vSyncSend.StringVar(&clientSecretFName, "client-secret", "", "(sync-send to a GSheet) set the client secret path and filename for GSheet access")
	vSyncRecieve = app.NewVerb("sync-recieve", "sync a collection using a data frame with recieve data from a table (e.g. CSV, GSheet)", fnSyncRecieve)
	vSyncRecieve.AddParams("COLLECTION", "FRAME_NAME", "CSV_FILENAME|GSHEET_ID SHEET_NAME")
	vSyncRecieve.StringVar(&clientSecretFName, "client-secret", "", "(sync-receive from a GSheet) set the client secret path and filename for GSheet access")

	// Search Verbs and options
	vFind = app.NewVerb("find", "find a JSON object base on a dot path and value", fnFind)
	vFind.AddParams("ONE_OR_MORE_INDEX_NAMES", "SEARCH_TERMS")
	vFind.StringVar(&sortBy, "sort", "", "a comma delimited list of field names to sort by")
	vFind.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	vFind.StringVar(&setHighlighter, "highlighter", "", "set the highlighter (ansi,html) for search results")
	vFind.StringVar(&resultFields, "fields", "", "comma delimited list of fields to display in the results")
	vFind.BoolVar(&jsonFormat, "json", false, "format results as a JSON document")
	vFind.BoolVar(&csvFormat, "csv", false, "format results as a CSV document, used with fields option")
	vFind.BoolVar(&csvSkipHeader, "csv-skip-header", false, "don't output a header row, only values for csv output")
	vFind.BoolVar(&idsOnly, "ids,ids-only", false, "output only a list of ids from results")
	vFind.IntVar(&from, "from", 0, "return the result starting with this result number")
	vFind.BoolVar(&explain, "explain", false, "explain results in a verbose JSON document")
	vFind.IntVar(&batchSize, "batch,size", 100, "set the number of records per response")
	vIndexer = app.NewVerb("indexer", "index a JSON object in a collection", fnIndexer)
	vIndexer.AddParams("COLLECTION", "INDEX_NAME", "INDEX_MAP_FILE")
	vIndexer.IntVar(&batchSize, "batch,size", 100, "set the number of records per response")
	vDeindexer = app.NewVerb("deindex", "remove a JSON object from an index", fnDeindexer)
	vDeindexer.AddParams("INDEX_NAME", "KEY")
	//vDeindexer.IntVar(&batchSize, "batch,size", 100, "set the number of records per response")

	// Frames and Grid
	vGrid = app.NewVerb("grid", "create a 2D JSON array from JSON objects", fnGrid)
	vGrid.AddParams("COLLECTION", "DOTPATH", "[DOTPATH ...]")
	vFrame = app.NewVerb("frame", "create a data frame", fnFrame)
	vFrame.AddParams("COLLECTION", "FRAME_NAME", "DOTPATH", "[DOTPATH ...]")
	vFrame.StringVar(&inputFName, "i,input", "", "frame only the keys listed in the file, one key per line")
	vFrames = app.NewVerb("frames", "list frames in a collection", fnFrames)
	vFrames.AddParams("COLLECTION")
	vReframe = app.NewVerb("reframe", "re-create an existing frame", fnReframe)
	vReframe.AddParams("COLLECTION", "FRAME_NAME")
	vReframe.StringVar(&inputFName, "i,input", "", "frame only the keys listed in the file, one key per line")
	vFrameLabels = app.NewVerb("frame-labels", "set labels for all columns in a frame", fnFrameLabels)
	vFrameLabels.AddParams("COLLECTION", "FRAME_NAME", "LABEL", "[LABEL ...]")
	vFrameTypes = app.NewVerb("frame-types", "set the types for all columns in a frame", fnFrameTypes)
	vFrameTypes.AddParams("COLLECTION", "FRAME_NAME", "TYPE", "[TYPE ...]")
	vFrameDelete = app.NewVerb("delete-frame", "delete a frame from a collection", fnFrameDelete)
	vFrameDelete.AddParams("COLLECTION", "FRAME_NAME")

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
