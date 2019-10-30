//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2019, Caltech
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
	"path"
	"strings"
	"sync"

	"github.com/caltechlibrary/dataset"
)

/**
 * libdataset.go provides the structures needed to embbed dataset
 * as a service via libdataset.so, libdataset.dll and librdataset.dynlib.
 * E.g. Supporting writing an asynchronous web service written in
 * Python via py_dataset needs save writes.
 */

// LibCollection holds a mutex and collection pointer so we can
// provide service like functionality when using libdataset.
type LibCollection struct {
	Mutex      *sync.Mutex
	Collection *Collection
}

// Lib holds a map of collection names to *LibCollection
type Lib struct {
	collections map[string]*LibCollection
}

var (
	lib *Lib
)

// LibIsOpen checks to see if a dataset collection is already opened.
func LibIsOpen(cName string) bool {
	if lib == nil {
		return false
	}
	if _, exists := lib.collections[cName]; exists == true &&
		lib.collections[cName] != nil &&
		lib.collections[cName].Collection != nil &&
		lib.collections[cName].Collection.Name == cName {
		return true
	}
	return false
}

// LibOpen opens a dataset collection for use in a service like
// context. Lib collections remain "open" until explicitly closed
// or closed via LibCloseAll().
// Writes to the collections are run through a mutex
// to prevent collisions. Subsequent LibOpen() will open
// additional collections under the the service.
func LibOpen(cName string) error {
	var (
		err error
	)
	if lib == nil {
		lib = new(Lib)
		lib.collections = make(map[string]*LibCollection)
	}
	if _, exists := lib.collections[cName]; exists == true {
		return fmt.Errorf("%q opened previously", cName)
	}
	lc := new(LibCollection)
	if lc.Collection, err = Open(cName); err != nil {
		return fmt.Errorf("%q failed to open, %s", cName, err)
	}
	lc.Mutex = new(sync.Mutex)
	lib.collections[cName] = lc
	return nil
}

// LibGetCollection takes a collection name, opens it if necessary and returns a handle
// to the LibCollection struct and error value.
func LibGetCollection(cName string) (*LibCollection, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc, nil
	}
	return nil, fmt.Errorf("%s not found", cName)
}

// LibCollections returns a list of collections previously
// opened with LibOpen()
func LibCollections() []string {
	cNames := []string{}
	if lib == nil {
		return cNames
	}
	for cName, c := range lib.collections {
		if c != nil && c.Collection != nil && path.Base(cName) == c.Collection.Name {
			cNames = append(cNames, cName)
		}
	}
	return cNames
}

// LibClose closes a dataset collections previously
// opened by LibOpen().  It will also set the internal
// lib variable to nil if there are no remaining collections.
func LibClose(cName string) error {
	if LibIsOpen(cName) {
		if lc, exists := lib.collections[cName]; exists == true {
			defer lc.Collection.Close()
			lc.Mutex = nil
			return nil
		}
	}
	return fmt.Errorf("%q not found", cName)
}

// LibCloseAll goes through the service collection list
// and closes each one.
func LibCloseAll() error {
	if lib == nil {
		return fmt.Errorf("Nothing to close")
	}
	errors := []string{}
	for cName, lc := range lib.collections {
		if lc.Collection != nil {
			if err := lc.Collection.Close(); err != nil {
				errors = append(errors, fmt.Sprintf("%q %s", cName, err))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

// LibKeys returns a list of keys for a collection opened with
// StartLib.
func LibKeys(cName string) []string {
	if LibIsOpen(cName) {
		if lc, found := lib.collections[cName]; found == true && lc != nil {
			return lc.Collection.Keys()
		}
	}
	return []string{}
}

// LibKeyExists returns true if the key exists in the collection or false otherwise
func LibKeyExists(cName string, key string) bool {
	if LibIsOpen(cName) {
		if lc, found := lib.collections[cName]; found == true && lc != nil {
			return lc.Collection.KeyExists(key)
		}
	}
	return false
}

// LibKeyFilter returns a list of keys given a list of keys and a filter expression.
func LibKeyFilter(cName string, keys []string, fitlerExpr string) ([]string, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.KeyFilter(keys, fitlerExpr)
	}
	return nil, fmt.Errorf("%q not found", cName)
}

// LibKeySortByExpression returns a list of sorted keys given a list of keys and expression
func LibKeySortByExpression(cName string, keys []string, sortExpr string) ([]string, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.KeySortByExpression(keys, sortExpr)
	}
	return nil, fmt.Errorf("%q not found", cName)
}

// LibCreateJSON takes a collection name, key and JSON object
// document and creates a new JSON object in the collection using
// the key.
func LibCreateJSON(cName string, key string, src []byte) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		err := lc.Collection.CreateJSON(key, src)
		lc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// LibReadJSON takes a collection name, key and returns a JSON object
// document.
func LibReadJSON(cName string, key string) ([]byte, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.ReadJSON(key)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// LibUpdateJSON takes a collection name, key and JSON object
// document and updates the collection.
func LibUpdateJSON(cName string, key string, src []byte) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		err := lc.Collection.UpdateJSON(key, src)
		lc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// LibDeleteJSON takes a collection name and key and removes
// and JSON object from the collection.
func LibDeleteJSON(cName string, key string) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		err := lc.Collection.Delete(key)
		lc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// LibFrameExists returns true if frame found in service collection,
// otherwise false
func LibFrameExists(cName string, fName string) bool {
	if LibIsOpen(cName) == true {
		if lc, found := lib.collections[cName]; found {
			return lc.Collection.FrameExists(fName)
		}
	}
	return false
}

// LibFrameCreate creates a frame in a service collection
func LibFrameCreate(cName string, fName string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		f, err := lc.Collection.FrameCreate(fName, keys, dotPaths, labels, verbose)
		lc.Mutex.Unlock()
		return f, err
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// LibFrameObjects returns a JSON document of a copy of the objects in a frame for
// the service collection. It is analogous to a dataset.ReadJSON but for a frame's
// object list
func LibFrameObjects(cName string, fName string) ([]map[string]interface{}, error) {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return nil, err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.FrameObjects(fName)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// LibFrameRefresh updates the frame object list's for the keys provided. Any new keys
//  cause a new object to be appended to the end of the list.
func LibFrameRefresh(cName string, fName string, keys []string, verbose bool) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.FrameRefresh(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// LibFrameReframe updates the frame object list. If a list of keys is provided then
// the object will be replaced with updated objects based on the keys provided.
func LibFrameReframe(cName string, fName string, keys []string, verbose bool) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.FrameReframe(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// LibFrameClear clears the object and key list from a frame
func LibFrameClear(cName string, fName string) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.FrameClear(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// LibFrameDelete deletes a frame from a service collection
func LibFrameDelete(cName string, fName string) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		return lc.Collection.FrameDelete(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// LibFrames returns a list of frame names in a service collection
func LibFrames(cName string) []string {
	if LibIsOpen(cName) == true {
		if lc, found := lib.collections[cName]; found {
			return lc.Collection.Frames()
		}
	}
	return []string{}
}

// LibCheck checks a dataset collection and reports error to console.
// NOTE: Collection is locked during check!
func LibCheck(cName string, verbose bool) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		dataset.Analyzer(cName, verbose)
		lc.Mutex.Unlock()
	}
	return fmt.Errorf("%q not found", cName)
}

// LibRepair repairs a collection
// NOTE: Collection is locked during repair!
func LibRepair(cName string, verbose bool) error {
	if lib == nil || LibIsOpen(cName) == false {
		if err := LibOpen(cName); err != nil {
			return err
		}
	}
	if lc, found := lib.collections[cName]; found {
		lc.Mutex.Lock()
		dataset.Repair(cName, verbose)
		lc.Mutex.Unlock()
	}
	return fmt.Errorf("%q not found", cName)
}
