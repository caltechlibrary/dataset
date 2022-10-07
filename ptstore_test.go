// ptstorage_test is a of the dataset package.
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
	"sort"
	"testing"
)

//
// Test the storage functionality
//

func TestPTStoreBasic(t *testing.T) {
	os.MkdirAll("testout", 0775)
	sName := path.Join("testout", "pt1.ds")
	if _, err := os.Stat(sName); err == nil {
		os.RemoveAll(sName)
	}
	store, err := PTStoreOpen(sName, "")
	if err != nil {
		t.Errorf(`Open(%q, ""), error %s`, sName, err)
		t.FailNow()
	}
	if store == nil {
		t.Errorf(`Open(%q, ""), store should not be nil`, sName)
		t.FailNow()
	}
	for _, setting := range []int{Major, Minor, Patch, None} {
		if err := store.SetVersioning(setting); err != nil {
			t.Errorf("store.SetVersioning(%d) failed, %s", setting, err)
			t.FailNow()
		}
	}
	objects := []map[string]interface{}{
		{"one": 1},
		{"two": 2},
		{"three": 3},
	}
	keys := []string{}
	for i, obj := range objects {
		key := fmt.Sprintf("%08d", i)
		keys = append(keys, key)
		src, err := json.MarshalIndent(obj, "", "    ")
		if err != nil {
			t.Errorf("json.MarshalIndent() failed for %q, %s", key, err)
			continue
		}
		if err := store.Create(key, src); err != nil {
			t.Errorf("Should be able to create %q in %q, %s", key, sName, err)
			t.FailNow()
		}
		if src, err := store.Read(key); err != nil {
			t.Errorf("store.Read(%q) error, %s", key, err)
			t.FailNow()
		} else if len(src) == 0 {
			t.Errorf("store.Read(%q), empty []byte", key)
			t.FailNow()
		}
		src = []byte(fmt.Sprintf(`{"one": 1, "alt": %q}`, key))
		if err := store.Update(key, src); err != nil {
			t.Errorf("store.Update(%q, %q), error, %s", key, src, err)
			t.FailNow()
		}
		if !store.HasKey(key) {
			t.Errorf("store.HasKey(%q) should have returned true", key)
			t.FailNow()
		}
	}
	tKeys, err := store.Keys()
	if err != nil {
		t.Errorf("store.Keys() error, %s", err)
		t.FailNow()
	}
	sort.Strings(keys)
	sort.Strings(tKeys)
	if len(keys) != len(tKeys) {
		t.Errorf("len(keys) == %d and len(tKeys) == %d, should be same", len(keys), len(tKeys))
		t.FailNow()
	}
	for i := range keys {
		if keys[i] != tKeys[i] {
			t.Errorf("Expected key %q, got %q", keys[i], tKeys[i])
		}
	}
	expectedL, gotL := int64(len(keys)), store.Length()
	if expectedL != gotL {
		t.Errorf("expected %d (%T), got %d (%T)for store.Length()", expectedL, expectedL, gotL, gotL)
	}

	for _, key := range keys {
		if err := store.Delete(key); err != nil {
			t.Errorf("store.Delete(%q) erorr, %s", key, err)
			t.FailNow()
		}
	}

	if err := store.SetVersioning(Patch); err != nil {
		t.Errorf("store.SetVersioning(%d) error, %s", Patch, err)
	}
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%07d", i)
		src := []byte(fmt.Sprintf(`{"one": %d, "two": "%s"}`, i, key))
		if err := store.Create(key, src); err != nil {
			t.Errorf(`store.Create(%q, %s) error, %s`, key, src, err)
			t.FailNow()
		}
		// Check for version after create
		if versions, err := store.Versions(key); err != nil {
			t.Errorf(`store.Versions(%q) error, %s`, key, err)
			t.FailNow()
		} else {
			if len(versions) != 1 {
				t.Errorf("expected 1 version number, %+v", versions)
				t.FailNow()
			}
			for i := range versions {
				expected := fmt.Sprintf("0.0.%d", i+1)
				if versions[i] != expected {
					t.Errorf("expected %q, got %q for version", expected, versions[i])
					t.FailNow()
				}
			}
		}
		// Check for version after update
		src = []byte(fmt.Sprintf(`{"one": %d, "two": "%s", "three": 3.0}`, i, key))
		if err := store.Update(key, src); err != nil {
			t.Errorf(`store.Update(%q, %s) error, %s`, key, src, err)
			t.FailNow()
		}
		if versions, err := store.Versions(key); err != nil {
			t.Errorf(`store.Versions(%q) error, %s`, key, err)
			t.FailNow()
		} else {
			if len(versions) != 2 {
				t.Errorf("expected 1 version number, %+v", versions)
				t.FailNow()
			}
			for i := range versions {
				expected := fmt.Sprintf("0.0.%d", i+1)
				if versions[i] != expected {
					t.Errorf("expected %q, got %q for version", expected, versions[i])
					t.FailNow()
				}
			}
		}
	}

	if err := store.Close(); err != nil {
		t.Errorf("store.Close() failed, %s", err)
		t.FailNow()
	}
}
