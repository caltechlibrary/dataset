//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/pairtree"
	"github.com/caltechlibrary/dataset/semver"
)

//
// Analyzer checks the collection version and analyzes current
// state of collection reporting on errors.
//
// NOTE: the collection MUST BE CLOSED when Analyzer is called otherwise
// the results will not be accurate.
//
func Analyzer(collectionName string, verbose bool) error {
	var (
		eCnt int
		wCnt int
		kCnt int
		data interface{}
		c    *Collection
		err  error
	)

	collectionPath := collectionName
	// Make sure collection exists
	_, err = os.Stat(collectionPath)
	if err != nil {
		return err
	}

	// Check for collections.json file.
	collection := path.Join(collectionPath, "collection.json")
	_, err = os.Stat(collection)
	if err != nil {
		return err
	}

	// Make sure the JSON documents in the collectionPath can be
	// parsed.
	files, err := os.ReadDir(collectionPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		isDir := file.IsDir()
		if !isDir && strings.HasSuffix(filename, ".json") {
			// Make sure we can JSON parse the file
			docPath := path.Join(collectionPath, filename)
			if src, err := os.ReadFile(docPath); err == nil {
				if err := json.Unmarshal(src, &data); err == nil {
					// release the memory
					data = nil
				} else {
					return fmt.Errorf("error parsing %s, %s", docPath, err)
				}
			} else {
				return fmt.Errorf("error opening %s, %s", docPath, err)
			}
		}
	}

	// NOTE: Check to see if we have a codemeta.json file
	codemeta := path.Join(collectionPath, "codemeta.json")
	_, err = os.Stat(codemeta)
	if err != nil {
		repairLog(verbose, "WARNING: Missing codemeta.json\n")
		wCnt++
	}

	// Now try to open the collection ...
	c, err = Open(collectionName)
	if err != nil {
		return err
	}
	defer c.Close()

	if c.StoreType == PTSTORE {
		keymap := path.Join(collectionPath, "keymap.json")
		if _, err := os.Stat(keymap); err != nil {
			repairLog(verbose, "WARNING: Missing keymap.json\n")
			wCnt++
		}
	}

	currentSV, err := semver.Parse([]byte(Version))
	if err != nil {
		return fmt.Errorf("Can't parse dataset version %q, %s", Version, err)
	}
	collectionSV, err := semver.Parse([]byte(c.DatasetVersion))
	if err != nil {
		return fmt.Errorf("Can't parse dataset version %q, %s", c.DatasetVersion, err)
	}
	if semver.Less(collectionSV, currentSV) {
		return fmt.Errorf("Migration required from %s to %s", c.DatasetVersion, Version)
	}

	if c.StoreType != PTSTORE {
		return fmt.Errorf("analyzer only supports pairtree storage")
	}

	// Set layout to PAIRTREE_LAYOUT
	// Make sure we have all the known pairs in the pairtree
	// Check to see if records can be found in their buckets
	keyMap := c.PTStore.Keymap()
	for k, v := range keyMap {
		// NOTE: as of 1.0.1 keys are forced to lower case internally.
		dirPath := path.Join(collectionPath, "pairtree", v)
		_, err := os.Stat(dirPath)
		if err != nil {
			repairLog(verbose, "ERROR: %s is missing (%q)", k, dirPath)
			eCnt++
		}
		// Is the key in mixed case?
		if strings.ToLower(k) != k {
			// Get lower case key
			lKey := strings.ToLower(k)
			if c.HasKey(lKey) {
				repairLog(true, "ERROR: key %q should be lower case, record CANNOT be merged for case sensitive path %q", k, dirPath)
				eCnt++
			} else {
				repairLog(verbose, "WARNING: key %q should have been lower case, record can be merged for path %q", k, dirPath)
				wCnt++
			}
		}
		// NOTE: k needs to be urlencoded before checking for file
		fname := url.QueryEscape(k) + ".json"
		docPath := path.Join(collectionPath, "pairtree", v, fname)
		_, err = os.Stat(docPath)
		if err != nil {
			repairLog(verbose, "ERROR: %s is missing (%q)", k, docPath)
			eCnt++
		}
		kCnt++
		if (kCnt % 5000) == 0 {
			repairLog(verbose, "%d of %d keys checked", kCnt, len(keyMap))
		}
	}
	if len(keyMap) > 0 {
		repairLog(verbose, "%d of %d keys checked", kCnt, len(keyMap))
	}

	// Check sub-directories in pairtree find but not in KeyMap
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil && len(keyMap) > 0 {
		repairLog(verbose, "ERROR: unable to walk pairtree, %s", err)
		eCnt++
	}
	// Pivot to pair path as key and key as value
	pathMap := map[string]string{}
	for _, pair := range pairs {
		key := pairtree.Decode(pair)
		pathMap[pair] = key
		if _, exists := keyMap[key]; exists == false {
			if _, exists := keyMap[strings.ToLower(key)]; exists {
				repairLog(verbose, "WARNING: key %q points at case sensitive path %q",
					strings.ToLower(key), path.Join(collectionName, "pairtree", pair, key+".json"))
				wCnt++
			} else {
				repairLog(verbose, "ERROR: %q found at %q not in collection", key, path.Join(collectionName, "pairtree", pair, key+".json"))
				eCnt++
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
// Repair takes a collection name and calls
// walks the pairtree and repairs collection.json as appropriate.
//
// NOTE: the collection MUST BE CLOSED when repair is called otherwise
// the repaired collection may revert.
//
func Repair(collectionName string, verbose bool) error {
	var (
		c   *Collection
		err error
	)
	// See if we can open a collection, if not then create an empty struct
	c, err = Open(collectionName)
	if err != nil {
		repairLog(verbose, "Open %s error, %s, attempting to re-create collection.json", collectionName, err)
		err = os.WriteFile(path.Join(collectionName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			repairLog(verbose, "Can't re-initilize %s, %s", collectionName, err)
			return err
		}
		repairLog(verbose, "Attempting to re-open %s", collectionName)

		c, err = Open(collectionName)
		if err != nil {
			repairLog(verbose, "Failed to re-open %s, %s", collectionName, err)
			return err
		}
	}
	defer c.Close()

	currentSV, err := semver.Parse([]byte(Version))
	if err != nil {
		return fmt.Errorf("Can't parse dataset version %q, %s", Version, err)
	}
	collectionSV, err := semver.Parse([]byte(c.DatasetVersion))
	if err != nil {
		return fmt.Errorf("Can't parse dataset version %q, %s", c.DatasetVersion, err)
	}
	if semver.Less(collectionSV, currentSV) {
		return fmt.Errorf("Migration required from %s to %s", c.DatasetVersion, Version)
	}

	if c.StoreType != PTSTORE {
		return fmt.Errorf("repair supports pairtree storage only")
	}

	c.DatasetVersion = Version
	repairLog(verbose, "Getting a list of pairs")
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil {
		repairLog(verbose, "ERROR: unable to walk pairtree, %s", err)
		return err
	}
	repairLog(verbose, "Adding missing pairs")
	keyMap := c.PTStore.Keymap()
	if keyMap == nil {
		keyMap = map[string]string{}
	}
	// Find any missing documents
	updateKeymap := false
	for _, pair := range pairs {
		// Make sure we're generating a lower case key
		key := pairtree.Decode(pair)
		if _, ok := keyMap[key]; ok {
			// We OK, just need to make sure document exists
			dirPath := path.Join(collectionName, "pairtree", pair)
			if _, err := os.Stat(dirPath); err == nil {
				fname := url.QueryEscape(key) + ".json"
				docPath := path.Join(dirPath, fname)
				if _, err := os.Lstat(docPath); err == nil {
					// Document path exists add it to the keymap
					keyMap[key] = pair
					updateKeymap = true
				}
			}
		}
	}
	repairLog(verbose, "%d keys in pairtree", len(keyMap))
	keyList, err := c.PTStore.Keys()
	if err != nil {
		repairLog(verbose, "chould not get keys to repair, %s", err)
	}
	repairLog(verbose, "checking that each key resolves to a value on disc")
	missingList := []string{}
	updateKeymap = false
	for _, key := range keyList {
		p, err := c.PTStore.DocPath(key)
		if err != nil {
			repairLog(verbose, "Missing document path %q (%q), %s", key, p, err)
			delete(keyMap, key)
			updateKeymap = true
			continue
		}
		if _, err := os.Stat(p); os.IsNotExist(err) == true {
			//NOTE: Mac OS X file system sensitivety handling can
			// mess this assumption up, need to see if we can find
			// the keys we remove and reattach walking the file system.
			repairLog(verbose, "Missing %s from %s, %s does not exist", key, collectionName, p)
			// We save the key to re-attach later...
			missingList = append(missingList, key)
			delete(keyMap, key)
			updateKeymap = true
		}
	}
	if updateKeymap {
		updateKeymap = false
		if err := c.PTStore.UpdateKeymap(keyMap); err != nil {
			repairLog(verbose, "Unable to update keymap for %q, %s", collectionName, err)
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
							keyMap[key] = kPath
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
		updateKeymap = false
		for key, value := range keyMap {
			// force paths to be POSIX version.
			if strings.Contains(value, "\\") {
				keyMap[key] = strings.ReplaceAll(value, "\\", "/")
				updateKeymap = true
			}
		}
		if updateKeymap {
			updateKeymap = false
			if err := c.PTStore.UpdateKeymap(keyMap); err != nil {
				repairLog(verbose, "failed to update keymap")
				return err
			}
		}
		if len(missingList) > 0 {
			repairLog(verbose, "Unable to find the following keys - %s", strings.Join(missingList, ", "))
		}
	}

	repairLog(verbose, "Saving metadata for %s", collectionName)
	// Save the collections' operational metadata
	c.Repaired = time.Now().Format("2006-01-02")
	src, err := json.MarshalIndent(c, "", "    ")
	filename := path.Join(c.workPath, "collection.json")
	err = ioutil.WriteFile(filename, src, 664)
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
