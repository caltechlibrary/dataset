package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// CaltechLibrary Packages
	"github.com/github.com/caltechlibrary/cli"
	"github.com/github.com/caltechlibrary/dataset"
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
	voc = map[string]string{
		"init": func(name string) error {
			if len(name) == 0 {
				return "", fmt.Errorf("missing a collection name")
			}
			collection, err := dataset.Create(name, dataset.GenerateBucketNames(alphabet, 2))
			if err != nil {
				return "", err
			}
			defer collection.Close()
			return "OK", nil
		},
		"create": func(name, src string) error {
			if len(collectionName) == 0 {
				return "", fmt.Errorf("missing a collection name")
			}
			if len(name) == 0 {
				return "", fmt.Errorf("missing document name")
			}
			if len(data) == 0 {
				return "", fmt.Errorf("missing JSON data")
			}
			collection, err := dataset.Open(collectionName)
			if err != nil {
				return "", err
			}
			defer collection.Close()

			m := map[string]interface{}{}
			if err := json.Unmarshal([]byte(src), m); err != nil {
				return fmt.Errorf("json formatting error, %s", err)
			}
			if err := collection.Create(name, m); err != nil {
				return "", err
			}
			return "OK", nil
		},
		"read": func(name) (string, error) {
			if len(collectionName) == 0 {
				return "", fmt.Errorf("missing a collection name")
			}
			if len(name) == 0 {
				return "", fmt.Errorf("missing document name")
			}
			collection, err := dataset.Open(collectionName)
			if err != nil {
				return "", err
			}
			defer collection.Close()

			m := map[string]interface{}{}
			if err := collection.Read(name, m); err != nil {
				return nil, err
			}
			src, err := json.Marshal(m)
			if err != nil {
				return "", err
			}
			return string(src), nil
		},
		"update": func(name, src string) (string, error) {
			if len(collectionName) == 0 {
				return "", fmt.Errorf("missing a collection name")
			}
			if len(name) == 0 {
				return "", fmt.Errorf("missing document name")
			}
			if len(data) == 0 {
				return "", fmt.Errorf("missing JSON data")
			}
			collection, err := dataset.Open(collectionName)
			if err != nil {
				return "", err
			}
			defer collection.Close()

			m := map[string]interface{}{}
			if err := json.Unmarshal([]byte(src), m); err != nil {
				return "", fmt.Errorf("json formatting error, %s", err)
			}
			if err := collection.Update(name, m); err != nil {
				return "", err
			}
			return "OK", nil
		},
		"delete": func(name string) (string, error) {
			if len(collectionName) == 0 {
				return "", fmt.Errorf("missing a collection name")
			}
			if len(name) == 0 {
				return "", fmt.Errorf("missing document name")
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
		},
		"keys": func() (string, error) {
			if len(collectionName) == 0 {
				return nil, fmt.Errorf("missing a collection name")
			}
			if len(name) == 0 {
				return nil, fmt.Errorf("missing document name")
			}
			collection, err := dataset.Open(collectionName)
			if err != nil {
				return nil, err
			}
			defer collection.Close()
			return strings.Join(collection.Keys, "\n"), nil
		},
	}
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(license, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
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

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
	}
	action, params := args[0], args[1:]

	if fn, ok := voc[action]; ok == true {
		output, err := fn(params...)
		if err != nil {
			fmt.Printf("Error %s", err)
			os.Exit(1)
		}
		fmt.Println(output)
	} else {
		fmt.Printf("Don't understand %s", action)
		os.Exit(1)
	}
}
