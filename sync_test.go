package dataset

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"
)

func TestMergeFromTable(t *testing.T) {
	src := []byte(`
"id","h1","h2","h3"
0,1,2,3
1,1,2,2
2,1,2,1
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

	// TODO:
	// Make sure we only have the three keys we expect
	// Make sure the values in reach JSON object only contian ._Key, .h1, .h3 and nothing else
	// Test updating table and merging without additional fields
	// Test updating adding the missing .h2 values

}

func TestMergeFrameIntoTable(t *testing.T) {
	// FIXME: Write c.MergeFrameIntoTable() implementation tests
	// FIXME: Write c.MergeFrameIntoTable()
	t.Errorf("c.MergeFrameIntoTable() not implemented.")
}
