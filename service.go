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
 * service.go provides the structures needed to embbed dataset
 * as a service via libdataset.go. E.g. to support writing
 * asynchronous web service written in Python via py_dataset.
 */

// ServiceCollection holds a mutex and collection pointer so we can provide service like
// functionality for via libdataset.
type ServiceCollection struct {
	Mutex      *sync.Mutex
	Collection *Collection
}

// Service holds a map of collection names to service collections
type Service struct {
	collections map[string]*ServiceCollection
}

var (
	service *Service
)

// ServiceOpen opens a data collection for use in a service like
// context. Service collections remain "open" until explicitly closed
// or the service is shutdown..
// Writes to the collections are run through a mutex
// to prevent collisions. Subsequent ServiceOpen() will open
// additional collections under the the service.
func ServiceOpen(cName string) error {
	var (
		err error
	)
	if service == nil {
		service = new(Service)
		service.collections = make(map[string]*ServiceCollection)
	}
	if _, exists := service.collections[cName]; exists == true {
		return fmt.Errorf("%q opened previously", cName)
	}
	sc := new(ServiceCollection)
	if sc.Collection, err = Open(cName); err != nil {
		return fmt.Errorf("%q failed to open, %s", cName, err)
	}
	sc.Mutex = new(sync.Mutex)
	service.collections[cName] = sc
	return nil
}

// ServiceCollections returns a list of collections previously
// opened with ServiceOpen()
func ServiceCollections() []string {
	cNames := []string{}
	if service == nil {
		return cNames
	}
	for cName, c := range service.collections {
		if c != nil && c.Collection != nil && path.Base(cName) == c.Collection.Name {
			cNames = append(cNames, cName)
		}
	}
	return cNames
}

// ServiceClose closes a dataset collections previously
// opened by ServiceOpen().  It will also set the internal
// service variable to nil if there are no remaining collections.
func ServiceClose(cName string) error {
	if service == nil {
		return fmt.Errorf("%q not open", cName)
	}
	if sc, exists := service.collections[cName]; exists == true {
		defer sc.Collection.Close()
		sc.Mutex = nil
		return nil
	}
	return fmt.Errorf("%q not found in service", cName)
}

// ServiceCloseAll goes through the service collection list
// and closes each one.
func ServiceCloseAll() error {
	if service == nil {
		return fmt.Errorf("Nothing to close")
	}
	errors := []string{}
	for cName, sc := range service.collections {
		if sc.Collection != nil {
			if err := sc.Collection.Close(); err != nil {
				errors = append(errors, fmt.Sprintf("%q %s", cName, err))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

// ServiceKeys returns a list of keys for a collection opened with
// StartService.
func ServiceKeys(cName string) []string {
	if service == nil {
		return []string{}
	}
	if sc, found := service.collections[cName]; found == true && sc != nil {
		return sc.Collection.Keys()
	}
	return []string{}
}

// ServiceCreateJSON takes a collection name, key and JSON object
// document and creates a new JSON object in the collection using
// the key.
func ServiceCreateJSON(cName string, key string, src []byte) error {
	if service == nil {
		return fmt.Errorf("Service not running")
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		err := sc.Collection.CreateJSON(key, src)
		sc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// ServiceReadJSON takes a collection name, key and returns a JSON object
// document.
func ServiceReadJSON(cName string, key string) ([]byte, error) {
	if service == nil {
		return nil, fmt.Errorf("Service not running")
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.ReadJSON(key)
	}
	return nil, fmt.Errorf("%q not available", cName)
}

// ServiceUpdateJSON takes a collection name, key and JSON object
// document and updates the collection.
func ServiceUpdateJSON(cName string, key string, src []byte) error {
	if service == nil {
		return fmt.Errorf("Service not running")
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		err := sc.Collection.UpdateJSON(key, src)
		sc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available", cName)
}

// ServiceDeleteJSON takes a collection name and key and removes
// and JSON object from the collection.
func ServiceDeleteJSON(cName string, key string) error {
	var (
		err error
	)
	if service == nil {
		return fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		err = sc.Collection.Delete(key)
		sc.Mutex.Unlock()
		return err
	}
	return fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameExists returns true if frame found in service collection,
// otherwise false
func ServiceFrameExists(cName string, fName string) bool {
	if service == nil {
		return false
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameExists(fName)
	}
	return false
}

// ServiceFrameCreate creates a frame in a service collection
func ServiceFrameCreate(cName string, fName string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	if service == nil {
		return nil, fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		f, err := sc.Collection.FrameCreate(fName, keys, dotPaths, labels, verbose)
		sc.Mutex.Unlock()
		return f, err
	}
	return nil, fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameObjects returns a JSON document of a copy of the objects in a frame for
// the service collection. It is analogous to a dataset.ReadJSON but for a frame's
// object list
func ServiceFrameObjects(cName string, fName string) ([]map[string]interface{}, error) {
	if service == nil {
		return nil, fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameObjects(fName)
	}
	return nil, fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameRefresh updates the frame object list's for the keys provided. Any new keys
//  cause a new object to be appended to the end of the list.
func ServiceFrameRefresh(cName string, fName string, keys []string, verbose bool) error {
	if service == nil {
		return fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameRefresh(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameReframe updates the frame object list. If a list of keys is provided then
// the object will be replaced with updated objects based on the keys provided.
func ServiceFrameReframe(cName string, fName string, keys []string, verbose bool) error {
	if service == nil {
		return fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameReframe(fName, keys, verbose)
	}
	return fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameClear clears the object and key list from a frame
func ServiceFrameClear(cName string, fName string) error {
	if service == nil {
		return fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameClear(fName)
	}
	return fmt.Errorf("%q not available from service", cName)
}

// ServiceFrameDelete deletes a frame from a service collection
func ServiceFrameDelete(cName string, fName string) error {
	if service == nil {
		return fmt.Errorf("Service not running, %q not available", cName)
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.FrameDelete(fName)
	}
	return fmt.Errorf("%q not available from service", cName)
}

// ServiceFrames returns a list of frame names in a service collection
func ServiceFrames(cName string) []string {
	if service == nil {
		return nil
	}
	if sc, found := service.collections[cName]; found {
		return sc.Collection.Frames()
	}
	return nil
}
