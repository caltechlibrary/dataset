//
// dsv1 is a submodule of dataset package.
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
package dsv1

import (
	"os"
	"path"
	"testing"
)

var (
	records = map[string]map[string]interface{}{
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
)

func toStringSlice(in []interface{}) []string {
	out := []string{}
	for _, s := range in {
		switch s.(type) {
		case string:
			v := s.(string)
			out = append(out, v)
		}
	}
	return out
}

func TestCollection(t *testing.T) {
	cName := path.Join("testout", "test_collection.ds")
	if _, err := os.Stat(cName); err == nil {
		// Remove any pre-existing test data
		os.RemoveAll(cName)
	}
	err := SetupV1TestCollection(cName, records)
	if err != nil {
		t.Errorf("Unable to setup test collection %q, %s", cName, err)
		t.FailNow()
	}

	// Now open the existing collection of colName
	c, err := Open(cName)
	if err != nil {
		t.Errorf("error Open() a collection %q", err)
		t.FailNow()
	}

	// Make sure I can read records
	keys := c.Keys()
	if len(keys) != len(records) {
		t.Errorf("expected %d keys, got %+v", len(records), keys)
		t.FailNow()
	}

	// Run through and make sure we can read all the records.
	for key, expectedObject := range records {
		gotObject := map[string]interface{}{}
		err = c.Read(key, gotObject, false)
		if err != nil {
			t.Errorf("%q: Read(), %s", c.Name, err)
			t.FailNow()
		}
		var (
			expectedS string
			expectedA []string
		)
		for field, val := range expectedObject {
			switch val.(type) {
			case string:
				expectedS = val.(string)
			case []string:
				expectedA = val.([]string)
			case []interface{}:
				expectedA = toStringSlice(val.([]interface{}))
			}

			if got, ok := gotObject[field]; ok {
				switch got.(type) {
				case string:
					gotS := got.(string)
					if expectedS != gotS {
						t.Errorf("expected %q, got %q for %q -> %q", expectedS, gotS, key, field)
					}
				case []interface{}:
					gotA := toStringSlice(got.([]interface{}))
					if len(expectedA) != len(gotA) {
						t.Errorf("expected %+v, got %+v", expectedA, gotA)
					} else {
						for i, es := range expectedA {
							if es != gotA[i] {
								t.Errorf("expected %q, got %q for %q -> %q", es, gotA[i], key, field)
							}
						}
					}
				case []string:
					gotA := got.([]string)
					if len(expectedA) != len(gotA) {
						t.Errorf("expected %+v, got %+v", expectedA, gotA)
					} else {
						for i, es := range expectedA {
							if es != gotA[i] {
								t.Errorf("expected %q, got %q for %q -> %q", es, gotA[i], key, field)
							}
						}
					}
				default:
					t.Errorf("Failed to compare %q type %T", field, got)
				}
			} else {
				t.Errorf("Missing %q from %+v", field, gotObject)
			}
		}
	}
}
