//
// Package dataset is a go package for managing JSON documents stored on disc
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"os"
	"testing"
)

func TestGenerateBucketNames(t *testing.T) {
	buckets := GenerateBucketNames(DefaultAlphabet, 3)
	for _, val := range buckets {
		if len(val) != 3 {
			t.Errorf("Should have a name of length 3. %q", val)
		}
	}
}

func TestPickBucketName(t *testing.T) {
	alphabet := "ab"
	buckets := GenerateBucketNames(alphabet, 2)
	expected := []string{"aa", "ab", "ba", "bb"}

	for i, expect := range expected {
		// simulate document count of doc added
		docNo := i
		result := pickBucket(buckets, docNo)
		if result != expect {
			t.Errorf("docNo %d expect %s, got %s", docNo, expect, result)
		}
	}
}

func TestCollection(t *testing.T) {
	colName := "testdata/col1"
	alphabet := "ab"
	buckets := GenerateBucketNames(alphabet, 2)
	if len(buckets) != 4 {
		t.Errorf("Should have four buckets %+v", buckets)
		t.FailNow()
	}

	// Remove any pre-existing test data
	os.RemoveAll(colName)

	// Create a new collection
	collection, err := Create(colName, buckets)
	if err != nil {
		t.Errorf("error Create() a collection %q", err)
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
	err = collection.Close()
	if err != nil {
		t.Errorf("error Close() a collection %q", err)
		t.FailNow()
	}

	// Now open the existing collection of colName
	collection, err = Open(colName)
	if err != nil {
		t.Errorf("error Open() a collection %q", err)
		t.FailNow()
	}

	if len(collection.KeyMap) > 0 {
		t.Errorf("expected 0 keys, got %d", len(collection.KeyMap))
	}
	rec1 := map[string]string{
		"name":  "freda",
		"email": "freda@inverness.example.org",
	}
	err = collection.Create("freda", rec1)
	if err != nil {
		t.Errorf("collection.Create(), %s", err)
		t.FailNow()
	}
	p, err := collection.DocPath("freda")
	if err != nil {
		t.Errorf("Should have docpath for %s, %s", "freda", err)
		t.FailNow()
	}
	if _, err := os.Stat(p); os.IsNotExist(err) == true {
		t.Errorf("Should have saved %s to disc at %s", "freda", p)
		t.FailNow()
	}
	if len(collection.KeyMap) != 1 {
		t.Errorf("expected 1 key, got %+v", collection)
		t.FailNow()
	}
	keys := collection.Keys()
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %+v", keys)
		t.FailNow()
	}

	// Create an empty record, then read it again to compare
	var rec2 map[string]string
	err = collection.Read("freda", &rec2)
	if err != nil {
		t.Errorf("Read(), %s", err)
		t.FailNow()
	}
	for k, expected := range rec1 {
		if val, ok := rec2[k]; ok == true {
			if expected != val {
				t.Errorf("expected %s in record, got, %s", expected, val)
				t.FailNow()
			}
		} else {
			t.Errorf("Read() missing %s in %+v, %+v", k, rec1, rec2)
			t.FailNow()
		}
	}
	// Should trigger update if a duplicate record
	err = collection.Create("freda", rec2)
	if err != nil {
		t.Errorf("Create on an existing record should just update it %+v", rec2)
		t.FailNow()
	}

	rec3 := map[string]string{}
	if err := collection.Read("freda", &rec3); err != nil {
		t.Errorf("Should have found freda in collection, %s", err)
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

	rec2["email"] = "freda@zbs.example.org"
	err = collection.Update("freda", rec2)
	if err != nil {
		t.Errorf("Could not update %s, %s", "freda", err)
		t.FailNow()
	}

	rec4 := map[string]string{}
	if err := collection.Read("freda", &rec4); err != nil {
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

	err = collection.Delete("freda")
	if err != nil {
		t.Errorf("Should be able to delete %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = collection.Read("freda", &rec2)
	if err == nil {
		t.Errorf("Record should have been deleted, %+v, %s", rec2, err)
	}

	err = Delete(colName)
	if err != nil {
		t.Errorf("Couldn't remove collection %s, %s", colName, err)
	}
}

func TestComplexKeys(t *testing.T) {
	colName := "testdata/col2"
	buckets := GenerateBucketNames("ab", 2)
	if len(buckets) != 4 {
		t.Errorf("Should have four buckets %+v", buckets)
		t.FailNow()
	}

	// Create a new collection
	collection, err := Create(colName, buckets)
	if err != nil {
		t.Errorf("error Create() a collection %q", err)
		t.FailNow()
	}
	testRecords := map[string]interface{}{
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
		err := collection.Create(k, v)
		if err != nil {
			t.Errorf("Can't create %s <-- %s", k, v)
		}
	}
}
