package dataset

import (
	"fmt"
	"strings"
)

// findLabel looks through an array of string for a specific label
func findLabel(labels []string, label string) (int, bool) {
	for pos, val := range labels {
		if val == label {
			return pos, true
		}
	}
	return -1, false
}

// strInArray checks to see if val is in an array of strings
func strInArray(a []string, val string) bool {
	for _, item := range a {
		if item == val {
			return true
		}
	}
	return false
}

// mergeKeys takes a keyList and an unordered list of keys
// appending the missing keys to the end of keyList
func mergeKeys(sorted []string, unsorted []string) []string {
	newKeys := []string{}
	for _, key := range unsorted {
		if strInArray(sorted, key) == false {
			newKeys = append(newKeys, key)
		}
	}
	if len(newKeys) > 0 {
		sorted = append(sorted, newKeys...)
	}
	return sorted
}

// dotPathToColumnMap creates a mapping from dotpath in collection
// to column number in table by matching header row values with
// a frame's labels. Returns an error if ._Key is not identified.
func dotPathToColumnMap(f *DataFrame, table [][]string) (map[string]int, error) {
	m := make(map[string]int)
	if len(f.Labels) != len(f.DotPaths) {
		return m, fmt.Errorf("corrupted frame, labels don't map to dot paths")
	}
	if len(table) < 2 {
		return m, fmt.Errorf("table is empty")
	}
	for colNo, label := range table[0] {
		if pos, hasLabel := findLabel(f.Labels, label); hasLabel == true {
			// Find the dotpath matching the label
			dotPath := f.DotPaths[pos]
			// Map the dotpath to a column number
			m[dotPath] = colNo
		}
	}
	// Sanity check the mapping for ._Key
	if _, hasID := m["._Key"]; hasID == false {
		return m, fmt.Errorf("table header row is missing %q column", f.Labels[0])
	}
	return m, nil
}

// rowToObj assembles a new JSON object from map into row and row values
// BUG: This is a naive map assumes all root level object properties
func rowToObj(key string, dotPathToCols map[string]int, row []string) map[string]interface{} {
	obj := map[string]interface{}{}
	for p, i := range dotPathToCols {
		if i < len(row) {
			attrName := strings.TrimPrefix(p, ".")
			obj[attrName] = row[i]
		}
	}
	obj["_Key"] = key
	return obj
}

// MergeFrameIntoTable - uses a DataFrame associated in the collection
// to map attributes into table appending new content and optionally
// overwriting existing content for rows with matching ids. Returns
// a new table (i.e. [][]string) or error.
func (c *Collection) MergeFrameIntoTable(frameName string, table [][]string, overwrite bool, verbose bool) ([][]string, error) {
	// MergeFrameIntoTable (overwrite) [][]interface{}{} with frame content
	//       adding extra columns/rows if needed
	return nil, fmt.Errorf("c.MergeFrameIntoTable() not implemented")
}

// MergeFromTable - uses a DataFrame associated in the collection
// to map columns from a table into JSON object attributes saving the
// JSON object in the collection.  If overwrite is true then JSON objects
// for matching keys will be updated, if false only new objects will be
// added to collection. Returns an error value
func (c *Collection) MergeFromTable(frameName string, table [][]string, overwrite bool, verbose bool) error {
	// Build Map dotpath to column position
	//
	// For each data row of table (i.e. row 1 through last row)
	//    get ID value
	//    if has ID && overwrite == true then join with overwrite
	//    else if has ID then join (append)
	//    else add object to collection
	// Regenerate the frame
	f, err := c.getFrame(frameName)
	if err != nil {
		return err
	}
	m, err := dotPathToColumnMap(f, table)
	if err != nil {
		return err
	}
	keyCol, _ := m["._Key"]
	key := ""
	keys := []string{}
	for _, row := range table[1:] {
		// get Key
		if keyCol < len(row) {
			key = row[keyCol]
		} else {
			key = ""
		}
		if key != "" {
			obj := rowToObj(key, m, row)
			if c.HasKey(key) {
				err = c.Join(key, obj, overwrite)
			} else {
				err = c.Create(key, obj)
			}
			if err != nil {
				return err
			}
			keys = append(keys, key)
		}
	}

	// Update the frame's keys
	f.Keys = mergeKeys(f.Keys, keys)
	// Save Frame so it can be regenerated later
	err = c.setFrame(frameName, f)
	if err != nil {
		return err
	}
	// Finally update the frame
	return c.Reframe(frameName, keys, verbose)
}
