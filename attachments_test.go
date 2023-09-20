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
package dataset

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	// Caltech Library packages
	"github.com/caltechlibrary/pairtree"
)

func TestAttachmentsPTStore(t *testing.T) {
	cName := path.Join("testout", "attachment_test_pt.ds")
	dsnURI := ""
	os.RemoveAll(cName)

	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Can't create collection %q (%s)", cName, err)
		t.FailNow()
	}
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Can't open collection %q (%s)", cName, err)
		t.FailNow()
	}

	// Create some temp text files to attach.
	key := "freda"
	version := "0.0.1"
	filename := path.Join("testout", "helloworld.txt")
	src := []byte("Hello World")
	size := int64(len(src))
	if size != int64(11) {
		t.Errorf("Expected 'Hello World' to be size 11, got %d", size)
	}
	expectedChecksum := "b10a8db164e0754105b7a99be72e3fe5"
	checksum := fmt.Sprintf("%x", md5.Sum(src))
	if strings.Compare(checksum, expectedChecksum) != 0 {
		t.Errorf("Expected a checksum %s, got %q", expectedChecksum, checksum)
	}
	if err := os.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("Can't create test helloworld.txt, %s", err)
		t.FailNow()
	}

	motto := []byte("Wowie Zowie!")
	data := &Attachment{
		Name: "impressed.txt",
		Size: int64(len(motto)),
		Checksums: map[string]string{
			version: fmt.Sprintf("%x", md5.Sum(motto)),
		},
	}
	if _, err := JSONMarshalIndent(data, "", "    "); err != nil {
		t.Errorf("marshal error %s", err)
		t.FailNow()
	}
	if err := os.WriteFile(path.Join("testout", data.Name), motto, 0664); err != nil {
		t.Errorf("Can't create test %q, %s", data.Name, err)
		t.FailNow()
	}

	record := map[string]interface{}{
		"name":  "freda",
		"motto": "it's all about what you sense when you've have sense to sense it",
	}
	if err := c.Create(key, record); err != nil {
		t.Errorf("failed to create %s in %s, %s", key, c.Name, err)
		t.FailNow()
	}

	//NOTE: Need to open the file as io.Reader. We use stream
	// buffers to limit the in-memory demands of copying the file.
	for _, filename := range []string{"helloworld.txt", "impressed.txt"} {
		filename = path.Join("testout", filename)
		buf, err := os.Open(filename)
		if err != nil {
			t.Errorf("failed to open test file %q as stream, %s", filename, err)
			t.FailNow()
		}
		if err := c.AttachStream(key, filename, buf); err != nil {
			buf.Close()
			t.Errorf("failed to add attachment to %q, %q, %s", key, filename, err)
			t.FailNow()
		}

		fPath, err := c.AttachmentPath(key, path.Base(filename))
		if err != nil {
			t.Errorf("Should be able to get a path to attachment %q, %s", filename, err)
			t.FailNow()
		}
		if fInfo, err := os.Lstat(fPath); os.IsNotExist(err) {
			t.Errorf("Attachment not created %q", filename)
			t.FailNow()
		} else if err != nil {
			t.Errorf("Stat error for %q, %s", fPath, err)
			t.FailNow()
		} else if fInfo.Size() == 0 {
			t.Errorf("Empty attachment, %q --> %q", filename, fPath)
			t.FailNow()
		}
	}

	filenames, err := c.Attachments(key)
	if err != nil {
		t.Errorf("can't list attachments for %s in %s, %s", key, c.Name, err)
		t.FailNow()
	}
	if len(filenames) != 2 {
		t.Errorf("Expected two file unversioned attachments, %+v", filenames)
		t.FailNow()
	}
	for _, filename := range filenames {
		if !(filename == "helloworld.txt" || filename == "impressed.txt") {
			t.Errorf("Unexpected filename in list, %s", filename)
			t.FailNow()
		}
		fPath, err := c.AttachmentPath(key, filename)
		if err != nil {
			t.Errorf("expected path c.AttachmentPath(%q, %q), %s", key, filename, err)
			continue
		}
		pairPath := pairtree.Encode(key)
		expectedPath, _ := filepath.Abs(path.Join("testout", "attachment_test_pt.ds", "attachments", pairPath, filename))
		if fPath != expectedPath {
			t.Errorf("expected %q, got %q", expectedPath, fPath)
			t.FailNow()
		}
		if fInfo, err := os.Lstat(fPath); err != nil {
			t.Errorf("Expected os.Stat(%q) for %q, %q, %q, %s", fPath, key, filename, version, err)
			continue
		} else {
			sName := fInfo.Name()
			size := fInfo.Size()
			if sName != filename {
				t.Errorf("expected %q, got %q for os.Stat(%q)", filename, sName, fPath)
			}
			if size == 0 {
				t.Errorf("expected size > 0, got %d for os.Stat(%q)", size, fPath)
			}
		}
	}
	key = "freda"
	if err := c.Prune(key, "helloworld.txt"); err != nil {
		t.Errorf("Delete helloworld.txt attachment, %s", err)
		t.FailNow()
	}
	names, err := c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) error, %s", key, err)
		t.FailNow()
	}
	if len(names) != 1 {
		t.Errorf("expected on attach in list, got %d, %+v", len(names), names)
	}
	if err := c.Prune("freda", "impressed.txt"); err != nil {
		t.Errorf("Delete impressed.txt attachment, %s", err)
		t.FailNow()
	}
	names, err = c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) expected, %s", key, err)
		t.FailNow()
	}
	if len(names) != 0 {
		t.Errorf("expected on attach in list, got %d, %+v", len(names), names)
	}

	if err := c.PruneAll("freda"); err != nil {
		t.Errorf("Delete all attachmemts, %s", err)
		t.FailNow()
	}

	// Make sure files have been removed from collection
	docDir := path.Join("testout", "attachment_test_pt.ds", "attachments", "fr", "ed", "a", "_")
	for _, filename := range []string{"impressed.txt", "helloworld.txt"} {
		aPath := path.Join(docDir, filename)
		if _, err := os.Stat(aPath); err == nil {
			t.Errorf("Should have deleted %s", aPath)
			t.FailNow()
		}
	}
	docDir = path.Join("testout", "attachment_test_pt.ds", "attachments", "fr", "ed", "a")
	for _, filename := range []string{"impressed.txt", "helloworld.txt"} {
		aPath := path.Join(docDir, filename)
		if _, err := os.Stat(aPath); err == nil {
			t.Errorf("Should have deleted %s", aPath)
			t.FailNow()
		}
	}
	

}

func TestAttachmentsSQLStore(t *testing.T) {
	cName := path.Join("testout", "attachment_test.ds")
	dsnURI := "sqlite://testout/attachment_test.ds/collection.db"
	os.RemoveAll(cName)

	c, err := Init(cName, dsnURI)
	if err != nil {
		t.Errorf("Can't create collection %q (%s)", cName, err)
		t.FailNow()
	}
	c.Close()
	c, err = Open(cName)
	if err != nil {
		t.Errorf("Can't open collection %q (%s)", cName, err)
		t.FailNow()
	}

	// Create some temp text files to attach.
	key := "freda"
	version := "0.0.1"
	filename := path.Join("testout", "helloworld.txt")
	src := []byte("Hello World")
	size := int64(len(src))
	if size != int64(11) {
		t.Errorf("Expected 'Hello World' to be size 11, got %d", size)
	}
	expectedChecksum := "b10a8db164e0754105b7a99be72e3fe5"
	checksum := fmt.Sprintf("%x", md5.Sum(src))
	if strings.Compare(checksum, expectedChecksum) != 0 {
		t.Errorf("Expected a checksum %s, got %q", expectedChecksum, checksum)
	}
	if err := os.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("Can't create test helloworld.txt, %s", err)
		t.FailNow()
	}

	motto := []byte("Wowie Zowie!")
	data := &Attachment{
		Name: "impressed.txt",
		Size: int64(len(motto)),
		Checksums: map[string]string{
			version: fmt.Sprintf("%x", md5.Sum(motto)),
		},
	}
	if _, err := JSONMarshalIndent(data, "", "    "); err != nil {
		t.Errorf("marshal error %s", err)
		t.FailNow()
	}
	if err := os.WriteFile(path.Join("testout", data.Name), motto, 0664); err != nil {
		t.Errorf("Can't create test %q, %s", data.Name, err)
		t.FailNow()
	}

	record := map[string]interface{}{
		"name":  "freda",
		"motto": "it's all about what you sense when you've have sense to sense it",
	}
	if err := c.Create(key, record); err != nil {
		t.Errorf("failed to create %s in %s, %s", key, c.Name, err)
		t.FailNow()
	}

	//NOTE: Need to open the file as io.Reader. We use stream
	// buffers to limit the in-memory demands of copying the file.
	for _, filename := range []string{"helloworld.txt", "impressed.txt"} {
		filename = path.Join("testout", filename)
		buf, err := os.Open(filename)
		if err != nil {
			t.Errorf("failed to open test file %q as stream, %s", filename, err)
			t.FailNow()
		}
		if err := c.AttachStream(key, filename, buf); err != nil {
			buf.Close()
			t.Errorf("failed to add attachment to %q, %q, %s", key, filename, err)
			t.FailNow()
		}

		fPath, err := c.AttachmentPath(key, path.Base(filename))
		if err != nil {
			t.Errorf("Should be able to get a path to attachment %q, %s", filename, err)
			t.FailNow()
		}
		if fInfo, err := os.Lstat(fPath); os.IsNotExist(err) {
			t.Errorf("Attachment not created %q", filename)
			t.FailNow()
		} else if err != nil {
			t.Errorf("Stat error for %q, %s", fPath, err)
			t.FailNow()
		} else if fInfo.Size() == 0 {
			t.Errorf("Empty attachment, %q --> %q", filename, fPath)
			t.FailNow()
		}
	}

	filenames, err := c.Attachments(key)
	if err != nil {
		t.Errorf("can't list attachments for %s in %s, %s", key, c.Name, err)
		t.FailNow()
	}
	if len(filenames) != 2 {
		t.Errorf("Expected two file unversioned attachments, %+v", filenames)
		t.FailNow()
	}
	for _, filename := range filenames {
		if !(filename == "helloworld.txt" || filename == "impressed.txt") {
			t.Errorf("Unexpected filename in list, %s", filename)
			t.FailNow()
		}
		fPath, err := c.AttachmentPath(key, filename)
		if err != nil {
			t.Errorf("expected path c.AttachmentPath(%q, %q), %s", key, filename, err)
			continue
		}
		pairPath := pairtree.Encode(key)
		expectedPath, _ := filepath.Abs(path.Join("testout", "attachment_test.ds", "attachments", pairPath, filename))
		if fPath != expectedPath {
			t.Errorf("expected %q, got %q", expectedPath, fPath)
			t.FailNow()
		}
		if fInfo, err := os.Lstat(fPath); err != nil {
			t.Errorf("Expected os.Stat(%q) for %q, %q, %q, %s", fPath, key, filename, version, err)
			continue
		} else {
			sName := fInfo.Name()
			size := fInfo.Size()
			if sName != filename {
				t.Errorf("expected %q, got %q for os.Stat(%q)", filename, sName, fPath)
			}
			if size == 0 {
				t.Errorf("expected size > 0, got %d for os.Stat(%q)", size, fPath)
			}
		}
	}
	key = "freda"
	if err := c.Prune(key, "helloworld.txt"); err != nil {
		t.Errorf("Delete helloworld.txt attachment, %s", err)
		t.FailNow()
	}
	names, err := c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) error, %s", key, err)
		t.FailNow()
	}
	if len(names) != 1 {
		t.Errorf("expected on attach in list, got %d, %+v", len(names), names)
	}
	if err := c.Prune("freda", "impressed.txt"); err != nil {
		t.Errorf("Delete impressed.txt attachment, %s", err)
		t.FailNow()
	}
	names, err = c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) expected, %s", key, err)
		t.FailNow()
	}
	if len(names) != 0 {
		t.Errorf("expected on attach in list, got %d, %+v", len(names), names)
	}

	if err := c.PruneAll("freda"); err != nil {
		t.Errorf("Delete all attachmemts, %s", err)
		t.FailNow()
	}

	// Make sure files have been removed from collection
	docDir := path.Join("testout", "attachment_test.ds", "attachments", "fr", "ed", "a", "_")
	for _, filename := range []string{"impressed.txt", "helloworld.txt"} {
		aPath := path.Join(docDir, filename)
		if _, err := os.Stat(aPath); err == nil {
			t.Errorf("Should have deleted %s", aPath)
			t.FailNow()
		}
	}

}

func TestVersionedAttachments(t *testing.T) {
	cName := path.Join("testout", "v_attachment_test.ds")
	dsnURI := "sqlite://testout/v_attachment_test.ds/collection.db"
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
	key := "freda"
	version := "0.0.1"
	filename := path.Join("testout", "helloworld.txt")
	src := []byte("Hello World")
	size := int64(len(src))
	if size != int64(11) {
		t.Errorf("Expected 'Hello World' to be size 11, got %d", size)
	}
	expectedChecksum := "b10a8db164e0754105b7a99be72e3fe5"
	checksum := fmt.Sprintf("%x", md5.Sum(src))
	if strings.Compare(checksum, expectedChecksum) != 0 {
		t.Errorf("Expected a checksum %s, got %q", expectedChecksum, checksum)
	}
	if err := os.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("Can't create test helloworld.txt, %s", err)
		t.FailNow()
	}

	motto := []byte("Wowie Zowie!")
	data := &Attachment{
		Name: "impressed.txt",
		Size: int64(len(motto)),
		Checksums: map[string]string{
			version: fmt.Sprintf("%x", md5.Sum(motto)),
		},
	}
	if _, err := JSONMarshalIndent(data, "", "    "); err != nil {
		t.Errorf("marshal error %s", err)
		t.FailNow()
	}
	if err := os.WriteFile(path.Join("testout", data.Name), motto, 0664); err != nil {
		t.Errorf("Can't create test %q, %s", data.Name, err)
		t.FailNow()
	}

	record := map[string]interface{}{
		"name":  "freda",
		"motto": "it's all about what you sense when you've have sense to sense it",
	}
	if err := c.Create(key, record); err != nil {
		t.Errorf("failed to create %s in %s, %s", key, c.Name, err)
		t.FailNow()
	}

	//NOTE: Need to open the file as io.Reader. We use stream
	// buffers to limit the in-memory demands of copying the file.
	for _, filename := range []string{"helloworld.txt", "impressed.txt"} {
		filename = path.Join("testout", filename)
		buf, err := os.Open(filename)
		if err != nil {
			t.Errorf("failed to open test file %q as stream, %s", filename, err)
			t.FailNow()
		}
		if err := c.AttachStream(key, filename, buf); err != nil {
			buf.Close()
			t.Errorf("failed to add attachment to %q, %q, %s", key, filename, err)
			t.FailNow()
		}

		fPath, err := c.AttachmentPath(key, path.Base(filename))
		if err != nil {
			t.Errorf("Should be able to get a path to attachment %q, %s", filename, err)
			t.FailNow()
		}
		if fInfo, err := os.Lstat(fPath); os.IsNotExist(err) {
			t.Errorf("Attachment not created %q", filename)
			t.FailNow()
		} else if err != nil {
			t.Errorf("Stat error for %q, %s", fPath, err)
			t.FailNow()
		} else if fInfo.Size() == 0 {
			t.Errorf("Empty attachment, %q --> %q", filename, fPath)
			t.FailNow()
		}
	}

	if filenames, err := c.Attachments(key); err != nil {
		t.Errorf("can't list attachments for %s in %s, %s", key, c.Name, err)
		t.FailNow()
	} else {
		if len(filenames) != 2 {
			t.Errorf("Expected one file attached, %+v", filenames)
			t.FailNow()
		}
		for _, filename := range filenames {
			if !(filename == "helloworld.txt" || filename == "impressed.txt") {
				t.Errorf("Unexpected filename in list, %s", filename)
			} else {
				versions, err := c.AttachmentVersions(key, filename)
				if err != nil {
					t.Errorf("Expected c.AttachmentVersions(%q, %q), %s", key, filename, err)
				}
				for _, version := range versions {
					fPath, err := c.AttachmentVersionPath(key, filename, version)
					if err != nil {
						t.Errorf("expected path c.AttachmentVersionPath(%q, %q), %s", key, filename, err)
						continue
					}
					if fInfo, err := os.Stat(fPath); err != nil {
						t.Errorf("Expected os.Stat(%q) for %q, %q, %q, %s", fPath, key, filename, version, err)
						continue
					} else {
						sName := fInfo.Name()
						size := fInfo.Size()
						if sName != version {
							t.Errorf("expected %q, got %q for os.Stat(%q)", filename, sName, fPath)
						}
						if size == 0 {
							t.Errorf("expected size > 0, got %d for os.Stat(%q)", size, fPath)
						}
					}
				}
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

	if err := c.PruneAll("freda"); err != nil {
		t.Errorf("Delete all attachmemts, %s", err)
		t.FailNow()
	}
	// Make sure files have been removed from collection
	docDir := path.Join("testout", "v_attachment_test.ds", "attachments", "fr", "ed", "a", "_")
	for _, filename := range []string{"impressed.txt", "helloworld.txt"} {
		aPath := path.Join(docDir, filename)
		if _, err := os.Stat(aPath); err == nil {
			t.Errorf("Should have deleted %s", aPath)
			t.FailNow()
		}
	}

	//
	// Now lets tests multiple versions of an attachment using both
	// AttachFile and AttachVersionStream.
	//
	key = "freda"
	filename = path.Join("testout", "motto.txt")
	motto = []byte("Wowie Zowie version 0.0.1\n")
	if err := os.WriteFile(filename, motto, 0664); err != nil {
		t.Errorf("Can't write test data %s, %s", filename, err)
		t.FailNow()
	}

	// Make first version of attachment
	if err := c.AttachFile(key, filename); err != nil {
		t.Errorf("Can't attach %s %s, %s", key, filename, err)
		t.FailNow()
	}
	names, err := c.Attachments(key)
	if err != nil {
		t.Errorf("Should see attachments for %q, %s", key, err)
		t.FailNow()
	}
	// Check to make sure first version is correct
	versions, err := c.AttachmentVersions(key, path.Base(filename))
	if err != nil {
		t.Errorf("c.AttachmentVersions(%q, %q), one expected, %+v %s", key, filename, versions, err)
		t.FailNow()
	}
	if len(versions) != 1 {
		t.Errorf("Attachments versions (1) expected, %+v", versions)
		t.FailNow()
	}

	motto = []byte("Wowie Zowie!! version 0.0.2\n")
	if err := os.WriteFile(filename, motto, 0664); err != nil {
		t.Errorf("Can't write test data %s, %s", filename, err)
	}

	// Make second version of attachment, 0.0.2.
	if err := c.AttachFile(key, filename); err != nil {
		t.Errorf("Can't attach %s %s, %s", key, filename, err)
	}

	// Check to make sure first version is correct
	versions, err = c.AttachmentVersions(key, path.Base(filename))
	if err != nil {
		t.Errorf("c.AttachmentVersions(%q, %q), one expected, %+v %s", key, filename, versions, err)
	}
	if len(versions) != 2 {
		t.Errorf("Attachments versions (2) expected, %+v", versions)
	}

	version = "0.0.1"
	motto = []byte("Wowie Zowie! version 0.0.1-rc2\n")
	if err := os.WriteFile(filename, motto, 0664); err != nil {
		t.Errorf("Can't write test data %s, %s", filename, err)
		t.FailNow()
	}
	if err := c.AttachVersionFile(key, filename, version); err != nil {
		t.Errorf("Can't attach %q %q %q, %s", key, filename, version, err)
	}

	// We should still have one attachment
	names, err = c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) one expected, %+v %s", key, names, err)
	}
	if len(names) != 1 {
		t.Errorf("Attachments (1) expected, %+v", names)
	}
	versions, err = c.AttachmentVersions(key, filename)
	if err != nil {
		t.Errorf("c.AttachmentVersions(%q, %q), three expected, %+v %s", key, filename, versions, err)
	}
	if len(versions) != 2 {
		t.Errorf("Attachments versions (2) expected, %+v", versions)
	}
	// The maximum version should remain 0.0.2
	expectedVersion := "0.0.2"
	gotVersion := versions[len(versions)-1]
	if expectedVersion != gotVersion {
		t.Errorf("expected version no %q, got %q", expectedVersion, gotVersion)
	}

	// Now add a third version via AttachFile()
	motto = []byte("Wowie Zowie!!! this is verion 3\n")
	if err := os.WriteFile(filename, motto, 0664); err != nil {
		t.Errorf("Can't write test data %s, %s", filename, err)
	}
	if err := c.AttachFile(key, filename); err != nil {
		t.Errorf("Can't attach %s %s, %s", key, filename, err)
	}
	names, err = c.Attachments(key)
	if err != nil {
		t.Errorf("c.Attachments(%q) one expected, %+v %s", key, names, err)
	}
	if len(names) != 1 {
		t.Errorf("Attachments (1) expected, %+v", names)
	}
	versions, err = c.AttachmentVersions(key, filename)
	if err != nil {
		t.Errorf("c.AttachmentVersions(%q, %q), three expected, %+v %s", key, filename, versions, err)
	}
	if len(versions) != 3 {
		t.Errorf("Attachments versions (3) expected, %+v", versions)
	}
	expectedVersion = "0.0.3"
	gotVersion = versions[len(versions)-1]
	if expectedVersion != versions[len(versions)-1] {
		t.Errorf("expected version no %q, got %q", expectedVersion, versions[len(versions)-1])
	}
}
