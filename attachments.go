//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	// Caltech Library packages
	//"github.com/caltechlibrary/storage"
)

// Attachment is a structure for holding non-JSON content you wish to store alongside a JSON document in a collection
type Attachment struct {
	// Name is the filename and path to be used inside the generated tar file
	Name string `json:"name"`
	// Body is a byte array for storing the content associated with Name
	Body []byte `json:"-"`
	// Size
	Size int64 `json:"size"`
}

// tarballName trims any .json from the record name.
func tarballName(docPath string) string {
	if path.Ext(docPath) == ".json" {
		return strings.TrimSuffix(docPath, ".json") + ".tar"
	}
	return docPath + ".tar"
}

// (DEPRECIATE) attach non-JSON documents to a JSON document in the collection.
// Attachments are stored in a tar file, if tar file exits then attachment(s)
// are appended to tar file.
func (c *Collection) attach(name string, attachments ...*Attachment) error {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then generate a new tarball
	// from attachments.
	docPath, err := c.DocPath(name)
	if err != nil {
		return err
	}

	// this is the name we will move two, we build the tarball as a tmp file.
	info := []map[string]interface{}{}
	docPath = tarballName(docPath)
	err = c.Store.WriteFilter(docPath, func(fp *os.File) error {
		tw := tar.NewWriter(fp)
		// For each attachtment add to the tar ball
		for _, attachment := range attachments {
			hdr := &tar.Header{
				Name: attachment.Name,
				Mode: 0664,
				Size: int64(len(attachment.Body)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if _, err := tw.Write(attachment.Body); err != nil {
				return err
			}
			m := map[string]interface{}{
				"name": attachment.Name,
				"size": int64(len(attachment.Body)),
			}
			info = append(info, m)
		}
		return tw.Close()
	})
	// Now update the _Attachment attribute in the JSON document.
	keyName, _ := keyAndFName(name)
	rec := make(map[string]interface{})
	err = c.Read(keyName, rec)
	if err != nil {
		return err
	}
	rec["_Attachments"] = info
	err = c.Update(name, rec)
	return err
}

// AttachFile is for attaching a single non-JSON document to a dataset record. It will replace
// ANY existing attached content (i.e. it creates an new tarball holding on this single document)
// This is a limitation of our storage package supporting a minimal set of operation across all
// storage environments (e.g. S3/Google Cloud Storage do not support append to file, only replacement).
// It takes the document key, name and an io.Reader reading content in and appending the results to
// the tar file updating the internal _Attributes metadata as needed.
func (c *Collection) AttachFile(keyName, fName string, buf io.Reader) error {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then generate a new tarball
	// from attachments.
	docPath, err := c.DocPath(keyName)
	if err != nil {
		return err
	}

	// FIXME: might make smore sense to just be a single map[string]interface{} rather than an array of maps unless we're
	// journalling the history of attachments.
	docListing := []map[string]interface{}{}
	docPath = tarballName(docPath)
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
		docInfo["name"] = fName
		docInfo["size"] = fSize
		docListing = append(docListing, docInfo)

		// Add our data to the tar ball
		hdr := &tar.Header{
			Name: fName,
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
	rec := map[string]interface{}{}
	err = c.Read(keyName, rec)
	if err != nil {
		return err
	}

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
			docInfo["name"] = fName
			docInfo["size"] = fSize
			docListing = append(docListing, docInfo)

			hdr := &tar.Header{
				Name: fName,
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
	rec := map[string]interface{}{}
	err = c.Read(keyName, rec)
	if err != nil {
		return err
	}
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
		}
	}
	return nil
}

// Detach a non-JSON document from a JSON document in the collection.
func (c *Collection) Detach(name string, filterNames ...string) error {
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
		return c.Store.RemoveAll(docPath)
	}

	// NOTE: If we're only removing some of the attached files then we need to re-write the tarball after reading it into memory
	buf, err := c.Store.ReadFile(docPath)
	if err != nil {
		return err
	}
	rd := bytes.NewBuffer(buf)
	tr := tar.NewReader(rd)

	// Read in old tarball and only write out files that match filterNames
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
			}
		}
		return nil
	})
	return err
}
