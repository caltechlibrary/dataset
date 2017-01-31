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
	"sort"
	"strings"
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
	if selectLists[0] != "keys" {
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
		if k != keys2[i] {
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
	if jackAndMojo.Len() != 2 {
		t.Errorf("Expected 2, got length of %d", jackAndMojo.Len())
		return false
	}

	// Test the non destructive operations
	jack := jackAndMojo.First()
	mojo := jackAndMojo.Last()
	restOfList := jackAndMojo.Rest()

	if "captainjack" != jack {
		t.Errorf("First() should have returned captainjack, %s", jack)
		return false
	}
	if "mojosam" != mojo {
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
	if jackAndMojo.Len() != 4 {
		t.Errorf("Expected length of 4, got %d for %+v", jackAndMojo.Len(), jackAndMojo)
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
	if "captainjack" != jack {
		t.Errorf("Last() should have returned captainjack, %s <- %+v", jack, jackAndMojo)
		return false
	}
	if "mojo" != mojo {
		t.Errorf("First() should have returned mojo, %s <- %+v", mojo, jackAndMojo)
		return false
	}
	jackAndMojo.Sort(ASC)
	for _, expected := range []string{"captainjack", "jack", "mojo", "mojosam"} {
		val := jackAndMojo.Shift()
		if expected != val {
			t.Errorf("Sort() failed, %+v\n", jackAndMojo)
			return false
		}
	}
	if jackAndMojo.Len() != 0 {
		t.Errorf("Shift didn't work, %+v", jackAndMojo)
		return false
	}
	jackAndMojo, _ = c.Select("jack-and-mojo", "captainjack", "jack", "mojo", "mojosam")
	jackAndMojo.Sort(DESC)
	for _, expected := range []string{"captainjack", "jack", "mojo", "mojosam"} {
		val := jackAndMojo.Pop()
		if expected != val {
			t.Errorf("Sort() failed %q != %q, %+v\n", expected, val, jackAndMojo)
			return false
		}
	}
	jackAndMojo.Unshift("littlefreda")
	if "littlefreda" != jackAndMojo.First() {
		t.Errorf("Unshift failed, %+v", jackAndMojo)
		return false
	}
	jackAndMojo.Push("littlefreda")
	if "littlefreda" != jackAndMojo.Last() {
		t.Errorf("Push failed, %+v", jackAndMojo)
		return false
	}

	// Make sure you cannot create a "collections" select list
	// Make sure you can not change the default "keys" select list
	return true
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

func TestSelectListSort(t *testing.T) {
	colName := "testdata/complex-sorting"
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

	testKeyList := []string{
		"A|2017-01-01|0",
		"B|2016-01-01|1",
		"C|2017-01-01|2",
		"D|2016-01-01|3",
		"A|2014-01-01|4",
		"A|2020-01-01|5",
		"B|1918-01-01|6",
		"B|1920-01-01|7",
		"C|2021-06-08|8",
	}
	collection.Clear("sorttests")
	sl, err := collection.Select(append([]string{"sorttests"}, testKeyList[:]...)...)
	if err != nil {
		t.Errorf("Cannot create select list simple: %s", err)
		t.FailNow()
	}

	// Setup simple sort expected results
	expectedSimpleSort := testKeyList[:]
	sort.Sort(sort.StringSlice(expectedSimpleSort))
	// Run simple sort of select list
	sl.Sort(ASC)

	// Compare results
	result := ""
	for i, expected := range expectedSimpleSort {
		result = sl.Keys[i]
		if expected != result {
			t.Errorf("for ith: %d, expected %s, got %s", i, expected, result)
		}
	}
	sl.CustomLessFn = func(s []string, i, j int) bool {
		k1, k2 := strings.Split(s[i], "|"), strings.Split(s[j], "|")
		// Compare each element of each key and sort zero-th element ascending, and first element descending
		if k1[0] <= k2[0] && k1[1] >= k2[1] {
			return true
		}
		return false
	}
	sl.Sort(ASC)
	expectedComplexSort := []string{
		"A|2020-01-01|5",
		"A|2017-01-01|0",
		"A|2014-01-01|4",
		"B|2016-01-01|1",
		"B|1920-01-01|7",
		"B|1918-01-01|6",
		"C|2021-06-08|8",
		"C|2017-01-01|2",
		"D|2016-01-01|3",
	}
	result = ""
	for i, expected := range expectedComplexSort {
		result = sl.Keys[i]
		if expected != result {
			t.Errorf("for ith: %d, expected %q, got %q\n", i, expected, result)
		}
	}

	sl.CustomLessFn = nil
	sl.Sort(ASC)
	result = ""
	for i, expected := range expectedSimpleSort {
		result = sl.Keys[i]
		if expected != result {
			t.Errorf("for ith: %d, expected %q, got %q\n", i, expected, result)
		}
	}

	test3PartKeys := []string{
		"0000-0001-5245-0538|2017-01-19|73753",
		"0000-0001-5245-0538|2017-01-18|73721",
		"0000-0001-5245-0538|2000-07-15|73689",
		"0000-0001-5245-0538|2000-05-01|73688",
		"0000-0001-5245-0538|2000-02-15|73679",
		"0000-0001-5245-0538|2004-09-01|73677",
	}
	expected3PartKeys := []string{
		"0000-0001-5245-0538|2017-01-19|73753",
		"0000-0001-5245-0538|2017-01-18|73721",
		"0000-0001-5245-0538|2004-09-01|73677",
		"0000-0001-5245-0538|2000-07-15|73689",
		"0000-0001-5245-0538|2000-05-01|73688",
		"0000-0001-5245-0538|2000-02-15|73679",
	}
	sl.Reset()
	sl.Keys = test3PartKeys[:]
	sl.SaveList()
	sl.CustomLessFn = func(s []string, i, j int) bool {
		k1, k2 := strings.Split(s[i], "|"), strings.Split(s[j], "|")
		if k1[0] <= k2[0] && k1[1] > k2[1] {
			return true
		}
		return false
	}
	sl.Sort(ASC)
	result = ""
	for i, expected := range expected3PartKeys {
		result = sl.Keys[i]
		if expected != result {
			t.Errorf("for ith: %d, expected %q, got %q\n", i, expected, result)
		}
	}
}
