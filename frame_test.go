//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2019, Caltech
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
