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
	"testing"
)

func TestRepair(t *testing.T) {
	o := map[string]interface{}{}
	o["a"] = 1

	// Setup a test collection and data
	cName := "test_repair.ds"
	os.RemoveAll(cName)
	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	err = c.Create("a", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["b"] = 2
	err = c.Create("b", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	o["c"] = 3
	err = c.Create("c", o)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	// Break the collection by removing a file from disc.
	p, err := c.DocPath("b")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	os.Remove(p)
	cnt := c.Length()
	if cnt != 3 {
		t.Errorf("Expected 3, got %d", cnt)
		t.FailNow()
	}
	c.Close()
	err = Repair(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	c, err = Open(cName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	defer c.Close()
	cnt = c.Length()
	if cnt != 2 {
		t.Errorf("Expected 2, got %d", cnt)
		t.FailNow()
	}
}
