// texts is part of dataset
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
	"io/ioutil"
	"strings"
)

// StringProcessor takes the a text and replaces all the keys
// (e.g. "{app_name}") with their value (e.g. "dataset"). It is
// used to prepare command line and daemon document for display.
func StringProcessor(varMap map[string]string, text string) string {
	src := text[:]
	for key, val := range varMap {
		if strings.Contains(src, key) {
			src = strings.ReplaceAll(src, key, val)
		}
	}
	return src
}

// BytesProcessor takes the a text and replaces all the keys
// (e.g. "{app_name}") with their value (e.g. "dataset"). It is
// used to prepare command line and daemon document for display.
func BytesProcessor(varMap map[string]string, text []byte) []byte {
	src := text[:]
	for key, val := range varMap {
		if bytes.Contains(src, []byte(key)) {
			src = bytes.ReplaceAll(src, []byte(key), []byte(val))
		}
	}
	return src
}

// ReadSource reads the source text from a filename or
// io.Reader (e.g. standard input) as fallback.
//
// ```
// src, err := ReadSource(inputName, os.Stdin)
//
//	if err != nil {
//	   ...
//	}
//
// ```
func ReadSource(fName string, in io.Reader) ([]byte, error) {
	var (
		src []byte
		err error
	)
	if fName == "" || fName == "-" {
		src, err = ioutil.ReadAll(in)
	} else {
		src, err = ioutil.ReadFile(fName)
	}
	return src, err
}

// WriteSource writes a source text to a file or to the io.Writer
// of filename not set.
func WriteSource(fName string, out io.Writer, src []byte) error {
	if fName == "" || fName == "-" {
		_, err := out.Write(src)
		return err
	}
	return ioutil.WriteFile(fName, src, 0664)
}

// ReadKeys reads a list of keys given filename or an io.Reader
// (e.g. standard input) as fallback. The key file should be formatted
// one key per line with a line delimited of "\n".
//
// ```
//
//	keys, err := dataset.ReadKeys(keysFilename, os.Stdin)
//	if err != nil {
//	}
//
// ```
func ReadKeys(keysName string, in io.Reader) ([]string, error) {
	src, err := ReadSource(keysName, in)
	if err != nil {
		return nil, err
	}
	keys := strings.Split(fmt.Sprintf("%s", bytes.TrimSpace(src)), "\n")
	return keys, nil
}

// WriteKeys writes a list of keys to given filename or to io.Writer
// as fallback. The key file is formatted as one key per line using
// "\n" as a separator.
//
// ```
//
//	keys := ...
//	if err := WriteKeys(out, keyFilename, keys); err != nil {
//	   ...
//	}
//
// ```
func WriteKeys(keyFilename string, out io.Writer, keys []string) error {
	src := []byte(strings.Join(keys, "\n"))
	return WriteSource(keyFilename, out, src)
}
