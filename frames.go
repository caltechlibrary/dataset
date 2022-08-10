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
	"os"
	"path"
	"sort"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/v2/dotpath"
)

//
// NOTE: frames.go presents an Object as the native go type map[string]interface{} and DataFrame is intended to let you work with an ordered list of objects.
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
	Labels []string `json:"labels"`

	// NOTE: Keys is an orded list of object keys in the frame.
	Keys []string `json:"keys"`

	// NOTE: Object map privides a quick index by key to object index.
	ObjectMap map[string]interface{} `json:"object_map"`

	// Created is the date the frame is originally generated and defined
	Created time.Time `json:"created"`

	// Updated is the date the frame is updated (e.g. reframed)
	Updated time.Time `json:"updated"`
}

// frameObject takes a list of dot paths, labels and object key
// then generates a new object based on that.
func (c *Collection) frameObject(key string, dotPaths []string, labels []string) (map[string]interface{}, error) {
	errors := []string{}
	obj := map[string]interface{}{}
	err := c.Read(key, obj)
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
			errors = append(errors, fmt.Sprintf("path %q not found, %q in %q for %+v", dpath, key, c.Name, obj))
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
//
// ```
//   var mapList []map[string]interface{}
//
//   keys := []string{ "123", "124", "125" }
//   dotPaths := []string{ ".title", ".description" }
//   labels := []string{ "Title", "Description" }
//   verbose := true
//   mapList, err = c.ObjectList(keys, dotPaths, labels, verbose)
// ```
//
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
				log.Printf("(pid: %d) WARNING: framing error, %s", pid, err)
			}
			if obj == nil {
				log.Printf("(pid: %d) WARNING: skipping object %q (%d) in %q, object is nil", pid, key, i, c.Name)
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

// HasFrame checks if a frame is defined already. Collection needs
// to previously been opened.
//
// ```
//   frameName := "journals"
//   if c.HasFrame(frameName) {
//      ...
//   }
// ```
//
func (c *Collection) HasFrame(frameName string) bool {
	framePath := path.Join(c.workPath, "_frames",
		path.Base(frameName)+".json")
	if _, err := os.Stat(framePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// getFrame retrieves a frame by frame name from a collection.
func (c *Collection) getFrame(key string) (*DataFrame, error) {
	framePath := path.Join(c.workPath, "_frames", key+".json")
	src, err := ioutil.ReadFile(framePath)
	if err != nil {
		return nil, err
	}
	f := new(DataFrame)
	if err := json.Unmarshal(src, &f); err != nil {
		return nil, err
	}
	// Double check if we have a bad object_map?
	if f.ObjectMap == nil {
		f.ObjectMap = map[string]interface{}{}
	}
	// return frame and error
	return f, err
}

// setFrame writes a DataFrame struct to the collection
func (c *Collection) setFrame(key string, f *DataFrame) error {
	// Check to see if we have a _frames directory to store our frames in
	if _, err := os.Stat(path.Join(c.workPath, "_frames")); err != nil {
		if err := os.MkdirAll(path.Join(c.workPath, "_frames"), 0775); err != nil {
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
	// save metadata and return
	err = os.WriteFile(path.Join(c.workPath, savedPath), src, 0664)
	if err != nil {
		return err
	}
	return nil
}

// rmFrame removes a frame from storage as well as from frames.json
func (c *Collection) rmFrame(key string) error {
	framePath := path.Join(c.workPath, "_frames", key+".json")
	if _, err := os.Stat(framePath); os.IsNotExist(err) {
		return fmt.Errorf("frame %q does not exist", key)
	}
	return os.RemoveAll(framePath)
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
// If you need to update a frame's objects use FrameRefresh(). If
// you need to change a frame's objects or ordering use FrameReframe().
//
// ```
//   frameName := "journals"
//   keys := []string{ "123", "124", "125" }
//   dotPaths := []string{ ".title", ".description" }
//   labels := []string{ "Title", "Description" }
//   verbose := true
//   frame, err := c.FrameCreate(frameName, keys, dotPaths, labels, verbose)
//   if err != nil {
//      ...
//   }
// ```
//
func (c *Collection) FrameCreate(name string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	frameDir := path.Join(c.workPath, "_frames")
	if _, err := os.Stat(frameDir); os.IsNotExist(err) {
		os.MkdirAll(frameDir, 0775)
	}
	// If frame exists return the existing frame
	if c.HasFrame(name) {
		return nil, fmt.Errorf("frame %q exists in %q", name, c.workPath)
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
	f.ObjectMap = make(map[string]interface{})
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
				log.Printf("(pid: %d) WARNING: skipping object for key %q (%d), object is nil", pid, key, i)
			}
		}
		if obj != nil {
			f.ObjectMap[key] = obj
			f.Keys = append(f.Keys, key)
		}
	}
	framePath := path.Join(c.workPath, "_frames", name+".json")
	src, err := json.Marshal(f)
	if err != nil {
		return f, fmt.Errorf("failed to encode fame %q, %s", name, err)
	}
	if err := ioutil.WriteFile(framePath, src, 0666); err != nil {
		return f, fmt.Errorf("failed to write frame %q, %s", name, err)
	}
	return f, nil
}

// Frames retrieves a list of available frame names associated with a
// collection.
//
// ```
//   frameNames := c.FrameNames()
//   for _, name := range frames {
//      // do something with each frame name
//      objects, err := c.FrameObjects(name)
//      ...
//   }
// ```
//
func (c *Collection) FrameNames() []string {
	framesDir := path.Join(c.workPath, "_frames")
	files, err := os.ReadDir(framesDir)
	if err != nil {
		return []string{}
	}
	keys := []string{}
	for _, file := range files {
		keys = append(keys, strings.TrimSuffix(path.Base(file.Name()), ".json"))
	}
	sort.Strings(keys)
	return keys
}

// FrameKeys retrieves a list of keys assocaited with a data frame
//
// ```
//   frameName := "journals"
//   keys := c.FrameKeys(frameName)
// ```
//
func (c *Collection) FrameKeys(name string) []string {
	frame, err := c.FrameRead(name)
	if err != nil {
		return []string{}
	}
	return frame.Keys
}

// FrameRead retrieves a frame from a collection.
// Returns the DataFrame and an error value
//
// ```
//   frameName := "journals"
//   data, err := c.FrameRead(frameName)
//   if err != nil {
//      ..
//   }
// ```
//
func (c *Collection) FrameRead(name string) (*DataFrame, error) {
	return c.getFrame(name)
}

// FrameDef retrieves the frame definition returns a
// a map string interface.
//
// ```
//   definition := map[string]interface{}{}
//   frameName := "journals"
//   definition, err := c.FrameDef(frameName)
//   if err != nil {
//      ..
//   }
// ```
//
func (c *Collection) FrameDef(name string) (map[string]interface{}, error) {
	frame, err := c.FrameRead(name)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{
		"name":      frame.Name,
		"dot_paths": frame.DotPaths,
		"labels":    frame.Labels,
	}
	return m, nil
}

// FrameRefresh updates a DataFrames' object list based on the
// existing keys in the frame.  It doesn't change the order of objects.
// It is used when objects in a collection that are included in the
// frame have been updated. It uses the frame's existing definition.
//
// NOTE: If an object is missing in the collection it gets pruned from
// the object list.
//
// ```
//   frameName, verbose := "journals", true
//   err := c.FrameRefresh(frameName, verbose)
//   if err != nil {
//      ...
//   }
// ```
//
func (c *Collection) FrameRefresh(name string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	for i, key := range f.Keys {
		obj, err := c.frameObject(key, f.DotPaths, f.Labels)
		if verbose == true {
			if err != nil {
				log.Printf("key %q (%d) frame error %s", key, i, err)
			}
			if obj == nil {
				log.Printf("key %q (%d) frame object is nil", key, i)
			}
		}
		if err != nil {
			if verbose {
				log.Printf("WARNING could not read %q from %q", key, c.workPath)
			}
			continue
		}
		if obj != nil {
			if _, ok := f.ObjectMap[key]; ok == false {
				f.Keys = append(f.Keys, key)
			}
			f.ObjectMap[key] = obj
		} else {
			// Remove the stale object
			delete(f.ObjectMap, key)
			for i, fkey := range f.Keys {
				if fkey == key {
					// Remove the stale key
					f.Keys = append(f.Keys[:i], f.Keys[i+1:]...)
					break
				}
			}
		}
	}
	return c.setFrame(name, f)
}

// FrameReframe **replaces** a frame's object list based on the
// keys provided. It uses the frame's existing definition.
//
// ```
//   frameName, verbose := "journals", false
//   keys := ...
//   err := c.FrameReframe(frameName, keys, verbose)
//   if err != nil {
//      ...
//   }
// ```
//
func (c *Collection) FrameReframe(name string, keys []string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	// GC the stale objects
	f.Keys = []string{}
	f.ObjectMap = make(map[string]interface{})

	// New Keys that will replace the values in f.Keys which are stale.
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
			f.ObjectMap[key] = obj
			f.Keys = append(f.Keys, key)
		}
	}
	// Now update the Keys list with the new keys
	f.Updated = time.Now()
	return c.setFrame(name, f)
}

// SaveFrame saves a frame in a collection or returns an error
//
// ```
//    frameName := "journals"
//    data, err := c.FrameRead(frameName)
//    if err != nil {
//       ...
//    }
//    // do stuff with the frame's data
//       ...
//    // Save the changed frame data
//    err = c.SaveFrame(frameName, data)
// ```
//
func (c *Collection) SaveFrame(name string, f *DataFrame) error {
	return c.setFrame(name, f)
}

// FrameClear empties the frame's object and key lists but
// leaves in place the Frame definition. Use Reframe()
// to re-populate a frame based on a new key list.
//
// ```
//   frameName := "journals"
//   err := c.FrameClear(frameName)
//   if err != nil  {
//      ...
//   }
//
func (c *Collection) FrameClear(name string) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	// Emtpy the key and Object list.
	f.Keys = []string{}
	f.ObjectMap = make(map[string]interface{})
	return c.setFrame(name, f)
}

// FrameDelete removes a frame from a collection, returns an error
// if frame can't be deleted.
//
// ```
//   frameName := "journals"
//   err := c.FrameDelete(frameName)
//   if err != nil {
//      ...
//   }
// ```
//
func (c *Collection) FrameDelete(name string) error {
	return c.rmFrame(name)
}

// FrameObjects returns a copy of a DataFrame's object list given a
// collection's frame name.
//
// ```
//    var (
//      err error
//      objects []map[string]interface{}
//    )
//    frameName := "journals"
//    objects, err = c.FrameObjects(frameName)
//    if err != nil  {
//       ...
//    }
// ```
//
func (c *Collection) FrameObjects(fName string) ([]map[string]interface{}, error) {
	f, err := c.FrameRead(fName)
	if err != nil {
		return nil, err
	}
	ol := f.Objects()
	return ol, nil
}

//
// The follow funcs define operations on the Frame struct.
//

// String renders the data structure DataFrame as JSON to a string
//
//  ```
//   frameName := "journals"
//   data, err := c.FrameRead(frameName)
//   if err != nil {
//      ...
//   }
//   fmt.Printf("\n%s\n", data.String())
//  ```
//
func (f *DataFrame) String() string {
	src, _ := json.MarshalIndent(f, "", "  ")
	return fmt.Sprintf("%s", src)
}

// Grid returns a table representaiton of a DataFrame's ObjectList
//
// ```
//   frameName, includeHeader := "journals", true
//   data, err := c.FrameRead(frameName)
//   if err != nil {
//      ...
//   }
//   rows, err := data.Grid(includeHeader)
//   if err != nil {
//      ...
//   }
//   ... /* now do something with the rows */ ...
// ```
//
func (f *DataFrame) Grid(includeHeaderRow bool) [][]interface{} {
	rowCnt := len(f.Keys)
	colCnt := len(f.Labels)
	if includeHeaderRow == true {
		rowCnt++
	}
	rows := [][]interface{}{}
	if includeHeaderRow {
		header := make([]interface{}, colCnt)
		for i, val := range f.Labels {
			header[i] = val
		}
		rows = append(rows, header)
	}
	// Now make reset of grid
	objectList := f.Objects()
	for i, obj := range objectList {
		rowNo := i
		if includeHeaderRow == true {
			rowNo++
		}
		row := make([]interface{}, colCnt)
		for colNo, label := range f.Labels {
			if val, OK := obj[label]; OK == true {
				row[colNo] = val
			} else {
				row[colNo] = ""
			}
		}
		if len(row) > 0 {
			rows = append(rows, row)
		}
	}
	return rows
}

// Objects returns a copy of DataFrame's object list (an array of map[string]interface{})
//
// ```
//   frameName := "journals"
//   data, err := c.FrameRead(frameName)
//   if err != nil {
//      ...
//   }
//   objectList, err := data.Objects()
//   if err != nil {
//      ...
//   }
// ```
//
func (f *DataFrame) Objects() []map[string]interface{} {
	ol := []map[string]interface{}{}
	for _, key := range f.Keys {
		if obj, found := f.ObjectMap[key]; found == true && obj != nil {
			rec := obj.(map[string]interface{})
			ol = append(ol, rec)
		}
	}
	return ol
}
