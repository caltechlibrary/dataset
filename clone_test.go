package dataset

import (
	"os"
	"path"
	"testing"
)

func TestCloneSample(t *testing.T) {
	testRecords := map[string]map[string]interface{}{
		"character:1": map[string]interface{}{
			"name": "Jack Flanders",
		},
		"character:2": map[string]interface{}{
			"name": "Little Frieda",
		},
		"character:3": map[string]interface{}{
			"name": "Mojo Sam the Yoodoo Man",
		},
		"character:4": map[string]interface{}{
			"name": "Kasbah Kelly",
		},
		"character:5": map[string]interface{}{
			"name": "Dr. Marlin Mazoola",
		},
		"character:6": map[string]interface{}{
			"name": "Old Far-Seeing Art",
		},
		"character:7": map[string]interface{}{
			"name": "Chief Wampum Stompum",
		},
		"character:8": map[string]interface{}{
			"name": "The Madonna Vampira",
		},
		"character:9": map[string]interface{}{
			"name": "Domenique",
		},
		"character:10": map[string]interface{}{
			"name": "Claudine",
		},
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

	c, err := Init(cName, "", PTSTORE)
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
	if err := c.CloneSample(trainingName, testName, keys, trainingSize, false); err != nil {
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
