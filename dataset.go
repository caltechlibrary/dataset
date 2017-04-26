//
// Package dataset is a go package for managing JSON documents stored on disc
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
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/storage"
)

const (
	// Version of the dataset package
	Version = "v0.0.1-beta11"

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

// GenerateBucketNames provides a list of permutations of requested length to use as bucket names
func GenerateBucketNames(alphabet string, length int) []string {
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
	Version string `json:"verison"`
	// Name of collection
	Name string `json:"name"`
	// Buckets is a list of bucket names used by collection
	Buckets []string `json:"buckets"`
	// KeyMap holds the document name to bucket map for the collection
	KeyMap map[string]string `json:"keymap"`
	// SelectLists holds the names of available select lists
	SelectLists []string `json:"select_lists"`
	// Store holds the storage system information (e.g. local disc, S3)
	// and related methods for interacting with it
	Store *storage.Store `json:"-"`
}

// SelectList is an ordered set of keys
type SelectList struct {
	// FName select list filename
	FName string `json:"name"`
	// Keys is the keys stored from a collection
	Keys []string `json:"keys"`
	// CustomLessFn points at the less than function used in sorting
	CustomLessFn func([]string, int, int) bool
	// Store is a pointer to the storage system available
	Store *storage.Store
}

// Len returns the number of keys in the select list
func (s *SelectList) Len() int {
	return len(s.Keys)
}

// Swap updates the position of two compared keys
func (s *SelectList) Swap(i, j int) {
	s.Keys[i], s.Keys[j] = s.Keys[j], s.Keys[i]
}

// Less compare two elements returning true if first is less than second, false otherwise
func (s *SelectList) Less(i, j int) bool {
	if s.CustomLessFn != nil {
		return s.CustomLessFn(s.Keys, i, j)
	}
	return s.Keys[i] < s.Keys[j]
}

// getStore returns a store object, collectionName from name
func getStore(name string) (*storage.Store, string, error) {
	var (
		collectionName string
		store          *storage.Store
		err            error
	)
	// Pick storage based on name
	if strings.HasPrefix(name, "s3://") == true {
		u, err := url.Parse(name)
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
		return store, collectionName, nil
	}

	// Regular file system storage.
	store, err = storage.Init(storage.FS, map[string]interface{}{})
	if err != nil {
		return nil, "", err
	}
	collectionName = name
	return store, collectionName, nil
}

// Create - create a new collection structure on disc
// name should be filesystem friendly
func Create(name string, bucketNames []string) (*Collection, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing a collection name")
	}
	store, collectionName, err := getStore(name)
	if err != nil {
		return nil, err
	}
	// See if we need an open or continue with create
	if _, err := store.Stat(collectionName); err == nil {
		return Open(name)
	}
	c := new(Collection)
	c.Version = Version
	c.Name = collectionName
	c.Buckets = bucketNames
	c.KeyMap = map[string]string{}
	c.SelectLists = []string{"keys"}
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
	if _, err := os.Stat(c.Name); err != nil {
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
	src, err = json.Marshal(c.Keys())
	if err != nil {
		return fmt.Errorf("Can't save key list, %s", err)
	}
	if err := c.Store.WriteFile(path.Join(c.Name, "keys.json"), src, 0664); err != nil {
		return fmt.Errorf("Can't store key list, %s", err)
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
	c.SelectLists = []string{}
	c.Store = nil
	return nil
}

// CreateAsJSON adds or replaces a JSON doc to a collection, if problem returns an error
// name must be unique (treated like a key in a key/value store)
func (c *Collection) CreateAsJSON(name string, src []byte) error {
	keyName, name := keyAndFName(name)

	if _, keyExists := c.KeyMap[keyName]; keyExists == true {
		return c.UpdateAsJSON(name, src)
	}
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection is not valid, zero buckets")
	}
	bucketName := pickBucket(c.Buckets, len(c.KeyMap))
	p := path.Join(c.Name, bucketName)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("mkdir %s %s", p, err)
	}
	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.KeyMap[keyName] = path.Join(bucketName)
	err = c.Store.WriteFile(path.Join(p, name), src, 0664)
	if err != nil {
		return err
	}
	return c.saveMetadata()
}

// Create a JSON doc from an interface{} and adds it  to a collection, if problem returns an error
// name must be unique
func (c *Collection) Create(name string, data interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%s, %s", name, err)
	}
	return c.CreateAsJSON(name, src)
}

func keyAndFName(name string) (string, string) {
	var keyName string
	if strings.HasSuffix(name, ".json") == true {
		keyName = strings.TrimSuffix(name, ".json")
		return keyName, name
	}
	return name, name + ".json"
}

// ReadAsJSON finds a the record in the collection and returns the JSON source
func (c *Collection) ReadAsJSON(name string) ([]byte, error) {
	var keyName string

	keyName, name = keyAndFName(name)
	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return nil, fmt.Errorf("%q does not exist", name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	p := path.Join(c.Name, bucketName)
	src, err := c.Store.ReadFile(path.Join(p, name))
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Read finds the record in a collection, updates the data interface provide and if problem returns an error
// name must exist or an error is returned
func (c *Collection) Read(name string, data interface{}) error {
	src, err := c.ReadAsJSON(name)
	if err != nil {
		return err
	}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal() failed for %s, %s", name, err)
	}
	return nil
}

// UpdateAsJSON takes a JSON doc and writes it to a collection (note: Record must exist or returns an error)
func (c *Collection) UpdateAsJSON(name string, src []byte) error {
	var keyName string

	keyName, name = keyAndFName(name)

	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist", name)
	}
	p := path.Join(c.Name, bucketName)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s", p, err)
	}
	return c.Store.WriteFile(path.Join(p, name), src, 0664)
}

// Update JSON doc in a collection from the provided data interface (note: JSON doc must exist or returns an error )
func (c *Collection) Update(name string, data interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s, %s", name, err)
	}
	return c.UpdateAsJSON(name, src)
}

// Delete removes a JSON doc from a collection
func (c *Collection) Delete(name string) error {
	var keyName string

	keyName, name = keyAndFName(name)

	bucketName, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found", name)
	}
	p := path.Join(c.Name, bucketName, name)
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

//
// Note: Select Lists are an array of keys (JSON documents in the collection but not in buckets)
//

func (c *Collection) hasList(name string) bool {
	keyName, _ := keyAndFName(name)
	for _, k := range c.SelectLists {
		if k == keyName {
			return true
		}
	}
	return false
}

func (c *Collection) getList(name string) (*SelectList, error) {
	var (
		data []string
	)

	_, name = keyAndFName(name)

	src, err := c.Store.ReadFile(path.Join(c.Name, name))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return nil, err
	}
	sl := &SelectList{
		FName: path.Join(c.Name, name),
		Keys:  data,
		Store: c.Store,
	}
	return sl, nil
}

// Select returns a select assocaited with a collection, it will be created if neccessary and
// any keys included will be added before returning the updated list
func (c *Collection) Select(params ...string) (*SelectList, error) {
	var (
		name     string
		listName string
		keys     []string
	)

	if len(params) == 0 {
		name = "keys"
	} else {
		name = params[0]
		if len(params) > 1 {
			keys = params[1:]
		}
	}

	listName, name = keyAndFName(name)
	if name == "collection.json" {
		return nil, fmt.Errorf("%s is not a valid select list", listName)
	}

	if name == "keys.json" {
		return c.getList("keys")
	}

	if c.hasList(listName) == true {
		sl, err := c.getList(listName)
		if err != nil {
			return nil, err
		}
		if len(keys) > 0 {
			sl.Keys = append(sl.Keys, keys[:]...)
			err = sl.SaveList()
		}
		return sl, err
	}

	sl := new(SelectList)
	sl.FName = path.Join(c.Name, name)
	sl.Keys = keys[:]
	sl.Store = c.Store
	err := sl.SaveList()
	if err != nil {
		return nil, err
	}

	c.SelectLists = append(c.SelectLists, listName)
	err = c.saveMetadata()
	if err != nil {
		return nil, err
	}
	return sl, nil
}

// Clear removes a select list from disc and the collection
func (c *Collection) Clear(name string) error {
	var (
		listName string
		err      error
	)

	listName, name = keyAndFName(name)

	if name == "collection.json" {
		return fmt.Errorf("%s is not a select list", listName)
	}
	if name == "keys.json" {
		return fmt.Errorf("cannot clear default select list")
	}

	removeItem := func(s []string, r string) ([]string, error) {
		for i, v := range s {
			if v == r {
				return append(s[:i], s[i+1:]...), nil
			}
		}
		return s, fmt.Errorf("%s not found", r)
	}

	c.SelectLists, err = removeItem(c.SelectLists, listName)
	if err != nil {
		return err
	}

	err = c.Store.Remove(path.Join(c.Name, name))
	if err != nil {
		return err
	}
	return c.saveMetadata()
}

// Lists returns a list of available select lists, should always contain the default keys list
func (c *Collection) Lists() []string {
	return c.SelectLists
}

//
// Select list operations and functions (operating on type SelectList)
//

// String returns the Keys portion of the select list as a string
// delimited with new lines
func (s *SelectList) String() string {
	return strings.Join(s.Keys, "\n")
}

// SaveList writes the .Keys to a JSON document named .FName
func (s *SelectList) SaveList() error {
	if len(s.Keys) == 0 {
		return s.Store.WriteFile(s.FName, []byte("[]"), 0664)
	}
	src, err := json.Marshal(s.Keys)
	if err != nil {
		return err
	}
	return s.Store.WriteFile(s.FName, src, 0664)
}

// First select list returns the first item in the list (non-destructively)
func (s SelectList) First() string {
	if len(s.Keys) > 0 {
		return s.Keys[0]
	}
	return ""
}

// Last select list returns the list item from the list (non-destructively)
func (s *SelectList) Last() string {
	l := len(s.Keys)
	if l > 0 {
		return s.Keys[l-1]
	}
	return ""
}

// Rest select list returns all but the first n items of the list (non-destructively)
func (s *SelectList) Rest() []string {
	l := len(s.Keys)
	if l > 1 {
		return s.Keys[1:]
	}
	return []string{}
}

// List returns all the keys in the select list (non-destructively)
func (s *SelectList) List() []string {
	return s.Keys[:]
}

//
// Mutable select list operations include
// Pop, Push, Shift, Unshift, Sort, Reverse
//

// Pop select list removes from the end of an array returning the element removed
func (s *SelectList) Pop() string {
	pos := len(s.Keys) - 1
	if pos < 0 {
		return ""
	}
	r := s.Keys[pos]
	if pos <= 0 {
		s.Keys = []string{}
	} else {
		s.Keys = s.Keys[0:pos]
	}
	return r
}

// Push select list appends an element to the end of an array
func (s *SelectList) Push(val string) {
	s.Keys = append(s.Keys, val)
}

// Shift select list removes from the beginning of and array returning the element removed
func (s *SelectList) Shift() string {
	l := len(s.Keys)
	if l > 0 {
		r := s.Keys[0]
		if l > 1 {
			s.Keys = s.Keys[1:]
		} else {
			s.Keys = []string{}
		}
		return r
	}
	return ""
}

// Unshift select list inserts an element at the start of an array
func (s *SelectList) Unshift(val string) {
	s.Keys = append([]string{val}, s.Keys[:]...)
}

// Sort sorts the keys in in ascending order alphabetically
func (s *SelectList) Sort(direction int) {
	if direction == DESC {
		sort.Sort(sort.Reverse(s))
		return
	}
	sort.Sort(s)
}

// Reverse flips the order of a select list
func (s *SelectList) Reverse() {
	last := len(s.Keys) - 1
	n := []string{}
	for i := last; i >= 0; i-- {
		n = append(n, s.Keys[i])
	}
	s.Keys = n
}

// Reset a select list to an empty state (file still exists on disc)
func (s *SelectList) Reset() {
	s.Keys = []string{}
}
