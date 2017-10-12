//
// dsfind is a command line utility that will search one or more Blevesearch indexes created by
// dsindexer. The output can be in a number formats included text, CSV and JSON.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
//
// Copyright (c) 2017, Caltech
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
	"os"
	"path"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	// Standard Options
	showHelp     bool
	showLicense  bool
	showVersion  bool
	showExamples bool

	// App Specific Options
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
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")

	// Application Options
	flag.StringVar(&indexList, "indexes", "", "colon or comma delimited list of index names")
	flag.StringVar(&sortBy, "sort", "", "a comma delimited list of field names to sort by")
	flag.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	flag.StringVar(&setHighlighter, "highlighter", "", "set the highlighter (ansi,html) for search results")
	flag.StringVar(&resultFields, "fields", "", "comma delimited list of fields to display in the results")
	flag.BoolVar(&jsonFormat, "json", false, "format results as a JSON document")
	flag.BoolVar(&csvFormat, "csv", false, "format results as a CSV document, used with fields option")
	flag.BoolVar(&csvSkipHeader, "csv-skip-header", false, "don't output a header row, only values for csv output")
	flag.BoolVar(&idsOnly, "ids", false, "output only a list of ids from results")
	flag.IntVar(&size, "size", 0, "number of results returned for request")
	flag.IntVar(&from, "from", 0, "return the result starting with this result number")
	flag.StringVar(&explain, "explain", "", "explain results in a verbose JSON document")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, appName, dataset.Version)
	cfg.LicenseText = fmt.Sprintf(dataset.License, appName, dataset.Version)
	cfg.UsageText = fmt.Sprintf("%s", Help["usage"])
	cfg.DescriptionText = fmt.Sprintf("%s", Help["description"])
	cfg.OptionText = "## OPTIONS\n\n"
	cfg.ExampleText = fmt.Sprintf("%s", Examples["index"])

	// Add help and examples
	for k, v := range Help {
		cfg.AddHelp(k, fmt.Sprintf("%s", v))
	}
	for k, v := range Examples {
		cfg.AddExample(k, fmt.Sprintf("%s", v))
	}

	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}

	if showExamples == true {
		if len(args) > 0 {
			fmt.Println(cfg.Example(args...))
		} else {
			fmt.Println(cfg.ExampleText)
		}
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	// We expect at least one arg, the search string
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, cfg.Usage())
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "Do not know what index to use")
		os.Exit(1)
	}

	options := map[string]string{}
	if explain != "" {
		options["explain"] = "true"
		jsonFormat = true
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

	idxAlias, idxFields, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open index %s, %s\n", strings.Join(indexNames, ", "), err)
		os.Exit(1)
	}
	defer idxAlias.Close()

	results, err := dataset.Find(os.Stdout, idxAlias, args, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't search index %s, %s\n", strings.Join(indexNames, ", "), err)
		os.Exit(1)
	}

	//
	// Handle results formatting choices
	//
	switch {
	case jsonFormat == true:
		if err := dataset.JSONFormatter(os.Stdout, results); err != nil {
			fmt.Fprintf(os.Stderr, "JSON formatting error, %s\n", err)
			os.Exit(1)
		}
	case csvFormat == true:
		var fields []string
		if resultFields == "" {
			fields = idxFields
		} else {
			fields = strings.Split(resultFields, ",")
		}
		if err := dataset.CSVFormatter(os.Stdout, results, fields, csvSkipHeader); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case idsOnly == true:
		for _, hit := range results.Hits {
			fmt.Fprintf(os.Stdout, "%s\n", hit.ID)
		}
	default:
		fmt.Fprintf(os.Stdout, "%s\n", results)
	}
}
