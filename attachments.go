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
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Attachment is a structure for holding non-JSON content you wish to store alongside a JSON document in a collection
type Attachment struct {
	// Name is the filename and path to be used inside the generated tar file
	Name string `json:"name"`

	// Content is a byte array for storing the content associated with Name
	// NOTE: It is NOT written out in the Attachment metadata, hence json:"-".
	Content []byte `json:"-"`

	// Size
	Size int64 `json:"size"`

	// Checksum, current implemented as a MD5 checksum for now
	Checksum string `json:"checksum"`

	// HRef points at last attached version of the attached document, e.g. v0.0.0/photo.png
	// If you moved an object out of the pairtree it should be a URL.
	HRef string `json:"href"`
	// VersionHRefs is a map to all versions of the attached document
	// {
	//    "v0.0.0": "photo.png",
	//    "v0.0.1": "photo.png",
	//    "v0.0.2": "photo.png"
	// }
	VersionHRefs map[string]string `json:"version_hrefs"`

	// Created a date string in RTC1123 format
	Created string `json:"created"`

	// Modified a date string in RFC1123 format
	Modified string `json:"modified"`

	// Metadata is a map to put application specific metadata in.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// attachmentNames takes a key, semver and filename and returns a path to the
// metadata and a path to where the file should be stored. Returns an error
// if the key is not found in collection.
func (c *Collection) attachmentNames(keyName, semver, fName string) (string, string, error) {
	var (
		attachmentFile     string
		attachmentMetadata string
	)
	if c.HasKey(keyName) == false {
		return "", "", fmt.Errorf("No key found for %q", keyName)
	}
	return attachmentFile, attachmentMetadata, nil
}

// getAttachmentList takes a JSON objects, pulls our "_Attachments" and reads them into
// and array of Attachment. Returns true if we have a list, false otherwise
func getAttachmentList(jsonObject map[string]interface{}) ([]*Attachment, book) {
	attachmentList := make([]*Attachment)
	attachments, ok := jsonObject["_Attachments"]
	if ok == false {
		return attachmentList, false
	}
	for i, obj := range attachments {
		attachment := new(Attachment)
		if name, ok := obj["name"]; ok == true {
			attachment.Name = string(name)
		}
		if size, ok := obj["size"]; ok == true {
			attachment.Size = int64(size)
		}
		if metadata, ok := obj["metadata"]; ok == true {
			attachment.Metadata = make(map[string]interface{})
			for k, v := range metadata {
				attachment.Metadata[k] = v
			}
		}
		if checksum, ok := obj["checksum"]; ok == true {
			attachment.Checksum = string(checksum)
		}
		if created, ok := obj["created"]; ok == true {
			attachment.Created = string(created)
		}
		if modified, ok := obj["modified"]; ok == true {
			attachment.Modified = string(modified)
		}
		if href, ok := obj["href"]; ok == true {
			attachment.HRef = string(href)
		}
		if versionHRefs, ok := obj["version_hrefs"]; ok == true {
			attachment.VersionHRefs = make(map[string]string)
			for k, v := range versionHRefs {
				attachment.VersionHRefs[string(k)] = string(v)
			}
		}
		attachmentList = append(attachmentList, attachment)
	}
	return attachmentList, true
}

// AttachFile is for attaching a single non-JSON document to a dataset record. It will replace
// ANY existing attached content with the same semver and basename.
func (c *Collection) AttachFile(keyName, semver string, fName string, buf io.Reader) error {
	if c.HasKey(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	if semver == "" {
		// We use version v0.0.0 for "unversioned" attachments.
		semver = "v0.0.0"
	}
	// Normalize fName to basename
	fName = path.Basename(fName)
	// Figure out where we're going to write things to
	attachmentFName, metadataFName, err := c.attachmentNames(keyName, semver, fName)
	if err != nil {
		return err
	}
	// Read in JSON object and metadata objects.
	jsonObject := map[string]interface{}{}
	attachmentObject := Attachment{}

	err := c.Read(keyName, jsonObject)
	if err != nil {
		return fmt.Errorf("Can't read %q, aborting, %s", keyName, err)
	}
	attachmentList, ok := getAttachmentList(jsonObject)
	if ok == true {
		for i, obj := range attachmentList {
			if obj.Name == fName {
				attachmentObject = obj
				break
			}
		}
	}

	// Update the metadata and write our metadata filename
	// Update JSON object record and write out
	// Write out attached filename
	// Write out attached metadata
	// Write out updated JSON Object

	// Normalize our fName to only be a basename
	baseName := path.Basename(fName)
	// Our metadata file will be basename plus .json
	metadataName := baseName + ".json"

	// FIXME: might make smore sense to just be a single map[string]interface{} rather than an array of maps unless we're
	// journalling the history of attachments.
	docListing := []map[string]interface{}{}
	attachmentPath, metadataPath = attachmentNames(docPath, fName, semver)
	err = c.Store.WriteFilter(docPath, func(fp *os.File) error {
		// NOTE: we always overwrite an existing tarball before our storage module
		// does can't assume an append operation is available for cloud storage.
		tw := tar.NewWriter(fp)

		// Read in our data
		data, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		fSize := int64(len(data))

		// Save our doc info for our metadata
		docInfo := map[string]interface{}{}
		docInfo["name"] = path.Base(fName)
		docInfo["size"] = fSize
		docListing = append(docListing, docInfo)

		// Add our data to the tar ball
		hdr := &tar.Header{
			Name: path.Base(fName),
			Mode: 0664,
			Size: fSize,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(data); err != nil {
			return err
		}
		return tw.Close()
	})

	// Now update the _Attachment attribute in the JSON document.

	//NOTE: Because we're always replacing the tarball (can't append in a cloud environment)
	// we must also always replace the attachments metadata.
	rec["_Attachments"] = docListing
	err = c.Update(keyName, rec)
	return err
}

// AttachFiles attaches non-JSON documents to a JSON document in the collection.
// Attachments are stored in a tar file, if tar file exits then attachment(s)
// are appended to tar file.
func (c *Collection) AttachFiles(name string, fileNames ...string) error {
	keyName, _ := keyAndFName(name)
	if c.HasKey(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	rec := map[string]interface{}{}
	err := c.Read(keyName, rec)
	if err != nil {
		return fmt.Errorf("Can't read %q, aborting, %s", keyName, err)
	}

	// NOTE: we normalize the keyName to omit a .json file extension,
	// make sure we have an associated JSON record, then generate a new tarball
	// from attachments.
	docPath, err := c.DocPath(keyName)
	if err != nil {
		return err
	}

	docListing := []map[string]interface{}{}
	docPath = tarballName(docPath)
	err = c.Store.WriteFilter(docPath, func(fp *os.File) error {
		tw := tar.NewWriter(fp)

		// For each filename add the file to the tar ball
		for _, fName := range fileNames {
			// Read in our data
			data, err := ioutil.ReadFile(fName)
			if err != nil {
				return err
			}
			fSize := int64(len(data))

			// Save Our Doc info for our metadata
			docInfo := map[string]interface{}{}
			docInfo["name"] = path.Base(fName)
			docInfo["size"] = fSize
			docListing = append(docListing, docInfo)

			hdr := &tar.Header{
				Name: path.Base(fName),
				Mode: 0664,
				Size: fSize,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}
		return tw.Close()
	})
	// Now update the _Attachment attribute in the JSON document.

	// Finally update our attachments metadata based on the new document(s) info
	rec["_Attachments"] = docListing
	err = c.Update(keyName, rec)
	return err
}

// Attachments returns a list of files in the attached tarball for a given name in the collection
func (c *Collection) Attachments(name string) ([]string, error) {
	keyName, _ := keyAndFName(name)
	rec := map[string]interface{}{}
	err := c.Read(keyName, rec)
	if err != nil {
		return nil, fmt.Errorf("Can't find %s", name)
	}
	fileNames := []string{}
	if l, ok := rec["_Attachments"]; ok == true {
		var (
			size  json.Number
			fname string
		)
		for _, valL := range l.([]interface{}) {
			m := valL.(map[string]interface{})
			fname = ""
			if s, ok := m["name"]; ok == true {
				fname = s.(string)
			}
			if i, ok := m["size"]; ok == true && fname != "" {
				size = i.(json.Number)
			}
			fileNames = append(fileNames, fmt.Sprintf("%s %s", fname, size))
		}
		return fileNames, nil
	}

	docPath, err := c.DocPath(name)
	if err != nil {
		return nil, err
	}
	docPath = tarballName(docPath)

	// Get the file and read into memory
	buf, err := c.Store.ReadFile(docPath)
	if err != nil {
		// FIXME: If no tarball, then return error "no attachments found"
		return nil, fmt.Errorf("no attachments found")
	}
	fp := bytes.NewBuffer(buf)

	tr := tar.NewReader(fp)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tarball
			break
		}
		if err != nil {
			// error reading tarball
			return fileNames, err
		}
		s := fmt.Sprintf("%s %d", hdr.Name, hdr.Size)
		fileNames = append(fileNames, s)
	}
	return fileNames, nil
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

// (DEPRECIATE) getAttached returns an Attachment array or error
// If no filterNames provided then return all attachments or error
func (c *Collection) getAttached(name string, filterNames ...string) ([]Attachment, error) {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then remove any tarball
	docPath, err := c.DocPath(name)
	if err != nil {
		return nil, err
	}
	docPath = tarballName(docPath)

	attachments := []Attachment{}

	buf, err := c.Store.ReadFile(docPath)
	if err != nil {
		return nil, err
	}
	fp := bytes.NewBuffer(buf)

	tr := tar.NewReader(fp)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tarball
			break
		}
		if err != nil {
			// error reading tarball
			return attachments, err
		}
		if filterNameFound(filterNames, hdr.Name) == true {
			buf := bytes.NewBuffer([]byte{})
			if _, err := io.Copy(buf, tr); err != nil {
				return attachments, err
			}
			attachments = append(attachments, Attachment{
				Name: hdr.Name,
				Body: buf.Bytes(),
			})
		}
	}
	return attachments, nil
}

// GetAttachedFiles returns an error if encountered, side effect is to write file to destination directory
// If no filterNames provided then return all attachments or error
func (c *Collection) GetAttachedFiles(name string, filterNames ...string) error {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then remove any tarball
	docPath, err := c.DocPath(name)
	if err != nil {
		return err
	}
	docPath = tarballName(docPath)

	buf, err := c.Store.ReadFile(docPath)
	if err != nil {
		return err
	}
	fp := bytes.NewBuffer(buf)

	tr := tar.NewReader(fp)
	found := []string{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tarball
			break
		}
		if err != nil {
			// error reading tarball
			return err
		}
		if filterNameFound(filterNames, hdr.Name) == true {
			// NOTE: write file to disc, not using defer because we want fp to close after each loop
			fp, err := os.Create(hdr.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(fp, tr); err != nil {
				fp.Close()
				return err
			}
			fp.Close()
			if len(filterNames) > 0 {
				found = append(found, hdr.Name)
			}
		}
	}
	if len(filterNames) > len(found) {
		missing := []string{}
		if len(found) > 0 {
			for _, fname := range filterNames {
				if filterNameFound(found, fname) == false {
					missing = append(missing, fname)
				}
			}
		} else {
			missing = filterNames
		}
		return fmt.Errorf("Missing attachments %q", strings.Join(missing, `", "`))
	}
	return nil
}

// Prune a non-JSON document from a JSON document in the collection.
func (c *Collection) Prune(name string, filterNames ...string) error {
	keyName, _ := keyAndFName(name)
	if c.HasKey(keyName) == false {
		return fmt.Errorf("No key found for %q", keyName)
	}
	rec := map[string]interface{}{}
	err := c.Read(keyName, rec)
	if err != nil {
		return fmt.Errorf("Can't read %q, aborting, %s", keyName, err)
	}
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then remove any tarball
	docPath, err := c.DocPath(name)
	if err != nil {
		return err
	}
	docPath = tarballName(docPath)
	if path.Ext(docPath) != ".tar" {
		return fmt.Errorf("Can't remove %q attachments", docPath)
	}

	// NOTE: If we're removing everything then just call Removeall on store for that tarball name
	if len(filterNames) == 0 {
		err := c.Store.RemoveAll(docPath)
		if err != nil {
			return err
		}
		delete(rec, "_Attachments")
		err = c.Update(keyName, rec)
		return err
	}

	// NOTE: If we're only removing some of the attached files then we need to re-write the tarball after reading it into memory
	buf, err := c.Store.ReadFile(docPath)
	if err != nil {
		return err
	}
	rd := bytes.NewBuffer(buf)
	tr := tar.NewReader(rd)

	docList := []map[string]interface{}{}

	// Read in old tarball and only write out files that match filterNames
	found := []string{}
	err = c.Store.WriteFilter(docPath, func(fp *os.File) error {
		tw := tar.NewWriter(fp)
		defer tw.Close()

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				// end of tar archive
				break
			}
			if err != nil {
				return err
			}
			if filterNameFound(filterNames, hdr.Name) == false {
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if _, err := io.Copy(tw, tr); err != nil {
					return err
				}
				//NOTE: Update the attachment list with remaining docs
				docInfo := map[string]interface{}{}
				docInfo["name"] = hdr.Name
				docInfo["size"] = hdr.Size
				docList = append(docList, docInfo)
			} else {
				found = append(found, hdr.Name)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	rec["_Attachments"] = docList
	err = c.Update(keyName, rec)
	if err != nil {
		return err
	}
	if len(found) < len(filterNames) {
		if len(found) == 0 {
			return fmt.Errorf("Unable to prune %q", strings.Join(filterNames, `", "`))
		}
		missing := []string{}
		for _, fName := range filterNames {
			if filterNameFound(found, fName) == false {
				missing = append(missing, fName)
			}
		}
		if len(missing) > 0 {
			return fmt.Errorf("Unable to prune %q", strings.Join(missing, `", "`))
		}
	}
	return nil
}
