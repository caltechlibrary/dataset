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
	"os"
	"path"
	"time"

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

// Collection is the holds both operational metadata
// for collection level operations on collections of JSON objects.
// General metadata is stored in a codemeta.json file in the root
// directory along side the collection.json file.
type Collection struct {
	// DatasetVersion of the collection
	DatasetVersion string `json:"dataset,omitempty"`

	// Name of collection
	Name string `json:"name"`

	// StoreType can be either "pairtree" (default or if attribute is
	// omitted) or "sqlstore".  If sqlstore the connection string, DSN URI,
	// will determine the type of SQL database being accessed.
	StoreType string `json:"storage_type,omitempty"`

	// DsnURI holds protocol plus dsn string. The protocol can be
	// "sqlite:", "mysql:" and the dsn string would conform to the Golang
	// database/sql driver for the SQL database.
	//
	// If blank the DSN value will be read from
	// the environment via `os.Getenv("DATASET_DSN_URI")`.
	//
	DsnURI string `json:"dsn_uri,omitempty"`

	// Created
	Created string `json:"created,omitempty"`

	// PTStore the point to the pairtree implementation of storage
	PTStore *ptstore.Storage `json:"-"`
	// SQLStore points to a SQL database with JSON column support
	SQLStore *sqlstore.Storage `json:"-"`

	//
	// Private varibles
	//

	// workPath holds the path the directory where the collection.json
	// file is found.
	workPath string `json:"-"`
}

//
// Public interface for dataset
//

// Open reads in a collection's operational metadata and returns
// a new collection structure and error value.
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
	c.workPath = name
	if c.DsnURI == "" {
		c.DsnURI = os.Getenv("DATASET_DSN_URI")
	}
	switch c.StoreType {
	case PTSTORE:
		c.PTStore, err = ptstore.Open(name, c.DsnURI)
	case SQLSTORE:
		c.SQLStore, err = sqlstore.Open(name, c.DsnURI)
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
//    c, err := dataset.Open("my_collection.ds")
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
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return nil
}

// initPTStore takes a *Collection and initializes a PTSTORE collection.
// For pairtrees this means create the directory structure and writing
// out the collection.json file, a skeleton codemeta.json and an empty
// keys.json file.
func (c *Collection) initPTStore() error {
	now := time.Now()
	today := now.Format(datestamp)
	c.DatasetVersion = Version
	c.Created = now.UTC().Format(timestampUTC)
	// Split see if c.Name path exists
	if _, err := os.Stat(c.Name); os.IsNotExist(err) {
		// Create directoriess if needed
		if err := os.MkdirAll(c.Name, 0775); err != nil {
			return fmt.Errorf("cannot create collect %q, %s", c.Name, err)
		}
	} else {
		return fmt.Errorf("%q already exists", c.Name)
	}
	fullName := c.Name
	basename := path.Base(c.Name)
	colName := path.Join(c.Name, "collection.json")
	c.Name = basename
	src, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot encode %q, %s", colName, err)
	}
	c.Name = fullName
	if err := ioutil.WriteFile(colName, src, 0664); err != nil {
		return fmt.Errorf("failed to create %q %s", colName, err)
	}
	// Create a default codemeta.json file in the directory
	//today := time.Now().Format(datestamp)
	m := map[string]string{
		"{today}":    today,
		"{c_name}":   path.Base(c.Name),
		"{app_name}": path.Base(os.Args[0]),
		"{version}":  Version,
	}
	src = BytesProcessor(m, []byte(`{
    "@context": "https://doi.org/10.5063/schema/codemeta-2.0",
    "@type": "SoftwareSourceCode",
    "dateCreated": "{today}",
    "name": "{c_name}",
    "version": "0.0.0",
    "releaseNotes": "This is a {app_name} {version} collection",
    "developmentStatus": "concept",
    "softwareRequirements": [
        "https://github.com/caltechlibrary/dataset"
    ]
}`))
	cmName := path.Join(c.Name, "codemeta.json")
	if err := ioutil.WriteFile(cmName, src, 0664); err != nil {
		return fmt.Errorf("failed to create %q, %s", cmName, err)
	}
	// Create the pairtree root
	ptName := path.Join(c.Name, "pairtree")
	if os.MkdirAll(ptName, 0775); err != nil {
		return fmt.Errorf("failed to create pairtree root, %s", err)
	}
	return nil
}

// initSQLStore initializes a new SQL based dataset storage system. It
// presumes that the database and been create and an appropriate
// database user has been created outside the dataset provided tooling.
// It uses a DNS in URI form where the "protocol" element identifies
// the type of SQL database, e.g. sqlite would use "sqlite:", MySQL
// would use "mysql:". The rest of the URI is formed from a Go style
// DSN (data source name). These are SQL system specific but usually
// include things like db name, user, password to access the database.
//
// A collection using SQL database for storage is split into one table
// per collection. There is a "_collection" table which holds collection
// wide properties, e.g. collection name, create, read, update, delete
// properties, versioning status, a copy of the codemeta.json document
// in a JSON column, the name of the collection in a column,
//
// InitSQLStore takes a *Collection and initializes the SQL database
// tables for a collection.
//
// NOTE: A SQL store still has a dataset named directory containing
// both the collection.json and codemeta.json file but it lacks a
// pairtree since that is where the object will be stored. Any
// attachments, if allowed, are stored in an S3 like bucket (e.g. via
// minio).
func (c *Collection) initSQLStore() error {
	now := time.Now()
	today := now.Format(datestamp)
	c.DatasetVersion = Version
	c.Created = now.UTC().Format(timestampUTC)
	// Split see if c.Name path exists
	if _, err := os.Stat(c.Name); os.IsNotExist(err) {
		// Create directoriess if needed
		if err := os.MkdirAll(c.Name, 0775); err != nil {
			return fmt.Errorf("cannot create collect %q, %s", c.Name, err)
		}
	} else {
		return fmt.Errorf("%q already exists", c.Name)
	}
	fullName := c.Name
	basename := path.Base(c.Name)
	colName := path.Join(c.Name, "collection.json")
	c.Name = basename
	src, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot encode %q, %s", colName, err)
	}
	c.Name = fullName
	if err := ioutil.WriteFile(colName, src, 0660); err != nil {
		return fmt.Errorf("failed to create %q %s", colName, err)
	}

	// NOTE: if the access value is left blank for  SQLstore then
	// the DSN URI is assumed to be held in the environment and
	// retrieved with `os.Getenv("DATASET_DSN_URI")`. That needs
	// to get assigned to the running collection but NOT get
	// written to disk in the collection.json file.
	if c.DsnURI == "" {
		c.DsnURI = os.Getenv("DATASET_DSN_URI")
	}

	// Create a default codemeta.json file in the directory
	//today := time.Now().Format(datestamp)
	m := map[string]string{
		"{today}":    today,
		"{c_name}":   path.Base(c.Name),
		"{app_name}": path.Base(os.Args[0]),
		"{version}":  Version,
	}
	src = BytesProcessor(m, []byte(`{
    "@context": "https://doi.org/10.5063/schema/codemeta-2.0",
    "@type": "SoftwareSourceCode",
    "dateCreated": "{today}",
    "name": "{c_name}",
    "version": "0.0.0",
    "releaseNotes": "This is a {app_name} {version} collection",
    "developmentStatus": "concept",
    "softwareRequirements": [
        "https://github.com/caltechlibrary/dataset"
    ]
}`))
	cmName := path.Join(c.Name, "codemeta.json")
	if err := ioutil.WriteFile(cmName, src, 0664); err != nil {
		return fmt.Errorf("failed to create %q, %s", cmName, err)
	}
	return nil
}

// Init - creates a new collection and opens it. It takes a name
// (e.g. directory holding the collection.json file), an access name
// (e.g. a file holding a DSN URI) and the storage type.
//
// For PTSTORE the access value can be left blank.
//
// ```
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
// ```
//
// For a sqlstore collection we need to pass the "access" value. This
// is the file containing a DNS or environment variables formating a DSN.
//
// ```
//   var (
//      c *Collection
//      err error
//   )
//   name := "my_collection"
//   dsnURI := os.Getenv("DATASET_DSN_URI")
//   c, err = dataset.Init(name, dsnURI, dataset.SQLSTORE)
//   if err != nil {
//     // ... handle error
//   }
//   defer c.Close()
// ```
//
// NOTE: if the dsnURI value is left blank for SQLstore then
// the dsnURI is assumed to be held in the environment and retrieved with
// `os.Getenv("DATASET_DSN_URI")`. A URI is formed by prefix a DNS
// (data source name) with a protocol, e.g. "sqlite:", "mysql:". Everything
// after the protocal is assumed to form a valid DSN for the given Go
// SQL driver.
//
func Init(name string, dsnURI string, storeType string) (*Collection, error) {
	var err error
	if storeType == "" {
		storeType = PTSTORE
	}
	c := new(Collection)
	c.Name = name
	c.DsnURI = dsnURI
	c.StoreType = storeType
	switch storeType {
	case PTSTORE:
		err = c.initPTStore()
	case SQLSTORE:
		err = c.initSQLStore()
	default:
		return nil, fmt.Errorf("%q storage type not supported", storeType)
	}
	if err != nil {
		return nil, err
	}
	if c.DsnURI == "" {
		c.DsnURI = os.Getenv("DATASET_DSN_URI")
	}
	return Open(name)
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
	return fmt.Errorf("ImportCodemeta(%q) is not implemented", fName)
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
func (c *Collection) Create(key string, obj map[string]interface{}) error {
	src, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
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
func (c *Collection) Read(key string, obj map[string]interface{}) error {
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
	return DecodeJSON(src, &obj)
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
func (c *Collection) Update(key string, obj map[string]interface{}) error {
	src, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
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
		return fmt.Errorf("%q not supported", c.StoreType)
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
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// Keys returns a array of strings holding all the keys
// in the collection.
//
// ```
//   keys, err := c.Keys()
//   for _, key := range keys {
//      ...
//   }
// ```
//
func (c *Collection) Keys() ([]string, error) {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Keys()
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Keys()
		}
	default:
		return nil, fmt.Errorf("%q not supported", c.StoreType)
	}
	return nil, fmt.Errorf("%s not open", c.Name)
}

// HasKey takes a collection and checks if a key exists. NOTE:
// collection must be open otherwise false will always be returned.
//
// ```
//   key := "123"
//   if c.HasKey(key) {
//      ...
//   }
// ```
//
func (c *Collection) HasKey(key string) bool {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.HasKey(key)
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.HasKey(key)
		}
	default:
		return false
	}
	// If we got here the collection isn't open ...
	return false
}

// Length returns the number of objects in a collection
// NOTE: Returns a -1 (as int64) on error, e.g. collection not open
// or Length not available for storage type.
func (c *Collection) Length() int64 {
	switch c.StoreType {
	case PTSTORE:
		if c.PTStore != nil {
			return c.PTStore.Length()
		}
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Length()
		}
	}
	return int64(-1)
}
