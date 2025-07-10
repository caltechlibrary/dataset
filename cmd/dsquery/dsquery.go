// dsquery.go is a command line program for working an dataset collections using the
// dataset v2 SQL store for the JSON documents (e.g. Postgres, MySQL, SQLite 3).
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
	"io"
	"os"
	"path"
	"strings"

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
	pretty, ptIndex, grid, csv, asYaml := false, false, "", "", ""
	sqlFName := ""
	datasetdHelpText, apiText, serviceText, yamlText := dataset.DatasetdHelpText, dataset.DatasetdApiText, dataset.DatasetdServiceText, dataset.DatasetdYAMLText
	helpText, dsimporterHelpText := dataset.DSQueryHelpText, dataset.DSImporterHelpText
	datasetHelpText := dataset.DatasetHelpText

	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&pretty, "pretty", false, "pretty JSON output")
	flag.StringVar(&grid, "grid", grid, "return JSON grid of values, requires a comma delimited string of attribute names")
	flag.StringVar(&csv, "csv", csv, "return csv file using the attribute names from list of objects")
	flag.StringVar(&asYaml, "yaml", asYaml, "return YAML file using the attribute names from list of objects")
	flag.StringVar(&sqlFName, "sql", sqlFName, "read SQL statement from a file")
	flag.BoolVar(&ptIndex, "index", ptIndex, "create a SQLite 3 'index' for a collection.")
	flag.Parse()
	args := flag.Args()

	//in := os.Stdin
	out := os.Stdout
	//eout := os.Stderr

	if showHelp {
		//NOTE: handle help for topic request
		topic := "help"
		if (len(args) > 0) {
			topic = args[0];
		}
		switch topic {
		case "dsimporter":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsimporterHelpText, "dsimporter", version, releaseDate, releaseHash))
		case "importer":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsimporterHelpText, "dsimporter", version, releaseDate, releaseHash))
		case "dataset":
			fmt.Fprintf(out, "%s\n", fmtHelp(datasetHelpText, "dataset", version, releaseDate, releaseHash))
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
			fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
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
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing C_NAME and SQL_STATEMENT")
		os.Exit(10)
	}

	// Create a DSQuery object and evaluate the command line options
	app := new(dataset.DSQuery)
	cName, stmt, params := "", "", []string{}
	if sqlFName != "" {
		in := os.Stdin
		if sqlFName != "-" {
			var err error
			in, err = os.Open(sqlFName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(10)
			}
			defer in.Close()
		}
		src, err := io.ReadAll(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(10)
		}
		stmt = fmt.Sprintf("%s", src)
	}
	for _, arg := range args {
		switch {
		case cName == "":
			cName = arg
		case stmt == "":
			stmt = arg
		default:
			params = append(params, arg)
		}
	}
	if cName == "" {
		fmt.Fprintf(os.Stderr, "missing C_NAME\n")
		os.Exit(10)
	}
	if stmt == "" {
		fmt.Fprintf(os.Stderr, "missing SQL_STATEMENT\n")
		os.Exit(10)
	}
	app.Pretty = pretty
	app.PTIndex = ptIndex
	attributes := []string{}
	if grid != "" {
		app.AsGrid = true
		app.AsCSV = false
		app.AsYAML = false
		attributes = strings.Split(grid, ",")
	}
	if csv != "" {
		app.AsGrid = false
		app.AsCSV = true
		app.AsYAML = false
		attributes = strings.Split(csv, ",")
	}
	if asYaml != "" {
		app.AsGrid = false
		app.AsCSV = false
		app.AsYAML = true
		attributes = strings.Split(asYaml, ",")
	}
	if app.AsGrid || app.AsCSV || app.AsYAML { 
		app.Attributes = []string{}
		for _, attr := range attributes {
			app.Attributes = append(app.Attributes, strings.TrimSpace(attr))
		}
	}
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr, cName, stmt, params); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(10)
	}
}
