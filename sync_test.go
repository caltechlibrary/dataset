package dataset

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	var (
		iVal int
		sVal string
	)
	src := []byte(`
"id","h1","h2","h3"
0,1,0,10
1,1,20,11
2,1,40,12
`)
	testKeys := []string{"0", "1", "2"}
	r := csv.NewReader(bytes.NewBuffer(src))
	table, err := r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	collectionName := "testdata/merge1.ds"
	frameName := "f1"
	overwrite := true
	verbose := true

	if _, err := os.Stat(collectionName); err == nil {
		err = os.RemoveAll(collectionName)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
	}
	c, err := InitCollection(collectionName, PAIRTREE_LAYOUT)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()
	f := new(DataFrame)
	f.AllKeys = true
	f.DotPaths = []string{"._Key", ".h1", ".h3"}
	f.Labels = []string{"id", "h1", "h3"}
	c.setFrame(frameName, f)

	err = c.MergeFromTable(frameName, table, overwrite, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	keys := c.Keys()
	if len(keys) != len(testKeys) {
		t.Errorf("expected %d keys, got %+v", len(testKeys), keys)
		t.FailNow()
	}

	// NOTE: Make sure grid dimensions match table minus header row
	f, err = c.getFrame(frameName)
	if len(f.Grid) != (len(table) - 1) {
		t.Errorf("expected %d rows, got %d rows", len(table), len(f.Grid))
		t.FailNow()
	}

	// NOTE: for non-header rows check the value against what we stored in
	// our collection
	for i, row := range table[1:] {
		key := row[0]
		obj := map[string]interface{}{}
		err = c.Read(key, obj)
		if err != nil {
			t.Errorf("Expected row %d, key %s in collection, %s", i, key, err)
		}
		// Check h1 value
		sVal = "1"
		if cell, ok := obj["h1"]; ok == true {
			if cell != sVal {
				t.Errorf("(h1) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("Missing h1 in row %d, key %s, obj -> %+v", i, key, obj)
		}
		// Check h2 doesn't exist
		if cell, ok := obj["h2"]; ok == true {
			t.Errorf("(h2) row %d, key %s, Unexpected value, got %s", i, key, cell)
		}
		// Check h3 value
		iVal, _ = strconv.Atoi(key)
		sVal = fmt.Sprintf("%d", iVal+10)
		if cell, ok := obj["h3"]; ok == true {
			if sVal != cell {
				t.Errorf("(h3) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("Missing h3 in row %d, key %s, obj -> %+v", i, key, obj)
		}

	}

	// NOTE: Update frame labels and dotpaths
	f.DotPaths = []string{"._Key", ".h1", ".h2", ".h4"}
	f.Labels = []string{"id", "h1", "h2", "h4"}
	c.setFrame(frameName, f)
	c.Reframe(frameName, f.Keys, false)

	// Update table values for next merge test
	for i, row := range table {
		if i == 0 {
			row = append(row, "h4", "h5")
		} else {
			row[1] = "3"
			row = append(row, "Y", "T")
		}
		table[i] = row
	}

	err = c.MergeFromTable(frameName, table, overwrite, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	// Now reconcile table and collection objects
	for i, row := range table[1:] {
		key := row[0]
		obj := map[string]interface{}{}
		err = c.Read(key, obj)
		if err != nil {
			t.Errorf("Expected row %d, key %s in collection, %s", i, key, err)
		}
		// Check h1 value
		sVal = "3"
		if cell, ok := obj["h1"]; ok == true {
			if cell != sVal {
				t.Errorf("(h1) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("Missing h1 in row %d, key %s, obj -> %+v", i, key, obj)
		}
		// Check h2 value
		iVal, _ = strconv.Atoi(key)
		sVal = fmt.Sprintf("%d", iVal*20)
		if cell, ok := obj["h2"]; ok == true {
			if sVal != cell {
				t.Errorf("(h2) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("(h2) row %d, key %s, expected a %s", i, key, sVal)
		}
		// Check h3 value
		iVal, _ = strconv.Atoi(key)
		sVal = fmt.Sprintf("%d", iVal+10)
		if cell, ok := obj["h3"]; ok == true {
			if sVal != cell {
				t.Errorf("(h3) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("Missing h3 in row %d, key %s, obj -> %+v", i, key, obj)
		}

		// Check h4 value
		sVal = "Y"
		if cell, ok := obj["h4"]; ok == true {
			//iVal, _ = strconv.Atoi(key)
			//sVal = fmt.Sprintf("%d", iVal*20)
			if sVal != cell {
				t.Errorf("(h4) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("(h4) row %d, key %s, expected a %s", i, key, sVal)
		}
		// Check h5 doesn't exist
		if cell, ok := obj["h5"]; ok == true {
			t.Errorf("(h5) row %d, key %s, Unexpected value, got %s", i, key, cell)
		}

	}
	if len(c.Keys()) != 3 {
		t.Errorf("Expected three keys in %s, got %d", collectionName, len(c.Keys()))
	}

	// TEST: Implement c.MergeIntoTable() tests
	sVal = time.Now().String()
	for _, key := range c.Keys() {
		obj := map[string]interface{}{}
		err = c.Read(key, obj)
		if err != nil {
			t.Errorf("Can't read %s from %s, %s", key, collectionName, err)
			t.FailNow()
		}
		obj["h2"] = sVal
		obj["h6"] = "F"
		err = c.Update(key, obj)
		if err != nil {
			t.Errorf("Can't update %s in %s, %s", key, collectionName, err)
			t.FailNow()
		}
	}

	rTable, err := c.MergeIntoTable(frameName, table, overwrite, verbose)
	if len(rTable) != len(table) {
		t.Errorf("expected %d rows, got %d", len(table), len(rTable))
	}
	// Make sure all table cells are accounted for
	for i, row := range table {
		if len(row) != len(rTable[i]) {
			t.Errorf("expected %d columns in row %d, got %d", len(row), i, len(rTable[i]))
		} else {
			for j, cell := range row {
				sVal = table[i][j]
				if cell != sVal {
					t.Errorf("expected row %d column %d %q, got %q", i, j, sVal, cell)
				}
			}
		}
	}

	// TEST: appending rows to table
	obj := map[string]interface{}{
		"_Key": "10",
		"h1":   "1",
		"h2":   "2",
		"h3":   "3",
		"h4":   "",
		"h5":   "4",
	}
	err = c.Create("10", obj)
	if err != nil {
		t.Errorf("%s\n", err)
		t.FailNow()
	}
	if len(c.Keys()) != 4 {
		t.Errorf("expected 4 keys, got %+v\n", c.Keys())
	}

	// NOTE: Final update frame labels and dotpaths
	f.DotPaths = []string{"._Key", ".h1", ".h2", ".h3", ".h4", ".h5", ".h6"}
	f.Labels = []string{"id", "h1", "h2", "h3", "h4", "h5", "h6"}
	c.setFrame(frameName, f)
	c.Reframe(frameName, c.Keys(), false)
	if err != nil {
		t.Errorf("%s\n", err)
		t.FailNow()
	}

	rTable, err = c.MergeIntoTable(frameName, table, overwrite, verbose)
	if len(rTable) != (len(table) + 1) {
		t.Errorf("expected %d rows, got %d", len(table)+1, len(rTable))
		//fmt.Printf("\n\trTable\n%+v\n", rTable)
		//fmt.Printf("\ttable\n%+v\n", table)
		t.FailNow()
	}
	for i, row := range table {
		for j, cell := range row {
			if rTable[i][j] != cell {
				t.Errorf("row %d, col %d, expected %q, got %q", i, j, cell, rTable[i][j])
			}
		}
	}
	lastRow := len(rTable) - 1
	tRow := []string{"10", "1", "2", "3", "", "4"}
	if len(tRow) != len(rTable[lastRow]) {
		t.Errorf("appended row, expected length %d, got %d", len(tRow), len(rTable[lastRow]))
		fmt.Printf("\texpected: %+v\n", tRow)
		fmt.Printf("\t     got: %+v\n", rTable[lastRow])
		t.FailNow()
	}
	row := rTable[lastRow]
	for j, tCell := range tRow {
		if row[j] != tCell {
			t.Errorf("row %d, col %d, excepted %q, got %q", lastRow, j, tCell, row[j])
		}
	}
}
