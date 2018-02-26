//
// dsindexer creates Blevesearch indexes for a dataset collection. These can be used by
// both dsfind and dsws (web server).
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	// Standard Options
	showHelp             bool
	showLicense          bool
	showVersion          bool
	showExamples         bool
	inputFName           string
	outputFName          string
	newLine              bool
	quiet                bool
	prettyPrint          bool
	generateMarkdownDocs bool

	// App Specific Options
	collectionName  string
	documentType    string
	batchSize       int
	updateIndex     bool
	idListFName     string
	deleteFromIndex bool
)

func init() {

	// Application Options
}

func main() {
	app := cli.NewCli(dataset.Version)
	appName := app.AppName()

	// Add Non-options docs
	app.AddParams("INDEX_DEF_JSON", "INDEX_NAME")

	// Add Help Docs
	for k, v := range Help {
		if k != "nav" {
			app.AddHelp(k, v)
		}
	}
	for k, v := range Examples {
		if k != "nav" {
			app.AddHelp(k, v)
		}
	}

	// Environment Options
	app.EnvStringVar(&collectionName, "DATASET", "", "Set the dataset collection you're working with")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "e,examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&newLine, "nl,newline", true, "if set to false suppress the trailing newline")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// Application Options
	app.StringVar(&collectionName, "c,collection", "", "sets the collection to be used")
	app.StringVar(&documentType, "t", "", "the label of the type of document you are indexing, e.g. accession, agent/person")
	app.IntVar(&batchSize, "batch", 0, "Set the size index batch, default is 100")
	app.BoolVar(&updateIndex, "update", false, "updating is slow, use this app if you want to update an exists")
	app.StringVar(&idListFName, "key-file", "", "Create/Update an index based on the keys provided in the file")
	app.BoolVar(&deleteFromIndex, "delete", false, "this will cause records to be deleted from an index, use with -key-file")

	// Action verbs (e.g. app.AddAction(STRING_VERB, FUNC_POINTER, STRING_DESCRIPTION)
	//FIXME: If the application is verb based add your verbs here

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
	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
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

	// Application Option's processing

	// Merge environment
	datasetEnv := os.Getenv("DATASET")
	if datasetEnv != "" && collectionName == "" {
		collectionName = datasetEnv
	}

	definitionFName := ""
	indexName := ""
	if len(args) == 1 {
		definitionFName = args[0]
		ext := path.Ext(definitionFName)
		if ext != "" {
			indexName = strings.TrimSuffix(definitionFName, ext) + ".bleve"
		} else {
			indexName = path.Base(definitionFName) + ".bleve"
		}
	} else if len(args) == 2 {
		definitionFName, indexName = args[0], args[1]
	} else {
		cli.ExitOnError(app.Eout, fmt.Errorf("See %s --help", appName), quiet)
	}

	collection, err := dataset.Open(collectionName)
	cli.ExitOnError(app.Eout, err, quiet)
	defer collection.Close()

	// Abort index build if it exists and updateIndex is false
	if updateIndex == false {
		if _, err := os.Stat(indexName); os.IsNotExist(err) == false {
			cli.ExitOnError(app.Eout, fmt.Errorf("Index exists, updating requires -update option (can be very slow)"), quiet)
		}
	}

	// NOTE: If a list of ids is provided create/update the index for those ids only
	var keys []string
	if idListFName != "" {
		src, err := ioutil.ReadFile(idListFName)
		cli.ExitOnError(app.Eout, err, quiet)

		klist := bytes.Split(src, []byte("\n"))
		for _, k := range klist {
			if len(k) > 0 {
				keys = append(keys, fmt.Sprintf("%s", k))
			}
		}
	} else {
		keys = collection.Keys()
	}

	if batchSize == 0 {
		if len(keys) > 10000 {
			batchSize = len(keys) / 100
		} else {
			batchSize = 100
		}
	}

	err = collection.Indexer(indexName, definitionFName, batchSize, keys)
	cli.ExitOnError(app.Eout, err, quiet)

	if newLine {
		fmt.Fprintln(app.Out, "")
	}
}
