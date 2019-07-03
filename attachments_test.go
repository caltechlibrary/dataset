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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestAttachments(t *testing.T) {
	cName := path.Join("testdata", "col3.ds")
	os.RemoveAll(cName)

	c, err := InitCollection(cName)
	if err != nil {
		t.Errorf("Can't create collection %q (%d)", cName, err)
		t.FailNow()
	}

	// Create some temp text files to attach.
	keyName := "freda"
	buf := []byte("Hello World")
	size := int64(len(buf))
	if size != int64(11) {
		t.Errorf("Expected 'Hello World' to be size 11, got %d", size)
	}
	expectedChecksum := "b10a8db164e0754105b7a99be72e3fe5"
	checksum := fmt.Sprintf("%x", md5.Sum(buf))
	if strings.Compare(checksum, expectedChecksum) != 0 {
		t.Errorf("Expected a checksum %s, got %q", expectedChecksum, checksum)
	}
	if err := ioutil.WriteFile(path.Join("testdata", "helloworld.txt"), buf, 0777); err != nil {
		t.Errorf("Can't create test helloworld.txt, %s", err)
		t.FailNow()
	}

	semver := "v0.0.0"
	motto := []byte("Wowie Zowie!")
	data := &Attachment{
		Name:    "impressed.txt",
		Content: motto,
		Size:    int64(len(motto)),
		Checksums: map[string]string{
			semver: fmt.Sprintf("%x", md5.Sum(motto)),
		},
	}
	if _, err := json.MarshalIndent(data, "", "    "); err != nil {
		t.Errorf("marshal error %s", err)
		t.FailNow()
	}
	if err := ioutil.WriteFile(path.Join("testdata", data.Name), data.Content, 0777); err != nil {
		t.Errorf("Can't create test %q, %s", data.Name, err)
		t.FailNow()
	}

	record := map[string]interface{}{
		"name":  "freda",
		"motto": "it's all about what you sense when you've have sense to sense it",
	}
	if err := c.Create(keyName, record); err != nil {
		t.Errorf("failed to create %s in %s, %s", keyName, c.Name, err)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, semver, path.Join("testdata", "helloworld.txt")); err != nil {
		t.Errorf("failed to add attachments to %s, %s", path.Join("testdata", "helloworld.txt"), err)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, semver, path.Join("testdata", data.Name)); err != nil {
		t.Errorf("failed to add attachments to %s, %s", c.Name, err)
		t.FailNow()
	}
	if files, err := c.Attachments(keyName); err != nil {
		t.Errorf("can't list attachments for %s in %s, %s", keyName, c.Name, err)
		t.FailNow()
	} else {
		if len(files) != 2 {
			t.Errorf("Expected one file attached, %+v", files)
			t.FailNow()
		}
		if files[0] != "helloworld.txt 11" {
			t.Errorf("Expected files[0] to be helloworld.txt, got %+v", files)
			t.FailNow()
		}
		if files[1] != "impressed.txt 12" {
			t.Errorf("Expected files[1] to be impressed.txt, got %+v", files)
			t.FailNow()
		}
	}

	if files, err := c.Attachments(keyName); err != nil {
		t.Errorf("Attachments (2) expected, %+v %s", files, err)
		t.FailNow()
	} else {
		if len(files) != 2 {
			t.Errorf("Should have two files attached")
		}
		for _, s := range files {
			if !(strings.HasPrefix(s, "helloworld.txt") || strings.HasPrefix(s, "impressed.txt")) {
				t.Errorf("Unexpected file in list, %s", s)
			}
		}
	}

	if err := c.Prune("freda", semver, "helloworld.txt"); err != nil {
		t.Errorf("Delete one file, %s", err)
	}
	if files, err := c.Attachments("freda"); err != nil {
		t.Errorf("Attachments (1) expected, %+v %s", files, err)
		t.FailNow()
	} else {
		if err := c.Prune("freda", semver, "impressed.txt"); err != nil {
			t.Errorf("Delete one file, %s", err)
		}
	}

	if err := c.Prune("freda", semver); err != nil {
		t.Errorf("Delete all attachmemts, %s", err)
		t.FailNow()
	}
	// Make sure files have been removed from collection
	docDir := path.Join("testdata", "col3.ds", "pairtree", "fr", "ed", "a", semver)
	for _, fName := range []string{"impressed.txt", "helloworld.txt"} {
		if _, err := os.Stat(path.Join(docDir, fName)); os.IsNotExist(err) == false {
			t.Errorf("Should have deleted %s, %s", path.Join(docDir, fName), err)
			t.FailNow()
		}
	}

	//
	// Now lets tests multiple versions of an attachment
	//
	keyName = "freda"
	semver = "v0.0.1"
	motto = []byte("Wowie Zowie")
	fName := path.Join("testdata", "motto.txt")
	if err := ioutil.WriteFile(fName, motto, 0777); err != nil {
		t.Errorf("Can't write test data %s, %s", fName, err)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, semver, fName); err != nil {
		t.Errorf("Can't attach %s %s, %s", semver, fName, err)
		t.FailNow()
	}
	semver = "v0.0.2"
	motto = []byte("Wowie Zowie!!")
	size = int64(len(motto))
	checksum = fmt.Sprintf("%x", md5.Sum(motto))
	fName = path.Join("testdata", "motto.txt")
	if err := ioutil.WriteFile(fName, motto, 0777); err != nil {
		t.Errorf("Can't write test data %s, %s", fName, err)
		t.FailNow()
	}
	// Check to make sure first version is correct
	if files, err := c.Attachments(keyName); err != nil {
		t.Errorf("Attachments (1) expected, %+v %s", files, err)
		t.FailNow()
	} else if len(files) != 1 {
		t.Errorf("Attachments (1) expected, %+v", files)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, semver, fName); err != nil {
		t.Errorf("Can't attach %s %s, %s", semver, fName, err)
		t.FailNow()
	}
	// We should stil have one attachment
	if files, err := c.Attachments(keyName); err != nil {
		t.Errorf("Attachments (1) expected, %+v %s", files, err)
		t.FailNow()
	} else if len(files) != 1 {
		t.Errorf("Attachments (1) expected, %+v", files)
		t.FailNow()
	}
	// Check JSON object
	jsonObject := map[string]interface{}{}
	err = c.Read(keyName, jsonObject, false)
	if err != nil {
		t.Errorf("Should be able to read %s, %s", keyName, err)
		t.FailNow()
	}
	// Make sure we have two semver items in Sizes, Checksums and
	// VersionHRefs
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == false {
		t.Errorf("Should be able to get attachment list %s", keyName)
		t.FailNow()
	}
	if len(attachmentList) != 1 {
		t.Errorf("Should have `1 attachments, got %d", len(attachmentList))
		t.FailNow()
	}
	if len(attachmentList[0].Checksums) != 2 {
		t.Errorf("Expected 2 checksums, got %d", len(attachmentList[0].Checksums))
		t.FailNow()
	}
	if len(attachmentList[0].Sizes) != 2 {
		t.Errorf("Expected 2 Sizes, got %d", len(attachmentList[0].Sizes))
		t.FailNow()
	}
	if len(attachmentList[0].VersionHRefs) != 2 {
		t.Errorf("Expected 2 VersionHRefs, got %d", len(attachmentList[0].VersionHRefs))
		t.FailNow()
	}

	// make sure we can marshal still
	if _, err := json.MarshalIndent(jsonObject, "", "    "); err != nil {
		t.Errorf("Could not marshal %s, %s", keyName, err)
		t.FailNow()
	}
}
