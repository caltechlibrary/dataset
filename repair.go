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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	// CaltechLibrary packages
	"github.com/caltechlibrary/storage"
)

//
// Helper functions
//

func keyFound(s string, l []string) bool {
	for _, ky := range l {
		if ky == s {
			return true
		}
	}
	return false
}

func findBuckets(p string) ([]string, error) {
	var buckets []string

	dirInfo, err := ioutil.ReadDir(p)
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

func findJSONDocs(p string) ([]string, error) {
	var jsonDocs []string

	dirInfo, err := ioutil.ReadDir(p)
	if err != nil {
		return jsonDocs, err
	}
	for _, item := range dirInfo {
		if item.IsDir() == false {
			jname := item.Name()
			if ext := path.Ext(jname); ext == ".json" {
				jsonDocs = append(jsonDocs, jname)
			}
		}
	}
	return jsonDocs, nil
}

func checkFileExists(p string) (string, bool) {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return p, false
	}
	return p, true
}

//
// Exported functionds for dataset cli usage
//

// Analyzer checks a collection for problems
//
// + checks if collection.json exists and is valid
// + checks version of collection and version of dataset tool running
// + checks if all collection.buckets exist
// + checks for unaccounted for buckets
// + checks if all keys in collection.keymap exist
// + checks for unaccounted for keys in buckets
// + checks for keys in multiple buckets and reports duplicate record modified times
//
func Analyzer(collectionName string) error {
	var (
		eCnt    int
		wCnt    int
		kCnt    int
		data    interface{}
		buckets []string
		c       *Collection
		err     error
	)

	if strings.HasPrefix(collectionName, "s3://") || strings.HasPrefix(collectionName, "gs://") {
		return fmt.Errorf("Analyzer only works on local disc storage")
	}

	// Check of collections.json
	for _, fname := range []string{"collection.json"} {
		if _, exists := checkFileExists(collectionName); exists == false {
			return fmt.Errorf("%q does not exist", collectionName)
		}
		if docPath, exists := checkFileExists(path.Join(collectionName, fname)); exists == false {
			log.Printf("Missing %s", docPath)
			eCnt++
		} else {
			// Make sure we can JSON parse the file
			if src, err := ioutil.ReadFile(docPath); err == nil {
				if err := json.Unmarshal(src, &data); err == nil {
					// release the memory
					data = nil
				} else {
					log.Printf("Error parsing %s, %s", docPath, err)
					eCnt++
				}
			} else {
				log.Printf("Error opening %s, %s", docPath, err)
				eCnt++
			}
		}
	}

	// See if we can open a collection, if not then create an empty struct
	if c, err = Open(collectionName); err == nil {
		defer c.Close()
	} else {
		log.Printf("ERROR: Open collection error, %s", err)
		c = new(Collection)
		eCnt++
	}
	if c.Store.Type != storage.FS {
		return fmt.Errorf("Analyzer only works on local disc storage")
	}
	if c.Version != Version {
		log.Printf("WARNING: Version mismatch collection %s, dataset %s", c.Version, Version)
		wCnt++
	}

	// Find buckets
	buckets, err = findBuckets(collectionName)
	if err != nil {
		log.Printf("No buckets found for %s, %s", collectionName, err)
		eCnt++
	}
	// Check if buckets match
	log.Printf("Checking buckets")
	for i, bck := range buckets {
		if keyFound(bck, c.Buckets) == false {
			log.Printf("%s is missing from collection bucket list", bck)
			eCnt++
		}
		if i > 0 && (i%100) == 0 {
			log.Printf("%d buckets matched", i)
		}
	}
	log.Printf("%d buckets matched", len(buckets))

	// Check to see if records can be found in their buckets
	log.Printf("Checking for %d keys from keymaps against their buckets", len(c.KeyMap))
	for ky, bucket := range c.KeyMap {
		if docPath, exists := checkFileExists(path.Join(collectionName, bucket, ky+".json")); exists == false {
			log.Printf("%s is missing", docPath)
			eCnt++
		}
		kCnt++
		if (kCnt % 5000) == 0 {
			log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))
		}
	}
	log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))

	// Check for duplicate records and orphaned records
	log.Printf("Scanning buckets for orphaned and duplicate records")
	kCnt = 0
	for j, bck := range buckets {
		if jsonDocs, err := findJSONDocs(path.Join(collectionName, bck)); err == nil {
			for _, jsonDoc := range jsonDocs {
				ky := strings.TrimSuffix(path.Base(jsonDoc), ".json")
				if val, ok := c.KeyMap[ky]; ok == true {
					if val != bck {
						log.Printf("%s is a duplicate", path.Join(collectionName, val, jsonDoc))
						wCnt++
					}
				} else {
					log.Printf("%s is an orphaned JSON Doc", path.Join(collectionName, bck, jsonDoc))
					eCnt++
				}
				kCnt++
			}
		} else {
			log.Printf("Can't open bucket %s, %s", bck, err)
			eCnt++
		}
		if (kCnt % 5000) == 0 {
			log.Printf("%d json docs in %d buckets processed", kCnt, j)
		}
	}
	log.Printf("%d docs in %d buckets processed", kCnt, len(buckets))

	if eCnt > 0 || wCnt > 0 {
		return fmt.Errorf("%d errors, %d warnings detected", eCnt, wCnt)
	}
	return nil
}

// Repair will take a collection name and attempt to recreate
// valid collection.json from content in discovered buckets and attached documents
func Repair(collectionName string) error {
	var (
		c   *Collection
		err error
	)

	// See if we can open a collection, if not then create an empty struct
	if c, err = Open(collectionName); err == nil {
		defer c.Close()
	} else {
		log.Printf("ERROR: Open collection error, %s", err)
		c = new(Collection)
	}
	if c.Store.Type != storage.FS {
		return fmt.Errorf("Repair only works on local disc storage")
	}
	if c.Version != Version {
		log.Printf("Migrating format from %s to %s", c.Version, Version)
	}
	c.Version = Version
	log.Printf("Getting a list of buckets")
	if buckets, err := findBuckets(path.Join(collectionName)); err == nil {
		c.Buckets = buckets
	} else {
		return err
	}
	log.Printf("Finding JSON docs in buckets")
	for j, bck := range c.Buckets {
		if jsonDocs, err := findJSONDocs(path.Join(collectionName, bck)); err == nil {
			for i, jsonDoc := range jsonDocs {
				ky := strings.TrimSuffix(jsonDoc, ".json")
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
	log.Printf("Saving metadata for %s", collectionName)
	return c.saveMetadata()
}
