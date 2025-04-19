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
package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// Test the Pairtree storage implementation
func TestPTStore(t *testing.T) {
	threeObjects := []map[string]interface{}{}
	threeObjects = append(threeObjects, map[string]interface{}{"one": 1})
	threeObjects = append(threeObjects, map[string]interface{}{"two": 2})
	threeObjects = append(threeObjects, map[string]interface{}{"three": 3})

	cName := path.Join("testout", "T1.ds")
	// Clear stale test output
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	// NOTE: Pairtree doesn't use an DSN URI so it is empty
	c, err := Init(cName, "pairtree")
	if err != nil {
		t.Errorf("Failed to create %q, %s", cName, err)
		t.FailNow()
	}
	if c == nil {
		t.Errorf("Failed to create a collection object")
		t.FailNow()
	}
	if err := c.Close(); err != nil {
		t.Errorf("Failed to close collection %q, %s", cName, err)
		t.FailNow()
	}
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Failed to open a collection object")
		t.FailNow()
	}
	for i, obj1 := range threeObjects {
		key := fmt.Sprintf("%d", i)
		// See if I can create and read back the object
		if err := c.Create(key, obj1); err != nil {
			t.Errorf("Expected to create %q, %s", key, err)
			t.FailNow()
		} else {
			obj2 := map[string]interface{}{}
			if err := c.Read(key, obj2); err != nil {
				t.Errorf("Expected to read %q, %s", key, err)
				t.FailNow()
			}
			for k, v := range obj1 {
				if v2, ok := obj2[k]; ok == false {
					t.Errorf("Expected %q in obj2 %+v", k, obj2)
				} else {
					// NOTE: c.Read() will use json.Number types for
					// integers and floats expressed in JSON. These
					// needed to be converted appropriately for comparison.
					x := fmt.Sprintf("%d", v)
					y := v2.(json.Number).String()
					if x != y {
						t.Errorf("Expected first value (%T) %q to equal second value (%T) %q", v, x, v2, y)
					}
				}
			}
		}
		// Make sure you can't overwrite a previously created object
		if err := c.Create(key, obj1); err == nil {
			t.Errorf("Expected Create to fail %q should already exist", key)
		}
		obj1["id"] = key
		if err := c.Update(key, obj1); err != nil {
			t.Errorf("Expected Update to succeed %q, %s", key, err)
		} else {
			obj2 := map[string]interface{}{}
			if err := c.Read(key, obj2); err != nil {
				t.Errorf("Expected update then Read to success, %q, %s", key, err)
			} else {
				for k, v := range obj1 {
					if v2, ok := obj2[k]; ok == false {
						t.Errorf("Expected %q in obj2 %+v", k, obj2)
					} else {
						switch v.(type) {
						case string:
							if v != v2 {
								t.Errorf("Expected first value (%T) %q to equal second value (%T) %q", v, v, v2, v)
							}
						case int:
							x := fmt.Sprintf("%d", v)
							y := v2.(json.Number).String()
							if x != y {
								t.Errorf("Expected first value (%T) %q to equal second value (%T) %q", v, x, v2, y)
							}
						default:
							t.Errorf("value was unexpected type %T -> %+v", v, v)
						}
					}
				}
			}
		}
	}
	if keys, err := c.Keys(); err != nil {
		t.Errorf("Expected to get a list of keys got error %s", err)
		t.FailNow()
	} else {
		if len(keys) != 3 {
			t.Errorf("Expected three keys for 3 objects got %+v", keys)
		}
		sort.Strings(keys)
		for i := 0; i < 3; i++ {
			expected := fmt.Sprintf("%d", i)
			if keys[i] != expected {
				t.Errorf("Expected key %s, got %s", expected, keys[i])
			}
		}
	}
	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("%d", i)
		if err := c.Delete(key); err != nil {
			t.Errorf("Expected to be able to delete %q, %s", key, err)
		}
	}
	if keys, err := c.Keys(); err != nil {
		t.Errorf("Expected to get keys back from List, %s", err)
	} else if len(keys) != 1 {
		t.Errorf("Expected one key left after delete, got %+v", keys)
	}
	start := "2001-01-01 00:00:00"
	end := time.Now().Format("2006-01-02") + " 23:23:59"
	if _, err := c.UpdatedKeys(start, end); err == nil {
		t.Errorf("expected error for unsupported UpdateKeys on pairtree collections")
	}
}

// Test the SQL storage implementation
func TestSQLStore(t *testing.T) {
	threeObjects := []map[string]interface{}{}
	threeObjects = append(threeObjects, map[string]interface{}{"one": 1})
	threeObjects = append(threeObjects, map[string]interface{}{"two": 2})
	threeObjects = append(threeObjects, map[string]interface{}{"three": 3})

	cName := path.Join("testout", "SQL1.ds")
	dsnURI := "sqlite://" + path.Join(cName, "sql1.db")
	// Clear stale test output
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	// NOTE: SQLStore requires a DSN URI so it is NOT empty
	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Failed to create %q, %s", cName, err)
	}
	if c == nil {
		t.Errorf("Failed to create a collection object")
		t.FailNow()
	}
	defer func() {
		if err := c.Close(); err != nil {
			t.Errorf("Failed to close collection %q, %s", cName, err)
			t.FailNow()
		}
	}()
	for i, obj1 := range threeObjects {
		key := fmt.Sprintf("%d", i)
		obj1["id"] = key
		// See if I can create and read back the object
		if err := c.Create(key, obj1); err != nil {
			t.Errorf("Expected to create %q, %s", key, err)
			t.FailNow()
		}
		// See if we can read the object just created
		obj2 := map[string]interface{}{}
		if err := c.Read(key, obj2); err != nil {
			t.Errorf("Expected to read %q, %s", key, err)
			t.FailNow()
		}
		for k, v := range obj1 {
			if v2, ok := obj2[k]; ok == false {
				t.Errorf("Expected %q in obj2 %+v", k, obj2)
			} else {
				x, y := "", ""
				switch v.(type) {
				case int:
					x = fmt.Sprintf("%d", v)
				case string:
					x = v.(string)
				}
				switch v2.(type) {
				case json.Number:
					y = v2.(json.Number).String()
				case string:
					y = v2.(string)
				}
				if x != y {
					t.Errorf("Expected (%T) %+v, got (%T) %+v", v, v, v2, v2)
				}
			}
		}
		// Make sure you can't overwrite a previously created object
		if err := c.Create(key, obj1); err == nil {
			t.Errorf("Expected Create to fail %q should already exist", key)
		}
		// Update the object
		obj1["greeting"] = "Hi There"
		if err := c.Update(key, obj1); err != nil {
			t.Errorf("Expected Update to succeed %q, %s", key, err)
		}
		// Now read back the updated record and check
		if err := c.Read(key, obj2); err != nil {
			t.Errorf("Expected update then Read to success, %q, %s", key, err)
		}
		// Check the attributes of the records
		for k, v := range obj1 {
			if v2, ok := obj2[k]; ok == false {
				t.Errorf("Expected %q in obj2 %+v", k, obj2)
			} else {
				x, y := "", ""
				switch v.(type) {
				case int:
					x = fmt.Sprintf("%d", v)
				case string:
					x = v.(string)
				}
				switch v2.(type) {
				case json.Number:
					y = v2.(json.Number).String()
				case string:
					y = v2.(string)
				}
				if x != y {
					t.Errorf("Expected (%T) %+v, got (%T) %+v", v, v, v2, v2)
				}
			}
		}
	}
	// Should have three records in collection.
	cnt := c.Length()
	if cnt != 3 {
		t.Errorf("Expect 3 records, got %d", cnt)
	}
	// Now check the keys stored
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("Expected to get a list of keys got error %s", err)
		t.FailNow()
	}
	k := int(cnt)
	l := len(keys)
	if k != l {
		t.Errorf("Expected three keys for (%T) %d objects got (%T) %d -> %+v", k, k, l, l, keys)
	}
	sort.Strings(keys)
	for i := 0; i < 3 && i < len(keys); i++ {
		expected := fmt.Sprintf("%d", i)
		got := keys[i]
		if expected != got {
			t.Errorf("Expected key (%T) %s, got (%T) %s", expected, expected, got, got)
		}
	}
	now := time.Now()
	start := "2001-01-01 00:00:00"
	end := now.Add(12*time.Hour).Format("2006-01-02") + " 23:23:59"
	updatedKeys, err := c.UpdatedKeys(start, end)
	if err != nil {
		t.Errorf("expected UpdateKeys to work, %s", err)
	}
	sort.Strings(updatedKeys)
	ul := len(updatedKeys)
	if ul == l {
		for i, k := range keys {
			if k != updatedKeys[i] {
				t.Errorf("expected %d-th key to be %q, got %q", i, k, updatedKeys[i])
			}
		}
	} else {
		t.Errorf("expected %d updated keys, got %d in %q", l, ul, cName)
	}

	// Let's set c.Query()
	sqlStmt := fmt.Sprintf(`select src
from %s
order by updated`, strings.TrimSuffix(path.Base(cName), ".ds"))
	src, err := c.QueryJSON(sqlStmt, false)
	if err != nil {
		t.Errorf(`ran c.Query(%q) did not expect error, %s`, sqlStmt, err)
		t.FailNow()
	}

	// Test deletes
	/*
		for i := 0; i < 2; i++ {
			key := fmt.Sprintf("%d", i)
			if err := c.Delete(key); err != nil {
				t.Errorf("Expected to be able to delete %q, %s", key, err)
			}
		}
		keys, err = c.Keys()
		if err != nil {
			t.Errorf("Expected to get keys back from List, %s", err)
			t.FailNow()
		}
		if len(keys) != 1 {
			t.Errorf("Expected one key left after delete, got %+v", keys)
			t.FailNow()
		}
	*/
}

func TestFredaExample(t *testing.T) {
	// Create a collection "mystuff" inside the directory called demo
	cName := path.Join("testout", "freda1.ds")
	if _, err := os.Stat(cName); err == nil {
		// Clear stale data if needed.
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "pairtree")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()
	// Create a JSON document
	docName := "freda.json"
	document := map[string]interface{}{
		"name":  "freda",
		"email": "freda@inverness.example.org",
	}
	if err := c.Create(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Keys should get a key back ...
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("expected c.Keys(), got error %s", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected one key, got %d", len(keys))
	} else if keys[0] != docName {
		t.Errorf("expected key %q, got %q", docName, keys[0])
	}
	// Read a JSON document
	if err := c.Read(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Update a JSON document
	document["email"] = "freda@zbs.example.org"
	if err := c.Update(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Delete a JSON document
	if err := c.Delete(docName); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
}

func TestSQLStoreFredaExample(t *testing.T) {
	// Create a collection "mystuff" inside the directory called demo
	cName := path.Join("testout", "freda2.ds")
	if _, err := os.Stat(cName); err == nil {
		// Clear stale data if needed.
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "sqlite://collection.db")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()
	// Create a JSON document
	docName := "freda.json"
	document := map[string]interface{}{
		"name":  "freda",
		"email": "freda@inverness.example.org",
	}
	if err := c.Create(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Keys should get a key back ...
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("expected c.Keys(), got error %s", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected one key, got %d", len(keys))
	} else if keys[0] != docName {
		t.Errorf("expected key %q, got %q", docName, keys[0])
	}
	start := "2000-01-01 00:00:00"
	end := "2100-12-31 23:23:59"
	keys, err = c.UpdatedKeys(start, end)
	if err != nil {
		t.Errorf("expected c.UpdatedKeys(%q, %q), got error %s", start, end, err)
	}
	if len(keys) != 1 {
		t.Errorf("expected one key, got %d", len(keys))
	} else if keys[0] != docName {
		t.Errorf("expected key %q, got %q", docName, keys[0])
	}

	// Read a JSON document
	if err := c.Read(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Update a JSON document
	document["email"] = "freda@zbs.example.org"
	if err := c.Update(docName, document); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Delete a JSON document
	if err := c.Delete(docName); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
}

//
// Ported tests from v1
//

func TestCollection(t *testing.T) {
	colName := "testout/test_collection.ds"
	// Remove any pre-existing test data
	if _, err := os.Stat(colName); err == nil {
		os.RemoveAll(colName)
	}

	// Create a new collection
	c, err := Init(colName, "pairtree")
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
	err = c.Close()
	if err != nil {
		t.Errorf("error Close() a collection %q", err)
		t.FailNow()
	}

	// Now open the existing collection of colName
	c, err = Open(colName)
	if err != nil {
		t.Errorf("error Open() a collection %q", err)
		t.FailNow()
	}

	keys, _ := c.Keys()
	if len(keys) > 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
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
			err = c.Create(id, rec)
			if err != nil {
				t.Errorf("%q: collection.Create(), %s", c.Name, err)
				t.FailNow()
			}
			if c.HasKey(id) == false {
				t.Errorf("%q was not created in %q, no error valuye returned", id, c.Name)
				t.FailNow()
			}
		}
	}
	keys, _ = c.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %+v", keys)
		t.FailNow()
	}

	// Create an empty record, then read it again to compare
	keyName := "Kahlo-F"
	rec2 := map[string]interface{}{}
	err = c.Read(keyName, rec2)
	if err != nil {
		t.Errorf("%q: Read(), %s", c.Name, err)
		t.FailNow()
	}
	rec1 := testData[0]
	for k, expected := range rec1 {
		if val, ok := rec2[k]; ok == true {
			if expected != val {
				t.Errorf("%q: expected %s in record, got, %s", c.Name, expected, val)
				t.FailNow()
			}
		} else {
			t.Errorf("%q: Read() missing %s in %+v, %+v", c.Name, k, rec1, rec2)
			t.FailNow()
		}
	}
	// Should trigger update if a duplicate record
	err = c.Create(keyName, rec2)
	if err == nil {
		t.Errorf("%q: Create not allow creationg on an existing record, %s --> %+v", c.Name, keyName, rec2)
		t.FailNow()
	}

	rec3 := map[string]interface{}{}
	if err := c.Read(keyName, rec3); err != nil {
		t.Errorf("%q: Should have found freda in collection, %s", c.Name, err)
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
	err = c.Update(keyName, rec2)
	if err != nil {
		t.Errorf("%s: Could not update %s, %s", c.Name, "freda", err)
		t.FailNow()
	}

	rec4 := map[string]interface{}{}
	if err := c.Read(keyName, rec4); err != nil {
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

	err = c.Delete(keyName)
	if err != nil {
		t.Errorf("Should be able to delete %s, %s", "freda.json", err)
		t.FailNow()
	}
	err = c.Read(keyName, rec2)
	if err == nil {
		t.Errorf("Record should have been deleted, %+v, %s", rec2, err)
	}

	/* FIXME: Do we need this really?
	err = deleteCollection(colName)
	if err != nil {
		t.Errorf("Couldn't remove collection %s, %s", colName, err)
	}
	*/
}

func TestComplexKeys(t *testing.T) {
	cName := path.Join("testout", "complex_keys.ds")
	// remove any stale test collection collection first...
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	// Create a new collection
	c, err := Init(cName, "pairtree")
	if err != nil {
		t.Errorf("error Create() a collection %q", err)
		t.FailNow()
	}
	defer c.Close()
	cnt := c.Length()
	if cnt > 0 {
		t.Errorf("expected 0 objects, got %d", cnt)
	}
	testRecords := map[string]map[string]interface{}{}
	testRecords["agent:person:1"] = map[string]interface{}{
		"name": "George",
		"id":   25,
	}
	testRecords["agent:person:2"] = map[string]interface{}{
		"name": "Carl",
		"id":   2523,
	}
	testRecords["agent:person:3333"] = map[string]interface{}{
		"name": "Mac",
		"id":   2,
	}
	testRecords["agent:person:29994"] = map[string]interface{}{
		"name": "Fred",
		"id":   9925,
	}

	testRecords["agent:person:29"] = map[string]interface{}{
		"name": "Mike",
		"id":   81,
	}
	testRecords["agent:person:100"] = map[string]interface{}{
		"name": "Tim",
		"id":   8,
	}
	testRecords["agent:person:101"] = map[string]interface{}{
		"name": "Kim",
		"id":   101,
	}

	for k, v := range testRecords {
		if err := c.Create(k, v); err != nil {
			t.Errorf("Can't create (%T) %s <-- (%T) %+v : %s", k, k, v, v, err)
		}
	}
}

func TestCaseHandling(t *testing.T) {
	// Setup a test collection and data
	cName := path.Join("testout", "test_case_handling.ds")
	os.RemoveAll(cName)
	c, err := Init(cName, "pairtree")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if _, err := os.Stat(cName); os.IsNotExist(err) {
		t.Errorf(`failed to create %q with c.Init(%q, "")`, cName, cName)
		t.FailNow()
	}
	o := map[string]interface{}{}
	o["a"] = 1
	err = c.Create("A", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["b"] = 2
	err = c.Create("B", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["c"] = 3
	err = c.Create("C", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Get back a list of keys, should all be lowercase.const
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("c.Keys() should not return error, %s", err)
		t.FailNow()
	}
	for _, key := range keys {
		if key == strings.ToUpper(key) {
			t.Errorf("Expected lower case %q, got %q in %q", strings.ToLower(key), key, cName)
		}
		p, err := c.DocPath(strings.ToUpper(key))
		if err != nil {
			t.Errorf("%s in %q", err, cName)
			t.FailNow()
		}
		// Check if p has the OS's separator.
		if !strings.Contains(p, fmt.Sprintf("%c", filepath.Separator)) {
			t.Errorf("Path seperator does not match host OS, %q <- %c", p, filepath.Separator)
			t.FailNow()
		}
	}
	cnt := c.Length()
	if cnt != 3 {
		t.Errorf("Expected 3, got %d", cnt)
		t.FailNow()
	}
	c.Close()
}
