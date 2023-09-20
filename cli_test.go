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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	// The main dataset package
	dsv1 "github.com/caltechlibrary/dsv1"
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
	DisplayLicense(out, appName)
	if out.Len() == 0 {
		t.Errorf("DisplayLicense() failed, nothing written to output buffer")
	}
	output = []byte{}
	out = bytes.NewBuffer(output)
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
	opt["has-key"] = []string{
		cName,
		"1",
	}

	// Check if basic CRUDL operations work from cli.
	for _, arg := range []string{"help", "init", "create", "read", "update", "delete", "keys", "has-key", "create", "keys", "has-key"} {
		args := []string{arg}
		if extra, ok := opt[arg]; ok {
			args = append(args, extra...)
		}
		if err := RunCLI(in, out, eout, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
}

func TestCloning(t *testing.T) {
	srcName := path.Join("testout", "cloning_src2.ds")
	dstName := path.Join("testout", "cloning_dst2.ds")
	if _, err := os.Stat(srcName); err == nil {
		os.RemoveAll(srcName)
	}
	if _, err := os.Stat(dstName); err == nil {
		os.RemoveAll(dstName)
	}

	// Populate our source repository
	source, err := Init(srcName, "")
	if err != nil {
		t.Errorf("unable to create source %q, %s", srcName, err)
		t.FailNow()
	}
	// Setup a collection to clone
	src := []byte(`[
	{ "one": 1 },
	{ "two": 2 },
	{ "three": 3 },
	{ "four": 4 }
]`)
	testData := []map[string]interface{}{}
	err = json.Unmarshal(src, &testData)
	if err != nil {
		t.Errorf("Can't create testdata")
		t.FailNow()
	}
	for i, obj := range testData {
		key := fmt.Sprintf("%+08d", i)
		if err := source.Create(key, obj); err != nil {
			t.Errorf("failed to create JSON doc for %q in %q, %s", key, srcName, err)
			t.FailNow()
		}
	}
	if source.Length() != 4 {
		t.Errorf("Expected 4 documents in our source repository")
		t.FailNow()
	}
	keys, err := source.Keys()
	if err != nil {
		t.Errorf("can't retrieve source keys, %s", err)
		t.FailNow()
	}
	source.Close()
	// Setup in, out and error buffers
	var (
		input, output []byte
	)
	// Write our keys to the "in" buffer
	input = []byte(strings.Join(keys, "\n"))

	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)

	// Clone repository
	args := []string{"clone", srcName, dstName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
}

func TestSampleCloning(t *testing.T) {
	input, output, errout := []byte{}, []byte{}, []byte{}
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	eout := bytes.NewBuffer(errout)

	srcName := path.Join("testout", "sample_src.ds")
	trainingName := path.Join("testout", "training_src.ds")
	trainingDsnURI := "sqlite://testout/training_src.ds/collection.db"
	testName := path.Join("testout", "test_src.ds")
	testDsnURI := "sqlite://testout/test_src.ds/collection.db"
	if _, err := os.Stat(srcName); err == nil {
		os.RemoveAll(srcName)
	}
	if _, err := os.Stat(trainingName); err == nil {
		os.RemoveAll(trainingName)
	}
	if _, err := os.Stat(testName); err == nil {
		os.RemoveAll(testName)
	}

	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}
	// Map IO for testing
	args := []string{"init", srcName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	// Add our records
	keys := []string{}
	for k, v := range testRecords {
		keys = append(keys, k)
		args = []string{"create", srcName, k}
		src, err := JSONMarshalIndent(v, "", "    ")
		if err != nil {
			t.Errorf("Can't marshal %q -> %+v", k, v)
			continue
		}
		in = bytes.NewBuffer(src)
		if err := RunCLI(in, out, eout, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
	// The keys will be read from the "in" for clone-sample.
	src := []byte(strings.Join(keys, "\n"))
	in = bytes.NewBuffer(src)
	args = []string{"clone-sample", "-size", "10", srcName, trainingName, trainingDsnURI, testName, testDsnURI}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
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

func TestCLIOnFrames(t *testing.T) {
	input, output, errout := []byte{}, []byte{}, []byte{}
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	eout := bytes.NewBuffer(errout)

	cName := path.Join("testout", "CLI_on_frames.ds")

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

	output = []byte{}
	out = bytes.NewBuffer(output)

	args = []string{"keys", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("expected to get a list of keys, %s", err)
		t.FailNow()
	}
	// Get keys for all of collection.
	src, err := ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("failed to read output for keys, %s", err)
		t.FailNow()
	}
	keys = strings.Split(string(src), "\n")
	if len(keys) == 0 {
		t.Errorf("expected a list of keys to frame, got none")
		t.FailNow()
	}

	frameName := "one-data"
	frameFile := path.Join(cName, "_frames", frameName+".json")
	args = []string{"frame", "-i", "-", cName, frameName, ".one=one"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	input = []byte(strings.Join(keys, "\n"))
	in = bytes.NewBuffer(input)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	if _, err := os.Stat(frameFile); err != nil {
		t.Errorf("failed to create frame %q at %q, %s", frameName, frameFile, err)
		t.FailNow()
	}

	// List frames
	args = []string{"frames", cName}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	expectedS := "one-data"
	src, err = ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("failed to read frames, %s", err)
		t.FailNow()
	}
	gotS := fmt.Sprintf("%s", src)
	if gotS != "one-data" {
		t.Errorf("Expected frame list %q, got %q", expectedS, gotS)
		t.FailNow()
	}

	// get keys from frame
	args = []string{"frame-keys", cName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("failed to read output of frame-keys, %s", err)
		t.FailNow()
	}
	okeys := strings.Split(string(src), "\n")
	if len(okeys) == 0 {
		t.Errorf("failed to read keys from frame-keys")
		t.FailNow()
	}

	// get definition from frame
	args = []string{"frame-def", cName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("failed to read frame-def, %s", err)
		t.FailNow()
	}
	if len(src) == 0 {
		t.Errorf("Failed to get frame definition")
		t.FailNow()
	}

	// get objects in frame
	args = []string{"frame-objects", cName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(out)
	if err != nil {
		t.Errorf("failed to read output of frame-objects, %s", err)
		t.FailNow()
	}
	if len(src) == 0 {
		t.Errorf("Failed to get frame objects")
		t.FailNow()
	}

	// refresh frame
	args = []string{"refresh", cName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		t.FailNow()
	}

	keys = okeys[0:3]

	// reframe
	args = []string{"reframe", cName, "one-data"}
	input = []byte(strings.Join(keys, "\n"))
	in = bytes.NewBuffer(input)
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		t.FailNow()
	}

	// delete frame
	/*
		args = []string{"delete-frame", cName, "one-data"}
		output = []byte{}
		out = bytes.NewBuffer(output)
		if err := RunCLI(in, out, eout, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	*/

}

func TestCLIOnAttachments(t *testing.T) {
	var (
		input, output []byte
	)
	input = []byte(`{
	"one": 1,
	"two": 2,
	"three": 3
}`)
	// Make test attachment file.
	attachmentFile := path.Join("testout", "test_attachment.txt")
	src := []byte(`
This is a plain text file for attachment.

1. One
2. Two
3. Three

The End.
`)
	if err := ioutil.WriteFile(attachmentFile, src, 0664); err != nil {
		t.Errorf("unable to create test attachment %q, %s", attachmentFile, err)
		t.FailNow()
	}
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)

	// Setup clean cName
	cName := path.Join("testout", "attached.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	args := []string{"init", cName, ""}
	if err := RunCLI(os.Stdin, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"create", cName, "uno"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"attach", cName, "uno", attachmentFile}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"attachments", cName, "uno"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	os.Rename(attachmentFile, attachmentFile+".bak")
	args = []string{"retrieve", cName, "uno", attachmentFile}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"prune", cName, "uno", attachmentFile}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
}

func TestCheckRepair(t *testing.T) {
	cName := path.Join("testout", "myfix.ds")
	csvName := path.Join("testout", "myfix.csv")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	if _, err := os.Stat(csvName); err == nil {
		os.RemoveAll(csvName)
	}

	// Setup test data
	data := map[string]map[string]interface{}{
		"freda": map[string]interface{}{
			"Name":   "Little Freda",
			"Email":  "freda@inverness.example.edu",
			"Office": "4th Tower",
			"Count":  1,
		},
		"mojo": map[string]interface{}{
			"Name":   "Mojo Same",
			"Email":  "mojo.sam@sams-cafe.example.org",
			"Office": "At the Piano",
			"Count":  2,
		},
	}

	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("Failed to create %q, %s", cName, err)
		t.FailNow()
	}
	for k, v := range data {
		if err := c.Create(k, v); err != nil {
			t.Errorf("Failed to setup record %q in %q -> %+v", k, cName, v)
			t.FailNow()
		}
	}
	expected64 := int64(2)
	got64 := c.Length()
	if expected64 != got64 {
		t.Errorf("Expected %d, got %d for count of test data", expected64, got64)
		t.FailNow()
	}
	c.Close()

	// IO Setup
	var (
		input  []byte
		output []byte
	)
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)

	// Run tests
	args := []string{"check", cName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("check should have been OK, %q, %s", strings.Join(args, " "), err)
	}

	// Test case of missing collections.json
	os.RemoveAll(path.Join(cName, "collection.json"))
	args = []string{"check", cName}
	if err := RunCLI(in, out, os.Stderr, args); err == nil {
		t.Errorf("check should have failed, %q", strings.Join(args, " "))
	}

	// Initiating a repair
	args = []string{"repair", cName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("repair should have succeeded, %q, %s", strings.Join(args, " "), err)
	}

	// Repair should have worked
	args = []string{"check", cName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("check should confirm repair worked, %q, %s", strings.Join(args, " "), err)
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
	src := []byte(`{"two": 2}`)
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

	src = []byte(`{"three": 3}`)
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
	src = []byte(`{"four": 4}`)
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

	expectedB = []byte(`{
    "one": 1
}`)
	key = "1"
	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
		t.Errorf("expected %q, got %q", expectedB, gotB)
		t.FailNow()
	}

	expectedB = []byte(`{
    "two": 2
}`)
	key = "2"
	args = []string{"read", cName, key}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("ran %q, got error %s", strings.Join(args, " "), err)
		t.FailNow()
	}
	gotB, _ = ioutil.ReadAll(out)
	if bytes.Compare(expectedB, gotB) != 0 {
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

func TestMigrateV1ToV2(t *testing.T) {
	srcName := path.Join("testout", "src_v1.ds")
	dstName := path.Join("testout", "dst_v2.ds")

	if _, err := os.Stat(srcName); err == nil {
		os.RemoveAll(srcName)
	}
	if _, err := os.Stat(dstName); err == nil {
		os.RemoveAll(dstName)
	}

	records := map[string]map[string]interface{}{
		"t1": {
			"one":   1,
			"two":   "two",
			"three": true,
		},
		"t2": {
			"two":   2,
			"three": false,
		},
		"t3": {
			"one":   3,
			"two":   3,
			"three": true,
		},
	}

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	if err := dsv1.SetupV1TestCollection(srcName, records); err != nil {
		t.Errorf("failed to setup v1 source collection, %s", err)
	}

	args := []string{"init", dstName, ""}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create empty destination collection, %s", err)
		t.FailNow()
	}

	args = []string{"migrate", srcName, dstName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("migration errors for %q to %q, %s", srcName, dstName, err)
		t.FailNow()
	}
}

func TestKeySample(t *testing.T) {
	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
	}
	dName := "testout"
	cName := path.Join(dName, "test_key_sample.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, "")
	if err != nil {
		t.Errorf("failed to setup %q, %s", cName, err)
		t.FailNow()
	}
	for key, val := range testRecords {
		if err := c.CreateObject(key, val); err != nil {
			t.Errorf("failed to setup record %q in %q, %s", key, cName, err)
			t.FailNow()
		}
	}
	// Close the collection so we can test from the runner
	c.Close()

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	args := []string{"keys", "-sample", "3", cName}

	//	t.Errorf(`Need a  test for "dataset keys -sample 3 testout/zbs2.ds"`)
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("key sample failed, %s", err)
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
	args = []string{"has-key", cName, key}
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

func TestCliAttachments(t *testing.T) {
	dName := "testout"
	cName := path.Join(dName, "test_mydata.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q, %s", cName, err)
		t.FailNow()
	}
	records := []string{
		`{"Name":"freda","EMail":"freda@inverness.example.edu","Office":"4th Tower","Count":1}`,
		`{"Name": "mojo", "EMail": "mojo.sam@sams-place.example.org", "Office": "piano", "Count": 2 }`,
	}

	key := "freda"
	filename := path.Join(dName, "freda.csv")
	src := []byte(`Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
`)
	if _, err := os.Stat(filename); err == nil {
		os.RemoveAll(filename)
	}
	if err := ioutil.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("failed to write test file %q, %s", filename, err)
		t.FailNow()
	}
	args = []string{"create", cName, key, records[0]}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q into %q, %s", filename, cName, err)
		t.FailNow()
	}
	args = []string{"attach", cName, key, filename}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to attach %q into %q in %q, %s", filename, key, cName, err)
		t.FailNow()
	}

	key = "mojo"
	filename = path.Join(dName, "mojo.csv")
	src = []byte(`Name,EMail,Office,Count
mojo,mojo.sam@sams-place.example.org,piano,2
`)
	if _, err := os.Stat(filename); err == nil {
		os.RemoveAll(filename)
	}
	if err := ioutil.WriteFile(filename, src, 0664); err != nil {
		t.Errorf("failed to write test file %q, %s", filename, err)
		t.FailNow()
	}
	args = []string{"create", cName, key, records[1]}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q into %q, %s", filename, cName, err)
		t.FailNow()
	}

	args = []string{"attach", cName, key, filename}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to attach %q into %q in %q, %s", filename, key, cName, err)
		t.FailNow()
	}

	args = []string{"retrieve", "-o", "-", cName, key, filename}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to retrieve %q from %q in %q, %s", filename, key, cName, err)
		t.FailNow()
	}
	src, _ = ioutil.ReadAll(out)
	if len(src) == 0 {
		t.Errorf("failed to retrieve %q from %q in %q, %s", filename, key, cName, "length is zero")
		t.FailNow()
	}

	key = "freda"
	filename = "freda.csv"
	args = []string{"prune", cName, key, filename}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to prune %q from %q in %q, %s", filename, key, cName, err)
		t.FailNow()
	}
}

func testCount(t *testing.T) {
	dName := "testout"
	cName := path.Join(dName, "test_count.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}

	in := bytes.NewBuffer([]byte{})
	out := bytes.NewBuffer([]byte{})
	eout := bytes.NewBuffer([]byte{})

	args := []string{"init", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to create %q, %s", cName, err)
		t.FailNow()
	}
	records := map[string]string{
		"freda": `{"Name":"freda","EMail":"freda@inverness.example.edu","Office":"4th Tower","Count":1,"published": true}`,
		"mojo":  `{"Name": mojo, "EMail": mojo.sam@sams-place.example.org, "Office": piano, "Count": 2,"published": false }`,
	}

	// Setup our test records
	for key, val := range records {
		args = []string{"create", cName, key, val}
		if err := RunCLI(in, out, eout, args); err != nil {
			t.Errorf("failed to create %q in %q, %s", key, cName, err)
			t.FailNow()
		}
	}

	// Set actual count of records
	args = []string{"count", cName}
	if err := RunCLI(in, out, eout, args); err != nil {
		t.Errorf("failed to count records in %q, %s", cName, err)
		t.FailNow()
	}
}
