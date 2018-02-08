//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Author R. S. Doiel, <rsdoiel@library.caltech.edu>
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
	Version = `v0.0.17-dev`

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

// Create a JSON doc from an map[string]interface{} and adds it  to a collection, if problem returns an error
// name must be unique. Document must be an JSON object (not an array).
func (c *Collection) Create(name string, data map[string]interface{}) error {
	// NOTE: Make sure collection exists before doing anything else!!
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection is not valid, zero buckets")
	}
	// Enforce the _Key attribute
	keyName, FName := keyAndFName(name)
	if _, keyExists := c.KeyMap[keyName]; keyExists == true {
		return fmt.Errorf("%s already exists in collection", name)
	}

	data["_Key"] = keyName
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s, %s", name, err)
	}

	bucketName := pickBucket(c.Buckets, len(c.KeyMap))
	p := path.Join(c.Name, bucketName)
	err = c.Store.MkdirAll(p, 0770)
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

// CreateFrom is a convienence function that takes an interface, converts it to a map[string]interface{} then calls Create.
func (c *Collection) CreateFrom(name string, obj interface{}) error {
	src, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return err
	}
	return c.Create(name, data)
}

// readAsJSON finds a the record in the collection and returns the JSON source, it is used in search.go
func (c *Collection) readAsJSON(name string) ([]byte, error) {
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

// Read finds the record in a collection, updates the data interface provide and if problem returns an error
// name must exist or an error is returned
func (c *Collection) Read(name string, data map[string]interface{}) error {
	src, err := c.readAsJSON(name)
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

// ReadInto is a convienence function where a Go stuct is converted into a map[string]interface{} then
// passed to Read.
func (c *Collection) ReadInto(name string, obj interface{}) error {
	m := map[string]interface{}{}
	err := c.Read(name, m)
	if err != nil {
		return err
	}
	src, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(src, &obj)
	if err != nil {
		return err
	}
	return nil
}

// Update JSON doc in a collection from the provided data interface (note: JSON doc must exist or returns an error )
func (c *Collection) Update(name string, data map[string]interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s, %s", name, err)
	}
	keyName, FName := keyAndFName(name)

	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist", keyName)
	}
	p := path.Join(c.Name, bucketName)
	err = c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s %s", p, err)
	}
	return c.Store.WriteFile(path.Join(p, FName), src, 0664)
}

// UpdateFrom is a convience function that converts an interface{} into a map[string]interface{}
// before calling update.
func (c *Collection) UpdateFrom(name string, obj interface{}) error {
	src, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return nil
	}
	return c.Update(name, data)
}

// Delete removes a JSON doc from a collection
func (c *Collection) Delete(name string) error {
	keyName, FName := keyAndFName(name)

	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found", keyName)
	}

	//NOTE: Need to remove any stale tarball before removing our record!
	tarball := keyName + ".xml"
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
func (c *Collection) ImportTable(table [][]string, skipHeaderRow bool, idCol int, useUUID bool, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		jsonFName  string
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
func (c *Collection) ExportCSV(fp io.Writer, eout io.Writer, filterExpr string, dotPaths []string, colNames []string, verboseLog bool) (int, error) {
	keys := c.Keys()
	f, err := tmplfn.ParseFilter(filterExpr)
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
			if ok, err := f.Apply(data); err == nil && ok == true {
				// write row out.
				row = []string{}
				for _, colPath := range dotPaths {
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
			}
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

// Extract takes a collection, a filter and a dot path and returns a list of unique values
// E.g. in a collection article records extracting orcid ids which are values in a authors field
func (c *Collection) Extract(filterExpr string, dotPath string) ([]string, error) {
	keys := c.Keys()
	f, err := tmplfn.ParseFilter(filterExpr)
	if err != nil {
		return nil, err
	}

	var (
		data map[string]interface{}
		rows []string
	)
	hash := make(map[string]bool)
	for _, key := range keys {
		if err := c.Read(key, data); err == nil {
			if ok, err := f.Apply(data); err == nil && ok == true {
				col, err := dotpath.Eval(dotPath, data)
				if err == nil {
					hash[colToString(col)] = true
				}
				data = nil
			}
		}
	}
	for ky := range hash {
		rows = append(rows, ky)
	}
	return rows, nil
}
