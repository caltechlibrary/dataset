//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
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
package dataset

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/namaste"
	"github.com/caltechlibrary/shuffle"
	"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"
)

const (
	// Version of the dataset package
	Version = `v0.0.45`

	// License is a formatted from for dataset package based command line tools
	License = `
%s %s

Copyright (c) 2018, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Sort directions
	ASC  = iota
	DESC = iota
)

// Supported file layout types
const (
	// Assume an unknown layout is zero, then add consts in order of adoption
	UNKNOWN_LAYOUT = iota

	// Buckets is the first file layout implemented when dataset started
	BUCKETS_LAYOUT = iota

	// Pairtree is the perferred file layout moving forward
	PAIRTREE_LAYOUT = iota
)

// Collection is the container holding buckets which in turn hold JSON docs
type Collection struct {
	// Version of collection being stored
	Version string `json:"version"`

	// Name of collection
	Name string `json:"name"`

	// Type allows for transitioning from bucket layout to pairtree layout for collections.
	Layout int `json:"layout"`

	// Buckets is a list of bucket names used by collection (depreciated, will be removed after migration to pairtree)
	Buckets []string `json:"buckets,omitempty"`

	// KeyMap holds the document key to path in the collection
	KeyMap map[string]string `json:"keymap"`

	// Store holds the storage system information (e.g. local disc, S3, GS)
	// and related methods for interacting with it
	Store *storage.Store `json:"-"`

	// FrameMap is a list of frame names and with rel path to the frame defined in the collection
	FrameMap map[string]string `json:"frames"`
}

//
// internal utility functions
//

// normalizeKeyName() trims leading and trailing spaces
func normalizeKeyName(s string) string {
	return strings.TrimSpace(s)
}

// collectionNameFromPath takes a path and normalized a collection name.
func collectionNameFromPath(p string) string {
	if strings.Contains(p, "://") {
		u, _ := url.Parse(p)
		return path.Base(u.Path)
	}
	return strings.TrimSpace(p)
}

// keyAndFName converts a key (which may have things like slashes) into a disc friendly name and key value
func keyAndFName(name string) (string, string) {
	var keyName string
	if strings.HasSuffix(name, ".json") == true {
		return keyName, name
	}
	return name, url.QueryEscape(name) + ".json"
}

// saveMetadata writes the collection's metadata to COLLECTION_NAME/collection.json
func (c *Collection) saveMetadata() error {
	// Check to see if collection exists, if not create it!
	if _, err := c.Store.Stat(c.Name); err != nil {
		if err := c.Store.MkdirAll(c.Name, 0775); err != nil {
			return err
		}
	}
	src, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Can't marshal metadata, %s", err)
	}
	if err := c.Store.WriteFile(path.Join(c.Name, "collection.json"), src, 0664); err != nil {
		return fmt.Errorf("Can't store collection metadata, %s", err)
	}
	return nil
}

//
// Public interface for dataset
//

// InitCollection - creates a new collection with default alphabet and names of length 2.
// NOTE: layoutType is provided to allow for future changes in the file layout of a collection.
func InitCollection(name string, layoutType int) (*Collection, error) {
	var (
		c   *Collection
		err error
	)
	switch layoutType {
	case PAIRTREE_LAYOUT:
		c, err = pairtreeCreateCollection(name)
	case BUCKETS_LAYOUT:
		c, err = bucketCreateCollection(name, DefaultBucketNames)
	default:
		c, err = bucketCreateCollection(name, DefaultBucketNames)
	}
	if err != nil {
		return nil, err
	}
	// Add Namaste type record
	namaste.DirType(name, fmt.Sprintf("dataset_%s", Version[1:]))
	namaste.When(name, time.Now().Format("2006-01-02"))
	return c, nil
}

// Open reads in a collection's metadata and returns and new collection structure and err
func Open(name string) (*Collection, error) {
	store, err := storage.GetStore(name)
	if err != nil {
		return nil, err
	}
	collectionName := collectionNameFromPath(name)
	src, err := store.ReadFile(path.Join(collectionName, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err != nil {
		return nil, err
	}
	//NOTE: we need to reset collectionName so we're working with a path useable to get to the JSON documents.
	c.Name = collectionName
	c.Store = store
	return c, nil
}

// Delete an entire collection
func Delete(name string) error {
	store, err := storage.GetStore(name)
	if err != nil {
		return err
	}
	collectionName := collectionNameFromPath(name)
	if err := store.RemoveAll(collectionName); err != nil {
		return err
	}
	return nil
}

// DocPath returns a full path to a key or an error if not found
func (c *Collection) DocPath(name string) (string, error) {
	keyName, name := keyAndFName(name)
	if p, ok := c.KeyMap[keyName]; ok == true {
		return path.Join(c.Name, p, name), nil
	}
	return "", fmt.Errorf("Can't find %q", name)
}

// Close closes a collection, writing the updated keys to disc
func (c *Collection) Close() error {
	// Cleanup c so it can't accidentally get reused
	c.Buckets = []string{}
	c.Name = ""
	c.KeyMap = map[string]string{}
	c.Store = nil
	return nil
}

// CreateJSON adds a JSON doc to a collection, if a problem occurs it returns an error
func (c *Collection) CreateJSON(key string, src []byte) error {
	switch c.Layout {
	case PAIRTREE_LAYOUT:
		return c.pairtreeCreateJSON(key, src)
	case BUCKETS_LAYOUT:
		return c.bucketCreateJSON(key, src)
	default:
		return c.bucketCreateJSON(key, src)
	}
}

// ReadJSON finds a the record in the collection and returns the JSON source
func (c *Collection) ReadJSON(name string) ([]byte, error) {
	switch c.Layout {
	case PAIRTREE_LAYOUT:
		return c.pairtreeReadJSON(name)
	case BUCKETS_LAYOUT:
		return c.bucketReadJSON(name)
	default:
		return c.bucketReadJSON(name)
	}
}

// UpdateJSON a JSON doc in a collection, returns an error if there is a problem
func (c *Collection) UpdateJSON(name string, src []byte) error {
	switch c.Layout {
	case PAIRTREE_LAYOUT:
		return c.pairtreeUpdateJSON(name, src)
	case BUCKETS_LAYOUT:
		return c.bucketUpdateJSON(name, src)
	default:
		return c.bucketUpdateJSON(name, src)
	}
}

// Create a JSON doc from an map[string]interface{} and adds it  to a collection, if problem returns an error
// name must be unique. Document must be an JSON object (not an array).
func (c *Collection) Create(name string, data map[string]interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s, %s", name, err)
	}
	return c.CreateJSON(name, src)
}

// Read finds the record in a collection, updates the data interface provide and if problem returns an error
// name must exist or an error is returned
func (c *Collection) Read(name string, data map[string]interface{}) error {
	src, err := c.ReadJSON(name)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	return nil
}

// Update JSON doc in a collection from the provided data interface (note: JSON doc must exist or returns an error )
func (c *Collection) Update(name string, data map[string]interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Update can't marshal into JSON %s, %s", name, err)
	}
	return c.UpdateJSON(name, src)
}

// Delete removes a JSON doc from a collection
func (c *Collection) Delete(name string) error {
	switch c.Layout {
	case PAIRTREE_LAYOUT:
		return c.pairtreeDelete(name)
	case BUCKETS_LAYOUT:
		return c.bucketDelete(name)
	default:
		return c.bucketDelete(name)
	}
}

// Keys returns a list of keys in a collection
func (c *Collection) Keys() []string {
	keys := []string{}
	for k := range c.KeyMap {
		keys = append(keys, k)
	}
	return keys
}

// HasKey returns true if key is in collection's KeyMap, false otherwise
func (c *Collection) HasKey(key string) bool {
	_, hasKey := c.KeyMap[key]
	//FIXME: if pairtree then we can also check by calculating the path and checking the storage system.
	return hasKey
}

// Length returns the number of keys in a collection
func (c *Collection) Length() int {
	return len(c.KeyMap)
}

// ImportCSV takes a reader and iterates over the rows and imports them as
// a JSON records into dataset.
func (c *Collection) ImportCSV(buf io.Reader, skipHeaderRow bool, idCol int, useUUID bool, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		jsonFName  string
		err        error
	)
	r := csv.NewReader(buf)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	lineNo := 0
	if skipHeaderRow == true {
		lineNo++
		fieldNames, err = r.Read()
		if err != nil {
			return lineNo, fmt.Errorf("Can't read csv table at %d, %s", lineNo, err)
		}
	}
	for {
		lineNo++
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lineNo, fmt.Errorf("Can't read csv table at %d, %s", lineNo, err)
		}
		var fieldName string
		record := map[string]interface{}{}
		if idCol < 0 {
			jsonFName = fmt.Sprintf("%d", lineNo)
		}
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					jsonFName = val
				}
			} else {
				fieldName = fmt.Sprintf("col_%d", i+1)
			}
			//Note: We need to convert the value
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				record[fieldName] = i
			} else if f, err := strconv.ParseFloat(val, 64); err == nil {
				record[fieldName] = f
			} else {
				record[fieldName] = val
			}
		}
		err = c.Create(jsonFName, record)
		if err != nil {
			return lineNo, fmt.Errorf("Can't write %+v to %s, %s", record, jsonFName, err)
		}
		if verboseLog == true && (lineNo%1000) == 0 {
			log.Printf("%d rows processed", lineNo)
		}
	}
	return lineNo, nil
}

// ImportTable takes a [][]string and iterates over the rows and imports them as
// a JSON records into dataset.
func (c *Collection) ImportTable(table [][]string, skipHeaderRow bool, idCol int, useUUID, overwrite, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		key        string
		err        error
	)
	if len(table) == 0 {
		return 0, fmt.Errorf("No data in table")
	}
	lineNo := 0
	// i.e. use the header row for field names
	if skipHeaderRow == true {
		for i, field := range table[lineNo] {
			if strings.TrimSpace(field) == "" {
				fieldNames = append(fieldNames, fmt.Sprintf("column_%03d", i))
			} else {
				fieldNames = append(fieldNames, strings.TrimSpace(field))
			}
		}
		lineNo++
	}
	rowCount := len(table)
	for {
		if lineNo >= rowCount {
			break
		}
		row := table[lineNo]
		lineNo++

		var fieldName string
		record := map[string]interface{}{}
		if idCol < 0 {
			key = fmt.Sprintf("%d", lineNo)
		}
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					key = val
				}
			} else {
				fieldName = fmt.Sprintf("column_%03d", i+1)
			}
			//Note: We need to convert the value
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				record[fieldName] = i
			} else if f, err := strconv.ParseFloat(val, 64); err == nil {
				record[fieldName] = f
			} else if strings.ToLower(val) == "true" {
				record[fieldName] = true
			} else if strings.ToLower(val) == "false" {
				record[fieldName] = false
			} else {
				record[fieldName] = val
			}
		}
		if overwrite == true && c.HasKey(key) == true {
			err = c.Update(key, record)
			if err != nil {
				return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
			}
		} else {
			err = c.Create(key, record)
			if err != nil {
				return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
			}
		}
		if verboseLog == true && (lineNo%1000) == 0 {
			log.Printf("%d rows processed", lineNo)
		}
	}
	return lineNo, nil
}

func colToString(cell interface{}) string {
	var s string
	switch cell.(type) {
	case string:
		s = fmt.Sprintf("%s", cell)
	case json.Number:
		s = fmt.Sprintf("%s", cell.(json.Number).String())
	default:
		src, err := json.Marshal(cell)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		s = fmt.Sprintf("%s", src)
	}
	return s
}

// ExportCSV takes a reader and iterates over the rows and exports then as a CSV file
func (c *Collection) ExportCSV(fp io.Writer, eout io.Writer, filterExpr string, dotExpr []string, colNames []string, verboseLog bool) (int, error) {
	keys, err := c.KeyFilter(c.Keys(), filterExpr)
	if err != nil {
		return 0, err
	}

	// write out colNames
	w := csv.NewWriter(fp)
	if err := w.Write(colNames); err != nil {
		return 0, err
	}

	var (
		cnt           int
		row           []string
		readErrors    int
		writeErrors   int
		dotpathErrors int
	)
	for i, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			// write row out.
			row = []string{}
			for _, colPath := range dotExpr {
				col, err := dotpath.Eval(colPath, data)
				if err == nil {
					row = append(row, colToString(col))
				} else {
					if verboseLog == true {
						log.Printf("error in dotpath %q for key %q in %s, %s\n", colPath, key, c.Name, err)
					}
					dotpathErrors++
					row = append(row, "")
				}
			}
			if err := w.Write(row); err == nil {
				cnt++
			} else {
				if verboseLog == true {
					log.Printf("error writing row %d from %s key %q, %s\n", i+1, c.Name, key, err)
				}
				writeErrors++
			}
			data = nil
		} else {
			log.Printf("error reading %s %q, %s\n", c.Name, key, err)
			readErrors++
		}
	}
	if readErrors > 0 || writeErrors > 0 || dotpathErrors > 0 && verboseLog == true {
		log.Printf("warning %d read error, %d write errors, %d dotpath errors in CSV export from %s", readErrors, writeErrors, dotpathErrors, c.Name)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return cnt, err
	}
	return cnt, nil
}

// KeyFilter takes a list of keys and  filter expression and returns the list of keys passing
// through the filter or an error
func (c *Collection) KeyFilter(keyList []string, filterExpr string) ([]string, error) {
	// Handle the trivial case of filter == "true"
	if filterExpr == "true" {
		return keyList, nil
	}

	// Some sort of filter is involved
	f, err := tmplfn.ParseFilter(filterExpr)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, key := range keyList {
		key = strings.TrimSpace(key)
		if len(key) > 0 {
			m := map[string]interface{}{}
			if err := c.Read(key, m); err == nil {
				if ok, err := f.Apply(m); err == nil && ok == true {
					keys = append(keys, key)
				}
			}
		}
	}
	return keys, nil
}

// Clone copies the current collection records into a newly initialized collection given a list of keys
// and new collection name. Returns an error value if there is a problem. Clone does NOT copy
// attachments, only the JSON records.
func (c *Collection) Clone(keys []string, cloneName string) error {
	if len(keys) == 0 {
		return fmt.Errorf("Zero keys clone from %s to %s", c.Name, cloneName)
	}
	//NOTE: this should create a collection using the same layout we cloning from
	clone, err := InitCollection(cloneName, c.Layout)
	if err != nil {
		return err
	}
	for _, key := range keys {
		src, err := c.ReadJSON(key)
		if err != nil {
			return err
		}
		err = clone.CreateJSON(key, src)
		if err != nil {
			return err
		}
	}
	return nil
}

// CloneSample takes the current collection, a sample size, a training collection name and a test collection
// name. The training collection will be created and receive a random sample of the records from the current
// collection based on the sample size provided. Sample size must be greater than zero and less than the total
// number of records in the current collection.
//
// If the test collection name is not an empty string it will be created and any records not in the training
// collection will be cloned from the current collection into the test collection.
func (c *Collection) CloneSample(sampleSize int, trainingCollectionName string, testCollectionName string) error {
	if sampleSize < 1 {
		return fmt.Errorf("sample size should be greater than zero")
	}
	keys := c.Keys()
	if sampleSize >= len(keys) {
		return fmt.Errorf("sample size too big, %s has %d keys", c.Name, len(keys))
	}
	if len(keys) == 0 {
		return fmt.Errorf("%s has zero keys", c.Name)
	}
	// Apply Sample Size
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffle.Strings(keys, random)
	trainingKeys := keys[0:sampleSize]
	if err := c.Clone(trainingKeys, trainingCollectionName); err != nil {
		return err
	}
	if len(testCollectionName) > 0 {
		testKeys := keys[sampleSize:]
		if err := c.Clone(testKeys, testCollectionName); err != nil {
			return err
		}
	}
	return nil
}

// IsCollection checks to see if a given path contains a
// collection.json file
func IsCollection(p string) bool {
	store, err := storage.GetStore(p)
	if err != nil {
		return false
	}
	if store.IsFile(path.Join(p, "collection.json")) {
		return true
	}
	return false
}

// CollectionLayout returns the numeric type
// association with the collection (e.g BUCKETS_LAYOUT,
// PAIRTREE_LAYOUT).
func CollectionLayout(p string) int {
	store, err := storage.GetStore(p)
	if err != nil {
		return UNKNOWN_LAYOUT
	}
	if store.IsDir(path.Join(p, "pairtree")) {
		return PAIRTREE_LAYOUT
	}
	if store.IsDir(path.Join(p, "aa")) {
		return BUCKETS_LAYOUT
	}
	if store.IsFile(path.Join(p, "collection.json")) {
		src, err := store.ReadFile(path.Join(p, "collection.json"))
		if err != nil {
			return UNKNOWN_LAYOUT
		}
		c := new(Collection)
		err = json.Unmarshal(src, &c)
		if err != nil {
			return UNKNOWN_LAYOUT
		}
		switch c.Layout {
		case BUCKETS_LAYOUT:
			return BUCKETS_LAYOUT
		case PAIRTREE_LAYOUT:
			return PAIRTREE_LAYOUT
		default:
			if len(c.Buckets) > 0 {
				return BUCKETS_LAYOUT
			}
			return UNKNOWN_LAYOUT
		}
	}
	return UNKNOWN_LAYOUT
}
