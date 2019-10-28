package dataset

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/tbl"
)

func TestTableColumnBehavior(t *testing.T) {
	src := []byte(`
first, second,third, fourth
1,,3,
2,2,2,
3,3,
`)
	r := csv.NewReader(bytes.NewBuffer(src))
	r.FieldsPerRecord = -1
	csvTable, err := r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	table := tbl.TableStringToInterface(csvTable)

	expectedObjs := map[string]map[string]interface{}{}
	expectedObjs["1"] = map[string]interface{}{
		"_Key":  "1",
		"first": 1,
		"third": 3,
	}
	expectedObjs["2"] = map[string]interface{}{
		"_Key":   "2",
		"first":  2,
		"second": 2,
		"third":  2,
	}
	expectedObjs["3"] = map[string]interface{}{
		"_Key":   "3",
		"first":  3,
		"second": 3,
	}

	for i, row := range table {
		if i > 0 {
			key, err := tbl.ValueInterfaceToString(row[0])
			if err != nil {
				t.Errorf("expected (%T) %+v to be convertable to string, %s", row[0], row[0], err)
				continue
			}
			for j, fieldName := range []string{"first", "second", "third", "fourth"} {
				if eObj, ok := expectedObjs[key]; ok == true {
					obj := eObj
					if len(obj) > len(row) {
						t.Errorf("row %d is short cells, obj %+v, row %+v", i, obj, row)
					}
					if val, ok := obj[fieldName]; ok == true && val != row[j] {
						t.Errorf("row %d, col %d, expected (%T) %+v, got (%T) %+v", i, j, val, val, row[j], row[j])
					}
				} else {
					t.Errorf("unexpected key %q in row %d table %+v\n", key, i, table)
				}
			}
		}
	}
}

func TestMerge(t *testing.T) {
	var (
		iVal int
		sVal string
	)
	overwrite := true
	verbose := true

	src := []byte(`
"id","h1","h2","h3"
0,1,0,10
1,1,20,11
2,1,40,12
`)
	testKeys := []string{"0", "1", "2"}
	r := csv.NewReader(bytes.NewBuffer(src))
	csvTable, err := r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	table := tbl.TableStringToInterface(csvTable)
	collectionName := "testdata/merge1.ds"
	frameName := "f1"

	if _, err := os.Stat(collectionName); err == nil {
		err = os.RemoveAll(collectionName)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
	}
	c, err := InitCollection(collectionName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()
	// Manually create a frame to merge with.
	f, err := c.FrameCreate(frameName, []string{}, []string{"._Key", ".h1", ".h3"}, []string{"id", "h1", "h3"}, verbose)

	err = c.MergeFromTable(frameName, table, overwrite, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	keys := c.Keys()
	if len(keys) != len(testKeys) {
		t.Errorf("expected %d keys, got %+v", len(testKeys), keys)

		w := csv.NewWriter(os.Stdout)
		fmt.Printf("table:\n")
		w.WriteAll(tbl.TableInterfaceToString(table))
		w.Flush()
		fmt.Printf("testKeys: %+v\n", testKeys)
		fmt.Printf("collection name: %s\n", collectionName)

		t.FailNow()
	}

	// NOTE: Make sure grid dimensions match table minus header row
	c.Reframe(frameName, keys, false)
	f, err = c.getFrame(frameName)
	if err != nil {
		t.Errorf("failed to get frame %s, %s", frameName, err)
		t.FailNow()
	}
	if len(f.ObjectMap) != len(f.Keys) {
		t.Errorf("Expected %d objects for %d keys, got %+v\n", len(f.ObjectMap), len(f.Keys), f)
		t.FailNow()
	}
	grid := f.Grid(false)
	if len(grid) != (len(table) - 1) {
		t.Errorf("expected %d rows, got %d rows", len(table), len(grid))
		t.FailNow()
	}

	// NOTE: for non-header rows check the value against what we stored in
	// our collection
	for i, row := range table[1:] {
		key, err := tbl.ValueInterfaceToString(row[0])
		if err != nil {
			t.Errorf("Expected row %d, key (%T) %+v to be string, %s", i, key, key, err)
		}
		obj := map[string]interface{}{}
		err = c.Read(key, obj, false)
		if err != nil {
			t.Errorf("Expected row %d, key %s in collection, %s", i, key, err)
		}
		// Check h1 value
		sVal := "1"
		if cell, ok := obj["h1"]; ok == true {
			if cell.(json.Number).String() != sVal {
				t.Errorf("(h1) row %d, key %s, expected %s, got (%T) %+v", i, key, sVal, cell, cell)
			}
		} else {
			t.Errorf("Missing h1 in row %d, key %s, obj -> %+v", i, key, obj)
		}
		// Check h2 doesn't exist
		if cell, ok := obj["h2"]; ok == true {
			t.Errorf("(h2) row %d, key %s, Unexpected value, got (%T) %+v", i, key, cell, cell)
		}
		// Check h3 value
		iVal, _ = strconv.Atoi(key)
		sVal = fmt.Sprintf("%d", iVal+10)
		if cell, ok := obj["h3"]; ok == true {
			if cell.(json.Number).String() != sVal {
				t.Errorf("(h3) row %d, key %s, expected %s, got (%T) %s", i, key, sVal, cell, cell)
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
		key, err := tbl.ValueInterfaceToString(row[0])
		if err != nil {
			t.Errorf("Expected row %d, key (%T) %+v, of type string,%s", i, key, key, err)
		}
		obj := map[string]interface{}{}
		err = c.Read(key, obj, false)
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
			if cell.(json.Number).String() != sVal {
				t.Errorf("(h2) row %d, key %s, expected %s, got %s", i, key, sVal, cell)
			}
		} else {
			t.Errorf("(h2) row %d, key %s, expected a %s", i, key, sVal)
		}
		// Check h3 value
		iVal, _ = strconv.Atoi(key)
		sVal = fmt.Sprintf("%d", iVal+10)
		if cell, ok := obj["h3"]; ok == true {
			if cell.(json.Number).String() != sVal {
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
		err = c.Read(key, obj, false)
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
	var cVal interface{}
	for i, row := range table {
		if len(row) != len(rTable[i]) {
			t.Errorf("expected %d columns in row %d, got %d", len(row), i, len(rTable[i]))
		} else {
			for j, cell := range row {
				cVal = table[i][j]
				if cell != cVal {
					t.Errorf("expected row %d column %d %+v, got %+v", i, j, cVal, cell)
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
		w := csv.NewWriter(os.Stdout)
		fmt.Printf("table:\n")
		w.WriteAll(tbl.TableInterfaceToString(table))
		w.Flush()
		fmt.Printf("rTable:\n")
		w.WriteAll(tbl.TableInterfaceToString(rTable))
		w.Flush()
		t.FailNow()
	}
	row := rTable[lastRow]
	for j, tCell := range tRow {
		if row[j] != tCell {
			t.Errorf("row %d, col %d, excepted %q, got %q", lastRow, j, tCell, row[j])
		}
	}
}

func TestAddedColumns(t *testing.T) {
	expectedCSV := []byte(`
id,one,two,three,four,five
0,A,B,C,D,E
1,B,C,D,E,F
2,C,D,E,F,G
3,D,E,F,G,H
4,E,F,G,H,I
`)

	initialCSV := []byte(`
id,one,two
0,A,B
1,B,C
2,C,D
3,D,E
4,E,F
`)

	r := csv.NewReader(bytes.NewBuffer(initialCSV))
	csvTable, err := r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	initialTbl := tbl.TableStringToInterface(csvTable)

	r = csv.NewReader(bytes.NewBuffer(expectedCSV))
	csvTable, err = r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	expectedTbl := tbl.TableStringToInterface(csvTable)

	collectionName := "testdata/merge2.ds"
	frameName := "f1"
	useHeaderRow := true
	overwrite := true
	verbose := true

	if _, err := os.Stat(collectionName); err == nil {
		err = os.RemoveAll(collectionName)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
	}
	c, err := InitCollection(collectionName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	buf := bytes.NewBuffer(initialCSV)
	_, err = c.ImportCSV(buf, 0, useHeaderRow, overwrite, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	keys := c.Keys()
	if len(keys) == 0 {
		t.Errorf("Import failed")
		t.FailNow()
	}

	f, err := c.FrameCreate(frameName, keys, []string{"._Key", ".one", ".two"}, []string{"id", "one", "two"}, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	err = c.SaveFrame(frameName, f)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	// Setup collection updates to test MergeIntoTable()
	keys = c.Keys()
	fieldVals := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	for i, key := range []string{"0", "1", "2", "3", "4"} {
		obj := map[string]interface{}{}
		err = c.Read(key, obj, false)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		obj["one"] = fieldVals[i]
		obj["two"] = fieldVals[i+1]
		obj["three"] = fieldVals[i+2]
		obj["four"] = fieldVals[i+3]
		obj["five"] = fieldVals[i+4]

		err = c.Update(key, obj)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
	}
	f.DotPaths = []string{"._Key", ".one", ".two", ".three", ".four", ".five"}
	f.Labels = []string{"id", "one", "two", "three", "four", "five"}
	err = c.SaveFrame(frameName, f)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	resultTbl, err := c.MergeIntoTable(frameName, initialTbl, overwrite, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	// Verify resulting table
	if len(expectedTbl) != len(resultTbl) {
		t.Errorf("expected table is not same size as result table")
		t.FailNow()
	}
	// Compare cell by cell
	for i, row := range expectedTbl {
		if len(resultTbl[i]) != len(row) {
			t.Errorf("row %d is different lengths, %+v, got %+v", i, expectedTbl[i], resultTbl[i])
			w := csv.NewWriter(os.Stdout)
			fmt.Printf("expectedTbl:\n")
			w.WriteAll(tbl.TableInterfaceToString(expectedTbl))
			w.Flush()
			fmt.Printf("resultTbl:\n")
			w.WriteAll(tbl.TableInterfaceToString(resultTbl))
			w.Flush()
			t.FailNow()
		}
		for j, cell := range row {
			// NOTE: In our test cases column 0 is key and
			// will be a string regardless of either it is numeric.
			if j == 0 {
				sCell, _ := tbl.ValueInterfaceToString(cell)
				rCell, _ := tbl.ValueInterfaceToString(resultTbl[i][j])
				if sCell != rCell {
					t.Errorf("row %d, col %d, expected (%T) %+v, got (%T) %+v", i, j, cell, cell, resultTbl[i][j], resultTbl[i][j])
				}
			} else if cell != resultTbl[i][j] {
				t.Errorf("row %d, col %d, expected (%T) %+v, got (%T) %+v", i, j, cell, cell, resultTbl[i][j], resultTbl[i][j])
				w := csv.NewWriter(os.Stdout)
				fmt.Printf("expectedTbl:\n")
				w.WriteAll(tbl.TableInterfaceToString(expectedTbl))
				w.Flush()
				fmt.Printf("resultTbl:\n")
				w.WriteAll(tbl.TableInterfaceToString(resultTbl))
				w.Flush()
				t.FailNow()
			}
		}
	}

}
