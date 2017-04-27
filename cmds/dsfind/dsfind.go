package main

import (
	"flag"
	"fmt"
	//	"io"
	"os"
	"path"
	//	"strconv"
	//"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	usage = `USAGE: %s [OPTIONS] INDEX_NAME SEARCH_STRING`

	description = `
SYNOPSIS

%s is a command line tool for querying a Bleve index based on records in a dataset 
collection. %s queries the index named based on search string. Results are written
to standard out by default. Options control how the query is processed and how 
results are handled.`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s -c characters email-index "Jack Flanders"

This would search the Bleve index named email-index for the string "Jack Flanders" 
returning records that matched based on how the index was defined.`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	collectionName string
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// Application Options
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
	flag.StringVar(&collectionName, "collection", "", "sets the collection to be used")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
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

	// Merge environment
	datasetEnv := os.Getenv("DATASET")
	if datasetEnv != "" && collectionName == "" {
		collectionName = datasetEnv
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}
	indexName, queryString := args[0], args[1]
	if err := dataset.Find(os.Stdout, indexName, queryString); err != nil {
		fmt.Fprintf(os.Stderr, "Can't search index %s, %s\n", indexName, err)
		os.Exit(1)
	}
}
