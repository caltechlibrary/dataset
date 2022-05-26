package sqlstore

import (
	"fmt"
	"os"
)

type Storage struct {
}

// Open opens the storage system and returns an storage struct and error
// It is passed either a filename. For a Pairtree the would be the
// path to collection.json and for a sql store file holding a DSN URI.
// The DSN URI is formed from a protocal prefixed to the DSN. E.g.
// for a SQLite connection to test.ds database the DSN URI might be
// "sqlite:file:test.ds?cache=shared".
//
// ```
//  store, err := c.Store.Open(c.Name, c.DsnURI)
//  if err != nil {
//     ...
//  }
// ```
//
func Open(name string, dsnURI string) (*Storage, error) {
	return nil, fmt.Errorf("DEBUG Open not working")
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
	return fmt.Errorf("Close() not implemented")
}

// Create stores a new JSON object in the collection
// It takes a string as a key and a byte slice of encoded JSON
//
//   err := storage.Create("123", []byte(`{"one": 1}`))
//   if err != nil {
//      ...
//   }
//
func (store *Storage) Create(string, []byte) error {
	return fmt.Errorf("Create() not implemented")
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
func (store *Storage) Read(string) ([]byte, error) {
	return nil, fmt.Errorf("Read() not implemented")
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
func (store *Storage) Update(string, []byte) error {
	return fmt.Errorf("Read() not implemented")
}

// Delete removes a JSON document from the collection
//
//   key := "123"
//   if err := storage.Delete(key); err != nil {
//      ...
//   }
//
func (store *Storage) Delete(string) error {
	return fmt.Errorf("Read() not implemented")
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
	return nil, fmt.Errorf("Read() not implemented")
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
	fmt.Fprintf(os.Stderr, "HasKey() not implemented")
	return false
}

// Length returns the number of records (count of rows in collection).
// Requires collection to be open.
func (store *Storage) Length() int64 {
	//FIXME: not emplimented.
	return int64(-1)
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
