//
// pairtree.go implements encoding/decoding of object identifiers and pairtree paths (paths) per
// https://confluence.ucop.edu/download/attachments/14254128/PairtreeSpec.pdf?version=2&modificationDate=1295552323000&api=v2
//
// Author R. S. Doiel, <rsdoiel@library.caltech.edu>
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
package pairtree

import (
	"os"
	"strings"
)

var (
	// Separator is the path separator used for your pairtree.
	// by default it is set to your operating system's path separator.
	Separator = '/'

	stepOneEncoding = map[rune][]rune{
		' ':  []rune("^20"),
		'"':  []rune("^22"),
		'<':  []rune("^3c"),
		'\\': []rune("^5c"),
		'*':  []rune("^2a"),
		'=':  []rune("^3d"),
		'^':  []rune("^5e"),
		'+':  []rune("^2b"),
		'>':  []rune("^3e"),
		'|':  []rune("^7c"),
		',':  []rune("^2c"),
		'?':  []rune("^3f"),
	}
	stepTwoEncoding = map[rune]rune{
		'/': '=',
		':': '+',
		'.': ',',
	}
)

func charEncode(src []rune) []rune {
	// NOTE: We run through stepOneEncoding map first, then stepTwoEncoding...
	results := []rune{}
	for i := 0; i < len(src); i++ {
		if val, ok := stepOneEncoding[src[i]]; ok == true {
			results = append(results, val...)
		} else {
			results = append(results, src[i])
		}
	}
	for i := 0; i < len(results); i++ {
		key := results[i]
		if val, ok := stepTwoEncoding[key]; ok == true {
			results[i] = val
		}
	}
	return results
}

func charDecode(s string) string {
	for replacement, target := range stepTwoEncoding {
		t := string(target)
		r := string(replacement)
		if strings.Contains(s, t) {
			s = strings.Replace(s, t, r, -1)
		}
	}
	for replacement, target := range stepOneEncoding {
		t := string(target)
		r := string(replacement)
		if strings.Contains(s, t) {
			s = strings.Replace(s, t, r, -1)
		}
	}
	return s
}

// Get will return the current separator in use in the package
func Get() rune {
	return Separator
}

// Set will set the separator used in encoding and decoding the Pairtree path
func Set(c rune) {
	Separator = c
}

// Encode takes a string and encodes it as a pairtree path.
func Encode(src string) string {
	s := charEncode([]rune(src))
	results := []rune{}
	for i := 0; i < len(s); i += 2 {
		if len(results) > 0 {
			results = append(results, Separator)
		}
		if (i + 2) < len(s) {
			//FIXME need to char encode here...
			t := s[i : i+2]
			results = append(results, t...)
		} else {
			//FIXME need to char encode here...
			t := s[i:]
			results = append(results, t...)
		}
	}
	if len(results) > 0 {
		results = append(results, Separator)
	}
	return string(results)
}

// Decode takes a pairtree path and returns the original string representation
func Decode(src string) string {
	s := []rune(src)
	results := []string{}
	prev, cur := 0, 0
	for ; cur < len(s); cur++ {
		if s[cur] == Separator {
			switch cur - prev {
			case 2:
				results = append(results, string(s[prev:cur]))
				prev = cur + 1
			case 1:
				results = append(results, string(s[prev:cur]))
			default:
				break
			}
		}
	}
	return charDecode(strings.Join(results, ""))
}

func init() {
	Separator = os.PathSeparator
}
