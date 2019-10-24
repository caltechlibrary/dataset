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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dotpath"
)

//
// NOTE: frame.go presents an Object as the native go type map[string]interface{} and DataFrame is intended to let you work with an ordered list of objects.
//

// DataFrame is the basic structure holding a list of objects as well as the definition
// of the list (so you can regenerate an updated list from a changed collection).
// It persists with the collection.
type DataFrame struct {
	// Explicit at creation
	Name string `json:"frame_name"`

	// CollectionName holds the name of the collection the frame was generated from. In theory you could
	// define a frame in one collection and use its results in another. A DataFrame can be rendered as a JSON
	// document.
	CollectionName string `json:"collection_name"`

	// DotPaths is a slice holding the definitions of what each Object attribute's data source is.
	DotPaths []string `json:"dot_paths"`

	// Labels are new attribute names for fields create from the provided
	// DotPaths.  Typically this is used to surface a deeper dotpath's
	// value as something more useful in the frame's context (e.g.
	// first_title from an array of titles might be labeled "title")
	Labels []string

	// NOTE: Keys is an orded list of object keys in the frame.
	Keys []string

	// NOTE: Object map privides a quick index by key to object index.
	objectMap map[string]interface{}

	// Created is the date the frame is originally generated and defined
	Created time.Time

	// Updated is the date the frame is updated (e.g. reframed)
	Updated time.Time
}

func (f *DataFrame) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["dot_paths"] = f.DotPaths
	m["labels"] = f.Labels
	m["keys"] = f.Keys
	m["object_list"] = f.Objects()
	m["created"] = f.Created
	m["updated"] = f.Updated
	return json.Marshal(m)
}

// frameObject takes a list of dot paths, labels and object key
// then generates a new object based on that.
func (c *Collection) frameObject(key string, dotPaths []string, labels []string) (map[string]interface{}, error) {
	errors := []string{}
	obj := map[string]interface{}{}
	err := c.Read(key, obj, false)
	if err != nil {
		return nil, err
	}
	o := map[string]interface{}{}
	for j, dpath := range dotPaths {
		value, err := dotpath.Eval(dpath, obj)
		if err == nil {
			key := labels[j]
			o[key] = value
		} else {
			errors = append(errors, fmt.Sprintf("%q path (%d) not found for key %q", dpath, j, key))
		}
	}
	if len(errors) > 0 {
		return o, fmt.Errorf("%s", strings.Join(errors, ", "))
	}
	return o, nil
}

// ObjectList (on a collection) takes a set of collection keys and builds
// an ordered array of objects from the array of keys, dot paths and
// labels provided.
func (c *Collection) ObjectList(keys []string, dotPaths []string, labels []string, verbose bool) ([]map[string]interface{}, error) {
	if len(dotPaths) != len(labels) {
		return nil, fmt.Errorf("dot paths and labels do not match")
	}
	pid := os.Getpid()
	objectList := make([]map[string]interface{}, len(keys))
	for i, key := range keys {
		obj, err := c.frameObject(key, dotPaths, labels)
		if verbose == true {
			if err != nil {
				log.Printf("(pid: %d) WARNING: framing error for key %q (%d), %s", pid, key, i, err)
			}
			if obj == nil {
				log.Printf("(pid: %d) WARNING: skipping object key %q (%d), object is nil", pid, key, i)
			}
		}
		if obj != nil {
			objectList = append(objectList, obj)
		}
		if verbose && (i > 0) && ((i % 1000) == 0) {
			log.Printf("(pid: %d) %d keys processed", pid, i)
		}
	}
	return objectList, nil
}

// hasFrame checks if a frame is defined already
func (c *Collection) hasFrame(key string) bool {
	if c.FrameMap == nil {
		return false
	}
	_, hasFrame := c.FrameMap[key]
	return hasFrame
}

// getFrame retrieves a frame by frame name from a collection.
func (c *Collection) getFrame(key string) (*DataFrame, error) {
	if c.FrameMap == nil {
		return nil, fmt.Errorf("no frames defined")
	}
	savedPath, ok := c.FrameMap[key]
	if ok == false {
		return nil, fmt.Errorf("frame %s not defined", key)
	}
	// read frame json from storage
	src, err := c.Store.ReadFile(path.Join(c.workPath, savedPath))
	if err != nil {
		return nil, err
	}
	// convert into DataFrame struct
	f := new(DataFrame)
	err = json.Unmarshal(src, &f)
	f.objectMap = make(map[string]interface{})
	// Update the internal objectMap attribute
	for i, key := range f.Keys {
		f.objectMap[key] = i
	}
	// return frame and error
	return f, err
}

// setFrame writes a DataFrame struct to the collection
func (c *Collection) setFrame(key string, f *DataFrame) error {
	// Check to see if we have a _frames directory to store our frames in
	if _, err := c.Store.Stat(path.Join(c.workPath, "_frames")); err != nil {
		if err := c.Store.MkdirAll(path.Join(c.workPath, "_frames"), 0775); err != nil {
			return err
		}
	}
	// Sanity check on frameName and collectionName
	f.CollectionName = c.Name
	f.Name = key

	// render DataFrame to JSON for storage
	src, err := json.Marshal(f)
	if err != nil {
		return err
	}
	// calculate the path to store the frame
	fName := key
	if strings.HasSuffix(fName, ".json") == false {
		fName = key + ".json"
	}
	savedPath := path.Join("_frames", fName)
	// update c.FrameMap with rel path to frame
	if c.FrameMap == nil {
		c.FrameMap = make(map[string]string)
	}
	c.FrameMap[key] = savedPath
	// save metadata and return
	err = c.Store.WriteFile(path.Join(c.workPath, savedPath), src, 0664)
	if err != nil {
		return err
	}
	err = c.SaveMetadata()
	return err
}

// rmFrame removes a frame from storage as well as from frames.json
func (c *Collection) rmFrame(key string) error {
	savedPath, ok := c.FrameMap[key]
	if ok == false {
		return fmt.Errorf("Frame %s not found", key)
	}
	delete(c.FrameMap, key)
	err := c.Store.Remove(path.Join(c.workPath, savedPath))
	err = c.SaveMetadata()
	return err
}

// FrameCreate takes a set of collection keys, dot paths and labels
// builds an ObjectList and assembles additional metadata returning
// a new Frame associated with the collection as well as an error value.
// If there is a mis-match in number of labels and dot paths an an error
// will be returned. If the frame already exists an error will be returned.
//
// Conceptually a frame is an ordered list of objects.  Frames are
// associated with a collection and the objects in a frame can
// easily be refreshed. Frames also serve as the basis for indexing
// a dataset collection and provide the data paths (expressed
// as a list of "dot paths"), labels (aka attribute names),
// and type information needed for indexing and search.
//
// If you need to update a frame's objects use Refresh(). If
// you need to change a frames object ordering use Reframe().
//
func (c *Collection) FrameCreate(name string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	// If frame exists return the existing frame
	if c.hasFrame(name) {
		return nil, fmt.Errorf("frame %q exists in %q", name, c.Name)
	}

	// Case of new Frame and with ObjectList
	if labels != nil && dotPaths != nil &&
		len(labels) != len(dotPaths) {
		return nil, fmt.Errorf("Mismatched dot paths and labels")
	}

	f := new(DataFrame)
	f.Name = name
	f.CollectionName = c.Name
	f.DotPaths = dotPaths[:]
	f.Labels = labels[:]
	f.Keys = []string{}
	f.objectMap = make(map[string]interface{})
	f.Created = time.Now()
	f.Updated = time.Now()

	// Populate our Object List
	pid := os.Getpid()
	for i, key := range keys {
		obj, err := c.frameObject(key, f.DotPaths, f.Labels)
		if verbose == true {
			if err != nil {
				log.Printf("(pid: %d) WARNING: framing error for key %q (%d), %s", pid, key, i, err)
			}
			if obj == nil {
				log.Printf("(pid: %d) WARNING: skipping object key %q (%d), object is nil", pid, key, i)
			}
		}
		if obj != nil {
			f.objectMap[key] = obj
			f.Keys = append(f.Keys, key)
		}
	}
	err := c.setFrame(name, f)
	return f, err
}

// FrameExists checkes to see if a frame is already defined.
// Returns true if it exists otherwise false
func (c *Collection) FrameExists(name string) bool {
	return c.hasFrame(name)
}

// Frames retrieves a list of available frames associated with a collection
func (c *Collection) Frames() []string {
	keys := []string{}
	if c.FrameMap == nil {
		return keys
	}
	for k := range c.FrameMap {
		keys = append(keys, k)
	}
	return keys
}

// Refresh updates the objects in the frame list
// from the related objects in the current collection
// without reordering the list of objects in the frame. If
// A new key is encountered the object is appended to the
// the frame's object list. No objects are removed from the
// frame. (non-destructive, only additive)
func (c *Collection) Refresh(name string, keys []string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	for i, key := range keys {
		obj, err := c.frameObject(key, f.DotPaths, f.Labels)
		if verbose == true {
			if err != nil {
				log.Printf("key %q (%d) frame error %s", key, i, err)
			}
			if obj == nil {
				log.Printf("key %q (%d) frame object is nil", key, i)
			}

		}
		if obj != nil {
			if _, ok := f.objectMap[key]; ok == true {
				f.objectMap[key] = obj
			} else {
				f.Keys = append(f.Keys, key)
				f.objectMap[key] = obj
			}
		} else {
			delete(f.objectMap, key)
			for i, fkey := range f.Keys {
				if fkey == key {
					f.Keys = append(f.Keys[:i], f.Keys[i+1:]...)
					break
				}
			}
		}
	}
	return c.setFrame(name, f)
}

// Reframe changes the order of a frame's objects and
// adds, updates or removes objects as needed from the frame
// based on the list of keys provided.
// (destructive update, will add or remove objects as needed)
func (c *Collection) Reframe(name string, keys []string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	// New Keys that will replace the values in f.Keys which are stale.
	nKeys := []string{}
	for _, key := range keys {
		obj, err := c.frameObject(key, f.DotPaths, f.Labels)
		if verbose == true {
			if err != nil {
				log.Printf("key %q frame error %s", key, err)
			}
			if obj == nil {
				log.Printf("key %q framed as nil object", key)
			}
		}
		if obj != nil {
			f.objectMap[key] = obj
			nKeys = append(nKeys, key)
		} else if _, ok := f.objectMap[key]; ok == true {
			// remove our stale object
			delete(f.objectMap, key)
		}
		// Figure out which objects to garbage collect
		for i, staleKey := range f.Keys {
			if key == staleKey {
				f.Keys = append(f.Keys[:i], f.Keys[i+1:]...)
			}
		}
	}
	// Now GC the objects in the stale key list
	for _, key := range f.Keys {
		if _, ok := f.objectMap[key]; ok == true {
			delete(f.objectMap, key)
		}
	}
	// Now update the Keys list with the new keys
	f.Keys = nKeys
	f.Updated = time.Now()
	return c.setFrame(name, f)
}

// SaveFrame saves a frame in a collection or returns an error
func (c *Collection) SaveFrame(name string, f *DataFrame) error {
	return c.setFrame(name, f)
}

// FrameClear empties the frame's object and key lists but
// leaves in place the Frame definition. Use Reframe()
// to re-populate a frame based on a new key list.
func (c *Collection) FrameClear(name string) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	// Emtpy the key and Object list.
	f.Keys = []string{}
	f.objectMap = make(map[string]interface{})
	return c.setFrame(name, f)
}

// FrameDelete removes a frame from a collection, returns an error if frame can't be deleted.
func (c *Collection) FrameDelete(name string) error {
	return c.rmFrame(name)
}

//
// The follow funcs define operations on the Frame struct.
//

// String renders the data structure DataFrame as JSON to a string
func (f *DataFrame) String() string {
	src, _ := json.MarshalIndent(f, "", "  ")
	return fmt.Sprintf("%s", src)
}

// Grid returns a Grid representaiton of a DataFrame's ObjectList
func (f *DataFrame) Grid(includeHeaderRow bool) [][]interface{} {
	rowCnt := len(f.Keys)
	colCnt := len(f.Labels)
	if includeHeaderRow == true {
		rowCnt++
	}
	rows := make([][]interface{}, rowCnt)
	if includeHeaderRow {
		header := make([]interface{}, colCnt)
		for i, val := range f.Labels {
			header[i] = val
		}
		rows[0] = header
	}
	// Now make reset of grid
	for i, key := range f.Keys {
		rowNo := i
		if includeHeaderRow == true {
			rowNo++
		}
		rows[rowNo] = make([]interface{}, colCnt)
		if obj, ok := f.objectMap[key]; ok == true {
			rec := obj.(map[string]interface{})
			for j, label := range f.Labels {
				if val, OK := rec[label]; OK == true {
					rows[rowNo][j] = val
				}
			}
		}
	}
	return rows
}

// Objects returns a copy of DataFrame's object list (an array of map[string]interface{})
func (f *DataFrame) Objects() []map[string]interface{} {
	ol := make([]map[string]interface{}, len(f.Keys))
	for _, key := range f.Keys {
		if obj, found := f.objectMap[key]; found == true && obj != nil {
			rec := obj.(map[string]interface{})
			ol = append(ol, rec)
		}
	}
	return ol
}
