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
	usage = `USAGE: %s [OPTIONS] COMMAND [COLLECTION PATH| FILENAME | JSON DATA]`

	description = `
SYNOPSIS

%s is a command line tool demonstrating %s package for managing 
JSON documents stored on disc. A dataset stores one of more collections, 
collections store the a buckted distribution of documents
as well as minimal metadata about the collection.

COMMANDS

+ init - initialize a new collection, requires a path to collection
	+ once collection is created, set the environment variable %s_COLLECTION
	  to collection name
+ create - creates a new JSON doc in collection
	+ requires JSON doc name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON doc to stdout
	+ requires JSON doc name
+ update - updates a JSON doc in collection
	+ requires JSON doc name, followed by replacement JSON blob or 
	  JSON blob read from stdin
+ delete - removes a JSON doc from collection
	+ requires JSON doc name
+ keys - returns the keys to stdout, one key per line
`

	examples = `
EXAMPLE

This is an example of creating a dataset called testdata/friends, saving
a record called "littlefreda.json" and reading it back.

   %s init testdata/friends
   export DATASET_COLLECTION=testdata/friends
   %s create littlefreda.json '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(%s keys); do
   	 %s read KY
   done   
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// App Specific Options
	collectionName string

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":   collectionInit,
		"create": createJSONDoc,
		"read":   readJSONDoc,
		"update": updateJSONDoc,
		"delete": deleteJSONDoc,
		"keys":   collectionKeys,
	}

	// alphabet to use for buckets
	alphabet = `abcdefghijklmnopqrstuvwxyz`
)

//
// These are verbs used in the command line utility
//

// collectionInit takes a name (e.g. directory path dataset/mycollection) and
// creates a new collection structure on disc
func collectionInit(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	name := args[0]
	if len(name) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	c, err := dataset.Create(name, dataset.GenerateBucketNames(alphabet, 2))
	if err != nil {
		return "", err
	}
	defer c.Close()
	return fmt.Sprintf("export DATASET_COLLECTION=%q", path.Join(c.Dataset, c.Name)), nil
}

// createJSONDoc adds a new JSON document to the collection
func createJSONDoc(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Expected a doc name and JSON blob")
	}
	name, src := args[0], args[1]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if strings.HasSuffix(name, ".json") == false {
		name = name + ".json"
	}
	if len(src) == 0 {
		return "", fmt.Errorf("missing JSON source")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if err := collection.CreateAsJSON(name, []byte(src)); err != nil {
		return "", err
	}
	return "OK", nil
}

// readJSONDoc returns the JSON from a document in the collection
func readJSONDoc(args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Missing document name")
	}
	name := args[0]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if strings.HasSuffix(name, ".json") == false {
		name = name + ".json"
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	src, err := collection.ReadAsJSON(name)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

// updateJSONDoc replaces a JSON document in the collection
func updateJSONDoc(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Expected document name and JSON blob")
	}
	name, src := args[0], args[1]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if len(src) == 0 {
		return "", fmt.Errorf("missing JSON source")
	}
	if strings.HasSuffix(name, ".json") == false {
		name = name + ".json"
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if err := collection.UpdateAsJSON(name, []byte(src)); err != nil {
		return "", err
	}
	return "OK", nil
}

// deleteJSONDoc removes a JSON document from the collection
func deleteJSONDoc(args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Missing document name")
	}
	name := args[0]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if strings.HasSuffix(name, ".json") == false {
		name = name + ".json"
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if err := collection.Delete(name); err != nil {
		return "", err
	}
	return "OK", nil
}

// collectionKeys returns the keys in a collection
func collectionKeys(args ...string) (string, error) {
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	return strings.Join(collection.Keys(), "\n"), nil
}

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")

	// Application Options
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName, appName, appName)

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
	collectionName = cfg.MergeEnv("collection", collectionName)

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}
	action, params := args[0], args[1:]

	if fn, ok := voc[action]; ok == true {
		output, err := fn(params...)
		if err != nil {
			fmt.Printf("Error %s\n", err)
			os.Exit(1)
		}
		fmt.Println(output)
	} else {
		fmt.Printf("Don't understand %s\n", action)
		os.Exit(1)
	}
}
