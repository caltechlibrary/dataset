package dsv1

import (
	"os"
	"path"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	cName := path.Join("testout", "check_ok.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	verbose := false // true

	if err := SetupV1TestCollection(cName, records); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("Analyzer failed unexpectedly, %s", err)
		t.FailNow()
	}

	collection := path.Join(cName, "collection.json")
	os.RemoveAll(collection)

	// Check should return an error
	if err := Analyzer(cName, verbose); err == nil {
		t.Errorf("Analyzer should have failed, %s removed", collection)
		t.FailNow()
	}
}
