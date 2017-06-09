package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/storage"

	// 3rd Party packages
	"github.com/google/uuid"
)

var (
	usage = `USAGE: %s [OPTIONS] COMMAND_AND_PARAMETERS`

	description = `
SYNOPSIS

dataset is a command line tool demonstrating dataset package for managing 
JSON documents stored on disc. A dataset is organized around collections,
collections contain buckets holding specific JSON documents and related content.
In addition to the JSON documents dataset maintains metadata for management
of the documents, their attachments as well as a ability to generate select lists
based JSON document keys (aka JSON document names).


COMMANDS

Collection and JSON Documant related--

+ init - initialize a new collection if none exists, requires a path to collection
  + once collection is created, set the environment variable DATASET
    to collection name
  + if you're using S3 for storing your dataset prefix your path with 's3://'
    'dataset init s3://mybucket/mydataset-collections'
+ create - creates a new JSON document or replace an existing one in collection
  + requires JSON document name followed by JSON blob or JSON blob read from stdin
+ read - displays a JSON document to stdout
  + requires JSON document name
+ update - updates a JSON document in collection
  + requires JSON document name, followed by replacement JSON document name or 
    JSON document read from stdin
  + JSON document must already exist
+ delete - removes a JSON document from collection
  + requires JSON document name
+ keys - returns the keys to stdout, one key per line
+ path - given a document name return the full path to document
+ attach - attaches a non-JSON content to a JSON record 
    + "dataset attach k1 stats.xlsx" would attach the stats.xlsx file to JSON document named k1
    + (stores content in a related tar file)
+ attachments - lists any attached content for JSON document
    + "dataset attachments k1" would list all the attachments for k1
+ attached - returns attachments for a JSON document 
    + "dataset attached k1" would write out all the attached files for k1
    + "dataset attached k1 stats.xlsx" would write out only the stats.xlsx file attached to k1
+ detach - remove attachments to a JSON document
    + "dataset detach k1 stats.xlsx" would rewrite the attachments tar file without including stats.xlsx
    + "dataset detach k1" would remove ALL attachments to k1
+ import - import a CSV file's rows as JSON documents
	+ "dataset import mydata.csv 1" would import the CSV file mydata.csv using column one's value as key

Select list related--

+ select - is the command for working with lists of collection keys
	+ "dataset select mylist k1 k2 k3" would create/update a select list 
	  mylist adding keys k1, k2, k3
+ lists - returns the select list names associated with a collection
	+ "dataset lists"
+ clear - removes a select list from the collection
	+ "dataset clear mylist"
+ first - writes the first key to stdout
	+ "dataset first mylist"
+ last would display the last key in the list
	+ "dataset last mylist"
+ rest displays all but the first key in the list
	+ "dataset rest mylist"
+ list displays a list of keys from the select list to stdout
	+ "dataet list mylist" 
+ shift writes the first key to stdout and remove it from list
	+ "dataset shift mylist" 
+ unshift would insert at the beginning 
	+ "dataset unshift mylist k4"
+ push would append the list
	+ "dataset push mylist k4"
+ pop removes last key form list and displays it
	+ "dataset pop mylist" 
+ sort orders the keys alphabetically in the list
	+ "dataset sort mylist asc" - sorts in ascending order
	+ "dataset sort mylist desc" - sorts in descending order
+ reverse flips the order of the list
	+ "dataset reverse mylists"
`

	examples = `
EXAMPLES

This is an example of creating a dataset called testdata/friends, saving
a record called "littlefreda.json" and reading it back.

   dataset init testdata/friends
   export DATASET=testdata/friends
   dataset create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of testdata/friends

   dataset -i mojosam.json create mojosam
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

   cat capt-jack.json | dataset create capt-jack
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done

Adding high-capt-jack.txt as an attachment to "capt-jack"

   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset attach capt-jack high-capt-jack.txt

List attachments for "capt-jack"

   dataset attachments capt-jack

Get the attachments for "capt-jack" (this will untar in your current directory)

   dataset attached capt-jack

Remove high-capt-jack.txt from "capt-jack"

    dataset detach capt-jack high-capt-jack.txt

Remove all attachments from "capt-jack"

   dataset detach capt-jack
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
	inputFName  string

	// App Specific Options
	collectionName string
	skipHeaderRow  bool
	useUUID        bool
	showVerbose    bool

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":        collectionInit,
		"create":      createJSONDoc,
		"read":        readJSONDoc,
		"update":      updateJSONDoc,
		"delete":      deleteJSONDoc,
		"keys":        collectionKeys,
		"path":        docPath,
		"select":      selectList,
		"lists":       lists,
		"clear":       clear,
		"first":       first,
		"last":        last,
		"rest":        rest,
		"list":        list,
		"push":        push,
		"pop":         pop,
		"shift":       shift,
		"unshift":     unshift,
		"length":      length,
		"sort":        sort,
		"reverse":     reverse,
		"attach":      addAttachments,
		"attachments": listAttachments,
		"attached":    getAttachments,
		"detach":      removeAttachments,
		"import":      importCSV,
		"check":       checkCollection,
		"repair":      repairCollection,
	}

	// alphabet to use for buckets
	alphabet = `abcdefghijklmnopqrstuvwxyz`
)

//
// These are verbs used in the command line utility
//

// checkCollection takes a collection name and checks for problems
func checkCollection(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	for _, cName := range args {
		if err := dataset.Analyzer(cName); err != nil {
			return "", err
		}
	}
	return "OK", nil
}

// repairCollection takes a collection name and recreates collection.json, keys.json
// based on what it finds on disc
func repairCollection(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	for _, cName := range args {
		if err := dataset.Repair(cName); err != nil {
			return "", err
		}
	}
	return "OK", nil
}

// collectionInit takes a name (e.g. directory path dataset/mycollection) and
// creates a new collection structure on disc
func collectionInit(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	name := args[0]
	collection, err := dataset.Create(name, dataset.GenerateBucketNames(alphabet, 2))
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if collection.Store.Type == storage.S3 {
		return fmt.Sprintf("export DATASET=s3://%s/%s", collection.Store.Config["AwsBucket"], collection.Name), nil
	}
	return fmt.Sprintf("export DATASET=%s", collection.Name), nil
}

// createJSONDoc adds a new JSON document to the collection
func createJSONDoc(args ...string) (string, error) {
	var (
		name string
		src  string
	)

	switch {
	case useUUID == true:
		name = uuid.New().String()
		if len(args) != 1 {
			return "", fmt.Errorf("Expected a JSON blob")
		}
		src = args[0]
	case len(args) == 2:
		name, src = args[0], args[1]
	default:
		return "", fmt.Errorf("Expected a doc name and JSON blob")
	}

	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name, set DATASET in the environment variable or use -c option")
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

	if useUUID == true {
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(src), &m); err != nil {
			return "", err
		}
		if _, ok := m["uuid"]; ok == true {
			m["_uuid"] = name
		} else {
			m["uuid"] = name
		}
		if err := collection.Create(name, m); err != nil {
			return "", err
		}
	} else if err := collection.CreateAsJSON(name, []byte(src)); err != nil {
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
		return "", fmt.Errorf("missing a collection name, set DATASET in the environment variable or use -c option")
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
	if err := sl.SaveList(); err != nil {
		return "", err
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
	r := sl.Pop()
	if err := sl.SaveList(); err != nil {
		return r, err
	}
	return r, nil
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
	r := sl.Shift()
	if err := sl.SaveList(); err != nil {
		return r, err
	}
	return r, nil
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
	if err := sl.SaveList(); err != nil {
		return "", err
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
	if err := sl.SaveList(); err != nil {
		return "", err
	}
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
	if err := sl.SaveList(); err != nil {
		return "", err
	}
	return "OK", nil
}

func addAttachments(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(params) < 2 {
		return "", fmt.Errorf("syntax: %s attach KEY PATH_TO_ATTACHMENT ...", os.Args[0])
	}
	key := params[0]
	err = collection.AttachFiles(key, params[1:]...)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func listAttachments(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) != 1 {
		return "", fmt.Errorf("syntax: %s attachments KEY", os.Args[0])
	}
	key := params[0]
	results, err := collection.Attachments(key)
	if err != nil {
		return "", err
	}
	return strings.Join(results, "\n"), nil
}

func getAttachments(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 1 {
		return "", fmt.Errorf("syntax: %s attached KEY [FILENAMES]", os.Args[0])
	}
	key := params[0]
	err = collection.GetAttachedFiles(key, params[1:]...)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func removeAttachments(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 1 {
		return "", fmt.Errorf("syntax: %s detach KEY", os.Args[0])
	}
	err = collection.Detach(params[0], params[1:]...)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func importCSV(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 1 {
		return "", fmt.Errorf("syntax: %s import CSV_FILENAME [COL_NUMBER_USED_FOR_ID]", os.Args[0])
	}
	idCol := -1
	csvFName := params[0]
	if len(params) > 1 {
		idCol, err = strconv.Atoi(params[1])
		if err != nil {
			return "", fmt.Errorf("Can't convert column number to integer, %s", err)
		}
		// NOTE: we need to adjust to zero based index
		idCol--
	}
	fp, err := os.Open(csvFName)
	if err != nil {
		return "", fmt.Errorf("Can't open %s, %s", csvFName, err)
	}
	defer fp.Close()

	if linesNo, err := collection.ImportCSV(fp, skipHeaderRow, idCol, useUUID, showVerbose); err != nil {
		return "", fmt.Errorf("Can't import CSV, %s", err)
	} else if showVerbose == true {
		log.Printf("%d total rows processed", linesNo)
	}
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
	flag.BoolVar(&skipHeaderRow, "skip-header-row", true, "skip the header row (use as property names)")
	flag.BoolVar(&useUUID, "uuid", false, "generate a UUID for a new JSON document name")
	flag.BoolVar(&showVerbose, "verbose", false, "output rows processed on importing from CSV")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, appName, fmt.Sprintf(dataset.License, appName, dataset.Version), dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = description
	cfg.ExampleText = examples

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
	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}
	action, params := args[0], args[1:]
	if fn, ok := voc[action]; ok == true {
		// Handle case of piping in or reading JSON from a file.
		if (action == "create" || action == "update") && len(params) <= 1 {
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
