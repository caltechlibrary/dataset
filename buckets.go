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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
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

//
// Analyzer and Repair functions for buckets layout
//
func bucketKeyFound(s string, l []string) bool {
	for _, ky := range l {
		if ky == s {
			return true
		}
	}
	return false
}

func findBuckets(p string) ([]string, error) {
	var buckets []string
	store, err := storage.Init(storage.StorageType(p), nil)
	if err != nil {
		return buckets, err
	}
	dirInfo, err := store.ReadDir(p)
	if err != nil {
		return buckets, err
	}
	for _, item := range dirInfo {
		if item.IsDir() == true {
			buckets = append(buckets, item.Name())
		}
	}
	return buckets, nil
}

// bucketAnalyzer checks a collection for problems
//
// + checks if collection.json exists and is valid
// + checks version of collection and version of dataset tool running
// + checks if all collection.buckets exist
// + checks for unaccounted for buckets
// + checks if all keys in collection.keymap exist
// + checks for unaccounted for keys in buckets
// + checks for keys in multiple buckets and reports duplicate record modified times
//
func bucketAnalyzer(collectionName string) error {
	var (
		eCnt    int
		wCnt    int
		kCnt    int
		data    interface{}
		buckets []string
		c       *Collection
		err     error
	)

	store, err := storage.GetStore(collectionName)
	if err != nil {
		return err
	}
	files, err := store.ReadDir(collectionName)
	if err != nil {
		return err
	}
	hasNamaste := false
	hasCollectionJSON := false
	for _, file := range files {
		fname := file.Name()
		switch {
		case strings.HasPrefix(fname, "0=dataset_"):
			hasNamaste = true
		case fname == "collection.json":
			hasCollectionJSON = true
		}
		if hasNamaste && hasCollectionJSON {
			break
		}
	}

	// NOTE: Check for Namaste 0=, warn if missing
	if hasNamaste == false {
		log.Printf("WARNING: Missing Namaste 0=dataset_%s\n", Version[1:])
		wCnt++
	}

	// NOTE: Check to see if we have a collections.json
	if hasCollectionJSON == false {
		log.Printf("WARNING: Missing collection.json\n")
		wCnt++
	} else {
		// Make sure we can JSON parse the file
		docPath := path.Join(collectionName, "collection.json")
		if src, err := store.ReadFile(docPath); err == nil {
			if err := json.Unmarshal(src, &data); err == nil {
				// release the memory
				data = nil
			} else {
				log.Printf("ERROR: parsing %s, %s", docPath, err)
				eCnt++
			}
		} else {
			log.Printf("ERROR: opening %s, %s", docPath, err)
			eCnt++
		}
	}

	// See if we can open a collection, if not then create an empty struct
	c, err = Open(collectionName)
	if err != nil {
		return fmt.Errorf("ERROR: Open %s, %s", collectionName, err)
	}
	defer c.Close()
	if c.Store.Type != storage.FS {
		return fmt.Errorf("Analyzer only works on local file system")
	}
	if c.Version != Version {
		log.Printf("WARNING: Version mismatch collection %s, dataset %s", c.Version, Version)
		wCnt++
	}

	// Find buckets
	buckets, err = findBuckets(collectionName)
	if err != nil {
		log.Printf("No buckets found for %s, %s", collectionName, err)
		wCnt++
	}
	// Check if buckets match
	for i, bck := range buckets {
		if bucketKeyFound(bck, c.Buckets) == false {
			log.Printf("ERROR: %s is missing from collection bucket list", bck)
			eCnt++
		}
		if i > 0 && (i%100) == 0 {
			log.Printf("%d buckets matched", i)
		}
	}
	if len(buckets) > 0 {
		log.Printf("%d buckets matched", len(buckets))
	}

	// Check to see if records can be found in their buckets
	for ky, bucket := range c.KeyMap {
		docPath := path.Join(collectionName, bucket, ky+".json")
		if store.IsFile(docPath) == false {
			log.Printf("ERROR: %s is missing", docPath)
			eCnt++
		}
		kCnt++
		if (kCnt % 5000) == 0 {
			log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))
		}
	}
	if len(c.KeyMap) > 0 {
		log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))
	}

	// Check for duplicate records and orphaned records
	kCnt = 0
	for j, bck := range buckets {
		if jsonDocs, err := store.FindByExt(path.Join(collectionName, bck), ".json"); err == nil {
			for _, jsonDoc := range jsonDocs {
				ky := strings.TrimSuffix(path.Base(jsonDoc), ".json")
				if val, ok := c.KeyMap[ky]; ok == true {
					if val != bck {
						log.Printf("%s is a duplicate", path.Join(collectionName, val, jsonDoc))
						wCnt++
					}
				} else {
					log.Printf("ERROR: %s is an orphaned JSON Doc", path.Join(collectionName, bck, jsonDoc))
					eCnt++
				}
				kCnt++
			}
		} else {
			log.Printf("ERROR: Can't open bucket %s, %s", bck, err)
			eCnt++
		}
		if (kCnt % 5000) == 0 {
			log.Printf("%d json docs in %d buckets processed", kCnt, j)
		}
	}
	if len(buckets) > 0 {
		log.Printf("%d docs in %d buckets processed", kCnt, len(buckets))
	}

	if eCnt > 0 || wCnt > 0 {
		return fmt.Errorf("%d errors, %d warnings detected", eCnt, wCnt)
	}
	return nil
}

func hasBucket(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

// bucketRepair will take a collection name and attempt to recreate
// valid collection.json from content in discovered buckets and attached documents
func bucketRepair(collectionName string) error {
	var (
		c   *Collection
		err error
	)

	store, err := storage.GetStore(collectionName)
	if err != nil {
		return fmt.Errorf("Repair only works supported storage types, %s", err)
	}
	if store.Type != storage.FS {
		return fmt.Errorf("Repair only works on local file system")
	}

	// See if we can open a collection, if not then create an empty struct
	c, err = Open(collectionName)
	if err != nil {
		log.Printf("Open %s error, %s, attempting to re-create collection.json", collectionName, err)
		err = store.WriteFile(path.Join(collectionName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			log.Printf("Can't re-initilize %s, %s", collectionName, err)
			return err
		}
		log.Printf("Attempting to re-open %s", collectionName)
		c, err = Open(collectionName)
		if err != nil {
			log.Printf("Failed to re-open %s, %s", collectionName, err)
			return err
		}
	}
	defer c.Close()

	if c.Version != Version {
		log.Printf("Migrating format from %s to %s", c.Version, Version)
	}
	c.Version = Version
	log.Printf("Getting a list of buckets")
	if buckets, err := findBuckets(collectionName); err == nil {
		c.Buckets = buckets
	} else {
		return err
	}
	log.Printf("Finding JSON docs in buckets")
	for j, bck := range c.Buckets {
		if c.KeyMap == nil {
			c.KeyMap = map[string]string{}
		}
		if jsonDocs, err := store.FindByExt(path.Join(collectionName, bck), ".json"); err == nil {
			for i, jsonDoc := range jsonDocs {
				ky := strings.TrimSuffix(jsonDoc, ".json")
				if strings.TrimSpace(ky) != "" {
					if val, ok := c.KeyMap[ky]; ok == true {
						if stat1, err := os.Stat(path.Join(collectionName, bck, ky+".json")); err == nil {
							if stat2, err := os.Stat(path.Join(collectionName, val, ky+".json")); err == nil {
								m1 := stat1.ModTime()
								m2 := stat2.ModTime()
								if m1.Unix() > m2.Unix() {
									log.Printf("Switching key %s from %s (%s) to  %s (%s)", ky, val, m2, bck, m1)
									c.KeyMap[ky] = bck
								}
							}
						}
					} else {
						c.KeyMap[ky] = bck
					}
				}
				if i > 0 && (i%5000) == 0 {
					log.Printf("Saving %d items in bucket %s", i, bck)
					if err := c.saveMetadata(); err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}
		log.Printf("Saving bucket %s (%d of %d)", bck, j, len(c.Buckets))
		if err := c.saveMetadata(); err != nil {
			return err
		}
	}
	log.Printf("%d keys in %d buckets", len(c.KeyMap), len(c.Buckets))
	keyList := c.Keys()
	log.Printf("checking that each key resolves to a value on disc")
	for _, key := range keyList {
		p, err := c.DocPath(key)
		if err != nil {
			break
		}
		if _, err := os.Stat(p); os.IsNotExist(err) == true {
			log.Printf("Removing %s from %s, %s does not exist", key, collectionName, p)
			delete(c.KeyMap, key)
		}
	}
	log.Printf("Saving metadata for %s", collectionName)
	if len(c.Buckets) < len(DefaultBucketNames) {
		log.Printf("Adding missing buckets")
		for _, bucket := range DefaultBucketNames {
			if hasBucket(c.Buckets, bucket) == false {
				c.Buckets = append(c.Buckets, bucket)
			}
		}
		log.Printf("Re-sorting buckets")
		sort.Strings(c.Buckets)
	}
	return c.saveMetadata()
}

func migrateToBuckets(collectionName string) error {
	// Open existing collection, get objects and attachments
	// and manually place in new layout updating nc.
	c, err := Open(collectionName)
	if err != nil {
		return err
	}
	oldKeyMap := map[string]string{}
	for k, v := range c.KeyMap {
		oldKeyMap[k] = v
	}
	c.Close()

	store, err := storage.GetStore(collectionName)
	if err != nil {
		return err
	}

	// Create a new collection struct, set to Buckets layout
	nc := new(Collection)
	nc.Name = collectionName
	nc.Layout = BUCKETS_LAYOUT
	nc.Buckets = DefaultBucketNames[:]
	nc.KeyMap = map[string]string{}
	nc.Store, _ = storage.GetStore(collectionName)

	for key, oldPath := range oldKeyMap {
		_, FName := keyAndFName(key)
		src, err := store.ReadFile(path.Join(collectionName, oldPath, FName))
		if err != nil {
			return err
		}
		// Write object to the new location
		err = nc.CreateJSON(key, src)
		if err != nil {
			return err
		}

		// Check for and handle any attachments
		tarballFName := strings.TrimSuffix(FName, ".json") + ".tar"
		oldTarballPath := path.Join(collectionName, oldPath, tarballFName)
		if store.IsFile(oldTarballPath) {
			fmt.Printf("Moving tarball %q\n", oldTarballPath)
			// Move the tarball from one layout to the other
			buf, err := store.ReadFile(oldTarballPath)
			if err != nil {
				return err
			}
			// Find the new location
			docPath, err := nc.DocPath(key)
			if err != nil {
				return err
			}
			newTarballPath := path.Join(collectionName, docPath, tarballFName)
			err = store.WriteFile(newTarballPath, buf, 0664)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
