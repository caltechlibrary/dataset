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
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/dotpath"
	"github.com/caltechlibrary/storage"
	"github.com/caltechlibrary/tmplfn"

	// 3rd Party packages
	"github.com/google/uuid"
)

const (
	// Version of the dataset package
	Version = `v0.0.29-dev`

	// License is a formatted from for dataset package based command line tools
	License = `
%s %s

Copyright (c) 2017, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	DefaultAlphabet = `abcdefghijklmnopqrstuvwxyz`

	ASC  = iota
	DESC = iota
)

// countToBucketID turns a count assigns it to a letter sequence (e.g. 0-999 is aa, 1000 - 1999 is ab, etc)
func countToBucketID(i int, bucketNames []string) string {
	bucketsize := len(bucketNames)
	// Calculate bucket number
	bucketIndex := i % bucketsize
	return bucketNames[bucketIndex]
}

// generateCombinations from an alphabet and length
//
// this function is based on example at https://play.golang.org/p/0bWDCibSUJ
func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c) // Once the iteration function is finished, we close the channel

		addLetter(c, "", alphabet, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// addLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
//
// this function is based on gist at https://play.golang.org/p/0bWDCibSUJ
func addLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

// pickBucket converts takes the number of picks and the
// count of JSON docs and returns a bucket name.
func pickBucket(buckets []string, docNo int) string {
	bucketCount := len(buckets)
	return buckets[(docNo % bucketCount)]
}

// generateBucketNames provides a list of permutations of requested length to use as bucket names
func generateBucketNames(alphabet string, length int) []string {
	l := []string{}
	for combo := range generateCombinations(alphabet, length) {
		if len(combo) == length {
			l = append(l, combo)
		}
	}
	return l
}

// Collection is the container holding buckets which in turn hold JSON docs
type Collection struct {
	// Version of collection being stored
	Version string `json:"version"`
	// Name of collection
	Name string `json:"name"`
	// Buckets is a list of bucket names used by collection
	Buckets []string `json:"buckets"`
	// KeyMap holds the document name to bucket map for the collection
	KeyMap map[string]string `json:"keymap"`
	// Store holds the storage system information (e.g. local disc, S3, GS)
	// and related methods for interacting with it
	Store *storage.Store `json:"-"`
	// FullPath is the fully qualified path on disc or URI to S3 or GS bucket
	FullPath string `json:"-"`
}

// getStore returns a store object, collectionName from name
func getStore(name string) (*storage.Store, string, error) {
	var (
		collectionName string
		store          *storage.Store
		err            error
	)
	// Pick storage based on name
	switch {
	case strings.HasPrefix(name, "s3://") == true:
		u, _ := url.Parse(name)
		opts := storage.EnvToOptions(os.Environ())
		opts["AwsBucket"] = u.Host
		store, err = storage.Init(storage.S3, opts)
		if err != nil {
			return nil, "", err
		}
		p := u.Path
		if strings.HasPrefix(p, "/") {
			p = p[1:]
		}
		collectionName = p
	case strings.HasPrefix(name, "gs://") == true:
		u, _ := url.Parse(name)
		opts := storage.EnvToOptions(os.Environ())
		opts["GoogleBucket"] = u.Host
		store, err = storage.Init(storage.GS, opts)
		if err != nil {
			return nil, "", err
		}
		p := u.Path
		if strings.HasPrefix(p, "/") {
			p = p[1:]
		}
		collectionName = p
	default:
		// Regular file system storage.
		store, err = storage.Init(storage.FS, map[string]interface{}{})
		if err != nil {
			return nil, "", err
		}
		collectionName = name
	}

	return store, collectionName, nil
}

// normalizeKeyName() trims leading and trailing spaces
func normalizeKeyName(s string) string {
	return strings.TrimSpace(s)
}

// InitCollection - creates a new collection with default alphabet and names of length 2.
func InitCollection(name string) (*Collection, error) {
	return create(name, DefaultBucketNames)
}

// create - create a new collection structure on disc
// name should be filesystem friendly
func create(name string, bucketNames []string) (*Collection, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing a collection name")
	}
	store, collectionName, err := getStore(name)
	if err != nil {
		return nil, err
	}
	// See if we need an open or continue with create
	if store.Type == storage.S3 || store.Type == storage.GS {
		if _, err := store.Stat(collectionName + "/collection.json"); err == nil {
			return Open(name)
		}
	} else {
		if _, err := store.Stat(collectionName); err == nil {
			return Open(name)
		}
	}
	c := new(Collection)
	c.Version = Version
	c.Name = collectionName
	c.Buckets = bucketNames
	c.KeyMap = map[string]string{}
	c.Store = store
	// Save the metadata for collection
	err = c.saveMetadata()
	return c, err
}

// Open reads in a collection's metadata and returns and new collection structure and err
func Open(name string) (*Collection, error) {
	store, collectionName, err := getStore(name)
	if err != nil {
		return nil, err
	}
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
	store, collectionName, err := getStore(name)
	if err != nil {
		return err
	}
	if err := store.RemoveAll(collectionName); err != nil {
		return err
	}
	return nil
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

// DocPath returns a full path to a key or an error if not found
func (c *Collection) DocPath(name string) (string, error) {
	keyName, name := keyAndFName(name)
	if bucketName, ok := c.KeyMap[keyName]; ok == true {
		return path.Join(c.Name, bucketName, name), nil
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

// keyAndFName converts a key (which may have things like slashes) into a disc friendly name and key value
func keyAndFName(name string) (string, string) {
	var keyName string
	if strings.HasSuffix(name, ".json") == true {
		return keyName, name
	}
	return name, url.QueryEscape(name) + ".json"
}

// CreateJSON adds a JSON doc to a collection, if a problem occurs it returns an error
func (c *Collection) CreateJSON(name string, src []byte) error {
	// NOTE: Make sure collection exists before doing anything else!!
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection is not valid, zero buckets")
	}

	// Enforce the _Key attribute is unique and does not exist in collection already
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)
	if _, keyExists := c.KeyMap[keyName]; keyExists == true {
		return fmt.Errorf("%s already exists in collection", name)
	}

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	bucketName := pickBucket(c.Buckets, len(c.KeyMap))
	p := path.Join(c.Name, bucketName)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("mkdir %s %s", p, err)
	}

	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.KeyMap[keyName] = path.Join(bucketName)
	err = c.Store.WriteFile(path.Join(p, FName), src, 0664)
	if err != nil {
		return err
	}
	return c.saveMetadata()
}

// ReadJSON finds a the record in the collection and returns the JSON source
func (c *Collection) ReadJSON(name string) ([]byte, error) {
	name = normalizeKeyName(name)
	// Handle potentially URL encoded names
	keyName, FName := keyAndFName(name)
	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return nil, fmt.Errorf("%q does not exist in %s", keyName, c.Name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	// we need to join c.Name + bucketName + name to get path do JSON document
	src, err := c.Store.ReadFile(path.Join(c.Name, bucketName, FName))
	if err != nil {
		return nil, err
	}
	return src, nil
}

// UpdateJSON a JSON doc in a collection, returns an error if there is a problem
func (c *Collection) UpdateJSON(name string, src []byte) error {
	// NOTE: Make sure collection exists before doing anything else!!
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection is not valid, zero buckets")
	}

	// Make sure Key exists before proceeding with update
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)
	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist", keyName)
	}

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	p := path.Join(c.Name, bucketName)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("Update (mkdir) %s %s", p, err)
	}
	return c.Store.WriteFile(path.Join(p, FName), src, 0664)
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
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)

	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found", keyName)
	}

	//NOTE: Need to remove any stale tarball before removing our record!
	tarball := keyName + ".tar"
	p := path.Join(c.Name, bucketName, tarball)
	if err := c.Store.RemoveAll(p); err != nil {
		return fmt.Errorf("Can't remove attachment for %q, %s", keyName, err)
	}
	p = path.Join(c.Name, bucketName, FName)
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
func (c *Collection) ImportCSV(buf io.Reader, skipHeaderRow bool, idCol int, useUUID bool, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		jsonFName  string
		err        error
	)
	r := csv.NewReader(buf)
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
		if idCol < 0 && useUUID == false {
			jsonFName = fmt.Sprintf("%d", lineNo)
		} else if useUUID == true {
			jsonFName = uuid.New().String()
			if _, ok := record["uuid"]; ok == true {
				record["_uuid"] = jsonFName
			} else {
				record["uuid"] = jsonFName
			}
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
		if idCol < 0 && useUUID == false {
			key = fmt.Sprintf("%d", lineNo)
		} else if useUUID == true {
			key = uuid.New().String()
			if _, ok := record["uuid"]; ok == true {
				record["_uuid"] = key
			} else {
				record["uuid"] = key
			}
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
		cnt        int
		row        []string
		readErrors int
	)
	for _, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			// write row out.
			row = []string{}
			for _, colPath := range dotExpr {
				col, err := dotpath.Eval(colPath, data)
				if err == nil {
					row = append(row, colToString(col))
				} else {
					row = append(row, "")
				}
			}
			if err := w.Write(row); err == nil {
				cnt++
			}
			data = nil
		} else {
			fmt.Fprintf(os.Stderr, "error reading %q, %s\n", key, err)
			readErrors += 1
		}
	}
	if readErrors > 0 {
		return cnt, fmt.Errorf("%d read errors encountered", readErrors)
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
		m := map[string]interface{}{}
		if err := c.Read(key, m); err == nil {
			if ok, err := f.Apply(m); err == nil && ok == true {
				keys = append(keys, key)
			}
		}
	}
	return keys, nil
}

// Extract takes a collection, a filter and a dot path and returns a list of unique values
// E.g. in a collection article records extracting orcid ids which are values in a authors field
func (c *Collection) Extract(filterExpr string, dotExpr string) ([]string, error) {

	keys, err := c.KeyFilter(c.Keys(), filterExpr)
	if err != nil {
		return nil, err
	}

	uniqueStrings := map[string]bool{}
	hKey := ""
	for _, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			if cell, err := dotpath.Eval(dotExpr, data); err == nil {
				switch cell.(type) {
				case []interface{}:
					for _, v := range cell.([]interface{}) {
						hKey = colToString(v)
						uniqueStrings[hKey] = true
					}
				case map[string]interface{}:
					for _, v := range cell.([]interface{}) {
						hKey = colToString(v)
						uniqueStrings[hKey] = true
					}
				default:
					hKey = colToString(cell)
					uniqueStrings[hKey] = true
				}
			} else if err != nil {
				return nil, fmt.Errorf("can't parse dotExpr %q", err)
			}
		} else {
			return nil, fmt.Errorf("c.Read() error, %s, %s", key, err)
		}
	}
	rows := []string{}
	for ky, _ := range uniqueStrings {
		rows = append(rows, ky)
	}
	sort.Strings(rows)
	return rows, nil
}
