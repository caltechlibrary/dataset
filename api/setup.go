package api

import (
	"fmt"
	"os"

	// Caltech Library packages
	ds "github.com/caltechlibrary/dataset/v2"
)

func SetupTestCollection(cName string, dsnURI string, records map[string]map[string]interface{}) error {
	// Create collection.json using v1 structures
	if len(cName) == 0 {
		return fmt.Errorf("missing a collection name")
	}
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := ds.Init(cName, dsnURI)
	if err != nil {
		return err
	}
	defer c.Close()
	// Now populate with some test records records.
	for key, obj := range records {
		if err := c.Create(key, obj); err != nil {
			return err
		}
	}
	return nil
}
