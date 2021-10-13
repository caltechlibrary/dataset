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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

// Attachment is a structure for holding non-JSON content metadata
// you wish to store alongside a JSON document in a collection
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

	// HRef points at last attached version of the attached document, e.g. v0.0.0/photo.png
	// If you moved an object out of the pairtree it should be a URL.
	HRef string `json:"href"`

	// VersionHRefs is a map to all versions of the attached document
	// {
	//    "v0.0.0": "... /photo.png",
	//    "v0.0.1": "... /photo.png",
	//    "v0.0.2": "... /photo.png"
	// }
	VersionHRefs map[string]string `json:"version_hrefs"`

	// Created a date string in RTC3339 format
	Created string `json:"created"`

	// Modified a date string in RFC3339 format
	Modified string `json:"modified"`

	// Metadata is a map for application specific metadata about attachments.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// getAttachmentList takes a JSON objects, pulls our "_Attachments" and reads them into
// and array of Attachment. Returns true if we have a list, false otherwise
func getAttachmentList(jsonObject map[string]interface{}) ([]*Attachment, bool) {
	attachmentList := []*Attachment{}
	attachments, ok := jsonObject["_Attachments"]
	if ok == false {
		return attachmentList, false
	}
	for _, obj := range attachments.([]interface{}) {
		attachment := new(Attachment)
		m := obj.(map[string]interface{})
		if name, ok := m["name"]; ok == true {
			attachment.Name = name.(string)
		}
		if val, ok := m["size"]; ok == true {
			if size, err := val.(json.Number).Int64(); err == nil {
				attachment.Size = size
			}
		}
		if sizes, ok := m["sizes"]; ok == true {
			m1 := sizes.(map[string]interface{})
			attachment.Sizes = make(map[string]int64)
			for k, v := range m1 {
				if size, err := v.(json.Number).Int64(); err == nil {
					attachment.Sizes[k] = size
				}
			}
		}
		// popagate our optional metadata
		if metadata, ok := m["metadata"]; ok == true {
			m2 := metadata.(map[string]interface{})
			attachment.Metadata = make(map[string]interface{})
			for k, v := range m2 {
				attachment.Metadata[k] = v
			}
		}
		if checksums, ok := m["checksums"]; ok == true {
			m3 := checksums.(map[string]interface{})
			attachment.Checksums = make(map[string]string)
			for k, v := range m3 {
				attachment.Checksums[k] = v.(string)
			}
		}
		if created, ok := m["created"]; ok == true {
			attachment.Created = created.(string)
		}
		if modified, ok := m["modified"]; ok == true {
			attachment.Modified = modified.(string)
		}
		if href, ok := m["href"]; ok == true {
			attachment.HRef = href.(string)
		}
		if version, ok := m["version"]; ok == true {
			attachment.Version = version.(string)
		}
		if versionHRefs, ok := m["version_hrefs"]; ok == true {
			m4 := versionHRefs.(map[string]interface{})
			attachment.VersionHRefs = make(map[string]string)
			for k, v := range m4 {
				attachment.VersionHRefs[k] = v.(string)
			}
		}
		attachmentList = append(attachmentList, attachment)
	}
	return attachmentList, true
}

func updateAttachmentList(attachmentList []*Attachment, newObj *Attachment) []*Attachment {
	isReplacement := false
	for i, oldObj := range attachmentList {
		if oldObj.Name == newObj.Name {
			isReplacement = true
			attachmentList[i] = newObj
			break
		}
	}
	// We append to the end of the list if we aren't replacinga object
	if !isReplacement {
		attachmentList = append(attachmentList, newObj)
	}
	return attachmentList
}

// AttachStream is for attaching open a non-JSON file buffer (via an io.Reader).
func (c *Collection) AttachStream(keyName, semver, fullName string, buf io.Reader) error {
	if c.KeyExists(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	if semver == "" {
		// We use version v0.0.0 for "unversioned" attachments.
		semver = "v0.0.0"
	}
	// Normalize fName to basename from fullName to be safe.
	fName := path.Base(fullName)

	// Read in JSON object and metadata objects.
	jsonObject := map[string]interface{}{}
	attachmentObject := &Attachment{}

	if err := c.Read(keyName, jsonObject, false); err != nil {
		return fmt.Errorf("Can't read %q, aborting, %s", keyName, err)
	}
	// This is the full path to the JSON Object document
	docPath, err := c.DocPath(keyName)
	if err != nil {
		return fmt.Errorf("Can't find document path %q, aborting, %s", keyName, err)
	}
	// This is JSON object's directory.
	docDir := path.Dir(docPath)

	// Now we're ready to get our attachment list.
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == true {
		for _, obj := range attachmentList {
			if obj.Name == fName {
				attachmentObject = obj
				break
			}
		}
	}

	// Update the metadata
	// Read in the attachment so I can compute the checksum as well as size.
	checksum, err := calcChecksum(fullName)
	if err != nil {
		return fmt.Errorf("Can't calc checksum, %q, %s", fullName, err)
	}
	attachmentObject.Name = fName
	// Add/update our version href
	attachmentObject.HRef = path.Join(docDir, semver, fName)
	// Write out attached filename and calc size and checksum
	err = os.MkdirAll(path.Dir(attachmentObject.HRef), 0777)

	if err != nil {
		return err
	}
	w, err := os.Create(attachmentObject.HRef)
	if err != nil {
		return fmt.Errorf("Can't create %q, %s", attachmentObject.HRef, err)
	}
	defer w.Close()
	l, err := io.Copy(w, buf)
	if err != nil {
		return fmt.Errorf("Can't write %q, %s", attachmentObject.HRef, err)
	}
	attachmentObject.Size = l
	attachmentObject.Version = semver
	if attachmentObject.Sizes == nil {
		attachmentObject.Sizes = make(map[string]int64)
	}
	attachmentObject.Sizes[semver] = l
	// Compute Checksum with md5 and store as a string
	if attachmentObject.Checksums == nil {
		attachmentObject.Checksums = make(map[string]string)
	}
	attachmentObject.Checksums[semver] = checksum
	// We need to make the semver directory if necessary
	if attachmentObject.VersionHRefs == nil {
		attachmentObject.VersionHRefs = make(map[string]string)
	}
	attachmentObject.VersionHRefs[semver] = attachmentObject.HRef
	now := time.Now()
	if attachmentObject.Created == "" {
		attachmentObject.Created = now.Format(time.RFC3339)
	}
	attachmentObject.Modified = now.Format(time.RFC3339)

	jsonObject["_Attachments"] = updateAttachmentList(attachmentList, attachmentObject)

	// Write out updated JSON Object and return any error
	err = c.Update(keyName, jsonObject)
	return err
}

// AttachFile is for attaching a single non-JSON document to a
// dataset record. It will replace ANY existing attached content
// with the same semver and basename.
func (c *Collection) AttachFile(keyName string, semver string, fullName string) error {
	// Normalize fName to basename of srcName
	dstName := path.Base(fullName)
	return c.AttachFileAs(keyName, semver, dstName, fullName)
}

// AttachFileAs is for attaching a single non-JSON document to a
// dataset record with a specific attachment name. It will replace ANY
// existing attached content with the same semver and destintation name.
func (c *Collection) AttachFileAs(keyName string, semver string, dstName string, srcName string) error {
	if c.KeyExists(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	if semver == "" {
		// We use version v0.0.0 for "unversioned" attachments.
		semver = "v0.0.0"
	}

	// Read in JSON object and metadata objects.
	jsonObject := map[string]interface{}{}
	attachmentObject := &Attachment{}

	if err := c.Read(keyName, jsonObject, false); err != nil {
		return fmt.Errorf("Can't read %q, aborting, %s", keyName, err)
	}
	// This is the full path to the JSON Object document
	docPath, err := c.DocPath(keyName)
	if err != nil {
		return fmt.Errorf("Can't find document path %q, aborting, %s", keyName, err)
	}
	// This is JSON object's directory.
	docDir := path.Dir(docPath)

	// Now we're ready to get our attachment list.
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == true {
		for _, obj := range attachmentList {
			if obj.Name == dstName {
				attachmentObject = obj
				break
			}
		}
	}

	// Update the metadata
	// Read in the attachment so I can compute the checksum as well as size.
	checksum, err := calcChecksum(srcName)
	if err != nil {
		return fmt.Errorf("Checksum error %q, %s", srcName, err)
	}

	attachmentObject.Name = dstName
	attachmentObject.Version = semver
	// Add/update our version href
	attachmentObject.HRef = path.Join(docDir, semver, dstName)
	// Write out attached filename and update size.
	err = os.MkdirAll(path.Dir(attachmentObject.HRef), 0777)
	if err != nil {
		return fmt.Errorf("Can't create directory for %q, %s", attachmentObject.HRef, err)
	}
	r, err := os.Open(srcName)
	if err != nil {
		return fmt.Errorf("Can't read %q, %s", srcName, err)
	}
	defer r.Close()
	w, err := os.Create(attachmentObject.HRef)
	if err != nil {
		return fmt.Errorf("Can't create %q, %s", attachmentObject.HRef, err)
	}
	defer w.Close()
	l, err := io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("Copy error, %s", err)
	}
	attachmentObject.Size = l
	if attachmentObject.Sizes == nil {
		attachmentObject.Sizes = make(map[string]int64)
	}
	attachmentObject.Sizes[semver] = l
	// Compute Checksum with md5 and store as a string
	if attachmentObject.Checksums == nil {
		attachmentObject.Checksums = make(map[string]string)
	}
	attachmentObject.Checksums[semver] = checksum
	// We need to make the semver directory if necessary
	if attachmentObject.VersionHRefs == nil {
		attachmentObject.VersionHRefs = make(map[string]string)
	}
	attachmentObject.VersionHRefs[semver] = attachmentObject.HRef
	now := time.Now()
	if attachmentObject.Created == "" {
		attachmentObject.Created = now.Format(time.RFC3339)
	}
	attachmentObject.Modified = now.Format(time.RFC3339)

	jsonObject["_Attachments"] = updateAttachmentList(attachmentList, attachmentObject)

	// Write out updated JSON Object and return any error
	err = c.Update(keyName, jsonObject)
	return err
}

// AttachFiles attaches non-JSON documents to a JSON document in the collection.
// Attachments are stored in a tar file, if tar file exits then attachment(s)
// are appended to tar file.
func (c *Collection) AttachFiles(keyName string, semver string, fileNames ...string) error {
	for _, fName := range fileNames {
		if err := c.AttachFile(keyName, semver, fName); err != nil {
			return err
		}
	}
	return nil
}

// Attachments returns a list of files and size attached for a key name in the collection
func (c *Collection) Attachments(keyName string) ([]string, error) {
	jsonObject := map[string]interface{}{}
	err := c.Read(keyName, jsonObject, false)
	if err != nil {
		return nil, fmt.Errorf("Can't find %s", keyName)
	}
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == false {
		return []string{}, nil
	}
	s := []string{}
	for _, attachment := range attachmentList {
		size := attachment.Size
		s = append(s, fmt.Sprintf("%s %d", attachment.Name, size))
	}
	return s, nil
}

// AttachmentPath takes a key, semver and filename and returns the path
// to the attached file (if found).
func (c *Collection) AttachmentPath(keyName string, semver string, filename string) (string, error) {
	if c.KeyExists(keyName) == false {
		return "", fmt.Errorf("No key found for %q", keyName)
	}
	objPath, err := c.DocPath(keyName)
	if err != nil {
		return "", fmt.Errorf("Cannot find object %s", keyName)
	}
	fName := path.Join(path.Dir(objPath), semver, filename)
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot find %s %s/%s", keyName, semver, filename)
	}
	return fName, nil
}

func filterNameFound(a []string, target string) bool {
	if len(a) == 0 {
		return true
	}
	for _, s := range a {
		if s == target {
			return true
		}
	}
	return false
}

// GetAttachedFiles returns an error if encountered, a side effect
// is the file(s) are written to the current work directory
// If no filterNames provided then return all attachments are written out
// An error value is always returned.
func (c *Collection) GetAttachedFiles(keyName string, semver string, filterNames ...string) error {
	if c.KeyExists(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	jsonObject := map[string]interface{}{}
	if err := c.Read(keyName, jsonObject, false); err != nil {
		return fmt.Errorf("Can't read %q, %s", keyName, err)
	}
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == false {
		return fmt.Errorf("No attachments")
	}
	version := semver
	for _, obj := range attachmentList {
		if filterNameFound(filterNames, obj.Name) {
			// Are we getting the current version?
			if semver == "" {
				version = obj.Version
			}
			// Retrieve the file by version
			if href, ok := obj.VersionHRefs[version]; ok == true {
				src, err := os.ReadFile(href)
				if err != nil {
					return err
				} else if err := os.WriteFile(obj.Name, src, 0777); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Can't find %s %q for key %q", semver, obj.Name, keyName)
			}
		}
	}
	return nil
}

// Prune a non-JSON document from a JSON document in the collection.
func (c *Collection) Prune(keyName string, semver string, filterNames ...string) error {
	if c.KeyExists(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	jsonObject := map[string]interface{}{}
	if err := c.Read(keyName, jsonObject, false); err != nil {
		return fmt.Errorf("Can't read %q, %s", keyName, err)
	}
	newAttachmentList := []*Attachment{}
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == false {
		return fmt.Errorf("No attachments found")
	}
	for _, obj := range attachmentList {
		if filterNameFound(filterNames, obj.Name) {
			// Are we getting the current version?
			// Check for a prior version
			if href, ok := obj.VersionHRefs[semver]; ok == true {
				if err := os.Remove(href); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("Can't find %s %q for key %q", semver, obj.Name, keyName)
			}

		} else {
			newAttachmentList = append(newAttachmentList, obj)
		}
	}
	// Now we need to update our attachments list and update the JSON document
	jsonObject["_Attachments"] = newAttachmentList
	if err := c.Update(keyName, jsonObject); err != nil {
		return err
	}
	return nil
}
