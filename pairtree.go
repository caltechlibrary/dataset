package dataset

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/pairtree"
	"github.com/caltechlibrary/storage"
)

//
// Pairtree file layout implementation for dataset collections.
//

// pairtreeCreateCollection - create a new collection structure on disc
// name should be filesystem friendly
func pairtreeCreateCollection(name string) (*Collection, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing a collection name")
	}
	collectionName := collectionNameFromPath(name)
	store, err := storage.GetStore(name)
	if err != nil {
		return nil, err
	}
	// See if we need an open or continue with create
	if store.Type == storage.S3 || store.Type == storage.GS {
		if _, err := store.Stat(collectionName + "/collection.json"); err == nil {
			return Open(name)
		}
	} else {
		if _, err := store.Stat(collectionName); err == nil {
			return Open(name)
		}
	}
	c := new(Collection)
	c.Version = Version
	c.Name = collectionName
	c.Layout = PAIRTREE_LAYOUT
	c.KeyMap = map[string]string{}
	c.Store = store
	// Save the metadata for collection
	err = c.saveMetadata()
	return c, err
}

// pairtreeCreateJSON adds a JSON doc to a collection, if a problem occurs it returns an error
func (c *Collection) pairtreeCreateJSON(key string, src []byte) error {
	if c.Layout != PAIRTREE_LAYOUT {
		return fmt.Errorf("Collection does not use a pairtree layout")
	}
	key = strings.TrimSpace(key)
	if key == "" || key == ".json" {
		return fmt.Errorf("must not be empty")
	}

	// Enforce the _Key attribute is unique and does not exist in collection already
	key = normalizeKeyName(key)
	keyName, FName := keyAndFName(key)
	if _, keyExists := c.KeyMap[keyName]; keyExists == true {
		return fmt.Errorf("%s already exists in collection %s", key, c.Name)
	}

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	pair := path.Join("pairtree", pairtree.Encode(key))
	err := c.Store.MkdirAll(path.Join(c.Name, pair), 0770)
	if err != nil {
		return fmt.Errorf("mkdir %s %s", pair, err)
	}

	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.KeyMap[keyName] = pair
	err = c.Store.WriteFile(path.Join(c.Name, pair, FName), src, 0664)
	if err != nil {
		return err
	}
	return c.saveMetadata()
}

// pairtreeReadJSON finds a the record in the collection and returns the JSON source
func (c *Collection) pairtreeReadJSON(name string) ([]byte, error) {
	if c.Layout != PAIRTREE_LAYOUT {
		return nil, fmt.Errorf("Collection does not use a pairtree layout")
	}
	name = normalizeKeyName(name)
	// Handle potentially URL encoded names
	keyName, FName := keyAndFName(name)
	p, ok := c.KeyMap[keyName]
	if ok != true {
		return nil, fmt.Errorf("%q does not exist in %s", keyName, c.Name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	// we need to join c.Name + bucketName + name to get path do JSON document
	src, err := c.Store.ReadFile(path.Join(c.Name, p, FName))
	if err != nil {
		return nil, err
	}
	return src, nil
}

// pairtreeUpdateJSON a JSON doc in a collection, returns an error if there is a problem
func (c *Collection) pairtreeUpdateJSON(name string, src []byte) error {
	if c.Layout != PAIRTREE_LAYOUT {
		return fmt.Errorf("Collection does not use a pairtree layout")
	}
	// Make sure Key exists before proceeding with update
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	//NOTE: KeyMap should include pairtree path (e.g. pairtree/AA/BB/CC...)
	pair, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist", keyName)
	}
	p := path.Join(c.Name, pair)
	err := c.Store.MkdirAll(p, 0770)
	if err != nil {
		return fmt.Errorf("Update (mkdir) %s %s", p, err)
	}
	return c.Store.WriteFile(path.Join(c.Name, p, FName), src, 0664)
}

// pairtreeDelete removes a JSON doc from a collection
func (c *Collection) pairtreeDelete(name string) error {
	if c.Layout != PAIRTREE_LAYOUT {
		return fmt.Errorf("Collection does not use a pairtree layout")
	}
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)

	pair, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found", keyName)
	}

	//NOTE: Need to remove any stale tarball before removing our record!
	tarball := keyName + ".tar"
	p := path.Join(c.Name, pair, tarball)
	if err := c.Store.RemoveAll(p); err != nil {
		return fmt.Errorf("Can't remove attachment for %q, %s", keyName, err)
	}
	p = path.Join(c.Name, pair, FName)
	if err := c.Store.Remove(p); err != nil {
		return fmt.Errorf("Error removing %q, %s", p, err)
	}

	delete(c.KeyMap, keyName)
	return c.saveMetadata()
}
