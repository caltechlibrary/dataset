//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2021, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollection(t *testing.T) {
	colName := "testdata/test_collection.ds"
	// Remove any pre-existing test data
	os.RemoveAll(colName)

	// Create a new collection
	c, err := Init(colName)
	if err != nil {
		t.Errorf("error create() a collection %q", err)
		t.FailNow()
	}
	// Make sure directories were create for col1
	if fInfo, err := os.Stat(colName); err != nil {
		t.Errorf("%s was not created, %s", colName, err)
		t.FailNow()
	} else if fInfo.IsDir() != true {
		t.Errorf("%s is supposed to be a directory!", colName)
		t.FailNow()
	}
	err = c.Close()
	if err != nil {
		t.Errorf("error Close() a collection %q", err)
		t.FailNow()
	}

	// Now open the existing collection of colName
	c, err = Open(colName)
	if err != nil {
		t.Errorf("error Open() a collection %q", err)
		t.FailNow()
	}

	if len(c.KeyMap) > 0 {
		t.Errorf("expected 0 keys, got %d", len(c.KeyMap))
	}
	testData := []map[string]interface{}{}
	src := `[
		{
			"id": "Kahlo-F",
			"given_name":  "Freda",
			"last_name": "Kahlo",
			"email": "freda@arts.example.org"
		},
		{
			"id": "Rivera-D",
			"given_name": "Diego",
			"family_name": "Rivera",
			"email": "deigo@arts.example.org"
		},
		{
			"id": "Dali-S",
			"given_name": "Salvador",
			"family_name": "Dali",
			"email": "salvador@collectivo.example.org"
		}
]`
	if err := json.Unmarshal([]byte(src), &testData); err != nil {
		t.Errorf("Failed to marshal test data, %s", err)
		t.FailNow()
	}

	for _, rec := range testData {
		if k, ok := rec["id"]; ok == true {
			id := k.(string)
			err = c.Create(id, rec)
			if err != nil {
				t.Errorf("%q: collection.Create(), %s", c.Name, err)
				t.FailNow()
			}
			if c.KeyExists(id) == false {
				t.Errorf("%q was not created in %q, no error valuye returned", id, c.Name)
				t.FailNow()
			}
			p, err := c.DocPath(id)
			if err != nil {
				t.Errorf("%q: Should have docpath for %s, %s", c.Name, id, err)
				t.FailNow()
			}
			if _, err := os.Stat(p); os.IsNotExist(err) == true {
				t.Errorf("%q: Should have saved %s to disc at %s", c.Name, id, p)
				t.FailNow()
			}
		}
	}

	if len(c.KeyMap) != 3 {
		t.Errorf("%q: expected 1 key, got %+v", c.Name, c)
		t.FailNow()
	}
	keys := c.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %+v", keys)
		t.FailNow()
	}

	// Create an empty record, then read it again to compare
	keyName := "Kahlo-F"
	rec2 := map[string]interface{}{}
	err = c.Read(keyName, rec2, false)
	if err != nil {
		t.Errorf("%q: Read(), %s", c.Name, err)
		t.FailNow()
	}
	rec1 := testData[0]
	for k, expected := range rec1 {
		if val, ok := rec2[k]; ok == true {
			if expected != val {
				t.Errorf("%q: expected %s in record, got, %s", c.Name, expected, val)
				t.FailNow()
			}
		} else {
			t.Errorf("%q: Read() missing %s in %+v, %+v", c.Name, k, rec1, rec2)
			t.FailNow()
		}
	}
	// Should trigger update if a duplicate record
	err = c.Create(keyName, rec2)
	if err == nil {
		t.Errorf("%q: Create not allow creationg on an existing record, %s --> %+v", c.Name, keyName, rec2)
		t.FailNow()
	}

	rec3 := map[string]interface{}{}
	if err := c.Read(keyName, rec3, false); err != nil {
		t.Errorf("%q: Should have found freda in collection, %s", c.Name, err)
		t.FailNow()
	}
	for k2, v2 := range rec2 {
		if v3, ok := rec3[k2]; ok == true {
			if v2 != v3 {
				t.Errorf("Expected v2 %+v, got v3 %+v", v2, v3)
			}
		} else {
			t.Errorf("missing key %s r3 in %+v <- r2: %+v \n", k2, rec3, rec2)
		}
	}

	rec2["email"] = "freda@collectivo.example.org"
	err = c.Update(keyName, rec2)
	if err != nil {
		t.Errorf("%s: Could not update %s, %s", c.Name, "freda", err)
		t.FailNow()
	}

	rec4 := map[string]interface{}{}
	if err := c.Read(keyName, rec4, false); err != nil {
		t.Errorf("Should have found freda in collection, %s", err)
		t.FailNow()
	}
	for k2, v2 := range rec2 {
		if v4, ok := rec4[k2]; ok == true {
			if v2 != v4 {
				t.Errorf("Expected v2 %+v, got v4 %+v", v2, v4)
			}
		} else {
			t.Errorf("missing key %s rec4 in %+v <- rec2: %+v \n", k2, rec4, rec2)
		}
	}

	err = c.Delete(keyName)
	if err != nil {
		t.Errorf("Should be able to delete %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = c.Read(keyName, rec2, false)
	if err == nil {
		t.Errorf("Record should have been deleted, %+v, %s", rec2, err)
	}

	err = deleteCollection(colName)
	if err != nil {
		t.Errorf("Couldn't remove collection %s, %s", colName, err)
	}
}

func TestComplexKeys(t *testing.T) {
	colName := "testdata/pairtree_layout/col2.ds"
	// remove any stale test collection collection first...
	os.RemoveAll(colName)

	// Create a new collection
	c, err := Init(colName)
	if err != nil {
		t.Errorf("error Create() a collection %q", err)
		t.FailNow()
	}
	testRecords := map[string]map[string]interface{}{
		"agent:person:1": map[string]interface{}{
			"name": "George",
			"id":   25,
		},
		"agent:person:2": map[string]interface{}{
			"name": "Carl",
			"id":   2523,
		},
		"agent:person:3333": map[string]interface{}{
			"name": "Mac",
			"id":   2,
		},
		"agent:person:29994": map[string]interface{}{
			"name": "Fred",
			"id":   9925,
		},
		"agent:person:29": map[string]interface{}{
			"name": "Mike",
			"id":   81,
		},
		"agent:person:100": map[string]interface{}{
			"name": "Tim",
			"id":   8,
		},
		"agent:person:101": map[string]interface{}{
			"name": "Kim",
			"id":   101,
		},
	}

	for k, v := range testRecords {
		err := c.Create(k, v)
		if err != nil {
			t.Errorf("Can't create %s <-- %s : %s", k, v, err)
		}
	}
}

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
	p := "testdata/pairtree_layout"
	cName := path.Join(p, "test_zbs_characters.ds")
	trainingName := path.Join(p, "test_zbs_training.ds")
	testName := path.Join(p, "test_zbs_test.ds")
	os.RemoveAll(cName)
	os.RemoveAll(trainingName)
	os.RemoveAll(testName)

	c, err := Init(cName)
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
	cnt := c.Length()
	trainingSize := 4
	testSize := cnt - trainingSize
	keys := c.Keys()
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

	if trainingSize != training.Length() {
		t.Errorf("Expected %d, got %d for %s", trainingSize, training.Length(), trainingName)
	}
	if testSize != test.Length() {
		t.Errorf("Expected %d, got %d for %s", testSize, test.Length(), testName)
	}

	keys = c.Keys()
	for _, key := range keys {
		switch {
		case training.KeyExists(key) == true:
			if test.KeyExists(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		case test.KeyExists(key) == true:
			if training.KeyExists(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		default:
			t.Errorf("Could not find %s in %s or %s", key, trainingName, testName)
		}
	}
}

func TestCaseHandling(t *testing.T) {
	// Setup a test collection and data
	cName := path.Join("testdata", "test_case_handling.ds")
	os.RemoveAll(cName)
	c, err := Init(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o := map[string]interface{}{}
	o["a"] = 1
	err = c.Create("A", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["b"] = 2
	err = c.Create("B", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["c"] = 3
	err = c.Create("C", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Get back a list of keys, should all be lowercase.const
	keys := c.Keys()
	for _, key := range keys {
		if key == strings.ToUpper(key) {
			t.Errorf("Expected lower case %q, got %q", strings.ToLower(key), key)
		}
		p, err := c.DocPath(strings.ToUpper(key))
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		// Check if p has the OS's separator.
		if !strings.Contains(p, fmt.Sprintf("%c", filepath.Separator)) {
			t.Errorf("Path seperator does not match host OS, %q <- %c", p, filepath.Separator)
			t.FailNow()
		}
	}
	cnt := c.Length()
	if cnt != 3 {
		t.Errorf("Expected 3, got %d", cnt)
		t.FailNow()
	}
	c.Close()
}
