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
	usage = `USAGE: %s [OPTIONS] COMMAND_AND_PARAMETERS`

	description = `
SYNOPSIS

%s is a command line tool demonstrating dataset package for managing 
JSON documents stored on disc. A dataset stores one of more collections, 
collections store the a bucketted distribution of JSON documents
as well as metadata about the collection (e.g. collection info,
select lists).

COMMANDS

Collection and JSON Documant related--

+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable %s_COLLECTION
    to collection name
+ create - creates a new JSON doc or replace an existing one in collection
  + requires JSON doc name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON doc to stdout
  + requires JSON doc name
+ update - updates a JSON doc in collection
  + requires JSON doc name, followed by replacement JSON blob or 
    JSON blob read from stdin
  + JSON document must already exist
+ delete - removes a JSON doc from collection
  + requires JSON doc name
+ keys - returns the keys to stdout, one key per line
+ path - given a document name return the full path to document

Select list related--

+ select - is the command for working with lists of collection keys
	+ "%s select mylist k1 k2 k3" would create/update a select list 
	  mylist adding keys k1, k2, k3
+ lists - returns the select list names associated with a collection
	+ "%s lists"
+ clear - removes a select list from the collection
	+ "%s clear mylist"
+ first - writes the first key to stdout
	+ "%s first mylist"
+ last would display the last key in the list
	+ "%s last mylist"
+ rest displays all but the first key in the list
	+ "%s rest mylist"
+ list displays a list of keys from the select list to stdout
	+ "dataet list mylist" 
+ shift writes the first key to stdout and remove it from list
	+ "%s shift mylist" 
+ unshift would insert at the beginning 
	+ "%s unshift mylist k4"
+ push would append the list
	+ "%s push mylist k4"
+ pop removes last key form list and displays it
	+ "%s pop mylist" 
+ sort orders the keys alphabetically in the list
	+ "%s sort mylist asc" - sorts in ascending order
	+ "%s sort mylist desc" - sorts in descending order
+ reverse flips the order of the list
	+ "%s reverse mylists"
`

	examples = `
EXAMPLE

This is an example of creating a dataset called testdata/friends, saving
a record called "littlefreda.json" and reading it back.

   %s init testdata/friends
   export DATASET_COLLECTION=testdata/friends
   %s create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(%s keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(%s read $KY)
   done

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of testdata/friends

   %s -i mojosam.json create mojosam
   for KY in $(%s keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(%s read $KY)
   done

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

   cat capt-jack.json | %s create capt-jack
   for KY in $(%s keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(%s read $KY)
   done
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
	inputFName  string

	// App Specific Options
	collectionName string

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":    collectionInit,
		"create":  createJSONDoc,
		"read":    readJSONDoc,
		"update":  updateJSONDoc,
		"delete":  deleteJSONDoc,
		"keys":    collectionKeys,
		"path":    docPath,
		"select":  selectList,
		"lists":   lists,
		"clear":   clear,
		"first":   first,
		"last":    last,
		"rest":    rest,
		"list":    list,
		"push":    push,
		"pop":     pop,
		"shift":   shift,
		"unshift": unshift,
		"length":  length,
		"sort":    sort,
		"reverse": reverse,
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
	collection, err := dataset.Create(name, dataset.GenerateBucketNames(alphabet, 2))
	if err != nil {
		return "", err
	}
	defer collection.Close()
	return fmt.Sprintf("export DATASET_COLLECTION=%q", path.Join(collection.Dataset, collection.Name)), nil
}

// createJSONDoc adds a new JSON document to the collection
func createJSONDoc(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Expected a doc name and JSON blob")
	}
	name, src := args[0], args[1]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name, set DATASET_COLLECTION in the environment variable or use -c option")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if len(src) == 0 {
		return "", fmt.Errorf("Can't create, no JSON source found in %s\n", name)
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
		return "", fmt.Errorf("missing a collection name, set DATASET_COLLECTION in the environment variable or use -c option")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	if len(src) == 0 {
		return "", fmt.Errorf("Can't update, no JSON source found in %s", name)
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

// docPath returns the path to a JSON document or an error
func docPath(args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Missing document name")
	}
	name := args[0]
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	return collection.DocPath(name)
}

func selectList(params ...string) (string, error) {
	if len(params) == 0 {
		params = []string{"keys"}
	}
	if params[0] == "collection" {
		return "", fmt.Errorf("collection is not a valid list name")
	}

	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	l, err := collection.Select(params...)
	if err != nil {
		return "", err
	}
	return strings.Join(l.Keys, "\n"), nil
}

func lists(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	return strings.Join(collection.Lists(), "\n"), nil
}

func clear(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("you can only clear one select list at a time")
	}
	if strings.Compare(params[0], "keys") == 0 {
		return "", fmt.Errorf("select list %s cannot be cleared", params[0])
	}
	if strings.Compare(params[0], "collection") == 0 {
		return "", fmt.Errorf("collection is not a valid select list name")
	}
	err = collection.Clear(params[0])
	if err != nil {
		return "", err
	}
	return "OK", nil

}

func first(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return sl.First(), nil
}

func last(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return sl.Last(), nil
}

func rest(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return strings.Join(sl.Rest(), "\n"), nil
}

func list(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return strings.Join(sl.List(), "\n"), nil
}

func length(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sl.Len()), nil
}

func push(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) < 2 {
		return "", fmt.Errorf("requires list name and one or more keys")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	for _, param := range params[1:] {
		l := sl.Len() + 1
		sl.Push(param)
		if l != sl.Len() {
			return "", fmt.Errorf("%s not added to %s", param, params[0])
		}
	}
	return "OK", nil
}

func pop(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return sl.Pop(), nil
}

func shift(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	return sl.Shift(), nil
}

func unshift(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) < 2 {
		return "", fmt.Errorf("requires list name and one or more keys")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	for _, param := range params[1:] {
		l := sl.Len() + 1
		sl.Unshift(param)
		if l != sl.Len() {
			return "", fmt.Errorf("%s not added to %s", param, params[0])
		}
	}
	return "OK", nil
}

func sort(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) < 2 {
		return "", fmt.Errorf("requires list name and direction (e.g. asc or desc)")
	}
	d := dataset.ASC
	direction := strings.ToLower(strings.TrimSpace(params[1]))
	switch {
	case strings.HasPrefix(direction, "asc"):
		d = dataset.ASC
	case strings.HasPrefix(direction, "desc"):
		d = dataset.DESC
	default:
		d = dataset.ASC
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	sl.Sort(d)
	return "OK", nil
}

func reverse(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) != 1 {
		return "", fmt.Errorf("requires a single list name")
	}
	sl, err := collection.Select(params[0])
	if err != nil {
		return "", err
	}
	sl.Reverse()
	return "OK", nil
}

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.StringVar(&inputFName, "i", "", "input filename")
	flag.StringVar(&inputFName, "input", "", "input filename")

	// Application Options
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
	flag.StringVar(&collectionName, "collection", "", "sets the collection to be used")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description,
		appName, appName, appName, appName, appName,
		appName, appName, appName, appName, appName,
		appName, appName, appName, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples,
		appName, appName, appName, appName,
		appName, appName, appName,
		appName, appName, appName)

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
		// Handle case of piping in or reading JSON from a file.
		if action == "create" && len(params) <= 1 {
			in, err := cli.Open(inputFName, os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer cli.CloseFile(inputFName, in)
			lines, err := cli.ReadLines(in)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			params = append(params, strings.Join(lines, "\n"))
		}

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
