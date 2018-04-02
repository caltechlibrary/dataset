//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
	"strings"
	"testing"
)

func TestBucketNames(t *testing.T) {
	buckets := generateBucketNames(DefaultAlphabet, 3)
	for _, val := range buckets {
		if len(val) != 3 {
			t.Errorf("Should have a name of length 3. %q", val)
		}
	}
}

func TestPickBucketName(t *testing.T) {
	alphabet := "ab"
	buckets := generateBucketNames(alphabet, 2)
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
	buckets := generateBucketNames(alphabet, 2)
	if len(buckets) != 4 {
		t.Errorf("Should have four buckets %+v", buckets)
		t.FailNow()
	}

	// Remove any pre-existing test data
	os.RemoveAll(colName)

	// Create a new collection
	collection, err := create(colName, buckets)
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
			err = collection.Create(id, rec)
			if err != nil {
				t.Errorf("collection.Create(), %s", err)
				t.FailNow()
			}
			p, err := collection.DocPath(id)
			if err != nil {
				t.Errorf("Should have docpath for %s, %s", id, err)
				t.FailNow()
			}
			if _, err := os.Stat(p); os.IsNotExist(err) == true {
				t.Errorf("Should have saved %s to disc at %s", id, p)
				t.FailNow()
			}
		}
	}

	if len(collection.KeyMap) != 3 {
		t.Errorf("expected 1 key, got %+v", collection)
		t.FailNow()
	}
	keys := collection.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %+v", keys)
		t.FailNow()
	}

	// Create an empty record, then read it again to compare
	keyName := "Kahlo-F"
	rec2 := map[string]interface{}{}
	err = collection.Read(keyName, rec2)
	if err != nil {
		t.Errorf("Read(), %s", err)
		t.FailNow()
	}
	rec1 := testData[0]
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
	err = collection.Create(keyName, rec2)
	if err == nil {
		t.Errorf("Create not allow creationg on an existing record, %s --> %+v", keyName, rec2)
		t.FailNow()
	}

	rec3 := map[string]interface{}{}
	if err := collection.Read(keyName, rec3); err != nil {
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

	rec2["email"] = "freda@collectivo.example.org"
	err = collection.Update(keyName, rec2)
	if err != nil {
		t.Errorf("Could not update %s, %s", "freda", err)
		t.FailNow()
	}

	rec4 := map[string]interface{}{}
	if err := collection.Read(keyName, rec4); err != nil {
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

	err = collection.Delete(keyName)
	if err != nil {
		t.Errorf("Should be able to delete %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = collection.Read(keyName, rec2)
	if err == nil {
		t.Errorf("Record should have been deleted, %+v, %s", rec2, err)
	}

	err = Delete(colName)
	if err != nil {
		t.Errorf("Couldn't remove collection %s, %s", colName, err)
	}
}

func TestExtract(t *testing.T) {
	colName := "testdata/test_extract.ds"
	os.RemoveAll(colName)
	c, err := InitCollection(colName)
	if err != nil {
		t.Errorf("failed to create %s, %s", colName, err)
		t.FailNow()
	}
	if err := c.CreateJSON("gutenberg:21489", []byte(`{"title": "The Secret of the Island", "formats": ["epub","kindle", "plain text", "html"], "authors": [{"given": "Jules", "family": "Verne"}], "url": "http://www.gutenberg.org/ebooks/21489"}`)); err != nil {
		t.Errorf("Can't add record %s", err)
		t.FailNow()
	}
	if err := c.CreateJSON("gutenberg:2488", []byte(`{"title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}`)); err != nil {
		t.Errorf("Can't add record %s", err)
		t.FailNow()
	}
	if err := c.CreateJSON("gutenberg:21839", []byte(`{ "title": "Sense and Sensibility", "formats": ["epub", "kindle", "plain text"], "authors": [{"given": "Jane", "family": "Austin"}], "url": "http://www.gutenberg.org/ebooks/21839" }`)); err != nil {
		t.Errorf("Can't add record %s", err)
		t.FailNow()
	}
	if err := c.CreateJSON("gutenberg:3186", []byte(`{"title": "The Mysterious Stranger, and Other Stories", "formats": ["epub","kindle", "plain text", "html"], "authors": [{ "given": "Mark", "family": "Twain"}], "url": "http://www.gutenberg.org/ebooks/3186"}`)); err != nil {
		t.Errorf("Can't add record %s", err)
		t.FailNow()
	}
	if err := c.CreateJSON("hathi:uc1321060001561131", []byte(`{ "title": "A year of American travel - Narrative of personal experience", "formats": ["pdf"], "authors": [{"given": "Jessie Benton", "family": "Fremont"}], "url": "https://babel.hathitrust.org/cgi/pt?id=uc1.32106000561131;view=1up;seq=9" }`)); err != nil {
		t.Errorf("Can't add record %s", err)
		t.FailNow()
	}
	if i := c.Length(); i != 5 {
		t.Errorf("Expected 5 records, got %d", i)
		t.FailNow()
	}
	keys := c.Keys()
	if len(keys) != 5 {
		t.Errorf("Expected 5 keys, got %+v\n", keys)
		t.FailNow()
	}

	l, err := c.Extract("true", ".authors[:].family")
	if err != nil {
		t.Errorf("Can't extract values, %s", err)
	}
	if len(l) != 4 {
		t.Errorf("expected four author last names, %+v", l)
	}
	for _, target := range []string{"Verne", "Austin", "Twain", "Fremont"} {
		found := false
		for _, v := range l {
			if v == target {
				found = true
				break
			}
		}
		if found == false {
			t.Errorf("Could not find %s in list %+v", target, l)
		}
	}

	src := []byte(`[
{"authorAffiliation":"California Institute of Technology","authorName":"Roberts, Ellis Earl"},
{"authorAffiliation":["Ariane Tracking Station, Ascension Island (SH)"],"authorName":"John, N."},
{"authorAffiliation":["California Institute of Technology, Pasadena, CA (US)"],"authorName":"Yavin, Y."},
{"authorAffiliation":["California Institute of Technology, Pasadena, CA, USA","University of Toronto, Toronto, ON, CAN"],"authorIdentifiers":[{"authorIdentifier":"0000-0003-2025-7519","authorIdentifierScheme":"ORCID"}],"authorName":"Jacob Hedelius"},
{"authorAffiliation":["California Institute of Technology, Pasadena, CA, USA"],"authorIdentifiers":[{"authorIdentifier":"0000-0002-6126-3854","authorIdentifierScheme":"ORCID"},{"authorIdentifier":"A-5460-2012","authorIdentifierScheme":"ResearcherID"}],"authorName":"Paul Wennberg"},
{"authorAffiliation":["Caltech Library"],"authorIdentifiers":[{"authorIdentifier":"0000-0003-0900-6903","authorIdentifierScheme":"ORCID"}],"authorName":"Doiel,Robert"},
{"authorAffiliation":["Caltech"],"authorName":"Doiel, Robert"},
{"authorAffiliation":["Wisconsin Educational Communications Board, Park Falls, WI (US)"],"authorName":"Ayers, J."},
{"authorName":"Neufeld, G."},
{"authorName":"Springett, S."},
{"authorName":"Yavin, Y."}
]`)
	data := []map[string]interface{}{}
	err = json.Unmarshal(src, &data)
	if err != nil {
		t.Errorf("Can't unmarshal test data, %s", err)
		t.FailNow()
	}
	for i, rec := range data {
		if err := c.Create(fmt.Sprintf("%d", i), rec); err != nil {
			t.Errorf("Can't create %d in %s, %s", i, c.Name, err)
		}

	}
	lines, err := c.Extract("true", ".authorAffiliation[:]")
	if err != nil {
		t.Errorf("Can't extract .authorAffiliation[:] from %s, %s", c.Name, err)
	}
	for i, line := range lines {
		fmt.Printf("DEBUG line (%d): %q\n", i, line)
		if strings.HasPrefix(line, "[") == true {
			t.Errorf("%d started as an array, %s, expecting simple string", i, line)
		}
	}
}

func TestComplexKeys(t *testing.T) {
	colName := "testdata/col2.ds"
	buckets := generateBucketNames("ab", 2)
	if len(buckets) != 4 {
		t.Errorf("Should have four buckets %+v", buckets)
		t.FailNow()
	}

	// remove any stale test collection collection first...
	os.RemoveAll(colName)

	// Create a new collection
	collection, err := create(colName, buckets)
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
		err := collection.Create(k, v)
		if err != nil {
			t.Errorf("Can't create %s <-- %s : %s", k, v, err)
		}
	}
}

func TestExtractOverVariableSchema(t *testing.T) {
	src := []byte(`[
{"name": "one", "test_path": 1},
{"name": "two", "field": "B"},
{"name": "three", "field": "C", "test_path": 2, "and_the_other_thing":true}
]`)

	data := []map[string]interface{}{}
	err := json.Unmarshal(src, &data)
	if err != nil {
		t.Errorf("Can't generate tests records, %s", err)
		t.FailNow()
	}
	cName := path.Join("testdata", "extract_test.ds")
	os.RemoveAll(cName)
	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("Can't create %s, %s", cName, err)
		t.FailNow()
	}
	for i, rec := range data {
		err = c.Create(fmt.Sprintf("%d", i), rec)
		if err != nil {
			t.Errorf("Can't create record %d in %s, %s", i, c.Name, err)
		}
	}
	result, err := c.Extract("true", ".field")
	if err != nil {
		t.Errorf("Expected no error for .field in %s, got %s", c.Name, err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 values, got %+v", result)
	}
	result, err = c.Extract("true", ".field_does_not_exist")
	if err != nil {
		t.Errorf("Expected no error for .field_does_not_exist in %s, got %s", c.Name, err)
	}
	if len(result) > 0 {
		t.Errorf("Expected an empty result for .field_does_not_exist, got %+v", result)
	}
	result, err = c.Extract("true", ".name")
	if len(result) != 3 {
		t.Errorf("Expected three values for .name, got %+v", result)
	}
	if err != nil {
		t.Errorf("Expected no errors for .name in %s, got, %s", c.Name, err)
	}
}
