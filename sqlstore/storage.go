package sqlstore

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	// Database specific drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "modernc.org/sqlite"
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
	var (
		ok  bool
		err error
	)

	// Check to see if the DSN coming from th environment
	if dsnURI == "" {
		dsnURI = os.Getenv("DATASET_DSN_URI")
	}
	store := new(Storage)
	store.WorkPath = name
	store.tableName = path.Base(store.WorkPath)
	store.driverName, store.dsn, ok = strings.Cut(dsnURI, "://")
	if !ok {
		return nil, fmt.Errorf(`DSN URI is malformed, expected DRIVER_NAME://DSN, got %q`, dsnURI)
	}
	// Validate the driver name as supported by sqlstore ...
	switch store.driverName {
	case "sqlite":
	case "mysql":
	case "pg":
	default:
		parts := strings.SplitN(dsnURI, "://", 2)
		return nil, fmt.Errorf("%q database not supported", parts[0])
	}
	store.db, err = sql.Open(store.driverName, store.dsn)
	if err != nil {
		return nil, err
	}
	// NOTE: These need to be tuned are suggested in the documentation at
	// https://pkg.go.dev/database/sql
	store.db.SetConnMaxLifetime(0)
	store.db.SetMaxIdleConns(50)
	store.db.SetMaxOpenConns(50)

	return store, err
}

// Init creates a table to hold the collection if it doesn't already
// exist.
func Init(name string, dsnURI string) (*Storage, error) {
	store, err := Open(name, dsnURI)
	if err != nil {
		return nil, err
	}

	// NOTE: need to make sure that store.tableName exists.
	// FIXME: This create statement is MySQL centric, needs to work
	// or be modified for all the supported databases
	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
  key VARCHAR(255) DEFAULT NOT NULL PRIMARY KEY,
  src JSON,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`, store.tableName)

	_, err = store.db.Exec(stmt)
	if err != nil {
		return nil, err
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
	case "pg":
		return store.db.Close()
	default:
		return fmt.Errorf("%q database not supported", store.driverName)
	}
	return nil
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
	stmt := fmt.Sprintf(`INSERT INTO %q (key, src) VALUES (?, ?)`, store.tableName)
	_, err := store.db.Exec(stmt, key, string(src))
	return err
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
	stmt := fmt.Sprintf(`SELECT src FROM %s WHERE key = ? LIMIT 1`, store.tableName)
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

// Update takes a key and encoded JSON object and updates a
// JSON document in the collection.
//
//   key := "123"
//   src := []byte(`{"one": 1, "two": 2}`)
//   if err := storage.Update(key, src); err != nil {
//      ...
//   }
//
func (store *Storage) Update(key string, src []byte) error {
	stmt := fmt.Sprintf(`REPLACE INTO %q (key, src) VALUES (?, ?)`, store.tableName)
	_, err := store.db.Exec(stmt, key, string(src))
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
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE key = ?`, store.tableName)
	_, err := store.db.Exec(stmt, key)
	return err
}

// List returns all keys in a collection as a slice of strings.
//
//   var keys []string
//   keys, _ = storage.Keys()
//   /* iterate over the keys retrieved */
//   for _, key := range keys {
//      ...
//   }
//
func (store *Storage) Keys() ([]string, error) {
	stmt := fmt.Sprintf(`SELECT key FROM %s ORDER BY key`, store.tableName)
	rows, err := store.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		value string
		keys  []string
	)
	if rows.Next() {
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
	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE key = ? LIMIT 1`, store.tableName)
	rows, err := store.db.Query(stmt)
	if err != nil {
		return false
	}
	defer rows.Close()

	var (
		value int
	)
	if rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return false
		}
	}
	if err := rows.Err(); err != nil {
		return false
	}
	return value > 0
}

// Length returns the number of records (count of rows in collection).
// Requires collection to be open.
func (store *Storage) Length() int64 {
	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, store.tableName)
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
