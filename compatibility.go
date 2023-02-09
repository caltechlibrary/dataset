/**
 * compatibility.go provides some wrapping methods for backward complatible
 * with v1 of dataset. These are likely to go away at some point.
 */
package dataset

import (
	"fmt"
)

// DocPath method provides access to a PTStore's document path. If the
// collection is not a PTStore then an empty path and error is returend.
func (c *Collection) DocPath(key string) (string, error) {
	if c.StoreType == PTSTORE && c.PTStore != nil {
		return c.PTStore.DocPath(key)
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

