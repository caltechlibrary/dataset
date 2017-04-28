package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	usage = `USAGE: %s [OPTIONS] INDEX_DEFINITION [INDEX_NAME]`

	description = `
SYNOPSIS

%s is a command line tool for creating a Bleve index based on records in a dataset 
collection. %s reads a JSON document for the index definition and uses that to
configure the Bleve index built based on the dataset collection. If an index
name is not provided then the index name will be the same as the collection with
the file extension of "bleve".

A index definition is JSON document where the indexable record is defined
along with dot paths into the JSON collection record being indexed.

If your collection has records that look like

    {"name":"Frieda Kahlo","occupation":"artist","id":"Frida_Kahlo","dob":"1907-07-06"}

and your wanted an index of names and occupation then your index definition file could
look like

   {
	   "name":{
		   "object_path": ".name",
	   },
	   "occupation": {
		   "object_path":".occupation"
	   }
   }

Based on this definition the "id" and "dob" fields would not be included in the index.`

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
	definitionFName := ""
	indexName := ""
	if len(args) == 1 {
		definitionFName = args[0]
		indexName = fmt.Sprintf("%s.bleve", collectionName)
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

	if err = collection.Indexer(indexName, definitionFName); err != nil {
		fmt.Fprintf(os.Stderr, "Can't build index %s, %s\n", indexName, err)
		os.Exit(1)
	}
}
