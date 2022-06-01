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
	"github.com/caltechlibrary/dataset/pairtree"
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
	// "sqlite://", "mysql://" and the dsn conforming to the Golang
	// database/sql driver name in the database/sql package.
	DsnURI string `json:"dsn_uri,omitempty"`

	// Created
	Created string `json:"created,omitempty"`

	// PTStore the point to the pairtree implementation of storage
	PTStore *ptstore.Storage `json:"-"`
	// SQLStore points to a SQL database with JSON column support
	SQLStore *sqlstore.Storage `json:"-"`

	// Versioning holds the type of versioning implemented in the collection.
	// It can be set to an empty string (the default) which means no versioning.
	// It can be set to "patch" which means objects and attachments are versioned by
	// a semver patch value (e.g. 0.0.X where X is incremented), "minor" where
	// the semver minor value is incremented (e.g. e.g. 0.X.0 where X is incremented),
	// or "major" where the semver major value is incremented (e.g. X.0.0 where X is
	// incremented). Versioning affects storage of JSON objects and their attachments
	// across the whole collection.
	Versioning string `json:"versioning,omitempty"`

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
		switch c.Versioning {
		case "major":
			c.PTStore.SetVersioning(ptstore.Major)
		case "minor":
			c.PTStore.SetVersioning(ptstore.Minor)
		case "patch":
			c.PTStore.SetVersioning(ptstore.Patch)
		default:
			c.PTStore.SetVersioning(ptstore.None)
		}
	case SQLSTORE:
		c.SQLStore, err = sqlstore.Open(name, c.DsnURI)
	default:
		return nil, fmt.Errorf("failed to open %s, %q storage type not supported", name, c.StoreType)
	}
	return c, err
}

// Close closes a collection. For a pairtree that means flushing the
// keymap to disk. For a SQL store it means closing a database connection.
// Close is often called in conjunction with "defer" keyword.
//
// ```
//    c, err := dataset.Open("my_collection.ds")
//    if err != nil { /* .. handle error ... */ }
//    /* do some stuff with the collection */
//    defer func() {
//      if err := c.Close(); err != nil {
//         /* ... handle closing error ... */
//      }
//    }()
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
	}
	return fmt.Errorf("%q not supported", c.StoreType)
}

// SetVersioning sets the versioning on a collection. The version string
// can be "major", "minor", "patch". Any other value (e.g. "", "off", "none")
// will turn off versioning for the collection.
func (c *Collection) SetVersioning(versioning string) error {
	switch versioning {
	case "major":
		c.Versioning = versioning
	case "minor":
		c.Versioning = versioning
	case "patch":
		c.Versioning = versioning
	default:
		c.Versioning = ""
	}
	colName := path.Join(c.workPath, "collection.json")
	src, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot encode %q, %s", colName, err)
	}
	if err := ioutil.WriteFile(colName, src, 0660); err != nil {
		return fmt.Errorf("failed to create %q %s", colName, err)
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
		if err := os.MkdirAll(c.Name, 0770); err != nil {
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
	if os.MkdirAll(ptName, 0770); err != nil {
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
		if err := os.MkdirAll(c.Name, 0700); err != nil {
			return fmt.Errorf("cannot create collect %q, %s", c.Name, err)
		}
	} else {
		return fmt.Errorf("%q already exists", c.Name)
	}
	fullName := c.Name
	basename := path.Base(c.Name)
	colName := path.Join(fullName, "collection.json")
	c.Name = basename
	src, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("cannot encode %q, %s", colName, err)
	}
	c.Name = fullName
	if err := ioutil.WriteFile(colName, src, 0600); err != nil {
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
	cmName := path.Join(fullName, "codemeta.json")
	if err := ioutil.WriteFile(cmName, src, 0664); err != nil {
		return fmt.Errorf("failed to create %q, %s", cmName, err)
	}
	//NOTE: the collection's table needs to be created using the
	// SQLStore's Init method..
	c.SQLStore, err = sqlstore.Init(basename, c.DsnURI)
	return err
}

// Init - creates a new collection and opens it. It takes a name
// (e.g. directory holding the collection.json and codemeta.josn files)
// and an optional DSN in URI form. The default storage engine is a
// pairtree (i.e. PTSTORE) but some SQL storage engines are supported.
//
// If a DSN URI is a non-empty string then it is the SQL storage engine
// is used. The database and user access in the SQL engine needs be setup
// before you can successfully intialized your dataset collection.
// Currently three SQL database engines are support, SQLite3 or MySQL 8.
// You select the SQL storage engine by forming a URI consisting of a
// "protocol" (e.g. "sqlite", "mysql"), the protocol delimiter "://" and
// a Go SQL supported DSN based on the database driver implementation.
//
// A MySQL 8 DSN URI would look something like
//
//    `mysql://DB_USER:DB_PASSWD@PROTOCAL_EXPR/DB_NAME`
//
// The one for SQLite3
//
//     `sqlite://PATH_TO_DATABASE`
//
// NOTE: The DSN URI is stored in the collections.json.  The file should
// NOT be world readable as that will expose your database password. You
// can remove the DSN URI after initializing your collection but will then
// need to provide the DATASET_DSN_URI envinronment variable so you can
// open your database successfully.
//
// For PTSTORE the access value can be left blank.
//
// ```
//   var (
//      c *Collection
//      err error
//   )
//   name := "my_collection.ds"
//   c, err = dataset.Init(name, "")
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
//   c, err = dataset.Init(name, dsnURI)
//   if err != nil {
//     // ... handle error
//   }
//   defer c.Close()
// ```
//
func Init(name string, dsnURI string) (*Collection, error) {
	var (
		err error
	)
	c := new(Collection)
	c.Name = name
	c.DsnURI = dsnURI
	if dsnURI == "" {
		c.StoreType = PTSTORE
	} else {
		c.StoreType = SQLSTORE
	}
	switch c.StoreType {
	case PTSTORE:
		err = c.initPTStore()
	case SQLSTORE:
		err = c.initSQLStore()
	default:
		return nil, fmt.Errorf("%q storage type not supported", c.StoreType)
	}
	if err != nil {
		return nil, err
	}
	return Open(name)
}

// Metadata returns a copy of the codemeta.json file content found
// in the collection directory.
func (c *Collection) Metadata() ([]byte, error) {
	fName := path.Join(c.Name, "codemeta.json")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q, %s", fName, err)
	}
	return src, nil
}

// UpdateMetadata imports new codemeta citation information replacing
// the previous version. Collection must be open.
//
// ```
//   c, err := dataset.Open("my_collection.ds")
//   if err != nil { /* ... handle error ... */ }
//   defer c.Close()
//   c.UpdateCodemeta("codemeta.json")
// ```
func (c *Collection) UpdateMetadata(fName string) error {
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return fmt.Errorf("failed to read %q, %s", fName, err)
	}
	if err := ioutil.WriteFile(path.Join(c.Name, "codemeta.json"), src, 0664); err != nil {
		return fmt.Errorf("failed to write %q, %s", path.Join(c.Name, "codemeta.json"), err)
	}
	return nil
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
			return c.SQLStore.Create(key, src)
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

// Versions retrieves a list of versions available for a JSON document if
// versioning is enabled for the collection.
//
// ```
//   key, version := "123", "0.0.1"
//   if versions, err := Versions(key); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) Versions(key string) ([]string, error) {
	var (
		versions []string
		err      error
	)
	switch c.StoreType {
	case PTSTORE:
		versions, err = c.PTStore.Versions(key)
	case SQLSTORE:
		versions, err = c.SQLStore.Versions(key)
	default:
		return nil, fmt.Errorf("%q not supported", c.StoreType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s, %s", key, err)
	}
	return versions, err
}

// ReadVersion retrieves a specific vesion from the collection for the given object.
//
// ```
//   key, version := "123", "0.0.1"
//   var obj map[string]interface{}
//
//   if err := ReadVersion(key, version, &obj); err != nil {
//      ...
//   }
// ```
//
func (c *Collection) ReadVersion(key string, version string, obj map[string]interface{}) error {
	var (
		src []byte
		err error
	)
	switch c.StoreType {
	case PTSTORE:
		src, err = c.PTStore.ReadVersion(key, version)
	case SQLSTORE:
		src, err = c.SQLStore.ReadVersion(key, version)
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

// Delete removes an object from the collection (this includes all versions and all attachments)
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

// HasKey takes a collection and checks if a key exists.
// NOTE: collection must be open otherwise false will always be returned.
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
