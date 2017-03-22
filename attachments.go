//
// Package dataset is a go package for managing JSON documents stored on disc
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
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
// Examples:
//
//     collection, err := dataset.Open("dataset/mystuff")
//     if err != nil {
//         log.Fatalf("%s", err)
//     }
//     defer collection.Close()
//     if err := collection.Attach("freda", "docs/helloworld.txt", []byte("Hello World!!!!")); err != nil {
//         log.Fatalf("%s", err)
//     }
package dataset

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path"
)

// Attachment is a structure for holding non-JSON content you wish to store alongside a JSON document in a collection
type Attachment struct {
	// Name is the filename and path to be used inside the generated tar file
	Name string
	// Body is a byte array for storing the content associated with Name
	Body []byte
}

// tarballName trims any .json from the record name.
func tarballName(docPath string) string {
	if path.Ext(docPath) == ".json" {
		return docPath[0:len(docPath)-5] + ".tar"
	}
	return docPath + ".tar"
}

// Attach a non-JSON document to a JSON document in the collection.
// Attachments are stored in a tar file, if tar file exits then attachment(s)
// are appended to tar file.
func (c *Collection) Attach(name string, attachments ...*Attachment) error {
	var (
		fp  *os.File
		err error
	)

	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then generate a new tarball
	// from attachments.
	docPath, err := c.DocPath(name)
	if err != nil {
		return err
	}

	docPath = tarballName(docPath)
	if _, err := os.Stat(docPath); os.IsNotExist(err) == true {
		fp, err = os.Create(docPath)
		if err != nil {
			return err
		}
	} else {
		fp, err = os.OpenFile(docPath, os.O_RDWR, 0664)
		if err != nil {
			return err
		}
		// Move to just before the trailer in the tarball
		if _, err = fp.Seek(-2<<9, os.SEEK_END); err != nil {
			return err
		}
	}
	defer fp.Close()
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
	}
	return tw.Close()
}

// Attachments returns a list of files in the attached tarball for a given name in the collection
func (c *Collection) Attachments(name string) ([]string, error) {
	fileNames := []string{}
	docPath, err := c.DocPath(name)
	if err != nil {
		return nil, err
	}
	docPath = tarballName(docPath)
	fp, err := os.Open(docPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
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
		fileNames = append(fileNames, hdr.Name)
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

// GetAttached returns an Attachment array or error
// If no filterNames provided then return all attachments or error
func (c *Collection) GetAttached(name string, filterNames ...string) ([]Attachment, error) {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then remove any tarball
	docPath, err := c.DocPath(name)
	if err != nil {
		return nil, err
	}
	docPath = tarballName(docPath)

	attachments := []Attachment{}

	fp, err := os.Open(docPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
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
		buf := bytes.NewBuffer([]byte{})
		if _, err := io.Copy(buf, tr); err != nil {
			return attachments, err
		}
		if filterNameFound(filterNames, hdr.Name) == true {
			attachments = append(attachments, Attachment{
				Name: hdr.Name,
				Body: buf.Bytes(),
			})
		}
	}
	return attachments, nil
}

// Detach a non-JSON document from a JSON document in the collection.
func (c *Collection) Detach(name string) error {
	// NOTE: we normalize the name to omit a .json file extension,
	// make sure we have an associated JSON record, then remove any tarball
	docPath, err := c.DocPath(name)
	if err != nil {
		return err
	}
	docPath = tarballName(docPath)
	return os.Remove(docPath)
}
