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
// NOTE: frame.go presents an Object as the native go type map[string]interface{} and
// ObjectList as a slice of Objects.
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

	// NOTE: Keys should hold the same values as column zero of the grid.
	// Keys controls the order of rows in a grid when reframing.
	Keys []string `json:"keys"`

	// NOTE: ObjectList is a replacement for Grid. This representaition more
	// closely mirrors data frames as used in Python and R. It also can help
	// avoid an inner loop on iteration because we don't need to track a relative
	// index to get a "cell" value from a column heading.
	ObjectList []map[string]interface{} `json:"object_list"`

	// Created is the date the frame is originally generated and defined
	Created time.Time `json:"created"`

	// Updated is the date the frame is updated (e.g. reframed)
	Updated time.Time `json:"updated,omitempty"`

	// AllKeys is a flag used to define a frame as operating over an entire collection,
	// this allows for simplier update.  NOTE: this value effects how Reframe works.
	AllKeys bool `json:"use_all_keys"`

	// FilterExpr is a the expression used to filter a collections keys to determine
	// how a frame is "reframed".  It generally is faster to create your key list outside
	// the frame but that approach has the disadvantage of not persisting with the frame.
	// NOTE: this value effects how Reframe works.
	FilterExpr string `json:"filter_expr,omitempty"`

	// SortExpr holds the sort expression so it persists with the frame. Often you can
	// get a faster sort outside the frame but that comes at a disadvantage of not being
	// persisted with the frame. NOTE: this value effects how Reframe works.
	SortExpr string `json:"sort_expr,omitempty"`
	// SampleSize is used to hold a frame intended to be a sample. It is used when re-generating
	// the same. NOTE: this value effects how Reframe works.
	SampleSize int `json:"sample_size"`

	// Labels are derived from the DotPaths provided but can be replaced without changing
	// the dotpaths. Typically this is used to surface a deeper dotpath's value as something more
	// useful in the frame's context (e.g. first_title from an array of titles might be labeled "title")
	Labels []string `json:"labels,omitempty"`
}

// ObjectList (on a collection) takes a set of collection keys and builds an array
// of objects (i.e. map[string]interface{}) from the array of keys, dot paths and
// labels provided.
func (c *Collection) ObjectList(keys []string, dotPaths []string, labels []string, verbose bool) ([]map[string]interface{}, error) {
	if len(dotPaths) != len(labels) {
		return nil, fmt.Errorf("dot paths and labels do not match")
	}
	pid := os.Getpid()
	objectList := make([]map[string]interface{}, len(keys))
	for i, key := range keys {
		rec := map[string]interface{}{}
		err := c.Read(key, rec, false)
		if err != nil {
			return nil, err
		}
		objectList[i] = make(map[string]interface{})
		for j, dpath := range dotPaths {
			value, err := dotpath.Eval(dpath, rec)
			if err == nil {
				key := labels[j]
				objectList[i][key] = value
			} else if verbose == true {
				log.Printf("(pid: %d) WARNING: skipped key %s, path %s for row %d and column %d, %s", pid, key, dpath, i, j, err)
			}
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

// Frame takes a set of collection keys, dotpaths and labels
// builds an ObjectList and assembles metadata returning
// a new CollectionFrame and error. Frames are
// associated with the collection and can be re-generated.
// If the length of labels and dotpaths mis-match an error will be
// returned. If the frame already exists the definition is NOT
// UPDATED and the existing frame is returned. If you need to
// update a frame use ReFrame().
func (c *Collection) Frame(name string, keys []string, dotPaths []string, labels []string, verbose bool) (*DataFrame, error) {
	// If frame exists return the existing frame
	if c.hasFrame(name) {
		return c.getFrame(name)
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
	f.Keys = keys[:]
	f.Created = time.Now()
	f.Updated = time.Now()

	// Populate our ObjectList
	ol, err := c.ObjectList(keys, dotPaths, labels, verbose)
	if err != nil {
		return nil, err
	}
	f.ObjectList = ol

	err = c.setFrame(name, f)
	return f, err
}

// HasFrame checkes to see if a frame is already defined.
func (c *Collection) HasFrame(name string) bool {
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

// Reframe will re-generate contents of a frame based on the current records in a collection.
// If a list of keys is supplied then the regenerated frame will be based on the new set of keys provided
func (c *Collection) Reframe(name string, keys []string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		if len(f.SortExpr) > 0 {
			keys, err = c.KeySortByExpression(keys, f.SortExpr)
			if err != nil {
				return err
			}
		}
		f.Keys = keys
	}
	f.Updated = time.Now()
	// NOTE: ObjectList is replaced Grid, RSD 2019-06-24, v0.0.64
	ol, err := c.ObjectList(f.Keys, f.DotPaths, f.Labels, verbose)
	if err != nil {
		return err
	}
	f.ObjectList = ol
	return c.setFrame(name, f)
}

// SaveFrame saves a frame in a collection or returns an error
func (c *Collection) SaveFrame(name string, f *DataFrame) error {
	return c.setFrame(name, f)
}

// DeleteFrame removes a frame from a collection, returns an error if frame can't be deleted.
func (c *Collection) DeleteFrame(name string) error {
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
	rowCnt := len(f.ObjectList)
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
	for i, rec := range f.ObjectList {
		rowNo := i
		if includeHeaderRow == true {
			rowNo++
		}
		rows[rowNo] = make([]interface{}, colCnt)
		for j, label := range f.Labels {
			if val, OK := rec[label]; OK == true {
				rows[rowNo][j] = val
			}
		}
	}
	return rows
}

// Objects returns a copy of DataFrame.ObjectList (array of map[string]interface{})
func (f *DataFrame) Objects() []map[string]interface{} {
	return f.ObjectList[:]
}
