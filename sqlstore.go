// sqlstore is a part of dataset
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
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	// Database specific drivers
	//_ "github.com/glebarez/go-sqlite"
	_ "modernc.org/sqlite"
	
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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

type SQLStore struct {
	// WorkPath holds the path to where the collection definition is held.
	WorkPath string

	// primaryTable holds the table name associated with the collection.
	// Usually the same as the "basename" in the WorkPath
	primaryName string

	// historyName holds the history table associated with the collection.
	historyName string

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

func ParseDSN(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if !(u.Scheme == "sqlite" || u.Scheme == "postgres" || u.Scheme == "mysql") {
		return "", fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}

	var kvs []string
	escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
	accrue := func(k, v string) {
		if v != "" {
			kvs = append(kvs, k+"='"+escaper.Replace(v)+"'")
		}
	}

	if u.User != nil {
		v := u.User.Username()
		accrue("user", v)

		v, _ = u.User.Password()
		accrue("password", v)
	}

	if host, port, err := net.SplitHostPort(u.Host); err != nil {
		accrue("host", u.Host)
	} else {
		accrue("host", host)
		accrue("port", port)
	}

	if u.Path != "" {
		accrue("dbname", u.Path[1:])
	}

	q := u.Query()
	for k := range q {
		accrue(k, q.Get(k))
	}
	sort.Strings(kvs) // Makes testing easier (not a performance concern)
	return strings.Join(kvs, " "), nil
}

func dsnFixUp(driverName string, dsn string, workPath string) string {
	switch driverName {
	case "postgres":
		return fmt.Sprintf("%s://%s", driverName, dsn)
	case "sqlite":
		// NOTE: the db needs to be stored in the dataset directory
		// to keep the dataset easily movable.
		dbName := filepath.Base(dsn)
		return path.Join(workPath, dbName)
	}
	return dsn
}

// SQLStoreInit creates a table to hold the collection if it doesn't already
// exist.
func SQLStoreInit(name string, dsnURI string) (*SQLStore, error) {
	var err error

	driverName, dsn, ok := strings.Cut(dsnURI, "://")
	if !ok {
		return nil, fmt.Errorf("could not parse DSN URI, got %q", dsnURI)
	}
	store := new(SQLStore)
	store.WorkPath = name
	store.dsn = dsnFixUp(driverName, dsn, name)
	store.driverName = driverName
	store.primaryName = strings.TrimSuffix(strings.ToLower(filepath.Base(name)), ".ds")
	store.historyName = store.primaryName + "_history"
	// Validate we support this SQL driver and form create statement.
	var (
		stmt string
		stmtHistory string
	)
	switch driverName {
	case "sqlite":
		stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  _key TEXT PRIMARY KEY NOT NULL,
  src TEXT,
  created TEXT,
  updated TEXT,
  version INTEGER NOT NULL DEFAULT 0
);`, store.primaryName)
		stmtHistory = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
_key TEXT NOT NULL,
src TEXT,
created TEXT,
updated TEXT,
version INTEGER NOT NULL DEFAULT 0,
PRIMARY KEY (_key, version)
);
`, store.historyName)
	    //FIXME: Need to add Trigger to handle upsert
	case "postgres":
		stmt = fmt.Sprintf(`CREATE TABLE %s (_key TEXT PRIMARY KEY,
src JSON,
created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
version INTEGER NOT NULL DEFAULT 0
)`, store.primaryName)
		stmtHistory = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
_key TEXT NOT NULL,
src TEXT,
created TEXT,
updated TEXT,
version INTEGER NOT NULL DEFAULT 0,
PRIMARY KEY (_key, version)
);
`, store.historyName)
	    //FIXME: Need to add Trigger to handle upsert
	case "mysql":
		stmt = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
_key VARCHAR(512) PRIMARY KEY,
src JSON,
created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
version INTEGER NOT NULL DEFAULT 0
)`, store.primaryName)
		stmtHistory = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
_key VARCHAR(512)  NOT NULL,
src JSON,
created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
version INTEGER NOT NULL DEFAULT 0,
PRIMARY KEY (_key, version)
)`, store.historyName)
	    //FIXME: Need to add Trigger to handle upsert
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
		return nil, fmt.Errorf("Failed to create table %q, %s", store.primaryName, err)
	}

	// Create the collection history table
	_, err = store.db.Exec(stmtHistory)
	if err != nil {
		return nil, fmt.Errorf("Failed to create history table %q, %s", store.primaryName, err)
	}
	return store, err
}

// SQLStoreOpen opens the storage system and returns an storage struct and error
// It is passed either a filename. For a Pairtree the would be the
// path to collection.json and for a sql store file holding a DSN URI.
// The DSN URI is formed from a protocal prefixed to the DSN. E.g.
// for a SQLite connection to test.ds database the DSN URI might be
// "sqlite://collections.db".
//
// ```
//
//	store, err := c.Store.Open(c.Name, c.DsnURI)
//	if err != nil {
//	   ...
//	}
//
// ```
func SQLStoreOpen(name string, dsnURI string) (*SQLStore, error) {
	var err error

	// Check to see if the DSN coming from th environment
	if dsnURI == "" {
		dsnURI = os.Getenv("DATASET_DSN_URI")
	}
	driverName, dsn, ok := strings.Cut(dsnURI, "://")
	if !ok {
		return nil, fmt.Errorf(`DSN URI is malformed, expected DRIVER_NAME://DSN, got %q`, dsnURI)
	}

	store := new(SQLStore)
	store.WorkPath = name
	store.primaryName = strings.TrimSuffix(strings.ToLower(filepath.Base(name)), ".ds")
	store.historyName = store.primaryName + "_history"
	store.driverName = driverName
	store.dsn = dsnFixUp(driverName, dsn, name)
	// Validate the driver name as supported by sqlstore ...
	switch store.driverName {
	case "postgres":
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
	return store, err
}

// Close closes the storage system freeing resources as needed.
//
// ```
//
//	if err := storage.Close(); err != nil {
//	   ...
//	}
//
// ```
func (store *SQLStore) Close() error {
	switch store.driverName {
	case "sqlite":
		return store.db.Close()
	case "postgres":
		return store.db.Close()
	case "mysql":
		return store.db.Close()
	default:
		return fmt.Errorf("%q database not supported", store.driverName)
	}
}

// Write stores a new JSON or replacement object in the collection
// It takes a string as a key and a byte slice of encoded JSON
//
//	err := storage.Create("123", []byte(`{"one": 1}`))
//	if err != nil {
//	   ...
//	}
func (store *SQLStore) Write(key string, src []byte) error {
	var (
		stmt string
	)
	// FIXME: this should happen as a transaction
	switch store.driverName {
	case "postgres":
		stmt = fmt.Sprintf(`INSERT INTO %s (_key, src, created, updated, version) VALUES ($1, $2, NOW(), NOW(), 0)`, store.primaryName)
	case "sqlite":
		stmt = fmt.Sprintf(`INSERT INTO %s (_key, src, created, updated, version) VALUES (?, ?, datetime(), datetime(), 0)`, store.primaryName)
	default:
		stmt = fmt.Sprintf(`INSERT INTO %s (_key, src, created, updated, version) VALUES ($1, $2, NOW(), NOW(), 0)`, store.primaryName)
	}
	// Insert the row in the primary table, then use that row to populate the history table
	_, err := store.db.Exec(stmt, key, string(src))
	if err != nil {
		return err
	}
	return nil
}

// Read retrieves takes a string as a key and returns the encoded
// JSON document from the collection
//
//	src, err := storage.Read("123")
//	if err != nil {
//	   ...
//	}
//	obj := map[string]interface{}{}
//	if err := json.Unmarshal(src, &obj); err != nil {
//	   ...
//	}
func (store *SQLStore) Read(key string) ([]byte, error) {
	var (
		stmt string
	)
	switch store.driverName {
	case "postgres":
		stmt = fmt.Sprintf(`SELECT src FROM %s WHERE _key = $1`, store.primaryName)
	default:
		stmt = fmt.Sprintf(`SELECT src FROM %s WHERE _key = ?`, store.primaryName)
	}
	rows, err := store.db.Query(stmt, key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		value string
	)

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

// Delete deaccession the JSON document form the collection by recording a null object
// associated with the key. This preserves the object history and allows for an "undelete"
// where the history of the object is preserved.
//
//	key := "123"
//	if err := storage.Delete(key); err != nil {
//	   ...
//	}
func (store *SQLStore) Delete(key string) error {
	return store.Write(key, []byte(""))
}

// Keys returns all keys in a collection as a slice of strings.
//
//	var keys []string
//	keys, _ = storage.Keys()
//	/* iterate over the keys retrieved */
//	for _, key := range keys {
//	   ...
//	}
func (store *SQLStore) Keys() ([]string, error) {
	var stmt string

	//FIXME: Need to decide key's behavior for keys that were deleted. Do I list them or exclude them?
	switch store.driverName {
	case "postgres":
		stmt = fmt.Sprintf(`SELECT _key FROM %s ORDER BY _key`, store.primaryName)
	default:
		stmt = fmt.Sprintf(`SELECT _key FROM %s ORDER BY _key`, store.primaryName)
	}
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

// HasKey will look up and make sure key is in collection.
// SQLStore must be open or zero false will always be returned.
//
// ```
//
//	key := "123"
//	if store.HasKey(key) {
//	   ...
//	}
//
// ```
func (store *SQLStore) HasKey(key string) bool {
	var stmt string

	switch store.driverName {
	case "postgres":
		stmt = fmt.Sprintf(`SELECT _key FROM %s WHERE _key = $1 LIMIT 1`, store.primaryName)
	default:
		stmt = fmt.Sprintf(`SELECT _key FROM %s WHERE _key = ? LIMIT 1`, store.primaryName)
	}
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
func (store *SQLStore) Length() int64 {
	var stmt string

	switch store.driverName {
	default:
		stmt = fmt.Sprintf(`SELECT COUNT(*) FROM %s`, store.primaryName)
	}
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
