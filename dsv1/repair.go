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
package dsv1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/caltechlibrary/dataset/v2/pairtree"
)

//
// Analyzer checks the collection version and analyzes current
// state of collection reporting on errors.
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
	files, err := os.ReadDir(collectionPath)
	if err != nil {
		return err
	}
	hasCollectionJSON := false
	for _, file := range files {
		fname := file.Name()
		switch {
		case fname == "collection.json":
			hasCollectionJSON = true
		}
		if hasCollectionJSON {
			break
		}
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
	c, err = Open(collectionName)
	if err != nil {
		return err
	}
	defer c.Close()

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
