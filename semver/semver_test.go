//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"strings"
	"testing"
)

func TestSemver(t *testing.T) {
	expected := "v1.1.1"
	v, err := ParseSemver([]byte(expected))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result := v.String()
	if expected != result {
		t.Errorf("expected %q, got %q", expected, result)
	}

	expected = "v1.1"
	v, err = ParseSemver([]byte(expected))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result = v.String()
	if expected != result {
		t.Errorf("expected %q, got %q", expected, result)
	}

	expected = "vA1.2.3"
	v, err = ParseSemver([]byte(expected))
	if err == nil {
		t.Errorf("expected an error, returns %s", v.ToJSON())
	}
}

func TestIncrement(t *testing.T) {
	b := []byte("v0.0.0")
	version, err := ParseSemver(b)
	if err != nil {
		t.Errorf("Failed to parse %q, got %s", b, err)
		t.FailNow()
	}
	expected := "v0.0.1"
	if err = version.IncPatch(); err != nil {
		t.Errorf("Error, increment patch version, %q, %s", version.String(), err)
		t.FailNow()
	}

	s := version.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Expected %q, got %q", expected, s)
	}

	if err = version.IncMinor(); err != nil {
		t.Errorf("Error, increment minor version, %q, %s", version.String(), err)
		t.FailNow()
	}
	expected = "v0.1.0"
	s = version.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Expected %q, got %q", expected, s)
	}

	if err = version.IncMajor(); err != nil {
		t.Errorf("Error, increment major version, %q, %s", version.String(), err)
		t.FailNow()
	}
	expected = "v1.0.0"
	s = version.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Expected %q, got %q", expected, s)
	}
}
