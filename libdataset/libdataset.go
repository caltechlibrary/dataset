//
// py/dataset.go is a C shared library for implementing a dataset module in Python3
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>

// Copyright (c) 2019, Caltech
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
	"C"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/gsheets"
	"github.com/caltechlibrary/dataset/tbl"
)

var (
	verbose          = false
	useStrictDotpath = true
	// NOTE: error state is shared because C doesn't easily pass multiple
	// return values without resorting to complex structures.
	errorValue error
)

// error_clear will set the global error state to nil.
//
//export error_clear
func error_clear() {
	errorValue = nil
}

// error_dispatch logs error messages to console based on string template
// Not exported.
//
func error_dispatch(err error, s string, values ...interface{}) {
	errorValue = err
	if verbose == true {
		log.Printf(s, values...)
	}
}

// error_message returns an error message previously recorded or
// an empty string if no errors recorded
//
//export error_message
func error_message() *C.char {
	if errorValue != nil {
		s := fmt.Sprintf("%s", errorValue)
		errorValue = nil
		return C.CString(s)
	}
	return C.CString("")
}

// use_strict_dotpath sets the library option value for
// enforcing strict dotpaths. 1 is true, any other value is false.
//
//export use_strict_dotpath
func use_strict_dotpath(v C.int) C.int {
	if int(v) == 1 {
		useStrictDotpath = true
		return C.int(1)
	}
	useStrictDotpath = false
	return C.int(0)
}

// is_verbose returns the library options' verbose value.
//
//export is_verbose
func is_verbose() C.int {
	if verbose == true {
		return C.int(1)
	}
	return C.int(0)
}

// verbose_on set library verbose to true
//
//export verbose_on
func verbose_on() {
	verbose = true
}

// verbose_off set library verbose to false
//
//export verbose_off
func verbose_off() {
	verbose = false
}

// messagef is an intertal library function for logging messages to
//
// the console. Not exported.
func messagef(s string, values ...interface{}) {
	if verbose == true {
		log.Printf(s, values...)
	}
}

// dataset_version returns the dataset version the libdataset presents.
//
//export dataset_version
func dataset_version() *C.char {
	return C.CString(dataset.Version)
}

// init_collection intializes a collection and records as much metadata
// as it can from the execution environment (e.g. username,
// datetime created)
//
//export init_collection
func init_collection(name *C.char) C.int {
	collectionName := C.GoString(name)
	if verbose == true {
		messagef("creating %s\n", collectionName)
	}
	error_clear()
	_, err := dataset.InitCollection(collectionName)
	if err != nil {
		error_dispatch(err, "Cannot create collection %s, %s", collectionName, err)
		return C.int(0)
	}
	messagef("%s initialized", collectionName)
	return C.int(1)
}

// is_open returns true (i.e. one) if a collection has been opened by libdataset, false (i.e. zero) otherwise
//
// export is_open
func is_open(cName *C.char) C.int {
	collectionName := C.GoString(cName)

	if dataset.IsOpen(collectionName) {
		return C.int(1)
	}
	return C.int(0)
}

// open returns 0 on successfully opening a collection 1 otherwise. Sets error messages if needed.
//
// export open
func open(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	_, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "Cannot open %q, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// collections returns a JSON list of collection names that are open otherwise an empty list.
//
// export collections
func collections() *C.char {
	cNames := dataset.Collections()
	src, err := json.Marshal(cNames)
	if err != nil {
		return C.CString("[]")
	}
	return C.CString(fmt.Sprintf("%s", src))
}

// close closes a collection previously opened.
//
// export close
func close(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if err := dataset.Close(collectionName); err != nil {
		error_dispatch(err, "Cannot close collection %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// close_all closes all collections previously opened
//
// export close_all
func close_all() C.int {
	error_clear()
	if err := dataset.CloseAll(); err != nil {
		error_dispatch(err, "Cannot close all collections, %s", err)
		return C.int(0)
	}
	return C.int(1)
}

// create_object takes JSON source and adds it to the collection with
// the provided key.
//
//export create_object
func create_object(cName, cKey, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	src := []byte(C.GoString(cSrc))

	error_clear()
	err := dataset.CreateJSON(collectionName, key, src)
	if err != nil {
		error_dispatch(err, "Create %s failed, %s", key, err)
		return C.int(0)
	}
	return C.int(1)
}

// read_object takes a key and returns JSON source of the record
//
//export read_object
func read_object(cName, cKey *C.char, cCleanObject C.int) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	cleanObject := (C.int(1) == cCleanObject)

	error_clear()
	var (
		src []byte
		err error
	)
	if src, err = dataset.ReadJSON(collectionName, key); err != nil {
		error_dispatch(err, "Can't read %s, %s", key, err)
		return C.CString("")
	}
	if cleanObject {
		object := map[string]interface{}{}
		if err := dataset.DecodeJSON(src, &object); err != nil {
			error_dispatch(err, "Can't decode %s, %s", key, err)
			return C.CString("")
		}
		if _, found := object["_Key"]; found {
			delete(object, "_Key")
			src, err = dataset.EncodeJSON(object)
			if err != nil {
				error_dispatch(err, "Can't encode %s, %s", key, err)
				return C.CString("")
			}
		}
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// THIS IS AN UGLY HACK, Python ctypes doesn't **easily** support
// undemensioned arrays of strings. So we will assume the array of
// keys has already been transformed into JSON before calling
// read_list.
//
//export read_object_list
func read_object_list(cName *C.char, cKeysAsJSON *C.char, cCleanObject C.int) *C.char {
	collectionName := C.GoString(cName)
	cleanObject := (C.int(1) == cCleanObject)
	l := []string{}
	errList := []string{}

	error_clear()
	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(cKeysAsJSON))
	keyList := []string{}
	err := json.Unmarshal(src, &keyList)
	if err != nil {
		error_dispatch(err, "Can't unmarshal key list, %s", err)
		return C.CString("")
	}
	if dataset.IsOpen(collectionName) == false {
		if err = dataset.Open(collectionName); err != nil {
			error_dispatch(err, "Cannot open collection %s, %s", collectionName, err)
			return C.CString("")
		}
	}

	for _, key := range keyList {
		src, err := dataset.ReadJSON(collectionName, key)
		if err != nil {
			errList = append(errList, fmt.Sprintf("(%s) %s", key, err))
			continue
		}
		if cleanObject == true {
			obj := map[string]interface{}{}
			if err = dataset.DecodeJSON(src, &obj); err != nil {
				errList = append(errList, fmt.Sprintf("(%s) %s", key, err))
			} else {
				delete(obj, "_Key")
				if src, err := json.Marshal(obj); err == nil {
					l = append(l, fmt.Sprintf("%s", src))
				} else {
					errList = append(errList, fmt.Sprintf("(%s) %s", key, err))
				}
			}
		} else {
			l = append(l, fmt.Sprintf("%s", src))
		}
	}
	if len(errList) > 0 {
		err = fmt.Errorf("%s", strings.Join(errList, "; "))
		error_dispatch(err, "Key read errors %s", err)
	}

	txt := fmt.Sprintf("[%s]", strings.Join(l, ","))
	return C.CString(txt)
}

// update_object takes a key and JSON source and replaces the record
// in the collection.
//
//export update_object
func update_object(cName, cKey, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	src := []byte(C.GoString(cSrc))

	error_clear()
	err := dataset.UpdateJSON(collectionName, key, src)
	if err != nil {
		error_dispatch(err, "Update %s failed, %s", key, err)
		return C.int(0)
	}
	return C.int(1)
}

// delete_object takes a key and removes a record from the collection
//
//export delete_object
func delete_object(cName, cKey *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	error_clear()
	err := dataset.DeleteJSON(collectionName, key)
	if err != nil {
		error_dispatch(err, "Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// join takes a collection name, a key, and merges JSON source with an
// existing JSON record. If overwrite is 1 it overwrites and replaces
// common values, if not 1 it only adds missing attributes.
//
//export join
func join(cName *C.char, cKey *C.char, cObjSrc *C.char, cOverwrite C.int) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	overwrite := (cOverwrite == 1)
	objectSrc := C.GoString(cObjSrc)

	error_clear()
	err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	newObject := map[string]interface{}{}
	if err := dataset.DecodeJSON([]byte(objectSrc), &newObject); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	src, err := dataset.ReadJSON(collectionName, key)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	outObject := map[string]interface{}{}
	if err := dataset.DecodeJSON(src, &outObject); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	// Merge out objects
	if overwrite {
		for k, v := range newObject {
			outObject[k] = v
		}
	} else {
		for k, v := range newObject {
			if _, ok := outObject[k]; ok != true {
				outObject[k] = v
			}
		}
	}
	src, err = dataset.EncodeJSON(outObject)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	if err := dataset.UpdateJSON(collectionName, key, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// key_exists returns 1 if the key exists in a collection or 0 if not.
//
//export key_exists
func key_exists(cName, cKey *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	if dataset.KeyExists(collectionName, key) {
		return C.int(1)
	}
	return C.int(0)
}

// keys returns JSON source of an array of keys from the collection
//
//export keys
func keys(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	error_clear()
	keyList := dataset.Keys(collectionName)
	src, err := json.Marshal(keyList)
	if err != nil {
		error_dispatch(err, "Can't marshal key list, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// key_filter returns JSON source of an array of keys passing
// through the filter of objects in the collection.
//
//export key_filter
func key_filter(cName, cKeyListExpr, cFilterExpr *C.char) *C.char {
	collectionName := C.GoString(cName)
	keyListExpr := C.GoString(cKeyListExpr)
	filterExpr := C.GoString(cFilterExpr)

	error_clear()
	keyList := []string{}
	if err := json.Unmarshal([]byte(keyListExpr), &keyList); err != nil {
		error_dispatch(err, "Unable to unmarshal keys", err)
		return C.CString("")
	}
	keys, err := dataset.KeyFilter(collectionName, keyList, filterExpr)
	if err != nil {
		error_dispatch(err, "filter error, %s", err)
		return C.CString("")
	}
	src, err := json.Marshal(keys)
	if err != nil {
		error_dispatch(err, "Can't marshal filtered keys, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// key_sort returns JSON source of an array of keys sorted by
// the sort expression applied to the objects in the collection.
//
//export key_sort
func key_sort(cName, cKeyListExpr, cSortExpr *C.char) *C.char {
	collectionName := C.GoString(cName)
	keyListExpr := C.GoString(cKeyListExpr)
	sortExpr := C.GoString(cSortExpr)

	error_clear()
	keyList := []string{}
	if err := json.Unmarshal([]byte(keyListExpr), &keyList); err != nil {
		error_dispatch(err, "Unable to unmarshal keys", err)
		return C.CString("")
	}

	keys, err := dataset.KeySortByExpression(collectionName, keyList, sortExpr)
	if err != nil {
		error_dispatch(err, "Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}
	src, err := json.Marshal(keys)
	if err != nil {
		error_dispatch(err, "Can't marshal sorted keys, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// count returns the number of objects (records) in a collection.
// if an error is encounter a -1 is returned.
//export count
func count(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if dataset.IsOpen(collectionName) {
		if c, err := dataset.GetCollection(collectionName); err == nil {
			i := c.Length()
			return C.int(i)
		}
	}
	return C.int(-1)
}

// import_csv - import a CSV file into a collection
// syntax: COLLECTION CSV_FILENAME ID_COL
//
// options that should support sensible defaults:
//
//     cUseHeaderRow
//     cOverwrite
//
//export import_csv
func import_csv(cName *C.char, cCSVFName *C.char, cIDCol C.int, cUseHeaderRow C.int, cOverwrite C.int) C.int {
	// Covert options
	collectionName := C.GoString(cName)
	csvFName := C.GoString(cCSVFName)
	idCol := int(cIDCol)
	useHeaderRow := (int(cUseHeaderRow) == 1)
	overwrite := (int(cOverwrite) == 1)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "Can't open %s, %s", collectionName, err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if idCol < 1 {
		error_dispatch(fmt.Errorf("invalid column number"), "Column number must be greater than zero, got %s", idCol)
		return C.int(0)
	}

	// NOTE: we need to adjust to zero based index
	idCol--

	// Now import our CSV file
	fp, err := os.Open(csvFName)
	if err != nil {
		error_dispatch(err, "Can't open %s, %s", csvFName, err)
		return C.int(0)
	}
	cnt, err := c.ImportCSV(fp, idCol, useHeaderRow, overwrite, verbose)
	if err != nil {
		error_dispatch(err, "%s\n", err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)

	return C.int(1)
}

// export_csv - export collection objects to a CSV file
// syntax: COLLECTION FRAME CSV_FILENAME
//
//export export_csv
func export_csv(cName *C.char, cFrameName *C.char, cCSVFName *C.char) C.int {
	// Convert out parameters
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFrameName)
	csvFName := C.GoString(cCSVFName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "Can't open %s, %s", collectionName, err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "Can't open %s, %s", collectionName, err)
		return C.int(0)
	}

	fp, err := os.Create(csvFName)
	if err != nil {
		error_dispatch(err, "Can't create %s, %s", csvFName, err)
		return C.int(0)
	}
	defer fp.Close()

	// Get Frame
	if c.FrameExists(frameName) == false {
		error_dispatch(err, "Missing frame %q in %s\n", frameName, collectionName)
		return C.int(0)
	}

	// Get dotpaths and column labels from frame
	f, err := c.FrameRead(frameName)
	if err != nil {
		error_dispatch(err, "%s\n", err)
		return C.int(0)
	}

	// Now export to CSV
	cnt, err := c.ExportCSV(fp, os.Stderr, f, verbose)
	if err != nil {
		error_dispatch(err, "Can't export CSV %s, %s", csvFName, err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)
	return C.int(1)
}

// import_gsheet - import a GSheet into a collection
// syntax: COLLECTION GSHEET_ID SHEET_NAME ID_COL CELL_RANGE
//
// options that should support sensible defaults:
//
//    cUseHeaderRow
//    cOverwrite
//
//export import_gsheet
func import_gsheet(cName *C.char, cSheetID *C.char, cSheetName *C.char, cIDCol C.int, cCellRange *C.char, cUseHeaderRow C.int, cOverwrite C.int) C.int {
	collectionName := C.GoString(cName)
	sheetID := C.GoString(cSheetID)
	sheetName := C.GoString(cSheetName)
	cellRange := C.GoString(cCellRange)
	idCol := int(cIDCol)
	useHeaderRow := (C.int(cUseHeaderRow) == 1)
	overwrite := (C.int(cOverwrite) == 1)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	// NOTE: we need to adjust to zero based index
	idCol--

	//FIXME: Need better search process for finding the google access key
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretJSON == "" {
		//clientSecretJSON = "client_secret.json"
		clientSecretJSON = "credentials.json"
	}

	table, err := gsheets.ReadSheet(clientSecretJSON, sheetID, sheetName, cellRange)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	cnt, err := c.ImportTable(table, idCol, useHeaderRow, overwrite, verbose)
	if err != nil {
		error_dispatch(err, "Errors importing %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)
	return C.int(1)
}

// export_gsheet - export collection objects to a GSheet
// syntax: COLLECTION FRAME GSHEET_ID GSHEET_NAME CELL_RANGE
//
//export export_gsheet
func export_gsheet(cName *C.char, cFrameName *C.char, cSheetID *C.char, cSheetName *C.char, cCellRange *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFrameName)
	sheetID := C.GoString(cSheetID)
	sheetName := C.GoString(cSheetName)
	cellRange := C.GoString(cCellRange)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "failed, %s %s, %s", sheetID, sheetName, err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "failed, %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}

	// Get Frame
	if c.FrameExists(frameName) == false {
		error_dispatch(err, "Missing frame %q in %s\n", frameName, collectionName)
		return C.int(0)
	}

	// Get dotpaths and column labels from frame
	f, err := c.FrameRead(frameName)
	if err != nil {
		error_dispatch(err, "%s\n", err)
		return C.int(0)
	}

	//FIXME: Need a better way to indentify the clientSecretName...
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretJSON == "" {
		//clientSecretJSON = "client_secret.json"
		clientSecretJSON = "credentials.json"
	}
	// gSheet expects a cell range, so we will generate one if needed.
	if cellRange == "" {
		lastCol := gsheets.ColNoToColLetters(len(f.Labels))
		lastRow := len(f.Keys) + 2
		cellRange = fmt.Sprintf("A1:%s%d", lastCol, lastRow)
	}

	//NOTE: we export to GSheet via creating a table [][]interface{}{}
	cnt, table, err := c.ExportTable(os.Stderr, f, verbose)
	if err != nil {
		error_dispatch(err, "%s\n", err)
		return C.int(0)
	}
	err = gsheets.WriteSheet(clientSecretJSON, sheetID, sheetName, cellRange, table)

	if err != nil {
		error_dispatch(err, "Failed to write %s %s, %s", sheetID, sheetName, err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)
	return C.int(1)
}

// status checks to see if a collection exists or not.
//
//export status
func status(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if _, err := dataset.GetCollection(collectionName); err != nil {
		error_dispatch(err, "failed, %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// list returns JSON array of objects in a collections based on a
// JSON array of keys.
//
//export list
func list(cName *C.char, cKeys *C.char) *C.char {
	collectionName := C.GoString(cName)
	sKeys := C.GoString(cKeys)

	error_clear()
	keys := []string{}
	err := json.Unmarshal([]byte(sKeys), &keys)
	if err != nil {
		error_dispatch(err, "Failed to unmarshal key list, %s", err)
		return C.CString("")
	}

	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}

	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "Can't find %q", collectionName)
		return C.CString("")
	}

	recs := []map[string]interface{}{}
	for _, name := range keys {
		m := map[string]interface{}{}
		err = c.Read(name, m, false)
		if err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
		recs = append(recs, m)
	}
	src, err := json.Marshal(recs)
	if err != nil {
		error_dispatch(err, "failed to marshal result, %s", err)
		return C.CString("")
	}
	return C.CString(string(src))
}

// path returns the path on disc to an JSON object document
// in the collection.
//
//export path
func path(cName *C.char, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}

	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "Can't find %q", collectionName)
		return C.CString("")
	}

	s, err := c.DocPath(key)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	return C.CString(s)
}

// check runs the analyzer over a collection and looks for
// problem records.
//
//export check
func check(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := dataset.Check(collectionName, verbose)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// repair runs the analyzer over a collection and repairs JSON
// objects and attachment discovered having a problem. Also is
// useful for upgrading a collection between dataset releases.
//
//export repair
func repair(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := dataset.Repair(collectionName, verbose)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// attach will attach a file to a JSON object in a collection. It takes
// a semver string (e.g. v0.0.1) and associates that with where it stores
// the file.  If semver is v0.0.0 it is considered unversioned, if v0.0.1
// or larger it is considered versioned.
//
//export attach
func attach(cName *C.char, cKey *C.char, cSemver *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	semver := C.GoString(cSemver)
	if semver == "" {
		semver = "v0.0.0"
	}
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			error_dispatch(err, "Can't unmarshal %q, %s", srcFNames, err)
			return C.int(0)
		}
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if c.KeyExists(key) == false {
		error_dispatch(fmt.Errorf("missing key"), "%q is not in collection", key)
		return C.int(0)
	}
	// NOTE: fName, fNames are file names NOT frame names
	for _, fname := range fNames {
		if _, err := os.Stat(fname); os.IsNotExist(err) {
			error_dispatch(err, "%s does not exist", fname)
			return C.int(0)
		}
	}
	err = c.AttachFiles(key, semver, fNames...)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// attachments returns a list of attachments and their size in
// associated with a JSON obejct in the collection.
//
//export attachments
func attachments(cName *C.char, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	if c.KeyExists(key) == false {
		error_dispatch(fmt.Errorf("missing key"), "%q is not in collection", key)
		return C.CString("")
	}
	results, err := c.Attachments(key)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	if len(results) > 0 {
		return C.CString(strings.Join(results, "\n"))
	}
	return C.CString("")
}

// detach exports the file associated with the semver from the JSON
// object in the collection. The file remains "attached".
//
//export detach
func detach(cName *C.char, cKey *C.char, cSemver *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	semver := C.GoString(cSemver)
	if semver == "" {
		semver = "v0.0.0"
	}
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			error_dispatch(err, "Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	}
	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if c.KeyExists(key) == false {
		error_dispatch(err, "%q is not in collection", key)
		return C.int(0)
	}
	err = c.GetAttachedFiles(key, semver, fNames...)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// prune removes an attachment by semver from a JSON object in the
// collection. This is destructive, the file is removed from disc.
//
//export prune
func prune(cName *C.char, cKey *C.char, cSemver *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	semver := C.GoString(cSemver)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			error_dispatch(err, "Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	err = c.Prune(key, semver, fNames...)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// clone takes a collection name, a JSON array of keys and creates
// a new collection with a new name based on the origin's collections'
// objects.
//
//export clone
func clone(cName *C.char, cKeys *C.char, dName *C.char) C.int {
	collectionName := C.GoString(cName)
	srcKeys := C.GoString(cKeys)
	destName := C.GoString(dName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	keys := []string{}
	err = json.Unmarshal([]byte(srcKeys), &keys)
	if err != nil {
		error_dispatch(err, "Can't unmarshal keys, %s", err)
		return C.int(0)
	}
	err = c.Clone(destName, keys, verbose)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// clone_sample is like clone both generates a sample or test and
// training set of sampled of the cloned collection.
//
//export clone_sample
func clone_sample(cName *C.char, cTrainingName *C.char, cTestName *C.char, cSampleSize C.int) C.int {
	collectionName := C.GoString(cName)
	sampleSize := int(cSampleSize)
	trainingName := C.GoString(cTrainingName)
	testName := C.GoString(cTestName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	keys := c.Keys()
	err = c.CloneSample(trainingName, testName, keys, sampleSize, verbose)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// grid generates a "Grid" structure from a collection.
//
//export grid
func grid(cName *C.char, cKeys *C.char, cDotPaths *C.char) *C.char {
	collectionName := C.GoString(cName)
	srcKeys := C.GoString(cKeys)
	srcDotpaths := C.GoString(cDotPaths)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	keys := []string{}
	err = json.Unmarshal([]byte(srcKeys), &keys)
	if err != nil {
		error_dispatch(err, "Can't unmarshal keys, %s", err)
		return C.CString("")
	}
	dotPaths := []string{}
	err = json.Unmarshal([]byte(srcDotpaths), &dotPaths)
	if err != nil {
		error_dispatch(err, "Can't unmarshal dot paths, %s", err)
		return C.CString("")
	}
	//NOTE: We're picking up the verbose flag from the modules global state
	g, err := c.Grid(keys, dotPaths, verbose)
	if err != nil {
		error_dispatch(err, "failed to create grid, %s", err)
		return C.CString("")
	}
	src, err := json.Marshal(g)
	if err != nil {
		error_dispatch(err, "failed to marshal grid, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// frame_exists returns 1 (true) if frame name exists in collection, 0 (false) otherwise
//
//export frame_exists
func frame_exists(cName *C.char, cFName *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	if dataset.FrameExists(collectionName, frameName) {
		return C.int(1)
	}
	return C.int(0)
}

// frame_keys takes a collection name and frame name and returns a list of keys from the frame or an empty list.
// The list is expressed as a JSON source.
//
//export frame_keys
func frame_keys(cName *C.char, cFName *C.char) *C.char {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	keys := dataset.FrameKeys(collectionName, frameName)
	src, err := json.Marshal(keys)
	if err != nil {
		return C.CString("[]")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// frame_create defines a new frame an populates it.
//
//export frame_create
func frame_create(cName *C.char, cFName *C.char, cKeysSrc *C.char, cDotPathsSrc *C.char, cLabelsSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	keysSrc := []byte(C.GoString(cKeysSrc))
	dotPathSrc := []byte(C.GoString(cDotPathsSrc))
	labelsSrc := []byte(C.GoString(cLabelsSrc))
	keys := []string{}
	dotPaths := []string{}
	labels := []string{}

	error_clear()
	if err := json.Unmarshal(keysSrc, &keys); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	if err := json.Unmarshal(dotPathSrc, &dotPaths); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	if err := json.Unmarshal(labelsSrc, &labels); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}

	_, err := dataset.FrameCreate(collectionName, frameName, keys, dotPaths, labels, verbose)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frame_objects retrieves a JSON source list of objects from a frame.
//
//export frame_objects
func frame_objects(cName *C.char, cFName *C.char) *C.char {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	error_clear()
	ol, err := dataset.FrameObjects(collectionName, frameName)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	src, err := json.Marshal(ol)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// frame_refresh refresh the contents of a frame given a list of keys.
//
//export frame_refresh
func frame_refresh(cName *C.char, cFName *C.char, cKeysSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	keysSrc := C.GoString(cKeysSrc)

	keys := []string{}

	error_clear()
	if err := json.Unmarshal([]byte(keysSrc), &keys); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	if err := dataset.FrameRefresh(collectionName, frameName, keys, verbose); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frame_reframe will change of object list in a frame based on the key list provided.
//
//export frame_reframe
func frame_reframe(cName *C.char, cFName *C.char, cKeysSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	keysSrc := C.GoString(cKeysSrc)

	keys := []string{}

	error_clear()
	if err := json.Unmarshal([]byte(keysSrc), &keys); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	if err := dataset.FrameReframe(collectionName, frameName, keys, verbose); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frame_clear will clear the object list and keys associated with a frame.
//
//export frame_clear
func frame_clear(cName *C.char, cFName *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	error_clear()
	if err := dataset.FrameClear(collectionName, frameName); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frame_delete will removes a frame from a collection
//
//export frame_delete
func frame_delete(cName *C.char, cFName *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	error_clear()
	if err := dataset.FrameClear(collectionName, frameName); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frames returns a JSON array of frames names in the collection.
//
//export frames
func frames(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}

	frameNames := dataset.Frames(collectionName)
	if len(frameNames) == 0 {
		return C.CString("[]")
	}
	src, err := json.Marshal(frameNames)
	if err != nil {
		error_dispatch(err, "failed to marshal frame names, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// sync_send_csv - synchronize a frame sending data to a CSV file
//
//export sync_send_csv
func sync_send_csv(cName *C.char, cFName *C.char, cCSVFilename *C.char, cSyncOverwrite C.int) C.int {
	var (
		c   *dataset.Collection
		src []byte
		err error
	)
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	csvFilename := C.GoString(cCSVFilename)
	syncOverwrite := (cSyncOverwrite == 1)

	src, err = ioutil.ReadFile(csvFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}
	if len(src) == 0 {
		fmt.Fprintf(os.Stderr, "No data in csv file %s\n", csvFilename)
		return C.int(1)
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err == nil {
			table = tbl.TableStringToInterface(csvTable)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(1)
		}
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(1)
		}
	}
	c, err = dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(1)

	}

	// Merge collection content into table
	table, err = c.MergeIntoTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}

	// Save the resulting table
	if len(src) > 0 {
		if err = os.Rename(csvFilename, csvFilename+".bak"); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(1)
		}
		if out, err := os.Create(csvFilename); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(1)
		} else {
			w := csv.NewWriter(out)
			w.WriteAll(tbl.TableInterfaceToString(table))
			err = w.Error()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return C.int(1)
			}
		}
	}
	return C.int(0)
}

// sync_recieve_csv - synchronize a frame recieving data from a CSV file
//
//export sync_recieve_csv
func sync_recieve_csv(cName *C.char, cFName *C.char, cCSVFilename *C.char, cSyncOverwrite C.int) C.int {
	var (
		src []byte
		err error
	)
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	csvFilename := C.GoString(cCSVFilename)
	syncOverwrite := (cSyncOverwrite == 1)

	src, err = ioutil.ReadFile(csvFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(1)
		}
		table = tbl.TableStringToInterface(csvTable)
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(1)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(1)

	}

	// Merge table contents into Collection and Frame
	err = c.MergeFromTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}
	return C.int(0)
}

// sync_send - synchronize a frame sending data to a GSheet
//
//export sync_send_gsheet
func sync_send_gsheet(cName, cFName, cGSheetID, cGSheetName, cCellRange *C.char, cSyncOverwrite C.int) C.int {

	var (
		err error
	)
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	gSheetID := C.GoString(cGSheetID)
	gSheetName := C.GoString(cGSheetName)
	cellRange := C.GoString(cCellRange)
	syncOverwrite := (cSyncOverwrite == 1)

	table := [][]interface{}{}
	// Populate table to sync
	// for GSheet
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretJSON == "" {
		//clientSecretJSON = "client_secret.json"
		clientSecretJSON = "credentials.json"
	}
	table, err = gsheets.ReadSheet(clientSecretJSON, gSheetID, gSheetName, cellRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(1)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(1)

	}

	// Merge collection content into table
	table, err = c.MergeIntoTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}

	// Save the resulting table
	clientSecretJSON = os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretJSON == "" {
		//clientSecretJSON = "client_secret.json"
		clientSecretJSON = "credentials.json"
	}
	// NOTE: WriteSheet expects a [][]interface{} not [][]string,
	// need to convert. This is a hack...
	t := [][]interface{}{}
	for _, row := range table {
		cells := []interface{}{}
		for _, cell := range row {
			cells = append(cells, cell)
		}
		t = append(t, cells)
	}
	err = gsheets.WriteSheet(clientSecretJSON, gSheetID, gSheetName, cellRange, t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}
	return C.int(0)
}

// sync_recieve_gsheet - synchronize a frame recieving data from a GSheet
//
//export sync_recieve_gsheet
func sync_recieve_gsheet(cName, cFName, cGSheetID, cGSheetName, cCellRange *C.char, cSyncOverwrite C.int) C.int {
	var (
		err error
	)
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	gSheetID := C.GoString(cGSheetID)
	gSheetName := C.GoString(cGSheetName)
	cellRange := C.GoString(cCellRange)
	syncOverwrite := (cSyncOverwrite == 1)

	if cellRange == "" {
		cellRange = "A1:Z"
	}

	table := [][]interface{}{}
	// Populate table to sync
	// for GSheet
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	if clientSecretJSON == "" {
		//clientSecretJSON = "client_secret.json"
		clientSecretJSON = "credentials.json"
	}
	table, err = gsheets.ReadSheet(clientSecretJSON, gSheetID, gSheetName, cellRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(1)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(1)

	}

	// Merge table contents into Collection and Frame
	err = c.MergeFromTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(1)
	}
	return C.int(0)
}

// frame_grid takes a frames object list and returns a grid
// (2D JSON array) representation of the object list.
// If the "header row" value is 1 a header row of labels is
// included, otherwise it is only the values of returned in the grid.
//
//export frame_grid
func frame_grid(cName *C.char, cFName *C.char, cIncludeHeaderRow C.int) *C.char {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	includeHeaderRow := false
	if cIncludeHeaderRow == C.int(1) {
		includeHeaderRow = true
	}

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.CString("")
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	f, err := c.FrameRead(frameName)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	g := f.Grid(includeHeaderRow)
	src, err := json.Marshal(g)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

//
// make_objects - is a function to creates empty a objects in batch.
// It requires a JSON list of keys to create. For each key present
// an attempt is made to create a new empty object based on the JSON
// provided (e.g. `{}`, `{"is_empty": true}`). The reason to do this
// is that it means the collection.json file is updated once for the
// whole call and that the keys are now reserved to be updated separately.
// Returns 1 on success, 0 if errors encountered.
//
//export make_objects
func make_objects(cName *C.char, keysAsJson *C.char, objectAsJson *C.char) C.int {
	collectionName := C.GoString(cName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(keysAsJson))
	keyList := []string{}
	err = json.Unmarshal(src, &keyList)
	if err != nil {
		error_dispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}
	objectSrc := []byte(C.GoString(objectAsJson))

	err = c.CreateObjectsJSON(keyList, objectSrc)
	if err != nil {
		error_dispatch(err, "Create objects failed, %s", err)
		return C.int(0)
	}
	return C.int(1)
}

//
// update_objects - is a function to update objects in batch.
// It requires a JSON array of keys and a JSON array of
// matching objects. The list of keys and objects are processed
// together with calls to update individual records. Returns 1 on
// success, 0 on error.
//
//export update_objects
func update_objects(cName *C.char, keysAsJson *C.char, objectsAsJson *C.char) C.int {
	collectionName := C.GoString(cName)

	error_clear()
	if dataset.IsOpen(collectionName) == false {
		if err := dataset.Open(collectionName); err != nil {
			error_dispatch(err, "%s", err)
			return C.int(0)
		}
	}
	c, err := dataset.GetCollection(collectionName)
	if err != nil {
		error_dispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(keysAsJson))
	keyList := []string{}
	err = json.Unmarshal(src, &keyList)
	if err != nil {
		error_dispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}
	src = []byte(C.GoString(objectsAsJson))
	objectList := []map[string]interface{}{}
	err = json.Unmarshal(src, &objectList)
	if err != nil {
		error_dispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}

	if len(keyList) != len(objectList) {
		error_dispatch(err, "expected %d keys for %d objects", len(keyList), len(objectList))
		return C.int(0)
	}

	errorNo := 1
	for i, key := range keyList {
		if c.KeyExists(key) {
			err = c.Update(key, objectList[i])
			if err != nil {
				error_dispatch(err, "Can't update key %q, %s", key, err)
				errorNo = 0
			}
		}
	}
	return C.int(errorNo)
}

// set_who will set the "who" value associated with the collection's metadata
//
//export set_who
func set_who(cName *C.char, cNamesSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	namesSrc := []byte(C.GoString(cNamesSrc))
	names := []string{}
	error_clear()
	if err := json.Unmarshal(namesSrc, &names); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)
	}
	if err := dataset.SetWho(collectionName, names); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_who will get the "who" value associated with the collection's metadata
//
//export get_who
func get_who(cName *C.char) *C.char {
	collectionName := C.GoString(cName)
	error_clear()

	names := dataset.GetWho(collectionName)
	src, err := json.Marshal(names)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// set_what will set the "what" value associated with the collection's metadata
//
//export set_what
func set_what(cName *C.char, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	src := C.GoString(cSrc)
	error_clear()

	if err := dataset.SetWhat(collectionName, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_what will get the "what" value associated with the collection's metadata
//
//export get_what
func get_what(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	src := dataset.GetWhat(collectionName)
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// set_when will set the "when" value associated with the collection's metadata
//
//export set_when
func set_when(cName *C.char, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	src := C.GoString(cSrc)
	error_clear()

	if err := dataset.SetWhen(collectionName, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_when will get the "what" value associated with the collection's metadata
//
//export get_when
func get_when(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	src := dataset.GetWhen(collectionName)
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// set_where will set the "where" value associated with the collection's metadata
//
//export set_where
func set_where(cName *C.char, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	src := C.GoString(cSrc)
	error_clear()

	if err := dataset.SetWhere(collectionName, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_where will get the "where" value associated with the collection's metadata
//
//export get_where
func get_where(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	src := dataset.GetWhere(collectionName)
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// set_version will set the "version" value associated with the collection's metadata
//
//export set_version
func set_version(cName *C.char, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	src := C.GoString(cSrc)
	error_clear()

	if err := dataset.SetVersion(collectionName, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_version will get the "version" value associated with the collection's metadata
//
//export get_version
func get_version(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	src := dataset.GetVersion(collectionName)
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// set_contact will set the "contact" value associated with the collection's metadata
//
//export set_contact
func set_contact(cName *C.char, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	src := C.GoString(cSrc)
	error_clear()

	if err := dataset.SetContact(collectionName, src); err != nil {
		error_dispatch(err, "%s", err)
		return C.int(0)

	}
	return C.int(1)
}

// get_contact will get the "contact" value associated with the collection's metadata
//
//export get_contact
func get_contact(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	src := dataset.GetContact(collectionName)
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

func main() {}
