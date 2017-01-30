//
// Package dataset is a go package for managing JSON documents stored on disc
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
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
// Examples:
//
//   // Create a collection
//   collection, err := dataset.Create("mystuff", "dataset", GenerateBucketNames("abc", 3))
//   if err != nil {
//       log.Fatalf("%s", err)
//	 }
//   defer collection.Close()
//   // Add a record
//   record := map[string]string{"name":"freda","email":"freda@inverness.example.org"}
//   if err := collection.Create("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Read a record
//   if err := collection.Read("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Update a record
//   record["email"] = "freda@zbs.example.org"
//   if err := collection.Update("freda", record); err != nil {
//       log.Fatalf("%s", err)
//   }
//   // Delete a record
//   if err := collection.Delete("freda"); err != nil {
//       log.Fatalf("%s", err)
//   }
//
package dataset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

const (
	// Version of the dataset package
	Version = "v0.0.1-beta3"

	// License for dataset package
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
	// Dataset is a directory name that holds collections
	Dataset string `json:"dataset"`
	// Buckets is a list of bucket names used by collection
	Buckets []string `json:"buckets"`
	// KeyMap holds the document name to bucket map for the collection
	KeyMap map[string]string `json:"keymap"`
	// SelectLists holds the names of available select lists
	SelectLists []string `json:"select_lists"`
}

// SelectList is an ordered set of keys
type SelectList struct {
	FName string   `json:"name"`
	Keys  []string `json:"keys"`
	// The following are where you add custom sort function for complex key select list
	Len  func() int
	Swap func(int, int)
	Less func(int, int) bool
}

// Create - create a new collection structure on disc
// name should be filesystem friendly
func Create(name string, bucketNames []string) (*Collection, error) {
	if _, err := os.Stat(name); err == nil {
		return Open(name)
	}
	c := new(Collection)
	c.Version = Version
	c.Name = path.Base(name)
	c.Dataset = path.Dir(name)
	c.Buckets = bucketNames
	c.KeyMap = map[string]string{}
	c.SelectLists = []string{"keys"}
	// Save the metadata for collection
	err := c.saveMetadata()
	return c, err
}

// Open reads in a collection's metadata and returns and new collection structure and err
func Open(name string) (*Collection, error) {
	dataPath := path.Dir(name)
	fname := path.Base(name)
	src, err := ioutil.ReadFile(path.Join(dataPath, fname, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err == nil {
		return c, err
	}
	return c, nil
}

// Delete an entire collection
func Delete(name string) error {
	if err := os.RemoveAll(name); err != nil {
		return err
	}
	return nil
}

// saveMetadata writes the collection's metadata to COLLECTION_NAME/collection.json
func (c *Collection) saveMetadata() error {
	if err := os.MkdirAll(path.Join(c.Dataset, c.Name), 0775); err != nil {
		return err
	}
	src, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Can't save metadata, %s", err)
	}
	if err := ioutil.WriteFile(path.Join(c.Dataset, c.Name, "collection.json"), src, 0664); err != nil {
		return err
	}
	src, err = json.Marshal(c.Keys())
	if err != nil {
		return fmt.Errorf("Can't save key list, %s", err)
	}
	if err := ioutil.WriteFile(path.Join(c.Dataset, c.Name, "keys.json"), src, 0664); err != nil {
		return err
	}
	return nil
}

// DocPath returns a full path to a key or an error if not found
func (c *Collection) DocPath(name string) (string, error) {
	keyName, name := keyAndFName(name)
	if bucketName, ok := c.KeyMap[keyName]; ok == true {
		return path.Join(c.Dataset, c.Name, bucketName, name), nil
	}
	return "", fmt.Errorf("Can't find %q", name)
}

// Close closes a collection, writing the updated keys to disc
func (c *Collection) Close() error {
	// Cleanup c so it can't accidentally get reused
	c.Dataset = ""
	c.Buckets = []string{}
	c.Name = ""
	c.KeyMap = map[string]string{}
	c.SelectLists = []string{}
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
	p := path.Join(c.Dataset, c.Name, bucketName)
	err := os.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("mkdir %s %s", p, err)
	}
	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.KeyMap[keyName] = path.Join(bucketName)
	err = ioutil.WriteFile(path.Join(p, name), src, 0664)
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
	p := path.Join(c.Dataset, c.Name, bucketName)
	src, err := ioutil.ReadFile(path.Join(p, name))
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
	p := path.Join(c.Dataset, c.Name, bucketName)
	err := os.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s", p, err)
	}
	return ioutil.WriteFile(path.Join(p, name), src, 0664)
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
	p := path.Join(c.Dataset, c.Name, bucketName, name)
	if err := os.Remove(p); err != nil {
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
		if strings.Compare(k, keyName) == 0 {
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

	src, err := ioutil.ReadFile(path.Join(c.Dataset, c.Name, name))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(src, &data)
	if err != nil {
		return nil, err
	}
	sl := &SelectList{
		FName: path.Join(c.Dataset, c.Name, name),
		Keys:  data,
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

	if strings.Compare(name, "keys.json") == 0 {
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
	sl.FName = path.Join(c.Dataset, c.Name, name)
	sl.Keys = keys[:]
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

	if strings.Compare(name, "keys.json") == 0 {
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

	err = os.Remove(path.Join(c.Dataset, c.Name, name))
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
	src, err := json.Marshal(s.Keys)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.FName, src, 0664)
}

// Length returns the number of items in the select list
func (s *SelectList) Length() int {
	return len(s.Keys)
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
	s.SaveList()
	return r
}

// Push select list appends an element to the end of an array
func (s *SelectList) Push(val string) {
	s.Keys = append(s.Keys, val)
	s.SaveList()
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
		s.SaveList()
		return r
	}
	return ""
}

// Unshift select list inserts an element at the start of an array
func (s *SelectList) Unshift(val string) {
	s.Keys = append([]string{val}, s.Keys[:]...)
	s.SaveList()
}

// Sort sorts the keys in in ascending order alphabetically
func (s *SelectList) Sort(direction int) {
	//FIXME: Need to allow for alternative sorts...
	if s.Swap == nil || s.Len == nil || s.Less == nil {
		if direction == DESC {
			sort.Sort(sort.Reverse(sort.StringSlice(s.Keys)))
			s.SaveList()
			return
		}
		sort.Strings(s.Keys)
	}
	s.SaveList()
}

// Reverse flips the order of a select list
func (s *SelectList) Reverse() {
	last := len(s.Keys) - 1
	n := []string{}
	for i := last; i >= 0; i-- {
		n = append(n, s.Keys[i])
	}
	s.Keys = n
	s.SaveList()
}
