//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2020, Caltech
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
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/tbl"
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/namaste"
	"github.com/caltechlibrary/pairtree"
	"github.com/caltechlibrary/shuffle"
	"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"
)

const (
	// Version of the dataset package
	Version = `v0.1.6`

	// License is a formatted from for dataset package based command line tools
	License = `
%s %s

Copyright (c) 2020, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Asc is used to identify ascending sorts
	Asc = iota
	// Desc is used to identify descending sorts
	Desc = iota

	// internal virtualize column name format string
	fmtColumnName = `column_%03d`
)

// Collection is the container holding a pairtree containing JSON docs
type Collection struct {
	// DatasetVersion of the collection
	DatasetVersion string `json:"dataset_version"`

	// Name (filename) of collection
	Name string `json:"name"`

	// workPath holds the path (i.e. non-protocol and hostname, in URI)
	workPath string // `json:"-"`

	// KeyMap holds the document key to path in the collection
	KeyMap map[string]string `json:"keymap"`

	// Store holds the storage system information (i.e. local disc)
	// and related methods for interacting with it
	Store *storage.Store `json:"-"`

	// FrameMap is a list of frame names and with rel path to the frame defined in the collection
	FrameMap map[string]string `json:"frames"`

	//
	// Metadata for collection.
	//

	// Created is the date/time the init command was run in
	// RFC1123 format.
	Created string `json:"created,omitempty"`

	// Version of collection being stored in semvar notation
	Version string `json:"version,omitempty"`

	// Contact info
	Contact string `json:"contact,omitempty"`

	// CodeMeta is a relative path or URL to a Code Meta
	// JSON document for the collection.  Often it'll be
	// in the collection's root and have the value "codemeta.json"
	// but also may be stored someplace else. It should be
	// an empty string if the codemeta.json file has not been
	// created.
	CodeMeta string `json:"codemeta,omitempty"`

	//
	// The following are the Namaste fields
	//

	// Who is the person(s)/organization(s) that created the collection
	Who []string `json:"who,omitempty"`
	// What - description of collection
	What string `json:"what,omitempty"`
	// When - date associated with collection (e.g. 2020,
	// 2020-10, 2020-10-02), should map to an approx date like in
	// archival work.
	When string `json:"when,omitempty"`
	// Where - location (e.g. URL, address) of collection
	Where string `json:"where,omitempty"`

	//
	// private attributes, experiments in performance tuning
	//

	// unsafeSaveMetadata is set to true when doing batch object
	// operations so we can save writing collections.json on each
	// create or delete.
	unsafeSaveMetadata bool

	// collectionMutex is used to prevent write collisions when writing
	// collection.json
	collectionMutex *sync.Mutex

	// objectMutex is used to sync on object writing (e.g. writes involving pairtree path)
	objectMutex *sync.Mutex

	// frameMutex is used to sync on frame writing (e.g. writes involving _frame path)
	frameMutex *sync.Mutex
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

// saveMetadata writes the collection's metadata to  c.Store and c.workPath
func (c *Collection) saveMetadata() error {
	if c.unsafeSaveMetadata == true {
		// NOTE: We're playing fast and loose with the collection metadata, skip saveMetadata().
		return nil
	}
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
	c.collectionMutex.Lock()
	defer c.collectionMutex.Unlock()
	if err := c.Store.WriteFile(path.Join(c.workPath, "collection.json"), src, 0664); err != nil {
		return fmt.Errorf("Can't store collection metadata, %s", err)
	}
	return nil
}

// addNamaste writes name as text files for collection in the collection's directory.
func (c *Collection) addNamaste() error {
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
func InitCollection(name string) (*Collection, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing a collection name")
	}
	store, err := storage.GetStore(name)
	if err != nil {
		return nil, err
	}
	collectionName := collectionNameAsPath(name)
	// See if we need an open or continue with create
	_, err = store.Stat(collectionName + "/collection.json")
	if err == nil {
		return openCollection(name)
	}

	if store.Type == storage.FS {
		err = os.MkdirAll(collectionName, 0775)
		if err != nil {
			return nil, err
		}
	}

	c := new(Collection)
	// Save the date and time
	dt := time.Now()
	// date and time is in RFC3339 format
	c.Created = dt.Format(time.RFC3339)
	// When is a date in YYYY-MM-DD format (can be approximate)
	// e.g. 2020, 2020-01, 2020-01-02
	c.When = dt.Format("2006-01-02")
	c.DatasetVersion = Version
	c.Name = path.Base(collectionName)
	c.Version = "v0.0.0"
	userinfo, err := user.Current()
	if err == nil {
		if userinfo.Name != "" {
			c.Who = []string{userinfo.Name}
		} else {
			c.Who = []string{userinfo.Username}
		}
	}
	if len(c.Who) > 0 {
		c.What = fmt.Sprintf("A dataset (%s) collection initilized on %s by %s.", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST."), strings.Join(c.Who, ", "))
	} else {
		c.What = fmt.Sprintf("A dataset %s collection initilized on %s", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST.."))
	}
	c.workPath = collectionName
	c.KeyMap = map[string]string{}
	c.Store = store

	c.collectionMutex = new(sync.Mutex)
	c.objectMutex = new(sync.Mutex)
	c.frameMutex = new(sync.Mutex)
	err = c.saveMetadata()
	if err != nil {
		return nil, err
	}
	err = c.addNamaste()
	return c, err
}

// openCollection reads in a collection's metadata and returns and new collection structure and err
func openCollection(name string) (*Collection, error) {
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
	if c.KeyMap == nil {
		c.KeyMap = make(map[string]string)
	}

	c.collectionMutex = new(sync.Mutex)
	c.objectMutex = new(sync.Mutex)
	c.frameMutex = new(sync.Mutex)
	return c, nil
}

// deleteCollection an entire collection
func deleteCollection(name string) error {
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
	if c != nil {
		c.Name = ""
		c.workPath = ""
		c.KeyMap = map[string]string{}
		c.Store = nil

		c.collectionMutex = nil
		c.objectMutex = nil
		c.frameMutex = nil
	}
	return nil
}

// CreateJSON adds a JSON doc to a collection, if a problem occurs it returns an error
func (c *Collection) CreateJSON(key string, src []byte) error {
	key = strings.TrimSpace(key)
	if key == "" || key == ".json" {
		return fmt.Errorf("must not be empty")
	}

	// Enforce the _Key attribute is unique and does not exist in collection already
	key = normalizeKeyName(key)
	keyName, FName := keyAndFName(key)
	if _, keyExists := c.KeyMap[keyName]; keyExists == true {
		return fmt.Errorf("%s already exists in collection %s", key, c.Name)
	}

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	var err error
	pair := pairtree.Encode(key)
	pairPath := path.Join("pairtree", pair)
	if c.Store.Type == storage.FS {
		err = c.Store.MkdirAll(path.Join(c.workPath, pairPath), 0770)
		if err != nil {
			return fmt.Errorf("mkdir %s %s", path.Join(c.workPath, pairPath), err)
		}
	}

	// We've almost made it, save the key's name and write the blob to pairtree
	err = c.Store.WriteFile(path.Join(c.workPath, pairPath, FName), src, 0664)
	if err != nil {
		return err
	}
	// We now need to update KeyMap
	c.KeyMap[key] = pairPath
	return c.saveMetadata()
}

// CreateObjectsJSON takes a list of keys and creates a default object for each key
// as quickly as possible. NOTE: if object already exist creation is skipped without
// reporting an error.
func (c *Collection) CreateObjectsJSON(keyList []string, src []byte) error {
	c.unsafeSaveMetadata = true
	defer func() {
		c.unsafeSaveMetadata = false
		c.saveMetadata()
	}()
	for _, key := range keyList {
		if c.KeyExists(key) == false {
			if err := c.CreateJSON(key, src); err != nil {
				return err
			}
		}
	}
	return nil
}

// IsKeyNotFound checks an error message and returns true if
// it is a key not found error.
func (c *Collection) IsKeyNotFound(e error) bool {
	if strings.Compare(e.Error(), "key not found") == 0 {
		return true
	}
	return false
}

// ReadJSON finds a the record in the collection and returns the JSON source
func (c *Collection) ReadJSON(name string) ([]byte, error) {
	if c.KeyExists(name) == false {
		return nil, fmt.Errorf("key not found")
	}
	name = normalizeKeyName(name)
	// Handle potentially URL encoded names
	keyName, FName := keyAndFName(name)
	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return nil, fmt.Errorf("%q does not exist in %s", keyName, c.Name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	// we need to join c.Name + bucketName + name to get path do JSON document
	src, err := c.Store.ReadFile(path.Join(c.workPath, pairPath, FName))
	if err != nil {
		return nil, err
	}
	return src, nil
}

// UpdateJSON a JSON doc in a collection, returns an error if there is a problem
func (c *Collection) UpdateJSON(name string, src []byte) error {
	var ()
	// Normalize key and filenames
	name = normalizeKeyName(name)
	keyName, fName := keyAndFName(name)
	// Make sure Key exists before proceeding with update
	if c.KeyExists(name) == false {
		return fmt.Errorf("key not found")
	}

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}

	// Make sure we preserve attachment metadata
	if bytes.Contains(src, []byte(`"_Attachments"`)) == false {
		obj := map[string]interface{}{}
		if err := c.Read(name, obj, false); err == nil {
			if val, ok := obj["_Attachments"]; ok == true {
				vArray := val.([]interface{})
				if vSrc, err := json.Marshal(vArray); err == nil {
					vSrc = append(append([]byte(`{"_Attachments":`), vSrc...), []byte(",")...)
					src = bytes.Replace(src, []byte(`{`), vSrc, 1)
				}
			}
		}
	}

	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	//NOTE: KeyMap should include pairtree path (e.g. pairtree/AA/BB/CC...)
	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist in %q", keyName, c.Name)
	}
	if c.Store.Type == storage.FS {
		err := c.Store.MkdirAll(path.Join(c.workPath, pairPath), 0770)
		if err != nil {
			return fmt.Errorf("Update (mkdir) %q, %s", path.Join(c.workPath, pairPath), err)
		}
	}
	return c.Store.WriteFile(path.Join(c.workPath, pairPath, fName), src, 0664)
}

// Create a JSON doc from an map[string]interface{} and adds it  to a collection, if problem returns an error
// name must be unique. Document must be an JSON object (not an array).
func (c *Collection) Create(name string, data map[string]interface{}) error {
	src, err := EncodeJSON(data)
	if err != nil {
		return fmt.Errorf("%s, %s", name, err)
	}
	return c.CreateJSON(name, src)
}

// Read finds the record in a collection, updates the data
// interface provide and if problem returns an error
// name must exist or an error is returned
func (c *Collection) Read(name string, data map[string]interface{}, cleanObject bool) error {
	src, err := c.ReadJSON(name)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if cleanObject == true {
		delete(data, "_Key")
		delete(data, "_Attachments")
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
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)

	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found in %q", keyName, c.Name)
	}

	//NOTE: Need to remove any stale tarball before removing our record!
	tarball := strings.TrimSuffix(FName, ".json") + ".tar"
	p := path.Join(c.workPath, pairPath, tarball)
	if err := c.Store.RemoveAll(p); err != nil {
		return fmt.Errorf("Can't remove attachment for %q, %s", keyName, err)
	}
	p = path.Join(c.workPath, pairPath, FName)
	if err := c.Store.Remove(p); err != nil {
		return fmt.Errorf("Error removing %q, %s", p, err)
	}

	delete(c.KeyMap, keyName)
	return c.saveMetadata()
}

// Keys returns a list of keys in a collection
func (c *Collection) Keys() []string {
	keys := []string{}
	for k := range c.KeyMap {
		keys = append(keys, k)
	}
	return keys
}

// KeyExists returns true if key is in collection's KeyMap, false otherwise
func (c *Collection) KeyExists(key string) bool {
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
			if c.KeyExists(key) {
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
			if c.KeyExists(key) == true {
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
	keys := f.Keys[:]
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
		if err := c.Read(key, data, false); err == nil {
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
	keys := f.Keys[:]
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
		if err := c.Read(key, data, false); err == nil {
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

// KeyFilter takes a list of keys and  filter expression and returns
// the list of keys passing through the filter or an error
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
			if err := c.Read(key, m, false); err == nil {
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
	clone, err := InitCollection(cloneName)
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

// Join takes a key, a map[string]interface{}{} and overwrite bool
// and merges the map with an existing JSON object in the collection.
// BUG: This is a naive join, it assumes the keys in object are top
// level properties.
func (c *Collection) Join(key string, obj map[string]interface{}, overwrite bool) error {
	if c.KeyExists(key) == false {
		return c.Create(key, obj)
	}
	record := map[string]interface{}{}
	err := c.Read(key, record, false)
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
