//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
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
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/caltechlibrary/dataset/v2/dsv1"
	"github.com/caltechlibrary/dataset/v2/pairtree"
)

func TestPairtree(t *testing.T) {
	cName := path.Join("testout", "pairtree_test.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key := "one"
	value := map[string]interface{}{"one": 1}
	err = c.Create(key, value)
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
}

func setupTestCollectionNamed(cName string) error {
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "")
	if err != nil {
		return err
	}
	defer c.Close()

	o := map[string]interface{}{}
	o["a"] = 1
	err = c.Create("a", o)
	if err != nil {
		return err
	}
	o["b"] = 2
	err = c.Create("b", o)
	if err != nil {
		return err
	}
	o["c"] = 3
	err = c.Create("c", o)
	if err != nil {
		return err
	}
	return nil
}

func removeRecord(cName string, key string) error {
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	p, err := c.PTStore.DocPath(key)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(p); err != nil {
		return err
	}
	if _, err := os.Stat(p); err == nil {
		return err
	}
	return nil
}

func TestCheck(t *testing.T) {

	// Setup a test collection and data
	cName := path.Join("testout", "test_check.ds")
	if err := setupTestCollectionNamed(cName); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	verbose := false
	// Confer check works for no problems.
	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("Expected Analyzer(%q, %t) return nil, got %s", cName, verbose, err)
		t.FailNow()
	}

	// Break the collection by removing a file from disc.
	if err := removeRecord(cName, "b"); err != nil {
		t.Errorf("failed to remove record, %s", err)
		t.FailNow()
	}

	// Check again, should return error
	if err := Analyzer(cName, verbose); err == nil {
		t.Errorf("Expected Analyzer(%q, %t) to return an error, got nil", cName, verbose)
		t.FailNow()
	}

	// Reset collection
	if err := setupTestCollectionNamed(cName); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	// Break collection by removing collection.json
	collection := path.Join(cName, "collection.json")
	if err := os.RemoveAll(collection); err != nil {
		t.Errorf("failed to remove %q, %s", collection, err)
	}

	// Recheck collection, should get error
	if err := Analyzer(cName, verbose); err == nil {
		t.Errorf("Expected Analyzer(%q, %t) to return an error, got nil", cName, verbose)
		t.FailNow()
	}
}

func collectionSize(cName string) (int64, error) {
	c, err := Open(cName)
	if err != nil {
		return int64(-1), err
	}
	defer c.Close()
	size := c.Length()
	return size, nil
}

func TestRepair(t *testing.T) {
	verbose := false
	o := map[string]interface{}{}
	o["one"] = 1

	// Setup a test collection and data
	cName := path.Join("testout", "test_repair.ds")
	if err := setupTestCollectionNamed(cName); err != nil {
		t.Errorf("Failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	// Break the collection
	x, err := collectionSize(cName)
	if err != nil {
		t.Errorf("unable to get size %q, %s", cName, err)
		t.FailNow()
	}

	// Break the collection by removing a file from disc.
	if err := removeRecord(cName, "b"); err != nil {
		t.Errorf("failed to remove record, %s", err)
		t.FailNow()
	}

	// Run repair
	err = Repair(cName, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	got, err := collectionSize(cName)
	if err != nil {
		t.Errorf("unable to get size %q, %s", cName, err)
		t.FailNow()
	}

	expected := x - 1
	if expected != got {
		t.Errorf("Expected size %d, got %d, check for %q", expected, got, cName)
		t.FailNow()
	}

	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("Expected no errors on Check of %q, got %s", cName, err)
		t.FailNow()
	}
}

func setupTestCollection2(cName string) error {
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	// Setup test data
	data := map[string]map[string]interface{}{
		"freda": map[string]interface{}{
			"Name":   "Little Freda",
			"Email":  "freda@inverness.example.edu",
			"Office": "4th Tower",
			"Count":  1,
		},
		"mojo": map[string]interface{}{
			"Name":   "Mojo Same",
			"Email":  "mojo.sam@sams-cafe.example.org",
			"Office": "At the Piano",
			"Count":  2,
		},
	}

	c, err := Init(cName, "")
	if err != nil {
		return err
	}
	defer c.Close()
	for k, v := range data {
		if err := c.Create(k, v); err != nil {
			return err
		}
	}
	expected64 := int64(2)
	got64 := c.Length()
	if expected64 != got64 {
		return err
	}
	return nil
}

func TestCheckMissingCollectionJSON(t *testing.T) {
	cName := path.Join("testout", "test_missing_collection_json.ds")
	if err := setupTestCollection2(cName); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}
	collectionJSON := path.Join(cName, "collection.json")
	if err := os.RemoveAll(collectionJSON); err != nil {
		t.Errorf("failed to remove %q, %s", collectionJSON, err)
		t.FailNow()
	}
	verbose := false
	err := Analyzer(cName, verbose)
	if err == nil {
		t.Errorf("expected Analyzer(%q, %t) to return an error after breaking collection", cName, verbose)
		t.FailNow()
	}
}

func breakCollectionJson(cName string) error {
	collectionJSON := path.Join(cName, "collection.json")
	err := os.RemoveAll(collectionJSON)
	if err != nil {
		return err
	}
	if _, err = os.Stat(collectionJSON); err == nil {
		return fmt.Errorf("removed collection.json but Stat(%q) returned nil", collectionJSON)
	}
	return nil
}

func TestRepairLikeCLI(t *testing.T) {
	cName := path.Join("testout", "test_repair_like_cli.ds")
	if err := setupTestCollection2(cName); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	// Analyzer should show everything is OK
	verbose := false
	err := Analyzer(cName, verbose)
	if err != nil {
		t.Errorf("Analyzer(%q) should have been, %s", cName, err)
		t.FailNow()
	}

	// Break the collection and see if Analyzer returns error.
	if err := breakCollectionJson(cName); err != nil {
		t.Errorf("unable to break collection %q, %s", cName, err)
		t.FailNow()
	}

	// Test case of missing collections.json
	err = Analyzer(cName, verbose)
	if err == nil {
		t.Errorf("2nd envocation Analyzer(%q) has NOT returned the expected error of missing collection.json", cName)
	}

	// Initiating a repair
	if err := Repair(cName, verbose); err != nil {
		t.Errorf("Repair(%q) should have succeeded, %s", cName, err)
		t.FailNow()
	}

	// Comfirm repair should have worked
	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("Analyzer(%q) after repair returned errors, %s", cName, err)
		t.FailNow()
	}
}

func TestV1Check(t *testing.T) {
	cName := path.Join("testout", "testv1.ds")

	// Test records
	records := map[string]map[string]interface{}{
		"kahlo-f": {
			"_Key":        "kahlo-f",
			"id":          "Kahlo-F",
			"given_name":  "Freda",
			"family_name": "Kahlo",
			"email":       "freda@arts.example.org",
			"genre":       []string{"painting"},
		},
		"rivera-d": {
			"_Key":        "rivera-d",
			"id":          "Rivera-D",
			"given_name":  "Diego",
			"family_name": "Rivera",
			"email":       "deigo@arts.example.org",
			"genre":       []string{"murials"},
		},
		"dali-s": {
			"_Key":        "dali-s",
			"id":          "Dali-S",
			"given_name":  "Salvador",
			"family_name": "Dali",
			"email":       "salvador@collectivo.example.org",
			"genre":       []string{"painting", "architecture"},
		},
		"lopez-t": {
			"_Key":        "lopez-t",
			"id":          "Lopez-T",
			"given_name":  "Thomas",
			"family_name": "Lopez",
			"email":       "mfulton@zbs.example.org",
			"genre":       []string{"playright", "speaker", "writer", "dj", "voice"},
		},
		"valdez-l": {
			"_Key":        "valdez-l",
			"id":          "Valdez-L",
			"given_name":  "Louis",
			"family_name": "Valdez",
			"email":       "lv@theatro-composeno.example.org",
			"genre":       []string{"playright", "speaker", "writer"},
		},
		"steinbeck-j": {
			"_Key":        "steinbeck-j",
			"id":          "Steinbeck-J",
			"given_name":  "John",
			"family_name": "Steinbeck",
			"email":       "jsteinbeck@shipharbor.example.org",
			"genre":       []string{"writer"},
		},
	}
	if err := dsv1.SetupV1TestCollection(cName, records); err != nil {
		t.Errorf("dsv1.SetupV1TestCollection(%q) -> %s", cName, err)
		t.FailNow()
	}

	verbose := false
	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("failed to analyze v1 dataset, %s", err)
	}
}

func TestV1Repair(t *testing.T) {
	cName := path.Join("testout", "testv1.ds")

	// Test records
	records := map[string]map[string]interface{}{
		"kahlo-f": {
			"_Key":        "kahlo-f",
			"id":          "Kahlo-F",
			"given_name":  "Freda",
			"family_name": "Kahlo",
			"email":       "freda@arts.example.org",
			"genre":       []string{"painting"},
		},
		"rivera-d": {
			"_Key":        "rivera-d",
			"id":          "Rivera-D",
			"given_name":  "Diego",
			"family_name": "Rivera",
			"email":       "deigo@arts.example.org",
			"genre":       []string{"murials"},
		},
		"dali-s": {
			"_Key":        "dali-s",
			"id":          "Dali-S",
			"given_name":  "Salvador",
			"family_name": "Dali",
			"email":       "salvador@collectivo.example.org",
			"genre":       []string{"painting", "architecture"},
		},
		"lopez-t": {
			"_Key":        "lopez-t",
			"id":          "Lopez-T",
			"given_name":  "Thomas",
			"family_name": "Lopez",
			"email":       "mfulton@zbs.example.org",
			"genre":       []string{"playright", "speaker", "writer", "dj", "voice"},
		},
		"valdez-l": {
			"_Key":        "valdez-l",
			"id":          "Valdez-L",
			"given_name":  "Louis",
			"family_name": "Valdez",
			"email":       "lv@theatro-composeno.example.org",
			"genre":       []string{"playright", "speaker", "writer"},
		},
		"steinbeck-j": {
			"_Key":        "steinbeck-j",
			"id":          "Steinbeck-J",
			"given_name":  "John",
			"family_name": "Steinbeck",
			"email":       "jsteinbeck@shipharbor.example.org",
			"genre":       []string{"writer"},
		},
	}
	if err := dsv1.SetupV1TestCollection(cName, records); err != nil {
		t.Errorf("dsv1.SetupV1TestCollection(%q) -> %s", cName, err)
		t.FailNow()
	}

	verbose := false
	if err := Repair(cName, verbose); err == nil {
		t.Errorf("Repair should fail for a version 1 dataset collection")
	}
}
