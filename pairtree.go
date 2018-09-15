package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	if store.Type == storage.FS {
		// Check to see if collectionName path exists, if not we need to create it.
		_, err := os.Stat(collectionName)
		if err != nil {
			os.MkdirAll(collectionName, 0775)
		}
	}
	// See if we need an open or continue with create
	_, err = store.Stat(collectionName + "/collection.json")
	if err == nil {
		return Open(name)
	}
	c := new(Collection)
	c.Version = Version
	c.Name = collectionName
	c.Layout = PAIRTREE_LAYOUT
	c.KeyMap = map[string]string{}
	c.Store = store
	err = c.saveMetadata()
	if err != nil {
		return nil, err
	}
	/*
		// Save the metadata for collection
		src, err := json.Marshal(c)
		if err != nil {
			return nil, fmt.Errorf("Can't marshal metadata for %s, %s", name, err)
		}
		err = store.WriteFile(path.Join(c.Name, "collection.json"), src, 0664)
		if err != nil {
			return nil, fmt.Errorf("Can't store collection metadata, %s", err)
		}
	*/
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

	pair := pairtree.Encode(key)
	pairPath := path.Join("pairtree", pair)
	err := c.Store.MkdirAll(path.Join(c.Name, pairPath), 0770)
	if err != nil {
		return fmt.Errorf("mkdir %s %s", path.Join(c.Name, pairPath), err)
	}

	// We've almost made it, save the key's bucket name and write the blob to bucket
	c.KeyMap[keyName] = pairPath
	err = c.Store.WriteFile(path.Join(c.Name, pairPath, FName), src, 0664)
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
	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return nil, fmt.Errorf("%q does not exist in %s", keyName, c.Name)
	}
	// NOTE: c.Name is the path to the collection not the name of JSON document
	// we need to join c.Name + bucketName + name to get path do JSON document
	src, err := c.Store.ReadFile(path.Join(c.Name, pairPath, FName))
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
	keyName, fName := keyAndFName(name)

	// Make sure we have an "object" not an array object in JSON notation
	if bytes.HasPrefix(src, []byte(`{`)) == false {
		return fmt.Errorf("dataset can only stores JSON objects")
	}
	// Add a _Key value if needed in the JSON source
	if bytes.Contains(src, []byte(`"_Key"`)) == false {
		src = bytes.Replace(src, []byte(`{`), []byte(`{"_Key":"`+keyName+`",`), 1)
	}

	//NOTE: KeyMap should include pairtree path (e.g. pairtree/AA/BB/CC...)
	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q does not exist in %q", keyName, c.Name)
	}
	err := c.Store.MkdirAll(path.Join(c.Name, pairPath), 0770)
	if err != nil {
		return fmt.Errorf("Update (mkdir) %q, %s", path.Join(c.Name, pairPath), err)
	}
	return c.Store.WriteFile(path.Join(c.Name, pairPath, fName), src, 0664)
}

// pairtreeDelete removes a JSON doc from a collection
func (c *Collection) pairtreeDelete(name string) error {
	if c.Layout != PAIRTREE_LAYOUT {
		return fmt.Errorf("Collection does not use a pairtree layout")
	}
	name = normalizeKeyName(name)
	keyName, FName := keyAndFName(name)

	pairPath, ok := c.KeyMap[keyName]
	if ok != true {
		return fmt.Errorf("%q key not found in %q", keyName, c.Name)
	}

	//NOTE: Need to remove any stale tarball before removing our record!
	tarball := strings.TrimSuffix(FName, ".json") + ".tar"
	p := path.Join(c.Name, pairPath, tarball)
	if err := c.Store.RemoveAll(p); err != nil {
		return fmt.Errorf("Can't remove attachment for %q, %s", keyName, err)
	}
	p = path.Join(c.Name, pairPath, FName)
	if err := c.Store.Remove(p); err != nil {
		return fmt.Errorf("Error removing %q, %s", p, err)
	}

	delete(c.KeyMap, keyName)
	return c.saveMetadata()
}

// pairtreeAnalyzer will scan a pairtree based collection for errors.
func pairtreeAnalyzer(collectionName string) error {
	var (
		eCnt int
		wCnt int
		kCnt int
		data interface{}
		c    *Collection
		err  error
	)

	store, err := storage.GetStore(collectionName)
	if err != nil {
		return err
	}
	files, err := store.ReadDir(collectionName)
	if err != nil {
		return err
	}
	hasNamaste := false
	hasCollectionJSON := false
	for _, file := range files {
		fname := file.Name()
		switch {
		case strings.HasPrefix(fname, "0=dataset_"):
			hasNamaste = true
		case fname == "collection.json":
			hasCollectionJSON = true
		}
		if hasNamaste && hasCollectionJSON {
			break
		}
	}

	// NOTE: Check for Namaste 0=, warn if missing
	if hasNamaste == false {
		log.Printf("WARNING: Missing Namaste 0=dataset_%s\n", Version[1:])
		wCnt++
	}

	// NOTE: Check to see if we have a collections.json
	if hasCollectionJSON == false {
		log.Printf("WARNING: Missing collection.json\n")
		wCnt++
	} else {
		// Make sure we can JSON parse the file
		docPath := path.Join(collectionName, "collection.json")
		if src, err := store.ReadFile(docPath); err == nil {
			if err := json.Unmarshal(src, &data); err == nil {
				// release the memory
				data = nil
			} else {
				log.Printf("ERROR: parsing %s, %s", docPath, err)
				eCnt++
			}
		} else {
			log.Printf("ERROR: opening %s, %s", docPath, err)
			eCnt++
		}
	}

	// Now try to open the collection ...
	c, err = Open(collectionName)
	if err != nil {
		return err
	}
	if c.Store.Type != storage.FS {
		return fmt.Errorf("Analyzer only works on local file system")
	}

	// Set layout to PAIRTREE_LAYOUT
	c.Layout = PAIRTREE_LAYOUT
	// Make sure we have all the known pairs in the pairtree
	// Check to see if records can be found in their buckets
	for k, v := range c.KeyMap {
		dirPath := path.Join(collectionName, v)
		// NOTE: k needs to be urlencoded before checking for file
		fname := url.QueryEscape(k) + ".json"
		docPath := path.Join(collectionName, v, fname)
		if store.IsDir(dirPath) == false {
			log.Printf("ERROR: %s is missing (%q)", k, dirPath)
			eCnt++
		} else if store.IsFile(docPath) == false {
			log.Printf("ERROR: %s is missing (%q)", k, docPath)
			eCnt++
		}
		kCnt++
		if (kCnt % 5000) == 0 {
			log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))
		}
	}
	if len(c.KeyMap) > 0 {
		log.Printf("%d of %d keys checked", kCnt, len(c.KeyMap))
	}

	// Check sub-directories in pairtree find but not in KeyMap
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil && len(c.KeyMap) > 0 {
		log.Printf("ERROR: unable to walk pairtree, %s", err)
		eCnt++
	} else {
		for _, pair := range pairs {
			key := pairtree.Decode(pair)
			if _, exists := c.KeyMap[key]; exists == false {
				log.Printf("WARNING: %s found at %q not in collection", key, path.Join(collectionName, "pairtree", pair, key+".json"))
				wCnt++
			}
		}
	}
	// FIXME: need to check for attachments and make sure they are record OK

	if eCnt > 0 || wCnt > 0 {
		return fmt.Errorf("%d errors, %d warnings detected", eCnt, wCnt)
	}
	return nil
}

func pairtreeRepair(collectionName string) error {
	var (
		c   *Collection
		err error
	)

	store, err := storage.GetStore(collectionName)
	if err != nil {
		return fmt.Errorf("Repair only works supported storage types, %s", err)
	}
	if store.Type != storage.FS {
		return fmt.Errorf("Repair only works on local file system")
	}

	// See if we can open a collection, if not then create an empty struct
	c, err = Open(collectionName)
	if err != nil {
		log.Printf("Open %s error, %s, attempting to re-create collection.json", collectionName, err)
		err = store.WriteFile(path.Join(collectionName, "collection.json"), []byte("{}"), 0664)
		if err != nil {
			log.Printf("Can't re-initilize %s, %s", collectionName, err)
			return err
		}
		log.Printf("Attempting to re-open %s", collectionName)
		c, err = Open(collectionName)
		if err != nil {
			log.Printf("Failed to re-open %s, %s", collectionName, err)
			return err
		}
	}
	defer c.Close()

	if c.Version != Version {
		log.Printf("Migrating format from %s to %s", c.Version, Version)
	}
	c.Version = Version
	if c.Layout == UNKNOWN_LAYOUT {
		c.Layout = PAIRTREE_LAYOUT
	}
	log.Printf("Getting a list of pairs")
	pairs, err := walkPairtree(path.Join(collectionName, "pairtree"))
	if err != nil {
		log.Printf("ERROR: unable to walk pairtree, %s", err)
		return err
	}
	log.Printf("Adding missing pairs")
	if c.KeyMap == nil {
		c.KeyMap = map[string]string{}
	}
	for _, pair := range pairs {
		key := pairtree.Decode(pair)
		if _, exists := c.KeyMap[key]; exists == false {
			c.KeyMap[key] = path.Join("pairtree", pair)
		}
	}
	log.Printf("%d keys in pairtree", len(c.KeyMap))
	keyList := c.Keys()
	log.Printf("checking that each key resolves to a value on disc")
	for _, key := range keyList {
		p, err := c.DocPath(key)
		if err != nil {
			break
		}
		if _, err := store.Stat(p); os.IsNotExist(err) == true {
			log.Printf("Removing %s from %s, %s does not exist", key, collectionName, p)
			delete(c.KeyMap, key)
		}
	}
	log.Printf("Saving metadata for %s", collectionName)
	return c.saveMetadata()
}

// migrateToPairtree will migrate JSON objects and attachments from
// a bucket oriented collection to a pairtree.
func migrateToPairtree(collectionName string) error {
	// Open existing collection, get objects and attachments
	// and manually place in new layout updating nc.
	c, err := Open(collectionName)
	if err != nil {
		return err
	}
	oldKeyMap := map[string]string{}
	for k, v := range c.KeyMap {
		oldKeyMap[k] = v
	}
	c.Close()
	store, err := storage.GetStore(collectionName)
	if err != nil {
		return err
	}

	// Create a new collection struct, set to Buckets layout
	nc := new(Collection)
	nc.Layout = PAIRTREE_LAYOUT
	nc.Name = collectionName
	nc.Version = Version
	nc.Buckets = nil
	nc.Store, _ = storage.GetStore(collectionName)
	nc.KeyMap = map[string]string{}

	for key, oldPath := range oldKeyMap {
		_, FName := keyAndFName(key)
		src, err := store.ReadFile(path.Join(collectionName, oldPath, FName))
		if err != nil {
			return err
		}
		// Write object to the new location
		err = nc.CreateJSON(key, src)
		if err != nil {
			return err
		}

		// Check for and handle any attachments
		tarballFName := strings.TrimSuffix(FName, ".json") + ".tar"
		oldTarballPath := path.Join(collectionName, oldPath, tarballFName)
		if store.IsFile(oldTarballPath) {
			// Move the tarball from one layout to the other
			buf, err := store.ReadFile(oldTarballPath)
			if err != nil {
				return err
			}
			pair := pairtree.Encode(key)
			pairPath := path.Join("pairtree", pair)
			newTarballPath := path.Join(collectionName, pairPath, tarballFName)
			err = nc.Store.WriteFile(newTarballPath, buf, 0664)
			if err != nil {
				return err
			}
		}
	}
	// OK, if all buckets processed, we can remove all the paths.
	for _, oldPath := range oldKeyMap {
		err = store.RemoveAll(path.Join(collectionName, oldPath))
		if err != nil {
			return fmt.Errorf("Cleaning after migration, %s", err)
		}
	}
	return nil
}

//
// Helper functions
//

// walkPairtree takes a store, a start path and returns a list
// of pairs found that also contain a pair's ${ID}.json file
func walkPairtree(startPath string) ([]string, error) {
	// pairs holds a list of discovered pairs
	pairs := []string{}
	err := filepath.Walk(startPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			f := path.Base(p)
			e := path.Ext(f)
			if e == ".json" {
				//NOTE: should be URL encoded at this point.
				key := strings.TrimSuffix(f, e)
				pair := pairtree.Encode(key)
				if strings.Contains(p, path.Join("pairtree", pair, f)) {
					pairs = append(pairs, pair)
				}
			}
		}
		return nil
	})
	return pairs, err
}
