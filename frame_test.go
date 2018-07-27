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
	"os"
	"path"
	"testing"
)

func TestFrame(t *testing.T) {
	layouts := []int{
		BUCKETS_LAYOUT,
		PAIRTREE_LAYOUT,
	}
	for _, cLayout := range layouts {
		os.RemoveAll(path.Join("testdata", "frame_test.ds"))
		cName := path.Join("testdata", "frame_test.ds")
		c, err := InitCollection(cName, cLayout)
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
			err := c.Create(key, rec)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		}

		f, err := c.Frame("frame-1", keys, []string{".id", ".one", ".two", ".three", ".four"}, false)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		expected := "frame-1"
		result := f.Name
		if expected != result {
			t.Errorf("expected %q, got %q, for %s", expected, result, f)
		}
		expected = cName // e.g. "testdata/frame_test.ds"
		result = f.CollectionName
		if expected != result {
			t.Errorf("expected %q, got %q, for %s", expected, result, f)
		}
	}
}
