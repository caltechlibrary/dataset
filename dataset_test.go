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
	"strings"
	"testing"
)

func TestGenerateBucketNames(t *testing.T) {
	alphabet := "abc"
	buckets := GenerateBucketNames(alphabet, 3)
	for _, val := range buckets {
		if len(val) != 3 {
			t.Errorf("Should have a name of length 3. %q", val)
		}
	}
}

func TestIntToBucketName(t *testing.T) {
	alphabet := "ab"
	buckets := GenerateBucketNames(alphabet, 2)

	for i, expected := range []string{"aa", "ab", "ba", "bb"} {
		result := intToBucketName(i, buckets)
		if strings.Compare(result, expected) != 0 {
			t.Errorf("%d expected %s, got %s", i, expected, result)
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
	err = collection.Close()
	if err != nil {
		t.Errorf("error Close() a collection %q", err)
		t.FailNow()
	}
	collection, err = Open(colName)
	if err != nil {
		t.Errorf("error Open() a collection %q", err)
		t.FailNow()
	}

	if len(collection.Keys) > 0 {
		t.Errorf("expected 0 keys, got %d", len(collection.Keys))
	}
	rec1 := map[string]string{
		"name":  "freda",
		"email": "freda@inverness.example.org",
	}
	err = collection.Create("freda.json", rec1)
	if err != nil {
		t.Errorf("collection.Create(), %s", err)
		t.FailNow()
	}
	if len(collection.Keys) != 1 {
		t.Errorf("expected 1 key, got %+v", collection)
		t.FailNow()
	}
	// Clear record, then read it again
	rec2 := map[string]string{}
	err = collection.Read("freda.json", &rec2)
	if err != nil {
		t.Errorf("Read(), %s", err)
	}
	for k, expected := range rec1 {
		if val, ok := rec2[k]; ok == true {
			if strings.Compare(expected, val) != 0 {
				t.Errorf("expected %s in record, got, %s", expected, val)
			}
		} else {
			t.Errorf("Read() missing %s in %+v", k, rec2)
			t.FailNow()
		}
	}
	rec2["email"] = "freda@zbs.example.org"
	// Should fail if we try to create a duplicate record
	err = collection.Create("freda.json", rec2)
	if err == nil {
		t.Errorf("Should not beable to create a duplicate %+v", rec2)
		t.FailNow()
	}
	err = collection.Update("freda.json", rec2)
	if err != nil {
		t.Errorf("Could not update %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = collection.Delete("freda.json")
	if err != nil {
		t.Errorf("Should be able to delete %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = collection.Read("freda.json", &rec2)
	if err == nil {
		t.Errorf("Record should have been deleted, %+v, %s", rec2, err)
	}

	err = Delete(colName)
	if err != nil {
		t.Errorf("Couldn't remove collection %s, %s", colName, err)
	}
}
