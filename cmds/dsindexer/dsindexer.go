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
	usage = `USAGE: %s [OPTIONS] INDEX_DEFINITION [INDEX_NAME]`

	description = `
SYNOPSIS

%s is a command line tool for creating a Bleve index based on records in a dataset 
collection. %s reads a JSON document for the index definition and uses that to
configure the Bleve index built based on the dataset collection. If an index
name is not provided then the index name will be the same as the definition file
with the .json replaced by .bleve.

A index definition is JSON document where the indexable record is defined
along with dot paths into the JSON collection record being indexed.

If your collection has records that look like

    {"name":"Frieda Kahlo","occupation":"artist","id":"Frida_Kahlo","dob":"1907-07-06"}

and your wanted an index of names and occupation then your index definition file could
look like

   {
	   "name":{
		   "object_path": ".name"
	   },
	   "occupation": {
		   "object_path":".occupation"
	   }
   }

Based on this definition the "id" and "dob" fields would not be included in the index.
`

	examples = `
EXAMPLES

In the example the index will be created for a collection called "characters".

    %s -c characters email-mapping.json email-index

This will build a Bleve index called "email-index" based on the index defined
in "email-mapping.json".
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

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

	if goMaxProcs > 0 {
		availableProcs := runtime.GOMAXPROCS(goMaxProcs)
		log.Printf("Using %d of %d CPU", goMaxProcs, availableProcs)
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
