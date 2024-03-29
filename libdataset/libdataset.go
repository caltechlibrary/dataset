//
// py/dataset.go is a C shared library for implementing a dataset module in Python3
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>

// Copyright (c) 2023, Caltech
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
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v2"
)

var (
	verbose          = false
	useStrictDotpath = true
	// NOTE: error state is shared because C doesn't easily pass multiple
	// return values without resorting to complex structures.
	errorValue error
)

/**
 * "Collections" data structures. Provides the structures needed
 * to embbed dataset in libdataset.so, libdataset.dll and
 * libdataset.dynlib.
 *
 * E.g. Supporting writing an asynchronous web service written in
 * Python via py_dataset needs save writes.
 */

// CMap holds a map of collection names to *dataset.Collection
type CMap struct {
	collections map[string]*dataset.Collection
}

var (
	cMap *CMap
)

// IsOpen checks to see if a dataset collection is already opened.
func IsOpen(cName string) bool {
	if cMap == nil {
		return false
	}
	if _, exists := cMap.collections[cName]; exists == true &&
		cMap.collections[cName] != nil &&
		cMap.collections[cName].Name == path.Base(cName) {
		return true
	}
	return false
}

// Open opens a dataset collection for use in a service like
// context. CMap collections remain "open" until explicitly closed
// or closed via CloseAll().
// Writes to the collections are run through a mutex
// to prevent collisions. Subsequent CMapOpen() will open
// additional collections under the the service.
func OpenCollection(cName string) error {
	var (
		err error
	)
	if cMap == nil {
		cMap = new(CMap)
		cMap.collections = make(map[string]*dataset.Collection)
	}
	if _, exists := cMap.collections[cName]; exists == true {
		return fmt.Errorf("%q opened previously", cName)
	}
	c, err := dataset.Open(cName)
	if err != nil {
		return fmt.Errorf("%q failed to open, %s", cName, err)
	}
	cMap.collections[cName] = c
	return nil
}

// InitCollection takes a collection name and initializes it.
func InitCollection(cName string, storageType string) error {
	if cMap == nil {
		cMap = new(CMap)
		cMap.collections = make(map[string]*dataset.Collection)
	}
	if _, exists := cMap.collections[cName]; exists == true {
		return fmt.Errorf("%q opened previously", cName)
	}
	c, err := dataset.Init(cName, storageType)
	if err != nil {
		return fmt.Errorf("%q failed to initalize, %s", cName, err)
	}
	cMap.collections[cName] = c
	return nil
}

// GetCollection takes a collection name, opens it if necessary and returns a handle
// to the CMapCollection struct and error value.
func GetCollection(cName string) (*dataset.Collection, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c, nil
	}
	return nil, fmt.Errorf("%s not found", cName)
}

// Collections returns a list of collections previously
// opened with CMapOpen()
func Collections() []string {
	cNames := []string{}
	if cMap == nil {
		return cNames
	}
	for cName, c := range cMap.collections {
		if c != nil && path.Base(cName) == c.Name {
			cNames = append(cNames, cName)
		}
	}
	return cNames
}

// Close closes a dataset collections previously
// opened by CMapOpen().  It will also set the internal
// cMap variable to nil if there are no remaining collections.
func Close(cName string) error {
	if IsOpen(cName) {
		if c, exists := cMap.collections[cName]; exists == true {
			err := c.Close()
			delete(cMap.collections, cName)
			return err
		}
	}
	return fmt.Errorf("%q not found", cName)
}

// CloseAll goes through the service collection list
// and closes each one.
func CloseAll() error {
	if cMap == nil {
		return fmt.Errorf("Nothing to close")
	}
	errors := []string{}
	for cName, c := range cMap.collections {
		if c != nil {
			if err := c.Close(); err != nil {
				errors = append(errors, fmt.Sprintf("%q %s", cName, err))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

// Keys returns a list of keys for a collection opened with
// StartCMap.
func Keys(cName string) []string {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return []string{}
		}
	}
	if c, found := cMap.collections[cName]; found == true && c != nil {
		keys, _ := c.Keys()
		return keys
	}
	return []string{}
}

// HasKey returns true if the key exists in the collection or false otherwise
func HasKey(cName string, key string) bool {
	c, err := GetCollection(cName)
	if err != nil {
		return false
	}
	return c.HasKey(key)
}

// CreateJSON takes a collection name, key and JSON object
// document and creates a new JSON object in the collection using
// the key.
func CreateJSON(cName string, key string, src []byte) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		err := c.CreateJSON(key, src)
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// ReadJSON takes a collection name, key and returns a JSON object
// document.
func ReadJSON(cName string, key string) ([]byte, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.ReadJSON(key)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// ReadJSONVersion takes a collection name, key, semver and 
// returns a JSON object document.
func ReadJSONVersion(cName string, key string, semver string) ([]byte, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.ReadJSONVersion(key, semver)
	}
	return nil, fmt.Errorf("%q not available", cName)
}


// UpdateJSON takes a collection name, key and JSON object
// document and updates the collection.
func UpdateJSON(cName string, key string, src []byte) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		err := c.UpdateJSON(key, src)
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// DeleteJSON takes a collection name and key and removes
// and JSON object from the collection.
func DeleteJSON(cName string, key string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		err := c.Delete(key)
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// HasFrame returns true if frame found in service collection,
// otherwise false
func HasFrame(cName string, fName string) bool {
	// We may need to open a dataset collection to check for a frame.
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return false
		}
	}
	if IsOpen(cName) == true {
		if c, found := cMap.collections[cName]; found {
			return c.HasFrame(fName)
		}
	}
	return false
}

// FrameKeys returns the ordered list of keys for the frame.
func FrameKeys(cName string, fName string) []string {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil
		}
	}
	if c, found := cMap.collections[cName]; found {
		f, err := c.FrameRead(fName)
		if err != nil {
			return nil
		}
		return f.Keys
	}
	return nil
}

// FrameCreate creates a frame in a service collection
func FrameCreate(cName string, fName string, keys []string, dotPaths []string, labels []string, verbose bool) (*dataset.DataFrame, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		f, err := c.FrameCreate(fName, keys, dotPaths, labels, verbose)
		return f, err
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// FrameObjects returns a JSON document of a copy of the objects in a
// frame for the service collection. It is analogous to a
// dataset.ReadJSON but for a frame's object list
func FrameObjects(cName string, fName string) ([]map[string]interface{}, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.FrameObjects(fName)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// FrameRefresh updates the frame object list's for the keys provided.
// Any new keys
//
// cause a new object to be appended to the end of the list.
func FrameRefresh(cName string, fName string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.Close()
		return c.FrameRefresh(fName, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameReframe updates the frame object list. If a list of keys is
// provided then the object will be replaced with updated objects based
// on the keys provided.
func FrameReframe(cName string, fName string, keys []string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.FrameReframe(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameClear clears the object and key list from a frame
func FrameClear(cName string, fName string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.FrameClear(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameDelete deletes a frame from a service collection
func FrameDelete(cName string, fName string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := OpenCollection(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.FrameDelete(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameNames returns a list of frame names in a service collection
func FrameNames(cName string) []string {
	if IsOpen(cName) == true {
		if c, found := cMap.collections[cName]; found {
			return c.FrameNames()
		}
	}
	return []string{}
}

// Check checks a dataset collection and reports error to console.
// NOTE: Collection is closed and objects are locked during check!
func Check(cName string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		return dataset.Analyzer(cName, verbose)
	}
	if _, found := cMap.collections[cName]; found {
		Close(cName)
		err := dataset.Analyzer(cName, verbose)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not found", cName)
}

// Repair repairs a collection
// NOTE: Collection objects are locked during repair!
func Repair(cName string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		return dataset.Repair(cName, verbose)
	}
	if _, found := cMap.collections[cName]; found {
		Close(cName)
		err := dataset.Repair(cName, verbose)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s not found", cName)
}

// error_clear will set the global error state to nil.
//
//export error_clear
func error_clear() {
	errorValue = nil
}

// errorDispatch logs error messages to console based on string template
// Not exported.
func errorDispatch(err error, s string, values ...interface{}) {
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

// dataset_version returns the version of libdataset.
//
//export dataset_version
func dataset_version() *C.char {
	txt := dataset.Version
	return C.CString(fmt.Sprintf("%q", txt))
}

/*
 * collection operations
 */

// init_collection intializes a collection and records as much metadata
// as it can from the execution environment (e.g. username,
// datetime created). NOTE: New parameter required, storageType. This
// can be either "pairtree" or "sqlstore".
//
//export init_collection
func init_collection(name *C.char, cStorageType *C.char) C.int {
	storageType := C.GoString(cStorageType)
	collectionName := C.GoString(name)
	if verbose == true {
		messagef("creating %s\n", collectionName)
	}
	error_clear()
	err := InitCollection(collectionName, storageType)
	if err != nil {
		errorDispatch(err, "Cannot create collection %s, %s", collectionName, err)
		return C.int(0)
	}
	messagef("%s initialized", collectionName)
	return C.int(1)
}

// is_collection_open returns true (i.e. one) if a collection has been opened by libdataset, false (i.e. zero) otherwise
//
//export is_collection_open
func is_collection_open(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	if IsOpen(collectionName) {
		return C.int(1)
	}
	return C.int(0)
}

// open_collection returns 0 on successfully opening a collection 1 otherwise. Sets error messages if needed.
//
//export open_collection
func open_collection(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	err := OpenCollection(collectionName)
	if err != nil {
		errorDispatch(err, "Cannot open %q, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// collections returns a JSON list of collection names that are open otherwise an empty list.
//
//export collections
func collections() *C.char {
	cNames := Collections()
	src, err := dataset.JSONMarshal(cNames)
	if err != nil {
		return C.CString("[]")
	}
	return C.CString(fmt.Sprintf("%s", src))
}

// close_collection closes a collection previously opened.
//
//export close_collection
func close_collection(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if err := Close(collectionName); err != nil {
		errorDispatch(err, "Cannot close collection %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// close_all_collections closes all collections previously opened
//
//export close_all_collections
func close_all_collections() C.int {
	error_clear()
	if err := CloseAll(); err != nil {
		errorDispatch(err, "Cannot close all collections, %s", err)
		return C.int(0)
	}
	return C.int(1)
}

// collection_exits checks to see if a collection exists or not.
//
//export collection_exists
func collection_exists(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "failed, %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// check_collection runs the analyzer over a collection and looks for
// problem records.
//
//export check_collection
func check_collection(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := Check(collectionName, verbose)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// repair_collection runs the analyzer over a collection and repairs JSON
// objects and attachment discovered having a problem. Also is
// useful for upgrading a collection between dataset releases.
//
//export repair_collection
func repair_collection(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	err := Repair(collectionName, verbose)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// clone_collection takes a collection name, a JSON array of keys and creates
// a new collection with a new name based on the origin's collections'
// objects. NOTE: If you are using pairtree dsn can be an empty string
// otherwise it needs to be a dsn to connect to the SQL store.
//
//export clone_collection
func clone_collection(cName *C.char, cDsn *C.char, cKeys *C.char, dName *C.char) C.int {
	collectionName := C.GoString(cName)
	srcKeys := C.GoString(cKeys)
	destName := C.GoString(dName)
	dsn := C.GoString(cDsn)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	keys := []string{}
	err = json.Unmarshal([]byte(srcKeys), &keys)
	if err != nil {
		errorDispatch(err, "Can't unmarshal keys, %s", err)
		return C.int(0)
	}
	err = c.Clone(destName, dsn, keys, verbose)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// clone_sample is like clone both generates a sample or test and
// training set of sampled of the cloned collection. NOTE: The
// training name and testing name are followed by their own dsn values.
// If the dsn is an empty string then a pairtree store is assumed.
//
//export clone_sample
func clone_sample(cName *C.char, cTrainingName *C.char, cTrainingDsn *C.char, cTestName *C.char, cTestDsn *C.char, cSampleSize C.int) C.int {
	collectionName := C.GoString(cName)
	sampleSize := int(cSampleSize)
	trainingName := C.GoString(cTrainingName)
	trainingDsn := C.GoString(cTrainingDsn)
	testName := C.GoString(cTestName)
	testDsn := C.GoString(cTestDsn)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	keys, _ := c.Keys()
	err = c.CloneSample(trainingName, trainingDsn, testName, testDsn, keys, sampleSize, verbose)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// import_csv - import a CSV file into a collection
// syntax: COLLECTION CSV_FILENAME ID_COL
//
// options that should support sensible defaults:
//
//	cUseHeaderRow
//	cOverwrite
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
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if idCol < 1 {
		errorDispatch(fmt.Errorf("invalid column number"), "Column number must be greater than zero, got %d", idCol)
		return C.int(0)
	}

	// NOTE: we need to adjust to zero based index
	idCol--

	// Now import our CSV file
	fp, err := os.Open(csvFName)
	if err != nil {
		errorDispatch(err, "Can't open %s, %s", csvFName, err)
		return C.int(0)
	}
	cnt, err := c.ImportCSV(fp, idCol, useHeaderRow, overwrite, verbose)
	if err != nil {
		errorDispatch(err, "%s\n", err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)

	return C.int(1)
}

// export_csv - export collection objects as a frame to a CSV file
// syntax: COLLECTION FRAME CSV_FILENAME
//
//export export_csv
func export_csv(cName *C.char, cFrameName *C.char, cCSVFName *C.char) C.int {
	// Convert out parameters
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFrameName)
	csvFName := C.GoString(cCSVFName)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "Can't open %s, %s", collectionName, err)
		return C.int(0)
	}

	fp, err := os.Create(csvFName)
	if err != nil {
		errorDispatch(err, "Can't create %s, %s", csvFName, err)
		return C.int(0)
	}
	defer fp.Close()

	// Get Frame
	if c.HasFrame(frameName) == false {
		errorDispatch(err, "Missing frame %q in %s\n", frameName, collectionName)
		return C.int(0)
	}

	// Get dotpaths and column labels from frame
	f, err := c.FrameRead(frameName)
	if err != nil {
		errorDispatch(err, "%s\n", err)
		return C.int(0)
	}

	// Now export to CSV
	cnt, err := c.ExportCSV(fp, os.Stderr, f, verbose)
	if err != nil {
		errorDispatch(err, "Can't export CSV %s, %s", csvFName, err)
		return C.int(0)
	}
	messagef("%d total rows processed", cnt)
	return C.int(1)
}

// sync_send_csv - synchronize a frame sending data to a CSV file
// returns 1 (True) on success, 0 (False) otherwise.
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
		return C.int(0)
	}
	if len(src) == 0 {
		fmt.Fprintf(os.Stderr, "No data in csv file %s\n", csvFilename)
		return C.int(0)
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err == nil {
			table = dataset.TableStringToInterface(csvTable)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(0)
		}
	}

	error_clear()
	c, err = GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)

	}

	// Merge collection content into table
	table, err = c.MergeIntoTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(0)
	}

	// Save the resulting table
	if len(src) > 0 {
		if err = os.Rename(csvFilename, csvFilename+".bak"); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(0)
		}
		out, err := os.Create(csvFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(0)
		}
		w := csv.NewWriter(out)
		w.WriteAll(dataset.TableInterfaceToString(table))
		if err = w.Error(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(0)
		}
	}
	return C.int(1)
}

// sync_recieve_csv - synchronize a frame recieving data from a CSV file
// returns 1 (True) on success, 0 (False) otherwise.
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
		return C.int(0)
	}

	table := [][]interface{}{}
	// Populate table to sync
	if len(src) > 0 {
		// for CSV
		r := csv.NewReader(bytes.NewReader(src))
		csvTable, err := r.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return C.int(0)
		}
		table = dataset.TableStringToInterface(csvTable)
	}

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	// Merge table contents into Collection and Frame
	err = c.MergeFromTable(frameName, table, syncOverwrite, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.int(0)
	}
	return C.int(1)
}

/*
 * Key operations
 */

// has_key returns 1 if the key exists in collection or 0 if not.
//
//export has_key
func has_key(cName, cKey *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	if HasKey(collectionName, key) {
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
	keyList := Keys(collectionName)
	src, err := dataset.JSONMarshal(keyList)
	if err != nil {
		errorDispatch(err, "Can't marshal key list, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

/*
 * Object operations
 */

// create_object takes JSON source and adds it to the collection with
// the provided key.
//
//export create_object
func create_object(cName, cKey, cSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	src := []byte(C.GoString(cSrc))

	error_clear()
	err := CreateJSON(collectionName, key, src)
	if err != nil {
		errorDispatch(err, "Create %s failed, %s", key, err)
		return C.int(0)
	}
	return C.int(1)
}

// read_object takes a key and returns JSON source of the record
//
//export read_object
func read_object(cName, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	error_clear()
	var (
		src []byte
		err error
	)
	if src, err = ReadJSON(collectionName, key); err != nil {
		errorDispatch(err, "Can't read %s, %s", key, err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// read_object_version takes a collection name, key and semver
// and returns JSON source of the record version.
//
//export read_object_version
func read_object_version(cName, cKey *C.char, cSemver *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	semver := C.GoString(cSemver)

	error_clear()
	var (
		src []byte
		err error
	)
	src, err = ReadJSONVersion(collectionName, key, semver)
	if err != nil {
		errorDispatch(err, "Can't read %s %s, %s", key, semver, err)
		return C.CString("")
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
func read_object_list(cName *C.char, cKeysAsJSON *C.char) *C.char {
	collectionName := C.GoString(cName)
	l := []string{}
	errList := []string{}

	error_clear()
	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(cKeysAsJSON))
	keyList := []string{}
	err := json.Unmarshal(src, &keyList)
	if err != nil {
		errorDispatch(err, "Can't unmarshal key list, %s", err)
		return C.CString("")
	}
	if _, err = GetCollection(collectionName); err != nil {
		errorDispatch(err, "Cannot open collection %s, %s", collectionName, err)
		return C.CString("")
	}

	for _, key := range keyList {
		src, err := ReadJSON(collectionName, key)
		if err != nil {
			errList = append(errList, fmt.Sprintf("(%s) %s", key, err))
			continue
		}
		l = append(l, fmt.Sprintf("%s", src))
	}
	if len(errList) > 0 {
		err = fmt.Errorf("%s", strings.Join(errList, "; "))
		errorDispatch(err, "Key read errors %s", err)
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
	err := UpdateJSON(collectionName, key, src)
	if err != nil {
		errorDispatch(err, "Update %s failed, %s", key, err)
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
	err := DeleteJSON(collectionName, key)
	if err != nil {
		errorDispatch(err, "Cannot open collection %s, %s", collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// join_objects takes a collection name, a key, and merges JSON source with an
// existing JSON record. If overwrite is 1 it overwrites and replaces
// common values, if not 1 it only adds missing attributes.
//
//export join_objects
func join_objects(cName *C.char, cKey *C.char, cObjSrc *C.char, cOverwrite C.int) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	overwrite := (cOverwrite == 1)
	objectSrc := C.GoString(cObjSrc)

	error_clear()
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	newObject := map[string]interface{}{}
	if err := dataset.JSONUnmarshal([]byte(objectSrc), &newObject); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	src, err := ReadJSON(collectionName, key)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	outObject := map[string]interface{}{}
	if err := dataset.JSONUnmarshal(src, &outObject); err != nil {
		errorDispatch(err, "%s", err)
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
	src, err = dataset.JSONMarshal(outObject)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	if err := UpdateJSON(collectionName, key, src); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// count_objects returns the number of objects (records) in a collection.
// if an error is encounter a -1 is returned.
//
//export count_objects
func count_objects(cName *C.char) C.int {
	collectionName := C.GoString(cName)
	error_clear()
	if c, err := GetCollection(collectionName); err == nil {
		i := c.Length()
		return C.int(i)
	}
	return C.int(-1)
}

// object_path returns the path on disc to an JSON object document
// in the collection.
//
//export object_path
func object_path(cName *C.char, cKey *C.char) *C.char {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "Can't find %q", collectionName)
		return C.CString("")
	}

	s, err := c.DocPath(key)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	return C.CString(s)
}

// create_objects - is a function to creates empty a objects in batch.
// It requires a JSON list of keys to create. For each key present
// an attempt is made to create a new empty object based on the JSON
// provided (e.g. `{}`, `{"is_empty": true}`). The reason to do this
// is that it means the collection.json file is updated once for the
// whole call and that the keys are now reserved to be updated separately.
// Returns 1 on success, 0 if errors encountered.
//
//export create_objects
func create_objects(cName *C.char, keysAsJSON *C.char, objectAsJSON *C.char) C.int {
	collectionName := C.GoString(cName)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(keysAsJSON))
	keyList := []string{}
	err = json.Unmarshal(src, &keyList)
	if err != nil {
		errorDispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}
	objectSrc := []byte(C.GoString(objectAsJSON))

	err = c.CreateObjectsJSON(keyList, objectSrc)
	if err != nil {
		errorDispatch(err, "Create objects failed, %s", err)
		return C.int(0)
	}
	return C.int(1)
}

// update_objects - is a function to update objects in batch.
// It requires a JSON array of keys and a JSON array of
// matching objects. The list of keys and objects are processed
// together with calls to update individual records. Returns 1 on
// success, 0 on error.
//
//export update_objects
func update_objects(cName *C.char, keysAsJSON *C.char, objectsAsJSON *C.char) C.int {
	collectionName := C.GoString(cName)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	// Now unpack our keys into an array of strings.
	src := []byte(C.GoString(keysAsJSON))
	keyList := []string{}
	err = json.Unmarshal(src, &keyList)
	if err != nil {
		errorDispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}
	src = []byte(C.GoString(objectsAsJSON))
	objectList := []map[string]interface{}{}
	err = json.Unmarshal(src, &objectList)
	if err != nil {
		errorDispatch(err, "Can't unmarshal key list, %s", err)
		return C.int(0)
	}

	if len(keyList) != len(objectList) {
		errorDispatch(err, "expected %d keys for %d objects", len(keyList), len(objectList))
		return C.int(0)
	}

	errorNo := 1
	for i, key := range keyList {
		if c.HasKey(key) {
			err = c.Update(key, objectList[i])
			if err != nil {
				errorDispatch(err, "Can't update key %q, %s", key, err)
				errorNo = 0
			}
		}
	}
	return C.int(errorNo)
}

// list_objects returns JSON array of objects in a collections based on a
// JSON array of keys.
//
//export list_objects
func list_objects(cName *C.char, cKeys *C.char) *C.char {
	collectionName := C.GoString(cName)
	sKeys := C.GoString(cKeys)

	error_clear()
	keys := []string{}
	err := json.Unmarshal([]byte(sKeys), &keys)
	if err != nil {
		errorDispatch(err, "Failed to unmarshal key list, %s", err)
		return C.CString("")
	}

	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "Can't find %q", collectionName)
		return C.CString("")
	}

	recs := []map[string]interface{}{}
	for _, name := range keys {
		m := map[string]interface{}{}
		err = c.Read(name, m)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.CString("")
		}
		recs = append(recs, m)
	}
	src, err := dataset.JSONMarshal(recs)
	if err != nil {
		errorDispatch(err, "failed to marshal result, %s", err)
		return C.CString("")
	}
	return C.CString(string(src))
}

/*
 * Attachment operations
 */

// attach will attach a file to a JSON object in a collection. If the
// collection is versioned then the semver will be managed automatically.
//
//export attach
func attach(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			errorDispatch(err, "Can't unmarshal %q, %s", srcFNames, err)
			return C.int(0)
		}
	}
	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if c.HasKey(key) == false {
		errorDispatch(fmt.Errorf("missing key"), "%q is not in collection", key)
		return C.int(0)
	}
	// NOTE: fName, fNames are file names NOT frame names
	for _, fName := range fNames {
		if _, err := os.Stat(fName); os.IsNotExist(err) {
			errorDispatch(err, "%s does not exist", fName)
			return C.int(0)
		}
		err = c.AttachFile(key, fName)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.int(0)
		}
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
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	if c.HasKey(key) == false {
		errorDispatch(fmt.Errorf("missing key"), "%q is not in collection", key)
		return C.CString("")
	}
	results, err := c.Attachments(key)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	if len(results) > 0 {
		return C.CString(strings.Join(results, "\n"))
	}
	return C.CString("")
}

// detach exports the file associated with the key and basenames 
// the JSON object in the collection. The file remains "attached".
//
//export detach
func detach(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			errorDispatch(err, "Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	} 

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if len(fNames) == 0 {
		// IF the list is empty detach expected to detached all files.
		results, err := c.Attachments(key)
		if err != nil {
			errorDispatch(err, "Can't determine attached files, %s",err)
			return C.int(0)
		}
		for _, fName := range results {
			fNames = append(fNames, fName)
		}
	}

	if c.HasKey(key) == false {
		errorDispatch(err, "%q is not in collection", key)
		return C.int(0)
	}
	for _, fName := range fNames {
		var src []byte
		src, err = c.RetrieveFile(key, fName)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.int(0)
		}
		err = os.WriteFile(fName, src, 0666)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.int(0)
		}
	}
	return C.int(1)
}

// detach_version exports the file associated with the semver from the JSON
// object in the collection. The file remains "attached".
//
//export detach_version
func detach_version(cName *C.char, cKey *C.char, cSemver *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	semver := C.GoString(cSemver)

	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			errorDispatch(err, "Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	} 

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}

	if len(fNames) == 0 {
		// IF the list is empty detach expected to detached all files.
		results, err := c.Attachments(key)
		if err != nil {
			errorDispatch(err, "Can't determine attached files, %s",err)
			return C.int(0)
		}
		for _, fName := range results {
			fNames = append(fNames, fName)
		}
	}

	if c.HasKey(key) == false {
		errorDispatch(err, "%q is not in collection", key)
		return C.int(0)
	}
	for _, fName := range fNames {
		var src []byte
		src, err = c.RetrieveVersionFile(key, fName, semver)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.int(0)
		}
		err = os.WriteFile(fName, src, 0666)
		if err != nil {
			errorDispatch(err, "%s", err)
			return C.int(0)
		}
	}
	return C.int(1)
}


// prune removes an attachment by basename from a JSON object in the
// collection. This is destructive, the file is removed from disc. 
// NOTE: If the collection is versioned prune removes ALL versions!!!
//
//export prune
func prune(cName *C.char, cKey *C.char, cFNames *C.char) C.int {
	collectionName := C.GoString(cName)
	key := C.GoString(cKey)
	srcFNames := C.GoString(cFNames)
	fNames := []string{}
	if len(srcFNames) > 0 {
		err := json.Unmarshal([]byte(srcFNames), &fNames)
		if err != nil {
			errorDispatch(err, "Can't unmarshal filename list, %s", err)
			return C.int(0)
		}
	}

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.int(0)
	}
	if len(fNames) == 0 {
		// IF the list is empty prune all files.
		if err := c.PruneAll(key); err != nil {
			errorDispatch(err, "prune all files failed, %s", err)
			return C.int(0)
		}
		return C.int(1)
	} 
	for i, fName := range fNames {
		err = c.Prune(key, fName)
		if err != nil {
			errorDispatch(err, "%s (%d) %s", fName, i, err)
			return C.int(0)
		}
	}
	return C.int(1)
}

/*
 * Frame operations
 */

// frame retrieves a frame including its metadata. NOTE:
// if you just want the object list, use frame_objects().
//
//export frame
func frame(cName *C.char, cFName *C.char) *C.char {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)

	error_clear()
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	f, err := c.FrameRead(frameName)
	if err != nil {
		errorDispatch(err, "failed to create frame, %s", err)
		return C.CString("")
	}
	src, err := dataset.JSONMarshal(f)
	if err != nil {
		errorDispatch(err, "failed to marshal frame, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}


// has_frame returns 1 (true) if frame name exists in collection, 0 (false) otherwise
//
//export has_frame
func has_frame(cName *C.char, cFName *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	if HasFrame(collectionName, frameName) {
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
	keys := FrameKeys(collectionName, frameName)
	src, err := dataset.JSONMarshal(keys)
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
	dotPathsSrc := []byte(C.GoString(cDotPathsSrc))
	labelsSrc := []byte(C.GoString(cLabelsSrc))
	keys := []string{}
	dotPaths := []string{}
	labels := []string{}

	error_clear()
	if err := json.Unmarshal(keysSrc, &keys); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	if err := json.Unmarshal(dotPathsSrc, &dotPaths); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	if err := json.Unmarshal(labelsSrc, &labels); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	_, err := FrameCreate(collectionName, frameName, keys, dotPaths, labels, verbose)
	if err != nil {
		errorDispatch(err, "%s", err)
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
	ol, err := FrameObjects(collectionName, frameName)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	src, err := dataset.JSONMarshal(ol)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

// frame_refresh refresh the contents of the frame using the
// existing keys associated with the frame and the current state
// of the collection.  NOTE: If a key is missing
// in the collection then the key and object is removed.
//
//export frame_refresh
func frame_refresh(cName *C.char, cFName *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)

	error_clear()
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}

	if err := FrameRefresh(collectionName, frameName, verbose); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}

// frame_reframe will change the key and object list in a frame based on
// the key list provided and the current state of the collection.
//
//export frame_reframe
func frame_reframe(cName *C.char, cFName *C.char, cKeysSrc *C.char) C.int {
	collectionName := C.GoString(cName)
	frameName := C.GoString(cFName)
	keysSrc := C.GoString(cKeysSrc)

	error_clear()
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	keys := []string{}
	if err := json.Unmarshal([]byte(keysSrc), &keys); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	//NOTE: We're picking up the verbose flag from the modules global state
	if err := FrameReframe(collectionName, frameName, keys, verbose); err != nil {
		errorDispatch(err, "%s", err)
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
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	if err := FrameClear(collectionName, frameName); err != nil {
		errorDispatch(err, "%s", err)
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
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	if err := FrameDelete(collectionName, frameName); err != nil {
		errorDispatch(err, "%s", err)
		return C.int(0)
	}
	return C.int(1)
}


// frame_names returns a JSON array of frames names in the collection.
//
//export frame_names
func frame_names(cName *C.char) *C.char {
	collectionName := C.GoString(cName)

	error_clear()
	if _, err := GetCollection(collectionName); err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}

	frameNames := FrameNames(collectionName)
	if len(frameNames) == 0 {
		return C.CString("[]")
	}
	src, err := dataset.JSONMarshal(frameNames)
	if err != nil {
		errorDispatch(err, "failed to marshal frame names, %s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
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
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "%q not found", collectionName)
		return C.CString("")
	}

	f, err := c.FrameRead(frameName)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	g := f.Grid(includeHeaderRow)
	src, err := dataset.JSONMarshal(g)
	if err != nil {
		errorDispatch(err, "%s", err)
		return C.CString("")
	}
	txt := fmt.Sprintf("%s", src)
	return C.CString(txt)
}

/*
 * metadata for collection
 */

// NOTE: Namaste methods removed, v2 of dataset supports a codemeta.json file for metadata about the collection. 2023-02-09

// get_codemeta returns any metadata associated with the collection as
// a codemeta JSON document
//
// export get_codemeta
func get_codemeta(cName *C.char) *C.char {
	collectionName := C.GoString(cName)
	txt := ""
	c, err := GetCollection(collectionName)
	if err == nil {
		if src, err := c.Codemeta(); err == nil {
			txt = fmt.Sprintf("%s", src)
		}
	}
	return C.CString(txt)
}

// set_versioning sets the versioning on a collection. versioning value
// can be "", "none",  "patch", /"minor", "major".
//
//export set_versioning
func set_versioning(cName *C.char, cVersioning *C.char) C.int {
	collectionName := C.GoString(cName)
	versioning := C.GoString(cVersioning)
	c, err := GetCollection(collectionName)
	if err != nil {
		errorDispatch(err, "error getting %q, %s", collectionName, err)
		return C.int(0)
	}
	switch versioning {
		case "":
			c.Versioning = ""
		case "none":
			c.Versioning = ""
		case "major":
			c.Versioning = "major"
		case "minor":
			c.Versioning = "minor"
		case "patch":
			c.Versioning = "patch"
		default:
			errorDispatch(err, "%q is not a valid versioning method", versioning)
			return C.int(0)
	}
	return C.int(1)
}

// get_versioning will returns the versioning setting (e.g. "", "patch",
// "minor", "major") on a collection.
//
//export get_versioning
func get_versioning(cName *C.char) *C.char {
	collectionName := C.GoString(cName)
	txt := ""
	c, err := GetCollection(collectionName)
	if err == nil {
		txt = c.Versioning
	}
	return C.CString(txt)
}

func main() {}
