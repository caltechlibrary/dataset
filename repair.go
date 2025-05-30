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
package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dsv1"
	"github.com/caltechlibrary/pairtree"
	"github.com/caltechlibrary/semver"
)

// sniffVersionNumber tries to get the dataset version
// string from collection.json file. Returns a semver
// or nil (on failure)
func sniffVersionNumber(cName string) *semver.Semver {
	collection := path.Join(cName, "collection.json")
	src, err := ioutil.ReadFile(collection)
	if err != nil {
		return nil
	}
	o := map[string]interface{}{}
	err = json.Unmarshal(src, &o)
	if err != nil {
		return nil
	}
	version, ok := o["dataset"]
	if ok {
		s := version.(string)
		sv, err := semver.ParseString(s)
		if err == nil && sv != nil {
			return sv
		}
		return nil
	}
	return nil
}

// Analyzer checks the collection version and analyzes current
// state of collection reporting on errors.
//
// NOTE: the collection MUST BE CLOSED when Analyzer is called otherwise
// the results will not be accurate.
func Analyzer(cName string, verbose bool) error {
	var (
		eCnt int
		wCnt int
		kCnt int
		data interface{}
		c    *Collection
		err  error
	)

	collectionPath := cName
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

	// Sniff the version number of the collection
	v2 := semver.NewSemver(2, 0, 0, "")
	currentSV := sniffVersionNumber(cName)
	if currentSV != nil && semver.Less(currentSV, v2) {
		repairLog(verbose, "WARNING: %q is a version 1 dataset collection", cName)
		return dsv1.Analyzer(cName, verbose)
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
	c, err = Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()

	if c.StoreType == SQLSTORE {
		_, err := c.SQLStore.Keys()
		if err != nil {
			return fmt.Errorf("WARNING: The collection.json's .name and .dsn_uri to not match the database connection and expected table name.")
		}
		return nil
	}

	if c.StoreType == PTSTORE {
		keymap := path.Join(collectionPath, "keymap.json")
		if _, err := os.Stat(keymap); err != nil {
			repairLog(verbose, "WARNING: Missing keymap.json\n")
			wCnt++
		}
	}

	if c.StoreType != PTSTORE && c.StoreType != SQLSTORE {
		return fmt.Errorf("analyzer only supports pairtree and SQL storage")
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
	pairs, err := walkPairtree(path.Join(cName, "pairtree"))
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
					strings.ToLower(key), path.Join(cName, "pairtree", pair, key+".json"))
				wCnt++
			} else {
				repairLog(verbose, "ERROR: %q found at %q not in collection", key, path.Join(cName, "pairtree", pair, key+".json"))
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

// FixMissingCollectionJson will scan the collection directory
// and environment making an educated guess to type of
// collection collection type
func FixMissingCollectionJson(cName string) error {
	collectionJson := path.Join(cName, "collection.json")
	//FIXME: Need to check to see if we should default to the old pairtree or SQLite3 database
	dsnURI := os.Getenv("DATASET_DSN_URI")
	pairPath := path.Join(cName, "pairtree")
	sqlitePath := path.Join(cName, "collection.db")
	keymapPath := path.Join(cName, "keymap.json")
	storeType := ""
	version := ""
	if _, err := os.Stat(sqlitePath); err == nil {
		storeType = SQLSTORE
		dsnURI = "sqlite://collection.db"
	} else if _, err := os.Stat(pairPath); err == nil {
		storeType = PTSTORE
	} else if dsnURI != "" {
		storeType = SQLSTORE
		version = Version
	}
	if storeType == "" {
		return fmt.Errorf("unable to determine storage type for %q", cName)
	}
	if _, err := os.Stat(keymapPath); err == nil {
		version = Version
	}
	c := &Collection{}
	c.Name = path.Base(cName)
	c.DatasetVersion = version
	c.StoreType = storeType
	c.DsnURI = dsnURI
	src, err := JSONMarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to encode %q, %s", collectionJson, err)
	}
	return ioutil.WriteFile(collectionJson, src, 0664)
}

// Repair a SQLite3 base collection.
func repairSqlite3(c *Collection) error {
	// Fixme see if SQLite3 is installed and in the path
	cmdPath, err := exec.LookPath("sqlite3")
	if err != nil {
		return fmt.Errorf("sqlite3 needs to be installed to repair %q, %s", c.Name, err)
	}
	dbName := path.Join(c.workPath, strings.TrimPrefix(c.DsnURI, "sqlite://"))
	cmdDump := exec.Command(cmdPath, dbName, ".dump")
	src, err := cmdDump.Output()
	if err != nil {
		return fmt.Errorf("failed to execute %q, %s", dbName + " .dump", err)
	}
	os.Rename(dbName, dbName + "-broken")
	cmdRestore := exec.Command(cmdPath, dbName)
	buffer := bytes.Buffer{}
	buffer.Write(src)
	cmdRestore.Stdin = &buffer
	cmdRestore.Stdout = os.Stdout
	cmdRestore.Stderr = os.Stderr
	if err := cmdRestore.Run(); err != nil {
		return fmt.Errorf("failed to retore %q, %s", dbName, err)
	}
	return nil
}

// Repair takes a collection name and calls
// walks the pairtree and repairs collection.json as appropriate.
//
// NOTE: the collection MUST BE CLOSED when repair is called otherwise
// the repaired collection may revert.
func Repair(cName string, verbose bool) error {
	var (
		c   *Collection
		err error
	)
	// Sniff the version number of the collection and delegate
	// if needed.
	v2 := semver.NewSemver(2, 0, 0, "")
	currentSV := sniffVersionNumber(cName)
	if currentSV != nil && semver.Less(currentSV, v2) {
		return fmt.Errorf("cannot repair %q dataset collections", currentSV.String())
	}

	collectionJson := path.Join(cName, "collection.json")
	// Check to see if we find a collection.json, if not see if we
	// can make a educated guess
	if _, err := os.Stat(collectionJson); err != nil {
		err := FixMissingCollectionJson(cName)
		if err != nil {
			return err
		}
	}
	// See if we can open a collection, if not then create an empty struct
	c, err = Open(cName)
	if err != nil {
		repairLog(verbose, "Open %s error, %s, attempting to re-create collection.json", cName, err)
		err = os.WriteFile(path.Join(cName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			repairLog(verbose, "Can't re-initilize %s, %s", cName, err)
			return err
		}
		repairLog(verbose, "Attempting to re-open %s", cName)

		c, err = Open(cName)
		if err != nil {
			repairLog(verbose, "Failed to re-open %s, %s", cName, err)
			return err
		}
	}
	defer c.Close()

	if c.StoreType != PTSTORE {
		// NOTE: in dataset 2.1.x the default storage is no SQLite3
		// We should handle the repair.
		if (strings.HasPrefix(c.DsnURI, "sqlite://")) {
			return repairSqlite3(c)
		}
		return fmt.Errorf("repair supports pairtree storage only")
	}

	c.DatasetVersion = Version
	repairLog(verbose, "Getting a list of pairs")
	pairs, err := walkPairtree(path.Join(cName, "pairtree"))
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
			dirPath := path.Join(cName, "pairtree", pair)
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
			repairLog(verbose, "Missing %s from %s, %s does not exist", key, cName, p)
			// We save the key to re-attach later...
			missingList = append(missingList, key)
			delete(keyMap, key)
			updateKeymap = true
		}
	}
	if updateKeymap {
		updateKeymap = false
		if err := c.PTStore.UpdateKeymap(keyMap); err != nil {
			repairLog(verbose, "Unable to update keymap for %q, %s", cName, err)
		}
	}
	if len(missingList) > 0 {
		sort.Strings(missingList)
		repairLog(verbose, "Trying to locate %d un-associated keys", len(missingList))
		err = filepath.Walk(path.Join(cName, "pairtree"), func(fPath string, info os.FileInfo, err error) error {
			if info.IsDir() == false {
				if key, err := url.QueryUnescape(strings.TrimSuffix(info.Name(), ".json")); err == nil {
					// Search our list of keys to see if we can fix path issue...
					for i, missingKey := range missingList {
						r := strings.Compare(key, missingKey)
						if r == 0 {
							kPath := strings.TrimPrefix(strings.TrimSuffix(fPath, info.Name()), cName)
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

	repairLog(verbose, "Saving metadata for %s", cName)
	// Save the collections' operational metadata
	c.Repaired = time.Now().Format("2006-01-02")
	src, err := JSONMarshalIndent(c, "", "    ")
	filename := path.Join(c.workPath, "collection.json")
	err = ioutil.WriteFile(filename, src, 664)
	if err != nil {
		return err
	}
	return nil
}

// Migrate a dataset v1 collection to a v2 collection.
// Both collections need to already exist. Records from v1
// will be read out of v1 and created in v2.
//
// NOTE: Migrate does not current copy attachments.
func Migrate(srcName string, dstName string, verbose bool) error {
	old, err := dsv1.Open(srcName)
	if err != nil {
		return err
	}
	defer old.Close()
	c, err := Open(dstName)
	if err != nil {
		return err
	}
	defer c.Close()

	keys := old.Keys()
	tot := len(keys)
	eCnt := 0
	for i, key := range keys {
		o := map[string]interface{}{}
		// Removed the v1 object cleanly
		if err := old.Read(key, o, true); err != nil {
			eCnt++
			repairLog(verbose, "failed to read %q from %q, %s", key, srcName, err)
		}
		// FIXME: Need to also handle attachments eventually
		if err := c.Create(key, o); err != nil {
			repairLog(verbose, "failed to write %q from %q, %s", key, dstName, err)
			eCnt++
		}
		if (i % 1000) == 0 {
			repairLog(verbose, "%d of %d processed", i, tot)
		}
	}
	if eCnt > 0 {
		return fmt.Errorf("%d error(s) encounterd in migration", eCnt)
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
