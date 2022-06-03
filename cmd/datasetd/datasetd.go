//
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
//
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
	"github.com/caltechlibrary/dataset/api"
	"github.com/caltechlibrary/dataset/cli"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	description = api.WebDescription
	examples    = api.WebExamples
	license     = api.License
)

func main() {
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display detailed help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		cli.DisplayUsage(os.Stdout, appName, flagSet, description, examples, license)
		os.Exit(0)
	}

	if showLicense {
		cli.DisplayLicense(os.Stdout, appName, license)
		os.Exit(0)
	}

	if showVersion {
		cli.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}

	/* Looking for settings.json */
	settings := "settings.json"
	if len(args) > 0 {
		settings = args[0]
	}
	if _, err := os.Stat(settings); err != nil {
		fmt.Fprintf(os.Stderr, `Could not find %s

Try %s --help for usage details
`, settings, appName)
		os.Exit(1)
	}

	cfg, err := api.LoadConfig(settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cound not read configuration %q, %s", settings, err)
		os.Exit(1)
	}

	/* Open SQL database holding collections */
	if err := api.OpenCollections(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "OpenCollections(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}

	/* Run API */
	if err := api.RunAPI(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "RunWebAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
}
