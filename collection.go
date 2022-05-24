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
	"path"

	// Dataset sub-modules
	"github.com/caltechlibrary/dataset/ptstore"
	"github.com/caltechlibrary/dataset/sqlstore"
)

const (
	//
	// Store types
	//

	// PTSTORE describes the storage type using a pairtree
	PTSTORE = "pairtree"

	// SQLSTORE describes the SQL storage type
	SQLSTORE = "sqlstore"
)

// Collection is the holds both metadata and associated metadata
// for collection level operations on collections of JSON objects.
type Collection struct {
	// DatasetVersion of the collection
	DatasetVersion string `json:"dataset,omitempty"`

	// Name of collection
	Name string `json:"name"`

	// StoreType can be either "pairtree" (default or if attribute is
	// omitted) or "sqlstore".  If sqlstore the connection string, DSN,
	// will determine the type of SQL database being accessed.
	StoreType string `json:"storage_type,omitempty"`

	// DSN holds the filename containing the DSN used to connect
	// to the SQL database.  If blank the DSN value will be read from
	// the environment via `os.Getenv("DSN")`.
	DSN string `json:"dsn,omitempty"`

	//
	// General Metadata for collection.
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
	// ??Store points to the storage system used in the collection.
	//

	// PTStore the point to the pairtree implementation of storage
	PTStore *ptstore.Storage `json:"-"`
	// SQLStore points to a SQL database with JSON column support
	SQLStore *sqlstore.Storage `json:"-"`
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
// Public interface for dataset
//

// Open reads in a collection's metadata and returns
// a new collection structure and error value. The collection
// structure includes a storage object that conforms to the
// StorageSystem interface (e.g. ptstore or sqlstore).
//
// ```
//    var (
//       c *Collection
//       err error
//    )
//    c, err = dataset.Open("collection.ds", dataset.PTSTORE)
//    if err != nil {
//       // ... handle error
//    }
//    defer c.Close()
// ```
//
// SQL store
//
// ```
//    var (
//       c *Collection
//       err error
//    )
//    // a sql/database DSN
//    dsn := ...
//    c, err = dataset.Open(dsn, dataset.SQLSTORE)
//    if err != nil {
//       // ... handle error
//    }
//    defer c.Close()
// ```
//
func Open(name string, dsn string, storeType string) (*Collection, error) {
	// NOTE: find the collection.json file then
	// open the appropriate store.
	src, err := ioutil.ReadFile(path.Join(name, "collection.json"))
	if err != nil {
		return nil, err
	}
	c := new(Collection)
	if err := json.Unmarshal(src, &c); err != nil {
		return nil, err
	}
	c.Name = name
	c.DSN = dsn
	c.StoreType = storeType
	switch storeType {
	case PTSTORE:
		c.PTStore, err = ptstore.Open(c.Name)
	case SQLSTORE:
		c.SQLStore, err = sqlstore.Open(c.DSN)
	default:
		return nil, fmt.Errorf("failed to open %s, %q storage type not supported", name, c.StoreType)
	}
	return c, err
}

// Close closes a collection. For a pairtree that means flushing
// keys to disk. for a SQL store it means closing a database connection.
// Close is often called in conjunction with "defer" keyword.
//
// ```
//    c, err := dataset.Open("my_collection.ds", dataset.PTSTORE)
//    if err != nil { /* .. handle error ... */ }
//    /* do some stuff with the collection */
//    if err := c.Close(); err != nil {
//       /* ... handle closing error ... */
//    }
// ```
//
func (c *Collection) Close() error {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Close()
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Close()
		}
	default:
		return fmt.Errorf("%q not supported")
	}
	return fmt.Errorf("nothing to close")
}

// initPTStore takes a *Collection and initializes a PTSTORE collection.
// For pairtrees this means create the directory structure and writing
// out the collection.json file, a skeleton codemeta.json and an empty
// keys.json file.
func (c *Collection) initPTStore() error {
	return fmt.Errorf("InitPTStore() not implemented")
}

// InitSQLStore takes a *Collection and initializes the SQL database
// tables for a collection.
func (c *Collection) initSQLStore() error {
	//
	// NOTE: if the access value is left blank for  SQLstore then
	// the DSN is assumed to be held in the environment and retrieved with
	// `os.Getenv("DATASETD_DSN")`. That needs to get assigned to the
	// running collection but NOT get written to disk in the
	// collection.json file.
	//
	return fmt.Errorf("InitSQLStore() not implemented")
}

// Init - creates a new collection and opens it. It takes a name
// (e.g. directory holding the collection.json file), an access name
// (e.g. a file holding a DSN) and the storage type.
//
// For PTSTORE the access value can be left blank.
//```
//   var (
//      c *Collection
//      err error
//   )
//   name := "my_collection.ds"
//   c, err = dataset.Init(name, "", dataset.PTSTORE)
//   if err != nil {
//     // ... handle error
//   }
//   defer c.Close()
//```
//
// For a sqlstore collection we need to pass the "access" value. This
// is the file containing a DNS or environment variables formating a DSN.
//
//```
//   var (
//      c *Collection
//      err error
//   )
//   name := "my_collection"
//   dsn := os.Getenv("DATASETD_DSN")
//   c, err = dataset.Init(name, dsn, dataset.SQLSTORE)
//   if err != nil {
//     // ... handle error
//   }
//   defer c.Close()
//```
//
// NOTE: if the access value is left blank for  SQLstore then
// the DSN is assumed to be held in the environment and retrieved with
// `os.Getenv("DATASETD_DSN")`.
//
func Init(name string, dsn string, storeType string) (*Collection, error) {
	var err error
	if storeType == "" {
		storeType = PTSTORE
	}
	c := new(Collection)
	c.Name = name
	c.DSN = dsn
	c.StoreType = storeType
	switch storeType {
	case PTSTORE:
		err = c.initPTStore()
	case SQLSTORE:
		err = c.initSQLStore()
	default:
		return nil, fmt.Errorf("%q storage type not supported", storeType)
	}
	return c, err
}

// Close closes a collection, writing the updated keys to disc
// Close removes the "lock.pid" file in the collection root.
// Close is often called in conjunction with "defer" keyword.
//
// ```
//    c, err := dataset.Open("my_collection.ds")
//    if err != nil { /* .. handle error ... */ }
//    /* do some stuff with the collection */
//    if err := c.Close(); err != nil {
//       /* ... handle closing error ... */
//    }
// ```
//

// ImportCodemeta imports codemeta citation information updating
// the collections metadata. Collection must be open.
//
// ```
//   c, err := dataset.Open("my_collection.ds")
//   if err != nil { /* ... handle error ... */ }
//   defer c.Close()
//   c.ImportCodemeta("codemeta.json")
// ```
func ImportCodemeta(fName string) error {
	return fmt.Errorf("ImportCodemeta(%q) is not implemented")
}

//
// The following are aliases to the storage system implementation.
//

// Create store a an object in the collection. Object will get
// converted to JSON source then stored. Collection must be open.
//
// ```
//   key := "123"
//   obj := map[]*interface{}{ "one": 1, "two": 2 }
//   if err := c.Create(key, obj); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) Create(key string, obj *interface{}) error {
	src, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", err)
	}
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Create(key, src)
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.PTStore.Create(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// Read retrieves an object from the collection, unmarshals it and
// updates the object pointed to by obj.
//
// ```
//   var obj map[string]interface{}
//
//   key := "123"
//   if err := c.Read(key, &obj); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) Read(key string, obj *interface{}) error {
	var (
		src []byte
		err error
	)
	switch c.StoreType {
	case PTSTORE:
		src, err = c.PTStore.Read(key)
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	if err != nil {
		return fmt.Errorf("failed to read %s, %s", key, err)
	}
	return json.Unmarshal(src, &obj)
}

// Update updates an existing JSON document.
//
// ```
//   key := "123"
//   obj["three"] = 3
//   if err := c.Update(key, obj); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) Update(key string, obj *interface{}) error {
	src, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", err)
	}
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Update(key, src)
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Update(key, src)
		}
	default:
		return fmt.Errorf("%q not supported")
	}
	return fmt.Errorf("%s not open", c.Name)
}

// Delete removes an object frmo the collection
//
// ```
//   key := "123"
//   if err := c.Delete(key); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) Delete(key string) error {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Delete(key)
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Delete(key)
		}
	default:
		return fmt.Errorf("%q not supported")
	}
	return fmt.Errorf("%s not open", c.Name)
}

// List returns a array of strings holding all the keys
// in the collection.
//
// ```
//   keys, err := c.List()
//   for _, key := range keys {
//      ...
//   }
// ```
//
func (c *Collection) List() ([]string, error) {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.List()
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.List()
		}
	default:
		return nil, fmt.Errorf("%q not supported")
	}
	return nil, fmt.Errorf("%s not open", c.Name)
}
