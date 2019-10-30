//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2019, Caltech
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
	"strings"
	"testing"
)

func TestService(t *testing.T) {
	cName := path.Join("testdata", "service0.ds")
	key := "k1"
	src := []byte(`{
	"title": "Orchids & Moonbeams",
	"cast": [
		{ 
			"last_name": "Lorick",
			"first_name": "Robert",
			"character": "Jack Flanders"
		},
		{
			"last_name": "Adams",
			"first_name": "Dave",
			"character": "Mojo Sam"
		},
		{
			"last_name": "Poirier",
			"first_name": "Pascale",
			"character": "Claudine"
		},
		{
			"last_name": "Donovan",
			"first_name": "Patrick",
			"character": "Pat Patternson"
		},
		{
			"last_name": "Goodhart Hebert",
			"first_name": "Camille",
			"character": "Bunny"
		},
		{
			"last_name": "Roth",
			"first_name": "Laura",
			"character": "Amber"
		}
	]
}`)

	os.RemoveAll(cName)
	InitCollection(cName)
	if err := ServiceOpen(cName); err != nil {
		t.Errorf("ServiceOpen(%q) failed, %s", cName, err)
		t.FailNow()
	}
	defer func() {
		if err := ServiceClose(cName); err != nil {
			t.Errorf("Failed CloseService(%q), %s", cName, err)
		}
	}()

	cNames := ServiceCollections()
	if len(cNames) != 1 {
		t.Errorf("Expected one cName %q, got (%d) %s", cName, len(cNames), strings.Join(cNames, ", "))
		t.FailNow()
	}
	if cNames[0] != cName {
		t.Errorf("Expected cName %q, got %q", cName, cNames[0])
	}

	if err := ServiceCreateJSON(cName, key, src); err != nil {
		t.Errorf("Expected a new record got error, %s", err)
		t.FailNow()
	}

	if src, err := ServiceReadJSON(cName, key); err != nil {
		t.Errorf("expected to find %q, got error, %s", key, err)
	} else {
		obj := map[string]interface{}{}
		if err := DecodeJSON(src, &obj); err != nil {
			t.Errorf("expected to decode json, got %s", err)
			t.FailNow()
		}
		//fmt.Printf("DEBUG %T -> %+v", obj, obj)
		if val, ok := obj["title"]; ok == false {
			t.Errorf("expected a title, found none.")
		} else {
			s := val.(string)
			if s != "Orchids & Moonbeams" {
				t.Errorf("expected a title 'Orchid & Moonbeams', got %q.", s)
			}
		}
		if val, ok := obj["cast"]; ok == false {
			t.Errorf("Expected a cast, found none")
		} else {
			ol := val.([]interface{})
			if len(ol) != 6 {
				t.Errorf("Expected six objects, got %+v", ol)
			}
		}
		obj["produced_by"] = "ZBS Fountation"
		obj["writer"] = "Tom Lopez, a.k.a. Meatball Fulton"
		src, err = EncodeJSON(obj)
		if err != nil {
			t.Errorf("expected json encoded, got error %s", err)
			t.FailNow()
		}
		if err := ServiceUpdateJSON(cName, key, src); err != nil {
			t.Errorf("Expected an updated record got error, %s", err)
			t.FailNow()
		}
		keys := ServiceKeys(cName)
		if len(keys) != 1 {
			t.Errorf("Expected one key, k1 got %d", len(keys))
			t.FailNow()
		}
		if keys[0] != "k1" {
			t.Errorf("Expected k1, got %q", keys[0])
			t.FailNow()
		}
		fName := "f1"
		dotPaths := []string{".title", ".cast[:].character"}
		labels := []string{"title", "characters"}
		verbose := false
		f, err := ServiceFrameCreate(cName, fName, keys, dotPaths, labels, verbose)
		if err != nil {
			t.Errorf("expected success for FrameCreate(), got %s", err)
			t.FailNow()
		}
		ol := f.Objects()
		if len(ol) != 1 {
			t.Errorf("Expected on object in frame, got %d", len(ol))
			t.FailNow()
		}
		ol2, err := ServiceFrameObjects(cName, fName)
		if err != nil {
			t.Errorf("Expected an object list for frame %q in %q", fName, cName)
			t.FailNow()
		}
		if len(ol) != len(ol2) {
			t.Errorf("expected len(ol) to equal len(ol2), got %d != %d", len(ol), len(ol2))
			t.FailNow()
		}
		src = []byte(`{
	"title": "The Incredible Adventures of Jack Flanders",
	"cast": [
		{ 
			"last_name": "Lorick",
			"first_name": "Robert",
			"character": "Jack Flanders"
		},
		{
			"last_name": "Adams",
			"first_name": "Dave",
			"character": "Mojo Sam"
		},
		{
			"last_name": "Orte",
			"first_name": "P. J.",
			"character": "Little Freda"
		}
	]
}`)
		key2 := "k0"
		if err := ServiceCreateJSON(cName, key2, src); err != nil {
			t.Errorf("expected success for CreateJSON k0, got %s", err)
			t.FailNow()
		}
		src, err := ServiceReadJSON(cName, key)
		if err != nil {
			t.Errorf("expected ReadJSON(%q, %q)", cName, key)
			t.FailNow()
		}
		if len(src) == 0 {
			t.Errorf("expected ReadJSON(%q, %q) -> %s", cName, key, src)
			t.FailNow()
		}
		if err := ServiceFrameRefresh(cName, fName, []string{key2, key}, verbose); err != nil {
			t.Errorf("expected success frame refresh, got %s", err)
			t.FailNow()
		}
		ol3, err := ServiceFrameObjects(cName, fName)
		if err != nil {
			t.Errorf("expected a new copy of object list, got %s", err)
			t.FailNow()
		}
		if len(ol3) != 2 {
			t.Errorf("expected length 2, got %d, %T -> %+v", len(ol3), ol3, ol3)
			t.FailNow()
		}
		if ol3[0]["title"] != "Orchids & Moonbeams" {
			t.Errorf("Expected first object to be Orchids & Moonbeams, got %+v", ol3[0])
		}
		if ol3[1]["title"] != "The Incredible Adventures of Jack Flanders" {
			t.Errorf("Expected first object to be The Incredible Adventures of Jack Flanders, got %+v", ol3[1])
		}
		if err := ServiceFrameReframe(cName, fName, []string{key2, key}, verbose); err != nil {
			t.Errorf("expected success frame refresh, got %s", err)
			t.FailNow()
		}
		ol3, err = ServiceFrameObjects(cName, fName)
		if err != nil {
			t.Errorf("expected a new copy of object list, got %s", err)
			t.FailNow()
		}
		if len(ol3) != 2 {
			t.Errorf("expected length 2, got %d, %T -> %+v", len(ol3), ol3, ol3)
			t.FailNow()
		}
		if ol3[0]["title"] != "The Incredible Adventures of Jack Flanders" {
			t.Errorf("Expected first object to be The Incredible Adventures of Jack Flanders, got %+v", ol3[0])
		}
		if ol3[1]["title"] != "Orchids & Moonbeams" {
			t.Errorf("Expected first object to be Orchids & Moonbeams, got %+v", ol3[1])
		}
		err = ServiceFrameClear(cName, fName)
		if err != nil {
			t.Errorf("expected a no error from ServiceFrameClear(%q, %q), got %s", cName, fName, err)
			t.FailNow()
		}
		ol3, err = ServiceFrameObjects(cName, fName)
		if err != nil {
			t.Errorf("expected a no error from ServiceFrameObjects(%q, %q), got %s", cName, fName, err)
			t.FailNow()
		}
		if len(ol3) != 0 {
			t.Errorf("expected zero objects, got %+v", ol3)
			t.FailNow()
		}
		err = ServiceFrameDelete(cName, fName)
		if err != nil {
			t.Errorf("expected a no error from ServiceFrameDelete(%q, %q), got %s", cName, fName, err)
			t.FailNow()
		}
		if fNames := ServiceFrames(cName); len(fNames) > 0 {
			t.Errorf("expected zsero frame names, got %s", strings.Join(fNames, ", "))
		}
	}
}
