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
	"fmt"
)

// Attachment is a structure for holding non-JSON content metadata
// you wish to store alongside a JSON document in a collection
// Attachments themselves reside in a pairtree of the collection
// (even when using a SQL store for the JSON document). They are
// stored in the object's expected pairtree location in a sub-directory
// called "_".
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

// AttachStream is for attaching open a non-JSON file buffer (via an io.Reader).
func (c *Collection) AttachStream(keyName, semver, fullName string, buf io.Reader) error {
	return fmt.Errorf("AttachStream() not implemented")
}

// AttachFile is for attaching a single non-JSON document to a
// dataset record. It will replace ANY existing attached content
// with the same semver and basename.
func (c *Collection) AttachFile(keyName string, semver string, fullName string) error {
	return fmt.Errorf("AttachFile() not implemented")
}

// AttachFileAs is for attaching a single non-JSON document to a
// dataset record with a specific attachment name. It will replace ANY
// existing attached content with the same semver and destintation name.
func (c *Collection) AttachFileAs(keyName string, semver string, dstName string, srcName string) error {
	return fmt.Errorf("AttachFileAs() not implemented")
}

// AttachFiles attaches non-JSON documents to a JSON document in the collection.
// Attachments are stored in a tar file, if tar file exits then attachment(s)
// are appended to tar file.
func (c *Collection) AttachFiles(keyName string, fNames ...string) error {
	return fmt.Errorf("AttachFiles() not implemented")
}

// Attachments returns a list of files and size attached for a key name in the collection
func (c *Collection) Attachments(keyName string) ([]string, error) {
	return nil, fmt.Errorf("Attachments() not implemented")
}

// AttachmentPath takes a key, semver and filename and returns the path
// to the attached file (if found).
func (c *Collection) AttachmentPath(keyName string, semver string, filename string) (string, error) {
	return "", fmt.Errorf("AttachmentPath() not implemented")
}

// GetAttachedFiles returns an error if encountered, a side effect
// is the file(s) are written to the current work directory
// If no filterNames provided then return all attachments are written out
// An error value is always returned.
func (c *Collection) GetAttachedFiles(keyName string, fNames ...string) error {
	return fmt.Errorf("GetAttachedFiles() not implemented")
}

// Prune a non-JSON document from a JSON document in the collection.
func (c *Collection) Prune(keyName string, fNames ...string) error {
	return fmt.Errorf("Prune() not implemented")
}
