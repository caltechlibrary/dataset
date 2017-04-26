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
	usage = `USAGE: %s [OPTIONS] INDEX_MAPPING_FILE INDEX_NAME`

	description = `
SYNOPSIS

%s is a command line tool for creating a Bleve index based on records in a dataset 
collection. %s reads a JSON document for the record structure of the index being 
built and saves the result s a bleve index.`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s -c characters email-mapping.json email-index

This will build a Bleve index called "email-index" based on the index defined
in "email-mapping.json".`

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
	definitionFName, indexName := args[0], args[1]

	collection, err := dataset.Open(collectionName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open dataset collection %s, %s\n", collectionName, err)
		os.Exit(1)
	}
	defer collection.Close()

	if mapping, err := dataset.ReadIndexMapFile(args[0]); err == nil {
		if err = collection.Indexer(indexName, mapping); err != nil {
			fmt.Fprintf(os.Stderr, "Can't build index %s, %s\n", indexName, err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Can't load index mapping %s, %s\n", definitionFName, err)
		os.Exit(1)
	}
}
