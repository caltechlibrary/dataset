//
// bucket.go is part of the dataset pacakge includes the operations needed for processing collections of JSON documents and their attachments
// using the bucket layout.
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
	"fmt"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/storage"
)

const (
	DefaultAlphabet = `abcdefghijklmnopqrstuvwxyz`
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

// bucketInitCollection - creates a new collection with default alphabet and names of length 2.
func bucketInitCollection(name string) (*Collection, error) {
	return bucketCreateCollection(name, DefaultBucketNames)
}

// bucketCreateCollection - create a new collection structure on disc
// name should be filesystem friendly
func bucketCreateCollection(name string, bucketNames []string) (*Collection, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing a collection name")
	}
	collectionName := collectionNameFromPath(name)
	store, err := storage.GetStore(name)
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
	c.Layout = BUCKETS_LAYOUT
	c.Name = collectionName
	c.Buckets = bucketNames
	c.KeyMap = map[string]string{}
	c.Store = store
	// Save the metadata for collection
	err = c.saveMetadata()
	return c, err
}

// bucketCreateJSON adds a JSON doc to a collection, if a problem occurs it returns an error
func (c *Collection) bucketCreateJSON(key string, src []byte) error {
	if c.Layout != BUCKETS_LAYOUT {
		return fmt.Errorf("Collection does not use a buckets layout")
	}
	key = strings.TrimSpace(key)
	if key == "" || key == ".json" {
		return fmt.Errorf("must not be empty")
	}
	// NOTE: Make sure collection exists before doing anything else!!
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection %q is not valid, zero buckets", c.Name)
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

// bucketReadJSON finds a the record in the collection and returns the JSON source
func (c *Collection) bucketReadJSON(name string) ([]byte, error) {
	if c.Layout != BUCKETS_LAYOUT {
		return nil, fmt.Errorf("Collection does not use a buckets layout")
	}
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

// bucketUpdateJSON a JSON doc in a collection, returns an error if there is a problem
func (c *Collection) bucketUpdateJSON(name string, src []byte) error {
	if c.Layout != BUCKETS_LAYOUT {
		return fmt.Errorf("Collection does not use a buckets layout")
	}
	// NOTE: Make sure collection exists before doing anything else!!
	if len(c.Buckets) == 0 {
		return fmt.Errorf("collection %q is not valid, zero buckets", c.Name)
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

	//NOTE: This is where Pairtree code would go instead of bucketName
	p := path.Join(c.Name, bucketName)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("Update (mkdir) %s %s", p, err)
	}
	return c.Store.WriteFile(path.Join(p, FName), src, 0664)
}

// bucketDelete removes a JSON doc from a collection
func (c *Collection) bucketDelete(name string) error {
	if c.Layout != BUCKETS_LAYOUT {
		return fmt.Errorf("Collection does not use a buckets layout")
	}
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
