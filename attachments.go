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
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/pairtree"
	"github.com/caltechlibrary/semver"
)

//
// Overview:
//
// As of v2 dataset attachments are held in a pairtree of their own.
// The reason is that SQL database support for directly storing
// binary blobs is limited so the file system remains a simple
// solution to storing attachments.
//
// Attachments can support versioned and unversioned collections.
// If the collection is unversioned then the attachment is placed
// in the pairtree directory formed by the collection's working path,
// "attachments", pairtree encoded JSON object key, filename.
// If the collection is versioned then the attachment is placed
// in the pairtree directory formed by the collection's working path,
// "attachments", pairtree encoded JSON object key, directory named
// for the filename being attached with the specific version of the
// file stored as a semver style version number.  The "current"
// semver is linked back to the unversioned directory as the filename.
//
// NOTE: the AttachVersionStream() func does not managed the symbolic
// link for the current version. That is managed by AttachStream() which
// handles but use case for attach an unversioned file or a versioned file.
//
// NOTE: The attachments methods do not lock. For the web service
// implementation of dataset (i.e. datasetd) the service must providing
// the necessary locking to avoid competing writes and deletes.
//

// Attachment is a structure for holding non-JSON content metadata
// you wish to store alongside a JSON document in a collection
// Attachments reside in a their own pairtree of the collection directory.
// (even when using a SQL store for the JSON document). The attachment
// metadata is read as needed from disk where the collection folder
// resides.
type Attachment struct {
	// Name is the filename and path to be used inside the generated tar file
	Name string `json:"name"`

	// Size remains to to help us migrate pre v0.0.61 collections.
	// It should reflect the last size added.
	Size int64 `json:"size"`

	// Sizes is the sizes associated with the version being attached
	Sizes map[string]int64 `json:"sizes"`

	// Current holds the semver to the last added version
	Version string `json:"version"`

	// Checksum, current implemented as a MD5 checksum for now
	// You should have one checksum per attached version.
	Checksums map[string]string `json:"checksums"`

	// HRef points at last attached version of the attached document
	// If you moved an object out of the pairtree it should be a URL.
	HRef string `json:"href"`

	// VersionHRefs is a map to all versions of the attached document
	// {
	//    "0.0.0": "... /photo.png",
	//    "0.0.1": "... /photo.png",
	//    "0.0.2": "... /photo.png"
	// }
	VersionHRefs map[string]string `json:"version_hrefs"`

	// Created a date string in RTC3339 format
	Created string `json:"created"`

	// Modified a date string in RFC3339 format
	Modified string `json:"modified"`

	// Metadata is a map for application specific metadata about attachments.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

//
// Private helper functions
//

// attachmentDir calculates a filepath's dir based on a collection
// and key.
func attachmentDir(c *Collection, key string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("collection isn't open")
	}
	workPath := c.workPath
	pairPath := pairtree.Encode(key)
	return path.Join(workPath, "attachments", pairPath), nil
}

// attachmentVersionDir calculates a filepath's dir based on a
// collection, key, version
func attachmentVersionDir(c *Collection, key string, filename string) (string, error) {
	if c == nil || c.workPath == "" {
		return "", fmt.Errorf("collection isn't open")
	}
	workPath := c.workPath
	pairPath := pairtree.Encode(key)
	return path.Join(workPath, "attachments", pairPath, "_", path.Base(filename)), nil
}

// Attachments returns a list of filenames for a key name in the collection
//
//	Example: "c" is a dataset collection previously opened,
//	"key" is a string.  The "key" is for a JSON document in
//	the collection. It returns an slice of filenames and err.
//
// ```
//
//	filenames, err := c.Attachments(key)
//	if err != nil {
//	   ...
//	}
//	// Print the names of the files attached to the JSON document
//	// referred to by "key".
//	for i, filename := ranges {
//	   fmt.Printf("key: %q, filename: %q", key, filename)
//	}
//
// ```
func (c *Collection) Attachments(key string) ([]string, error) {
	aPath, err := attachmentDir(c, key)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(aPath); os.IsNotExist(err) {
		return []string{}, nil
	}
	attachments := []string{}
	dir, err := os.ReadDir(aPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range dir {
		if entry != nil {
			filename := path.Base(entry.Name())
			if !entry.IsDir() {
				attachments = append(attachments, filename)
			}
		}
	}
	return attachments, nil
}

// AttachmentVersions returns a list of versions for an attached file
// to a JSON document in the collection.
//
//	Example: retrieve a list of versions of an attached file.
//	"key" is a key in the collection, filename is name of an
//	attached file for the JSON document referred to by key.
//
// ```
//
//	versions, err := c.AttachmentVersions(key, filename)
//	if err != nil {
//	   ...
//	}
//	for i, version := range versions {
//	   fmt.Printf("key: %q, filename: %q, version: %q", key, filename, version)
//	}
//
// ```
func (c *Collection) AttachmentVersions(key string, filename string) ([]string, error) {
	aPath, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return nil, err
	}
	names, err := os.ReadDir(aPath)
	if err != nil {
		return nil, err
	}
	versions := []string{}
	for _, entry := range names {
		version := entry.Name()
		if entry.IsDir() == false {
			versions = append(versions, path.Base(version))
		}
	}
	return versions, nil
}

// AttachStream is for attaching a non-JSON file via a io buffer.
// It requires the JSON document key, the filename and a io.Reader.
// It does not close the reader. If the collection is versioned then
// the document attached is automatically versioned per collection
// versioning setting.
//
//	Example: attach the file "report.pdf" to JSON document "123"
//	in an open collection.
//
// ```
//
//	key, filename := "123", "report.pdf"
//	buf, err := os.Open(filename)
//	if err != nil {
//	   ...
//	}
//	err := c.AttachStream(key, filename, buf)
//	if err != nil {
//	   ...
//	}
//	buf.Close()
//
// ```
func (c *Collection) AttachStream(key string, filename string, buf io.Reader) error {
	aDir, err := attachmentDir(c, key)
	if err != nil {
		return fmt.Errorf("Can't figure out attachment path for %q, %q, %s", key, filename, err)
	}

	// Create the attachment path if neccessary
	if _, err := os.Stat(aDir); os.IsNotExist(err) {
		if err := os.MkdirAll(aDir, 0775); err != nil {
			return err
		}
	}
	if c.Versioning == "" || c.Versioning == "none" {
		attachmentFilename := path.Join(aDir, path.Base(filename))
		out, err := os.Create(attachmentFilename)
		if err != nil {
			return fmt.Errorf("failed to create %q, %q, %s", key, filename, err)
		}
		defer out.Close()
		if _, err := io.Copy(out, buf); err != nil {
			return fmt.Errorf("failed to write %q, %q to stream, %s", key, filename, err)
		}
	} else {
		// Get version
		version := "0.0.0"
		versions, err := c.AttachmentVersions(key, filename)
		if err == nil && len(versions) > 0 {
			version = versions[len(versions)-1]
		}
		sv, err := semver.Parse([]byte(version))
		switch c.Versioning {
		case "major":
			sv.IncMajor()
		case "minor":
			sv.IncMinor()
		case "patch":
			sv.IncPatch()
		}
		version = strings.TrimPrefix(sv.String(), "v")
		vDir, err := attachmentVersionDir(c, key, path.Base(filename))
		if _, err := os.Stat(vDir); os.IsNotExist(err) {
			os.MkdirAll(vDir, 0775)
		}
		// "old" name
		linkTo := path.Join("_", path.Base(filename), version)
		// "new" name
		target := path.Join(aDir, path.Base(filename))
		err = c.AttachVersionStream(key, filename, version, buf)
		// NOTE: A link is used to the versioned file to save space
		// If a link exists we need to remove it, then create a new
		// link.
		if _, err := os.Lstat(target); err == nil {
			os.Remove(target)
		}
		if err := os.Symlink(linkTo, target); err != nil {
			return fmt.Errorf("failed to link attachment %q, %q, %q, %s", key, filename, version, err)
		}
	}
	return nil
}

// ```
//
//	key, filename := "123", "report.pdf"
//	err := c.AttachFile(key, filename)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) AttachFile(key string, filename string) error {
	buf, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer buf.Close()
	return c.AttachStream(key, filename, buf)
}

// AttachVersionStream is for attaching open a non-JSON file buffer
// (via an io.Reader) to a specific version of a file. If attached
// file exists it is replaced.
//
//	Example: attach the file "report.pdf", version "0.0.3" to
//	JSON document "123" in an open collection.
//
// ```
//
//	key, filename, version := "123", "helloworld.txt", "0.0.3"
//	buf, err := os.Open(filename)
//	if err != nil {
//	   ...
//	}
//	err := c.AttachVersionStream(key, filename, version, buf)
//	if err != nil {
//	   ...
//	}
//	buf.Close()
//
// ```
func (c *Collection) AttachVersionStream(key string, filename string, version string, buf io.Reader) error {
	vDir, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return fmt.Errorf("failed to calculate version path, %s", err)
	}
	if _, err := os.Stat(vDir); os.IsNotExist(err) {
		os.MkdirAll(vDir, 0775)
	}
	vPath := path.Join(vDir, version)
	out, err := os.Create(vPath)
	if err != nil {
		return fmt.Errorf("failed to create versioned attachment, %q, %q, %q, %s", key, filename, version, err)
	}
	defer out.Close()
	if _, err := io.Copy(out, buf); err != nil {
		return fmt.Errorf("failed to write %q, %q, %q to output stream, %s", key, filename, version, err)
	}
	return nil
}

// AttachVersionFile attaches a file to a JSON document in the collection.
// This does NOT increment the version number of attachment(s). It is used
// to explicitly replace a attached version of a file. It does not update
// the symbolic link to the "current" attachment.
//
// ```
//
//	key, filename, version := "123", "report.pdf", "0.0.3"
//	err := c.AttachVersionFile(key, filename, version)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) AttachVersionFile(key string, filename string, version string) error {
	buf, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer buf.Close()
	return c.AttachVersionStream(key, filename, version, buf)
}

// AttachmentPath takes a key and filename and returns the path file
// system path to the attached file (if found). For versioned collections
// this is the path the symbolic link for the "current" version.
//
// ```
//
//	key, filename := "123", "report.pdf"
//	docPath, err := c.AttachmentPath(key, filename)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) AttachmentPath(key string, filename string) (string, error) {
	aDir, err := attachmentDir(c, key)
	if err != nil {
		return "", err
	}
	aPath := path.Join(aDir, path.Base(filename))
	if _, err := os.Lstat(aPath); err != nil {
		return "", err

	}
	return aPath, nil
}

// AttachmentVersionPath takes a key, filename and semver returning
// the path to the attached versioned file (if found).
//
// ```
//
//	key, filename, version := "123", "report.pdf", "0.0.3"
//	docPath, err := c.AttachmentVersionPath(key, filename, version)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) AttachmentVersionPath(key string, filename string, version string) (string, error) {
	vDir, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return "", err
	}
	vPath := path.Join(vDir, version)
	if _, err := os.Stat(vPath); err != nil {
		return "", err

	}
	return vPath, nil
}

// RetrieveStream takes a key and filename then returns an io.Reader,
// and error. If the collection is versioned then the stream is for the
// "current" version of the attached file.
//
// ```
//
//	key, filename := "123", "report.pdf"
//	src := []byte{}
//	buf := bytes.NewBuffer(src)
//	err := c.Retrieve(key, filename, buf)
//	if err != nil {
//	   ...
//	}
//	ioutil.WriteFile(filename, src, 0664)
//
// ```
func (c *Collection) RetrieveStream(key string, filename string, out io.Writer) error {
	aDir, err := attachmentDir(c, key)
	if err != nil {
		return err
	}
	aPath := path.Join(aDir, path.Base(filename))
	in, err := os.Open(aPath)
	if err != nil {
		return err
	}
	defer in.Close()
	size, err := io.Copy(out, in)
	if err != nil {
		return err
	}
	if size == 0 {
		return fmt.Errorf("zero bytes copied")
	}
	return err
}

// RetrieveFile retrieves a file attached to a JSON document in the
// collection.
//
// ```
//
//	key, filename := "123", "report.pdf"
//	src, err := c.RetrieveFile(key, filename)
//	if err != nil {
//	   ...
//	}
//	err = ioutil.WriteFile(filename, src, 0664)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) RetrieveFile(key string, filename string) ([]byte, error) {
	src := []byte{}
	buf := bytes.NewBuffer(src)
	err := c.RetrieveStream(key, filename, buf)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// RetrieveVersionStream takes a key, filename and version then
// returns an io.Reader and error.
//
// ```
//
//	key, filename, version := "123", "helloworld.txt", "0.0.3"
//	src := []byte{}
//	buf := bytes.NewBuffer(src)
//	err := c.RetrieveVersion(key, filename, version, buf)
//	if err != nil {
//	   ...
//	}
//	ioutil.WriteFile(filename + "_" + version, src, 0664)
//
// ```
func (c *Collection) RetrieveVersionStream(key string, filename string, version string, buf io.Writer) error {
	vDir, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return err
	}
	vPath := path.Join(vDir, version)
	in, err := os.Open(vPath)
	if err != nil {
		return err
	}
	defer in.Close()
	if _, err := io.Copy(buf, in); err != nil {
		return err
	}
	return nil
}

// RetrieveVersionFile retrieves a file version attached to a JSON
// document in the collection.
//
// ```
//
//	key, filename, version := "123", "report.pdf", "0.0.3"
//	src, err := c.RetrieveVersionFile(key, filename, version)
//	if err != nil  {
//	   ...
//	}
//	err = ioutil.WriteFile(filename + "_" + version, src, 0664)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) RetrieveVersionFile(key string, filename string, version string) ([]byte, error) {
	src := []byte{}
	buf := bytes.NewBuffer(src)
	err := c.RetrieveVersionStream(key, filename, version, buf)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// Prune removes a an attached document from the JSON record given a key and
// filename. NOTE: In versioned collections this include removing all
// versions of the attached document.
//
// ```
//
//	key, filename := "123", "report.pdf"
//	err := c.Prune(key, filename)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) Prune(key string, filename string) error {
	vDir, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return err
	}
	if _, err := os.Stat(vDir); err == nil {
		if err := os.RemoveAll(vDir); err != nil {
			return err
		}
	}
	aDir, err := attachmentDir(c, key)
	if err != nil {
		return err
	}
	aPath := path.Join(aDir, path.Base(filename))
	if _, err := os.Lstat(aPath); err == nil {
		if err := os.RemoveAll(aPath); err != nil {
			return err
		}
	}
	return nil
}

// PruneVersion removes an attached version of a document.
//
// ```
//
//	key, filename, version := "123", "report.pdf, "0.0.3"
//	err := c.PruneVersion(key, filename, version)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) PruneVersion(key string, filename string, version string) error {
	vDir, err := attachmentVersionDir(c, key, filename)
	if err != nil {
		return err
	}
	vPath := path.Join(vDir, version)
	return os.RemoveAll(vPath)
}

// PruneAll removes attachments from a JSON record in the collection.
// When the collection is versioned it removes all versions of all too.
//
// ```
//
//	key := "123"
//	err := c.PruneAll(key)
//	if err != nil {
//	   ...
//	}
//
// ```
func (c *Collection) PruneAll(key string) error {
	if c == nil {
		return fmt.Errorf("collection isn't open")
	}
	workPath := c.workPath
	pairPath := pairtree.Encode(key)
	vDir := path.Join(workPath, "attachments", pairPath)
	return os.RemoveAll(vDir)
}
