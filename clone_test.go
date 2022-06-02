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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestClone(t *testing.T) {
	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}

	cName, dsnURI := path.Join("testout", "zbs1.ds"), ""
	// cleanup stale data
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Failed to create seed collection %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()

	// Populate our seed collection
	for k, v := range testRecords {
		if err := c.Create(k, v); err != nil {
			t.Errorf("Could not create %q in %q (seed collection), %s", k, cName, err)
		}
	}
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("Could not get keys from %q, %s", cName, err)
		t.FailNow()
	}

	// Make clone collection
	ncName, ncDsnURI := path.Join("testout", "zbs2.ds"), "sqlite://testout/zbs2.ds/collections.db"
	if _, err := os.Stat(ncName); err == nil {
		os.RemoveAll(ncName)
	}
	if err := c.Clone(ncName, ncDsnURI, keys[0:5], false); err != nil {
		t.Errorf("clone failed, %q to %q, %s", cName, ncName, err)
		t.FailNow()
	}

	// Make sure clone has records.
	nc, err := Open(ncName)
	if err != nil {
		t.Errorf("failed to open clone %q, %s", ncName, err)
		t.FailNow()
	}
	defer nc.Close()
	for _, key := range keys[0:5] {
		if c.HasKey(key) != nc.HasKey(key) {
			t.Errorf("Expected %q in %q %t, got %q in %q %t", key, cName, c.HasKey(key), key, ncName, nc.HasKey(key))
		}
	}
}

func TestCloneSample(t *testing.T) {
	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}
	p := "testout"
	cName := path.Join(p, "test_zbs_characters.ds")
	trainingName := path.Join(p, "test_zbs_training.ds")
	testName := path.Join(p, "test_zbs_test.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	if _, err := os.Stat(trainingName); err == nil {
		os.RemoveAll(trainingName)
	}
	if _, err := os.Stat(testName); err == nil {
		os.RemoveAll(testName)
	}

	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("Can't create %s, %s", cName, err)
		t.FailNow()
	}
	for key, value := range testRecords {
		err := c.Create(key, value)
		if err != nil {
			t.Errorf("Can't add %s to %s, %s", key, cName, err)
			t.FailNow()
		}
	}
	cnt := int(c.Length())
	trainingSize := 4
	testSize := cnt - trainingSize
	keys, err := c.Keys()
	if err != nil {
		t.Errorf("Expected keys in collection to clone, %s", err)
		t.FailNow()
	}
	if err := c.CloneSample(trainingName, "", testName, "", keys, trainingSize, false); err != nil {
		t.Errorf("Failed to create samples %s (%d) and %s, %s", trainingName, trainingSize, testName, err)
	}
	training, err := Open(trainingName)
	if err != nil {
		t.Errorf("Could not open %s, %s", trainingName, err)
		t.FailNow()
	}
	defer training.Close()
	test, err := Open(testName)
	if err != nil {
		t.Errorf("Could not open %s, %s", testName, err)
		t.FailNow()
	}
	defer test.Close()

	if trainingSize != int(training.Length()) {
		t.Errorf("Expected %d, got %d for %s", trainingSize, training.Length(), trainingName)
	}
	if testSize != int(test.Length()) {
		t.Errorf("Expected %d, got %d for %s", testSize, test.Length(), testName)
	}

	keys, _ = c.Keys()
	for _, key := range keys {
		switch {
		case training.HasKey(key) == true:
			if test.HasKey(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		case test.HasKey(key) == true:
			if training.HasKey(key) == true {
				t.Errorf("%s and %s has key %s", trainingName, testName, key)
			}
		default:
			t.Errorf("Could not find %s in %s or %s", key, trainingName, testName)
		}
	}
}

// TestCloneLongKeys added test case
func TestCloneLongKeys(t *testing.T) {
	// Create a source dataset
	cName := path.Join("testout", "src1.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	source, err := Init(cName, "")
	if err != nil {
		t.Errorf("Can't create source dataset, %s", err)
		t.FailNow()
	}
	defer source.Close()

	// Setup a collection to clone
	src := []byte(`[
	{ "one": 1 },
	{ "two": 2 },
	{ "three": 3 },
	{ "four": 4 }
]`)
	testData := []map[string]interface{}{}
	err = json.Unmarshal(src, &testData)
	if err != nil {
		t.Errorf("Can't create testdata")
		t.FailNow()
	}
	for i, obj := range testData {
		key := fmt.Sprintf("%+08d", i)
		if err := source.Create(key, obj); err != nil {
			t.Errorf("failed to create JSON doc for %q in %q, %s", key, cName, err)
			t.FailNow()
		}
	}
	if source.Length() != 4 {
		t.Errorf("Expected 4 documents in our source repository")
		t.FailNow()
	}
	keys, err := source.Keys()
	if err != nil {
		t.Errorf("can't retrieve source keys, %s", err)
		t.FailNow()
	}
	// Clone connection
	tName := path.Join("testout", "dst1.ds")
	if _, err := os.Stat(tName); err == nil {
		os.RemoveAll(tName)
	}
	err = source.Clone(tName, "", keys, false)
	if err != nil {
		t.Errorf(`expected source.Clone(%q, "") to succeed, %s`, tName, err)
		t.FailNow()
	}
	// Open cloned repository
	dest, err := Open(tName)
	if err != nil {
		t.Errorf(`failed to open clone %q, %s`, tName, err)
		t.FailNow()
	}
	// Check if they got copied successful

	for _, key := range keys {
		if !dest.HasKey(key) {
			t.Errorf("Expected %q to have key %q, missing", tName, key)
		} else {
			obj := map[string]interface{}{}
			if err := dest.Read(key, obj); err != nil {
				t.Errorf("Expected dest.Read(%q, obj) to succeed, %s", key, err)
				t.FailNow()
			}
			if len(obj) == 0 {
				t.Errorf("Expected object %q to have content", key)
				t.FailNow()
			}
		}
	}
}
