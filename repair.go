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
	"fmt"
	"log"
	"strings"
	"time"

	// CaltechLibrary packages
	"github.com/caltechlibrary/namaste"
	"github.com/caltechlibrary/storage"
)

//
// Exported functions for dataset cli usage
//

//
// Analyzer checks the collection version and either calls
// bucketAnalyzer or pairtreeAnalyzer as appropriate.
//
func Analyzer(collectionName string) error {
	var err error
	switch CollectionLayout(collectionName) {
	case BUCKETS_LAYOUT:
		err = bucketAnalyzer(collectionName)
	case PAIRTREE_LAYOUT:
		err = pairtreeAnalyzer(collectionName)
	default:
		err = fmt.Errorf("Unknown layout for %s\n", collectionName)
	}
	return err
}

//
// Repair takes a collection name and calls
// wither bucketRepair or pairtreeRepair as appropriate.
//
func Repair(collectionName string) error {
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
	hasPairtree := false
	hasBuckets := false
	for _, file := range files {
		fname := file.Name()
		switch {
		case strings.HasPrefix(fname, "0=dataset"):
			hasNamaste = true
		case fname == "collection.json":
			hasCollectionJSON = true
		case fname == "pairtree" && file.IsDir() == true:
			hasPairtree = true
		case fname == "aa" && file.IsDir() == true:
			hasBuckets = true
		}
	}
	// NOTE: Check for Namaste 0=, warn and create if missing
	if hasNamaste == false {
		// Add Namaste type record
		namaste.DirType(collectionName, fmt.Sprintf("dataset_%s\n", Version[1:]))
		namaste.When(collectionName, time.Now().Format("2006-01-02"))
	}
	// NOTE: Check to see if we have a collections.json, warn and create if missing
	if hasCollectionJSON == false {
		log.Printf("Missing collection.json, will be regenerating it")
	}

	// NOTE: We're working with buckets (e.g. aa, ab, exists)
	if hasBuckets {
		if err := bucketRepair(collectionName); err != nil {
			return err
		}
	}

	// NOTE: if we're this fair we should repair the pairtree
	if hasPairtree {
		if err := pairtreeRepair(collectionName); err != nil {
			return err
		}
	}
	return nil
}
