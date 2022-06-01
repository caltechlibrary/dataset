//
// pairtree_test.go implements testing for encoding/decoding object identifiers and pairtree paths (ppaths) per
// https://confluence.ucop.edu/download/attachments/14254128/PairtreeSpec.pdf?version=2&modificationDate=1295552323000&api=v2
//
// Author R. S. Doiel, <rsdoiel@library.caltech.edu>
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
package pairtree

import (
	"os"
	"strings"
	"testing"
)

func sepJoin(sep rune, parts ...string) string {
	pathSep := string([]rune{sep})
	return strings.Join(parts, pathSep) + pathSep
}

func TestCharEncoding(t *testing.T) {
	//FIXME: test is posix-centric, need to handle other
	// paths delimiters.
	testCharEncoding := map[string]string{
		"ark:/13030/xt12t3":                     "ark+=13030=xt12t3",
		"http://n2t.info/urn:nbn:se:kb:repos-1": "http+==n2t,info=urn+nbn+se+kb+repos-1",
		"what-the-*@?#!^!?":                     "what-the-^2a@^3f#!^5e!^3f",
	}

	for src, expected := range testCharEncoding {
		result := string(charEncode([]rune(src)))
		if result != expected {
			t.Errorf("%q, expected %q, got %q", src, expected, result)
		}
	}
	for expected, src := range testCharEncoding {
		result := charDecode(src)
		if result != expected {
			t.Errorf("%q, expected %q, got %q", src, expected, result)
		}
	}
}

func TestBasic(t *testing.T) {
	sep := Separator
	// Test Basic encoding
	testEncodings := map[string]string{
		"abcd":       sepJoin(sep, "ab", "cd"),
		"abcdefg":    sepJoin(sep, "ab", "cd", "ef", "g"),
		"12-986xy4":  sepJoin(sep, "12", "-9", "86", "xy", "4"),
		"2018-06-01": sepJoin(sep, "20", "18", "-0", "6-", "01"),
		"a":          sepJoin(sep, "a"),
		"ab":         sepJoin(sep, "ab"),
		"abc":        sepJoin(sep, "ab", "c"),
		"abcde":      sepJoin(sep, "ab", "cd", "e"),
		"mnopqz":     sepJoin(sep, "mn", "op", "qz"),
	}
	for src, expected := range testEncodings {
		result := Encode(src)
		if result != expected {
			t.Errorf("encoding %q, expected %q, got %q", src, expected, result)
		}
	}

	testDecodings := map[string]string{}
	for val, key := range testEncodings {
		testDecodings[key] = val
	}

	// Test Basic decoding
	for src, expected := range testDecodings {
		result := Decode(src)
		if result != expected {
			t.Errorf("decoding %q, expected %q, got %q", src, expected, result)
		}
	}
}

func TestCustomSeparator(t *testing.T) {
	sep := os.PathSeparator
	if Separator != os.PathSeparator {
		t.Errorf("separator not set, expected %c, got %c", sep, Separator)
		t.FailNow()
	}
	sep = ':'
	Set(':')
	if Separator != sep {
		t.Errorf("separator not set, expected %c, got %c", sep, Separator)
		t.FailNow()
	}
	testEncodings := map[string]string{
		"abcd":       sepJoin(sep, "ab", "cd"),
		"abcdefg":    sepJoin(sep, "ab", "cd", "ef", "g"),
		"12-986xy4":  sepJoin(sep, "12", "-9", "86", "xy", "4"),
		"2018-06-01": sepJoin(sep, "20", "18", "-0", "6-", "01"),
		"a":          sepJoin(sep, "a"),
		"ab":         sepJoin(sep, "ab"),
		"abc":        sepJoin(sep, "ab", "c"),
		"abcde":      sepJoin(sep, "ab", "cd", "e"),
		"mnopqz":     sepJoin(sep, "mn", "op", "qz"),
	}
	testDecodings := map[string]string{}
	for val, key := range testEncodings {
		testDecodings[key] = val
	}

	for src, expected := range testEncodings {
		result := Encode(src)
		if result != expected {
			t.Errorf("encoding %q, expected %q, got %q", src, expected, result)
		}
	}

	// Test Basic decoding
	for src, expected := range testDecodings {
		result := Decode(src)
		if result != expected {
			t.Errorf("decoding %q, expected %q, got %q", src, expected, result)
		}
	}
	sep = os.PathSeparator
	Set(os.PathSeparator)
	if Separator != os.PathSeparator {
		t.Errorf("separator not set, expected %c, got %c", sep, Separator)
		t.FailNow()
	}
}

func TestAdvanced(t *testing.T) {
	testData := map[string]string{
		"ark:/13030/xt12t3": sepJoin(Separator, "ar", "k+", "=1", "30", "30",
			"=x", "t1", "2t", "3"),
		"http://n2t.info/urn:nbn:se:kb:repos-1": sepJoin(Separator, "ht", "tp",
			"+=", "=n", "2t", ",i", "nf", "o=", "ur", "n+", "nb", "n+",
			"se", "+k", "b+", "re", "po", "s-", "1"),
		"what-the-*@?#!^!?": sepJoin(Separator, "wh", "at", "-t", "he", "-^",
			"2a", "@^", "3f", "#!", "^5", "e!", "^3",
			"f"),
	}
	for src, expected := range testData {
		result := Encode(src)
		if result != expected {
			t.Errorf("encode %q, expected %q, got %q", src, expected, result)
		}
	}
	for expected, src := range testData {
		result := Decode(src)
		if result != expected {
			t.Errorf("decode %q, expected %q, got %q", src, expected, result)
		}
	}
}

func TestUTF8Names(t *testing.T) {
	testData := map[string]string{
		"Hänggi-P": sepJoin(Separator, "Hä", "ng", "gi", "-P"),
	}
	for src, expected := range testData {
		result := Encode(src)
		if result != expected {
			t.Errorf("encode %q, expected %q, got %q", src, expected, result)
		}
	}
	for expected, src := range testData {
		result := Decode(src)
		if result != expected {
			t.Errorf("decode %q, expected %q, got %q", src, expected, result)
		}
	}
}
