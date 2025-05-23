/**
 * compatibility.go provides some wrapping methods for backward complatible
 * with v1 of dataset. These are likely to go away at some point.
 */
package dataset

import (
)

// CreateObjectsJSON takes a list of keys and creates a default object
// for each key as quickly as possible. This is useful in vary narrow
// situation like quickly creating test data. Use with caution.
//
// NOTE: if object already exist creation is skipped without
// reporting an error.
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
