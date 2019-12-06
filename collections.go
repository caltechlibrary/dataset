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
)

/**
 * collections.go provides the structures needed to embbed dataset
 * as a service via libdataset.so, libdataset.dll and libdataset.dynlib
 * as well as write cli like the `dataset` command.
 * E.g. Supporting writing an asynchronous web service written in
 * Python via py_dataset needs save writes.
 */

// CMap holds a map of collection names to *Collection
type CMap struct {
	collections map[string]*Collection
}

var (
	cMap *CMap
)

// IsOpen checks to see if a dataset collection is already opened.
func IsOpen(cName string) bool {
	if cMap == nil {
		return false
	}
	if _, exists := cMap.collections[cName]; exists == true &&
		cMap.collections[cName] != nil &&
		cMap.collections[cName].Name == path.Base(cName) {
		return true
	}
	return false
}

// Open opens a dataset collection for use in a service like
// context. CMap collections remain "open" until explicitly closed
// or closed via CloseAll().
// Writes to the collections are run through a mutex
// to prevent collisions. Subsequent CMapOpen() will open
// additional collections under the the service.
func Open(cName string) error {
	var (
		err error
	)
	if cMap == nil {
		cMap = new(CMap)
		cMap.collections = make(map[string]*Collection)
	}
	if _, exists := cMap.collections[cName]; exists == true {
		return fmt.Errorf("%q opened previously", cName)
	}
	c, err := openCollection(cName)
	if err != nil {
		return fmt.Errorf("%q failed to open, %s", cName, err)
	}
	cMap.collections[cName] = c
	return nil
}

// GetCollection takes a collection name, opens it if necessary and returns a handle
// to the CMapCollection struct and error value.
func GetCollection(cName string) (*Collection, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c, nil
	}
	return nil, fmt.Errorf("%s not found", cName)
}

// Collections returns a list of collections previously
// opened with CMapOpen()
func Collections() []string {
	cNames := []string{}
	if cMap == nil {
		return cNames
	}
	for cName, c := range cMap.collections {
		if c != nil && path.Base(cName) == c.Name {
			cNames = append(cNames, cName)
		}
	}
	return cNames
}

// Close closes a dataset collections previously
// opened by CMapOpen().  It will also set the internal
// cMap variable to nil if there are no remaining collections.
func Close(cName string) error {
	if IsOpen(cName) {
		if c, exists := cMap.collections[cName]; exists == true {
			return c.Close()
		}
	}
	return fmt.Errorf("%q not found", cName)
}

// CloseAll goes through the service collection list
// and closes each one.
func CloseAll() error {
	if cMap == nil {
		return fmt.Errorf("Nothing to close")
	}
	errors := []string{}
	for cName, c := range cMap.collections {
		if c != nil {
			if err := c.Close(); err != nil {
				errors = append(errors, fmt.Sprintf("%q %s", cName, err))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

// Keys returns a list of keys for a collection opened with
// StartCMap.
func Keys(cName string) []string {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return []string{}
		}
	}
	if c, found := cMap.collections[cName]; found == true && c != nil {
		return c.Keys()
	}
	return []string{}
}

// KeyExists returns true if the key exists in the collection or false otherwise
func KeyExists(cName string, key string) bool {
	/*
		if cMap == nil || IsOpen(cName) == false {
			if err := Open(cName); err != nil {
				return false
			}
		}
		if c, found := cMap.collections[cName]; found == true && c != nil {
			return c.KeyExists(key)
		}
		return false
	*/
	c, err := GetCollection(cName)
	if err != nil {
		return false
	}
	return c.KeyExists(key)
}

// KeyFilter returns a list of keys given a list of keys and a filter expression.
func KeyFilter(cName string, keys []string, fitlerExpr string) ([]string, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.KeyFilter(keys, fitlerExpr)
	}
	return nil, fmt.Errorf("%q not found", cName)
}

// KeySortByExpression returns a list of sorted keys given a list of keys and expression
func KeySortByExpression(cName string, keys []string, sortExpr string) ([]string, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.KeySortByExpression(keys, sortExpr)
	}
	return nil, fmt.Errorf("%q not found", cName)
}

// CreateJSON takes a collection name, key and JSON object
// document and creates a new JSON object in the collection using
// the key.
func CreateJSON(cName string, key string, src []byte) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		err := c.CreateJSON(key, src)
		c.objectMutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// ReadJSON takes a collection name, key and returns a JSON object
// document.
func ReadJSON(cName string, key string) ([]byte, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.ReadJSON(key)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// UpdateJSON takes a collection name, key and JSON object
// document and updates the collection.
func UpdateJSON(cName string, key string, src []byte) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		err := c.UpdateJSON(key, src)
		c.objectMutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// DeleteJSON takes a collection name and key and removes
// and JSON object from the collection.
func DeleteJSON(cName string, key string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		err := c.Delete(key)
		c.objectMutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameExists returns true if frame found in service collection,
// otherwise false
func FrameExists(cName string, fName string) bool {
	if IsOpen(cName) == true {
		if c, found := cMap.collections[cName]; found {
			return c.FrameExists(fName)
		}
	}
	return false
}

// FrameKeys returns the ordered list of keys for the frame.
func FrameKeys(cName string, fName string) []string {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil
		}
	}
	if c, found := cMap.collections[cName]; found {
		f, err := c.FrameRead(fName)
		if err != nil {
			return nil
		}
		return f.Keys
	}
	return nil
}

// FrameCreate creates a frame in a service collection
func FrameCreate(cName string, fName string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		f, err := c.FrameCreate(fName, keys, dotPaths, labels, verbose)
		c.objectMutex.Unlock()
		return f, err
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// FrameObjects returns a JSON document of a copy of the objects in a frame for
// the service collection. It is analogous to a dataset.ReadJSON but for a frame's
// object list
func FrameObjects(cName string, fName string) ([]map[string]interface{}, error) {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return nil, err
		}
	}
	if c, found := cMap.collections[cName]; found {
		return c.FrameObjects(fName)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// FrameRefresh updates the frame object list's for the keys provided. Any new keys
//  cause a new object to be appended to the end of the list.
func FrameRefresh(cName string, fName string, keys []string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.collectionMutex = new(sync.Mutex)
		c.objectMutex = new(sync.Mutex)
		c.frameMutex = new(sync.Mutex)
		return c.FrameRefresh(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameReframe updates the frame object list. If a list of keys is provided then
// the object will be replaced with updated objects based on the keys provided.
func FrameReframe(cName string, fName string, keys []string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.frameMutex.Lock()
		defer c.frameMutex.Unlock()
		return c.FrameReframe(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameClear clears the object and key list from a frame
func FrameClear(cName string, fName string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.frameMutex.Lock()
		defer c.frameMutex.Unlock()
		return c.FrameClear(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// FrameDelete deletes a frame from a service collection
func FrameDelete(cName string, fName string) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.frameMutex.Lock()
		defer c.frameMutex.Unlock()
		return c.FrameDelete(fName)
	}
	return fmt.Errorf("%q not available", cName)
}

// Frames returns a list of frame names in a service collection
func Frames(cName string) []string {
	if IsOpen(cName) == true {
		if c, found := cMap.collections[cName]; found {
			return c.Frames()
		}
	}
	return []string{}
}

// Check checks a dataset collection and reports error to console.
// NOTE: Collection objects are locked during check!
func Check(cName string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		err := analyzer(cName, verbose)
		c.objectMutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not found", cName)
}

// Repair repairs a collection
// NOTE: Collection objects are locked during repair!
func Repair(cName string, verbose bool) error {
	if cMap == nil || IsOpen(cName) == false {
		if err := Open(cName); err != nil {
			return err
		}
	}
	if c, found := cMap.collections[cName]; found {
		c.objectMutex.Lock()
		err := repair(cName, verbose)
		c.objectMutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not found", cName)
}

// SetWho sets the collection's Who metadata value for a collection
func SetWho(cName string, names []string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.Who = names
	if err = c.saveMetadata(); err != nil {
		return err
	}
	return c.addNamaste()
}

// GetWho get the Who metadata value for a collection.
func GetWho(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return strings.Join(c.Who, "\n")
}

// SetWhat sets the What value for a collection
func SetWhat(cName string, what string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.What = what
	if err = c.saveMetadata(); err != nil {
		return err
	}
	return c.addNamaste()
}

// GetWhat get the What metadata value for a collection.
func GetWhat(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return c.What
}

// SetWhen sets the When value of a collection
func SetWhen(cName string, when string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.When = when
	if err = c.saveMetadata(); err != nil {
		return err
	}
	return c.addNamaste()
}

// GetWhen gets the When value for a collection
func GetWhen(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return c.When
}

// SetWhere sets the Where value of a collection
func SetWhere(cName string, where string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.Where = where
	if err = c.saveMetadata(); err != nil {
		return err
	}
	return c.addNamaste()
}

// GetWhere gets the Where value for a collection
func GetWhere(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return c.Where
}

// SetVersion sets the metadata value for the collection's version.
func SetVersion(cName string, version string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.Version = version
	return c.saveMetadata()
}

// GetVersion gets the version info for the collection.
func GetVersion(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return c.Version
}

// SetContact sets the metadata value for the collection's version.
func SetContact(cName string, contact string) error {
	c, err := GetCollection(cName)
	if err != nil {
		return err
	}
	c.Contact = contact
	return c.saveMetadata()
}

// GetContact gets the contact info for the collection.
func GetContact(cName string) string {
	c, err := GetCollection(cName)
	if err != nil {
		return ""
	}
	return c.Contact
}
