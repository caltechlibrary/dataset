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
	cName := path.Join("testout", "col3.ds")
	dsnURI := "sqlite://testout/col3.ds/collection.db"
	os.RemoveAll(cName)

	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Can't create collection %q (%s)", cName, err)
		t.FailNow()
	}
	c.SetVersioning("patch")
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Can't open collection %q (%s)", cName, err)
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
	if err := ioutil.WriteFile(path.Join("testout", "helloworld.txt"), buf, 0777); err != nil {
		t.Errorf("Can't create test helloworld.txt, %s", err)
		t.FailNow()
	}

	version := "0.0.1"
	motto := []byte("Wowie Zowie!")
	data := &Attachment{
		Name: "impressed.txt",
		Size: int64(len(motto)),
		Checksums: map[string]string{
			version: fmt.Sprintf("%x", md5.Sum(motto)),
		},
	}
	if _, err := json.MarshalIndent(data, "", "    "); err != nil {
		t.Errorf("marshal error %s", err)
		t.FailNow()
	}
	if err := ioutil.WriteFile(path.Join("testout", data.Name), motto, 0777); err != nil {
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
	if err := c.AttachFile(keyName, path.Join("testout", "helloworld.txt")); err != nil {
		t.Errorf("failed to add attachments to %s, %s", path.Join("testout", "helloworld.txt"), err)
		t.FailNow()
	}
	fPath, err := c.AttachmentPath(keyName, "helloworld.txt")
	if err != nil {
		t.Errorf("Should be able to get a path to attachment, %s", err)
		t.FailNow()
	}
	if fInfo, err := os.Stat(fPath); os.IsNotExist(err) {
		t.Errorf("Attachment not created %q", path.Join("testout", "helloworld.txt"))
		t.FailNow()
	} else if err != nil {
		t.Errorf("Stat error for %q, %s", fPath, err)
		t.FailNow()
	} else if fInfo.Size() == 0 {
		t.Errorf("Empty attachment, %q --> %q", fPath)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, path.Join("testout", data.Name)); err != nil {
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

	if err := c.Prune("freda", "helloworld.txt"); err != nil {
		t.Errorf("Delete one file, %s", err)
	}
	if files, err := c.Attachments("freda"); err != nil {
		t.Errorf("Attachments (1) expected, %+v %s", files, err)
		t.FailNow()
	} else {
		if err := c.Prune("freda", "impressed.txt"); err != nil {
			t.Errorf("Delete one file, %s", err)
		}
	}

	if err := c.Prune("freda"); err != nil {
		t.Errorf("Delete all attachmemts, %s", err)
		t.FailNow()
	}
	// Make sure files have been removed from collection
	docDir := path.Join("testout", "col3.ds", "attached", "fr", "ed", "a", version)
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
	version = "0.0.1"
	motto = []byte("Wowie Zowie")
	fName := path.Join("testout", "motto.txt")
	if err := ioutil.WriteFile(fName, motto, 0777); err != nil {
		t.Errorf("Can't write test data %s, %s", fName, err)
		t.FailNow()
	}
	if err := c.AttachFile(keyName, fName); err != nil {
		t.Errorf("Can't attach %s %s, %s", fName, err)
		t.FailNow()
	}
	version = "v0.0.2"
	motto = []byte("Wowie Zowie!!")
	size = int64(len(motto))
	checksum = fmt.Sprintf("%x", md5.Sum(motto))
	fName = path.Join("testout", "motto.txt")
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
	if err := c.AttachFile(keyName, fName); err != nil {
		t.Errorf("Can't attach %s %s, %s", version, fName, err)
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
	err = c.Read(keyName, jsonObject)
	if err != nil {
		t.Errorf("Should be able to read %s, %s", keyName, err)
		t.FailNow()
	}
	// Make sure we have two semver items in Sizes, Checksums and
	// VersionHRefs
	attachmentList, err := c.Attachments(keyName)
	if err != nil {
		t.Errorf("Should be able to get attachment list %s", keyName)
		t.FailNow()
	}
	for _, fName := range attachmentList {
		meta, err := c.AttachmentInfo(keyName, fName)
		if meta.Checksum == "" {
			t.Errorf("Expected a checksum for %q, %q", keyName, fName)
			t.FailNow()
		}
		if meta.Sizes == 0 {
			t.Errorf("Expected file size > 0 for %q, %q", keyName, fName)
			t.FailNow()
		}
		if meta.VersionHRefs == "" {
			t.Errorf("Expected VersionHRefs for %q, %q", keyName, fName)
			t.FailNow()
		}
	}
}
