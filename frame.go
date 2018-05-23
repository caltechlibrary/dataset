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
	// Explicit at creator
	Name           string          `json:"frame_name"`
	CollectionName string          `json:"collection_name"`
	DotPaths       []string        `json:"dot_paths"`
	Keys           []string        `json:"keys"`
	Grid           [][]interface{} `json:"grid"`
	Created        time.Time       `json:"created"`
	Updated        time.Time       `json:"updated,omitempty"`

	// Derived or explicitly set after creation
	Labels      []string `json:"labels,omitempty"`
	ColumnTypes []string `json:"column_types,omitempty"`
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
	src, err := c.Store.ReadFile(path.Join(c.Name, savedPath))
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
	if _, err := c.Store.Stat(path.Join(c.Name, "_frames")); err != nil {
		if err := c.Store.MkdirAll(path.Join(c.Name, "_frames"), 0775); err != nil {
			return err
		}
	}

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
	err = c.Store.WriteFile(path.Join(c.Name, savedPath), src, 0664)
	if err != nil {
		return err
	}
	err = c.saveMetadata()
	return err
}

// rmFrame removes a frame from storage as well as from frames.json
func (c *Collection) rmFrame(key string) error {
	savedPath, ok := c.FrameMap[key]
	if ok == false {
		return fmt.Errorf("Frame %s not found", key)
	}
	delete(c.FrameMap, key)
	err := c.Store.Remove(path.Join(c.Name, savedPath))
	err = c.saveMetadata()
	return err
}

// Frame takes a set of collection keys and dotpaths, builds a grid and assembles
// the grid and metadata returning a new CollectionFrame and error. Frames are
// assoicated with the collection.
func (c *Collection) Frame(name string, keys []string, dotPaths []string, verbose bool) (*DataFrame, error) {
	// Read an existing frame or return error
	if len(keys) == 0 || len(dotPaths) == 0 {
		return c.getFrame(name)
	}
	// Case of new Frame Build our Grid.

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
	g, err := c.Grid(keys, dotPaths, verbose)
	if err != nil {
		return nil, err
	}
	f.Grid = g
	labels := []string{}
	colTypes := []string{}
	// NOTE: derive labels from dotPaths and default column types to string
	for _, p := range dotPaths {
		l := dotpath.ToLabel(p)
		labels = append(labels, l)
		// Set a default column type of string
		colTypes = append(colTypes, "string")
	}
	f.Labels = labels[:]
	// FIXME: Derive column types from grid values
	f.ColumnTypes = colTypes[:]
	err = c.setFrame(name, f)
	return f, err
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
	if len(keys) > 0 {
		f.Keys = keys[:]
	}
	f.Updated = time.Now()
	g, err := c.Grid(f.Keys, f.DotPaths, verbose)
	if err != nil {
		return err
	}
	f.Grid = g
	return c.setFrame(name, f)
}

// FrameLabels sets the labels for a frame, the number of labels must match the number of dot paths (columns) in the frame.
func (c *Collection) FrameLabels(name string, labels []string) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	if len(f.DotPaths) != len(labels) {
		return fmt.Errorf("number of columns (%d) does not match the number of labels (%d)", len(f.DotPaths), len(labels))
	}
	f.Labels = labels[:]
	f.Updated = time.Now()
	return c.setFrame(name, f)
}

// FrameTypes sets the types associated with a frame's columns, types list must match the number of dot paths (columns) in the frame.
func (c *Collection) FrameTypes(name string, columnTypes []string) error {
	f, err := c.getFrame(name)
	if err != nil {
		return err
	}
	if len(f.DotPaths) != len(columnTypes) {
		return fmt.Errorf("number of columns (%d) does not match the number of column types (%d)", len(f.DotPaths), len(columnTypes))
	}
	f.ColumnTypes = columnTypes[:]
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
