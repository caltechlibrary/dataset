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
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAttachments(t *testing.T) {
	cName := "testdata/pairtree_layout/col3.ds"
	os.RemoveAll(cName)

	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("Can't create collection %q (%d)", cName, err)
		t.FailNow()
	}

	record := map[string]interface{}{
		"name":  "freda",
		"motto": "it's all about what you sense when you've have sense to sense it",
	}
	data := &Attachment{
		Name: "impressed.txt",
		Body: []byte("Wowie Zowie!"),
	}

	if err := c.Create("freda", record); err != nil {
		t.Errorf("failed to create freda in %s, %s", c.Name, err)
		t.FailNow()
	}
	if err := c.attach("freda", data); err != nil {
		t.Errorf("failed to add attachments to %s, %s", c.Name, err)
		t.FailNow()
	}
	if files, err := c.Attachments("freda"); err != nil {
		t.Errorf("can't list attachments for freda in %s, %s", c.Name, err)
		t.FailNow()
	} else {
		if len(files) != 1 {
			t.Errorf("Expected one file attached, %+v", files)
			t.FailNow()
		}
		if files[0] != "impressed.txt 12" {
			t.Errorf("Expected files[0] to be impressed, got %+v", files)
			t.FailNow()
		}
	}
	if attachments, err := c.getAttached("freda"); err != nil {
		t.Errorf("Expected attachments, %s", err)
		t.FailNow()
	} else {
		if len(attachments) != 1 {
			t.Errorf("Expected one attachment, %+v\n", attachments)
			t.FailNow()
		}
		for _, a := range attachments {
			if (a.Name == "impressed.txt" && bytes.Compare(a.Body, []byte("Wowie Zowie!")) == 0) == false {
				t.Errorf("Expected impressed.txt, got %+v", a)
				t.FailNow()
			}
		}
	}

	if attachments, err := c.getAttached("freda", "impressed.txt"); err != nil {
		t.Errorf("Expected attachments, %s", err)
		t.FailNow()
	} else {
		if len(attachments) != 1 {
			t.Errorf("Expected one attachment, %+v\n", attachments)
			t.FailNow()
		}
		for _, a := range attachments {
			if (a.Name == "impressed.txt" && bytes.Compare(a.Body, []byte("Wowie Zowie!")) == 0) == false {
				t.Errorf("Expected impressed.txt, got %+v", a)
				t.FailNow()
			}
		}
	}

	if err := c.attach("freda", &Attachment{Name: "what/she/smokes.txt", Body: []byte("A Havana Cigar")}); err != nil {
		t.Errorf("Appending attachment, %s", err)
		t.FailNow()
	}

	if files, err := c.Attachments("freda"); err != nil {
		t.Errorf("Attachments after append, %+v %s", files, err)
		t.FailNow()
	} else {
		if len(files) != 1 {
			t.Errorf("Should have one file after appending an attachment (each call to attach should generate a fresh tarball)")
		}
		for _, s := range files {
			if s != "impressed.txt" && s != "what/she/smokes.txt 14" {
				t.Errorf("Unexpected file in list, %s", s)
			}
		}
	}

	if attachments, err := c.getAttached("freda", "what/she/smokes.txt"); err != nil {
		t.Errorf("Expected attachments, %s", err)
		t.FailNow()
	} else {
		if len(attachments) != 1 {
			t.Errorf("Expected one attachment, %+v\n", attachments)
			t.FailNow()
		}
		for _, a := range attachments {
			if (a.Name == "what/she/smokes.txt" && bytes.Compare(a.Body, []byte("A Havana Cigar")) == 0) == false {
				t.Errorf("Expected what/she/smokes.txt, got %+v", a)
				t.FailNow()
			}
		}
	}

	if err := c.Prune("freda", "what/she/smokes.txt"); err != nil {
		t.Errorf("Delete one file, %s", err)
	}
	tarDocPath, err := c.DocPath("freda")
	if err != nil {
		t.Errorf("Should have gotten docpath for freda, %s", err)
		t.FailNow()
	}
	tarDocPath = strings.TrimSuffix(tarDocPath, ".json") + ".tar"

	if _, err := os.Stat(tarDocPath); err != nil {
		t.Errorf("Shouldn't have deleted %s, %s", tarDocPath, err)
		t.FailNow()
	}

	if err := c.Prune("freda"); err != nil {
		t.Errorf("Delete whole tarball, %s", err)
		t.FailNow()
	}

	if _, err := os.Stat(tarDocPath); os.IsNotExist(err) == false {
		t.Errorf("Should have deleted %s, %s", tarDocPath, err)
		t.FailNow()
	}
}
