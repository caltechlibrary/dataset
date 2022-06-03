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
	"os"
	"path"
	"testing"
)

func TestPTStoreVersioning(t *testing.T) {
	cName := path.Join("testout", "vmajor.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	if err := c.SetVersioning("major"); err != nil {
		t.Errorf("failed to set versioning to major for %q, %s", cName, err)
		t.FailNow()
	}
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()

	if c.Versioning != "major" {
		t.Errorf("Versioning wasn't set to major, %q", c.Versioning)
		t.FailNow()
	}
	key := "123"
	obj := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "Hi There!",
	}
	if err := c.Create(key, obj); err != nil {
		t.Errorf("Expected c.Create(%q, obj) to work, %s", key, err)
		t.FailNow()
	}
	versions, err := c.Versions(key)
	if len(versions) != 1 {
		t.Errorf("Expected a single version number got %+v", versions)
		t.FailNow()
	}
	if versions[0] != "1.0.0" {
		t.Errorf("Expected 1.0.0 version got %q for %s", versions[0], key)
		t.FailNow()
	}
}

func TestSQLStoreVersioning(t *testing.T) {
	cName := path.Join("testout", "vmajor_sql.ds")
	dsnURI := "sqlite://testout/vmajor_sql.ds/collection.db"
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	if err := c.SetVersioning("major"); err != nil {
		t.Errorf("failed to set versioning to major for %q, %s", cName, err)
		t.FailNow()
	}
	if _, err := os.Stat(path.Join(cName, "versioning.json")); os.IsNotExist(err) {
		t.Errorf("Missing versioning.json from %q", cName)
		t.FailNow()
	}
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Unabled to create %q, %s", cName, err)
		t.FailNow()
	}
	defer c.Close()

	if c.Versioning != "major" {
		t.Errorf("Versioning wasn't set to major, %q", c.Versioning)
		t.FailNow()
	}
	key := "123"
	obj := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": "Hi There!",
	}
	if err := c.Create(key, obj); err != nil {
		t.Errorf("Expected c.Create(%q, obj) to work, %s", key, err)
		t.FailNow()
	}
	versions, err := c.Versions(key)
	if len(versions) != 1 {
		t.Errorf("Expected a single version number got %+v", versions)
		t.FailNow()
	}
	if versions[0] != "1.0.0" {
		t.Errorf("Expected 1.0.0 version got %q for %s", versions[0], key)
		t.FailNow()
	}
}
