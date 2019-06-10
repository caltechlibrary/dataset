package dataset

import (
	"fmt"
	"log"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/tbl"
	"github.com/caltechlibrary/dotpath"
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

// labelsToHeaderRow checks the labels of a frame to make sure
// all labels are in table's header row. If not it appends the
// missing columns to the end of the header row and returns
// new header row and true if a change is needed.
func labelsToHeaderRow(f *DataFrame, table [][]interface{}) ([]string, bool) {
	changed := false
	header := []string{}
	for i, cell := range table[0] {
		val, err := tbl.ValueInterfaceToString(cell)
		if err == nil {
			header = append(header, val)
		} else {
			header = append(header, fmt.Sprintf(fmtColumnName, i+1))
			changed = true
		}
	}
	for _, label := range f.Labels {
		if strInArray(header, label) == false {
			header = append(header, label)
			changed = true
		}
	}
	return header, changed
}

// dotPathToColumnMap creates a mapping from dotpath in collection
// to column number in table by matching header row values with
// a frame's labels. Returns an error if ._Key is not identified.
func dotPathToColumnMap(f *DataFrame, table [][]interface{}) (map[string]int, error) {
	m := make(map[string]int)
	if len(f.Labels) != len(f.DotPaths) {
		return m, fmt.Errorf("corrupted frame, labels don't map to dot paths")
	}
	if len(table) < 2 {
		return m, fmt.Errorf("table is empty")
	}
	for i, val := range table[0] {
		label, err := tbl.ValueInterfaceToString(val)
		if err == nil && strings.TrimSpace(label) != "" {
			if pos, hasLabel := findLabel(f.Labels, label); hasLabel == true {
				// Find the dotpath matching the label
				dotPath := f.DotPaths[pos]
				// Map the dotpath to a column number
				m[dotPath] = i
			}
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
func rowToObj(key string, dotPathToCols map[string]int, row []interface{}) map[string]interface{} {
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

// hasKey takes a list of keys (string) and sees if key is in list
func hasKey(keys []string, key string) bool {
	for _, item := range keys {
		if item == key {
			return true
		}
	}
	return false
}

// MergeIntoTable - uses a DataFrame associated in the collection
// to map attributes into table appending new content and optionally
// overwriting existing content for rows with matching ids. Returns
// a new table (i.e. [][]interface{}) or error.
func (c *Collection) MergeIntoTable(frameName string, table [][]interface{}, overwrite bool, verbose bool) ([][]interface{}, error) {
	// Build Map dotpath to column position
	//
	// For each data row of table (i.e. row 1 through last row)
	//    get ID value
	//    if has ID && overwrite == true then replace cells values
	//    else save id for append to table
	// Update table
	f, err := c.getFrame(frameName)
	if err != nil {
		return table, err
	}
	// Makesure we have a header that supports all the Frame's
	// dotPaths and label
	headerRow, changed := labelsToHeaderRow(f, table)
	if changed {
		table[0] = tbl.RowStringToInterface(headerRow)
	}

	// Based on table's new header, calc the map of dotpath to
	// column no.
	m, err := dotPathToColumnMap(f, table)
	if err != nil {
		return table, err
	}
	dotPaths := f.DotPaths
	keyCol, _ := m["._Key"]
	key := ""
	tableKeys := []string{}
	for i, row := range table {
		//NOTE: we skip the header row
		if i == 0 {
			continue
		}
		// Get ID from row
		if keyCol < len(row) {
			key, err = tbl.ValueInterfaceToString(row[keyCol])
			if err == nil && key != "" {
				// collect the tables' row keys
				tableKeys = append(tableKeys, key)
			} else {
				if verbose {
					log.Printf("skipping row %d, invalid key found in column %d, %+v, %T", i, keyCol, row[keyCol], row[keyCol])
				}
				continue
			}
		} else {
			if verbose {
				log.Printf("skipping row %d, no key found in column %d", i, keyCol)
			}
			continue
		}
		if c.HasKey(key) {
			// Pad cells in row if necessary
			for i := len(row); i < len(headerRow); i++ {
				row = append(row, "")
			}
			obj := map[string]interface{}{}
			err := c.Read(key, obj, false)
			if err != nil {
				return table, fmt.Errorf("Can't read %s from row %d in collection", key, i)
			}
			// For each row replace cells in dotPath map to column number
			for _, p := range dotPaths {
				//NOTE: need to do this in order, so iterate over
				// f.DotPaths then get j from map m.
				j, ok := m[p]
				if ok == false {
					continue
				}
				val, err := dotpath.Eval(p, obj)
				if err == nil {
					row[j] = val
				}
			}
			// update row in table
			table[i] = row
		} else if verbose {
			log.Printf("skipping row %d, key %s not found in collection %s", i, key, c.Name)
		}
	}
	// Append rows to table if needed
	for _, key := range f.Keys {
		if hasKey(tableKeys, key) == false {
			// Generate a row to add
			row := make([]interface{}, len(headerRow)-1)
			// Get the data for the row
			obj := map[string]interface{}{}
			err = c.Read(key, obj, false)
			if err != nil {
				return table, fmt.Errorf("failed to read %q in %s, %s\n", key, c.Name, err)
			}
			// For each row replace cells in dotPath map to column number
			for p, j := range m {
				val, err := dotpath.Eval(p, obj)
				if err == nil {
					// Pad cells in row if necessary
					for j >= len(row) {
						row = append(row, nil)
					}
					row[j] = val
				}
			}
			table = append(table, row)
		}
	}
	return table, nil
}

// MergeFromTable - uses a DataFrame associated in the collection
// to map columns from a table into JSON object attributes saving the
// JSON object in the collection.  If overwrite is true then JSON objects
// for matching keys will be updated, if false only new objects will be
// added to collection. Returns an error value
func (c *Collection) MergeFromTable(frameName string, table [][]interface{}, overwrite bool, verbose bool) error {
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
	for i, row := range table[1:] {
		// get Key
		if keyCol < len(row) {
			key, err = tbl.ValueInterfaceToString(row[keyCol])
			if err != nil {
				if verbose {
					log.Printf("skipping row %d, invalid key found in column %d, %+v, %T, %s", i+2, keyCol, row[keyCol], row[keyCol], err)
				}
			} else if key == "" {
				if verbose {
					log.Printf("skipping row %d, missing key found in column %d", i+2, keyCol)
				}
			} else {
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
		} else {
			if verbose {
				log.Printf("skipping row %d, unknown key column", i)
			}
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
