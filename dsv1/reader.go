//
// dsv1 is a submodule of dataset package.
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
package dsv1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	// Version 1.1.1-legacy
	Version = `1.1.1-legacy`

	// Asc is used to identify ascending sorts
	Asc = iota
	// Desc is used to identify descending sorts
	Desc = iota

	// internal virtualize column name format string
	fmtColumnName = `column_%03d`
)

// Collection is the container holding a pairtree containing JSON docs
type Collection struct {
	// DatasetVersion of the collection
	DatasetVersion string `json:"dataset,omitempty"`

	// Name (filename) of collection
	Name string `json:"name"`

	// workPath holds the path (i.e. non-protocol and hostname, in URI)
	workPath string // `json:"-"`

	// KeyMap holds the document key to path in the collection
	KeyMap map[string]string `json:"keymap,omitempty"`

	// FrameMap is a list of frame names and with rel path to the frame defined in the collection
	FrameMap map[string]string `json:"frames,omitempty"`

	//
	// Metadata for collection.
	//

	// Description describes what is in the collection.
	Description string `json:"description,omitempty"`

	// Created is the date/time the init command was run in
	// RFC1123 format.
	Created string `json:"created,omitempty"`

	// Version of collection being stored in semvar notation
	Version string `json:"version,omitempty"`

	// Contact info
	Contact string `json:"contact,omitempty"`

	// Author holds a list of PersonOrOrg
	Author []*PersonOrOrg `json:"author,omitempty"`

	// Contributors holds a list of PersonOrOrg
	Contributor []*PersonOrOrg `json:"contributor,omitempty"`

	// Funder holds a list of PersonOrOrg
	Funder []*PersonOrOrg `json:"funder,omitempty"`

	// DOI holds the digital object identifier if defined.
	DOI string `json:"doi,omitempty"`

	// License holds a pointer to the license information for
	// the collection. E.g. CC0 URL
	License string `json:"license,omitempty"`

	// Annotation is a map to any addition metadata associated with
	// the Collection's metadata.
	Annotation map[string]interface{} `json:"annotation,omitempty"`

	//
	// The following are the Namaste fields are depreciated
	// they are left in place so we can easily migrate their
	// content into more appropriate fields or annotations.
	//

	// Who is the person(s)/organization(s) that created the collection
	Who []string `json:"who,omitempty"`
	// What - description of collection
	What string `json:"what,omitempty"`
	// When - date associated with collection (e.g. 2021,
	// 2021-10, 2021-10-02), should map to an approx date like in
	// archival work.
	When string `json:"when,omitempty"`
	// Where - location (e.g. URL, address) of collection
	Where string `json:"where,omitempty"`

	//
	// private attributes, experiments in performance tuning
	//

	// unsafeSaveMetadata is set to true when doing batch object
	// operations so we can save writing collection.json on each
	// create or delete.
	unsafeSaveMetadata bool
}

// PersonOrOrg holds a the description of a person or organizaion
// associated with the dataset collection. e.g. author, contributor
// or funder.
type PersonOrOrg struct {
	// Type is either "Person" or "Organization"
	Type string `json:"@type,omitempty"`
	// ID is either an ORCID or ROR
	ID string `json:"@id,omitempty"`
	// Name of an organization, empty if person
	Name string `json:"name,omitempty"`
	// Given name for a person, empty of organization
	GivenName string `json:"givenName,omitempty"`
	// Family name for a person, empty of organization
	FamilyName string `json:"familyName,omitempty"`
	// Affiliation holds the intitution affiliation of a person.
	Affiliation []*PersonOrOrg `json:"affiliation,omitempty"`
	// Annotation holds custom fields, e.g. a grant number of a funder
	Annotation map[string]interface{} `json:"annotation,omitempty"`
}

//
// internal utility functions
//

// normalizeKeyName() trims leading and trailing spaces
func normalizeKeyName(s string) string {
	// NOTE: As of 1.0.1 release, keys are normalized to lowercase.
	return strings.ToLower(strings.TrimSpace(s))
}

// collectionNameAsPath takes a uri and normalizes collection name
// to a path
func collectionNameAsPath(p string) string {
	if strings.Contains(p, "://") {
		u, _ := url.Parse(p)
		return u.Path
	}
	return strings.TrimSpace(p)
}

// keyAndFName converts a key (which may have things like slashes) into a disc friendly name and key value
func keyAndFName(name string) (string, string) {
	if strings.HasSuffix(name, ".json") == true {
		return name, url.QueryEscape(name)
	}
	return name, url.QueryEscape(name) + ".json"
}

// localizePairPath checks if map has value and adjusts value
// to localized OS path separator (e.g. for Windows).
func localizePairPath(key string, m map[string]string) (string, bool) {
	value, ok := m[key]
	if ok && (os.PathSeparator != '/') {
		parts := strings.Split(value, "/")
		value = path.Join(parts...)
	}
	return value, ok
}

// saveMetadata writes the collection's metadata to c.workPath
func (c *Collection) saveMetadata() error {
	if c.unsafeSaveMetadata == true {
		// NOTE: We're playing fast and loose with the collection metadata, skip saveMetadata().
		return nil
	}
	// Check to see if collection exists, if not create it!
	if _, err := os.Stat(c.workPath); err != nil {
		if err := os.MkdirAll(c.workPath, 0775); err != nil {
			return err
		}
	}
	// Make sure pair paths in c.KeyMap are encoded POSIX style
	for key, value := range c.KeyMap {
		if strings.Contains(value, "\\") {
			c.KeyMap[key] = strings.ReplaceAll(value, "\\", "/")
		}
	}
	src, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Can't marshal metadata, %s", err)
	}
	if err := os.WriteFile(path.Join(c.workPath, "collection.json"), src, 0664); err != nil {
		return fmt.Errorf("Can't store collection metadata, %s", err)
	}
	return nil
}

// deleteCollection an entire collection
func deleteCollection(name string) error {
	_, err := os.Stat(name)
	if err != nil {
		return err
	}
	collectionName := collectionNameAsPath(name)
	if err := os.RemoveAll(collectionName); err != nil {
		return err
	}
	return nil
}

//
// Public interface for dataset
//

// Open reads in a collection's metadata and returns
// and new collection structure or error.
//
// ```
//    var (
//       c *Collection
//       err error
//    )
//    c, err = dataset.Open("collection.ds")
//    if err != nil {
//       // ... handle error
//    }
//    defer c.Close()
// ```
//
func Open(name string) (*Collection, error) {
	_, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	collectionName := collectionNameAsPath(name)
	src, err := os.ReadFile(path.Join(collectionName, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err != nil {
		return nil, err
	}
	//NOTE: we need to reset collectionName so we're working with a path useable to get to the JSON documents.
	c.Name = path.Base(collectionName)
	c.workPath = collectionName
	if c.KeyMap == nil {
		c.KeyMap = make(map[string]string)
	}

	return c, nil
}

// DocPath returns a full path to a key or an error if not found
//
// ```
//    c, err := dataset.Open("my_collection.ds")
//    if err != nil {
//       // ... handle error ...
//    }
//    defer c.Close()
//    key := "my-object-key"
//    docPath := c.DocPath(key)
// ```
func (c *Collection) DocPath(name string) (string, error) {
	name = normalizeKeyName(name)
	keyName, name := keyAndFName(name)
	if p, ok := localizePairPath(keyName, c.KeyMap); ok == true {
		return path.Join(c.workPath, p, name), nil
	}
	return "", fmt.Errorf("Can't find %q", name)
}

// Close closes a collection, writing the updated keys to disc
// Close removes the "lock.pid" file in the collection root.
// Close is often called in conjunction with "defer" keyword.
//
// ```
//    c, err := dataset.Open("my_collection.ds")
//    if err != nil { // .. handle error ...
//    }
//    // do some stuff with the collection
//    if err := c.Close(); err != nil {
//       // ... handle closing error ...
//    }
// ```
//
func (c *Collection) Close() error {
	// Cleanup c so it can't accidentally get reused
	if c != nil {
		//dName := c.workPath // collectionNameAsPath(c.Name)
		lockName := path.Join(c.workPath, "lock.pid")
		if err := os.Remove(lockName); err != nil {
			return fmt.Errorf("could not remove %s, %s", lockName, err)
		}
		c.Name = ""
		c.workPath = ""
		c.KeyMap = map[string]string{}
	}
	return nil
}

// IsKeyNotFound checks an error message and returns true if
// it is a key not found error.
func (c *Collection) IsKeyNotFound(e error) bool {
	if strings.Compare(e.Error(), "key not found") == 0 {
		return true
	}
	return false
}

// ReadJSON finds a the record in the collection and
// returns the JSON source or an error.
//
// ```
//    var (
//       c *Collection
//    )
//    // ... collection previously opened and assigned to "c" ...
//    key := "object-1"
//    src, err := c.ReadJSON(key)
//    if err != nil {
//       // ... handle error ...
//    }
//    // ... do something with the JSON encoded "src" value ...
// ```
//
func (c *Collection) ReadJSON(name string) ([]byte, error) {
	name = normalizeKeyName(name)
	// Handle potentially URL encoded names
	keyName, FName := keyAndFName(name)
	pairPath, ok := localizePairPath(keyName, c.KeyMap)
	if ok != true {
		return nil, fmt.Errorf("%q does not exist in %s", keyName, c.Name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	// we need to join c.Name + bucketName + name to get path do JSON document
	src, err := os.ReadFile(path.Join(c.workPath, pairPath, FName))
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Read finds the record in a collection, updates the data
// interface provide and if problem returns an error
// name must exist or an error is returned
//
// ```
//    var (
//       c *dataset.Collection
//    )
//    // ... collection previously opened and assigned to "c" ...
//    key := "object-2"
//    obj, err := c.Read(key)
//    if err != nil { // ... handle error ...
//    }
// ```
func (c *Collection) Read(name string, data map[string]interface{}, cleanObject bool) error {
	src, err := c.ReadJSON(name)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if cleanObject == true {
		delete(data, "_Key")
		delete(data, "_Attachments")
	}
	return nil
}

// Keys returns a list of keys in a collection
//    var (
//       c *dataset.Collection
//       keys []string
//    )
//    // ... collection previously opened and assigned to "c" ...
//
//    keys := c.Keys()
//    for _, key := range keys {
//       // ... do something with the list of keys ...
//    }
// ```
//
func (c *Collection) Keys() []string {
	keys := []string{}
	for k := range c.KeyMap {
		keys = append(keys, k)
	}
	return keys
}

// KeyExists returns true if key is in collection's KeyMap, false otherwise
//
//    var (
//       c *dataset.Collection
//    )
//    // ... collection previously opened and assigned to "c" ...
//
//    key := "object-1"
//    if c.KeyExists(key) == true {
//       // ... do something with the key ...
//    }
// ```
//
func (c *Collection) KeyExists(key string) bool {
	_, hasKey := c.KeyMap[normalizeKeyName(key)]
	return hasKey
}

// Length returns the number of keys in a collection
//
//    var (
//       c *dataset.Collection
//    )
//    // ... collection previously opened and assigned to "c" ...
//
//    l := c.Length()
//    // ... do something with the number of itemsin the collection ...
// ```
//
func (c *Collection) Length() int {
	return len(c.KeyMap)
}

// IsCollection checks to see if a given path contains a
// collection.json file
func IsCollection(p string) bool {
	finfo, err := os.Stat(path.Join(p, "collection.json"))
	if (err == nil) && (finfo.IsDir() == false) {
		return true
	}
	return false
}

// MetadataJSON() returns a collection's metadata fields as a
// JSON encoded byte array.
func (c *Collection) MetadataJSON() []byte {
	meta := new(Collection)
	meta.DatasetVersion = c.DatasetVersion
	meta.Name = c.Name
	meta.Description = c.Description
	meta.DOI = c.DOI
	meta.Created = c.Created
	meta.Version = c.Version
	meta.Contact = c.Contact
	if c.Author != nil {
		meta.Author = []*PersonOrOrg{}
		for _, obj := range c.Author {
			meta.Author = append(meta.Author, obj)
		}
	}
	if c.Contributor != nil {
		meta.Contributor = []*PersonOrOrg{}
		for _, obj := range c.Contributor {
			meta.Contributor = append(meta.Contributor, obj)
		}
	}
	if c.Funder != nil {
		meta.Funder = []*PersonOrOrg{}
		for _, obj := range c.Funder {
			meta.Funder = append(meta.Funder, obj)
		}
	}
	meta.License = c.License
	if c.Annotation != nil {
		meta.Annotation = map[string]interface{}{}
		for key, value := range c.Annotation {
			meta.Annotation[key] = value
		}
	}
	src, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		src = []byte{}
	}
	return src
}
