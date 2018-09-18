package dataset

import (
	"os"
	"path"
	"testing"

	// Caltech Library packages
	"github.com/caltechlibrary/pairtree"
)

func TestPairtree(t *testing.T) {
	cName := "testdata/pairtree_test.ds"
	c, err := InitCollection(cName, PAIRTREE_LAYOUT)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key := "one"
	value := []byte(`{"one":1}`)
	err = c.CreateJSON(key, value)
	if err != nil {
		t.Errorf("failed to create %q, %s", key, err)
		t.FailNow()
	}
	pair := pairtree.Encode(key)
	stat, err := os.Stat(path.Join(c.workPath, "pairtree", pair))
	if err != nil {
		t.Errorf("failed to find %q, %s", key, err)
		t.FailNow()
	}
	if stat.IsDir() == false {
		t.Errorf("expected true, got false for %q", path.Join(c.workPath, "pairtree", pair))
	}
	stat, err = os.Stat(path.Join(cName, "pairtree", pair, key+".json"))
	if err != nil {
		t.Errorf("Expected to find %q, errored out with %s", path.Join(c.workPath, "pairtree", pair, key+"json"), err)
		t.FailNow()
	}
	// Looks like we passed so clean things up...
	os.RemoveAll(path.Join(c.workPath))
}
