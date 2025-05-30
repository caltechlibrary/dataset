// This is part of the dataset package.
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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func setupCliTestCollectionWithMappedObjects(cName string, dsnURI string, mappedObjects map[string]map[string]interface{}) error {
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, dsnURI)
	if err != nil {
		return err
	}
	defer c.Close()
	for k, v := range mappedObjects {
		if err := c.Create(k, v); err != nil {
			return err
		}
	}
	return nil
}

func TestDisplay(t *testing.T) {
	appName := "TestDisplay"
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	output := []byte{}
	out := bytes.NewBuffer(output)
	CliDisplayUsage(out, appName, flagSet)
	if out.Len() == 0 {
		t.Errorf("DisplayUsage() failed, nothing written to output buffer")
	}
}

func TestRunCLIOnCRUDL(t *testing.T) {
	var (
		input, output, errout []byte
	)
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	eout := bytes.NewBuffer(errout)
	// Cleanup stale test data
	cName := path.Join("testout", "RunCLIOnCRUDL_1.ds")
	if _, err := os.Stat(cName); err == nil {
		if err := os.RemoveAll(cName); err != nil {
			t.Errorf("cannot remove stale %q, %s", cName, err)
			t.FailNow()
		}
	}

	// Setup command line args

	// Try intializing a collection

	opt := make(map[string][]string)
	opt["init"] = []string{
		cName,
	}
	opt["create"] = []string{
		cName,
		"1",
		`{"one": 1}`,
	}
	opt["read"] = []string{
		cName,
		"1",
	}
	opt["update"] = []string{
		cName,
		"1",
		`{"one": 1, "two": "2"}`,
	}
	opt["delete"] = []string{
		cName,
		"1",
	}
	opt["keys"] = []string{
		cName,
	}
	opt["haskey"] = []string{
		cName,
		"1",
	}

	// Check if basic CRUDL operations work from cli.
	for _, arg := range []string{"help", "init", "create", "read", "update", "delete", "keys", "haskey", "create", "keys", "haskey"} {
		args := []string{arg}
		if extra, ok := opt[arg]; ok {
			args = append(args, extra...)
		}
		if err := RunCLI(in, out, eout, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
}

func TestCLIKeys(t *testing.T) {
	input, output, errout := []byte{}, []byte{}, []byte{}
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	eout := bytes.NewBuffer(errout)

	cName := path.Join("testout", "CLI_keys.ds")

	mappedObjects := map[string]map[string]interface{}{}
	mappedObjects["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
		"one":  1,
	}
	mappedObjects["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
		"one":  2,
	}
	mappedObjects["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
		"one":  3,
	}
	mappedObjects["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
		"one":  4,
	}
	mappedObjects["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
		"one":  3,
	}
	mappedObjects["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
		"one":  2,
	}
	mappedObjects["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
		"one":  1,
	}
	mappedObjects["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
		"one":  0,
	}
	mappedObjects["character:9"] = map[string]interface{}{
		"name": "Domenique",
		"one":  1,
	}
	mappedObjects["character:10"] = map[string]interface{}{
		"name": "Claudine",
		"one":  1,
	}
	if err := setupCliTestCollectionWithMappedObjects(cName, "", mappedObjects); err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}

	// Setup I/O, args and keys
	args := []string{}
	keys := []string{}

	input = []byte{}
	in = bytes.NewBuffer(input)
	output = []byte{}
	out = bytes.NewBuffer(output)
	errout = []byte{}
	eout = bytes.NewBuffer(errout)

	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("expected to get a list of keys, %s", err)
		t.FailNow()
	}
	src, err := ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("could not read output of keys, %s", err)
		t.FailNow()
	}
	// Get keys for all of collection.
	keys = strings.Split(string(src), "\n")
	if len(keys) == 0 {
		t.Errorf("expected a list of keys to frame, got none")
		t.FailNow()
	}
}

func TestDataset(t *testing.T) {
	cName := path.Join("testout", "test1.ds")

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	expectedB := []byte(``)
	gotB, _ := ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	key := "1"
	args = []string{"create", cName, key, `{"one": 1}`}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	expectedB = []byte(``)
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	key = "2"
	src := []byte(`{"two":2}`)
	if _, err := in.Write(src); err != nil {
		t.Errorf("could not write the input, %s", err)
		t.FailNow()
	}
	args = []string{"create", "-i", "-", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	expectedB = []byte(``)
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	src = []byte(`{"three":3}`)
	key = "3"
	testFile := path.Join("testout", "test3.json")
	ioutil.WriteFile(testFile, src, 0664)

	args = []string{"create", "-i", testFile, cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	expectedB = []byte(``)
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	key = "4"
	src = []byte(`{"four":4}`)
	testFile = path.Join("testout", "test4.json")
	ioutil.WriteFile(testFile, src, 0664)

	args = []string{"create", "-i", testFile, cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	expectedB = []byte(``)
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	expectedB = []byte(`{"one":1}`)
	key = "1"
	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(bytes.TrimSpace(expectedB), bytes.TrimSpace(gotB)) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	expectedB = []byte(`{"two":2}`)
	key = "2"
	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(bytes.TrimSpace(expectedB), bytes.TrimSpace(gotB)) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	expectedB = []byte(`1
2
3
4`)
	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

}

func TestIssue19(t *testing.T) {
	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	cName := path.Join("testout", "test_issue19.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	expectedB := []byte(``)
	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ := ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	expectedB = []byte(``)
	key := "freda"
	src := []byte(`{
    "name": "freda",
	"email": "freda@inverness.example.org",
	"try": 1
}`)
	in.Write(src)
	args = []string{"create", "-i", "-", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	// Now attempt to create the record a second time without -overwrite
	in.Write(src)
	args = []string{"create", "-i", "-", cName, key}
	if err := RunCLI(in, out, eout, args); err == nil {
		t.Errorf("Expected an error when creating a ducplicated record %q in %q", key, cName)
		t.FailNow()
	}
	/* NOTE: eout will be null as that is handle via
	   cmd/dataset/dataset.go taking the return error value of
	   RunCLI(). */

	// Now include the -overwrite option
	expectedB = []byte(``)
	in.Write(src)
	args = []string{"create", "-overwrite", "-i", "-", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(eout)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}
}

func TestReadme(t *testing.T) {
	dName := "testout"
	cName := path.Join(dName, "test_readme.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	// Create our test collection
	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q, %s", cName, err)
		t.FailNow()
	}

	// Create a JSON document
	key := "freda"
	doc := fmt.Sprintf(`{"id": %q, "name": %q, "email":"%s@inverness.example.org"}`, key, key, key)
	args = []string{"create", cName, key, doc}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	// Make sure we have a record called freda
	args = []string{"haskey", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to find key %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	// Read a JSON document back
	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to read key %q in %q, %s", key, cName, err)
		t.FailNow()
	}
	src, err := ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("Reading buffer of JSON document error, %s", err)
	}
	if len(src) == 0 {
		t.Errorf("No data found in %q from %q", key, cName)
	}

	// Update a JSON document
	doc = fmt.Sprintf(`{ "id": %q, "name": %q, "email":"%s@zbs.example.org", "count": 2}`, key, key, key)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to update key %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	// List the keys in the collection
	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to get keys %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	// Delete a JSON document
	args = []string{"delete", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to get keys %q in %q, %s", key, cName, err)
		t.FailNow()
	}
}

func testGettingStarted(t *testing.T) {
	dName := "testout"
	cName := path.Join(dName, "FavoriteThings.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	// test_getting_started

	// Create our test collection
	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q, %s", cName, err)
		t.FailNow()
	}

	key := "beverage"
	args = []string{"create", cName, key, `{"thing":"coffee"}`}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to read %q in %q, %s", key, cName, err)
		t.FailNow()
	}

	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to get keys %q, %s", cName, err)
		t.FailNow()
	}

	filename := path.Join(dName, "jazz-notes.json")
	src := []byte(`{
    "songs": ["Blue Rondo al la Turk", "Bernie's Tune", "Perdido"],
    "pianist": [ "Dave Brubeck" ],
    "trumpet": [ "Dirk Fischer", "Dizzy Gillespie" ]
}
}`)
	if err := ioutil.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("failed to write test file %q, %s", filename, err)
		t.FailNow()
	}

	key = "jazz-notes"
	args = []string{"create", "-i", filename, cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q in %q from %q, %s", key, cName, filename, err)
		t.FailNow()
	}

	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to get keys from %q, %s", cName, err)
		t.FailNow()
	}

	args = []string{"read", cName, "beverage", "jazz-notes"}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to read multiple keys from %q, %s", cName, err)
		t.FailNow()
	}
}
