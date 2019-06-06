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
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dotpath"
)

type DataFrame struct {
	// Explicit at creation
	Name           string   `json:"frame_name"`
	CollectionName string   `json:"collection_name"`
	DotPaths       []string `json:"dot_paths"`
	// NOTE: Keys should hold the same values as column zero of the grid.
	// Keys controls the order of rows in a grid when reframing.
	Keys []string        `json:"keys"`
	Grid [][]interface{} `json:"grid"`

	// NOTE: Objects is a replacement for Grid, it is an objects
	// which base on use from Python and shell is easier to work with
	// accurately then a 2D array which usually leads to at leats two
	// or more inner loops in scripts.
	ObjectList []map[string]interface{} `json:"object_list"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated,omitempty"`

	// NOTE: these values effect how Reframe works
	AllKeys    bool   `json:"use_all_keys"`
	FilterExpr string `json:"filter_expr,omitempty"`
	SortExpr   string `json:"sort_expr,omitempty"`
	SampleSize int    `json:"sample_size"`

	// Derived or explicitly set after creation
	Labels []string `json:"labels,omitempty"`
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

// Frame takes a set of collection keys and dotpaths, builds a grid and assembles
// the grid and metadata returning a new CollectionFrame and error. Frames are
// associated with the collection and can be re-generated.
func (c *Collection) Frame(name string, keys []string, dotPaths []string, verbose bool) (*DataFrame, error) {
	// If frame exists return the existing frame
	if c.hasFrame(name) {
		return c.getFrame(name)
	}

	// Case of new Frame and building our Grid.

	// NOTE: we need to enforce that column zero is explicitly ._Key
	hasKeyColumn := false
	for _, key := range dotPaths {
		if key == "._Key" {
			hasKeyColumn = true
			break
		}
	}
	if hasKeyColumn == false {
		dotPaths = append(dotPaths, "")
		copy(dotPaths[1:], dotPaths)
		dotPaths[0] = "._Key"
	}

	f := new(DataFrame)
	f.Name = name
	f.CollectionName = c.Name
	f.DotPaths = dotPaths[:]
	f.Keys = keys[:]
	f.Created = time.Now()
	f.Updated = time.Now()

	// NOTE: derive labels from dotPaths and
	// default column types to string
	labels := []string{}
	for _, p := range dotPaths {
		l := dotpath.ToLabel(p)
		labels = append(labels, l)
	}
	f.Labels = labels[:]

	// Populate our ObjectList
	ol, err := c.ObjectList(keys, dotPaths, labels, verbose)
	if err != nil {
		return nil, err
	}
	f.ObjectList = ol

	// Populate our Grid (Grid is depreciated, RSD 2019-06-06)
	g, err := c.Grid(keys, dotPaths, verbose)
	if err != nil {
		return nil, err
	}
	f.Grid = g
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
	for k, _ := range c.FrameMap {
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
	// NOTE: we need to enforce that column zero is explicitly ._Key
	hasKeyColumn := false
	for _, key := range f.DotPaths {
		if key == "._Key" {
			hasKeyColumn = true
			break
		}
	}
	if hasKeyColumn == false {
		dotPaths := f.DotPaths
		// Update DotPaths
		dotPaths = append(dotPaths, "")
		copy(dotPaths[1:], dotPaths)
		dotPaths[0] = "._Key"
		f.DotPaths = dotPaths[:]
		// Update Labels
		labels := f.Labels
		labels = append(labels, "")
		copy(labels[1:], labels)
		labels[0] = "_Key"
		f.Labels = labels[:]
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
	// NOTE: ObjectList is replacing Grid, RSD 2019-06-06
	ol, err := c.ObjectList(f.Keys, f.DotPaths, f.Labels, verbose)
	if err != nil {
		return err
	}
	f.ObjectList = ol
	// NOTE: Grid is depreciated, RSD 2019-06-06
	g, err := c.Grid(f.Keys, f.DotPaths, verbose)
	if err != nil {
		return err
	}
	f.Grid = g
	return c.setFrame(name, f)
}

// SaveFrame saves a frame in a collection or returns an error
func (c *Collection) SaveFrame(name string, f *DataFrame) error {
	return c.setFrame(name, f)
}

// FrameLabels sets the labels for a frame, the number of labels
// must match the number of dot paths (columns) in the frame.
// NOTE: FrameLabels will cause the ObjectList to be regenerated from
// the current state of the collection.
func (c *Collection) FrameLabels(name string, labels []string, verbose bool) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	if len(f.DotPaths) != len(labels) {
		return fmt.Errorf("number of columns (%d) does not match the number of labels (%d)", len(f.DotPaths), len(labels))
	}
	f.Labels = labels[:]
	// NOW we need to regenerate our ObjectList
	ol, err := c.ObjectList(f.Keys, f.DotPaths, f.Labels, verbose)
	if err != nil {
		return err
	}
	f.ObjectList = ol
	// NOTE: we need to refenerate our Grid to match (this is depreciated, RSD 2019-06-06)
	// NOTE: Grid is depreciated, RSD 2019-06-06
	g, err := c.Grid(f.Keys, f.DotPaths, verbose)
	if err != nil {
		return err
	}
	f.Grid = g

	f.Updated = time.Now()
	return c.setFrame(name, f)
}

// DeleteFrame removes a frame from a collection, returns an error if frame can't be deleted.
func (c *Collection) DeleteFrame(name string) error {
	return c.rmFrame(name)
}

// String renders the data structure DataFrame as JSON to a string
func (f *DataFrame) String() string {
	src, _ := json.MarshalIndent(f, "", "  ")
	return fmt.Sprintf("%s", src)
}
