//
// dataset is a command line utility to manage content stored in a dataset collection.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
//
// Copyright (c) 2018, Caltech
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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheet"
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/shuffle"
	//"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"

	// 3rd Party packages
	"github.com/google/uuid"
)

var (
	// Standard Options
	showHelp             bool
	showLicense          bool
	showVersion          bool
	showExamples         bool
	inputFName           string
	outputFName          string
	newLine              bool
	quiet                bool
	prettyPrint          bool
	generateMarkdownDocs bool

	// App Specific Options
	collectionName    string
	useHeaderRow      bool
	useUUID           bool
	showVerbose       bool
	sampleSize        int
	clientSecretFName string
	overwrite         bool
	batchSize         int
	keyFName          string

	// Search specific options, application Options
	showHighlight  bool
	setHighlighter string
	resultFields   string
	sortBy         string
	jsonFormat     bool
	csvFormat      bool
	csvSkipHeader  bool
	idsOnly        bool
	size           int
	from           int
	explain        bool // Note: will be force results to be in JSON format

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":          collectionInit,
		"status":        collectionStatus,
		"create":        createJSONDoc,
		"read":          readJSONDoc,
		"list":          listJSONDoc,
		"update":        updateJSONDoc,
		"delete":        deleteJSONDoc,
		"join":          joinJSONDoc,
		"keys":          collectionKeys,
		"haskey":        hasKey,
		"count":         collectionCount,
		"path":          docPath,
		"attach":        addAttachments,
		"attachments":   listAttachments,
		"detach":        getAttachments,
		"prune":         removeAttachments,
		"import-csv":    importCSV,
		"export-csv":    exportCSV,
		"extract":       extract,
		"check":         checkCollection,
		"repair":        repairCollection,
		"import-gsheet": importGSheet,
		"export-gsheet": exportGSheet,
		"indexer":       indexer,
		"deindexer":     deindexer,
		"find":          find,
	}
)

//
// These are verbs used in the command line utility
//

// checkCollection takes a collection name and checks for problems
func checkCollection(params ...string) (string, error) {
	if len(params) == 0 && collectionName == "" {
		return "", fmt.Errorf("syntax: %s COLLECTION_NAME [COLLECTION_NAME ...]", os.Args[0])
	}
	if collectionName != "" {
		if err := dataset.Analyzer(collectionName); err != nil {
			return "", err
		}
	}
	for _, cName := range params {
		if err := dataset.Analyzer(cName); err != nil {
			return "", err
		}
	}
	return "OK", nil
}

// repairCollection takes a collection name and recreates collection.json, keys.json
// based on what it finds on disc
func repairCollection(params ...string) (string, error) {
	if len(params) == 0 && collectionName == "" {
		return "", fmt.Errorf("syntax: %s COLLECTION_NAME [COLLECTION_NAME ...]", os.Args[0])
	}
	if collectionName != "" {
		if err := dataset.Repair(collectionName); err != nil {
			return "", err
		}
	}
	for _, cName := range params {
		if err := dataset.Repair(cName); err != nil {
			return "", err
		}
	}
	return "OK", nil
}

// collectionInit takes a name (e.g. directory path dataset/mycollection) and
// creates a new collection structure on disc
func collectionInit(params ...string) (string, error) {
	if collectionName == "" && len(params) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if collectionName != "" {
		params = append(params, collectionName)
	}
	for _, cName := range params {
		c, err := dataset.InitCollection(cName)
		if err != nil {
			return "", err
		}
		c.Close()
	}
	return "OK", nil
}

// collectionStatus sees if we can find the dataset collection given the path
func collectionStatus(params ...string) (string, error) {
	if len(params) == 1 && collectionName == "" {
		return "", fmt.Errorf("syntax: %s status COLLECTION_NAME [COLLECTION_NAME ...]", os.Args[0])
	}
	if len(params) == 0 {
		params = []string{collectionName}
	}
	for _, cName := range params {
		c, err := dataset.Open(cName)
		if err != nil {
			return "", fmt.Errorf("%s: %s", cName, err)
		}
		c.Close()
	}
	return "OK", nil
}

// createJSONDoc adds a new JSON document to the collection
func createJSONDoc(params ...string) (string, error) {
	var (
		key       string
		objectSrc string
		src       []byte
		err       error
	)
	if len(params) != 2 {
		return "", fmt.Errorf("Expected a key and a JSON document %q", strings.Join(params, " "))
	}

	key, objectSrc = params[0], params[1]

	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name, set DATASET in the environment variable or use -c option")
	}
	if len(key) == 0 {
		return "", fmt.Errorf("Missing document key")
	}
	if len(objectSrc) == 0 {
		return "", fmt.Errorf("Missing JSON document for %s\n", key)
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if strings.HasSuffix(objectSrc, ".json") {
		src, err = ioutil.ReadFile(objectSrc)
		if err != nil {
			return "", err
		}
	} else {
		src = []byte(objectSrc)
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(src, &m); err != nil {
		return "", fmt.Errorf("%s must be a valid JSON Object", key)
	}
	if useUUID == true {
		m["_UUID"] = key
	}
	if overwrite == true && collection.HasKey(key) == true {
		if err := collection.Update(key, m); err != nil {
			return "", err
		}
		return "OK", nil
	}
	if err := collection.Create(key, m); err != nil {
		return "", err
	}
	return "OK", nil
}

// readJSONDoc returns the JSON from a document in the collection, if more than one key is provided
// it returns an array of JSON docs ordered by the keys provided
func readJSONDoc(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Missing document name")
	}
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	if len(args) == 0 {
		return "", fmt.Errorf("missing document name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if len(args) == 1 {
		m := map[string]interface{}{}
		err := collection.Read(args[0], m)
		if err != nil {
			return "", err
		}
		if prettyPrint {
			src, err := json.MarshalIndent(m, "", "    ")
			if err != nil {
				return "", err
			}
			return string(src), nil
		}
		src, err := json.Marshal(m)
		return string(src), err
	}

	recs := []map[string]interface{}{}
	for _, name := range args {
		m := map[string]interface{}{}
		err := collection.Read(name, m)
		if err != nil {
			return "", err
		}
		recs = append(recs, m)
	}
	if prettyPrint {
		src, err := json.MarshalIndent(recs, "", "    ")
		return string(src), err
	}
	src, err := json.Marshal(recs)
	return string(src), err
}

// listJSONDoc returns a JSON array from a document in the collection
// if not matching records returns an empty list
func listJSONDoc(args ...string) (string, error) {
	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name")
	}
	if len(args) == 0 && len(keyFName) == 0 {
		return "[]", nil
	}
	var keyList []string
	if len(keyFName) > 0 {
		src, err := ioutil.ReadFile(keyFName)
		if err != nil {
			return "", fmt.Errorf("Cannot read key file %s, %s", keyFName, err)
		}
		txt := fmt.Sprintf("%s", src)
		for _, key := range strings.Split(txt, "\n") {
			key = strings.TrimSpace(key)
			if len(key) > 0 {
				keyList = append(keyList, key)
			}
		}
	}
	if len(args) > 0 {
		keyList = append(keyList, args...)
	}

	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	recs := []map[string]interface{}{}
	for _, name := range keyList {
		m := map[string]interface{}{}
		err := collection.Read(name, m)
		if err != nil {
			return "", err
		}
		recs = append(recs, m)
	}
	if prettyPrint {
		src, err := json.MarshalIndent(recs, "", "    ")
		return string(src), err
	}
	src, err := json.Marshal(recs)
	return string(src), err
}

// updateJSONDoc replaces a JSON document in the collection
func updateJSONDoc(params ...string) (string, error) {
	var (
		key       string
		objectSrc string
		src       []byte
		err       error
	)
	if len(params) != 2 {
		return "", fmt.Errorf("Expected document name and JSON blob")
	}
	key, objectSrc = params[0], params[1]

	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name, set DATASET in the environment variable or use -c option")
	}
	if len(key) == 0 {
		return "", fmt.Errorf("Missing document key")
	}
	if len(objectSrc) == 0 {
		return "", fmt.Errorf("Can't update, no JSON source found for %s", key)
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	if strings.HasSuffix(objectSrc, ".json") {
		src, err = ioutil.ReadFile(objectSrc)
		if err != nil {
			return "", err
		}
	} else {
		src = []byte(objectSrc)
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(src, &data); err != nil {
		return "", err
	}
	if err := collection.Update(key, data); err != nil {
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
func joinJSONDoc(params ...string) (string, error) {
	var (
		src        []byte
		err        error
		adverb     string
		key        string
		objectSrcs []string
	)
	if len(params) < 3 {
		return "", fmt.Errorf("expected append or overwrite, collection key, one or more JSON Objects, got %s", strings.Join(params, ", "))
	}
	adverb, key, objectSrcs = strings.ToLower(params[0]), params[1], params[2:]

	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	outObject := map[string]interface{}{}
	newObject := map[string]interface{}{}

	if err := collection.Read(key, outObject); err != nil {
		return "", err
	}

	for _, objectSrc := range objectSrcs {
		if strings.HasSuffix(objectSrc, ".json") {
			src, err = ioutil.ReadFile(objectSrc)
			if err != nil {
				return "", err
			}
		} else {
			src = []byte(objectSrc)
		}
		if err := json.Unmarshal(src, &newObject); err != nil {
			return "", err
		}
		switch adverb {
		case "append":
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
			return "", fmt.Errorf("Unknown join type %q", adverb)
		}
	}
	if err := collection.Update(key, outObject); err != nil {
		return "", err
	}
	return "OK", nil
}

// collectionKeys returns the keys in a collection
// If a 'filter expression' is provided it will return a filtered list of keys.
// Filters with like Go's text/template if statement where the 'filter expression' is
// the condititional expression in a if/else statement. If the expression evaluates to "true"
// then the kehy is included in the list of keys If the expression evaluates to "false" then
// it is excluded for the list of keys.
// If a 'sort expression' is provided then the resulting keys are ordered by that expression.
func collectionKeys(args ...string) (string, error) {
	var (
		keyList  []string
		sortExpr string
	)
	// NOTE: We ignore args because this function always returns the full list
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	// Set our filter
	filterExpr := "true"
	if len(args) > 0 {
		filterExpr = args[0]
	}

	// Some sort of Sort is involved
	if len(args) > 1 {
		sortExpr = args[1]
	}

	// Some sort of sub selection of keys is involved
	if len(args) > 2 {
		keyList = args[2:]
	} else {
		keyList = collection.Keys()
	}

	// Save the resulting keys in a separate list
	keys := []string{}

	// Process the filter
	keys, err = collection.KeyFilter(keyList, filterExpr)
	if err != nil {
		return "", err
	}

	// Apply Sample Size
	if sampleSize > 0 {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			keys = keys[0:sampleSize]
		}
	}

	// If now sort we're done
	if len(sortExpr) == 0 {
		return strings.Join(keys, "\n"), nil
	}

	// We still have sorting to do.
	keys, err = collection.KeySortByExpression(keys, args[1])
	return strings.Join(keys, "\n"), err
}

// hasKey returns true if keys are found in collection.json, false otherwise
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

// collectionCount returns the number of keys in a collection
// can optionally accept a filter to return a subset count of keys
func collectionCount(args ...string) (string, error) {
	var keyList []string

	// NOTE: We ignore args because this function always returns a count
	if len(collectionName) == 0 {
		return "", fmt.Errorf("missing a collection name")
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	// Trivial case where we want length of whole collection
	if len(args) == 0 || (len(args) == 1 && args[0] == "true") {
		return fmt.Sprintf("%d", collection.Length()), nil
	}

	// Some sort of filter is involved.
	f, err := tmplfn.ParseFilter(args[0])
	if err != nil {
		return "", err
	}
	if len(args) > 1 {
		keyList = args[1:]
	} else {
		keyList = collection.Keys()
	}
	cnt := 0
	for _, key := range keyList {
		m := map[string]interface{}{}
		if err := collection.Read(key, m); err == nil {
			if ok, err := f.Apply(m); err == nil && ok == true {
				cnt++
			}
		}
	}
	return fmt.Sprintf("%d", cnt), nil
}

// streamFilterResults works like filter but outputs the results as it find them
func streamFilterResults(w *os.File, keyList []string, filterExp string, sampleSize int) error {
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

	if len(keyList) == 0 {
		keyList = collection.Keys()
		if sampleSize > 0 {
			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			shuffle.Strings(keyList, random)
			keyList = keyList[0:sampleSize]
		}
	}
	for _, key := range keyList {
		m := map[string]interface{}{}
		if err := collection.Read(key, m); err == nil {
			if ok, err := f.Apply(m); err == nil && ok == true {
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
	if collection.HasKey(key) == false {
		return "", fmt.Errorf("%q is not in collection", key)
	}
	for _, fname := range params[1:] {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			return "", fmt.Errorf("%s does not exist", fname)
		}
	}
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
	if collection.HasKey(key) == false {
		return "", fmt.Errorf("%q is not in collection", key)
	}
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
		return "", fmt.Errorf("syntax: %s detach KEY [FILENAMES]", os.Args[0])
	}
	key := params[0]
	if collection.HasKey(key) == false {
		return "", fmt.Errorf("%q is not in collection", key)
	}
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
		return "", fmt.Errorf("syntax: %s prune KEY", os.Args[0])
	}
	err = collection.Prune(params[0], params[1:]...)
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
	if len(params) < 2 {
		return "", fmt.Errorf("syntax: %s import-csv CSV_FILENAME COL_NUMBER_USED_FOR_ID", os.Args[0])
	}
	csvFName := params[0]
	idCol, err := strconv.Atoi(params[1])
	if err != nil {
		return "", fmt.Errorf("Can't convert column number to integer, %s", err)
	}
	if idCol < 1 {
		return "", fmt.Errorf("Column number must be greater than zero, got %s", idCol)
	}

	// NOTE: we need to adjust to zero based index
	idCol--
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
	if clientSecretFName != "" {
		clientSecretJSON = clientSecretFName
	}
	if clientSecretJSON == "" {
		clientSecretJSON = "client_secret.json"
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 4 {
		return "", fmt.Errorf("syntax: %s import-gsheet SHEET_ID SHEET_NAME CELL_RANGE COL_NUMBER_USED_FOR_ID", os.Args[0])
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

	if linesNo, err := collection.ImportTable(table, useHeaderRow, idCol, useUUID, overwrite, showVerbose); err != nil {
		return "", fmt.Errorf("Errors importing %s, %s", sheetName, err)
	} else if showVerbose == true {
		log.Printf("%d total rows processed", linesNo)
	}
	return "OK", nil
}

func exportGSheet(params ...string) (string, error) {
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretFName != "" {
		clientSecretJSON = clientSecretFName
	}
	if clientSecretJSON == "" {
		clientSecretJSON = "client_secret.json"
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	if len(params) < 5 {
		return "", fmt.Errorf("syntax: %s export-gsheet SHEET_ID SHEET_NAME CELL_RANGE FILTER_EXPR EXPORT_FIELD_LIST [COLUMN_NAMES]", os.Args[0])
	}
	spreadSheetId := params[0]
	sheetName := params[1]
	cellRange := params[2]
	filterExpr := params[3]
	dotExprs := strings.Split(params[4], ",")
	colNames := []string{}
	if len(params) <= 5 {
		for _, val := range dotExprs {
			colNames = append(colNames, val)
		}
	} else {
		colNames = strings.Split(params[5], ",")
	}
	// Trim the any spaces for paths and column names
	for i, val := range dotExprs {
		dotExprs[i] = strings.TrimSpace(val)
	}
	for i, val := range colNames {
		colNames[i] = strings.TrimSpace(val)
	}

	table := [][]interface{}{}
	if len(colNames) > 0 {
		row := []interface{}{}
		for _, name := range colNames {
			row = append(row, name)
		}
		table = append(table, row)
	}
	keys := collection.Keys()

	if strings.ToLower(filterExpr) == "true" {
		for _, key := range keys {
			m := map[string]interface{}{}
			if err := collection.Read(key, m); err == nil {
				row := []interface{}{}
				for _, colPath := range dotExprs {
					col, err := dotpath.Eval(colPath, m)
					if err == nil {
						row = append(row, col)
					} else {
						row = append(row, "")
						return "", fmt.Errorf("failed to evaluate dot path, %s", err)
					}
				}
				table = append(table, row)
			}
		}
	} else {
		f, err := tmplfn.ParseFilter(filterExpr)
		if err != nil {
			return "", err
		}

		for _, key := range keys {
			m := map[string]interface{}{}
			if err := collection.Read(key, m); err == nil {
				if ok, err := f.Apply(m); err == nil && ok == true {
					// save row out.
					row := []interface{}{}
					for _, colPath := range dotExprs {
						col, err := dotpath.Eval(colPath, m)
						if err == nil {
							row = append(row, col)
						} else {
							row = append(row, "")
							return "", fmt.Errorf("failed to evaluate dot path, %s", err)
						}
					}
					table = append(table, row)
				}
			}
		}
	}
	if err := gsheets.WriteSheet(clientSecretJSON, spreadSheetId, sheetName, cellRange, table); err != nil {
		return "", err
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
		return "", fmt.Errorf("syntax: %s export-csv CSV_FILENAME FILTER_EXPR DOTPATHS [COLUMN_NAMES]", os.Args[0])
	}
	csvFName := params[0]
	filterExpr := params[1]
	dotExprs := strings.Split(params[2], ",")
	colNames := []string{}
	if len(params) < 4 {
		for _, val := range dotExprs {
			colNames = append(colNames, strings.TrimPrefix(val, "."))
		}
	} else {
		colNames = strings.Split(params[3], ",")
	}
	// Trim the any spaces for paths and column names
	for i, val := range dotExprs {
		dotExprs[i] = strings.TrimSpace(val)
	}
	for i, val := range colNames {
		colNames[i] = strings.TrimSpace(val)
	}

	fp, err := os.Create(csvFName)
	if err != nil {
		return "", fmt.Errorf("Can't create %s, %s", csvFName, err)
	}
	defer fp.Close()

	if linesNo, err := collection.ExportCSV(fp, os.Stderr, filterExpr, dotExprs, colNames, showVerbose); err != nil {
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
		return "", fmt.Errorf("syntax: %s extract FILTER_EXPR DOTPATH_EXPR", os.Args[0])
	}
	filterExpr := strings.TrimSpace(params[0])
	dotExpr := strings.TrimSpace(params[1])
	lines, err := collection.Extract(filterExpr, dotExpr)
	return strings.Join(lines, "\n"), err
}

// indexer replaces dsindexer command and is used to build a Bleve index for a collection
func indexer(params ...string) (string, error) {
	var (
		indexName    string
		indexMapName string
		keyList      []string
	)
	if len(params) < 2 {
		return "", fmt.Errorf("syntax: %s [OPTIONS] indexer INDEX_MAP_FILENAME INDEX_NAME", os.Args[0])
	}
	if len(params) > 0 {
		if strings.HasSuffix(params[0], "bleve") {
			indexName = params[0]
		} else {
			indexMapName = params[0]
		}
	}
	if len(params) > 1 {
		if strings.HasSuffix(params[1], "bleve") {
			indexName = params[1]
		} else {
			indexMapName = params[1]
		}
	}

	c, err := dataset.Open(collectionName)
	if err != nil {
		return "", fmt.Errorf("Cannot open collection %s, %s", collectionName, err)
	}
	defer c.Close()

	if len(keyFName) > 0 {
		src, err := ioutil.ReadFile(keyFName)
		if err != nil {
			return "", fmt.Errorf("Cannot read key file %s, %s", keyFName, err)
		}
		txt := fmt.Sprintf("%s", src)
		for _, key := range strings.Split(txt, "\n") {
			key = strings.TrimSpace(key)
			if len(key) > 0 {
				keyList = append(keyList, key)
			}
		}
	} else {
		keyList = c.Keys()
	}

	if batchSize == 0 {
		if len(keyList) > 100000 {
			batchSize = 1000
		} else if len(keyList) > 10000 {
			batchSize = len(keyList) / 100
		} else if len(keyList) > 1000 {
			batchSize = len(keyList) / 10
		} else {
			batchSize = 100
		}
	}

	err = c.Indexer(indexName, indexMapName, keyList, batchSize)
	if err != nil {
		return "", fmt.Errorf("Indexing error %s %s, %s", collectionName, indexName, err)
	}
	// return success
	return "OK", nil
}

// deindexer replaces dsindexer command and is used to build a Bleve index for a collection
func deindexer(params ...string) (string, error) {
	var (
		indexName string
		keyFName  string
		keyList   []string
	)
	if len(params) == 0 {
		return "", fmt.Errorf("syntax: %s deindexer INDEX_NAME KEY_FILENAME", os.Args[0])
	}
	if len(params) > 0 {
		if strings.HasSuffix(params[0], ".bleve") {
			indexName = params[0]
		} else {
			keyFName = params[0]
		}
	}
	if len(params) > 1 {
		if strings.HasSuffix(params[1], ".bleve") {
			indexName = params[1]
		} else {
			keyFName = params[1]
		}
	}

	if len(keyFName) > 0 {
		src, err := ioutil.ReadFile(keyFName)
		if err != nil {
			return "", fmt.Errorf("Cannot read key file %s, %s", keyFName, err)
		}
		txt := fmt.Sprintf("%s", src)
		for _, key := range strings.Split(txt, "\n") {
			key = strings.TrimSpace(key)
			if len(key) > 0 {
				keyList = append(keyList, key)
			}
		}
	}
	if len(keyList) == 0 {
		return "", fmt.Errorf("Deindexing requires a list of keys to de-index")
	}

	if batchSize == 0 {
		if len(keyList) > 100000 {
			batchSize = 1000
		} else if len(keyList) > 10000 {
			batchSize = len(keyList) / 100
		} else if len(keyList) > 1000 {
			batchSize = len(keyList) / 10
		} else {
			batchSize = 100
		}
	}
	if err := dataset.Deindexer(indexName, keyList, batchSize); err != nil {
		return "", fmt.Errorf("Deindexing error %s %s, %s", collectionName, indexName, err)
	}
	// return success
	return "OK", nil
}

func find(params ...string) (string, error) {
	if len(params) < 2 {
		return "", fmt.Errorf("syntax: %s [OPTIONS] INDEX_NAMES QUERY_STRING", os.Args[0])
	}
	indexNames := []string{}
	queryString := ""
	for _, param := range params {
		if len(param) > 0 {
			if strings.HasSuffix(param, ".bleve") {
				if strings.Contains(param, ":") == true {
					indexNames = append(indexNames, strings.Split(param, ":")...)
				} else {
					indexNames = append(indexNames, param)
				}
			} else {
				queryString = param
			}
		}
	}
	options := map[string]string{}
	if explain == true {
		options["explain"] = "true"
		jsonFormat = true
	}

	if sampleSize > 0 {
		options["sample"] = fmt.Sprintf("%d", sampleSize)
	}
	if from != 0 {
		options["from"] = fmt.Sprintf("%d", from)
	}
	if batchSize > 0 {
		options["size"] = fmt.Sprintf("%d", batchSize)
	}
	if sortBy != "" {
		options["sort"] = sortBy
	}
	if showHighlight == true {
		options["highlight"] = "true"
		options["highlighter"] = "ansi"
	}
	if setHighlighter != "" {
		options["highlight"] = "true"
		options["highlighter"] = setHighlighter
	}

	if resultFields != "" {
		options["fields"] = strings.TrimSpace(resultFields)
	} else {
		options["fields"] = "*"
	}

	idxList, idxFields, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		return "", fmt.Errorf("Can't open index %s, %s", strings.Join(indexNames, ", "), err)
	}

	results, err := dataset.Find(idxList.Alias, queryString, options)
	if err != nil {
		return "", fmt.Errorf("Find error %s, %s", strings.Join(indexNames, ", "), err)
	}
	err = idxList.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't close indexes %s, %s", strings.Join(indexNames, ", "), err)
	}

	//
	// Handle results formatting choices
	//
	switch {
	case jsonFormat == true:
		if prettyPrint {
			src, err := json.MarshalIndent(results, "", "    ")
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s", src), err
		}
		src, err := json.Marshal(results)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s", src), err
	case csvFormat == true:
		var fields []string
		if resultFields == "" {
			fields = idxFields
		} else {
			fields = strings.Split(resultFields, ",")
		}
		var buf bytes.Buffer
		fp := bufio.NewWriter(&buf)
		err = dataset.CSVFormatter(fp, results, fields, csvSkipHeader)
		if err != nil {
			return "", err
		}
		if err := fp.Flush(); err != nil {
			return "", err
		}
		return buf.String(), nil
	case idsOnly == true:
		ids := []string{}
		for _, hit := range results.Hits {
			ids = append(ids, hit.ID)
		}
		return strings.Join(ids, "\n"), nil
	}
	return results.String(), nil
}

func main() {
	app := cli.NewCli(dataset.Version)
	appName := app.AppName()
	// We require an "ACTION" or verb for command to work.
	app.ActionsRequired = true

	// Add command line parameters.
	app.AddParams("COLLECTION_NAME")

	// Add Help Docs
	for k, v := range Help {
		app.AddHelp(k, v)
	}
	for k, v := range Examples {
		app.AddHelp(k, v)
	}

	// Add Environment options
	app.EnvStringVar(&collectionName, "DATASET", "", "Set the working path to your dataset collection")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "e,examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&newLine, "nl,newline", true, "if set to false suppress the trailing newline")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// Application Options
	app.StringVar(&collectionName, "c,collection", "", "sets the collection to be used")
	app.BoolVar(&useHeaderRow, "use-header-row", true, "(import) use the header row as attribute names in the JSON document")
	app.BoolVar(&useUUID, "uuid", false, "(import) generate a UUID for a new JSON document name")
	app.BoolVar(&showVerbose, "verbose", false, "output rows processed on importing from CSV")
	app.IntVar(&sampleSize, "sample", 0, "set the sample size when listing keys")
	app.StringVar(&clientSecretFName, "client-secret", "", "(import-gsheet, export-gsheet) set the client secret path and filename for GSheet access")
	app.BoolVar(&overwrite, "overwrite", false, "overwrite will treat a create as update if the record exists")
	app.IntVar(&batchSize, "batch,size", 0, "(indexer, deindexer, find) set the number of records per response")
	app.StringVar(&keyFName, "key-file", "", "operate on the record keys contained in file, one key per line")

	// Search specific application options
	app.StringVar(&sortBy, "sort", "", "(find) a comma delimited list of field names to sort by")
	app.BoolVar(&showHighlight, "highlight", false, "(find) display highlight in search results")
	app.StringVar(&setHighlighter, "highlighter", "", "(find) set the highlighter (ansi,html) for search results")
	app.StringVar(&resultFields, "fields", "", "(find) comma delimited list of fields to display in the results")
	app.BoolVar(&jsonFormat, "json", false, "(find) format results as a JSON document")
	app.BoolVar(&csvFormat, "csv", false, "(find) format results as a CSV document, used with fields option")
	app.BoolVar(&csvSkipHeader, "csv-skip-header", false, "(find) don't output a header row, only values for csv output")
	app.BoolVar(&idsOnly, "ids,ids-only", false, "(find) output only a list of ids from results")
	app.IntVar(&from, "from", 0, "(find) return the result starting with this result number")
	app.BoolVar(&explain, "explain", false, "(find) explain results in a verbose JSON document")

	// Action verbs (e.g. app.AddAction(STRING_VERB, FUNC_POINTER, STRING_DESCRIPTION)
	// NOTE: Sense dataset cli was developed pre-existed cli v0.0.6 we're only document our actions and not run them via cli.
	app.AddVerb("init", "Initialize a dataset collection")
	app.AddVerb("status", "Checks to see if a collection name contains a 'collection.json' file")
	app.AddVerb("create", "Create a JSON record in a collection")
	app.AddVerb("read", "Read back a JSON record from a collection")
	app.AddVerb("list", "List the JSON records as an array for provided record ids")
	app.AddVerb("update", "Update a JSON record in a collection")
	app.AddVerb("delete", "Delete a JSON record (and attachments) from a collection")
	app.AddVerb("join", "Join a JSON record with a new JSON object in a collection")
	app.AddVerb("keys", "List the keys in a collection, support filtering and sorting")
	app.AddVerb("haskey", "Returns true if key is in collection, false otherwise")
	app.AddVerb("count", "Counts the number of records in a collection, accepts a filter for sub-counts")
	app.AddVerb("path", "Show the file system path to a JSON record in a collection")
	app.AddVerb("attach", "Attach a document (file) to a JSON record in a collection")
	app.AddVerb("attachments", "List of attachments associated with a JSON record in a collection")
	app.AddVerb("detach", "Copy an attach out of an associated JSON record in a collection")
	app.AddVerb("prune", "Remove attachments from a JSON record in a collection")
	app.AddVerb("import", "Import a CSV file's rows as JSON records into a collection")
	app.AddVerb("export", "Export a JSON records from a collection to a CSV file")
	app.AddVerb("extract", "Extract unique values from JSON records in a collection based on a dot path expression")
	app.AddVerb("check", "Check the health of a dataset collection")
	app.AddVerb("repair", "Try to repair a damaged dataset collection")
	app.AddVerb("import-gsheet", "Import a GSheet rows as JSON records into a collection")
	app.AddVerb("export-gsheet", "Export a collection's JSON records to a GSheet")
	app.AddVerb("indexer", "Create/Update a Bleve index of a collection")
	app.AddVerb("deindexer", "Remove record(s) from a Bleve index for a collection")
	app.AddVerb("find", "Query a bleve index(es) associated with a collection")

	// We're ready to process args
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
		os.Exit(0)
	}
	if showHelp {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showExamples {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			keys := []string{}
			for k, _ := range Examples {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			fmt.Fprintf(app.Out, "try \"%s -examples TOPIC\" for any of these topics:\n\t%s\n\n", appName, strings.Join(keys, "\n\t"))
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	// Application Option's processing
	if len(args) == 0 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("See %s --help for usage", appName), quiet)
	}

	// Check for collectionName in environment
	// Trival check, look for *.ds, s3://, gs:// in the args and use that for collection name if present.
	for i, arg := range args {
		if strings.HasSuffix(arg, ".ds") || strings.HasSuffix(arg, ".dataset") || strings.HasPrefix(arg, "gs://") || strings.HasPrefix(arg, "s3://") {
			collectionName = arg
			if i < len(args) {
				if i == 0 {
					args = args[1:]
				} else {
					args = append(args[:i], args[i+1:]...)
				}
			}
		}
	}

	// Merge environment if colleciton name not set
	if collectionName == "" {
		datasetEnv := os.Getenv("DATASET")
		if datasetEnv != "" {
			collectionName = strings.TrimSpace(datasetEnv)
		}
	}

	in, err := cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(inputFName, in)

	out, err := cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(outputFName, out)

	var (
		action string
		params []string
	)
	if len(args) > 1 {
		action, params = args[0], args[1:]
	} else if len(args) == 1 {
		action = args[0]
		params = []string{}
	}

	// NOTE: Special case of when -useUUID flag set when action is create, import or
	// import-gsheet, we need to auto-generate the UUID as key and add to our args
	// appropriately
	if useUUID {
		id := uuid.New().String()
		switch action {
		case "create":
			params = append([]string{id}, args[1:]...)
		case "import":
			params = append(args, id)
		case "import-gsheet":
		}
	}

	var data string
	if inputFName != "" {
		src, err := ioutil.ReadAll(in)
		cli.ExitOnError(os.Stderr, err, quiet)
		data = string(src)
	}

	fn, ok := voc[action]
	if ok == false {
		cli.ExitOnError(os.Stderr, fmt.Errorf("do not understand %s for %q", action, strings.Join(os.Args, " ")), quiet)
	}

	if (action == "create" || action == "update" || action == "join") && len(data) > 0 {
		params = append(params, data)
	}

	if (action == "read" || action == "list") && len(data) > 0 {
		// Split the input if available
		lines := strings.Split(data, "\n")
		for _, key := range lines {
			params = append(params, key)
		}
	}

	if (action == "keys" || action == "count") && len(data) > 0 {
		// Split the input if available
		lines := strings.Split(data, "\n")
		// If filter we want to output the ids as a stream as they are found
		filterExpr := "true"
		sortExpr := ""
		keyList := []string{}

		if len(params) > 0 {
			filterExpr = params[0]
		}
		if len(params) > 1 {
			sortExpr = params[1]
		}

		// Get any key list that might be passed in (either via cli or stdin)
		if len(params) > 3 {
			keyList = params[2:]
		} else if len(lines) > 0 {
			for _, line := range lines {
				keyList = append(keyList, line)
			}
		}

		// If we are NOT sorting we can just filter the output now and be done.
		if len(sortExpr) == 0 {
			if err := streamFilterResults(app.Out, keyList, filterExpr, sampleSize); err != nil {
				fmt.Fprintf(app.Out, "%s\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		// We need to make sure our params are setup with sane defaults for keys and count
		params = []string{
			filterExpr,
			sortExpr,
		}
		for _, k := range keyList {
			params = append(params, k)
		}
	}

	output, err := fn(params...)
	cli.ExitOnError(os.Stderr, err, quiet)
	if quiet == false || showVerbose == true && output != "" {
		fmt.Fprintf(app.Out, "%s", output)
	}
	if newLine {
		fmt.Fprintln(app.Out, "")
	}
}
