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
	"path"
	"strings"
	"testing"
)

func TestService(t *testing.T) {
	cName := path.Join("testdata", "service0.ds")
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

	if cNames, err := ServiceCollections(); err != nil {
		t.Errorf("ServiceCollections() failed, %s", err)
		t.FailNow()
	} else {
		if len(cNames) != 1 {
			t.Errorf("Expected one cName %q, got (%d) %s", cName, len(cNames), strings.Join(cNames, ", "))
			t.FailNow()
		}
		if cNames[0] != cName {
			t.Errorf("Expected cName %q, got %q", cName, cNames[0])
		}
	}

}
