//
// dataset is a command line utility to manage content stored in a dataset collection.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheets"
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/shuffle"
	"github.com/caltechlibrary/storage"
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
	collectionName string
	useHeaderRow   bool
	useUUID        bool
	showVerbose    bool
	sampleSize     int

	// Vocabulary
	voc = map[string]func(...string) (string, error){
		"init":          collectionInit,
		"status":        collectionStatus,
		"create":        createJSONDoc,
		"read":          readJSONDocs,
		"list":          listJSONDocs,
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
		"import":        importCSV,
		"export":        exportCSV,
		"extract":       extract,
		"check":         checkCollection,
		"repair":        repairCollection,
		"import-gsheet": importGSheet,
		"export-gsheet": exportGSheet,
	}
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
	collection, err := dataset.InitCollection(name)
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

// collectionStatus sees if we can find the dataset collection given the path
func collectionStatus(args ...string) (string, error) {
	if len(args) == 0 && collectionName == "" {
		return "", fmt.Errorf("missing a collection name")
	}
	args = append(args, collectionName)
	for _, collectionName := range args {
		_, err := dataset.Open(collectionName)
		if err != nil {
			return "", fmt.Errorf("%s: %s", collectionName, err)
		}
	}
	return "OK", nil
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
		return "", fmt.Errorf("Expected a document name and a JSON document")
	}

	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name, set DATASET in the environment variable or use -c option")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("Missing document name")
	}
	if len(src) == 0 {
		return "", fmt.Errorf("Missing JSON document %s\n", name)
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(src), &m); err != nil {
		return "", fmt.Errorf("%s must be a valid JSON Object", name)
	}
	if useUUID == true {
		m["_uuid"] = name
	}
	if err := collection.Create(name, m); err != nil {
		return "", err
	}
	return "OK", nil
}

// readJSONDocs returns the JSON from a document in the collection, if more than one key is provided
// it returns an array of JSON docs ordered by the keys provided
func readJSONDocs(args ...string) (string, error) {
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
		data := map[string]interface{}{}
		err := collection.Read(args[0], data)
		if err != nil {
			return "", err
		}
		if prettyPrint {
			src, err := json.MarshalIndent(data, "", "    ")
			if err != nil {
				return "", err
			}
			return string(src), nil
		}
		src, err := json.Marshal(data)
		return string(src), err
	}

	var rec map[string]interface{}
	recs := []map[string]interface{}{}
	for _, name := range args {
		err := collection.Read(name, rec)
		if err != nil {
			return "", err
		}
		recs = append(recs, rec)
	}
	if prettyPrint {
		src, err := json.MarshalIndent(recs, "", "    ")
		return string(src), err
	}
	src, err := json.Marshal(recs)
	return string(src), err
}

// listJSONDocs returns a JSON array from a document in the collection
// if not matching records returns an empty list
func listJSONDocs(args ...string) (string, error) {
	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name")
	}
	if len(args) == 0 {
		return "[]", nil
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()

	var rec map[string]interface{}
	recs := []map[string]interface{}{}
	for _, name := range args {
		err := collection.Read(name, rec)
		if err != nil {
			return "", err
		}
		recs = append(recs, rec)
	}
	if prettyPrint {
		src, err := json.MarshalIndent(recs, "", "    ")
		return string(src), err
	}
	src, err := json.Marshal(recs)
	return string(src), err
}

// updateJSONDoc replaces a JSON document in the collection
func updateJSONDoc(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Expected document name and JSON blob")
	}
	name, src := args[0], []byte(args[1])
	if len(collectionName) == 0 {
		return "", fmt.Errorf("Missing a collection name, set DATASET in the environment variable or use -c option")
	}
	if len(name) == 0 {
		return "", fmt.Errorf("Missing document name")
	}
	if len(src) == 0 {
		return "", fmt.Errorf("Can't update, no JSON source found in %s", name)
	}
	collection, err := dataset.Open(collectionName)
	if err != nil {
		return "", err
	}
	defer collection.Close()
	data := map[string]interface{}{}
	if err := json.Unmarshal(src, &data); err != nil {
		return "", err
	}
	if err := collection.Update(name, data); err != nil {
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
		return "", fmt.Errorf("expected update or overwrite, collection key, one or more JSON Objects, got %s", strings.Join(args, ", "))
	}
	action := strings.ToLower(args[0])
	key := args[1]
	objects_src := args[2:]

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
	for _, src := range objects_src {
		if err := json.Unmarshal([]byte(src), &newObject); err != nil {
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

	// Trivial case of return all keys
	if len(args) == 0 || (len(args) == 1 && args[0] == "true") {
		if sampleSize == 0 {
			return strings.Join(collection.Keys(), "\n"), nil
		}
		keys := collection.Keys()
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		shuffle.Strings(keys, random)
		if sampleSize <= len(keys) {
			return strings.Join(keys[0:sampleSize], "\n"), nil
		}
		return strings.Join(collection.Keys(), "\n"), nil
	}

	// Some sort of filter is involved
	f, err := tmplfn.ParseFilter(args[0])
	if err != nil {
		return "", err
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
	for _, key := range keyList {
		data := map[string]interface{}{}
		if err := collection.Read(key, data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
				keys = append(keys, key)
			}
		}
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
	keys, err = collection.SortKeysByExpression(keys, args[1])
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
		data := map[string]interface{}{}
		if err := collection.Read(key, data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
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
		data := map[string]interface{}{}
		if err := collection.Read(key, data); err == nil {
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
	if collection.HasKey(key) == false {
		return "", fmt.Errorf("%q is not in collection", key)
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

func exportGSheet(params ...string) (string, error) {
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
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
	dotPaths := strings.Split(params[4], ",")
	colNames := []string{}
	if len(params) < 5 {
		for _, val := range dotPaths {
			colNames = append(colNames, val)
		}
	} else {
		colNames = strings.Split(params[5], ",")
	}
	// Trim the any spaces for paths and column names
	for i, val := range dotPaths {
		dotPaths[i] = strings.TrimSpace(val)
	}
	for i, val := range colNames {
		colNames[i] = strings.TrimSpace(val)
	}

	keys := collection.Keys()
	f, err := tmplfn.ParseFilter(filterExpr)
	if err != nil {
		return "", err
	}

	var (
		data map[string]interface{}
	)

	table := [][]interface{}{}
	if len(colNames) > 0 {
		row := []interface{}{}
		for _, name := range colNames {
			row = append(row, name)
		}
		table = append(table, row)
	}
	for _, key := range keys {
		if err := collection.Read(key, data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
				// save row out.
				row := []interface{}{}
				for _, colPath := range dotPaths {
					col, err := dotpath.Eval(colPath, data)
					if err == nil {
						row = append(row, col)
					} else {
						row = append(row, "")
					}
				}
				table = append(table, row)
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
		return "", fmt.Errorf("syntax: %s export CSV_FILENAME FILTER_EXPR DOTPATHS [COLUMN_NAMES]", os.Args[0])
	}
	csvFName := params[0]
	filterExpr := params[1]
	dotPaths := strings.Split(params[2], ",")
	colNames := []string{}
	if len(params) < 4 {
		for _, val := range dotPaths {
			colNames = append(colNames, val)
		}
	} else {
		colNames = strings.Split(params[3], ",")
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

	if linesNo, err := collection.ExportCSV(fp, os.Stderr, filterExpr, dotPaths, colNames, showVerbose); err != nil {
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

func main() {
	app := cli.NewCli(dataset.Version)
	appName := app.AppName()

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
	app.BoolVar(&newLine, "nl,newline", true, "if set to false to suppress a trailing newline")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// Application Options
	app.StringVar(&collectionName, "c,collection", "", "sets the collection to be used")
	app.BoolVar(&useHeaderRow, "use-header-row", true, "use the header row as attribute names in the JSON document")
	app.BoolVar(&useUUID, "uuid", false, "generate a UUID for a new JSON document name")
	app.BoolVar(&showVerbose, "verbose", false, "output rows processed on importing from CSV")
	app.IntVar(&sampleSize, "sample", 0, "set the sample size when listing keys")

	// Action verbs (e.g. app.AddAction(STRING_VERB, FUNC_POINTER, STRING_DESCRIPTION)
	// NOTE: Sense this pre-existed cli v0.0.6 we're going to stick with what we evolved.

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
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
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

	// Merge environment
	datasetEnv := os.Getenv("DATASET")
	if datasetEnv != "" && collectionName == "" {
		collectionName = datasetEnv
	}

	if len(args) == 0 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("See %s --help for usage", appName), quiet)
	}

	in, err := cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(inputFName, in)

	out, err := cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(outputFName, out)

	action, params := args[0], args[1:]

	var data string
	if inputFName != "" {
		src, err := ioutil.ReadAll(in)
		cli.ExitOnError(os.Stderr, err, quiet)
		data = string(src)
	}

	fn, ok := voc[action]
	if ok == false {
		cli.ExitOnError(os.Stderr, fmt.Errorf("do not understand %s", action), quiet)
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
