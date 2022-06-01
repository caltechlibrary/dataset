package dataset

import (
	"os"
	"path"
	"testing"
)

func TestPTStoreVersioning(t *testing.T) {
	cName := path.Join("testout", "vmajor.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	if err := c.SetVersioning("major"); err != nil {
		t.Errorf("failed to set versioning to major for %q, %s", cName, err)
		t.FailNow()
	}
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()

	if c.Versioning != "major" {
		t.Errorf("Versioning wasn't set to major, %q", c.Versioning)
		t.FailNow()
	}
	key := "123"
	obj := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "Hi There!",
	}
	if err := c.Create(key, obj); err != nil {
		t.Errorf("Expected c.Create(%q, obj) to work, %s", key, err)
		t.FailNow()
	}
	versions, err := c.Versions(key)
	if len(versions) != 1 {
		t.Errorf("Expected a single version number got %+v", versions)
		t.FailNow()
	}
	if versions[0] != "1.0.0" {
		t.Errorf("Expected 1.0.0 version got %q for %s", versions[0], key)
		t.FailNow()
	}
}

func TestSQLStoreVersioning(t *testing.T) {
	t.Errorf("TestSQLStoreVersioning() not implemented")
}
