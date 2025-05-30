// sqlstore is a sub module of the dataset package.
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
	"fmt"
	"os"
	"path"
	"sort"
	"testing"
)

//
// Test the storage functionality
//

func TestSQLStoreBasic(t *testing.T) {
	os.MkdirAll("testout", 0775)
	sName := path.Join("testout", "sqls1.ds")
	sDsnURI := "sqlite://testout/sqls1.ds/collection.db"
	if _, err := os.Stat(sName); err == nil {
		os.RemoveAll(sName)
	}
	os.MkdirAll(sName, 0775)
	store, err := SQLStoreInit(sName, sDsnURI)
	if err != nil {
		t.Errorf("failed to create table %q, %s", sName, err)
		t.FailNow()
	}
	if store == nil {
		t.Errorf(`Init(%q, %q), store should not be nil`, sName, sDsnURI)
		t.FailNow()
	}
	store.Close()

	store, err = SQLStoreOpen(sName, sDsnURI)
	if err != nil {
		t.Errorf(`SQLStoreOpen(%q, %q), error %s`, sName, sDsnURI, err)
		t.FailNow()
	}
	if store == nil {
		t.Errorf(`SQLStoreOpen(%q, %q), store should not be nil`, sName, sDsnURI)
		t.FailNow()
	}
	// Setup main databases
	objects := []map[string]interface{}{
		{"one": 1},
		{"two": 2},
		{"three": 3},
	}
	keys := []string{}
	for i, obj := range objects {
		key := fmt.Sprintf("%08d", i)
		keys = append(keys, key)
		src, err := JSONMarshalIndent(obj, "", "    ")
		if err != nil {
			t.Errorf("JSONMarshalIndent() failed for %q, %s", key, err)
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

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("%07d", i)
		src := []byte(fmt.Sprintf(`{"one": %d, "two": "%s"}`, i, key))
		if err := store.Create(key, src); err != nil {
			t.Errorf(`store.Create(%q, %s) error, %s`, key, src, err)
			t.FailNow()
		}
		// Check for update
		src = []byte(fmt.Sprintf(`{"one": %d, "two": "%s", "three": 3.0}`, i, key))
		if err := store.Update(key, src); err != nil {
			t.Errorf(`store.Update(%q, %s) error, %s`, key, src, err)
			t.FailNow()
		}
	}

	if err := store.Close(); err != nil {
		t.Errorf("store.Close() failed, %s", err)
		t.FailNow()
	}
}
