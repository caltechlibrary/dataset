//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/caltechlibrary/pairtree"
)

//
// analyzer checks the collection version and either calls
// bucketAnalyzer or pairtreeAnalyzer as appropriate.
//
func analyzer(collectionName string, verbose bool) error {
	var (
		eCnt int
		wCnt int
		kCnt int
		data interface{}
		c    *Collection
		err  error
	)

	collectionPath := collectionName
	files, err := os.ReadDir(collectionPath)
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
		repairLog(verbose, "WARNING: Missing Namaste 0=dataset_%s\n", Version[1:])
		wCnt++
	}

	// NOTE: Check to see if we have a collection.json
	if hasCollectionJSON == false {
		repairLog(verbose, "WARNING: Missing collection.json\n")
		wCnt++
	} else {
		// Make sure we can JSON parse the file
		docPath := path.Join(collectionPath, "collection.json")
		if src, err := os.ReadFile(docPath); err == nil {
			if err := json.Unmarshal(src, &data); err == nil {
				// release the memory
				data = nil
			} else {
				repairLog(verbose, "ERROR: parsing %s, %s", docPath, err)
				eCnt++
			}
		} else {
			repairLog(verbose, "ERROR: opening %s, %s", docPath, err)
			eCnt++
		}
	}

	// Now try to open the collection ...
	c, err = openCollection(collectionName)
	if err != nil {
		return err
	}

	// Set layout to PAIRTREE_LAYOUT
	// Make sure we have all the known pairs in the pairtree
	// Check to see if records can be found in their buckets
	for k, v := range c.KeyMap {
		// NOTE: as of 1.0.1 keys are forced to lower case internally.
		dirPath := path.Join(collectionPath, v)
		_, err := os.Stat(dirPath)
		if err != nil {
			repairLog(verbose, "ERROR: %s is missing (%q)", k, dirPath)
			eCnt++
		}
		if strings.ToLower(k) != k {
			if c.KeyExists(strings.ToLower(k)) {
				repairLog(true, "ERROR: key %q should be lower case, record CANNOT be merged for case sensitive path %q", k, dirPath)
				eCnt++

			} else {
				repairLog(verbose, "WARNING: key %q should be lower case, record can be merged for case sensitive path %q", k, dirPath)
				wCnt++
			}
		}
		// NOTE: k needs to be urlencoded before checking for file
		fname := url.QueryEscape(k) + ".json"
		docPath := path.Join(collectionPath, v, fname)
		_, err = os.Stat(docPath)
		if err != nil {
			repairLog(verbose, "ERROR: %s is missing (%q)", k, docPath)
			eCnt++
		}
		kCnt++
		if (kCnt % 5000) == 0 {
			repairLog(verbose, "%d of %d keys checked", kCnt, len(c.KeyMap))
		}
	}
	if len(c.KeyMap) > 0 {
		repairLog(verbose, "%d of %d keys checked", kCnt, len(c.KeyMap))
	}

	// Check sub-directories in pairtree find but not in KeyMap
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil && len(c.KeyMap) > 0 {
		repairLog(verbose, "ERROR: unable to walk pairtree, %s", err)
		eCnt++
	} else {
		for _, pair := range pairs {
			key := pairtree.Decode(pair)
			if _, exists := c.KeyMap[key]; exists == false {
				if _, exists := c.KeyMap[strings.ToLower(key)]; exists {
					repairLog(verbose, "WARNING: key %q points at case sensitive path %q",
						strings.ToLower(key), path.Join(collectionName, "pairtree", pair, key+".json"))
					wCnt++
				} else {
					repairLog(verbose, "ERROR: %q found at %q not in collection", key, path.Join(collectionName, "pairtree", pair, key+".json"))
					eCnt++
				}
			}
		}
	}
	// FIXME: need to check for attachments and make sure they are recorded OK

	if eCnt > 0 || wCnt > 0 {
		return fmt.Errorf("%d errors, %d warnings detected", eCnt, wCnt)
	}
	return nil
}

//
// repair takes a collection name and calls
// walks the pairtree and repairs collection.json as appropriate.
//
func repair(collectionName string, verbose bool) error {
	var (
		c   *Collection
		err error
	)

	// See if we can open a collection, if not then create an empty struct
	c, err = openCollection(collectionName)
	if err != nil {
		repairLog(verbose, "Open %s error, %s, attempting to re-create collection.json", collectionName, err)
		err = os.WriteFile(path.Join(collectionName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			repairLog(verbose, "Can't re-initilize %s, %s", collectionName, err)
			return err
		}
		repairLog(verbose, "Attempting to re-open %s", collectionName)

		c, err = openCollection(collectionName)
		if err != nil {
			repairLog(verbose, "Failed to re-open %s, %s", collectionName, err)
			return err
		}
	}
	defer c.Close()

	if c.DatasetVersion != Version {
		repairLog(verbose, "Migrating format from %s to %s", c.DatasetVersion, Version)
	}
	c.DatasetVersion = Version
	repairLog(verbose, "Getting a list of pairs")
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil {
		repairLog(verbose, "ERROR: unable to walk pairtree, %s", err)
		return err
	}
	repairLog(verbose, "Adding missing pairs")
	if c.KeyMap == nil {
		c.KeyMap = map[string]string{}
	}
	for _, pair := range pairs {
		key := pairtree.Decode(pair)
		if strings.ToLower(key) != key {
			if c.KeyExists(key) && c.KeyExists(strings.ToLower(key)) == false {
				tKey := strings.ToLower(key)
				tValue, _ := c.KeyMap[key]
				repairLog(true, "WARNING: moving key %q to %q is being saved lowercase for case sensitive path %q",
					key, tKey, tValue)
				delete(c.KeyMap, key)
				c.KeyMap[tKey] = tValue
			} else if c.KeyExists(key) && c.KeyExists(strings.ToLower(key)) {
				pairPath1, _ := c.KeyMap[key]
				pairPath2, _ := c.KeyMap[strings.ToLower(key)]
				if pairPath1 == "" {
					delete(c.KeyMap, key)
					repairLog(verbose, "WARNING: key %q points at %q.", strings.ToLower(key), pairPath2)
				} else if pairPath1 != pairPath2 {
					repairLog(true, "ERROR: key %q cannot merged as %q for case sensitive path %q.",
						key, strings.ToLower(key), path.Join(c.Name, "pairtree", pair))
				} else {
					repairLog(verbose, "WARNING: previously merged key %q for %q.",
						key, path.Join(c.Name, "paritree", pairPath1))
				}
			} else {
				repairLog(true, "WARNING: key %q added for case sensitive path %q.",
					key, path.Join(c.Name, "pairtree", pair))
				c.KeyMap[strings.ToLower(key)] = path.Join("pairtree", pair)
			}
		} else if _, exists := c.KeyMap[key]; exists == false {
			c.KeyMap[key] = path.Join("pairtree", pair)
		}
	}
	repairLog(verbose, "%d keys in pairtree", len(c.KeyMap))
	keyList := c.Keys()
	repairLog(verbose, "checking that each key resolves to a value on disc")
	missingList := []string{}
	for _, key := range keyList {
		p, err := c.DocPath(key)
		if err != nil {
			break
		}
		if _, err := os.Stat(p); os.IsNotExist(err) == true {
			//NOTE: Mac OS X file system sensitivety handling can
			// mess this assumption up, need to see if we can find
			// the keys we remove and reattach walking the file system.
			repairLog(verbose, "Missing %s from %s, %s does not exist", key, collectionName, p)
			// We save the key to re-attach later...
			missingList = append(missingList, key)
			delete(c.KeyMap, key)
		}
	}
	if len(missingList) > 0 {
		sort.Strings(missingList)
		repairLog(verbose, "Trying to locate %d un-associated keys", len(missingList))
		err = filepath.Walk(path.Join(collectionName, "pairtree"), func(fPath string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				if key, err := url.QueryUnescape(strings.TrimSuffix(info.Name(), ".json")); err == nil {
					// Search our list of keys to see if we can fix path issue...
					for i, missingKey := range missingList {
						r := strings.Compare(key, missingKey)
						if r == 0 {
							kPath := strings.TrimPrefix(strings.TrimSuffix(fPath, info.Name()), collectionName)
							// trim leading separator ...
							kPath = kPath[1:]
							repairLog(verbose, "Fixing path for key %q", key)
							c.KeyMap[key] = kPath
							// Now remove key from missingList
							missingList = append(missingList[:i], missingList[i+1:]...)
							continue
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			repairLog(verbose, "Walking file path error, %s", err)
		}
		// NOTE: the pairtree path in collection.json should be
		// using POSIX path separator.
		for key, value := range c.KeyMap {
			// force paths to be POSIX version.
			if strings.Contains(value, "\\") {
				c.KeyMap[key] = strings.ReplaceAll(value, "\\", "/")
			}
		}
		if len(missingList) > 0 {
			repairLog(verbose, "Unable to find the following keys - %s", strings.Join(missingList, ", "))
		}
	}
	repairLog(verbose, "Saving metadata for %s", collectionName)
	if c.When == "" {
		c.When = time.Now().Format("2006-01-02")
	}
	err = c.saveMetadata()
	if err != nil {
		return err
	}
	return nil
}

//
// Helper functions
//

func repairLog(verbose bool, rest ...interface{}) {
	if verbose == true {
		s := rest[0].(string)
		log.Printf(s, rest[1:]...)
	}
}

// walkPairtree takes a store, a start path and returns a list
// of pairs found that also contain a pair's ${ID}.json file
func walkPairtree(startPath string) ([]string, error) {
	var err error
	// pairs holds a list of discovered pairs
	pairs := []string{}
	err = filepath.Walk(startPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			f := path.Base(p)
			e := path.Ext(f)
			if e == ".json" {
				//NOTE: should be URL encoded at this point.
				key := strings.TrimSuffix(f, e)
				pair := pairtree.Encode(key)
				if strings.Contains(p, path.Join("pairtree", pair, f)) {
					pairs = append(pairs, pair)
				}
			}
		}
		return nil
	})
	return pairs, err
}
