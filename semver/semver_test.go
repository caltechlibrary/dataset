//
// semver is a semantic version number package used by dataset.
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
package semver

import (
	"strings"
	"testing"
)

func TestSemver(t *testing.T) {
	expected := "1.1.1"
	v, err := Parse([]byte(expected))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result := v.String()
	if expected != result {
		t.Errorf("expected %q, got %q", expected, result)
	}

	expected = "1.1"
	v, err = Parse([]byte(expected))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result = v.String()
	if expected != result {
		t.Errorf("expected %q, got %q", expected, result)
	}

	expected = "A1.2.3"
	v, err = Parse([]byte(expected))
	if err == nil {
		t.Errorf("expected an error, returns %s", v.ToJSON())
	}

	expected = "2.0.0-next"
	v, err = Parse([]byte(expected))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result = v.String()
	if expected != result {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestIncrement(t *testing.T) {
	b := []byte("v0.0.0")
	version, err := Parse(b)
	if err != nil {
		t.Errorf("Failed to parse %q, got %s", b, err)
		t.FailNow()
	}
	expected := "0.0.1"
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
	expected = "0.1.0"
	s = version.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Expected %q, got %q", expected, s)
	}

	if err = version.IncMajor(); err != nil {
		t.Errorf("Error, increment major version, %q, %s", version.String(), err)
		t.FailNow()
	}
	expected = "1.0.0"
	s = version.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Expected %q, got %q", expected, s)
	}
}

func TestLess(t *testing.T) {
	a := []byte("100.0.0")
	b := []byte("9.9.9")
	aSemver, err := Parse(a)
	if err != nil {
		t.Errorf("failed to parse %q, %s", a, err)
		t.FailNow()
	}
	bSemver, err := Parse(b)
	if err != nil {
		t.Errorf("failed to parse %q, %s", b, err)
		t.FailNow()
	}

	expected := false
	got := Less(aSemver, bSemver)
	if expected != got {
		t.Errorf(`Expected %t, got %t for %s < %s`, expected, got, aSemver.String(), bSemver.String())
		t.FailNow()
	}
}

func TestSorting(t *testing.T) {
	inStrings := []string{
		"0.101.1",
		"1.1.1",
		"111.111.111",
		"1.1.0",
		"1.0.1",
		"0.0.0",
		"0.1.1",
		"1.1.1",
		"10.9.1",
		"0.9.9",
		"1.9.1",
		"9.5.3",
		"0.1.0",
		"12.1.11",
		"10.10.10",
		"11.11.11",
		"2.0.0-next",
	}
	expectedStrings := []string{
		"0.0.0",
		"0.1.0",
		"0.1.1",
		"0.9.9",
		"0.101.1",
		"1.0.1",
		"1.1.0",
		"1.1.1",
		"1.1.1",
		"1.9.1",
		"2.0.0-next",
		"9.5.3",
		"10.9.1",
		"10.10.10",
		"11.11.11",
		"12.1.11",
		"111.111.111",
	}
	gotStrings := SortStrings(inStrings)
	for i := 0; i < len(expectedStrings); i++ {
		expectedVersion := expectedStrings[i]
		gotVersion := gotStrings[i]
		if expectedVersion != gotVersion {
			t.Errorf("expected (string) %q at %d, got %q", expectedVersion, i, gotVersion)
		}
	}

	inSemvers := []*Semver{}
	for _, val := range inStrings {
		sv, err := Parse([]byte(val))
		if err != nil {
			t.Errorf("inStrings semvar parse failed for %q, %s", val, err)
			t.FailNow()
		}
		inSemvers = append(inSemvers, sv)
	}
	Sort(inSemvers)
	for i := 0; i < len(expectedStrings); i++ {
		expectedVersion := expectedStrings[i]
		gotVersion := inSemvers[i].String()
		if expectedVersion != gotVersion {
			t.Errorf("expected (*Semver) %q at %d, got %q", expectedVersion, i, gotVersion)
		}
	}
}
