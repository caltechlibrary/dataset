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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestDisplay(t *testing.T) {
	appName := "TestDisplay"
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)

	output := []byte{}
	out := bytes.NewBuffer(output)
	DisplayLicense(out, appName, License)
	if out.Len() == 0 {
		t.Errorf("DisplayLicense() failed, nothing written to output buffer")
	}
	output = []byte{}
	out = bytes.NewBuffer(output)
	DisplayUsage(out, appName, flagSet, "This is a description", "This is examples", "This would be license text")
	if out.Len() == 0 {
		t.Errorf("DisplayUsage() failed, nothing written to output buffer")
	}
}

func TestRunCLIOnCRUDL(t *testing.T) {
	var (
		input, output []byte
	)
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	// Cleanup stale test data
	cName := path.Join("testout", "C1.ds")
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
		if err := RunCLI(in, out, os.Stderr, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
}

func TestCloning(t *testing.T) {
	srcName := path.Join("testout", "src2.ds")
	dstName := path.Join("testout", "dst2.ds")
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
	srcName := path.Join("testout", "zbs_characters.ds")
	trainingName := path.Join("testout", "zbs_training.ds")
	trainingDsnURI := "sqlite://testout/zbs_training.ds/collection.db"
	testName := path.Join("testout", "zbs_test.ds")
	testDsnURI := "sqlite://testout/zbs_test.ds/collection.db"
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
	var (
		input, output []byte
	)
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	args := []string{"init", srcName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	// Add our records
	keys := []string{}
	for k, v := range testRecords {
		keys = append(keys, k)
		args = []string{"create", srcName, k}
		src, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			t.Errorf("Can't marshal %q -> %+v", k, v)
			continue
		}
		in = bytes.NewBuffer(src)
		if err := RunCLI(in, out, os.Stderr, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
	// The keys will be read from the "in" for clone-sample.
	src := []byte(strings.Join(keys, "\n"))
	in = bytes.NewBuffer(src)
	args = []string{"clone-sample", "-size", "10", srcName, trainingName, trainingDsnURI, testName, testDsnURI}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
}

func TestCLIOnFrames(t *testing.T) {
	srcName := path.Join("testout", "zbs_frames.ds")
	if _, err := os.Stat(srcName); err == nil {
		os.RemoveAll(srcName)
	}

	testRecords := map[string]map[string]interface{}{}
	testRecords["character:1"] = map[string]interface{}{
		"name": "Jack Flanders",
		"one":  1,
	}
	testRecords["character:2"] = map[string]interface{}{
		"name": "Little Frieda",
		"one":  2,
	}
	testRecords["character:3"] = map[string]interface{}{
		"name": "Mojo Sam the Yoodoo Man",
		"one":  3,
	}
	testRecords["character:4"] = map[string]interface{}{
		"name": "Kasbah Kelly",
		"one":  4,
	}
	testRecords["character:5"] = map[string]interface{}{
		"name": "Dr. Marlin Mazoola",
		"one":  3,
	}
	testRecords["character:6"] = map[string]interface{}{
		"name": "Old Far-Seeing Art",
		"one":  2,
	}
	testRecords["character:7"] = map[string]interface{}{
		"name": "Chief Wampum Stompum",
		"one":  1,
	}
	testRecords["character:8"] = map[string]interface{}{
		"name": "The Madonna Vampira",
		"one":  0,
	}
	testRecords["character:9"] = map[string]interface{}{
		"name": "Domenique",
		"one":  1,
	}
	testRecords["character:10"] = map[string]interface{}{
		"name": "Claudine",
		"one":  1,
	}
	var (
		input, output []byte
	)
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	args := []string{"init", srcName}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	// Add our records
	keys := []string{}
	for k, v := range testRecords {
		keys = append(keys, k)
		args = []string{"create", srcName, k}
		src, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			t.Errorf("Can't marshal %q -> %+v", k, v)
			continue
		}
		in = bytes.NewBuffer(src)
		if err := RunCLI(in, out, os.Stderr, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
		}
	}
	args = []string{"frame", srcName, "one-data", ".one"}
	input = []byte(strings.Join(keys, "\n"))
	in = bytes.NewBuffer(input)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}

	// List frames
	args = []string{"frames", srcName}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	expectedS := "one-data"
	gotS := strings.TrimSpace(fmt.Sprintf("%s", output))
	if gotS != "one-data" {
		t.Errorf("Expected %q, got %q", expectedS, gotS)
	}

	// get keys from frame
	args = []string{"frame-keys", srcName, "uno"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	okeys := strings.Split(fmt.Sprintf("%s", output), "\n")
	if len(okeys) == 0 {
		t.Errorf("failed to read keys from frame-keys")
	}

	// get definition from frame
	args = []string{"frame-def", srcName, "uno"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	if len(output) == 0 {
		t.Errorf("Failed to get frame definition")
	} else {
		fmt.Printf("DEBUG frame-def -> %q\n", output)
	}

	// get objects in frame
	args = []string{"frame-objects", srcName, "uno"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	if len(output) == 0 {
		t.Errorf("Failed to get frame objects")
	} else {
		fmt.Printf("DEBUG frame-objects -> %q\n", output)
	}

	// refresh frame
	// reframe
	// delete frame

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
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	cName := path.Join("testout", "attached.ds")
	args := []string{"init", cName}
	if err := RunCLI(os.Stdin, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"create", cName, "uno"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"attach", cName, "uno", "README.md"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"attachments", cName, "uno"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"retreive", cName, "uno", "README.md"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
	args = []string{"prune", cName, "uno"}
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}
}
