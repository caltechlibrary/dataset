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
	"strings"
	"testing"
)

func TestFrame(t *testing.T) {
	verbose := false
	os.RemoveAll(path.Join("testdata", "frame_test.ds"))
	cName := path.Join("testdata", "frame_test.ds")
	c, err := InitCollection(cName)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer c.Close()

	//NOTE: test data and to load into collection and generate grid
	tRecords := []map[string]interface{}{
		map[string]interface{}{
			"_Key":  "A",
			"id":    "A",
			"one":   "one",
			"two":   22,
			"three": 3.0,
			"four":  []string{"one", "two", "three"},
		},
		map[string]interface{}{
			"_Key":  "B",
			"id":    "B",
			"two":   2000,
			"three": 3000.1,
		},
		map[string]interface{}{
			"_Key": "C",
			"id":   "C",
		},
		map[string]interface{}{
			"_Key":  "D",
			"id":    "D",
			"one":   "ONE",
			"two":   20,
			"three": 334.1,
			"four":  []string{},
		},
	}
	keys := []string{}
	for _, rec := range tRecords {
		key := rec["_Key"].(string)
		keys = append(keys, key)
		src, _ := EncodeJSON(rec)
		err := c.CreateJSON(key, src)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	f, err := c.FrameCreate("frame-1", keys, []string{".id", ".one", ".two", ".three", ".four"}, []string{"id", "one", "two", "three", "four"}, verbose)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(f.ObjectMap) == 0 {
		t.Errorf("Expect > 0 objects in ObjectMap")
		t.FailNow()
	}
	if len(f.ObjectMap) != len(tRecords) {
		t.Errorf("Expected tRecords (%d) to be same length as objectList (%d) -> %s", len(tRecords), len(f.ObjectMap), f.String())
		t.FailNow()
	}
	if len(f.Keys) != len(tRecords) {
		t.Errorf("Expected tRecords (%d) to be same length as keys (%d) -> %s", len(tRecords), len(f.Keys), f.String())
		t.FailNow()
	}
	expected := "frame-1"
	result := f.Name
	if expected != result {
		t.Errorf("expected %q, got %q, for %s", expected, result, f)
	}
	expected = c.Name // e.g. "frame_test.ds from  "testdata/frame_test.ds"
	result = f.CollectionName
	if expected != result {
		t.Errorf("expected %q, got %q, for %s", expected, result, f)
	}
	//FIXME: need some tests on frame structure.
	objectList := f.Objects()
	for i, obj := range objectList {
		if len(obj) == 0 {
			t.Errorf("object in object list (%d) should have content, %+v\n", i, objectList)
			t.FailNow()
		}
		rec := tRecords[i]
		for j, key := range f.Labels {
			if val, ok := obj[key]; ok != true {
				if _, ok := rec[key]; ok == true {
					t.Errorf("(%d, %d) missing key %q in obj, %+v\n", i, j, key, obj)
				}
			} else if expected, ok := rec[key]; ok != true {
				t.Errorf("(%d, %d) missing key %q in record, %+v\n", i, j, key, expected)
			} else {
				switch val.(type) {
				case string:
					if strings.Compare(val.(string), expected.(string)) != 0 {
						t.Errorf("(%d, %d, %s) expected %q, got %q", i, j, key, expected, val)
					}
				case int:
					if val.(int) != expected.(int) {
						t.Errorf("(%d, %d, %s) expected int %d, got %d", i, j, key, expected, val)
					}
				case int64:
					if val.(int64) != expected.(int64) {
						t.Errorf("(%d, %d, %s) expected int %d, got %d", i, j, key, expected, val)
					}
				case float64:
					if val.(float64) != expected.(float64) {
						t.Errorf("(%d, %d, %s) expected int %f, got %f", i, j, key, expected, val)
					}
				case json.Number:
					n1 := val.(json.Number).String()
					n2 := ""
					switch expected.(type) {
					case int:
						n2 = fmt.Sprintf("%d", expected)
					case int64:
						n2 = fmt.Sprintf("%d", expected)
					case float64:
						n2 = fmt.Sprintf("%1.1f", expected)
						// Handle the case that json.Number returns float
						// as int for valued stored, e.g. 3.0 returned as 3
						if len(n1) < len(n2) {
							n2 = n2[0:len(n1)]
						}
					}
					if strings.Compare(n1, n2) != 0 {
						t.Errorf("(%d, %d, %s) expected %s, got %s", i, j, key, n2, n1)
					}
				case []interface{}:
					e := len(expected.([]string))
					v := len(val.([]interface{}))
					if e != v {
						t.Errorf("(%d, %d, %s) expected length %d, got %d", i, j, key, e, v)
					}
				default:
					t.Errorf("(%d, %d, %s) something didn't match, expected (%T) %+v, got (%T) %+v", i, j, key, expected, expected, val, val)
					t.FailNow()
				}
			}
		}
	}
}

func TestIssue9PyDataset(t *testing.T) {
	verbose := false
	os.RemoveAll(path.Join("testdata", "frame_test2.ds"))
	cName := path.Join("testdata", "frame_test2.ds")
	c, err := InitCollection(cName)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer c.Close()

	jsonSrc := []byte(`[
        { "id":    "A", "nameIdentifiers": [
                {
                    "nameIdentifier": "0000-000X-XXXX-XXXX",
                    "nameIdentifierScheme": "ORCID",
                    "schemeURI": "http://orcid.org/"
                },
                {
                    "nameIdentifier": "H-XXXX-XXXX",
                    "nameIdentifierScheme": "ResearcherID",
                    "schemeURI": "http://www.researcherid.com/rid/"
                }], "two":   22, "three": 3.0, "four":  ["one", "two", "three"] 
},
        { "id":    "B", "two":   2000, "three": 3000.1 },
        { "id": "C" },
        { "id":    "D", "nameIdentifiers": [
                {
                    "nameIdentifier": "0000-000X-XXXX-XXXX",
                    "nameIdentifierScheme": "ORCID",
                    "schemeURI": "http://orcid.org/"
                }], "two":   20, "three": 334.1, "four":  [] }
    ]`)
	listObjects := []map[string]interface{}{}
	// FIXME: setup a custom marshaller so numbers are json.Number()
	err = json.Unmarshal(jsonSrc, &listObjects)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	for i, obj := range listObjects {
		if id, ok := obj["id"]; ok == true {
			key := id.(string)
			src, _ := EncodeJSON(obj)
			if err := c.CreateJSON(key, src); err != nil {
				t.Errorf("(%d) key %s, error: %s", i, key, err)
				t.FailNow()
			}
		}
	}
	// Now let's see if our frame works ...
	keys := c.Keys()
	f, err := c.FrameCreate("f1", keys,
		[]string{
			"._Key",
			".nameIdentifiers",
			".nameIdentifiers[:].nameIdentifier",
			".two",
			".three",
			".four",
		}, []string{
			"id",
			"nameIdentifiers",
			"nameIdentifier",
			"two",
			"three",
			"four",
		}, verbose)
	if err != nil {
		t.Errorf("Can't make frame f1, %s", err)
		t.FailNow()
	}
	if f == nil {
		t.Errorf("Expected a frame named f1, got nil")
	}
}

func TestFrameRefresh(t *testing.T) {
	cName := path.Join("testdata", "frame10.ds")
	os.RemoveAll(cName)
	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("expected to create %q, got %s", cName, err)
		t.FailNow()
	}
	key := "k1"
	src := []byte(`{
		"title": "Orchids & Moonbeams",
		"cast": [
			{ 
				"last_name": "Lorick",
				"first_name": "Robert",
				"character": "Jack Flanders"
			},
			{
				"last_name": "Adams",
				"first_name": "Dave",
				"character": "Mojo Sam"
			},
			{
				"last_name": "Poirier",
				"first_name": "Pascale",
				"character": "Claudine"
			},
			{
				"last_name": "Donovan",
				"first_name": "Patrick",
				"character": "Pat Patternson"
			},
			{
				"last_name": "Goodhart Hebert",
				"first_name": "Camille",
				"character": "Bunny"
			},
			{
				"last_name": "Roth",
				"first_name": "Laura",
				"character": "Amber"
			}
		]
		}`)
	if err := c.CreateJSON(key, src); err != nil {
		t.Errorf("expected to create %q, got %s", key, err)
		t.FailNow()
	}
	src = []byte(`{
		"title": "The incredible Adventures of Jack Flanders",
		"cast": [
			{ 
				"last_name": "Lorick",
				"first_name": "Robert",	
				"character": "Jack Flan"
			},
			{
				"last_name": "Adams",
				"first_name": "Dave",
				"character": "Mojo Sam"
			},
			{
				"last_name": "Orte",
				"first_name": "P. J.",
				"character": "Little Freda"
			}
		]
	}`)
	key = "k0"
	if err := c.CreateJSON(key, src); err != nil {
		t.Errorf("expected to create %q, got %s", key, err)
		t.FailNow()
	}
	for _, key := range []string{"k0", "k1"} {
		if _, err := c.ReadJSON(key); err != nil {
			t.Errorf("expected %q, got error %s", key, err)
			t.FailNow()
		}
	}

	fName := "f1"
	verbose := false
	dotPaths := []string{".title", ".cast"}
	labels := []string{"title", "cast"}
	keys := []string{"k1"}
	f, err := c.FrameCreate(fName, keys, dotPaths, labels, verbose)
	if err != nil {
		t.Errorf("expected to create frame %q, got %s", key, err)
		t.FailNow()
	}
	ol := f.Objects()
	if len(ol) != 1 {
		t.Errorf("expected one object, got %d", len(ol))
		t.FailNow()
	}
	if c.FrameExists(fName) == false {
		t.Errorf("expected %q, none was found", fName)
		t.FailNow()
	}
	if err := c.FrameRefresh(fName, verbose); err != nil {
		t.Errorf("expected successful refresh %q, got %s", fName, err)
		t.FailNow()
	}
	ol2, err := c.FrameObjects(fName)
	if err != nil {
		t.Errorf("expected object list, got error %s", err)
		t.FailNow()
	}
	if len(ol2) != 1 {
		t.Errorf("expected 1 object, got %d -> %+v", len(ol2), ol2)
		t.FailNow()
	}
}

func TestIssue12PyDataset(t *testing.T) {
	cName := path.Join("testdata", "test_issue12.ds")
	os.RemoveAll(cName)
	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("failed to create %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()
	// Build some test data ...
	keys := []string{"1", "2", "3", "4", "5"}
	for i, key := range keys {
		obj := map[string]interface{}{}
		src := []byte(fmt.Sprintf(`{"id": "%d", "c1": %d, "c2": %d, "c3": %d}`, i, (i + 1), (i + 3), (i + 5)))
		if err := json.Unmarshal(src, &obj); err != nil {
			t.Errorf("failed to marshal %s, %s", src, err)
			t.FailNow()
		}
		if c.KeyExists(key) == true {
			if err = c.Update(key, obj); err != nil {
				t.Errorf("failed to update %q in %q, %s", key, cName, err)
			}
		} else {
			if err = c.Create(key, obj); err != nil {
				t.Errorf("failed to create %q in %q, %s", key, cName, err)
			}
		}
	}
	// Clear out any stale frames.
	for i, fName := range c.Frames() {
		if err := c.FrameDelete(fName); err != nil {
			t.Errorf("Failed to delete frame (%d) %q in %q, %s", i, fName, cName, err)
			t.FailNow()
		}
	}
	fName := "issue12"
	dotPaths := []string{".c1", ".c3"}
	labels := []string{".col1", ".col3"}
	verbose := true
	f, err := c.FrameCreate(fName, keys, dotPaths, labels, verbose)
	if err != nil {
		t.Errorf("FrameCreate failed, %q in %q, %s", fName, cName, err)
		t.FailNow()
	}
	if len(f.Keys) != len(keys) {
		t.Errorf("expected %d keys, got %d", len(keys), len(f.Keys))
		t.FailNow()
	}
	fObjects := f.Objects()
	if len(fObjects) != len(keys) {
		t.Errorf("expected %d objects, got %d", len(keys), len(fObjects))
		t.FailNow()
	}
	if err := c.FrameDelete(fName); err != nil {
		t.Errorf("expected no errors for delete frame, got %s", err)
		t.FailNow()
	}
}
