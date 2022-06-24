//
// semver is a semantic version number package used by dataset.
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
package semver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
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

// NewSemver takes a major number, minor number, patch number
// and suffix string and returns a populated Semver.
//
// ```
//    version := semver.NewSemver(0, 0, 1, "-alpha")
//    fmt.Printf("Version is now %q", version.String())
// ```
//
func NewSemver(major int, minor int, patch int, suffix string) *Semver {
	version := new(Semver)
	version.Major = fmt.Sprintf("%d", major)
	version.Minor = fmt.Sprintf("%d", minor)
	version.Patch = fmt.Sprintf("%d", patch)
	if suffix != "" {
		version.Suffix = suffix
	}
	return version
}

// Return a *Semver as a encoded semver string.
//
// ```
//   sv := new(semver.Semver)
//   sv.Major = "1"
//   sv.Minor = "0"
//   sv.Patch = "3"
//   fmt.Printf("display semver 1.0.3 -> %q\n", sv.String())
// ```
//
func (sv *Semver) String() string {
	if sv.Patch == "" {
		return sv.Major + "." + sv.Minor
	}
	if sv.Suffix == "" {
		return sv.Major + "." + sv.Minor + "." + sv.Patch
	}
	return sv.Major + "." + sv.Minor + "." + sv.Patch + "-" + sv.Suffix
}

// ToJSON takes a version struct and returns JSON as byte slice
//
// ```
//   sv := new(semver.Semver)
//   sv.Major = "1"
//   sv.Minor = "0"
//   sv.Patch = "3"
//   fmt.Printf("JSON semver -> %s\n", sv.ToJSON())
// ```
//
func (sv *Semver) ToJSON() []byte {
	src, _ := json.Marshal(sv)
	return src
}

// Parse takes a byte slice and returns a version struct,
// and an error value.
//
// ```
//   version := []byte("1.0.3")
//   sv, err := semver.Parse(version)
//   if err != nil {
//      ...
//   }
// ```
//
func Parse(src []byte) (*Semver, error) {
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
		if strings.Contains(parts[2], "-") {
			patch := strings.Split(parts[2], "-")
			parts[2] = patch[0]
			parts = append(parts, patch[1])
		}
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

// ParseString accepts a string instead of a byte slice.
func ParseString(version string) (*Semver, error) {
	return Parse([]byte(version))
}

// IncPatch increments the patch level if it is numeric
// or returns an error.
//
// ```
//   version := []byte("1.0.3")
//   sv, err := semver.Parse(version)
//   if err != nil {
//      ...
//   }
//   sv.IncPatch()
//   fmt.Printf("display 1.0.4 -> %q\n", sv.String())
// ```
//
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
//
// ```
//   version := []byte("1.0.3")
//   sv, err := semver.Parse(version)
//   if err != nil {
//      ...
//   }
//   sv.IncMinor()
//   fmt.Printf("display 1.1.0 -> %q\n", sv.String())
// ```
//
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
//
// ```
//   version := []byte("1.0.3")
//   sv, err := semver.Parse(version)
//   if err != nil {
//      ...
//   }
//   sv.IncMajor()
//   fmt.Printf("display 2.0.0 -> %q\n", sv.String())
// ```
//
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

// Less compares two semvers and if semver "a" is less than semver "b"
// returns true false otherwise.
//
// ```
//   versionA, versionB := "1.0.3", "1.4.2"
//   svA, err := semver.Parse(versionA)
//   ...
//   svB, err := semver.Parse(versionB)
//   ...
//   if semver.Less(a, b) {
//      fmt.Printf("%q is less than %q\n", versionA, versionB)
//   } else if versionA == versionB {
//      fmt.Printf("%q is equal to %q\n", versionA, versionB)
//   } else {
//      fmt.Printf("%q is greater than %q\n", versionA, versionB)
//   }
// ```
func Less(a *Semver, b *Semver) bool {
	// Compare Major value
	x, err := strconv.ParseInt(a.Major, 10, 64)
	if err != nil {
		return false
	}
	y, err := strconv.ParseInt(b.Major, 10, 64)
	if err != nil {
		return false
	}
	if x != y {
		return x < y
	}
	// Compare Minor value
	x, err = strconv.ParseInt(a.Minor, 10, 64)
	if err != nil {
		return false
	}
	y, err = strconv.ParseInt(b.Minor, 10, 64)
	if err != nil {
		return false
	}
	if x != y {
		return x < y
	}
	// Compare Patch
	x, err = strconv.ParseInt(a.Patch, 10, 64)
	if err != nil {
		return false
	}
	y, err = strconv.ParseInt(b.Patch, 10, 64)
	if err != nil {
		return false
	}
	return x < y
}

// Sort semvers takes a slice of Semver, sorts them in ascending order.
//
// ```
//   versionA, versionB, vesionC := "1.0.3", "9.2.1", "0.3.6"
//   svList := []*semver.Semver{}
//   sv, _ := semver.Parse(versionA)
//   svA, _ := semver.Parse(versionA)
// ```
func Sort(values []*Semver) {
	sort.Slice(values, func(i, j int) bool {
		return Less(values[i], values[j])
	})
}

// SortStrings takes a list of strings which contain semvers
// convers the strings to a list of Semver structures, sorts the list
// and returns a new the list of strings and an error value
//
// NOTE: Any unparsable semver values passed are skipped and
// not returned in the sortes list of strings.
func SortStrings(values []string) []string {
	l := []*Semver{}
	for _, val := range values {
		sv, err := Parse([]byte(val))
		if err == nil {
			l = append(l, sv)
		}
	}
	Sort(l)
	out := []string{}
	for i := 0; i < len(l); i++ {
		out = append(out, l[i].String())
	}
	return out
}
