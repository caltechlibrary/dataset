// datasetd implements a web service for working with dataset collections.
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
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
package main

//
// datasetd provide collection access/management via a simple
// HTTP/HTTPS API
//

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/v2"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool
)

func main() {
	appName := path.Base(os.Args[0])
	version, releaseDate, releaseHash, licenseText  := dataset.Version, dataset.ReleaseDate, dataset.ReleaseHash, dataset.LicenseText
	helpText, apiText, serviceText, yamlText := dataset.DatasetdHelpText, dataset.DatasetdApiText, dataset.DatasetdServiceText, dataset.DatasetdYAMLText
	dsqueryHelpText, dsimporterHelpText := dataset.DSQueryHelpText, dataset.DSImporterHelpText
	datasetHelpText := dataset.DatasetHelpText
	fmtHelp := dataset.FmtHelp

	// Standard Options
	debug := false
	flag.BoolVar(&showHelp, "help", false, "display detailed help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&debug, "debug", debug, "log debugging information")

	flag.Parse()
	args := flag.Args()

	out := os.Stdout
	//eout := os.Stderr

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
		case "dsimporter":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsimporterHelpText, "dsimporter", version, releaseDate, releaseHash))
		case "importer":
			fmt.Fprintf(out, "%s\n", fmtHelp(dsimporterHelpText, "dsimporter", version, releaseDate, releaseHash))
		case "dataset":
			fmt.Fprintf(out, "%s\n", fmtHelp(datasetHelpText, "dataset", version, releaseDate, releaseHash))
		case "api":
			fmt.Fprintf(out, "%s\n", fmtHelp(apiText, appName, version, releaseDate, releaseHash))	
		case "service":
			fmt.Fprintf(out, "%s\n", fmtHelp(serviceText, appName, version, releaseDate, releaseHash))	
		case "yaml":
			fmt.Fprintf(out, "%s\n", fmtHelp(yamlText, appName, version, releaseDate, releaseHash))	
		case "config":
			fmt.Fprintf(out, "%s\n", fmtHelp(yamlText, appName, version, releaseDate, releaseHash))	
		default:
			fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		}
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintf(out, "%s\n", licenseText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintf(out, "%s %s (%s %s)\n", appName, version, releaseDate, releaseHash)
		os.Exit(0)
	}

	/* Looking for settings.json */
	settings := ""
	if len(args) > 0 {
		settings = args[0]
	}
	if _, err := os.Stat(settings); err != nil {
		fmt.Fprintf(os.Stderr, `Could not find %s

Try %s -help for usage details
`, settings, appName)
		os.Exit(1)
	}

	/* Run API */
	if err := dataset.RunAPI(appName, settings, debug); err != nil {
		fmt.Fprintf(os.Stderr, "RunAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
}
