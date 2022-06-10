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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"testing"
	"time"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset/pairtree"
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
			"email":       "jsteinbeck@shiphardbor.example.org",
			"genre":       []string{"writer"},
		},
	}
)

func setupTestCollection(cName string) error {
	// Create collection.json using v1 structures
	if len(cName) == 0 {
		return fmt.Errorf("missing a collection name")
	}
	// Make root directory to hold collection.
	os.MkdirAll(cName, 0775)
	// Generate a v1 collection
	c := new(Collection)
	// Save the date and time
	dt := time.Now()
	// date and time is in RFC3339 format
	c.Created = dt.Format(time.RFC3339)
	// When is a date in YYYY-MM-DD format (can be approximate)
	// e.g. 2021, 2021-01, 2021-01-02
	c.When = dt.Format("2006-01-02")
	c.DatasetVersion = Version
	c.Name = path.Base(cName)
	c.Version = "v0.0.0"
	userinfo, err := user.Current()
	if err == nil {
		if userinfo.Name != "" {
			c.Who = []string{userinfo.Name}
		} else {
			c.Who = []string{userinfo.Username}
		}
	}
	if len(c.Who) > 0 {
		c.What = fmt.Sprintf("A dataset (%s) collection initilized on %s by %s.", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST."), strings.Join(c.Who, ", "))
	} else {
		c.What = fmt.Sprintf("A dataset %s collection initilized on %s", Version, dt.Format("Monday, January 2, 2006 at 3:04pm MST.."))
	}
	c.workPath = cName
	c.KeyMap = make(map[string]string)

	if err := c.saveMetadata(); err != nil {
		return err
	}
	// Now populate with some test records records.
	for key, obj := range records {
		src, err := json.MarshalIndent(obj, "", "    ")
		pair := pairtree.Encode(key)
		pPath := path.Join("pairtree", pair)
		c.KeyMap[key] = pPath
		filename := path.Join(cName, pPath, key+".json")
		if err := os.MkdirAll(path.Dir(filename), 0775); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filename, src, 0664); err != nil {
			return err
		}
		// Update collection.json to have keys for added records
		if err = c.saveMetadata(); err != nil {
			return err
		}
	}
	return nil
}

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
	err := setupTestCollection(cName)
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
