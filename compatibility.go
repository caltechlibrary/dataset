/**
 * compatibility.go provides some wrapping methods for backward complatible
 * with v1 of dataset. These are likely to go away at some point.
 */
package dataset

import (
	"fmt"
	"path"
	"strings"
)

// DocPath method provides access to a PTStore's document path. If the
// collection is not a PTStore then an empty path and error is returned
// with an error message.
// NOTE: the path returned is a full path including the JSON document
// stored.
//
// ```
//    c, err := dataset.Open(cName, "")
//    // ... handle error ...
//    key := "2488"
//    s, err := c.DocPath(key)
//    // ... handle error ...
//    fmt.Printf("full path to JSON document %q is %q\n", key, s)
// ```
func (c *Collection) DocPath(key string) (string, error) {
	if c.StoreType == PTSTORE && c.PTStore != nil {
		s, err := c.PTStore.DocPath(strings.ToLower(key))
		if err != nil {
			return "", err
		}
		return path.Join(s, key + ".json"), nil
	}
	return "", fmt.Errorf("%q does not support document paths", c.StoreType)
}


// CreateObjectsJSON takes a list of keys and creates a default object
// for each key as quickly as possible. This is useful in vary narrow
// situation like quickly creating test data. Use with caution.
//
// NOTE: if object already exist creation is skipped without
// reporting an error.
//
func (c *Collection) CreateObjectsJSON(keyList []string, src []byte) error {
	for _, key := range keyList {
		if c.HasKey(key) == false {
			if err := c.CreateJSON(key, src); err != nil {
				return err
			}
		}
	}
	return nil
}

