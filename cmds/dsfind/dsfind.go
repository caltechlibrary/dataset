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
	usage = `USAGE: %s [OPTIONS] SEARCH_STRINGS`

	description = `
SYNOPSIS

%s is a command line tool for querying a Bleve indexes based on the records in a 
dataset collection. By default %s is assumed there is an index named after the 
collection. An option lets you choose different indexes to query. Results are 
written to standard out and are paged. Options can be used to modify the type
of queries submitted as well as indexes and what is output.`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s -c characters "Jack Flanders"

This would search the Bleve index named characters.bleve for the string "Jack Flanders" 
returning records that matched based on how the index was defined.`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	collectionName string
	indexNames     string
	showHighlight  bool
	resultFields   string
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
	flag.StringVar(&indexNames, "indexes", "", "a colon delimited list of index names")
	flag.BoolVar(&showHighlight, "highlight", false, "display highlight in search results")
	flag.StringVar(&resultFields, "fields", "*", "colon delimited list of fields to display in the results, defaults to *")
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
	if len(indexNames) == 0 {
		indexNames = fmt.Sprintf("%s.bleve", collectionName)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}
	options := map[string]string{}
	if showHighlight == true {
		options["highlight"] = "true"
		options["highlighter"] = "ansi"
	}
	if resultFields != "" {
		options["result_fields"] = strings.TrimSpace(resultFields)
	}

	if err := dataset.Find(os.Stdout, strings.Split(indexNames, ":"), args, options); err != nil {
		fmt.Fprintf(os.Stderr, "Can't search index %s, %s\n", indexNames, err)
		os.Exit(1)
	}
}
