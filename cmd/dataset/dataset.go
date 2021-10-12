//
// dataset is a command line tool, Go package, shared library and Python package for working with JSON objects as collections on local disc.
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/cli"
	"github.com/caltechlibrary/dataset/tbl"
)

var (
	synopsis = `
_dataset_ is a command line tool for working with JSON objects as
collections on local disc.
`

	description = `
The [dataset](docs/dataset.html) command line tool supports
common data management operations for JSON objects stored
as collections.

Features:

- Basic storage actions (*create*, *read*, *update* and *delete*)
- Listing of collection *keys*
- Import/Export to/from CSV files
- The ability to reshape data by performing simple object *joins*
- The ability to create data *grids* and *frames* from
  keys lists and "dot paths" using a collection's JSON objects

Limitations:

_dataset_ has many limitations, some are listed below

- it is not a multi-process, multi-user data store
`

	examples = `
Below is a simple example of shell based interaction with dataset
a collection using the command line dataset tool.

` + "```shell" + `
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use.
    dataset init friends.ds
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document
    dataset create friends.ds frieda '{"name":"frieda","email":"frieda@inverness.example.org"}'
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
too.
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
	cleanObject      bool
	generateMarkdown bool
	generateManPage  bool
	showVerbose      bool

	// Application Options
	//collectionName string
	// header row defaults to true.
	allKeys           = false
	useHeaderRow      = true
	clientSecretFName string
	overwrite         bool
	syncOverwrite     bool
	batchSize         int
	sampleSize        int
	keyFName          string
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
	setValue       bool // Note: set a collection level metadata value

	// Application Verbs
	vInit         *cli.Verb // init
	vStatus       *cli.Verb // status
	vCreate       *cli.Verb // create
	vRead         *cli.Verb // read
	vUpdate       *cli.Verb // update
	vDelete       *cli.Verb // delete
	vJoin         *cli.Verb // join
	vKeys         *cli.Verb // keys
	vHasKey       *cli.Verb // haskey
	vCount        *cli.Verb // count
	vPath         *cli.Verb // path
	vAttach       *cli.Verb // attach
	vAttachments  *cli.Verb // attachments
	vDetach       *cli.Verb // detach (depreciated)
	vRetrieve     *cli.Verb // retrieve (replaces detach)
	vPrune        *cli.Verb // prune
	vImport       *cli.Verb // import
	vExport       *cli.Verb // export
	vCheck        *cli.Verb // check
	vRepair       *cli.Verb // repair
	vCloneSample  *cli.Verb // clone-sample
	vClone        *cli.Verb // clone
	vFrame        *cli.Verb // frame
	vFrameObjects *cli.Verb // frame-objects
	vFrameGrid    *cli.Verb // frame-grid
	vFrameExists  *cli.Verb // has-frame
	vFrames       *cli.Verb // frames
	vRefresh      *cli.Verb // refresh
	vReframe      *cli.Verb // reframe
	vFrameDelete  *cli.Verb // delete-frame
	vSyncSend     *cli.Verb // sync-send
	vSyncRecieve  *cli.Verb // sync-recieve
	vWho          *cli.Verb // who
	vWhat         *cli.Verb // what
	vWhen         *cli.Verb // when
	vWhere        *cli.Verb // where
	vVersion      *cli.Verb // version of collection (semvar)
	vContact      *cli.Verb // contact info for collection

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
	for _, cName := range args {
		c, err = dataset.Init(cName)
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

// fnWho - given a collection path, sets the names for the Who list.
func fnWho(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	if len(args) < 1 {
		fmt.Fprintf(eout,
			"expected a collection name and/or person(s) name\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
	}
	defer c.Close()
	if setValue {
		who := []string{}
		if len(args) > 1 {
			who = args[1:]
		} else {
			src, err := ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read names, %s\n", err)
				return 1
			}
			who = strings.Split(fmt.Sprintf("%s", src), "\n")
		}
		c.Who = who
		err = c.Save()
		if err != nil {
			fmt.Fprintf(eout, "%s", err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", strings.Join(c.Who, "\n"))
	}
	return 0
}

// fnWhat - given a collection path, add description of collection
func fnWhat(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	if len(args) < 1 {
		fmt.Fprintf(eout, "expected a collection name and description\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Unable to open collection %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if setValue {
		if len(args) > 1 {
			c.What = strings.Join(args[1:], "\n")
			if err := c.Save(); err != nil {
				fmt.Fprintf(eout, "Failed to save metadata for %q, %q\n", cName, err)
			}
		} else {
			src, err := ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read description, %s\n", err)
				return 1
			}
			c.What = fmt.Sprintf("%s", src)
			if err := c.Save(); err != nil {
				fmt.Fprintf(eout, "Failed to save metadata for %q, %q\n", cName, err)
				return 1
			}
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", c.What)
	}
	return 0
}

// fnWhen - given a collection path, add date for collection
func fnWhen(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	if len(args) < 1 {
		fmt.Fprintf(eout, "expected a collection name and date(s)\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "could not open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()

	if setValue {
		if len(args) > 1 {
			c.When = strings.Join(args[1:], "\n")
			if err = c.Save(); err != nil {
				fmt.Fprintf(eout, "failed to save metadata for %q, %s\n", cName, err)
				return 1
			}

		} else {
			src, err := ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read date(s), %s\n", err)
				return 1
			}
			c.When = fmt.Sprintf("%s", src)
			if err := c.Save(); err != nil {
				fmt.Fprintf(eout, "failed to save metadata for %q, %s\n", cName, err)
				return 1
			}
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", c.When)
	}
	return 0
}

// fnWhere - given a collection path, add location for collection
// (e.g. url)
func fnWhere(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	if len(args) < 1 {
		fmt.Fprintf(eout, "expected a collection name and location\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if setValue {
		if len(args) > 1 {
			c.Where = strings.Join(args[1:], "\n")
			if err := c.Save(); err != nil {
				fmt.Fprintf(eout, "failed to save metadata for %q, %s\n", cName, err)
				return 1
			}
		} else {
			src, err := ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read location, %s\n", err)
				return 1
			}
			c.Where = fmt.Sprintf("%s", src)
			if err := c.Save(); err != nil {
				fmt.Fprintf(eout, "failed to save metadata for %q, %s\n", cName, err)
				return 1
			}
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", c.Where)
	}
	return 0
}

// fnVersion - given a collection path, add date for semvar version for collection
func fnVersion(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	if len(args) < 1 {
		fmt.Fprintf(eout, "expected a collection name and semvar verion string\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if setValue {
		src := []byte("")
		if len(args) > 1 {
			src = []byte(strings.Join(args[1:], " "))
		} else {
			src, err = ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read semvar version string, %s\n", err)
				return 1
			}
		}
		semver, err := dataset.ParseSemver(src)
		if err != nil {
			fmt.Fprintf(eout, "failed to parse semvar version string %q, %s\n", src, err)
			return 1
		}
		c.DatasetVersion = semver.String()
		if err := c.Save(); err != nil {
			fmt.Fprintf(eout, "Failed to save metadata %q, %s\n", cName, err)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", c.DatasetVersion)
	}
	return 0
}

// fnContact - given a collection path, add contact info
func fnContact(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()
	if len(args) < 1 {
		fmt.Fprintf(eout, "expected a collection name and/or contact info\n")
		return 1
	}
	cName := args[0]
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if setValue {
		src := []byte("")
		if len(args) > 1 {
			src = []byte(strings.Join(args[1:], "\n"))
		} else {
			src, err = ioutil.ReadAll(in)
			if err != nil {
				fmt.Fprintf(eout, "failed to read contact info, %s\n", err)
				return 1
			}
		}
		c.Contact = fmt.Sprintf("%s", src)
		if err := c.Save(); err != nil {
			fmt.Fprintf(eout, "Failed to save metadata %q, %s\n", cName, err)
			return 1
		}
	} else {
		fmt.Fprintf(out, "%s", c.Contact)
	}
	return 0
}

// fnStatus - given a path see if it is a collection by attempting to "open" it
func fnStatus(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		err error
	)
	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

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
	for _, cName := range args {
		c, err := dataset.Open(cName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if showVerbose {
			fmt.Fprintf(out, "%s, dataset version %s, collection version %s\n", cName, c.DatasetVersion, c.Version)
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
		cName string
		key   string
		src   []byte
		err   error
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
		cName, key = args[0], args[1]
		if inputFName == "-" || inputFName == "" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	case 3:
		cName, key = args[0], args[1]
		// Need to decide if args[2] is JSON source or filename
		if strings.HasPrefix(args[2], "{") && strings.HasSuffix(args[2], "}") {
			src = []byte(args[2])
		} else {
			src, err = ioutil.ReadFile(args[2])
			if err != nil {
				fmt.Fprintf(eout, "Can't read %s, %s\n", args[2], err)
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

	m := map[string]interface{}{}
	if err := dataset.DecodeJSON(src, &m); err != nil {
		fmt.Fprintf(eout, "%s must be a valid JSON Object, %s", key, err)
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open collection %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if c.KeyExists(key) == true && overwrite == true {
		if err := c.UpdateJSON(key, src); err != nil {
			fmt.Fprintf(eout, "failed to update %q in %s, %s\n", key, cName, err)
			return 1
		}
		if quiet == false {
			fmt.Fprintf(out, "OK")
		}
		return 0
	}
	if err := c.CreateJSON(key, src); err != nil {
		fmt.Fprintf(eout, "failed to create %q in %s, %s\n", key, cName, err)
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
		cName string
		keys  []string
		src   []byte
		err   error
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
		cName = args[0]
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
		cName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()
	if len(keys) == 1 {
		m := map[string]interface{}{}
		if err := c.Read(keys[0], m, cleanObject); err != nil {
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
		err := c.Read(key, m, cleanObject)
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
		cName string
		key   string
		src   []byte
		err   error
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
		cName, key = args[0], args[1]
		if inputFName == "-" || inputFName == "" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(inputFName)
		}
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	case 3:
		cName, key = args[0], args[1]
		//NOTE: Check if src is file or a object literal string
		if strings.HasPrefix(args[2], "{") && strings.HasSuffix(args[2], "}") {
			src = []byte(args[2])
		} else {
			src, err = ioutil.ReadFile(args[2])
			if err != nil {
				fmt.Fprintf(eout, "Can't read %s, %s\n", args[2], err)
			}
		}
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	if strings.HasSuffix(key, ".json") {
		key = strings.TrimSuffix(key, ".json")
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()
	m := map[string]interface{}{}
	if err := json.Unmarshal(src, &m); err != nil {
		fmt.Fprintf(eout, "%s must be a valid JSON Object, %s", key, err)
		return 1
	}
	if err := c.Update(key, m); err != nil {
		fmt.Fprintf(eout, "failed to update %s in %s, %s\n", key, cName, err)
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
		cName string
		keys  []string
		src   []byte
		err   error
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
		cName = args[0]
	case len(args) >= 2:
		cName, keys = args[0], args[1:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
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
		cName string
		key   string
		src   []byte
		err   error
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
		cName, key = args[0], args[1]
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
			cName, key, src = args[0], args[1], []byte(args[2])
		} else {
			cName, key = args[0], args[1]
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
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()
	// unmarshal new object
	newObj := map[string]interface{}{}
	if err := json.Unmarshal(src, &newObj); err != nil {
		fmt.Fprintf(eout, "%s must be a valid JSON Object, %s", key, err)
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
		cName string
		keys  []string
		err   error
		src   []byte
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
		cName = args[0]
	case len(args) == 2:
		cName, filterExpr = args[0], args[1]
	case len(args) == 3:
		cName, filterExpr, sortExpr = args[0], args[1], args[2]
	case len(args) > 3:
		cName, filterExpr, sortExpr, keys = args[0], args[1], args[2], args[3:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
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

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		random.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})
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
		cName string
		keys  []string
		err   error
		src   []byte
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
		cName = args[0]
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
	case len(args) >= 2:
		cName, keys = args[0], args[1:]
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

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	for i, key := range keys {
		if i > 0 {
			fmt.Fprintf(out, "\n")
		}
		if c.KeyExists(key) {
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
		cName string
		keys  []string
		err   error
		src   []byte
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
		cName = args[0]
	case len(args) == 2:
		cName, filterExpr = args[0], args[1]
	case len(args) > 2:
		cName, filterExpr, keys = args[0], args[1], args[2:]
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

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// If we don't have a sub selection of keys, get a list of keys
	if len(keys) == 0 {
		keys = c.Keys()
	}

	fmt.Fprintf(out, "%d", len(keys))
	return 0
}

// fnPath - return a path(s) to an object(s) given a key(s)
func fnPath(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName   string
		keys    []string
		src     []byte
		docPath string
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
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		cName = args[0]
	case len(args) >= 2:
		cName, keys = args[0], args[1:]
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

	c, err := dataset.Open(cName)
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
		cName  string
		key    string
		src    []byte
		fNames []string
		err    error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	semver := "v0.0.0"
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name, key, semver and attachment name(s)\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key, semver and attachment name(s)\n")
		return 1
	case len(args) == 2:
		if len(fNames) == 0 {
			fmt.Fprintf(eout, "Missing attachment name(s)\n")
			return 1
		}
		cName, key = args[0], args[1]
	case len(args) == 3:
		cName, key, fNames = args[0], args[1], args[2:]
	case len(args) > 3:
		//Is args[2] a semver or a filename?
		if val, err := dataset.ParseSemver([]byte(args[2])); err == nil {
			semver = val.String()
			cName, key, fNames = args[0], args[1], args[3:]
		} else {
			cName, key, fNames = args[0], args[1], args[2:]
		}
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

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.KeyExists(key) == false {
		fmt.Fprintf(eout, "%q is not in %s\n", key, cName)
		return 1
	}
	for _, fname := range fNames {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			fmt.Fprintf(eout, "%s does not exist\n", fname)
			return 1
		}
	}
	err = c.AttachFiles(key, semver, fNames...)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	return 0
}

// fnAttachments - list the attachments of an object(s) given a key(s)
func fnAttachments(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName       string
		keys        []string
		src         []byte
		attachments []string
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
		fmt.Fprintf(eout, "Missing collection name, key(s)\n")
		return 1
	case len(args) == 1:
		if len(keys) == 0 {
			fmt.Fprintf(eout, "Missing key(s)\n")
			return 1
		}
		cName = args[0]
	case len(args) >= 2:
		cName, keys = args[0], args[1:]
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

	c, err := dataset.Open(cName)
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
		cName  string
		key    string
		src    []byte
		fNames []string
		err    error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	semver := ""
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and key\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key\n")
		return 1
	case len(args) == 2:
		cName, key = args[0], args[1]
	case len(args) == 3:
		cName, key, fNames = args[0], args[1], args[2:]
	case len(args) > 3:
		//Is args[2] a semver or a filename?
		if val, err := dataset.ParseSemver([]byte(args[2])); err == nil {
			semver = val.String()
			cName, key, fNames = args[0], args[1], args[3:]
		} else {
			cName, key, fNames = args[0], args[1], args[2:]
		}
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

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.KeyExists(key) == false {
		fmt.Fprintf(eout, "%q is not in %s", key, cName)
		return 1
	}
	err = c.GetAttachedFiles(key, semver, fNames...)
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
		cName  string
		key    string
		src    []byte
		fNames []string
		err    error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	semver := "v0.0.0"
	switch {
	case len(args) == 0:
		fmt.Fprintf(eout, "Missing collection name and key\n")
		return 1
	case len(args) == 1:
		fmt.Fprintf(eout, "Missing key\n")
		return 1
	case len(args) == 2:
		cName, key = args[0], args[1]
	case len(args) == 3:
		cName, key, fNames = args[0], args[1], args[2:]
	case len(args) >= 3:
		//Is args[2] a semver or a filename?
		if val, err := dataset.ParseSemver([]byte(args[2])); err == nil {
			semver = val.String()
			cName, key, fNames = args[0], args[1], args[3:]
		} else {
			cName, key, fNames = args[0], args[1], args[2:]
		}
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

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	if c.KeyExists(key) == false {
		fmt.Fprintf(eout, "%q is not in %s", key, cName)
		return 1
	}
	err = c.Prune(key, semver, fNames...)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprint(out, "OK")
	}
	return 0
}

// fnFrame - define a data frame and populate it with a list of keys,
// dotpaths and label pairs
//
//     dataset keys collection.ds |\
//        dataset frame collection.ds my-frame \
//             ".creator.given=given_name" \
// 			   ".creator.family=family_name" \
//             ".popular_color[0]=favorite_color"
//
// Verb Options: filter-expression (e.g. -filter),
// key list filename (e.g. -i), sample size (e.g. -sample)
// labels (e.g. -labels)
func fnFrame(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName        string
		frameName    string
		keys         []string
		keyPathPairs []string
		dotPaths     []string
		labels       []string
		src          []byte
		err          error
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
		cName, frameName = args[0], args[1]
	case len(args) >= 3:
		cName, frameName, keyPathPairs = args[0], args[1], args[2:]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	if len(keyPathPairs) > 0 {
		for _, item := range keyPathPairs {
			if strings.Contains(item, "=") == true {
				kp := strings.SplitN(item, "=", 2)
				dotPaths = append(dotPaths, strings.TrimSpace(kp[0]))
				labels = append(labels, strings.TrimSpace(kp[1]))
			} else {
				item = strings.TrimSpace(item)
				dotPaths = append(dotPaths, item)
				labels = append(labels, strings.TrimPrefix(item, "."))
			}
		}
	}

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// Check to see if frame exists...
	if c.FrameExists(frameName) {
		if len(labels) > 0 || len(dotPaths) > 0 || len(filterExpr) > 0 {
			fmt.Fprintf(eout, "frame %q already exists\n", frameName)
			return 1
		}
		f, err := c.FrameRead(frameName)
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
	keys = c.Keys()
	if allKeys == false && len(inputFName) > 0 {
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

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		random.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
	}

	// Run a sanity check before we create a new frame...
	if len(dotPaths) == 0 {
		fmt.Fprintf(eout, "No dotpaths, frame creation aborted\n")
		return 1
	}
	if len(labels) == 0 {
		fmt.Fprintf(eout, "No labels, frame creation aborted\n")
		return 1
	}
	//NOTE: We need to be able to frame an empty collection so we
	// can bring mapped content in from a spreadsheet or CSV file easily.
	if len(keys) == 0 && len(c.KeyMap) > 0 {
		fmt.Fprintf(eout, "No keys, frame creation aborted\n")
		return 1
	}

	// NOTE: We defining a new frame now.
	f, err := c.FrameCreate(frameName, keys, dotPaths, labels, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// NOTE: Make need to make sure we save our additional
	// settings - sampleSize

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

// fnFrameObjects - list the frames object list .
func fnFrameObjects(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
		src       []byte
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
		fmt.Fprintf(eout, "Missing frame name for %s\n", args[0])
		return 1
	case len(args) == 2:
		cName = args[0]
		frameName = args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	f, err := c.FrameRead(frameName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// Handle pretty printing
	if prettyPrint {
		src, err = json.MarshalIndent(f.Objects(), "", "    ")
	} else {
		src, err = json.Marshal(f.Objects())
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	fmt.Fprintf(out, "%s", src)
	return 0
}

// fnFrameGrid - get a 2D JSON array of a frame's object list.
func fnFrameGrid(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
		src       []byte
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
		fmt.Fprintf(eout, "Missing frame name for %s\n", args[0])
		return 1
	case len(args) == 2:
		cName = args[0]
		frameName = args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	f, err := c.FrameRead(frameName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	// Handle pretty printing
	if prettyPrint {
		src, err = json.MarshalIndent(f.Grid(useHeaderRow), "", "    ")
	} else {
		src, err = json.Marshal(f.Grid(useHeaderRow))
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
		cName      string
		frameNames []string
		err        error
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
		cName = args[0]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
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

// fnFrameExists - check if a frame has been defined in collection
func fnFrameExists(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
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
		cName, frameName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	if c.FrameExists(frameName) {
		fmt.Fprintf(out, "true")
	} else {
		fmt.Fprintf(out, "false")
	}
	return 0
}

// fnFrameDelete - remove a frame from a collection
func fnFrameDelete(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
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
		cName, frameName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()
	err = c.FrameDelete(frameName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnReframe updates a Frame's object list from the current state
// of collection using the existing keys or the keys supplied.
//
//    dataset reframe -i keys.txt collections.ds my-frame
//
func fnReframe(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		keys      []string
		src       []byte
		err       error
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
		cName, frameName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()

	// Check to see if frame exists...
	if c.FrameExists(frameName) == false {
		fmt.Fprintf(eout, "Frame %q not defined in %s\n", frameName, cName)
		return 1
	}

	keys = c.FrameKeys(frameName)

	// Read from inputFName, update frame's keys
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

	// Apply Sample Size
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	if sampleSize <= len(keys) && sampleSize > 0 {
		keys = keys[0:sampleSize]
	}

	if len(keys) == 0 {
		fmt.Fprintf(eout, "No keys available to update frame\n")
		return 1
	}

	// Now regenerate the objects in the frame
	err = c.FrameReframe(frameName, keys, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnRefresh updates a Frame's object list from the current state
// of collection using the existing frame keys.
//
//    dataset refresh collections.ds my-frame
//
func fnRefresh(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
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
		cName, frameName = args[0], args[1]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}
	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "Failed to open %q, %s\n", cName, err)
		return 1
	}
	defer c.Close()

	// Check to see if frame exists...
	if c.FrameExists(frameName) == false {
		fmt.Fprintf(eout, "Frame %q not defined in %s\n", frameName, cName)
		return 1
	}

	// Now regenerate object list with Refresh
	err = c.FrameRefresh(frameName, showVerbose)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnImport - import a CSV file into a collection
// syntax: COLLECTION CSV_FILENAME ID_COL
//
// options:
// -overwrite
// -use-header-row
// -verbose
//
func fnImport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName         string
		csvFName      string
		idColNoString string
		idCol         int
		err           error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) <= 1:
		fmt.Fprintf(eout, "Missing CSV filename and table details\n")
		return 1
	case len(args) < 3:
		fmt.Fprintf(eout, "Missing table details (e.g. ID_COL_NO) \n")
		return 1
	case len(args) == 3:
		cName, csvFName, idColNoString = args[0], args[1], args[2]
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
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

	// See if we have a CSV filename
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
	}
	if quiet == false {
		fmt.Fprintf(out, "OK")
	}
	return 0
}

// fnExport - export collection objects to a CSV file (or stdout)
// syntax examples: COLLECTION FRAME [CSV_FILENAME]
//
// options:
// -verbose
// -client-secret
func fnExport(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName     string
		frameName string
		err       error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) <= 1:
		fmt.Fprintf(eout, "Missing collection name, frame name, CSV filename\n")
		return 1
	case len(args) == 2:
		cName, frameName = args[0], args[1]
	case len(args) == 3:
		cName, frameName, outputFName = args[0], args[1], args[2]
		if outputFName != "-" {
			fp, err := os.Create(outputFName)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			defer fp.Close()
			out = fp
		}
	default:
		fmt.Fprintf(eout, "Don't understand parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	c, err := dataset.Open(cName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	defer c.Close()

	// for CSV: COLLECTION FRAME_NAME FILENAME

	// Get Frame
	if c.FrameExists(frameName) == false {
		fmt.Fprintf(eout, "Missing frame %q in %s\n", frameName, cName)
		return 1
	}
	// Get dotpaths and column labels from frame
	f, err := c.FrameRead(frameName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	cnt := 0
	cnt, err = c.ExportCSV(out, eout, f, showVerbose)
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

// fnSyncSend - synchronize a frame sending data to a CSV filename
// syntax: COLLECTION FRAME [CSV_FILENAME]
//
func fnSyncSend(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName       string
		frameName   string
		csvFilename string
		src         []byte
		err         error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) <= 1:
		fmt.Fprintf(eout, "Missing collection name and csv filename\n")
		return 1
	case len(args) == 2:
		cName, frameName = args[0], args[1]
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing csv filename\n")
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
	case len(args) == 3:
		cName, frameName, csvFilename = args[0], args[1], args[2]
		src, err = ioutil.ReadFile(csvFilename)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if len(src) == 0 {
			fmt.Fprintf(eout, "No data in csv file %s\n", csvFilename)
			return 1
		}
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err == nil {
			table = tbl.TableStringToInterface(csvTable)
		}
	}
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	c, err := dataset.Open(cName)
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
		if csvFilename != "" {
			fp, err := os.Create(csvFilename)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			defer fp.Close()
			w := csv.NewWriter(fp)
			w.WriteAll(tbl.TableInterfaceToString(table))
			err = w.Error()
		} else {
			w := csv.NewWriter(out)
			w.WriteAll(tbl.TableInterfaceToString(table))
			err = w.Error()
		}
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

// fnSyncRecieve - synchronize a frame receiving data from a CSV file
func fnSyncRecieve(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
	var (
		cName       string
		frameName   string
		csvFilename string
		src         []byte
		err         error
	)

	err = flagSet.Parse(args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	args = flagSet.Args()

	switch {
	case len(args) <= 1:
		fmt.Fprintf(eout, "Missing collection name, frame name and CSV filename\n")
		return 1
	case len(args) == 2:
		cName, frameName = args[0], args[1]
		if inputFName == "" {
			fmt.Fprintf(eout, "Missing csv filename\n")
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
	case len(args) == 3:
		cName, frameName, csvFilename = args[0], args[1], args[2]
		src, err = ioutil.ReadFile(csvFilename)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
	default:
		fmt.Fprintf(eout, "Too many parameters, %s\n", strings.Join(args, " "))
		return 1
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		table = tbl.TableStringToInterface(csvTable)
	}

	c, err := dataset.Open(cName)
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
	for _, cName := range args {
		err = dataset.Analyzer(cName, showVerbose)
		if err != nil {
			fmt.Fprintf(eout, "error in %q, %s\n", cName, err)
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
	for _, cName := range args {
		err = dataset.Repair(cName, showVerbose)
		if err != nil {
			fmt.Fprintf(eout, "error in %q, %s\n", cName, err)
			return 1
		}
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
	app.BoolVar(&showVerbose, "V,verbose", false, "output rows processed on importing from CSV")

	// Application Verbs
	app.VerbsRequired = true

	// Collection oriented functions
	vInit = app.NewVerb("init", "initialize a collection", fnInit)
	vInit.SetParams("COLLECTION")
	vStatus = app.NewVerb("status", "collection status", fnStatus)
	vStatus.SetParams("COLLECTION")
	vCheck = app.NewVerb("check", "check a collection for errors", fnCheck)
	vCheck.SetParams("COLLECTION", "[COLLECTION ...]")
	vRepair = app.NewVerb("repair", "repair a collection", fnRepair)
	vRepair.SetParams("COLLECTION")
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
	vRead.BoolVar(&cleanObject, "c,clean", false, "Remove dataset underscore variables before returning object")
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
	vAttach.SetParams("COLLECTION", "KEY", "[SEMVER]", "[FILENAMES]")
	vAttach.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	vAttachments = app.NewVerb("attachments", "list attachments for a JSON object", fnAttachments)
	vAttachments.SetParams("COLLECTION", "KEY")
	vAttachments.StringVar(&inputFName, "i,input", "", "read keys(s), one per line, from a file")

	vDetach = app.NewVerb("detach", "(depreciated in favor of retrieve) detach a copy of the attachment from a JSON object", fnDetach)
	vDetach.SetParams("COLLECTION", "KEY", "[SEMVER]", "[FILENAMES]")
	vDetach.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	vRetrieve = app.NewVerb("retrieve", "detach a copy of the attachment from a JSON object", fnDetach)
	vRetrieve.SetParams("COLLECTION", "KEY", "[SEMVER]", "[FILENAMES]")
	vRetrieve.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	vPrune = app.NewVerb("prune", "prune an the attachment to a JSON object", fnPrune)
	vPrune.SetParams("COLLECTION", "KEY", "[SEMVER]", "[FILENAMES]")
	vPrune.StringVar(&inputFName, "i,input", "", "read filename(s), one per line, from a file")

	// Frames
	vFrame = app.NewVerb("frame", "create or retrieve a data frame", fnFrame)
	vFrame.SetParams("COLLECTION", "FRAME_NAME", "DOTPATH", "[DOTPATH ...]")
	vFrame.StringVar(&inputFName, "i,input", "", "use only the keys, one per line, from a file")
	vFrame.StringVar(&filterExpr, "filter", "", "apply filter for inclusion in frame")
	vFrame.StringVar(&sortExpr, "sort", "", "apply sort expression for keys/grid in frame")
	vFrame.IntVar(&sampleSize, "s,sample", -1, "make frame based on a key sample of a given size")
	vFrame.BoolVar(&allKeys, "a,all", allKeys, "Use all collection keys for frame")
	vFrame.BoolVar(&showVerbose, "v,verbose", showVerbose, "verbose reporting for frame generation")
	vFrame.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vFrameObjects = app.NewVerb("frame-objects", "return the object list of a frame", fnFrameObjects)
	vFrameObjects.SetParams("COLLECTION", "FRAME_NAME")
	vFrameObjects.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vFrameGrid = app.NewVerb("frame-grid", "return the object list as a 2D array", fnFrameGrid)
	vFrameGrid.SetParams("COLLECTION", "FRAME_NAME")
	vFrameGrid.BoolVar(&useHeaderRow, "use-header-row", useHeaderRow, "Include labels as a header row")
	vFrameGrid.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vReframe = app.NewVerb("reframe", "replace the objects in a frame using the keys provided.", fnReframe)
	vReframe.SetParams("COLLECTION", "FRAME_NAME")
	vReframe.StringVar(&inputFName, "i,input", "", "frame only the keys listed in the file, one key per line")
	vReframe.IntVar(&sampleSize, "s,sample", -1, "reframe based on a key sample of a given size")
	vReframe.BoolVar(&showVerbose, "v,verbose", false, "use verbose output")
	vReframe.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vRefresh = app.NewVerb("refresh", "refresh an existing frame based on current collection state.", fnRefresh)
	vRefresh.SetParams("COLLECTION", "FRAME_NAME")
	vRefresh.BoolVar(&showVerbose, "v,verbose", false, "use verbose output")
	vRefresh.BoolVar(&prettyPrint, "p,pretty", prettyPrint, "pretty print JSON output")

	vFrames = app.NewVerb("frames", "list frames in a collection", fnFrames)
	vFrames.SetParams("COLLECTION")

	vFrameExists = app.NewVerb("hasframe", "see if a frame has been defined", fnFrameExists)
	vFrameExists.SetParams("COLLECTION", "FRAME_NAME")

	vFrameDelete = app.NewVerb("delete-frame", "delete a frame from a collection", fnFrameDelete)
	vFrameDelete.SetParams("COLLECTION", "FRAME_NAME")

	// Import/export collections from/into tables
	vImport = app.NewVerb("import", "import from a CSV file into a collection of JSON objects", fnImport)
	vImport.SetParams("COLLECTION", "CSV_FILENAME", "ID_COL_NO", "[CELL_RANGE]")
	vImport.BoolVar(&useHeaderRow, "use-header-row", useHeaderRow, "use the header row as attribute names in the JSON object")
	vImport.BoolVar(&overwrite, "O,overwrite", false, "overwrite existing JSON objects")
	vImport.BoolVar(&showVerbose, "v,verbose", false, "verbose output")
	vExport = app.NewVerb("export", "export a collection's frame of JSON objects into a CSV file", fnExport)
	vExport.SetParams("COLLECTION", "FRAME_NAME", "CSV_FILENAME")
	vExport.BoolVar(&useHeaderRow, "use-header-row", useHeaderRow, "insert a header row in sheet")
	vExport.BoolVar(&overwrite, "O,overwrite", false, "overwrite existing cells")
	vExport.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	// Synchronize (send/receive) collections of objects with tables using frames
	vSyncSend = app.NewVerb("sync-send", "sync a frame of objects sending data to a CSV file", fnSyncSend)
	vSyncSend.SetParams("COLLECTION", "FRAME_NAME", "[CSV_FILENAME]")
	vSyncSend.StringVar(&inputFName, "i,input", "", "read CSV content from a file")
	vSyncSend.StringVar(&outputFName, "o,output", "", "write CSV content to a file")
	vSyncSend.BoolVar(&syncOverwrite, "O,overwrite", true, "overwrite existing cells in table")
	vSyncSend.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	vSyncRecieve = app.NewVerb("sync-recieve", "sync a frame of objects recieving data from a CSV file", fnSyncRecieve)
	vSyncRecieve.SetParams("COLLECTION", "FRAME_NAME", "CSV_FILENAME")
	vSyncRecieve.StringVar(&inputFName, "i,input", "", "read CSV content from a file")
	vSyncRecieve.BoolVar(&syncOverwrite, "O,overwrite", true, "overwrite existing cells in frame")
	vSyncRecieve.BoolVar(&showVerbose, "v,verbose", false, "verbose output")

	vWho = app.NewVerb("who", "authorship, owner or maintainer name(s)", fnWho)
	vWho.SetParams("COLLECTION", "[WHO]")
	vWho.BoolVar(&setValue, "set", false, "set the value(s)")
	vWhat = app.NewVerb("what", "description of collection", fnWhat)
	vWhat.SetParams("COLLECTION", "[WHAT]")
	vWhat.BoolVar(&setValue, "set", false, "set the value(s)")
	vWhen = app.NewVerb("when", "created or publication data", fnWhen)
	vWhen.SetParams("COLLECTION", "[WHEN]")
	vWhen.BoolVar(&setValue, "set", false, "set the value(s)")
	vWhere = app.NewVerb("where", "url or description of where to find collection", fnWhere)
	vWhere.SetParams("COLLECTION", "[WHERE]")
	vWhere.BoolVar(&setValue, "set", false, "set the value(s)")
	vVersion = app.NewVerb("version", "version of collection in semvar format", fnVersion)
	vVersion.SetParams("COLLECTION", "[SEMVAR]")
	vVersion.BoolVar(&setValue, "set", false, "set the value(s)")
	vContact = app.NewVerb("contact", "contact info for questions and support", fnContact)
	vContact.SetParams("COLLECTION", "[CONTACT_INFO]")
	vContact.BoolVar(&setValue, "set", false, "set the value(s)")

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
		fmt.Fprintln(app.Out, "%s\n", dataset.License)
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
