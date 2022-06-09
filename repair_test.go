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
	"strings"
	"testing"

	"github.com/caltechlibrary/dataset/pairtree"
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

func TestRepair(t *testing.T) {
	verbose := false
	o := map[string]interface{}{}
	o["a"] = 1

	// Setup a test collection and data
	cName := path.Join("testout", "test_repair.ds")
	os.RemoveAll(cName)
	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	err = c.Create("a", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["b"] = 2
	err = c.Create("b", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["c"] = 3
	err = c.Create("c", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Break the collection by removing a file from disc.
	p, err := c.PTStore.DocPath("b")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if err := os.RemoveAll(p); err != nil {
		t.Errorf("failed to remove %q, %s", p, err)
		t.FailNow()
	}
	if _, err := os.Stat(p); err == nil {
		t.Errorf("removed failed, for %q", p)
		t.FailNow()
	}
	cnt := c.Length()
	if cnt != 3 {
		t.Errorf("Expected 3, got %d", cnt)
		t.FailNow()
	}
	c.Close()

	err = Repair(cName, verbose)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	c, err = Open(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	cnt = c.Length()
	if cnt != 2 {
		t.Errorf("Expected 2, got %d, check for %q", cnt, p)
		keys, _ := c.Keys()
		t.Errorf("Found the following keys\n%s", strings.Join(keys, "\n"))
		t.Errorf("Detecting missing JSON document failed")
		t.FailNow()
	}
	c.Close()

	if err := Analyzer(cName, verbose); err != nil {
		t.Errorf("Expected no errors on Check of %q, got %s", cName, err)
		t.FailNow()
	}
}

func TestRepairLikeCLI(t *testing.T) {
	cName := path.Join("testout", "myfix.ds")
	csvName := path.Join("testout", "myfix.csv")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	if _, err := os.Stat(csvName); err == nil {
		os.RemoveAll(csvName)
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
		t.Errorf("Failed to create %q, %s", cName, err)
		t.FailNow()
	}
	for k, v := range data {
		if err := c.Create(k, v); err != nil {
			t.Errorf("Failed to setup record %q in %q -> %+v", k, cName, v)
			t.FailNow()
		}
	}
	expected64 := int64(2)
	got64 := c.Length()
	if expected64 != got64 {
		t.Errorf("Expected %d, got %d for count of test data", expected64, got64)
		t.FailNow()
	}
	c.Close()

	// Analyzer should show everything is OK
	verbose := false
	err = Analyzer(cName, verbose)
	if err != nil {
		t.Errorf("Analyzer(%q) should have been, %s", cName, err)
		t.FailNow()
	}

	// Break the collection and see if Analyzer returns error.
	collectionJSON := path.Join(cName, "collections.json")
	err = os.RemoveAll(collectionJSON)
	if err != nil {
		t.Errorf("Failed to remove %q for testing analysis, %s", collectionJSON, err)
		t.FailNow()
	}
	fmt.Printf("\n\nDEBUG collection.json %q was removed\n", collectionJSON)
	_, err = os.Stat(collectionJSON)                              // DEBUG
	fmt.Printf("DEBUG os.Stat(%q) -> %+v\n", collectionJSON, err) // DEBUG

	// Test case of missing collections.json
	fmt.Printf("DEBUG start analyzer on %q\n", cName)
	err = Analyzer(cName, verbose)
	if err == nil {
		fmt.Printf("DEBUG Analyzer(%q, %t) -> %+v\n", cName, verbose, err)
		t.Errorf("Analyzer(%q) has NOT returned the expected error of missing collection.json", cName)
		_, err = os.Stat(collectionJSON)
		if err == nil {
			t.Errorf("%q should be missing.", collectionJSON)
		}
	}
	fmt.Printf("\nDEBUG done analyzer on %q\n\n", cName)

	if _, err := os.Stat(collectionJSON); err == nil {
		t.Errorf("%q was magically recreated\n", collectionJSON)
		t.FailNow()
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
