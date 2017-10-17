//
// dataset is a command line utility to manage content stored in a dataset collection.
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheets"
	"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"

	// 3rd Party packages
	"github.com/google/uuid"
)

var (
	// Standard Options
	showHelp     bool
	showLicense  bool
	showVersion  bool
	showExamples bool
	inputFName   string
	outputFName  string

	// App Specific Options
	collectionName string
	useHeaderRow   bool
	useUUID        bool
	showVerbose    bool
	quietMode      bool
	noNewLine      bool

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":          collectionInit,
		"create":        createJSONDoc,
		"read":          readJSONDoc,
		"update":        updateJSONDoc,
		"delete":        deleteJSONDoc,
		"join":          joinJSONDoc,
		"keys":          collectionKeys,
		"haskey":        hasKey,
		"filter":        filter,
		"path":          docPath,
		"attach":        addAttachments,
		"attachments":   listAttachments,
		"attached":      getAttachments,
		"detach":        removeAttachments,
		"import":        importCSV,
		"export":        exportCSV,
		"extract":       extract,
		"check":         checkCollection,
		"repair":        repairCollection,
		"import-gsheet": importGSheet,
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
		return fmt.Sprintf("export DATASET=\"s3://%s/%s\"", collection.Store.Config["AwsBucket"], collection.Name), nil
	}
	if collection.Store.Type == storage.GS {
		return fmt.Sprintf("export DATASET=\"gs://%s/%s\"", collection.Store.Config["GoogleBucket"], collection.Name), nil
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

// joinJSONDoc addes/copies fields from another JSON document into the one in the collection.
func joinJSONDoc(args ...string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("either update or overwrite, collection key, one or more JSON document names")
	}
	action := strings.ToLower(args[0])
	key := args[1]
	args = args[2:]

	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	outObject := map[string]interface{}{}
	newObject := map[string]interface{}{}

	if err := collection.Read(key, &outObject); err != nil {
		return "", err
	}
	for _, arg := range args {
		src, err := ioutil.ReadFile(arg)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(src, &newObject); err != nil {
			return "", err
		}
		switch action {
		case "update":
			for k, v := range newObject {
				if _, ok := outObject[k]; ok != true {
					outObject[k] = v
				}
			}
		case "overwrite":
			for k, v := range newObject {
				outObject[k] = v
			}
		default:
			return "", fmt.Errorf("Unknown join type %q", action)
		}
	}
	if err := collection.Update(key, outObject); err != nil {
		return "", err
	}
	return "OK", nil
}

// collectionKeys returns the keys in a collection
func collectionKeys(args ...string) (string, error) {
	// NOTE: We ignore args because this function always returns the full list
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

// hasKey returns true if key is found in collection.json, false otherwise
// If more than one key is provided then each key is checked and an array
// of true/false values will be returned matching the order of the keys provided
// one key state per line
func hasKey(args ...string) (string, error) {
	keyState := []string{}
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	for _, arg := range args {
		keyState = append(keyState, fmt.Sprintf("%t", collection.HasKey(arg)))
	}
	return strings.Join(keyState, "\n"), nil
}

// filter returns a list of collection ids where the filter value returns true.
// the filter notation is based on that Go text/template pipelines that would return
// true in an if/else block.
func filter(args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("filter requires a single filter expression")
	}

	f, err := tmplfn.ParseFilter(args[0])
	if err != nil {
		return "", err
	}

	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	keys := []string{}
	for _, key := range collection.Keys() {
		data := map[string]interface{}{}
		if err := collection.Read(key, &data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
				keys = append(keys, key)
			}
		}
	}
	return strings.Join(keys, "\n"), nil
}

// streamFilterResults works like filter but outputs the results as it find them
func streamFilterResults(w *os.File, filterExp string) error {
	f, err := tmplfn.ParseFilter(filterExp)
	if err != nil {
		return err
	}

	if len(collectionName) == 0 {
		return fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return err
	}
	defer collection.Close()

	for _, key := range collection.Keys() {
		data := map[string]interface{}{}
		if err := collection.Read(key, &data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
				fmt.Fprintln(w, key)
			}
		}
	}
	return nil
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

	if linesNo, err := collection.ImportCSV(fp, useHeaderRow, idCol, useUUID, showVerbose); err != nil {
		return "", fmt.Errorf("Can't import CSV, %s", err)
	} else if showVerbose == true {
		log.Printf("%d total rows processed", linesNo)
	}
	return "OK", nil
}

func importGSheet(params ...string) (string, error) {
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 3 {
		return "", fmt.Errorf("syntax: %s import-gsheet SHEET_ID SHEET_NAME CELL_RANGE [COL_NO_FOR_ID]", os.Args[0])
	}
	spreadSheetId := params[0]
	sheetName := params[1]
	cellRange := params[2]
	idCol := -1
	if len(params) == 4 {
		idCol, err = strconv.Atoi(params[3])
		if err != nil {
			return "", fmt.Errorf("Can't convert column number to integer, %s", err)
		}
		// NOTE: we need to adjust to zero based index
		idCol--
	}

	table, err := gsheets.ReadSheet(clientSecretJSON, spreadSheetId, sheetName, cellRange)
	if err != nil {
		return "", err
	}

	if linesNo, err := collection.ImportTable(table, useHeaderRow, idCol, useUUID, showVerbose); err != nil {
		return "", fmt.Errorf("Can't import Google Sheet, %s", err)
	} else if showVerbose == true {
		log.Printf("%d total rows processed", linesNo)
	}
	return "OK", nil
}

func exportCSV(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 3 {
		return "", fmt.Errorf("syntax: %s export CSV_FILENAME FILTER_EXPR DOTPATHS [COLUMN_NAMES]", os.Args[0])
	}
	csvFName := params[0]
	filterExpr := params[1]
	dotPaths := strings.Split(params[2], ",")
	colNames := []string{}
	if len(params) == 4 {
		colNames = strings.Split(params[3], ",")
	} else {
		for _, val := range dotPaths {
			colNames = append(colNames, val)
		}
	}
	// Trim the any spaces for paths and column names
	for i, val := range dotPaths {
		dotPaths[i] = strings.TrimSpace(val)
	}
	for i, val := range colNames {
		colNames[i] = strings.TrimSpace(val)
	}

	fp, err := os.Create(csvFName)
	if err != nil {
		return "", fmt.Errorf("Can't create %s, %s", csvFName, err)
	}
	defer fp.Close()

	if linesNo, err := collection.ExportCSV(fp, filterExpr, dotPaths, colNames, showVerbose); err != nil {
		return "", fmt.Errorf("Can't export CSV, %s", err)
	} else if showVerbose == true {
		log.Printf("%d total rows processed", linesNo)
	}
	return "OK", nil
}

// extract returns a list of unique values from nested arrays across collection based on
// the filter expression provided.
func extract(params ...string) (string, error) {
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 2 {
		return "", fmt.Errorf("syntax: %s extract FILTER_EXPR DOTPATH", os.Args[0])
	}
	filterExpr := strings.TrimSpace(params[0])
	dotPaths := strings.TrimSpace(params[1])
	lines, err := collection.Extract(filterExpr, dotPaths)
	if err != nil {
		return "", fmt.Errorf("Can't export CSV, %s", err)
	}
	return strings.Join(lines, "\n"), nil
}

func handleError(err error, exitCode int) {
	if quietMode == false {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	if exitCode >= 0 {
		os.Exit(exitCode)
	}
}

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")
	flag.StringVar(&inputFName, "i", "", "input filename")
	flag.StringVar(&inputFName, "input", "", "input filename")
	flag.StringVar(&outputFName, "o", "", "output filename")
	flag.StringVar(&outputFName, "output", "", "output filename")

	// Application Options
	flag.StringVar(&collectionName, "c", "", "sets the collection to be used")
	flag.StringVar(&collectionName, "collection", "", "sets the collection to be used")
	flag.BoolVar(&useHeaderRow, "use-header-row", true, "use the header row as attribute names in the JSON document")
	flag.BoolVar(&useUUID, "uuid", false, "generate a UUID for a new JSON document name")
	flag.BoolVar(&showVerbose, "verbose", false, "output rows processed on importing from CSV")
	flag.BoolVar(&quietMode, "quiet", false, "suppress error and status output")
	flag.BoolVar(&noNewLine, "no-newline", false, "suppress a trailing newline on output")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, strings.ToUpper(appName), dataset.Version)
	cfg.LicenseText = fmt.Sprintf(dataset.License, appName, dataset.Version)
	cfg.UsageText = fmt.Sprintf("%s", Help["usage"])
	cfg.DescriptionText = fmt.Sprintf("%s", Help["description"])
	cfg.OptionText = "## OPTIONS\n\n"
	cfg.ExampleText = fmt.Sprintf("%s", Examples["examples"])

	// Add help and example pages
	for k, v := range Help {
		if k != "nav" {
			cfg.AddHelp(k, fmt.Sprintf("%s", v))
		}
	}
	for k, v := range Examples {
		if k != "nav" {
			cfg.AddExample(k, fmt.Sprintf("%s\n", v))
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
		if len(args) > 0 {
			fmt.Println(cfg.Example(args...))
		} else {
			fmt.Printf("\n%s", cfg.Example())
		}
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

	if len(args) == 0 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}

	in, err := cli.Open(inputFName, os.Stdin)
	if err != nil {
		handleError(err, 1)
	}
	defer cli.CloseFile(inputFName, in)

	out, err := cli.Create(outputFName, os.Stdout)
	if err != nil {
		handleError(err, 1)
	}
	defer cli.CloseFile(outputFName, out)

	action, params := args[0], args[1:]
	if fn, ok := voc[action]; ok == true {
		// If filter we want to output the ids as a stream as they are found
		if action == "filter" {
			var filterExp string
			if len(params) > 0 {
				filterExp = params[0]
			} else {
				buf, err := ioutil.ReadAll(in)
				if err != nil {
					handleError(err, 1)
				}
				filterExp = fmt.Sprintf("%s", buf)
			}
			log.Fatal(streamFilterResults(out, filterExp))
			os.Exit(0)
		}
		// Handle case of piping in or reading JSON from a file.
		if (action == "create" || action == "update") && len(params) <= 1 {
			lines, err := cli.ReadLines(in)
			if err != nil {
				handleError(err, 1)
			}
			params = append(params, strings.Join(lines, "\n"))
		}

		output, err := fn(params...)
		if err != nil {
			handleError(err, 1)
		}
		if quietMode == false || showVerbose == true {
			nl := "\n"
			if noNewLine == true {
				nl = ""
			}
			fmt.Fprintf(out, "%s%s", output, nl)
		}
	} else {
		handleError(fmt.Errorf("Don't understand %q\n", action), 1)
	}
}
