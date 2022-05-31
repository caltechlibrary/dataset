package dataset

import (
	"os"
	"path"
	"testing"
)

func TestClone(t *testing.T) {
	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}

	cName, dsnURI := path.Join("testout", "zbs1.ds"), ""
	// cleanup stale data
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Failed to create seed collection %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()

	// Populate our seed collection
	for k, v := range testRecords {
		if err := c.Create(k, v); err != nil {
			t.Errorf("Could not create %q in %q (seed collection), %s", k, cName, err)
		}
	}
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("Could not get keys from %q, %s", cName, err)
		t.FailNow()
	}

	// Make clone collection
	ncName, ncDsnURI := path.Join("testout", "zbs2.ds"), "sqlite://testout/zbs2.ds/collections.db"
	if _, err := os.Stat(ncName); err == nil {
		os.RemoveAll(ncName)
	}
	if err := c.Clone(ncName, ncDsnURI, keys[0:5], false); err != nil {
		t.Errorf("clone failed, %q to %q, %s", cName, ncName, err)
		t.FailNow()
	}

	// Make sure clone has records.
	nc, err := Open(ncName)
	if err != nil {
		t.Errorf("failed to open clone %q, %s", ncName, err)
		t.FailNow()
	}
	defer nc.Close()
	for _, key := range keys[0:5] {
		if c.HasKey(key) != nc.HasKey(key) {
			t.Errorf("Expected %q in %q %t, got %q in %q %t", key, cName, c.HasKey(key), key, ncName, nc.HasKey(key))
		}
	}
}

func TestCloneSample(t *testing.T) {
	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}
	p := "testout"
	cName := path.Join(p, "test_zbs_characters.ds")
	trainingName := path.Join(p, "test_zbs_training.ds")
	testName := path.Join(p, "test_zbs_test.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	if _, err := os.Stat(trainingName); err == nil {
		os.RemoveAll(trainingName)
	}
	if _, err := os.Stat(testName); err == nil {
		os.RemoveAll(testName)
	}

	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("Can't create %s, %s", cName, err)
		t.FailNow()
	}
	for key, value := range testRecords {
		err := c.Create(key, value)
		if err != nil {
			t.Errorf("Can't add %s to %s, %s", key, cName, err)
			t.FailNow()
		}
	}
	cnt := int(c.Length())
	trainingSize := 4
	testSize := cnt - trainingSize
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("Expected keys in collection to clone, %s", err)
		t.FailNow()
	}
	if err := c.CloneSample(trainingName, "", testName, "", keys, trainingSize, false); err != nil {
		t.Errorf("Failed to create samples %s (%d) and %s, %s", trainingName, trainingSize, testName, err)
	}
	training, err := Open(trainingName)
	if err != nil {
		t.Errorf("Could not open %s, %s", trainingName, err)
		t.FailNow()
	}
	defer training.Close()
	test, err := Open(testName)
	if err != nil {
		t.Errorf("Could not open %s, %s", testName, err)
		t.FailNow()
	}
	defer test.Close()

	if trainingSize != int(training.Length()) {
		t.Errorf("Expected %d, got %d for %s", trainingSize, training.Length(), trainingName)
	}
	if testSize != int(test.Length()) {
		t.Errorf("Expected %d, got %d for %s", testSize, test.Length(), testName)
	}

	keys, _ = c.Keys()
	for _, key := range keys {
		switch {
		case training.HasKey(key) == true:
			if test.HasKey(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		case test.HasKey(key) == true:
			if training.HasKey(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		default:
			t.Errorf("Could not find %s in %s or %s", key, trainingName, testName)
		}
	}
}
