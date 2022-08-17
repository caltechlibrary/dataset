// dataset is a command line tool, Go package, shared library and Python package for working with JSON objects as collections on local disc.
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

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v2"
)

var (
	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
)

func main() {
	appName := path.Base(os.Args[0])

	// Standard Options
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")

	// We're ready to process args
	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		dataset.DisplayUsage(out, appName, flagSet)
		os.Exit(0)
	}
	if showLicense {
		dataset.DisplayLicense(out, appName)
		os.Exit(0)
	}
	if showVersion {
		dataset.DisplayVersion(out, appName)
		os.Exit(0)
	}

	// Application Logic
	err := dataset.RunCLI(in, out, eout, args)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		os.Exit(1)
	}
}
