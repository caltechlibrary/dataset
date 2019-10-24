package dataset

import (
	"fmt"
	"sync"
)

/**
 * service.go provides the structures needed to embbed dataset
 * as a service via libdataset.go. E.g. to support writing
 * asynchronous web service written in Python via py_dataset.
 */

type ServiceCollection struct {
	Mutex      *sync.Mutex
	Collection *Collection
}

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
	} else {
		return fmt.Errorf("%q not found in service", cName)
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
	var (
		err error
	)
	if service == nil {
		return fmt.Errorf("Service not running")
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		err = sc.Collection.CreateJSON(key, src)
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
	var (
		err error
	)
	if service == nil {
		return fmt.Errorf("Service not running")
	}
	if sc, found := service.collections[cName]; found {
		sc.Mutex.Lock()
		err = sc.Collection.UpdateJSON(key, src)
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

// ServiceFrameExists returns true if frame found in service collection, otherwise false
// ServiceFrameCreate creates a frame in a service collection
// ServiceFrameObjects returns a JSON document of a copy of the objects in a frame for the service collection. It is analogous to a dataset.ReadJSON but for a frame's object list
// ServiceReframe updates the frame object list. If a list of keys is provided then the object will be replaced with updated objects based on the keys provided.
// ServiceFrameClear clears the object and key list from a frame
// ServiceFrameDelete deletes a frame from a service collection
// ServiceFrames returns a list of frame names in a service collection
