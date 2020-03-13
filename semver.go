//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2020, Caltech
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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Err holds Semver's error messages
type Err struct {
	Msg string
}

func (err *Err) Error() string {
	return err.Msg
}

// Semver holds the information to generate a semver string
type Semver struct {
	// Major version number (required, must be an integer as string)
	Major string `json:"major"`
	// Minor version number (required, must be an integer as string)
	Minor string `json:"minor"`
	// Patch level (optional, must be an integer as string)
	Patch string `json:"patch,omitempty"`
	// Suffix string, (optional, any string)
	Suffix string `json:"suffix,omitempty"`
	// Timestamp (optional, a timestamp in form of YYYY-MM-DD HH:MM:SS)
}

func (sv *Semver) String() string {
	if sv.Patch == "" {
		return "v" + sv.Major + "." + sv.Minor
	}
	if sv.Suffix == "" {
		return "v" + sv.Major + "." + sv.Minor + "." + sv.Patch
	}
	return "v" + sv.Major + "." + sv.Minor + "." + sv.Patch + sv.Suffix
}

// ToJSON takes a version struct and returns JSON as byte slice
func (sv *Semver) ToJSON() []byte {
	src, _ := json.Marshal(sv)
	return src
}

// ParseSemver takes a byte slice and returns a version struct,
// and an error value.
func ParseSemver(src []byte) (*Semver, error) {
	var (
		i   int
		err error
	)
	sv := new(Semver)
	if bytes.HasPrefix(src, []byte("v")) {
		src = bytes.TrimPrefix(src, []byte("v"))
	}
	parts := strings.Split(string(src), ".")
	if len(parts) > 0 {
		i, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, &Err{Msg: "Major value must be an integer"}
		}
		sv.Major = strconv.Itoa(i)
	} else {
		return nil, &Err{Msg: "Invalid version, expecting semver string"}
	}
	if len(parts) > 1 {
		i, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, &Err{Msg: "Minor value must be an integer"}
		}
		sv.Minor = strconv.Itoa(i)
	} else {
		return nil, &Err{Msg: "Invalid version, expecting semver string"}
	}
	if len(parts) > 2 {
		i, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, &Err{Msg: "Patch value must be an integer"}
		}
		sv.Patch = strconv.Itoa(i)
	}
	if len(parts) > 3 {
		sv.Suffix = parts[3]
	}
	return sv, nil
}

// IncPatch increments the patch level if it is numeric
// or returns an error.
func (sv *Semver) IncPatch() error {
	i, err := strconv.Atoi(sv.Patch)
	if err != nil {
		return err
	}
	i++
	sv.Patch = fmt.Sprintf("%d", i)
	return nil
}

// IncMinor increments a minor version number and zeros the
// patch level or returns an error. Returns an error if increment fails.
func (sv *Semver) IncMinor() error {
	i, err := strconv.Atoi(sv.Minor)
	if err != nil {
		return err
	}
	i++
	sv.Patch = "0"
	sv.Minor = fmt.Sprintf("%d", i)
	return nil
}

// IncMajor increments a major version number, zeros minor
// and patch values. Returns an error if increment fails.
func (sv *Semver) IncMajor() error {
	i, err := strconv.Atoi(sv.Major)
	if err != nil {
		return err
	}
	i++
	sv.Patch = "0"
	sv.Minor = "0"
	sv.Major = fmt.Sprintf("%d", i)
	return nil
}
