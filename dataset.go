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
	"github.com/caltechlibrary/dataset/tbl"
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/namaste"
	"github.com/caltechlibrary/shuffle"
	"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"
)

const (
	// Version of the dataset package
	Version = `v0.0.57`

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

	// Supported file layout types
	// Assume an unknown layout is zero, then add consts in order of adoption
	UNKNOWN_LAYOUT = iota

	// Buckets is the first file layout implemented when dataset started
	BUCKETS_LAYOUT

	// Pairtree is the perferred file layout moving forward
	PAIRTREE_LAYOUT

	// internal virtualize column name format string
	fmtColumnName = `column_%03d`
)

// Collection is the container holding buckets which in turn hold JSON docs
type Collection struct {
	// DatasetVersion of the collection
	DatasetVersion string `json:"dataset_version"`

	// Name of collection
	Name string `json:"name"`

	// workPath holds the path (i.e. non-protocol and hostname, in URI)
	workPath string `json:"-"`

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

	//
	// Metadata for collection
	//

	// Who - creator, owner, maintainer name(s)
	Who []string `json:"who,omitempty"`
	// What - description of collection
	What string `json:"what,omitempty"`
	// When - date associated with collection (e.g. 2018, 2018-10, 2018-10-02)
	When string `json:"when,omitempty"`
	// Where - location (e.g. URL, address) of collection
	Where string `json:"where,omitempty"`
	// Version of collection being stored in semvar notation
	Version string `json:"version,omitempty"`
	// Contact info
	Contact string `json:"contact,omitempty"`
}

//
// internal utility functions
//

// normalizeKeyName() trims leading and trailing spaces
func normalizeKeyName(s string) string {
	return strings.TrimSpace(s)
}

// collectionNameAsPath takes a uri and normalizes collection name
// to a path
func collectionNameAsPath(p string) string {
	if strings.Contains(p, "://") {
		u, _ := url.Parse(p)
		return u.Path
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

// SaveMetadata writes the collection's metadata to  c.Store and c.workPath
func (c *Collection) SaveMetadata() error {
	// Check to see if collection exists, if not create it!
	if c.Store.Type == storage.FS {
		if _, err := c.Store.Stat(c.workPath); err != nil {
			if err := c.Store.MkdirAll(c.workPath, 0775); err != nil {
				return err
			}
		}
	}
	src, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Can't marshal metadata, %s", err)
	}
	if err := c.Store.WriteFile(path.Join(c.workPath, "collection.json"), src, 0664); err != nil {
		return fmt.Errorf("Can't store collection metadata, %s", err)
	}
	// Add/Update Namaste
	loc, err := c.Store.Location(c.workPath)
	if err == nil {
		namaste.DirType(loc, fmt.Sprintf("dataset_%s", Version[1:]))
		if len(c.Who) > 0 {
			for _, who := range c.Who {
				namaste.Who(loc, who)
			}
		}
		if c.What != "" {
			if strings.Contains(c.What, "\n") {
				s := strings.Split(c.What, "\n")
				namaste.What(loc, s[0]+"...")
			} else {
				namaste.What(loc, c.What)
			}
		}
		if c.When != "" {
			if strings.Contains(c.When, "\n") {
				s := strings.Split(c.When, "\n")
				namaste.When(loc, s[0]+"...")
			} else {
				namaste.When(loc, c.When)
			}
		}
		if c.Where != "" {
			if strings.Contains(c.Where, "\n") {
				s := strings.Split(c.Where, "\n")
				namaste.Where(loc, s[0]+"...")
			} else {
				namaste.Where(loc, c.Where)
			}
		}
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
	case BUCKETS_LAYOUT:
		c, err = bucketCreateCollection(name, DefaultBucketNames)
	case PAIRTREE_LAYOUT:
		c, err = pairtreeCreateCollection(name)
	default:
		c, err = pairtreeCreateCollection(name)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Open reads in a collection's metadata and returns and new collection structure and err
func Open(name string) (*Collection, error) {
	store, err := storage.GetStore(name)
	if err != nil {
		return nil, err
	}
	collectionName := collectionNameAsPath(name)
	src, err := store.ReadFile(path.Join(collectionName, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err != nil {
		return nil, err
	}
	//NOTE: we need to reset collectionName so we're working with a path useable to get to the JSON documents.
	c.Name = path.Base(collectionName)
	c.workPath = collectionName
	c.Store = store
	return c, nil
}

// Delete an entire collection
func Delete(name string) error {
	store, err := storage.GetStore(name)
	if err != nil {
		return err
	}
	collectionName := collectionNameAsPath(name)
	if err := store.RemoveAll(collectionName); err != nil {
		return err
	}
	return nil
}

// DocPath returns a full path to a key or an error if not found
func (c *Collection) DocPath(name string) (string, error) {
	keyName, name := keyAndFName(name)
	if p, ok := c.KeyMap[keyName]; ok == true {
		return path.Join(c.workPath, p, name), nil
	}
	return "", fmt.Errorf("Can't find %q", name)
}

// Close closes a collection, writing the updated keys to disc
func (c *Collection) Close() error {
	// Cleanup c so it can't accidentally get reused
	c.Buckets = []string{}
	c.Name = ""
	c.workPath = ""
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
	if c.HasKey(name) == false {
		return nil, fmt.Errorf("key not found")
	}
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
	if c.HasKey(name) == false {
		return fmt.Errorf("key not found")
	}
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
	return hasKey
}

// Length returns the number of keys in a collection
func (c *Collection) Length() int {
	return len(c.KeyMap)
}

// ImportCSV takes a reader and iterates over the rows and imports them as
// a JSON records into dataset.
//BUG: returns lines processed should probably return number of rows imported
func (c *Collection) ImportCSV(buf io.Reader, idCol int, skipHeaderRow bool, overwrite bool, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		key        string
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
			key = fmt.Sprintf("%d", lineNo)
		}
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					key = val
				}
			} else {
				fieldName = fmt.Sprintf(fmtColumnName, i+1)
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
				val = strings.TrimSpace(val)
				if len(val) > 0 {
					record[fieldName] = val
				}
			}
		}
		if len(key) > 0 && len(record) > 0 {
			if c.HasKey(key) {
				if overwrite == true {
					err = c.Update(key, record)
					if err != nil {
						return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
					}
				} else if verboseLog {
					log.Printf("Skipping row %d, key %q, already exists", lineNo, key)
				}
			} else {
				err = c.Create(key, record)
				if err != nil {
					return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
				}
			}
		} else if verboseLog {
			log.Printf("Skipping row %d, key value missing", lineNo)
		}
		if verboseLog == true && (lineNo%1000) == 0 {
			log.Printf("%d rows processed", lineNo)
		}
	}
	return lineNo, nil
}

// ImportTable takes a [][]interface{} and iterates over the rows and
// imports them as a JSON records into dataset.
func (c *Collection) ImportTable(table [][]interface{}, idCol int, useHeaderRow bool, overwrite, verboseLog bool) (int, error) {
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
	if useHeaderRow == true {
		for i, val := range table[0] {
			cell, err := tbl.ValueInterfaceToString(val)
			if err == nil && strings.TrimSpace(cell) != "" {
				fieldNames = append(fieldNames, cell)
			} else {
				fieldNames = append(fieldNames, fmt.Sprintf(fmtColumnName, i))
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
		// Find the key and setup record to save
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					key, err = tbl.ValueInterfaceToString(val)
					if err != nil {
						key = ""
					}
				}
			} else {
				fieldName = fmt.Sprintf(fmtColumnName, i+1)
			}
			record[fieldName] = val
		}
		if len(key) > 0 && len(record) > 0 {
			if c.HasKey(key) == true {
				if overwrite == true {
					err = c.Update(key, record)
					if err != nil {
						return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
					}
				} else if verboseLog == true {
					log.Printf("Skipped row %d, key %s exists in %s", lineNo, key, c.Name)
				}
			} else {
				err = c.Create(key, record)
				if err != nil {
					return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
				}
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

// ExportCSV takes a reader and frame and iterates over the objects
// generating rows and exports then as a CSV file
func (c *Collection) ExportCSV(fp io.Writer, eout io.Writer, f *DataFrame, verboseLog bool) (int, error) {
	//, filterExpr string, dotExpr []string, colNames []string, verboseLog bool) (int, error) {
	if f.AllKeys == true {
		f.Keys = c.Keys()
	}
	keys, err := c.KeyFilter(f.Keys, f.FilterExpr)
	if err != nil {
		return 0, err
	}
	dotExpr := f.DotPaths
	colNames := f.Labels

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
						log.Printf("error in dotpath %q for key %q in %s, %s\n", colPath, key, c.workPath, err)
					}
					dotpathErrors++
					row = append(row, "")
				}
			}
			if err := w.Write(row); err == nil {
				cnt++
			} else {
				if verboseLog == true {
					log.Printf("error writing row %d from %s key %q, %s\n", i+1, c.workPath, key, err)
				}
				writeErrors++
			}
			data = nil
		} else {
			log.Printf("error reading %s %q, %s\n", c.workPath, key, err)
			readErrors++
		}
	}
	if readErrors > 0 || writeErrors > 0 || dotpathErrors > 0 && verboseLog == true {
		log.Printf("warning %d read error, %d write errors, %d dotpath errors in CSV export from %s", readErrors, writeErrors, dotpathErrors, c.workPath)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return cnt, err
	}
	return cnt, nil
}

// ExportTable takes a reader and frame and iterates over the objects
// generating rows and exports then as a CSV file
func (c *Collection) ExportTable(eout io.Writer, f *DataFrame, verboseLog bool) (int, [][]interface{}, error) {
	//, filterExpr string, dotExpr []string, colNames []string, verboseLog bool) (int, error) {
	if f.AllKeys == true {
		f.Keys = c.Keys()
	}
	keys, err := c.KeyFilter(f.Keys, f.FilterExpr)
	if err != nil {
		return 0, nil, err
	}
	dotExpr := f.DotPaths
	colNames := f.Labels

	var (
		cnt           int
		row           []interface{}
		readErrors    int
		dotpathErrors int
	)
	table := [][]interface{}{}
	// Copy column names to table
	for _, colName := range colNames {
		row = append(row, colName)
	}
	table = append(table, row)

	for _, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			// write row out.
			row = []interface{}{}
			for _, colPath := range dotExpr {
				col, err := dotpath.Eval(colPath, data)
				if err == nil {
					row = append(row, col)
				} else {
					if verboseLog == true {
						log.Printf("error in dotpath %q for key %q in %s, %s\n", colPath, key, c.workPath, err)
					}
					dotpathErrors++
					row = append(row, nil)
				}
			}
			table = append(table, row)
			cnt++
			data = nil
		} else {
			log.Printf("error reading %s %q, %s\n", c.workPath, key, err)
			readErrors++
		}
	}
	if (readErrors > 0 || dotpathErrors > 0) && verboseLog == true {
		log.Printf("warning %d read error, %d dotpath errors in table export from %s", readErrors, dotpathErrors, c.workPath)
	}
	return cnt, table, nil
}

// KeyFilter takes a list of keys and  filter expression and returns the list of keys passing
// through the filter or an error
func (c *Collection) KeyFilter(keyList []string, filterExpr string) ([]string, error) {
	// Handle the trivial case of filter resolving to true
	// NOTE: empty filter is treated as "true"
	if filterExpr == "true" || filterExpr == "" {
		return keyList, nil
	}

	// Some sort of filter is involved
	filter, err := tmplfn.ParseFilter(filterExpr)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, key := range keyList {
		key = strings.TrimSpace(key)
		if len(key) > 0 {
			m := map[string]interface{}{}
			if err := c.Read(key, m); err == nil {
				if ok, err := filter.Apply(m); err == nil && ok == true {
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
func (c *Collection) Clone(cloneName string, keys []string, verbose bool) error {
	if len(keys) == 0 {
		return fmt.Errorf("Zero keys clone from %s to %s", c.Name, cloneName)
	}
	//NOTE: this should create a collection using the same layout we cloning from
	clone, err := InitCollection(cloneName, c.Layout)
	if err != nil {
		return err
	}
	i := 0
	for _, key := range keys {
		src, err := c.ReadJSON(key)
		if err != nil {
			return err
		}
		err = clone.CreateJSON(key, src)
		if err != nil {
			return err
		}
		i++
		if verbose && (i%100) == 0 {
			log.Printf("%d objects processed\n", i)
		}
	}
	if verbose {
		log.Printf("%d total objects processed\n", i)
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
func (c *Collection) CloneSample(trainingCollectionName string, testCollectionName string, keys []string, sampleSize int, verbose bool) error {
	if sampleSize < 1 {
		return fmt.Errorf("sample size should be greater than zero")
	}
	if len(keys) == 0 {
		keys = c.Keys()
	}
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
	if err := c.Clone(trainingCollectionName, trainingKeys, verbose); err != nil {
		return err
	}
	if len(testCollectionName) > 0 {
		testKeys := keys[sampleSize:]
		if err := c.Clone(testCollectionName, testKeys, verbose); err != nil {
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
	workPath := collectionNameAsPath(p)
	store, err := storage.GetStore(p)
	if err != nil {
		return UNKNOWN_LAYOUT
	}
	src, err := store.ReadFile(path.Join(workPath, "collection.json"))
	if err == nil {
		c := new(Collection)
		err = json.Unmarshal(src, &c)
		if err != nil {
			return UNKNOWN_LAYOUT
		}
		if len(c.Buckets) > 0 {
			return BUCKETS_LAYOUT
		}
		return PAIRTREE_LAYOUT
	}
	l, err := store.FindByExt(path.Join(workPath, "aa"), ".json")
	if err != nil {
		return PAIRTREE_LAYOUT
	}
	if len(l) > 0 {
		return BUCKETS_LAYOUT
	}
	return UNKNOWN_LAYOUT
}

// Join takes a key, a map[string]interface{}{} and overwrite bool
// and merges the map with an existing JSON object in the collection.
// BUG: This is a naive join, it assumes the keys in object are top
// level properties.
func (c *Collection) Join(key string, obj map[string]interface{}, overwrite bool) error {
	if c.HasKey(key) == false {
		return c.Create(key, obj)
	}
	record := map[string]interface{}{}
	err := c.Read(key, record)
	if err != nil {
		return err
	}

	// Merge object
	for k, v := range obj {
		if overwrite == true {
			record[k] = v
		} else if _, hasProperty := record[k]; hasProperty == false {
			record[k] = v
		}
	}

	// Update record and return
	return c.Update(key, record)
}
