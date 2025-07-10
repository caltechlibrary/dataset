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

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following will be set when version.go is generated
	version := dataset.Version
	releaseDate := dataset.ReleaseDate
	releaseHash := dataset.ReleaseHash
	fmtHelp := dataset.FmtHelp
	overwrite, comma, comment := false, "", ""
	lazyQuotes, trimSpace := false, false
	dsqueryHelpText, helpText := dataset.DSQueryHelpText, dataset.DSImporterHelpText
	datasetdHelpText, apiText, serviceText, yamlText := dataset.DatasetdHelpText, dataset.DatasetdApiText, dataset.DatasetdServiceText, dataset.DatasetdYAMLText
	datasetHelpText := dataset.DatasetHelpText
	
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

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		//NOTE: handle help for topic request
		topic := "help"
		if (len(args) > 0) {
			topic = args[0];
		}
		switch topic {
		case "dsquery":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsqueryHelpText, "dsquery", version, releaseDate, releaseHash))
		case "query":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsqueryHelpText, "dsquery", version, releaseDate, releaseHash))
		case "dataset":
			fmt.Fprintf(out, "%s\n", fmtHelp(datasetHelpText, "datasetd", version, releaseDate, releaseHash))
		case "datasetd":
			fmt.Fprintf(out, "%s\n", fmtHelp(datasetdHelpText, "datasetd", version, releaseDate, releaseHash))
		case "api":
			fmt.Fprintf(out, "%s\n", fmtHelp(apiText, "datasetd", version, releaseDate, releaseHash))	
		case "service":
			fmt.Fprintf(out, "%s\n", fmtHelp(serviceText, "datasetd", version, releaseDate, releaseHash))	
		case "yaml":
			fmt.Fprintf(out, "%s\n", fmtHelp(yamlText, "datasetd", version, releaseDate, releaseHash))	
		case "config":
			fmt.Fprintf(out, "%s\n", fmtHelp(yamlText, "datasetd", version, releaseDate, releaseHash))	
		default:
			fmt.Fprintf(os.Stdout, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		}
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
