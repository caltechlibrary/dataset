//
// dsindexer creates Blevesearch indexes for a dataset collection. These can be used by
// both dsfind and dsws (web server).
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
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
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
	collectionName string
	documentType   string
	batchSize      int
	updateIndex    bool
	idListFName    string
	goMaxProcs     int
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
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
	flag.StringVar(&collectionName, "collection", "", "sets the collection to be used")
	flag.StringVar(&documentType, "t", "", "the label of the type of document you are indexing, e.g. accession, agent/person")
	flag.IntVar(&batchSize, "batch", 100, "Set the size index batch, default is 100")
	flag.BoolVar(&updateIndex, "update", false, "updating is slow, use this flag if you want to update an exists")
	flag.StringVar(&idListFName, "id-file", "", "Create/Update an index for the ids in file")
	flag.IntVar(&goMaxProcs, "max-procs", -1, "Change the maximum number of CPUs that can executing simultaniously")
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
		if k != "nav" {
			cfg.AddHelp(k, fmt.Sprintf("%s", v))
		}
	}
	for k, v := range Examples {
		if k != "nav" {
			cfg.AddExample(k, fmt.Sprintf("%s", v))
		}
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
		/*
			if len(args) > 0 {
				fmt.Println(cfg.Example(args...))
			} else {
				fmt.Printf("\n%s", cfg.Example())
			}
		*/
		fmt.Println(cfg.ExampleText)
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

	if goMaxProcs > 0 {
		availableProcs := runtime.GOMAXPROCS(goMaxProcs)
		log.Printf("Using %d of %d CPU", goMaxProcs, availableProcs)
	}

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
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}

	collection, err := dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open dataset collection %s, %s\n", collectionName, err)
		os.Exit(1)
	}
	defer collection.Close()

	// Abort index build if it exists and updateIndex is false
	if updateIndex == false {
		if _, err := os.Stat(indexName); os.IsNotExist(err) == false {
			fmt.Fprintf(os.Stderr, "Index exists, updating requires -update option (can be very slow)\n")
			os.Exit(1)
		}
	}

	// NOTE: If a list of ids is provided create/update the index for those ids only
	var keys []string
	if idListFName != "" {
		if src, err := ioutil.ReadFile(idListFName); err == nil {
			klist := bytes.Split(src, []byte("\n"))
			for _, k := range klist {
				if len(k) > 0 {
					keys = append(keys, fmt.Sprintf("%s", k))
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Can't read %s, %s", idListFName, err)
			os.Exit(1)
		}
	}

	if err = collection.Indexer(indexName, definitionFName, batchSize, keys); err != nil {
		fmt.Fprintf(os.Stderr, "Can't build index %s, %s\n", indexName, err)
		os.Exit(1)
	}
}
