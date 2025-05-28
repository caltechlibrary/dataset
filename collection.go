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
package dataset

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/models"

)

const (
	//
	// Store types
	//

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

	// StoreType is alwasy "sqlstore". v3 does not support "pairtree". 
	// It remains for backward compatibility purposes.
	StoreType string `json:"storage_type,omitempty"`

	// DsnURI holds protocol plus dsn string. The protocol can be
	// "sqlite://", dsn conforming to the Golang database/sql driver name in the database/sql package.
	// As of v3 this string should be "sqlite://collection.db".
	DsnURI string `json:"dsn_uri,omitempty"`

	// Created
	Created string `json:"created,omitempty"`

	// Repaired
	Repaired string `json:"repaired,omitempty"`

	// SQLStore points to a SQL database with JSON column support
	SQLStore *SQLStore `json:"-"`

	// History determines of objects are written to the history table on
	// create, update and delete. By default it should be true.
	History bool `json:"history,omitempty"`
	
	//
	// Models is a placeholder for eventually integrating the models package support
	//
	Model *models.Model `json:"-"`

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
//
//	var (
//	   c *dataset.Collection
//	   err error
//	)
//	c, err = dataset.Open("collection.ds")
//	if err != nil {
//	   // ... handle error
//	}
//	defer c.Close()
//
// ```
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
	fullPath, err := filepath.Abs(name)
	if err == nil {
		c.workPath = fullPath
	} else {
		c.workPath = name
	}
	if c.DsnURI == "" {
		c.DsnURI = os.Getenv("DATASET_DSN_URI")
	}
	switch c.StoreType {
	case SQLSTORE:
		c.SQLStore, err = SQLStoreOpen(c.workPath, c.DsnURI)
	default:
		return nil, fmt.Errorf("failed to open %s, %q storage type not supported", name, c.StoreType)
	}
	return c, err
}

// Close closes a collection. For a SQL store it means closing a database connection.
// Close is often called in conjunction with "defer" keyword.
//
// ```
//
//	c, err := dataset.Open("my_collection.ds")
//	if err != nil { /* .. handle error ... */ }
//	/* do some stuff with the collection */
//	defer func() {
//	  if err := c.Close(); err != nil {
//	     /* ... handle closing error ... */
//	  }
//	}()
//
// ```
func (c *Collection) Close() error {
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Close()
		}
	}
	return fmt.Errorf("%q not supported", c.StoreType)
}

// WorkPath returns the working path to the collection.
func (c *Collection) WorkPath() string {
	return c.workPath
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
	src, err := JSONMarshal(c)
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
	c.SQLStore, err = SQLStoreInit(fullName, c.DsnURI)
	return err
}

// Init - creates a new collection and opens it. It takes a name
// (e.g. directory holding the collection.json and codemeta.josn files)
// and an optional DSN in URI form. The default storage engine is SQLite3.
//
// For a sqlstore collection we need to pass the "access" value. This
// is the file containing a DNS or environment variables formating a DSN.
//
// ```
//
//	var (
//	   c *Collection
//	   err error
//	)
//	name := "my_collection.ds"
//	dsnURI := "sqlite://collection.db"
//	c, err = dataset.Init(name, dsnURI)
//	if err != nil {
//	  // ... handle error
//	}
//	defer c.Close()
//
// ```
func Init(name string, dsnURI string) (*Collection, error) {
	var (
		err error
	)
	c := new(Collection)
	c.Name = name
	c.DsnURI = dsnURI
	if dsnURI == "" {
		c.StoreType = SQLSTORE
		c.DsnURI = "sqlite://collection.db"
	} else {
		c.StoreType = SQLSTORE
	}
	switch c.StoreType {
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

// Codemeta returns a copy of the codemeta.json file content found
// in the collection directory. The collection must be previous open.
//
// ```
//
//	name := "my_collection.ds"
//	c, err := dataset.Open(name)
//	if err != nil {
//	   ...
//	}
//	defer c.Close()
//	src, err := c.Metadata()
//	if err != nil {
//	   ...
//	}
//	ioutil.WriteFile("codemeta.json", src, 664)
//
// ```
func (c *Collection) Codemeta() ([]byte, error) {
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
//
//	name := "my_collection.ds"
//	codemetaFilename := "../codemeta.json"
//	c, err := dataset.Open(name)
//	if err != nil {
//	   ...
//	}
//	defer c.Close()
//	c.UpdateMetadata(codemetaFilename)
//
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
// A Go `map[string]interface{}` is a common way to handle ad-hoc
// JSON data in gow. Use `CreateObject()` to store structured
// data.
//
// ```
//
//	key := "123"
//	obj := map[]*interface{}{ "one": 1, "two": 2 }
//	if err := c.Create(key, obj); err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) Create(key string, obj map[string]interface{}) error {
	src, err := JSONMarshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
	}
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Create(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// CreateObject is used to store structed data in a dataset collection.
// The object needs to be defined as a Go struct notated approriately
// with the domain markup for working with json.
//
// ```
//
//	import (
//	  "encoding/json"
//	  "fmt"
//	  "os"
//	)
//
//	type Record struct {
//	    ID string `json:"id"`
//	    Name string `json:"name,omitempty"`
//	    EMail string `json:"email,omitempty"`
//	}
//
//	func main() {
//	    c, err := dataset.Open("friends.ds")
//	    if err != nil {
//	         fmt.Fprintf(os.Stderr, "%s", err)
//	         os.Exit(1)
//	    }
//	    defer c.Close()
//
//	    obj := &Record{
//	        ID: "mojo",
//	        Name: "Mojo Sam",
//	        EMail: "mojo.sam@cosmic-cafe.example.org",
//	    }
//	    if err := c.CreateObject(obj.ID, obj); err != nil {
//	         fmt.Fprintf(os.Stderr, "%s", err)
//	         os.Exit(1)
//	    }
//	    fmt.Printf("OK\n")
//	    os.Exit(0)
//	}
//
// ```
func (c *Collection) CreateObject(key string, obj interface{}) error {
	src, err := JSONMarshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
	}
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Create(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// CreateJSON is used to store JSON directory into a dataset collection.
// NOTE: the JSON is NOT validated.
//
// ```
//
//	import (
//	  "fmt"
//	  "os"
//	)
//
//	func main() {
//	    c, err := dataset.Open("friends.ds")
//	    if err != nil {
//	         fmt.Fprintf(os.Stderr, "%s", err)
//	         os.Exit(1)
//	    }
//	    defer c.Close()
//
//	    src := []byte(`{ "ID": "mojo", "Name": "Mojo Sam", "EMail": "mojo.sam@cosmic-cafe.example.org" }`)
//	    if err := c.CreateJSON("modo", src); err != nil {
//	         fmt.Fprintf(os.Stderr, "%s", err)
//	         os.Exit(1)
//	    }
//	    fmt.Printf("OK\n")
//	    os.Exit(0)
//	}
//
// ```
func (c *Collection) CreateJSON(key string, src []byte) error {
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Create(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// Read retrieves a map[string]inteferface{} from the collection,
// unmarshals it and updates the object pointed to by the map.
//
// ```
//
//	obj := map[string]interface{}{}
//
//	key := "123"
//	if err := c.Read(key, &obj); err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) Read(key string, obj map[string]interface{}) error {
	var (
		src []byte
		err error
	)
	switch c.StoreType {
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	if err != nil {
		return fmt.Errorf("failed to read %s, %s", key, err)
	}
	return JSONUnmarshal(src, &obj)
}

// ReadObject retrieves structed data via Go's general inteferface{} type.
// The JSON document is retreived from the collection,
// unmarshaled and variable holding the struct is updated.
//
// ```
//
//	type Record struct {
//	    ID string `json:"id"`
//	    Name string `json:"name,omitempty"`
//	    EMail string `json:"email,omitempty"`
//	}
//
//	// ...
//
//	var obj *Record
//
//	key := "123"
//	if err := c.Read(key, &obj); err != nil {
//	   // ... handle error
//	}
//
// ```
func (c *Collection) ReadObject(key string, obj interface{}) error {
	var (
		src []byte
		err error
	)
	switch c.StoreType {
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	if err != nil {
		return fmt.Errorf("failed to read %s, %s", key, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(src))
	decoder.UseNumber()
	if err := decoder.Decode(&obj); err != nil {
		return err
	}
	return nil
}

// ReadJSON retrieves JSON stored in a dataset collection for
// a given key. NOTE: It does not validate the JSON
//
// ```
//
//		key := "123"
//		src, err := c.ReadJSON(key)
//	 if err != nil {
//		   // ... handle error
//		}
//
// ```
func (c *Collection) ReadJSON(key string) ([]byte, error) {
	var (
		src []byte
		err error
	)
	switch c.StoreType {
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return nil, fmt.Errorf("%q not supported", c.StoreType)
	}
	if err != nil {
		return src, fmt.Errorf("failed to read %s, %s", key, err)
	}
	return src, nil
}


// Versions retrieves a list of versions available for a JSON document if
// versioning is enabled for the collection.
//
// ```
//
//	key, version := "123", "0.0.1"
//	if versions, err := Versions(key); err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) Versions(key string) ([]string, error) {
	var (
		versions []string
		err      error
	)
	switch c.StoreType {
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

// Update replaces a JSON document in the collection with a new one.
// If the collection is versioned then it creates a new versioned copy
// and updates the "current" version to use it.
//
// ```
//
//	key := "123"
//	obj["three"] = 3
//	if err := c.Update(key, obj); err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) Update(key string, obj map[string]interface{}) error {
	src, err := JSONMarshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
	}
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Update(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// UpdateObject replaces a JSON document in the collection with a new one.
// If the collection is versioned then it creates a new versioned copy
// and updates the "current" version to use it.
//
// ```
//
//	type Record struct {
//	    // ... structure def goes here.
//	    Three int `json:"three"`
//	}
//
//	var obj = *Record
//
//	key := "123"
//	obj := &Record {
//	  Three: 3,
//	}
//	if err := c.Update(key, obj); err != nil {
//	   // ... handle error
//	}
//
// ```
func (c *Collection) UpdateObject(key string, obj interface{}) error {
	src, err := JSONMarshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for %s, %s", key, err)
	}
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Update(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// UpdateJSON replaces a JSON document in the collection with a new one.
// NOTE: It does not validate the JSON
//
// ```
//
//	src := []byte(`{"Three": 3}`)
//	key := "123"
//	if err := c.UpdateJSON(key, src); err != nil {
//	   // ... handle error
//	}
//
// ```
func (c *Collection) UpdateJSON(key string, src []byte) error {
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Update(key, src)
		}
	default:
		return fmt.Errorf("%q not supported", c.StoreType)
	}
	return fmt.Errorf("%s not open", c.Name)
}

// Delete removes an object from the collection. If the collection is
// versioned then all versions are deleted. Any attachments to the
// JSON document are also deleted including any versioned attachments.
//
// ```
//
//	key := "123"
//	if err := c.Delete(key); err != nil {
//	   // ... handle error
//	}
//
// ```
func (c *Collection) Delete(key string) error {
	switch c.StoreType {
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
//
//	keys, err := c.Keys()
//	for _, key := range keys {
//	   ...
//	}
//
// ```
func (c *Collection) Keys() ([]string, error) {
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Keys()
		}
	default:
		return nil, fmt.Errorf("%q not supported", c.StoreType)
	}
	return nil, fmt.Errorf("%s not open", c.Name)
}

// KeysJSON returns a JSON encoded list of Keys
//
// ```
//
//	src, err := c.KeysJSON()
//  if err != nil {
//      // ... handle error ...
//  }
//  fmt.Printf("%s\n", src)
//
// ```
func (c *Collection) KeysJSON() ([]byte, error) {
	keys, err := c.Keys()
	if err != nil {
		return nil, err
	}
	return JSONMarshal(keys)
}


// HasKey takes a collection and checks if a key exists.
// NOTE: collection must be open otherwise false will always be returned.
//
// ```
//
//	key := "123"
//	if c.HasKey(key) {
//	   ...
//	}
//
// ```
func (c *Collection) HasKey(key string) bool {
	switch c.StoreType {
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
//
// ```
//
//	var x int64
//	x = c.Length()
//
// ```
func (c *Collection) Length() int64 {
	switch c.StoreType {
	case SQLSTORE:
		if c.SQLStore != nil {
			return c.SQLStore.Length()
		}
	}
	return int64(-1)
}


// Query implement the SQL query against a SQLStore or SQLties3 index of pairtree.
func (c *Collection) Query(sqlStmt string, debug bool) ([]interface{}, error) {
	src, err := c.QueryJSON(sqlStmt, debug)
	if err != nil {
		return nil, err
	}
	l := []interface{}{}
	if err := JSONUnmarshal(src, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// Query implement the SQL query against a SQLStore and return JSON results.
func (c *Collection) QueryJSON(sqlStmt string, debug bool) ([]byte, error) {
	if strings.Compare(c.StoreType, SQLSTORE) == 0 {
		if c.SQLStore == nil {
			return nil, fmt.Errorf("sqlstore failed to open")
		}
	} else {
		return nil, fmt.Errorf("not implemented for pairtree storage")
	}
	var (
		rows *sql.Rows
		err error
	)
	if debug {
		fmt.Fprintf(os.Stderr, "SQL: %s\n", sqlStmt)
	}
	rows, err = c.SQLStore.db.Query(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("sql: %s, %s", sqlStmt, err)
	}
	src := []byte(`[`)
	i := 0
	for rows.Next() {
		// Get our row values
		cell := []byte{}
		if err := rows.Scan(&cell); err != nil {
			return nil, err
		}
		if src == nil || len(src) == 0 {
			fmt.Fprintf(os.Stderr, "DEBUG expect src to have content\n")
		}
		if (i > 0) {
			src = append(src, ',')
		}
		src = append(src, cell...)
		i++
	}
	src = append(src, ']')
	err = rows.Err()
	if err != nil {
		return src, err
	}
	return src, nil
}
