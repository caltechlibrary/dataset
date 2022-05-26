package ptstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/pairtree"
)

type Storage struct {
	// Working path to the directory where the collections.json is found.
	WorkPath string

	// keyMapName holds the path to the keys.json key map file.
	keyMapName string

	// The KeyMap holds a map from key to JSON document in pairtree.
	// The pairtree is relative to the WorkPath. It is read from
	// keys.json in the WorkPath directory
	keyMap map[string]string

	// keys holds a sorted list of keys from the map
	keys []string
}

// Open opens the storage system and returns an storage struct and error
// It is passed a directory name that holds collection.json.
// The second parameter is for a DSN URI which is ignored in a pairtree
// implementation.
//
// ```
//  name := "testout/T1.ds" // a collection called "T1.ds"
//  store, err := c.Store.Open(name, "")
//  if err != nil {
//     ...
//  }
//  defer store.Close()
// ```
//
func Open(name string, dsnURI string) (*Storage, error) {
	store := new(Storage)
	store.WorkPath = name
	// Find the key map file and read it
	store.keyMapName = path.Join(name, "keymap.json")
	store.keyMap = map[string]string{}
	src, err := ioutil.ReadFile(store.keyMapName)
	if err == nil {
		// We have data so we need to decode it.
		if err := json.Unmarshal(src, &store.keyMap); err != nil {
			return nil, fmt.Errorf("failed to decode key map for %q, %s", name, err)
		}
		for key := range store.keyMap {
			store.keys = append(store.keys, key)
		}
		sort.Strings(store.keys)
	}
	// FIXME: Frames need implementation
	return store, nil
}

// Close closes the storage system freeing resources as needed.
//
// ```
//   if err := store.Close(); err != nil {
//      ...
//   }
// ```
//
func (store *Storage) Close() error {
	src, err := json.Marshal(store.keyMap)
	if err != nil {
		return fmt.Errorf("could not encode key map for %q, %s", store.WorkPath, err)
	}
	if err := ioutil.WriteFile(store.keyMapName, src, 0664); err != nil {
		return fmt.Errorf("failed to write kep map for %q, %s", store.WorkPath, err)
	}
	//FIXME: Implement frames save
	return nil
}

// Create stores a new JSON object in the collection
// It takes a string as a key and a byte slice of encoded JSON
//
//   err := store.Create("123", []byte(`{"one": 1}`))
//   if err != nil {
//      ...
//   }
//
func (store *Storage) Create(key string, src []byte) error {
	// NOTE: Keys are always normalized to lower case due to
	// naming issues in case insensitive file systems.
	key = strings.ToLower(key)
	if _, foundIt := store.keyMap[key]; foundIt {
		return fmt.Errorf("%s exists in %s", key, store.WorkPath)
	}
	sep := pairtree.Separator
	if sep != '/' {
		pairtree.Set('/')
	}

	// Figure out the map key (i.e. always / delimited)
	ptKey, ptPath := pairtree.Encode(key), ""
	if os.IsPathSeparator('/') {
		ptPath = ptKey
	} else {
		// OS dependent pairtree path
		pairtree.Set(os.PathSeparator)
		ptPath = pairtree.Encode(key)
	}
	// Return the seperator to the original state
	pairtree.Set(sep)

	// Generate the path to store document
	dName := path.Join(store.WorkPath, "pairtree", ptPath)
	if _, err := os.Stat(dName); err == nil {
		return fmt.Errorf("%s exists in %s", key, store.WorkPath)
	}
	if err := os.MkdirAll(dName, 0775); err != nil {
		return fmt.Errorf("Unable to create %q, %s", dName, err)
	}

	// Save the document to the ptPath location
	fName := path.Join(dName, fmt.Sprintf("%s.json", key))
	if err := ioutil.WriteFile(fName, src, 0664); err != nil {
		return fmt.Errorf("failed to write %q, %s", fName, err)
	}
	// Update keyMap
	store.keyMap[key] = ptKey

	// Insert into store's keys list and re-sort
	store.keys = append(store.keys, key)
	sort.Strings(store.keys)

	// Save the metadata for the updated key map
	src, err := json.Marshal(store.keyMap)
	if err != nil {
		return fmt.Errorf("unable to encode key map for %q in %q, %s", key, store.WorkPath, err)
	}
	if err := ioutil.WriteFile(store.keyMapName, src, 0664); err != nil {
		return fmt.Errorf("failed to write %q, %s", store.keyMapName, err)
	}
	return nil
}

// Read retrieves takes a string as a key and returns the encoded
// JSON document from the collection
//
//   src, err := store.Read("123")
//   if err != nil {
//      ...
//   }
//   obj := map[string]interface{}{}
//   if err := json.Unmarshal(src, &obj); err != nil {
//      ...
//   }
func (store *Storage) Read(key string) ([]byte, error) {
	// NOTE: Keys are always normalized to lower case due to
	// naming issues in case insensitive file systems.
	key = strings.ToLower(key)
	ptPath, ok := store.keyMap[key]
	if !ok {
		return nil, fmt.Errorf("%q not found in %q", key, store.WorkPath)
	}
	// Normalize the disk path if necessary
	if !os.IsPathSeparator('/') {
		ptPath = path.Join(strings.Split(ptPath, "/")...)
	}
	fName := path.Join(store.WorkPath, "pairtree", ptPath, fmt.Sprintf("%s.json", key))
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q in %q, %s", key, store.WorkPath, err)
	}
	return src, nil
}

// Update takes a key and encoded JSON object and updates a
// JSON document in the collection.
//
//   key := "123"
//   src := []byte(`{"one": 1, "two": 2}`)
//   if err := store.Update(key, src); err != nil {
//      ...
//   }
//
func (store *Storage) Update(key string, src []byte) error {
	// NOTE: Keys are always normalized to lower case due to
	// naming issues in case insensitive file systems.
	key = strings.ToLower(key)
	ptPath, ok := store.keyMap[key]
	if !ok {
		return fmt.Errorf("%q does not exists in %q", key, store.WorkPath)
	}

	// Save the document to the ptPath location
	fName := path.Join(store.WorkPath, "pairtree", ptPath, fmt.Sprintf("%s.json", key))
	if err := ioutil.WriteFile(fName, src, 0664); err != nil {
		return fmt.Errorf("failed to write %q, %s", fName, err)
	}
	return nil
}

// Delete removes a JSON document and any attachments from the collection
//
//   key := "123"
//   if err := store.Delete(key); err != nil {
//      ...
//   }
//
func (store *Storage) Delete(key string) error {
	// NOTE: Keys are always normalized to lower case due to
	// naming issues in case insensitive file systems.
	key = strings.ToLower(key)
	ptPath, ok := store.keyMap[key]
	if !ok {
		return fmt.Errorf("%q does not exists in %q", key, store.WorkPath)
	}

	// Save the document to the ptPath location
	dName := path.Join(store.WorkPath, "pairtree", ptPath)
	if err := os.RemoveAll(dName); err != nil {
		return fmt.Errorf("failed to delete %q in %q, %s", key, store.WorkPath, err)
	}
	delete(store.keyMap, key)
	// Save the metadata for the updated key map
	src, err := json.Marshal(store.keyMap)
	if err != nil {
		return fmt.Errorf("unable to encode key map for %q in %q, %s", key, store.WorkPath, err)
	}
	if err := ioutil.WriteFile(store.keyMapName, src, 0664); err != nil {
		return fmt.Errorf("failed to write %q, %s", store.keyMapName, err)
	}

	// Remove key from store.keys, could be more efficient ...
	l := len(store.keys) - 1
	for i, val := range store.keys {
		if val == key {
			if i < l {
				store.keys = append(store.keys[:i], store.keys[i+1:]...)
			} else {
				store.keys = store.keys[:i]
			}
			break
		}
	}
	return nil
}

// List returns all keys in a collection as a slice of strings.
//
//   var keys []string
//   keys, _ = store.Keys()
//   /* iterate over the keys retrieved */
//   for _, key := range keys {
//      ...
//   }
//
// NOTE: the error will always be nil, this func signature needs to match
// the other storage engines.
func (store *Storage) Keys() ([]string, error) {
	return store.keys, nil
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
	key = strings.ToLower(key)
	ok := false
	if len(store.keyMap) > 0 {
		_, ok = store.keyMap[key]
	}
	return ok
}

// Length returns the number of records (len(store.keys)) in the collection
// Requires collection to be open.
func (store *Storage) Length() int64 {
	return int64(len(store.keys))
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
