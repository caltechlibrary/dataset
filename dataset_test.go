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
	if _, err := os.Stat(colName); os.IsNotExist(err) == true {
		t.Errorf("%s was not created", colName)
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
			if strings.Compare(expected, val) != 0 {
				t.Errorf("expected %s in record, got, %s", expected, val)
				t.FailNow()
			}
		} else {
			t.Errorf("Read() missing %s in %+v, %+v", k, rec1, rec2)
			t.FailNow()
		}
	}
	rec2["email"] = "freda@zbs.example.org"
	// Should fail if we try to create a duplicate record
	err = collection.Create("freda", rec2)
	if err == nil {
		t.Errorf("Should not beable to create a duplicate %+v", rec2)
		t.FailNow()
	}
	err = collection.Update("freda", rec2)
	if err != nil {
		t.Errorf("Could not update %s, %s", "freda", err)
		t.FailNow()
	}

	// Run subtests of select list behavior
	if ok := selectListBehavior(t, collection); ok == false {
		t.FailNow()
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

func selectListBehavior(t *testing.T, c *Collection) bool {
	// Select collection level sellect lists
	keys1 := c.Keys()
	selectLists := c.Lists()
	if len(selectLists) != 1 {
		t.Errorf("Have unexpected select lists, %+v", selectLists)
		return false
	}
	if strings.Compare(selectLists[0], "keys") != 0 {
		t.Errorf("Should find keys in %+v", selectLists)
		return false
	}
	sl, err := c.Select("keys")
	if err != nil {
		t.Errorf("select failed, %s", err)
		return false
	}
	keys2 := sl.Keys[:]
	if len(keys2) != 1 {
		t.Errorf("Should only have one key in collection, %+v", keys2)
		return false
	}
	if len(keys1) != len(keys2) {
		t.Errorf("select list does match collection keys")
		return false
	}
	for i, k := range keys1 {
		if strings.Compare(k, keys2[i]) != 0 {
			t.Errorf("Select list does not match key at %d, %q != %q", i, k, keys2[i])
			return false
		}
	}

	records := map[string]interface{}{
		"littlefreda": map[string]string{
			"name": "Little Freda",
		},
		"mojosam": map[string]string{
			"name": "Mojo Sam",
		},
		"captainjack": map[string]string{
			"name": "Jack Flanders",
		},
	}

	// Create some additional records to work with
	for name, rec := range records {
		err = c.Create(name, rec)
		if err != nil {
			t.Errorf("Could not create test record %s %s -> %s", name, rec, err)
			return false
		}
	}

	// Try create jack-and-mojo select list
	jackAndMojo, err := c.Select("jack-and-mojo", "captainjack", "mojosam")
	if err != nil {
		t.Errorf("create jack-and-mojo select list %s", err)
		return false
	}
	if jackAndMojo.Length() != 2 {
		t.Errorf("Expected 2, got length of %d", jackAndMojo.Length())
		return false
	}

	// Test the non destructive operations
	jack := jackAndMojo.First()
	mojo := jackAndMojo.Last()
	restOfList := jackAndMojo.Rest()

	if strings.Compare("captainjack", jack) != 0 {
		t.Errorf("First() should have returned captainjack, %s", jack)
		return false
	}
	if strings.Compare("mojosam", mojo) != 0 {
		t.Errorf("Last() should have returned mojosam, %s", mojo)
		return false
	}
	if len(restOfList) != 1 {
		t.Errorf("Rest() should return a list of 1 key, %+v", restOfList)
		return false
	}

	// Now we'll update the list by re-selecting it, should add two elements to it
	jackAndMojo, err = c.Select("jack-and-mojo", "jack", "mojo")
	if err != nil {
		t.Errorf("updating jack-and-mojo list with two more elements - jack and mojo, %+v", err)
		return false
	}
	if jackAndMojo.Length() != 4 {
		t.Errorf("Expected length of 4, got %d for %+v", jackAndMojo.Length(), jackAndMojo)
		return false
	}
	restOfList = jackAndMojo.Rest()
	if len(restOfList) != 3 {
		t.Errorf("Rest() should now return 3 elements, %+v", restOfList)
		return false
	}

	// Try out reverse, sort, shift, pop, unshift and push operations
	jackAndMojo.Reverse()
	mojo = jackAndMojo.First()
	jack = jackAndMojo.Last()
	if strings.Compare("captainjack", jack) != 0 {
		t.Errorf("Last() should have returned captainjack, %s <- %+v", jack, jackAndMojo)
		return false
	}
	if strings.Compare("mojo", mojo) != 0 {
		t.Errorf("First() should have returned mojo, %s <- %+v", mojo, jackAndMojo)
		return false
	}
	jackAndMojo.Sort(ASC)
	for _, expected := range []string{"captainjack", "jack", "mojo", "mojosam"} {
		val := jackAndMojo.Shift()
		if strings.Compare(expected, val) != 0 {
			t.Errorf("Sort() failed, %+v\n", jackAndMojo)
			return false
		}
	}
	if jackAndMojo.Length() != 0 {
		t.Errorf("Shift didn't work, %+v", jackAndMojo)
		return false
	}
	jackAndMojo, _ = c.Select("jack-and-mojo", "captainjack", "jack", "mojo", "mojosam")
	jackAndMojo.Sort(DESC)
	for _, expected := range []string{"captainjack", "jack", "mojo", "mojosam"} {
		val := jackAndMojo.Pop()
		if strings.Compare(expected, val) != 0 {
			t.Errorf("Sort() failed %q != %q, %+v\n", expected, val, jackAndMojo)
			return false
		}
	}
	jackAndMojo.Unshift("littlefreda")
	if strings.Compare("littlefreda", jackAndMojo.First()) != 0 {
		t.Errorf("Unshift failed, %+v", jackAndMojo)
		return false
	}
	jackAndMojo.Push("littlefreda")
	if strings.Compare("littlefreda", jackAndMojo.Last()) != 0 {
		t.Errorf("Push failed, %+v", jackAndMojo)
		return false
	}

	// Make sure you cannot create a "collections" select list
	// Make sure you can not change the default "keys" select list
	return true
}
