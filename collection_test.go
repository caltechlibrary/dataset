//
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
//
package dataset

import (
	"fmt"
	"os"
	"path"
	"sort"
	"testing"
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
	c, err := Init(cName, "", PTSTORE)
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
				} else if v != v2 {
					t.Errorf("Expected value 1 %+v to equal value 2 %+v", v, v2)
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
					} else if v != v2 {
						t.Errorf("Expected value 1 %+v to equal value 2 %+v", v, v2)
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
}

// Test the SQL storage implementation
func TestSQLStore(t *testing.T) {
	threeObjects := []map[string]interface{}{}
	threeObjects = append(threeObjects, map[string]interface{}{"one": 1})
	threeObjects = append(threeObjects, map[string]interface{}{"two": 2})
	threeObjects = append(threeObjects, map[string]interface{}{"three": 3})

	cName := path.Join("T2.ds")
	dsnURI := "sqlite:file:testout/T2.db?cache=shared"
	// Clear stale test output
	if _, err := os.Stat(path.Join("testout", cName)); err == nil {
		os.RemoveAll(cName)
	}
	// NOTE: SQLStore requires a DSN URI so it is NOT empty
	c, err := Init(cName, dsnURI, SQLSTORE)
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
				} else if v != v2 {
					t.Errorf("Expected value 1 %+v to equal value 2 %+v", v, v2)
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
					} else if v != v2 {
						t.Errorf("Expected value 1 %+v to equal value 2 %+v", v, v2)
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
}
