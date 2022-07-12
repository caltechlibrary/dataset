//
// sqlstore is a sub module of the dataset package.
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
package sqlstore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v2/semver"

	// Database specific drivers
	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// None means versioning is turned off for collection
	None = iota
	// Major means increment the major semver value on creation or update
	Major
	// Minor means increment the minor semver value on creation or update
	Minor
	// Patach means increment the patch semver value on creation or update
	Patch

	versionPrefix = "_v_"
)

type Storage struct {
	// WorkPath holds the path to where the collection definition is held.
	WorkPath string

	// TableName holds the table name associated with the collection.
	// Usually the same as the "basename" in the WorkPath
	tableName string

	// dsn holds the SQL connection information needed to access
	// the SQL stored collection, it is everything after the protocol
	// in the DSN URI of the collection.
	dsn string

	// driverName
	driverName string

	// db database handle
	db *sql.DB

	// versioning
	Versioning int
}

// Init creates a table to hold the collection if it doesn't already
// exist.
func Init(name string, dsnURI string) (*Storage, error) {
	var err error

	driverName, dsn, ok := strings.Cut(dsnURI, "://")
	if !ok {
		return nil, fmt.Errorf("could not parse DSN URI, got %q", dsnURI)
	}
	store := new(Storage)
	store.WorkPath = name
	store.dsn = dsn
	store.driverName = driverName
	store.tableName = strings.TrimSuffix(strings.ToLower(path.Base(name)), ".ds")
	// Validate we support this SQL driver and form create statement.
	var stmt string
	switch driverName {
	case "sqlite":
		stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  key VARCHAR(255) PRIMARY KEY,
  src JSON,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP
)`, store.tableName)
	case "mysql":
		stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS `+"`"+`%s`+"`"+` (
  `+"`"+`key`+"`"+` VARCHAR(255) PRIMARY KEY,
  src JSON,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`, store.tableName)
	default:
		return nil, fmt.Errorf("%q database not supported", store.driverName)
	}
	// Open the DB
	db, err := sql.Open(store.driverName, store.dsn)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("%s opened and returned nil", store.driverName)
	}
	store.db = db

	// Create the collection table
	_, err = store.db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("Failed to create table %q, %s", store.tableName, err)
	}
	return store, err
}

// saveNewVersion saves an object to the version table for collection
func (store *Storage) saveNewVersion(key string, src []byte) error {
	// Figure out the next version number in sequence
	var (
		sv *semver.Semver
	)
	l, err := store.Versions(key)
	if err != nil {
		return err
	}
	if len(l) > 0 {
		versions := []*semver.Semver{}
		for _, val := range l {
			sv, err := semver.Parse([]byte(val))
			if err == nil {
				versions = append(versions, sv)
			}
		}
		semver.Sort(versions)
		sv = versions[len(versions)-1]
	} else {
		sv, _ = semver.Parse([]byte("0.0.0"))
	}
	switch store.Versioning {
	case Major:
		sv.IncMajor()
	case Minor:
		sv.IncMinor()
	default:
		sv.IncPatch()
	}
	version := sv.String()
	versionTable := versionPrefix + store.tableName
	stmt := fmt.Sprintf(`INSERT INTO `+"`"+`%s`+"`"+` (`+"`"+`key`+"`"+`, version, src) VALUES (?, ?, ?)`, versionTable)
	_, err = store.db.Exec(stmt, key, version, string(src))
	if err != nil {
		return fmt.Errorf(`failed to save version %q for %q in %q, %s`, key, version, store.WorkPath, err)
	}
	return nil
}

// saveVersioning() is a help function to store current versioning settings.
func (store *Storage) saveVersioning() error {
	versioningName := path.Join(store.WorkPath, "versioning.json")
	src := []byte(fmt.Sprintf(`{"versioning": %d}`, store.Versioning))
	if _, err := os.Stat(store.WorkPath); os.IsNotExist(err) {
		os.MkdirAll(store.WorkPath, 775)
	}
	if err := ioutil.WriteFile(versioningName, src, 0664); err != nil {
		return err
	}
	return nil
}

// getVersioning() reads the versioning information for collection
// and returns the integer value it finds.
func (store *Storage) getVersioning() error {
	versioningName := path.Join(store.WorkPath, "versioning.json")
	if _, err := os.Stat(versioningName); os.IsNotExist(err) {
		store.Versioning = None
		return nil
	}
	src, err := ioutil.ReadFile(versioningName)
	if err != nil {
		return err
	}
	m := map[string]int{}
	if err := json.Unmarshal(src, &m); err != nil {
		return err
	}
	if val, ok := m["versioning"]; ok {
		switch val {
		case None:
			store.Versioning = None
		case Major:
			store.Versioning = Major
		case Minor:
			store.Versioning = Minor
		case Patch:
			store.Versioning = Patch
		default:
			store.Versioning = None
			return fmt.Errorf("Unknown/unsupported version type")
		}
	}
	return nil
}

// SetVersioning sets versioning to Major, Minor, Patch or None
// If versioning is set to Major, Minor or Patch a table in the
// open SQL storage engine will be created.
func (store *Storage) SetVersioning(setting int) error {
	switch setting {
	case None:
		store.Versioning = None
	case Major:
		store.Versioning = setting
	case Minor:
		store.Versioning = setting
	case Patch:
		store.Versioning = setting
	default:
		return fmt.Errorf("Unknown/unsupported version type")
	}
	if store.Versioning != None {
		var (
			stmt         string
			versionTable = versionPrefix + store.tableName
		)
		switch store.driverName {
		case "sqlite":
			stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  key VARCHAR(255) NOT NULL,
  version VARCHAR(255) NOT NULL,
  src JSON,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (key, version)
)`, versionTable)
		case "mysql":
			stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  key VARCHAR(255) NOT NULL,
  version VARCHAR(255) NOT NULL,
  src JSON,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (key, version)
)`, versionTable)
		default:
			return fmt.Errorf("%q database not supported", store.driverName)
		}
		// Create the collection table
		if _, err := store.db.Exec(stmt); err != nil {
			return fmt.Errorf("Failed to create version table %q, %s", versionTable, err)
		}
	}
	if err := store.saveVersioning(); err != nil {
		return err
	}
	return nil
}

// Open opens the storage system and returns an storage struct and error
// It is passed either a filename. For a Pairtree the would be the
// path to collection.json and for a sql store file holding a DSN URI.
// The DSN URI is formed from a protocal prefixed to the DSN. E.g.
// for a SQLite connection to test.ds database the DSN URI might be
// "sqlite://file:test.ds?cache=shared".
//
// ```
//  store, err := c.Store.Open(c.Name, c.DsnURI)
//  if err != nil {
//     ...
//  }
// ```
//
func Open(name string, dsnURI string) (*Storage, error) {
	var err error

	// Check to see if the DSN coming from th environment
	if dsnURI == "" {
		dsnURI = os.Getenv("DATASET_DSN_URI")
	}
	driverName, dsn, ok := strings.Cut(dsnURI, "://")
	if !ok {
		return nil, fmt.Errorf(`DSN URI is malformed, expected DRIVER_NAME://DSN, got %q`, dsnURI)
	}

	store := new(Storage)
	store.WorkPath = name
	store.tableName = strings.TrimSuffix(strings.ToLower(path.Base(name)), ".ds")
	store.driverName = driverName
	store.dsn = dsn
	// Validate the driver name as supported by sqlstore ...
	switch store.driverName {
	case "sqlite":
	case "mysql":
	default:
		return nil, fmt.Errorf("%q database not supported", store.driverName)
	}
	store.db, err = sql.Open(store.driverName, store.dsn)
	if err != nil {
		return nil, err
	}
	if store.db == nil {
		return nil, fmt.Errorf("%s opened and returned nil", store.driverName)
	}
	// NOTE: These need to be tuned are suggested in the documentation at
	// https://pkg.go.dev/database/sql
	store.db.SetConnMaxLifetime(0)
	store.db.SetMaxIdleConns(50)
	store.db.SetMaxOpenConns(50)

	if err := store.getVersioning(); err != nil {
		return store, err
	}
	return store, err
}

// Close closes the storage system freeing resources as needed.
//
// ```
//   if err := storage.Close(); err != nil {
//      ...
//   }
// ```
//
func (store *Storage) Close() error {
	switch store.driverName {
	case "sqlite":
		return store.db.Close()
	case "mysql":
		return store.db.Close()
	default:
		return fmt.Errorf("%q database not supported", store.driverName)
	}
}

// Create stores a new JSON object in the collection
// It takes a string as a key and a byte slice of encoded JSON
//
//   err := storage.Create("123", []byte(`{"one": 1}`))
//   if err != nil {
//      ...
//   }
//
func (store *Storage) Create(key string, src []byte) error {
	stmt := fmt.Sprintf(`INSERT INTO `+"`"+`%s`+"`"+` (`+"`"+`key`+"`"+`, src) VALUES (?, ?)`, store.tableName)
	_, err := store.db.Exec(stmt, key, string(src))
	if err != nil {
		return err
	}
	if store.Versioning != None {
		return store.saveNewVersion(key, src)
	}
	return nil
}

// Read retrieves takes a string as a key and returns the encoded
// JSON document from the collection
//
//   src, err := storage.Read("123")
//   if err != nil {
//      ...
//   }
//   obj := map[string]interface{}{}
//   if err := json.Unmarshal(src, &obj); err != nil {
//      ...
//   }
func (store *Storage) Read(key string) ([]byte, error) {
	stmt := fmt.Sprintf(`SELECT src FROM `+"`"+`%s`+"`"+` WHERE `+"`"+`key`+"`"+` = ? LIMIT 1`, store.tableName)
	rows, err := store.db.Query(stmt, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var value string

	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return []byte(value), nil
}

// Versions return a list of semver strings for a versioned object.
func (store *Storage) Versions(key string) ([]string, error) {
	stmt := fmt.Sprintf(`SELECT version FROM `+"`"+`%s`+"`"+` WHERE `+"`"+`key`+"`"+` = ?`, versionPrefix+store.tableName)
	rows, err := store.db.Query(stmt, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := []string{}
	value := ""
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	values = semver.SortStrings(values)
	return values, nil
}

// ReadVersion returns a specific version of a JSON object.
func (store *Storage) ReadVersion(key string, version string) ([]byte, error) {
	stmt := fmt.Sprintf(`SELECT src FROM `+"`"+`%s`+"`"+` WHERE `+"`"+`key`+"`"+` = ? AND version = ?`, versionPrefix+store.tableName)
	rows, err := store.db.Query(stmt, key, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	value := ""
	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
	}
	return []byte(value), nil
}

// Update takes a key and encoded JSON object and updates a
//
//   key := "123"
//   src := []byte(`{"one": 1, "two": 2}`)
//   if err := storage.Update(key, src); err != nil {
//      ...
//   }
//
func (store *Storage) Update(key string, src []byte) error {
	stmt := fmt.Sprintf(`REPLACE INTO `+"`"+`%s`+"`"+` (`+"`"+`key`+"`"+`, src) VALUES (?, ?)`, store.tableName)
	_, err := store.db.Exec(stmt, key, string(src))
	if err != nil {
		return err
	}
	if store.Versioning != None {
		return store.saveNewVersion(key, src)
	}
	return err
}

// Delete removes a JSON document from the collection
//
//   key := "123"
//   if err := storage.Delete(key); err != nil {
//      ...
//   }
//
func (store *Storage) Delete(key string) error {
	stmt := fmt.Sprintf(`DELETE FROM `+"`"+`%s`+"`"+` WHERE key = ?`, store.tableName)
	_, err := store.db.Exec(stmt, key)
	// FIXME: Remove attachments
	// FIXME: remove versions
	return err
}

// Keys returns all keys in a collection as a slice of strings.
//
//   var keys []string
//   keys, _ = storage.Keys()
//   /* iterate over the keys retrieved */
//   for _, key := range keys {
//      ...
//   }
//
func (store *Storage) Keys() ([]string, error) {
	stmt := fmt.Sprintf(`SELECT `+"`"+`key`+"`"+` FROM `+"`"+`%s`+"`"+` ORDER BY `+"`"+`key`+"`", store.tableName)
	rows, err := store.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		value string
		keys  []string
	)
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		if value != "" {
			keys = append(keys, value)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// UpdatedKeys returns all keys updated in a time range
//
//```
//   var (
//      keys []string
//      start = "2022-06-01 00:00:00"
//      end = "20022-06-30 23:23:59"
//   )
//   keys, _ = storage.UpdatedKeys(start, end)
//   /* iterate over the keys retrieved */
//   for _, key := range keys {
//      ...
//   }
//```
//
func (store *Storage) UpdatedKeys(start string, end string) ([]string, error) {
	if start == "" {
		return nil, fmt.Errorf("missing start time value")
	}
	if end == "" {
		return nil, fmt.Errorf("missing end time value")
	}
	stmt := fmt.Sprintf(`SELECT `+"`"+`key`+"`"+` FROM `+"`"+`%s`+"`"+` WHERE (updated >= ? AND updated <= ?) ORDER BY updated ASC`, store.tableName)

	rows, err := store.db.Query(stmt, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		value string
		keys  []string
	)
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return nil, err
		}
		if value != "" {
			keys = append(keys, value)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// HasKey will look up and make sure key is in collection.
// Storage must be open or zero false will always be returned.
//
// ```
//   key := "123"
//   if store.HasKey(key) {
//      ...
//   }
// ```
func (store *Storage) HasKey(key string) bool {
	stmt := fmt.Sprintf(`SELECT `+"`"+`key`+"`"+` FROM `+"`"+`%s`+"`"+` WHERE `+"`"+`key`+"`"+` = ? LIMIT 1`, store.tableName)
	rows, err := store.db.Query(stmt, key)
	if err != nil {
		return false
	}
	defer rows.Close()

	var (
		value string
	)
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return false
		}
	}
	if err := rows.Err(); err != nil {
		return false
	}
	return key == value
}

// Length returns the number of records (count of rows in collection).
// Requires collection to be open.
func (store *Storage) Length() int64 {
	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM `+"`"+`%s`+"`", store.tableName)
	rows, err := store.db.Query(stmt)
	if err != nil {
		return int64(-1)
	}
	defer rows.Close()

	var (
		value int64
	)
	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return int64(-1)
		}
	}
	if err := rows.Err(); err != nil {
		return int64(-1)
	}
	return value
}

// Frames
// Frame
// FrameDef
// FrameObjects
// Refresh
// Reframe
// DeleteFrame
// HasFrame

// Attachments
// Attach
// Retrieve
// Prune

// Sample
// Clone
// CloneSample

// Check
// Repair
