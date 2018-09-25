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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheets"
	"github.com/caltechlibrary/shuffle"
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
    dataset read friends.ds frieda

    # Path to JSON document
    dataset path friends.ds frieda

    # Update a JSON document
    dataset update friends.ds frieda '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys friends.ds

    # Get keys filtered for the name "frieda"
    dataset keys friends.ds '(eq .name "frieda")'

    # Join frieda-profile.json with "frieda" adding unique key/value pairs
    dataset join friends.ds frieda frieda-profile.json

    # Join frieda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from frieda-profile.json
    dataset join -overwrite friends.ds frieda frieda-profile.json

    # Delete a JSON document
    dataset delete friends.ds frieda

    # Import data from a CSV file using column 1 as key
    dataset import friends.ds my-data.csv 1

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
	collectionName string
	// NOTE: setAllKeys is used by Reframe to persist an useAllKeys value
	setAllKeys bool
	// One time use of all keys in grid, a newly created Frame or Reframe
	useAllKeys        bool
	useHeaderRow      bool
	clientSecretFName string
	overwrite         bool
	syncOverwrite     bool
	batchSize         int
	sampleSize        int
	keyFName          string
	collectionLayout  = "pairtree" // Default collection file layout
	filterExpr        string
	sortExpr          string

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
	vHasFrame    *cli.Verb // has-frame
	vFrames      *cli.Verb // frames
	vReframe     *cli.Verb // reframe
	vFrameLabels *cli.Verb // frame-labels
	vFrameTypes  *cli.Verb // frame-types
	vFrameDelete *cli.Verb // delete-frame
	vSyncSend    *cli.Verb // sync-send
	vSyncRecieve *cli.Verb // sync-recieve
)

// keysFromSrc takes a byte splice, splits them on "\n" and converts any
// non-empty line string appended to the keys slice
func keysFromSrc(src []byte) []string {
	var keys []string
	for _, line := range strings.Split(string(src), "\n") {
		s := strings.TrimSpace(line)
		if len(s) > 0 {
			keys = append(keys, s)
		}
	}
	return keys
}

// fnInit - create a dataset collection
func fnInit(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		c   *dataset.Collection
		err error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

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
			fmt.Fprintf(eout, "%s is an unknown layout\n", collectionLayout)
			return 1
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		c.Close()
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnStatus - given a path see if it is a collection by attempting to "open" it
func fnStatus(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		c   *dataset.Collection
		err error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

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
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnCreate - add a new JSON document in  collection
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
		fmt.Fprintf(eout, "%s\n", err)
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
		collectionName, key = args[0], args[1]
		// Need to decide if args[2] is JSON source or filename
		if strings.HasPrefix(args[2], "{") && strings.HasSuffix(args[2], "}") {
			src = []byte(args[2])
		} else {
			src, err = ioutil.ReadFile(args[2])
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
		}

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
			fmt.Fprintf(eout, "failed to create %q in %s, %s\n", key, collectionName, err)
			return 1
		}
		if quiet == false {
			fmt.Fprintf(out, "OK")
		}
		return 0
	}

	if err := c.Create(key, m); err != nil {
		fmt.Fprintf(eout, "failed to create %q in %s, %s\n", key, collectionName, err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnRead - retreive a JSON document from a collection
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
		fmt.Fprintf(eout, "%s\n", err)
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
		keys = keysFromSrc(src)
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
			fmt.Fprintf(eout, "%s, %s\n", keys[0], err)
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

// fnUpdate - replace a JSON document in a collection
func fnUpdate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		key            string
		src            []byte
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
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
	if err := c.Update(key, m); err != nil {
		fmt.Fprintf(eout, "failed to update %s in %s, %s\n", key, collectionName, err)
		return 1
	}
	if quiet == false {
		fmt.Fprint(out, "OK")
	}
	return 0
}

// fnDelete - remove a JSON document from a collection
func fnDelete(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		keys           []string
		src            []byte
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = keysFromSrc(src)
	}
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		collectionName = args[0]
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

	for _, key := range keys {
		err := c.Delete(key)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnJoin - joins a JSON object in the collection with a new JSON object appending
// new attributes and optionally overwriting existing attribute in common.
func fnJoin(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		key            string
		src            []byte
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
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
		if strings.HasPrefix(args[2], "{") && strings.HasSuffix(args[2], "}") {
			collectionName, key, src = args[0], args[1], []byte(args[2])
		} else {
			collectionName, key = args[0], args[1]
			src, err = ioutil.ReadFile(args[2])
			if err != nil {
				fmt.Fprintf(eout, "%s", err)
				return 1
			}
		}
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
	// unmarshal new object
	newObj := map[string]interface{}{}
	if err := json.Unmarshal(src, &newObj); err != nil {
		fmt.Fprintf(eout, "%s must be a valid JSON Object", key)
		return 1
	}
	// Join the object
	err = c.Join(key, newObj, overwrite)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if quiet == false {
		fmt.Fprint(out, "OK")
	}
	return 0
}

// fnKeys returns the keys in a collection
// If a 'filter expression' is provided it will return a filtered list of keys.
// Filters with like Go's text/template if statement where the 'filter expression' is
// the condititional expression in a if/else statement. If the expression evaluates to "true"
// then the key is included in the list of keys If the expression evaluates to "false" then
// it is excluded for the list of keys.
// If a 'sort expression' is provided then the resulting keys are ordered by that expression.
func fnKeys(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		keys           []string
		c              *dataset.Collection
		err            error
		src            []byte
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		collectionName = args[0]
	case len(args) == 2:
		collectionName, filterExpr = args[0], args[1]
	case len(args) == 3:
		collectionName, filterExpr, sortExpr = args[0], args[1], args[2]
	case len(args) > 3:
		collectionName, filterExpr, sortExpr, keys = args[0], args[1], args[2], args[3:]
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

	// If we don't have a sub selection of keys, get a complete list of keys
	if len(keys) == 0 {
		keys = c.Keys()
	}
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = keysFromSrc(src)
	}

	// Apply Filter Expression
	if len(filterExpr) > 0 && filterExpr != "true" {
		keys, err = c.KeyFilter(keys[:], filterExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
	}

	// If now sort we're done
	if len(sortExpr) == 0 {
		fmt.Fprintf(out, "%s", strings.Join(keys, "\n"))
		return 0
	}

	// We still have sorting to do.
	keys, err = c.KeySortByExpression(keys, sortExpr)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	fmt.Fprintf(out, strings.Join(keys, "\n"))
	return 0
}

// fnHasKey - check if a key to an object exists in a collection optionally matching keys and a filter expression
func fnHasKey(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		c              *dataset.Collection
		collectionName string
		keys           []string
		err            error
		src            []byte
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	// Process positional parameters
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		collectionName = args[0]
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
	case len(args) >= 2:
		collectionName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read in any key list from a file.
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = append(keys, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	for i, key := range keys {
		if i > 0 {
			fmt.Fprintf(out, "\n")
		}
		if c.HasKey(key) {
			fmt.Fprintf(out, "true")
		} else {
			fmt.Fprintf(out, "false")
		}
	}
	return 0
}

// fnCount - count objects in a collection, optionally matching keys and a filter expression
func fnCount(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		keys           []string
		c              *dataset.Collection
		err            error
		src            []byte
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		collectionName = args[0]
	case len(args) == 2:
		collectionName, filterExpr = args[0], args[1]
	case len(args) > 2:
		collectionName, filterExpr, keys = args[0], args[1], args[2:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read keys from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = append(keys, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// If we don't have a sub selection of keys, get a list of keys
	if len(keys) == 0 {
		keys = c.Keys()
	}

	// Process the filter against the keys if necessary
	if len(filterExpr) > 0 && filterExpr != "true" {
		keys, err = c.KeyFilter(keys[:], filterExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}
	fmt.Fprintf(out, "%d", len(keys))
	return 0
}

// fnPath - return a path(s) to an object(s) given a key(s)
func fnPath(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		keys           []string
		src            []byte
		docPath        string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		collectionName = args[0]
	case len(args) >= 2:
		collectionName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read keys from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = append(keys, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	errCnt := 0
	for i, key := range keys {
		if i > 0 {
			fmt.Fprintf(out, "\n")
		}
		docPath, err = c.DocPath(key)
		if err != nil {
			fmt.Fprintf(eout, "key %q, %s\n", key, err)
			errCnt++
		} else {
			fmt.Fprintf(out, "%s", docPath)
		}
	}
	return errCnt
}

// fnAttach - attach a file(s) to an object
func fnAttach(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		key            string
		src            []byte
		fNames         []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key and attachment name(s)\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key and attachment name(s)\n")
		return 1
	case len(args) == 2:
		if len(fNames) == 0 {
			fmt.Fprintf(eout, "Missing attachment name(s)\n")
			return 1
		}
		collectionName, key = args[0], args[1]
	case len(args) >= 3:
		collectionName, key, fNames = args[0], args[1], args[2:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read filenames from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fNames = append(fNames, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.HasKey(key) == false {
		fmt.Fprintf(eout, "%q is not in %s\n", key, collectionName)
		return 1
	}
	for _, fname := range fNames {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			fmt.Fprintf(eout, "%s does not exist\n", fname)
			return 1
		}
	}
	err = c.AttachFiles(key, fNames...)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	return 0
}

// fnAttachments - list the attachments of an object(s) given a key(s)
func fnAttachments(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		keys           []string
		src            []byte
		attachments    []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		collectionName = args[0]
	case len(args) >= 2:
		collectionName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read keys from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = append(keys, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	errCnt := 0
	for i, key := range keys {
		if i > 0 {
			fmt.Fprintf(out, "\n")
		}
		attachments, err = c.Attachments(key)
		if err != nil {
			fmt.Fprintf(eout, "key %q, %s\n", key, err)
			errCnt++
		} else {
			fmt.Fprintf(out, "%s", strings.Join(attachments, "\n"))
		}
	}
	return errCnt
}

// fnDetach - return a file(s) attached to an object(s) for a given key
func fnDetach(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		key            string
		src            []byte
		fNames         []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and key\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key\n")
		return 1
	case len(args) == 2:
		collectionName, key = args[0], args[1]
	case len(args) >= 3:
		collectionName, key, fNames = args[0], args[1], args[2:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read filenames from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fNames = append(fNames, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.HasKey(key) == false {
		fmt.Fprintf(eout, "%q is not in %s", key, collectionName)
		return 1
	}
	err = c.GetAttachedFiles(key, fNames...)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprint(out, "OK")
	}
	return 0
}

// fnPrune - remove a file(s) attached to an object for a given key
func fnPrune(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		key            string
		src            []byte
		fNames         []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and key\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key\n")
		return 1
	case len(args) == 2:
		collectionName, key = args[0], args[1]
	case len(args) >= 3:
		collectionName, key, fNames = args[0], args[1], args[2:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	// Read filenames from inputFName
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fNames = append(fNames, keysFromSrc(src)...)
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.HasKey(key) == false {
		fmt.Fprintf(eout, "%q is not in %s", key, collectionName)
		return 1
	}
	err = c.Prune(key, fNames...)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprint(out, "OK")
	}
	return 0
}

// fnGrid - generate a grid (2D array) based on a list of key(s) and dotpath(s).
// Keys map to rows, dotpaths map to columns
//
// Command Syntax: [VERB_OPTIONS] COLLECTION_NAME DOTPATH [DOTPATH ...]
// Verb Options: filter-expression (-filter) , key list filename (-i,-input), sample size (-sample), verbose (-v, verbose)
func fnGrid(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		keys           []string
		dotPaths       []string
		src            []byte
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key list filename, and dot path(s)\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing dot paths\n")
		return 1
	case len(args) >= 2:
		collectionName, dotPaths = args[0], args[1:]
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

	// Get all keys or read from inputFName
	if useAllKeys && len(inputFName) == 0 {
		keys = c.Keys()
	}
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = append(keys, keysFromSrc(src)...)
	}

	// Apply Filter Expression
	if len(filterExpr) > 0 && filterExpr != "true" {
		keys, err = c.KeyFilter(keys[:], filterExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
	}

	g, err := c.Grid(keys, dotPaths, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if prettyPrint {
		src, err = json.MarshalIndent(g, "", "    ")
	} else {
		src, err = json.Marshal(g)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s", src)
	return 0
}

// fnFrame - define a data frame and populate it with a list of keys and doptpaths
// syntax: [VERB_OPTIONS] COLLECTION_NAME FRAME_NAME DOTPATH [DOTPATH ...]
// Verb Options: filter-expression (e.g. -filter) , key list filename (e.g. -i), sample size (e.g. -sample)
// labels (e.g. -labels)
func fnFrame(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		f              *dataset.DataFrame
		c              *dataset.Collection
		keys           []string
		dotPaths       []string
		//labels         []string
		src []byte
		err error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and frame name\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name\n")
		return 1
	case len(args) == 2:
		collectionName, frameName = args[0], args[1]
	case len(args) >= 3:
		collectionName, frameName, dotPaths = args[0], args[1], args[2:]
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

	// Check to see if frame exists...
	if c.HasFrame(frameName) {
		if len(dotPaths) > 0 || len(filterExpr) > 0 {
			fmt.Fprintf(eout, "frame %q already exists\n", frameName)
			return 1
		}
		f, err = c.Frame(frameName, nil, nil, showVerbose)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		// Handle pretty printing
		if prettyPrint {
			src, err = json.MarshalIndent(f, "", "    ")
		} else {
			src, err = json.Marshal(f)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		fmt.Fprintf(out, "%s", src)
		return 0
	}

	// Get all keys or read from inputFName
	if useAllKeys {
		keys = c.Keys()
	}
	if len(inputFName) > 0 {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = keysFromSrc(src)
	}

	// Apply Filter Expression
	if len(filterExpr) > 0 {
		keys, err = c.KeyFilter(keys[:], filterExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
	}

	// Apply SortExpr to key order
	if sortExpr != "" {
		keys, err = c.KeySortByExpression(keys, sortExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}

	// Run a sanity check before we create a new frame...
	if len(dotPaths) == 0 {
		fmt.Fprintf(eout, "No dotpaths, frame creation aborted\n")
		return 1
	}
	if len(keys) == 0 {
		fmt.Fprintf(eout, "No keys, frame creation aborted\n")
		return 1
	}

	// NOTE: See if we are reading a frame back or define one.
	f, err = c.Frame(frameName, keys, dotPaths, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// NOTE: Make need to make sure we save our additional
	// settings - useAllKeys, sampleSize, sortExpr and filterExpr
	f.AllKeys = useAllKeys
	f.FilterExpr = filterExpr
	f.SortExpr = sortExpr
	if sampleSize >= 0 {
		f.SampleSize = sampleSize
	}

	// Handle pretty printing
	if prettyPrint {
		src, err = json.MarshalIndent(f, "", "    ")
	} else {
		src, err = json.Marshal(f)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s", src)
	return 0
}

// fnFrames - list the frames defined in a collection.
func fnFrames(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		frameNames     []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and frame name\n")
		return 1
	case len(args) == 1:
		collectionName = args[0]
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

	frameNames = c.Frames()
	if len(frameNames) > 0 {
		fmt.Fprintf(out, "%s", strings.Join(frameNames, "\n"))
	}
	return 0
}

// fnFrameLabels - set the labels (column headings) associated with a frame's grid.
func fnFrameLabels(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		frameName      string
		labels         []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, frame name and labels\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name and labels\n")
		return 1
	case len(args) == 2:
		fmt.Fprintf(eout, "labels\n")
		return 1
	case len(args) >= 3:
		collectionName, frameName, labels = args[0], args[1], args[2:]
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

	err = c.FrameLabels(frameName, labels)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnFrameTypes - set the column types (for column values) associated with a frame's grid.
func fnFrameTypes(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		c              *dataset.Collection
		frameName      string
		types          []string
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, frame name and types\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name and types\n")
		return 1
	case len(args) == 2:
		fmt.Fprintf(eout, "types\n")
		return 1
	case len(args) >= 3:
		collectionName, frameName, types = args[0], args[1], args[2:]
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

	err = c.FrameTypes(frameName, types)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnHasFrame - check if a frame has been defined in collection
func fnHasFrame(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and frame name\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name\n")
		return 1
	case len(args) == 2:
		collectionName, frameName = args[0], args[1]
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

	if c.HasFrame(frameName) {
		fmt.Fprintf(out, "true")
	} else {
		fmt.Fprintf(out, "false")
	}
	return 0
}

// fnFrameDelete - remove a frame from a collection
func fnFrameDelete(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and frame name\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name\n")
		return 1
	case len(args) == 2:
		collectionName, frameName = args[0], args[1]
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

	err = c.DeleteFrame(frameName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnReframe(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		f              *dataset.DataFrame
		c              *dataset.Collection
		keys           []string
		dotPaths       []string
		frameUpdated   bool
		src            []byte
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and frame name\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name\n")
		return 1
	case len(args) == 2:
		collectionName, frameName = args[0], args[1]
	case len(args) >= 3:
		collectionName, frameName, dotPaths = args[0], args[1], args[2:]
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

	// Check to see if frame exists...
	if c.HasFrame(frameName) == false {
		fmt.Fprintf(eout, "Frame %q not defined in %s\n", frameName, collectionName)
		return 1
	}

	f, err = c.Frame(frameName, nil, nil, false)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// Update frame settings
	frameUpdated = false

	if setAllKeys == true {
		f.AllKeys = useAllKeys
		frameUpdated = true
	}

	// Now we can rebuild frame...
	if f.AllKeys || useAllKeys {
		keys = c.Keys()
		frameUpdated = true
	} else {
		keys = f.Keys[:]
	}

	// Read from inputFName, update frame's keys
	if len(inputFName) > 0 {
		f.AllKeys = false
		frameUpdated = true
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = keysFromSrc(src)
	}

	// Apply Filter Expression
	if len(filterExpr) > 0 {
		keys, err = c.KeyFilter(keys[:], filterExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		f.FilterExpr = filterExpr
		frameUpdated = true
	}

	// Apply Sample Size
	if sampleSize > -1 {
		f.SampleSize = sampleSize
	}
	if f.SampleSize > 0 {
		frameUpdated = true
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
		f.SampleSize = sampleSize
		frameUpdated = true
	}

	// Apply SortExpr
	if sortExpr != "" {
		f.SortExpr = sortExpr
		keys, err = c.KeySortByExpression(keys, sortExpr)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		frameUpdated = true
	}

	// NOTE: See if we are reading a frame back or define one.
	if len(dotPaths) > 0 {
		f.DotPaths = dotPaths
		frameUpdated = true
	}

	// Finally run some sanity checks before updating frame...
	if len(f.DotPaths) == 0 {
		fmt.Fprintf(eout, "Missing dotpaths in frame definition\n")
		return 1
	}
	if len(keys) == 0 {
		fmt.Fprintf(eout, "No keys available to update frame\n")
		return 1
	}

	// Save the updated frame definition
	if frameUpdated {
		err = c.SaveFrame(frameName, f)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	}

	// Now regenerate grid content with Reframe
	err = c.Reframe(frameName, keys, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnImport - import a CSV file or GSheet into a collection
// syntax: COLLECTION CSV_FILENAME ID_COL CELL_RANGE
//         COLLECTION GSHEET_ID SHEET_NAME ID_COL [CELL_RANGE]
// options:
// -overwrite
// -use-header-row
// -verbose
// -client-secret
func fnImport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		csvFName       string
		gSheetID       string
		gSheetName     string
		idColNoString  string
		idCol          int
		cellRange      string
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, filename (gSheet ID and Sheet name), ID col no\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing filename and table details\n")
		return 1
	case len(args) < 3:
		fmt.Fprintf(eout, "Missing table details (e.g. ID_COL_NO) \n")
		return 1
	case len(args) == 3:
		collectionName, csvFName, idColNoString = args[0], args[1], args[2]
	case len(args) == 4:
		cellRange = "A1:Z"
		collectionName, gSheetID, gSheetName, idColNoString = args[0], args[1], args[2], args[3]
	case len(args) == 5:
		collectionName, gSheetID, gSheetName, idColNoString, cellRange = args[0], args[1], args[2], args[3], args[4]
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

	idCol, err = strconv.Atoi(idColNoString)
	if err != nil {
		fmt.Fprintf(eout, "expected column id number, %s\n", err)
		return 1
	}
	// NOTE: We need to convert column number to zero based columns
	idCol--
	if idCol < 0 {
		fmt.Fprintf(eout, "column number must be greater than zero")
		return 1
	}

	// See if we have a GSheet ID or CSV filename
	if len(csvFName) > 0 {
		fp, err := os.Open(csvFName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		cnt, err := c.ImportCSV(fp, idCol, useHeaderRow, overwrite, showVerbose)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		} else if showVerbose {
			fmt.Fprintf(out, "%d total rows processed\n", cnt)
		}
	} else {
		//FIXME: Need better search process for finding the google access key
		clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
		if clientSecretFName != "" {
			clientSecretJSON = clientSecretFName
		}
		if clientSecretJSON == "" {
			//clientSecretJSON = "client_secret.json"
			clientSecretJSON = "credentials.json"
		}
		table, err := gsheets.ReadSheet(clientSecretJSON, gSheetID, gSheetName, cellRange)
		if err != nil {
			fmt.Fprintf(eout, "Errors importing %s, %s", gSheetName, err)
			return 1
		}
		if cnt, err := c.ImportTable(table, idCol, useHeaderRow, overwrite, showVerbose); err != nil {
			fmt.Fprintf(eout, "Errors importing %s, %s", gSheetName, err)
			return 1
		} else if showVerbose == true {
			fmt.Fprintf(out, "%d total rows processed\n", cnt)
		}
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnExport - export collection objects to a CSV file or GSheet
// syntax examples: COLLECTION FRAME [CSV_FILENAME]
//                  COLLECTION FRAME CSV_FILENAME
//                  COLLECTION FRAME GSHEET_ID GSHEET_NAME [CELL_RANGE]
// options:
// -overwrite
// -use-header-row
// -verbose
// -client-secret
func fnExport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		gSheetID       string
		gSheetName     string
		cellRange      string
		c              *dataset.Collection
		f              *dataset.DataFrame
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, frame name, filename (or gSheet ID and Sheet name)\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing frame name and filename (or gSheet ID and Sheet name)\n")
		return 1
	case len(args) == 2:
		collectionName, frameName = args[0], args[1]
	case len(args) == 3:
		collectionName, frameName, outputFName = args[0], args[1], args[2]
		if outputFName != "-" {
			fp, err := os.Create(outputFName)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			defer fp.Close()
			out = fp
		}
	case len(args) == 4:
		collectionName, frameName, gSheetID, gSheetName = args[0], args[1], args[2], args[3]
	case len(args) == 5:
		collectionName, frameName, gSheetID, gSheetName, cellRange = args[0], args[1], args[2], args[3], args[4]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	if outputFName == "" && gSheetID == "" {
		fmt.Fprintf(eout, "Missing output name or gSheet ID with Sheet Name\n")
		return 1
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// for GSheet: COLLECTION FRAME_NAME SHEET_ID SHEET_NAME
	// for CSV: COLLECTION FRAME_NAME FILENAME

	// Get Frame
	if c.HasFrame(frameName) == false {
		fmt.Fprintf(eout, "Missing frame %q in %s\n", frameName, collectionName)
		return 1
	}
	// Get dotpaths and column labels from frame
	f, err = c.Frame(frameName, nil, nil, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if f.FilterExpr == "" {
		f.FilterExpr = "true"
	}

	cnt := 0
	table := [][]interface{}{}
	if len(gSheetID) == 0 {
		cnt, err = c.ExportCSV(out, eout, f, showVerbose)
	} else {
		//FIXME: Need a better way to indentify the clientSecretName...
		clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
		if clientSecretFName != "" {
			clientSecretJSON = clientSecretFName
		}
		if clientSecretJSON == "" {
			//clientSecretJSON = "client_secret.json"
			clientSecretJSON = "credentials.json"
		}
		// gSheet expects a cell range, so we will generate one if needed.
		if cellRange == "" {
			lastCol := gsheets.ColNoToColLetters(len(f.Labels))
			lastRow := len(f.Keys) + 2
			cellRange = fmt.Sprintf("A1:%s%d", lastCol, lastRow)
		}

		//NOTE: we export to GSheet via creating a table [][]interface{}{}
		cnt, table, err = c.ExportTable(eout, f, showVerbose)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		err = gsheets.WriteSheet(clientSecretJSON, gSheetID, gSheetName, cellRange, table)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if showVerbose {
		fmt.Fprintf(out, "%d total objects processed\n", cnt)
	}
	if outputFName != "" && outputFName != "-" {
		if quiet == false {
			fmt.Fprintf(out, "OK")
		}
	}
	return 0
}

// fnSyncSend - synchronize a frame sending data to a CSV file or GSheet
// syntax: COLLECTION FRAME [CSV_FILENAME|GSHEET SHEET ID]
//
func fnSyncSend(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		csvFilename    string
		gSheetID       string
		gSheetName     string
		cellRange      string
		c              *dataset.Collection
		src            []byte
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch len(args) {
	case 0:
		fmt.Fprintf(eout, "Missing collection name, frame name and csv filename or gsheet id and gsheet name\n")
		return 1
	case 1:
		fmt.Fprintf(eout, "Missing frame name and csv filename or gsheet id with sheet name\n")
		return 1
	case 2:
		collectionName, frameName = args[0], args[1]
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing csv filename or gsheet id with sheet name\n")
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
		collectionName, frameName, csvFilename = args[0], args[1], args[2]
		src, err = ioutil.ReadFile(csvFilename)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if len(src) == 0 {
			fmt.Fprintf(eout, "No data in csv file %s\n", csvFilename)
			return 1
		}
	case 4:
		collectionName, frameName, gSheetID, gSheetName = args[0], args[1], args[2], args[3]
		cellRange = "A1:Z"
	case 5:
		collectionName, frameName, gSheetID, gSheetName, cellRange = args[0], args[1], args[2], args[3], cellRange
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	table := [][]string{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		table, err = r.ReadAll()
	} else {
		// for GSheet
		clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
		if clientSecretFName != "" {
			clientSecretJSON = clientSecretFName
		}
		if clientSecretJSON == "" {
			//clientSecretJSON = "client_secret.json"
			clientSecretJSON = "credentials.json"
		}
		table, err = gsheets.ReadSheet(clientSecretJSON, gSheetID, gSheetName, cellRange)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// Merge collection content into table
	table, err = c.MergeIntoTable(frameName, table, syncOverwrite, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// Save the resulting table
	if len(src) > 0 {
		w := csv.NewWriter(out)
		w.WriteAll(table)
		err = w.Error()
	} else {
		clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
		if clientSecretFName != "" {
			clientSecretJSON = clientSecretFName
		}
		if clientSecretJSON == "" {
			//clientSecretJSON = "client_secret.json"
			clientSecretJSON = "credentials.json"
		}
		// NOTE: WriteSheet expects a [][]interface{} not [][]string,
		// need to convert. This is a hack...
		t := [][]interface{}{}
		for _, row := range table {
			cells := []interface{}{}
			for _, cell := range row {
				cells = append(cells, cell)
			}
			t = append(t, cells)
		}
		err = gsheets.WriteSheet(clientSecretJSON, gSheetID, gSheetName, cellRange, t)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnSyncRecieve - synchronize a frame receiving data from a CSV file or GSheet
func fnSyncRecieve(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		frameName      string
		csvFilename    string
		gSheetID       string
		gSheetName     string
		cellRange      string
		c              *dataset.Collection
		src            []byte
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch len(args) {
	case 0:
		fmt.Fprintf(eout, "Missing collection name, frame name and csv filename or gsheet id and gsheet name\n")
		return 1
	case 1:
		fmt.Fprintf(eout, "Missing frame name and csv filename or gsheet id with sheet name\n")
		return 1
	case 2:
		collectionName, frameName = args[0], args[1]
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing csv filename or gsheet id with sheet name\n")
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
		collectionName, frameName, csvFilename = args[0], args[1], args[2]
		src, err = ioutil.ReadFile(csvFilename)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	case 4:
		collectionName, frameName, gSheetID, gSheetName = args[0], args[1], args[2], args[3]
		cellRange = "A1:Z"
	case 5:
		collectionName, frameName, gSheetID, gSheetName, cellRange = args[0], args[1], args[2], args[3], cellRange
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	table := [][]string{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		table, err = r.ReadAll()
	} else {
		// for GSheet
		clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
		if clientSecretFName != "" {
			clientSecretJSON = clientSecretFName
		}
		if clientSecretJSON == "" {
			//clientSecretJSON = "client_secret.json"
			clientSecretJSON = "credentials.json"
		}
		table, err = gsheets.ReadSheet(clientSecretJSON, gSheetID, gSheetName, cellRange)
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	c, err = dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// Merge table contents into Collection and Frame
	err = c.MergeFromTable(frameName, table, syncOverwrite, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnIndexer(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		defName        string
		idxName        string
		src            []byte
		keys           []string
		c              *dataset.Collection
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, index definition filename and index name\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing index definition name and index name\n")
		return 1
	case len(args) == 2:
		collectionName, defName = args[0], args[1]
		ext := path.Ext(defName)
		idxName = strings.TrimSuffix(defName, ext) + ".bleve"
	case len(args) == 3:
		collectionName, defName, idxName = args[0], args[1], args[2]
	case len(args) == 4:
		collectionName, defName, idxName, keys = args[0], args[1], args[2], args[3:]
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

	if inputFName != "" {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = strings.Split(string(src), "\n")
	}
	if len(keys) == 0 {
		keys = c.Keys()
	}

	if batchSize == 0 {
		if len(keys) > 100000 {
			batchSize = 1000
		} else if len(keys) > 10000 {
			batchSize = len(keys) / 100
		} else if len(keys) > 1000 {
			batchSize = len(keys) / 10
		} else {
			batchSize = 100
		}
	}

	err = c.Indexer(idxName, defName, keys, batchSize)
	if err != nil {
		fmt.Fprintf(eout, "Indexing error %s %s, %s\n", collectionName, idxName, err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnFind(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		idxNames    []string
		queryString string
		err         error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing index name(s) and query string\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "query string\n")
		return 1
	case len(args) >= 2:
		i := len(args) - 1
		queryString = args[i]
		idxNames = args[0:i]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	options := map[string]string{}
	if explain == true {
		options["explain"] = "true"
		jsonFormat = true
	}

	if sampleSize > 0 {
		options["sample"] = fmt.Sprintf("%d", sampleSize)
	}
	if from != 0 {
		options["from"] = fmt.Sprintf("%d", from)
	}
	if batchSize > 0 {
		options["size"] = fmt.Sprintf("%d", batchSize)
	}
	if sortBy != "" {
		options["sort"] = sortBy
	}
	if showHighlight == true {
		options["highlight"] = "true"
		options["highlighter"] = "ansi"
	}
	if setHighlighter != "" {
		options["highlight"] = "true"
		options["highlighter"] = setHighlighter
	}

	if resultFields != "" {
		options["fields"] = strings.TrimSpace(resultFields)
	} else {
		options["fields"] = "*"
	}

	idxList, idxFields, err := dataset.OpenIndexes(idxNames)
	if err != nil {
		fmt.Fprintf(eout, "Can't open index %s, %s", strings.Join(idxNames, ", "), err)
		return 1
	}

	results, err := dataset.Find(idxList.Alias, queryString, options)
	if err != nil {
		fmt.Fprintf(eout, "Find error %s, %s", strings.Join(idxNames, ", "), err)
		return 1
	}
	err = idxList.Close()
	if err != nil {
		fmt.Fprintf(eout, "Can't close indexes %s, %s", strings.Join(idxNames, ", "), err)
	}

	//
	// Handle results formatting choices
	//
	switch {
	case jsonFormat == true:
		if prettyPrint {
			src, err := json.MarshalIndent(results, "", "    ")
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s", src)
		} else {
			src, err := json.Marshal(results)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s", src)
		}
	case csvFormat == true:
		var fields []string
		if resultFields == "" {
			fields = idxFields
		} else {
			fields = strings.Split(resultFields, ",")
		}
		err = dataset.CSVFormatter(out, results, fields, csvSkipHeader)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	case idsOnly == true:
		//NOTE: idsOnly should never be quiet, per issue 46.
		ids := []string{}
		for _, hit := range results.Hits {
			ids = append(ids, hit.ID)
		}
		fmt.Fprintf(out, strings.Join(ids, "\n"))
	}
	return 0
}

func fnDeindexer(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		idxName string
		src     []byte
		keys    []string
		err     error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing index name and key(s)\n")
		return 1
	case len(args) == 1:
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		idxName = args[0]
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = strings.Split(string(src), "\n")
	case len(args) >= 2:
		idxName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	if len(keys) == 0 {
		fmt.Fprintf(eout, "Missing any keys to de-index in %s\n", idxName)
		return 1
	}

	if batchSize == 0 {
		if len(keys) > 100000 {
			batchSize = 1000
		} else if len(keys) > 10000 {
			batchSize = len(keys) / 100
		} else if len(keys) > 1000 {
			batchSize = len(keys) / 10
		} else {
			batchSize = 100
		}
	}

	err = dataset.Deindexer(idxName, keys, batchSize)
	if err != nil {
		fmt.Fprintf(eout, "Deindexing error %s, %s\n", idxName, err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnCheck(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	if len(args) == 0 {
		fmt.Fprintf(eout, "Missing collection name(s) to check\n")
		return 1
	}
	for _, collectionName := range args {
		err = dataset.Analyzer(collectionName)
		if err != nil {
			fmt.Fprintf(eout, "error in %q, %s\n", collectionName, err)
			return 1
		}
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnRepair(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	if len(args) == 0 {
		fmt.Fprintf(eout, "Missing collection name(s) to check\n")
		return 1
	}
	for _, collectionName := range args {
		err = dataset.Repair(collectionName)
		if err != nil {
			fmt.Fprintf(eout, "error in %q, %s\n", collectionName, err)
			return 1
		}
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnMigrate(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		collectionName string
		layoutName     string
		layout         int
		err            error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and layout (e.g. pairtree, buckets\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing layout (e.g. pairtree, buckets\n")
		return 1
	case len(args) == 2:
		collectionName, layoutName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	switch strings.ToLower(layoutName) {
	case "buckets":
		layout = dataset.BUCKETS_LAYOUT
	case "pairtree":
		layout = dataset.PAIRTREE_LAYOUT
	default:
		fmt.Fprintf(eout, "Unsported layout %q", layoutName)
		return 1
	}
	// Repair collection as prep for migration
	if err = dataset.Repair(collectionName); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if err = dataset.Migrate(collectionName, layout); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnClone(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		srcCollectionName  string
		destCollectionName string
		keys               []string
		src                []byte
		err                error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch len(args) {
	case 0:
		fmt.Fprintf(eout, "Missing source and destination collections name\n")
		return 1
	case 1:
		fmt.Fprintf(eout, "Missing destination collection name\n")
		return 1
	case 2:
		srcCollectionName, destCollectionName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Too many parameters %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(srcCollectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if inputFName != "" {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = strings.Split(string(src), "\n")
	} else {
		keys = c.Keys()
	}
	if len(keys) == 0 {
		fmt.Fprintf(eout, "No objects to clone\n")
		return 1
	}

	err = c.Clone(destCollectionName, keys, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func fnCloneSample(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		srcCollectionName      string
		trainingCollectionName string
		testCollectionName     string
		keys                   []string
		src                    []byte
		err                    error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch len(args) {
	case 0:
		fmt.Fprintf(eout, "Missing source, training and test collections name\n")
		return 1
	case 1:
		fmt.Fprintf(eout, "Missing training and test collections name\n")
		return 1
	case 2:
		srcCollectionName, trainingCollectionName = args[0], args[1]
	case 3:
		srcCollectionName, trainingCollectionName, testCollectionName = args[0], args[1], args[2]
	default:
		fmt.Fprintf(eout, "Too many parameters %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(srcCollectionName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if inputFName != "" {
		if inputFName == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		keys = strings.Split(string(src), "\n")
	} else {
		keys = c.Keys()
	}
	if len(keys) == 0 {
		fmt.Fprintf(eout, "No objects to clone\n")
		return 1
	}
	// NOTE: Default Sample size is 10% of keys rounded down to nearest in
	if size == 0 {
		size = int(math.Floor(float64(len(keys)) * 0.10))
	}

	err = c.CloneSample(trainingCollectionName, testCollectionName, keys, size, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

func main() {
	app := cli.NewCli(dataset.Version)
	app.SetParams("COLLECTION", "VERB", "[VERB OPTIONS]", "[VERB PARAMETERS ...]")

	// Add Help Docs
	app.SectionNo = 1 // The manual page section number
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))

	topics := []string{}
	for k, v := range Examples {
		app.AddHelp(k, v)
		topics = append(topics, k)
	}
	if len(Examples) > 0 {
		app.AddHelp("examples", []byte(fmt.Sprintf(`%s

To view a specific example use --help EXAMPLE\_NAME where EXAMPLE\_NAME is one of the following: %s`, examples, strings.Join(topics, ", "))))
	}

	app.AddHelp("bugs", []byte(bugs))

	// Add Help Docs
	for k, v := range Help {
		app.AddHelp(k, v)
	}

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

	// Application Verbs
	app.VerbsRequired = true

	// Collection oriented functions
	vInit = app.NewVerb("init", "initialize a collection", fnInit)
	vInit.SetParams("COLLECTION")
	vInit.StringVar(&collectionLayout, "layout", "pairtree", "set file layout for a new collection (i.e. \"buckets\" or \"pairtree\")")
	vStatus = app.NewVerb("status", "collection status", fnStatus)
	vStatus.SetParams("COLLECTION")
	vCheck = app.NewVerb("check", "check a collection for errors", fnCheck)
	vCheck.SetParams("COLLECTION", "[COLLECTION ...]")
	vRepair = app.NewVerb("repair", "repair a collection", fnRepair)
	vRepair.SetParams("COLLECTION")
	vMigrate = app.NewVerb("migrate", "migrate a collection's layout", fnMigrate)
	vMigrate.SetParams("COLLECTION", "LAYOUT")
	vClone = app.NewVerb("clone", "clone a collection", fnClone)
	vClone.SetParams("SRC_COLLECTION", "DEST_COLLECTION")
	vClone.StringVar(&inputFName, "i,input", "", "read key(s), one per line, from a file")

	vClone.BoolVar(&showVerbose, "v,verbose", false, "verbose output")
	vCloneSample = app.NewVerb("clone-sample", "clone a sample from a collection", fnCloneSample)
	vCloneSample.SetParams("SOURCE_COLLECTION", "SAMPLE_COLLECTION", "[TEST_COLLECTION]")
	vCloneSample.StringVar(&inputFName, "i,input", "", "read key(s), one per line, from a file")
	vCloneSample.IntVar(&size, "sample", -1, "set sample size")
	vCloneSample.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	// Object oriented functions
	vCreate = app.NewVerb("create", "create a JSON object", fnCreate)
	vCreate.SetParams("COLLECTION", "KEY", "[JSON_SRC|JSON_FILENAME]")
	vCreate.StringVar(&inputFName, "i,input", "", "input file to read JSON object source from")
	vCreate.BoolVar(&overwrite, "O,overwrite", false, "treat as an update if object already exists")

	vRead = app.NewVerb("read", "read a JSON object from key(s)", fnRead)
	vRead.SetParams("COLLECTION", "[KEY]", "[KEY ...]")
	vRead.StringVar(&inputFName, "i,input", "", "read key(s), one per line, from a file")
	vRead.BoolVar(&prettyPrint, "p,pretty", false, "pretty print JSON output")

	vUpdate = app.NewVerb("update", "update a JSON object", fnUpdate)
	vUpdate.SetParams("COLLECTION", "KEY", "[JSON_SRC|JSON_FILENAME]")
	vUpdate.StringVar(&inputFName, "i,input", "", "input file to read JSON object source from")

	vDelete = app.NewVerb("delete", "delete a JSON object", fnDelete)
	vDelete.SetParams("COLLECTION", "[KEY]", "[KEY ...]")
	vDelete.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vJoin = app.NewVerb("join", "join attributes to a JSON object", fnJoin)
	vJoin.SetParams("COLLECTION", "KEY", "[JSON_SRC|JSON_FILENAME]")
	vJoin.StringVar(&inputFName, "i,input", "", "read JSON source from file")
	vJoin.BoolVar(&overwrite, "overwrite", false, "if true replace attributes otherwise append only new attributes")

	vKeys = app.NewVerb("keys", "list keys in collection", fnKeys)
	vKeys.SetParams("COLLECTION", "[FILTER_EXPR]", "[SORT_EXPR]", "[KEY ...]")
	vKeys.IntVar(&sampleSize, "sample", -1, "set a sample size for keys returned")
	vKeys.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vHasKey = app.NewVerb("haskey", "check for key(s) in collection", fnHasKey)
	vHasKey.SetParams("COLLECTION", "[KEY]", "[KEY ...]")
	vHasKey.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vCount = app.NewVerb("count", "count JSON objects", fnCount)
	vCount.SetParams("COLLECTION", "[FILTER_EXPR]", "[KEY ...]")
	vCount.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	vPath = app.NewVerb("path", "path to JSON object", fnPath)
	vPath.SetParams("COLLECTION", "[KEY]", "[KEY ...]")
	vPath.StringVar(&inputFName, "i,input", "", "read keys, one per line, from a file")

	// Attachment handling
	vAttach = app.NewVerb("attach", "attach a file to JSON object", fnAttach)
	vAttach.SetParams("COLLECTION", "KEY", "[FILENAMES]")
	vAttach.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	vAttachments = app.NewVerb("attachments", "list attachments for a JSON object", fnAttachments)
	vAttachments.SetParams("COLLECTION", "KEY")
	vAttachments.StringVar(&inputFName, "i,input", "", "read keys(s), one per line, from a file")

	vDetach = app.NewVerb("detach", "detach a copy of the attachment from a JSON object", fnDetach)
	vDetach.SetParams("COLLECTION", "KEY", "[FILENAMES]")
	vDetach.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	vPrune = app.NewVerb("prune", "prune an the attachment to a JSON object", fnPrune)
	vPrune.SetParams("COLLECTION", "KEY", "[FILENAMES]")
	vPrune.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	// Frames and Grid
	vGrid = app.NewVerb("grid", "create a 2D JSON array from JSON objects", fnGrid)
	vGrid.SetParams("COLLECTION", "DOTPATH", "[DOTPATH ...]")
	vGrid.BoolVar(&useAllKeys, "a,all", false, "use all keys in a collection")
	vGrid.StringVar(&inputFName, "i,input", "", "use only the keys, one per line, from a file")
	vGrid.StringVar(&filterExpr, "filter", "", "apply filter for inclusion in grid")
	vGrid.IntVar(&sampleSize, "s,sample", -1, "make grid based on a key sample of a given size")
	vGrid.BoolVar(&showVerbose, "v,verbose", showVerbose, "verbose reporting for grid generation")
	vGrid.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vFrame = app.NewVerb("frame", "create or retrieve a data frame", fnFrame)
	vFrame.SetParams("COLLECTION", "FRAME_NAME", "DOTPATH", "[DOTPATH ...]")
	vFrame.BoolVar(&useAllKeys, "a,all", false, "use all keys for frame (one time)")
	vFrame.StringVar(&inputFName, "i,input", "", "use only the keys, one per line, from a file")
	vFrame.StringVar(&filterExpr, "filter", "", "apply filter for inclusion in frame")
	vFrame.StringVar(&sortExpr, "sort", "", "apply sort expression for keys/grid in frame")
	vFrame.IntVar(&sampleSize, "s,sample", -1, "make frame based on a key sample of a given size")
	vFrame.BoolVar(&showVerbose, "v,verbose", showVerbose, "verbose reporting for frame generation")
	vFrame.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	// NOTE: Labels are used with sync-send/sync-receive to map dotpaths to column names
	vFrameLabels = app.NewVerb("frame-labels", "set labels for all columns in a frame", fnFrameLabels)
	vFrameLabels.SetParams("COLLECTION", "FRAME_NAME", "LABEL", "[LABEL ...]")

	// NOTE: Types are used  when defining search indexes
	vFrameTypes = app.NewVerb("frame-types", "set the types for all columns in a frame", fnFrameTypes)
	vFrameTypes.SetParams("COLLECTION", "FRAME_NAME", "TYPE", "[TYPE ...]")

	vReframe = app.NewVerb("reframe", "re-generate an existing frame", fnReframe)
	vReframe.SetParams("COLLECTION", "FRAME_NAME")
	vReframe.BoolVar(&useAllKeys, "a,all", false, "use all keys for frame (one time)")
	vReframe.BoolVar(&setAllKeys, "set-all-keys", false, "persist all keys setting")
	vReframe.StringVar(&inputFName, "i,input", "", "frame only the keys listed in the file, one key per line")
	vReframe.StringVar(&filterExpr, "filter", "", "apply and replace filter expression")
	vReframe.StringVar(&sortExpr, "sort", "", "apply sort expression for keys/grid in frame")
	vReframe.IntVar(&sampleSize, "s,sample", -1, "reframe based on a key sample of a given size")
	vReframe.BoolVar(&showVerbose, "v,verbose", false, "use verbose output")
	vReframe.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vFrames = app.NewVerb("frames", "list frames in a collection", fnFrames)
	vFrames.SetParams("COLLECTION")

	vHasFrame = app.NewVerb("hasframe", "see if a frame has been defined", fnHasFrame)
	vHasFrame.SetParams("COLLECTION", "FRAME_NAME")

	vFrameDelete = app.NewVerb("delete-frame", "delete a frame from a collection", fnFrameDelete)
	vFrameDelete.SetParams("COLLECTION", "FRAME_NAME")

	// Search and indexing
	vFind = app.NewVerb("find", "find a JSON object base on a dot path and value", fnFind)
	vFind.SetParams("INDEX_NAME(S)", "QUERY_STRING")
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
	vFind.IntVar(&batchSize, "batch,size", 100, "set the number of objects per response")
	vIndexer = app.NewVerb("indexer", "index JSON object(s) in a collection", fnIndexer)
	vIndexer.SetParams("COLLECTION", "INDEX_NAME", "INDEX_MAP_FILE", "[KEY ...]")
	vIndexer.StringVar(&inputFName, "i,input", "", "read key(s), one per line, from a file")
	vIndexer.IntVar(&batchSize, "batch,size", 100, "set the number of objects per response")
	vIndexer.BoolVar(&showVerbose, "v,verbose", false, "verbose output")
	vDeindexer = app.NewVerb("deindexer", "remove a JSON object from an index", fnDeindexer)
	vDeindexer.SetParams("INDEX_NAME", "[KEY]", "[KEY ...]")
	vDeindexer.StringVar(&inputFName, "i,input", "", "read key(s), one per line, from a file")
	vDeindexer.IntVar(&batchSize, "batch,size", 0, "set the number of objects per response")

	// Import/export collections from/into tables
	vImport = app.NewVerb("import", "import from a table (CSV, GSheet) into a collection of JSON objects", fnImport)
	vImport.SetParams("COLLECTION", "CSV_FILENAME|GSHEET_ID SHEET_NAME", "ID_COL_NO", "[CELL_RANGE]")
	vImport.StringVar(&clientSecretFName, "client-secret", "", "(import from GSheet) set the client secret path and filename for GSheet access")
	vImport.BoolVar(&useHeaderRow, "use-header-row", true, "use the header row as attribute names in the JSON object")
	vImport.BoolVar(&overwrite, "O,overwrite", false, "overwrite existing JSON objects")
	vImport.BoolVar(&showVerbose, "v,verbose", false, "verbose output")
	vExport = app.NewVerb("export", "export a collection's frame of JSON objects into a table (CSV, GSheet)", fnExport)
	vExport.SetParams("COLLECTION", "CSV_FILENAME|GSHEET_ID SHEET_NAME", "FRAME_NAME")
	vExport.StringVar(&clientSecretFName, "client-secret", "", "(export into a GSheet) set the client secret path and filename for GSheet access")
	vExport.BoolVar(&useHeaderRow, "use-header-row", true, "insert a header row in sheet")
	vExport.BoolVar(&overwrite, "O,overwrite", false, "overwrite existing cells")
	vExport.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	// Synchronize (send/receive) collections of objects with tables using frames
	vSyncSend = app.NewVerb("sync-send", "sync a frame of objects sending data to a table (e.g. CSV, GSheet)", fnSyncSend)
	vSyncSend.SetParams("COLLECTION", "FRAME_NAME", "[CSV_FILENAME|GSHEET_ID SHEET_NAME [CELL_RANGE]]")
	vSyncSend.StringVar(&clientSecretFName, "client-secret", "", "(sync-send to a GSheet) set the client secret path and filename for GSheet access")
	vSyncSend.StringVar(&inputFName, "i,input", "", "read CSV content from a file")
	vSyncSend.StringVar(&outputFName, "o,output", "", "write CSV content to a file")
	vSyncSend.BoolVar(&syncOverwrite, "O,overwrite", true, "overwrite existing cells in table")
	vSyncSend.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	vSyncRecieve = app.NewVerb("sync-recieve", "sync a frame of objects recieving data from a table (e.g. CSV, GSheet)", fnSyncRecieve)
	vSyncRecieve.SetParams("COLLECTION", "FRAME_NAME", "CSV_FILENAME|GSHEET_ID SHEET_NAME")
	vSyncRecieve.StringVar(&clientSecretFName, "client-secret", "", "(sync-receive from a GSheet) set the client secret path and filename for GSheet access")
	vSyncRecieve.StringVar(&inputFName, "i,input", "", "read CSV content from a file")
	vSyncRecieve.BoolVar(&syncOverwrite, "O,overwrite", true, "overwrite existing cells in frame")
	vSyncRecieve.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

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
