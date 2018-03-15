//
// py/dataset.go is a C shared library for implementing a dataset module in Python3
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>

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

// #cgo pkg-config: python-3.6
// #define Py_LIMITED_API
// #include <Python.h>
import (
	"C"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheet"
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/tmplfn"
)

var (
	verbose          = false
	useStrictDotpath = true
	errorValue       error
)

func error_dispatch(err error, s string, values ...interface{}) {
	errorValue = err
	if verbose == true {
		log.Printf(s, values...)
	}
}

//export error_message
func error_message() *C.char {
	if errorValue != nil {
		s := fmt.Sprintf("%s", errorValue)
		errorValue = nil
		return C.CString(s)
	}
	return C.CString("")
}

//export use_strict_dotpath
func use_strict_dotpath(v C.int) C.int {
	if int(v) == 1 {
		useStrictDotpath = true
		return C.int(1)
	}
	useStrictDotpath = false
	return C.int(0)
}

//export is_verbose
func is_verbose() C.int {
	if verbose == true {
		return C.int(1)
	}
	return C.int(0)
}

//export verbose_on
func verbose_on() {
	verbose = true
}

//export verbose_off
func verbose_off() {
	verbose = false
}

func messagef(s string, values ...interface{}) {
	if verbose == true {
		log.Printf(s, values...)
	}
}

//export version
func version() *C.char {
	return C.CString(dataset.Version)
}

//export init_collection
func init_collection(name *C.char) C.int {
	collectionName := C.GoString(name)
	if verbose == true {
		messagef("creating %s\n", collectionName)
	}
	_, err := dataset.InitCollection(collectionName)
	if err != nil {
		messagef("Cannot create collection %s, %s", collectionName, err)
		return C.int(0)
	}
	messagef("%s initialized", collectionName)
	return C.int(1)
}

//export has_key
func has_key(name, key *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	if c.HasKey(k) {
		return C.int(1)
	}
	return C.int(0)
}

//export create_record
func create_record(name, key, src *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)
	v := []byte(C.GoString(src))

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	err = c.CreateJSON(k, v)
	if err != nil {
		messagef("Create %s failed, %s", k, err)
		return C.int(0)
	}
	return C.int(1)
}

//export read_record
func read_record(name, key *C.char) *C.char {
	collectionName := C.GoString(name)
	k := C.GoString(key)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()

	src, err := c.ReadJSON(k)
	if err != nil {
		messagef("Can't read %s, %s", k, err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//export update_record
func update_record(name, key, src *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)
	v := []byte(C.GoString(src))

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	err = c.UpdateJSON(k, v)
	if err != nil {
		messagef("Update %s failed, %s", k, err)
		return C.int(0)
	}
	return C.int(1)
}

//export delete_record
func delete_record(name, key *C.char) C.int {
	collectionName := C.GoString(name)
	k := C.GoString(key)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	err = c.Delete(k)
	if err != nil {
		messagef("Update %s failed, %s", k, err)
		return C.int(0)
	}
	return C.int(1)
}

//export join
func join(cName *C.char, cKey *C.char, cAdverb *C.char, cObjSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	adverb := C.GoString(cAdverb)
	objectSrc := C.GoString(cObjSrc)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer c.Close()

	outObject := map[string]interface{}{}
	newObject := map[string]interface{}{}

	if err := c.Read(key, outObject); err != nil {
		messagef("%s", err)
		return C.int(0)
	}

	if err := json.Unmarshal([]byte(objectSrc), &newObject); err != nil {
		messagef("%s", err)
		return C.int(0)
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
		messagef("Unknown join type %q", adverb)
		return C.int(0)
	}
	if err := c.Update(key, outObject); err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

//export keys
func keys(cname, cFilterExpr, cSortExpr *C.char) *C.char {
	collectionName := C.GoString(cname)
	filterExpr := C.GoString(cFilterExpr)
	sortExpr := C.GoString(cSortExpr)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()

	keyList := c.Keys()
	if filterExpr != "" {
		keyList, err = c.KeyFilter(keyList, filterExpr)
		if err != nil {
			messagef("Filter error, %s", err)
			return C.CString("")
		}
	}
	if sortExpr != "" {
		keyList, err = c.KeySortByExpression(keyList, sortExpr)
		if err != nil {
			messagef("Sort error, %s", err)
			return C.CString("")
		}
	}
	src, err := json.Marshal(keyList)
	if err != nil {
		messagef("Can't marshal key list, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//export key_filter
func key_filter(cname, cKeyListExpr, cFilterExpr *C.char) *C.char {
	collectionName := C.GoString(cname)
	keyListExpr := C.GoString(cKeyListExpr)
	filterExpr := C.GoString(cFilterExpr)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()

	keyList := []string{}
	if err := json.Unmarshal([]byte(keyListExpr), &keyList); err != nil {
		messagef("Unable to unmarshal keys", err)
		return C.CString("")
	}
	keys, err := c.KeyFilter(keyList, filterExpr)
	if err != nil {
		messagef("filter error, %s", err)
		return C.CString("")
	}
	src, err := json.Marshal(keys)
	if err != nil {
		messagef("Can't marshal filtered keys, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//export key_sort
func key_sort(cname, cKeyList, cSortExpr *C.char) *C.char {
	collectionName := C.GoString(cname)
	keyList := C.GoString(cKeyList)
	sortExpr := C.GoString(cSortExpr)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()

	keys := []string{}
	if err := json.Unmarshal([]byte(keyList), &keys); err != nil {
		messagef("Unable to unmarshal keys", err)
		return C.CString("")
	}
	keys, err = c.KeySortByExpression(keys, sortExpr)
	if err != nil {
		messagef("filter error, %s", err)
		return C.CString("")
	}
	src, err := json.Marshal(keys)
	if err != nil {
		messagef("Can't marshal sorted keys, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//export count
func count(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()
	i := c.Length()
	return C.int(i)
}

//export extract
func extract(cName, filterExpr, dotExpr *C.char) *C.char {
	collectionName := C.GoString(cName)
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()
	values, err := c.Extract(C.GoString(filterExpr), C.GoString(dotExpr))
	if err != nil {
		messagef("Extract failed for %s, %q, %q:  %s", collectionName, filterExpr, dotExpr, err)
		return C.CString("")
	}
	src, err := json.Marshal(values)
	if err != nil {
		messagef("Can't marshal extracted values for %s, %s", collectionName, err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//export indexer
func indexer(cName, cIndexName, cIndexMapName, cKeyList *C.char, cBatchSize C.int) C.int {
	collectionName := C.GoString(cName)
	indexName := C.GoString(cIndexName)
	indexMapName := C.GoString(cIndexMapName)
	keyList := C.GoString(cKeyList)
	batchSize := int(cBatchSize)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		// return 0 (false)
		return C.int(0)
	}
	defer c.Close()

	keys := []string{}
	if keyList != "" {
		err = json.Unmarshal([]byte(keyList), &keys)
		if err != nil {
			messagef("Can't unmarshal key list, %s", err)
			// return 0 (false)
			return C.int(0)
		}
	}

	err = c.Indexer(indexName, indexMapName, keys, batchSize)
	if err != nil {
		messagef("Indexing error %s %s, %s", collectionName, indexName, err)
		// return 0 (false)
		return C.int(0)
	}
	// return 1 (true) for success
	return C.int(1)
}

//export deindexer
func deindexer(cName, cIndexName, cKeyList *C.char, cBatchSize C.int) C.int {
	collectionName := C.GoString(cName)
	indexName := C.GoString(cIndexName)
	keyList := C.GoString(cKeyList)
	batchSize := int(cBatchSize)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Cannot open collection %s, %s", collectionName, err)
		// return 0 (false), failed
		return C.int(0)
	}
	defer c.Close()

	keys := []string{}
	if keyList != "" {
		err = json.Unmarshal([]byte(keyList), &keys)
		if err != nil {
			messagef("Can't unmarshal key list, %s", err)
			// return 0 (false), failed
			return C.int(0)
		}
	}

	err = c.Deindexer(indexName, keys, batchSize)
	if err != nil {
		messagef("Deindexing error %s %s, %s", collectionName, indexName, err)
		// return 0 (false), failed
		return C.int(0)
	}
	// return 1 (true) for success
	return C.int(1)
}

//export find
func find(cIndexNames, cQueryString, cOptionsMap *C.char) *C.char {
	indexNamesSrc := C.GoString(cIndexNames)
	queryString := C.GoString(cQueryString)
	optionsSrc := C.GoString(cOptionsMap)

	indexNames := []string{}
	if strings.HasPrefix(indexNamesSrc, "[") {
		err := json.Unmarshal([]byte(indexNamesSrc), &indexNames)
		if err != nil {
			messagef("Can't unmarshal index names, %s", err)
			return C.CString("")
		}
	} else if strings.Contains(indexNamesSrc, ":") {
		indexNames = strings.Split(indexNamesSrc, ":")
	} else {
		indexNames = []string{indexNamesSrc}
	}
	options := map[string]string{}
	if optionsSrc != "" {
		err := json.Unmarshal([]byte(optionsSrc), &options)
		if err != nil {
			messagef("Options error, %s", err)
			// return "", failed
			return C.CString("")
		}
	}

	idxList, _, err := dataset.OpenIndexes(indexNames)
	if err != nil {
		messagef("Can't open index %s, %s", strings.Join(indexNames, ", "), err)
		return C.CString("")
	}

	result, err := dataset.Find(idxList.Alias, queryString, options)
	if err != nil {
		messagef("Find error %s, %s", strings.Join(indexNames, ", "), err)
		// return "", failed
		return C.CString("")
	}
	err = idxList.Close()
	if err != nil {
		messagef("Can't close indexes %s, %s", strings.Join(indexNames, ", "), err)
	}

	src, err := json.Marshal(result)
	if err != nil {
		messagef("Can't marshal results, %s", err)
		// return "", failed
		return C.CString("")
	}

	txt := fmt.Sprintf("%s", src)
	// return our encoded results, success
	return C.CString(txt)
}

//export import_csv
func import_csv(cName *C.char, cCSVFName *C.char, cIDCol C.int, cUseHeaderRow C.int, cUseUUID C.int) C.int {
	// Covert options
	collectionName := C.GoString(cName)
	csvFName := C.GoString(cCSVFName)
	idCol := int(cIDCol)
	useHeaderRow := (int(cUseHeaderRow) == 1)
	useUUID := (int(cUseUUID) == 1)

	collection, err := dataset.Open(collectionName)
	if err != nil {
		messagef("Can't open %s, %s", collectionName, err)
		return C.int(0)
	}
	defer collection.Close()

	if idCol < 1 {
		messagef("Column number must be greater than zero, got %s", idCol)
		return C.int(0)
	}

	// NOTE: we need to adjust to zero based index
	idCol--
	fp, err := os.Open(csvFName)
	if err != nil {
		messagef("Can't open %s, %s", csvFName, err)
		return C.int(0)
	}
	defer fp.Close()

	if linesNo, err := collection.ImportCSV(fp, useHeaderRow, idCol, useUUID, verbose); err != nil {
		messagef("Can't import CSV, %s", err)
		return C.int(0)
	} else {
		messagef("%d total rows processed", linesNo)
	}
	return C.int(1)
}

//export export_csv
func export_csv(cName, cCSVFName, cFilterExpr, cDotExprs, cColNames *C.char) C.int {
	// Convert out parameters
	collectionName := C.GoString(cName)
	csvFName := C.GoString(cCSVFName)
	filterExpr := C.GoString(cFilterExpr)
	dotExprs := strings.Split(C.GoString(cDotExprs), ",")
	colNames := strings.Split(C.GoString(cColNames), ",")
	if len(colNames) == 0 {
		for _, val := range dotExprs {
			colNames = append(colNames, strings.TrimPrefix(val, "."))
		}
	}

	// Trim the any spaces for paths and column names
	for i, val := range dotExprs {
		dotExprs[i] = strings.TrimSpace(val)
	}
	for i, val := range colNames {
		colNames[i] = strings.TrimSpace(val)
	}

	collection, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer collection.Close()

	fp, err := os.Create(csvFName)
	if err != nil {
		messagef("Can't create %s, %s", csvFName, err)
		return C.int(0)
	}
	defer fp.Close()

	if linesNo, err := collection.ExportCSV(fp, os.Stderr, filterExpr, dotExprs, colNames, verbose); err != nil {
		messagef("Can't export CSV, %s", err)
		return C.int(0)
	} else {
		messagef("%d total rows processed", linesNo)
	}
	return C.int(1)
}

//export import_gsheet
func import_gsheet(cName, cClientSecretJSON, cSheetID, cSheetName, cCellRange *C.char, cIDCol C.int, cUseHeaderRow C.int, cUseUUID C.int, cOverwrite C.int) C.int {
	collectionName := C.GoString(cName)
	clientSecretJSON := C.GoString(cClientSecretJSON)
	sheetID := C.GoString(cSheetID)
	sheetName := C.GoString(cSheetName)
	cellRange := C.GoString(cCellRange)
	idCol := int(cIDCol)
	useHeaderRow := (C.int(cUseHeaderRow) == 1)
	useUUID := (C.int(cUseUUID) == 1)
	overwrite := (C.int(cOverwrite) == 1)

	collection, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer collection.Close()

	// NOTE: we need to adjust to zero based index
	idCol--

	table, err := gsheets.ReadSheet(clientSecretJSON, sheetID, sheetName, cellRange)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}

	linesNo, err := collection.ImportTable(table, useHeaderRow, idCol, useUUID, overwrite, verbose)
	if err != nil {
		messagef("Errors importing %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}
	messagef("%d total rows processed", linesNo)
	return C.int(1)
}

//export export_gsheet
func export_gsheet(cName, cClientSecretJSON, cSheetID, cSheetName, cCellRange, cFilterExpr, cDotExprs, cColNames *C.char) C.int {
	collectionName := C.GoString(cName)
	clientSecretJSON := C.GoString(cClientSecretJSON)
	sheetID := C.GoString(cSheetID)
	sheetName := C.GoString(cSheetName)
	cellRange := C.GoString(cCellRange)
	filterExpr := C.GoString(cFilterExpr)
	dotExprs := strings.Split(C.GoString(cDotExprs), ",")
	colNames := strings.Split(C.GoString(cColNames), ",")

	collection, err := dataset.Open(collectionName)
	if err != nil {
		messagef("failed, %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}
	defer collection.Close()

	if len(colNames) == 0 {
		for _, val := range dotExprs {
			colNames = append(colNames, strings.TrimPrefix(val, "."))
		}
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
		for i, key := range keys {
			m := map[string]interface{}{}
			if err := collection.Read(key, m); err == nil {
				row := []interface{}{}
				for j, colExpr := range dotExprs {
					col, err := dotpath.Eval(colExpr, m)
					if err != nil {
						if useStrictDotpath == true {
							messagef("failed, cell (%s: %d, %d), %s %s to evaluate %q, %s", key, i, j, sheetID, sheetName, colExpr, err)
							return C.int(0)
						}
						messagef("warning, cell (%s: %d, %d), %s %s to evaluate %q, %s", key, i, j, sheetID, sheetName, colExpr, err)
					}
					row = append(row, col)
				}
				table = append(table, row)
			}
		}
	} else {
		f, err := tmplfn.ParseFilter(filterExpr)
		if err != nil {
			messagef("failed, %s %s filter expression %q, %s", sheetID, sheetName, filterExpr, err)
			return C.int(0)
		}

		for i, key := range keys {
			m := map[string]interface{}{}
			if err := collection.Read(key, m); err == nil {
				if ok, err := f.Apply(m); err == nil && ok == true {
					// save row out.
					row := []interface{}{}
					for j, colExpr := range dotExprs {
						col, err := dotpath.Eval(colExpr, m)
						if err != nil {
							if useStrictDotpath == true {
								messagef("failed, cell (%s: %d, %d), %s %s to evaluate %q, %s", key, i, j, sheetID, sheetName, colExpr, err)
								return C.int(0)
							}
							messagef("warning, cell (%s: %d, %d), %s %s to evaluate %q, %s", key, i, j, sheetID, sheetName, colExpr, err)
						}
						row = append(row, col)
					}
					table = append(table, row)
				}
			}
		}
	}
	err = gsheets.WriteSheet(clientSecretJSON, sheetID, sheetName, cellRange, table)
	if err != nil {
		messagef("Failed to write %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}
	return C.int(1)
}

//export status
func status(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("failed, %s, %s", collectionName, err)
		return C.int(0)
	}
	c.Close()
	return C.int(1)
}

//export list
func list(cName *C.char, cKeys *C.char) *C.char {
	collectionName := C.GoString(cName)
	sKeys := C.GoString(cKeys)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.CString("")
	}
	defer c.Close()

	keys := []string{}
	err = json.Unmarshal([]byte(sKeys), &keys)
	if err != nil {
		messagef("Failed to unmarshal key list, %s", err)
		return C.CString("")
	}

	recs := []map[string]interface{}{}
	for _, name := range keys {
		m := map[string]interface{}{}
		err = c.Read(name, m)
		if err != nil {
			messagef("%s", err)
			return C.CString("")
		}
		recs = append(recs, m)
	}
	src, err := json.Marshal(recs)
	if err != nil {
		messagef("failed to marshal result, %s", err)
		return C.CString("")
	}
	return C.CString(string(src))
}

//export path
func path(cName *C.char, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.CString("")
	}
	defer c.Close()
	s, err := c.DocPath(key)
	if err != nil {
		messagef("%s", err)
		return C.CString("")
	}
	return C.CString(s)
}

//export check
func check(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := dataset.Analyzer(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

//export repair
func repair(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := dataset.Repair(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

//export attach
func attach(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			messagef("Can't unmarshal %q, %s", srcFNames, err)
			return C.int(0)
		}
	}

	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer c.Close()

	if c.HasKey(key) == false {
		messagef("%q is not in collection", key)
		return C.int(0)
	}
	for _, fname := range fNames {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			messagef("%s does not exist", fname)
			return C.int(0)
		}
	}
	err = c.AttachFiles(key, fNames...)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

//export attachments
func attachments(cName *C.char, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.CString("")
	}
	defer c.Close()
	if c.HasKey(key) == false {
		messagef("%q is not in collection", key)
		return C.CString("")
	}
	results, err := c.Attachments(key)
	if err != nil {
		messagef("%s", err)
		return C.CString("")
	}
	if len(results) > 0 {
		return C.CString(strings.Join(results, "\n"))
	}
	return C.CString("")
}

//export detach
func detach(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			messagef("Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer c.Close()
	if c.HasKey(key) == false {
		messagef("%q is not in collection", key)
		return C.int(0)
	}
	err = c.GetAttachedFiles(key, fNames...)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

//export prune
func prune(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			messagef("Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	defer c.Close()

	err = c.Prune(key, fNames...)
	if err != nil {
		messagef("%s", err)
		return C.int(0)
	}
	return C.int(1)
}

func main() {}
