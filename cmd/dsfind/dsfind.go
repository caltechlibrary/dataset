//
// dsfind is a command line utility that will search one or more Blevesearch indexes created by
// dsindexer. The output can be in a number formats included text, CSV and JSON.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
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
	"os"
	"path"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	description = `uses a Bleve index to search a dataset returning results to the command line
`

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

	// Application Options
	indexList      string
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
	explain        string // Note: will be converted to boolean so expecting 1,0,T,F,true,false, etc.
	sampleSize     int
)

func main() {
	app := cli.NewCli(dataset.Version)
	appName := app.AppName()

	// Add Params
	app.AddParams("[INDEX_LIST]", "QUERY_STRING")

	// Add Help Docs
	app.AddHelp("description", []byte(description))
	for k, v := range Help {
		app.AddHelp(k, v)
	}
	for k, v := range Examples {
		app.AddHelp(k, v)
	}

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
	app.StringVar(&indexList, "indexes", "", "colon or comma delimited list of index names")
	app.StringVar(&sortBy, "sort", "", "a comma delimited list of field names to sort by")
	app.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	app.StringVar(&setHighlighter, "highlighter", "", "set the highlighter (ansi,html) for search results")
	app.StringVar(&resultFields, "fields", "", "comma delimited list of fields to display in the results")
	app.BoolVar(&jsonFormat, "json", false, "format results as a JSON document")
	app.BoolVar(&csvFormat, "csv", false, "format results as a CSV document, used with fields option")
	app.BoolVar(&csvSkipHeader, "csv-skip-header", false, "don't output a header row, only values for csv output")
	app.BoolVar(&idsOnly, "ids", false, "output only a list of ids from results")
	app.IntVar(&size, "size", 0, "number of results returned for request")
	app.IntVar(&from, "from", 0, "return the result starting with this result number")
	app.StringVar(&explain, "explain", "", "explain results in a verbose JSON document")
	app.IntVar(&sampleSize, "sample", 0, "return a sample of size N of results")

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

	// We expect at least one arg, the search string
	if len(args) == 0 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("See %s --help", appName), quiet)
	}

	// Handle the case where indexes were listed with the -indexes option like dsfind
	var indexNames []string
	if indexList != "" {
		var delimiter = ","
		if strings.Contains(indexList, ":") {
			delimiter = ":"
		}
		indexNames = strings.Split(indexList, delimiter)
	}

	// Collect any additional index names from the remaining args
	for _, arg := range args {
		if path.Ext(arg) == ".bleve" {
			indexNames = append(indexNames, arg)
		}
	}
	if len(indexNames) == 0 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("Do not know what index to use"), quiet)
	}

	options := map[string]string{}
	if explain != "" {
		options["explain"] = "true"
		jsonFormat = true
	}

	if sampleSize > 0 {
		options["sample"] = fmt.Sprintf("%d", sampleSize)
	}
	if from != 0 {
		options["from"] = fmt.Sprintf("%d", from)
	}
	if size > 0 {
		options["size"] = fmt.Sprintf("%d", size)
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

	idxList, idxFields, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		cli.ExitOnError(os.Stderr, fmt.Errorf("Can't open index %s, %s", strings.Join(indexNames, ", "), err), quiet)
	}
	defer idxList.Close()

	results, err := dataset.Find(idxList.Alias, args, options)
	if err != nil {
		cli.ExitOnError(os.Stderr, fmt.Errorf("Can't search index %s, %s", strings.Join(indexNames, ", "), err), quiet)
	}

	//
	// Handle results formatting choices
	//
	switch {
	case jsonFormat == true:
		err := dataset.JSONFormatter(app.Out, results, prettyPrint)
		cli.ExitOnError(os.Stderr, err, quiet)
	case csvFormat == true:
		var fields []string
		if resultFields == "" {
			fields = idxFields
		} else {
			fields = strings.Split(resultFields, ",")
		}
		err := dataset.CSVFormatter(app.Out, results, fields, csvSkipHeader)
		cli.ExitOnError(os.Stderr, err, quiet)
	case idsOnly == true:
		for _, hit := range results.Hits {
			fmt.Fprintf(app.Out, "%s", hit.ID)
		}
	default:
		fmt.Fprintf(app.Out, "%s", results)
	}

	if newLine {
		fmt.Fprintln(app.Out, "")
	}
}
