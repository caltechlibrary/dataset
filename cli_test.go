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
	args = []string{"frame-keys", srcName, "one-data"}
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
	args = []string{"frame-def", srcName, "one-data"}
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
	args = []string{"frame-objects", srcName, "one-data"}
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

	// reframe
	args = []string{"reframe", srcName, "one-data", "field_one=.one"}
	input = []byte{}
	in = bytes.NewBuffer(input)
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}

	// refresh frame
	args = []string{"refresh", srcName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}

	// delete frame
	args = []string{"delete-frame", srcName, "one-data"}
	output = []byte{}
	out = bytes.NewBuffer(output)
	if err := RunCLI(in, out, os.Stderr, args); err != nil {
		t.Errorf("unexpected error when running %q, %s", strings.Join(args, " "), err)
	}

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

/**FIXME: convert this Bash script into Go based cli testing
#!/bin/bash

function assert_exists() {
	if [ "$#" != "2" ]; then
		echo "wrong number of parameters for $1, $*"
		exit 1
	fi
	if [[ ! -f "$2" && ! -d "$2" ]]; then
		echo "$1: $2 does not exists"
		exit 1
	fi

}

function assert_equal() {
	if [ "$#" != "3" ]; then
		echo "wrong number of parameters for $1, $*"
		exit 1
	fi
	if [ "$2" != "$3" ]; then
		echo "$1: expected |$2| got |$3|"
		exit 1
	fi
}

function fail_now() {
    echo $@
    exit 1
}

#
# Tests
#
function test_dataset() {
	if [[ -f "testdata/test1.ds/collection.json" ]]; then
		rm -fR testdata/test1.ds
	fi
	EXT=".exe"
	OS=$(uname)
	if [ "$OS" != "Windows" ]; then
		EXT=""
	fi
	echo "checking for bin/dataset${EXT}"
	if [[ ! -f "bin/dataset${EXT}" || ! -f "cmd/dataset/assets.go" ]]; then
		# We need to build
		pkgassets -o cmd/dataset/assets.go \
			-p main -ext=".md" -strip-prefix="/" \
			-strip-suffix=".md" \
			Examples examples/dataset \
			Help docs/dataset
		go build -o "bin/dataset${EXT}" cmd/dataset/dataset.go cmd/dataset/assets.go
	fi

	# Test init
	EXPECTED="OK"
	RESULT=$(bin/dataset init testdata/test1.ds)
	assert_equal "init testdata/test1.ds" "$EXPECTED" "$RESULT"
	export DATASET="testdata/test1.ds"

	assert_exists "collection create" "testdata/test1.ds"
	assert_exists "collection created metadata" "testdata/test1.ds/collection.json"


	# Test create
	if ! bin/dataset create "${DATASET}" 1 '{"one":1}'; then
       fail_now "Failed to create record 1 in ${DATASET}"
    fi
	if ! echo '{"two":2}' | bin/dataset create -i - "${DATASET}" 2; then
       fail_now "Failed to create record 2 in ${DATASET}"
    fi
	echo '{"three":3}' >"testdata/test3.json"
	if ! bin/dataset create -i testdata/test3.json "${DATASET}" 3; then
       fail_now "Failed to create record 3 in ${DATASET}"
    fi
	echo '{"four":4}' >"testdata/test4.json"
    if ! bin/dataset create "${DATASET}" 4 testdata/test4.json; then
       fail_now "Failed to create record 4 in ${DATASET}"
    fi

	# Test read
	EXPECTED='{"_Key":"1","one":1}'
	RESULT=$(bin/dataset read "${DATASET}" 1)
	assert_equal "read 1:" "$EXPECTED" "$RESULT"
	EXPECTED='{"_Key":"2","two":2}'
	RESULT=$(echo -n '2' | bin/dataset -i - read "${DATASET}")
	assert_equal "read 2:" "$EXPECTED" "$RESULT"


	# Test keys
	EXPECTED="1 2 3 4 "
	RESULT=$(bin/dataset keys "${DATASET}" | sort | tr "\n" " ")
	assert_equal "keys:" "$EXPECTED" "$RESULT"
	EXPECTED="1 "
	RESULT=$(bin/dataset keys "${DATASET}" '(eq .one 1)' | sort | tr "\n" " ")

	if [ -f "testdata/test1.ds/collection.json" ]; then
		rm -fR testdata/test1.ds
	fi
	echo "test_dataset, OK"
}


function test_issue19() {
	echo "test_issue19"
	if [[ -d "testdata/test_issue19.ds" ]]; then
		rm -fR testdata/test_issue19.ds
	fi
	bin/dataset -nl=false -quiet init "testdata/test_issue19.ds"
	bin/dataset -nl=false -quiet create testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":1}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record in an empty collection"
		exit 1
	fi

	# Now try creating the record again without -overwrite
    echo "NOTE: expecting an error message on next line"
	bin/dataset -nl=false -quiet create testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":2}'
	if [[ "$?" != "1" ]]; then
        echo "Expected return value 1, got $?"
		echo "Failed, should NOT be able to create the record when it exists in a collection without -overwrite"
		exit 1
	fi

	# Now try to create the record with -overwrite
	bin/dataset -nl=false -quiet create -overwrite testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":3}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record with -overwrite!"
		exit 1
	fi

	echo "test_issue19, OK"
	rm -fR "testdata/test_issue19.ds"
}

function test_readme () {
    echo "test_readme"
    mkdir -p testdata
    if [[ -d "testdata/mystuff.ds" ]]; then
        rm -fR testdata/mystuff.ds
    fi

    # Create a collection "testdata/mystuff.ds", the ".ds" lets the bin/dataset command know that's the collection to use.
    bin/dataset -quiet -nl=false init testdata/mystuff.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (206): could not init mystuff.ds'
        exit 1
    fi
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document
    bin/dataset -quiet -nl=false create testdata/mystuff.ds freda '{"name":"freda","email":"freda@inverness.example.org"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (216): could not create freda.json'
        exit 1
    fi
    # If successful then you should see an OK otherwise an error message

    # Make sure we have a record called freda
    bin/dataset -quiet -nl="false" haskey testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (222): (failed) testdata/mystuff.ds haskey freda'
        exit 1
    fi


    # Read a JSON document
    bin/dataset -quiet -nl=false read testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (232): could not read freda.json'
        exit 1
    fi

    # Path to JSON document
    bin/dataset -quiet -nl=false path testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (237): could not path freda.json'
        exit 1
    fi

    # Update a JSON document
    bin/dataset -quiet -nl=false update testdata/mystuff.ds freda '{"name":"freda","email":"freda@zbs.example.org", "count": 2}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (244): could not update freda.json'
        exit 1
    fi

    # If successful then you should see an OK or an error message

    # List the keys in the collection
    bin/dataset -quiet -nl=false keys testdata/mystuff.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (253): could not keys'
        exit 1
    fi

    # Get keys filtered for the name "freda"
    bin/dataset -nl=false -quiet keys testdata/mystuff.ds '(eq .name "freda")' > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (260): could not keys'
        exit 1
    fi

    # Join freda-profile.json with "freda" adding unique key/value pairs
    cat << EOT > testdata/freda-profile.json
{"name": "little freda", "office": "SFL", "count": 3}
EOT

    bin/dataset -quiet -nl=false join testdata/mystuff.ds freda testdata/freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (271): could not join update'
        exit 1
    fi

    # Join freda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from freda-profile.json
    cat << EOT > testdata/freda-profile.json
{"name": "little freda", "office": "SFL", "count": 4}
EOT

    bin/dataset -quiet -nl=false join -overwrite testdata/mystuff.ds freda testdata/freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (283): could not join overwrite'
        exit 1
    fi


    # Delete a JSON document
    bin/dataset -quiet -nl=false delete testdata/mystuff.ds freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (290): could not join overwrite'
        exit 1
    fi

    # Import from a CSV file
    cat << EOT > testdata/my-data.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    bin/dataset -quiet -nl=false import testdata/mystuff.ds testdata/my-data.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (302): (failed) import testdata/mystuff.ds testdata/my-data.csv 1'
        exit 1
    fi

    echo "test_readme, OK"
    # To remove the collection just use the Unix shell command
    rm -fR testdata/mystuff.ds
    rm testdata/freda-profile.json
    rm testdata/my-data.csv
}

function test_getting_started() {
    echo "test_getting_started"
    if [[ -d "testdata/FavoriteThings.ds" ]]; then
        rm -fR testdata/FavoriteThings.ds
    fi
    bin/dataset -quiet -nl=false init testdata/FavoriteThings.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not init testdata/FavoriteThings.ds'
        exit 1
    fi

    bin/dataset -quiet -nl=false create testdata/FavoriteThings.ds beverage '{"thing":"coffee"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds create beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false read testdata/FavoriteThings.ds beverage > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds read beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false keys testdata/FavoriteThings.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds keys'
        exit 1
    fi

    cat << EOT > testdata/jazz-notes.json
{
    "songs": ["Blue Rondo al la Turk", "Bernie's Tune", "Perdido"],
    "pianist": [ "Dave Brubeck" ],
    "trumpet": [ "Dirk Fischer", "Dizzy Gillespie" ]
}
EOT
    bin/dataset -quiet -nl=false create testdata/FavoriteThings.ds "jazz-notes" testdata/jazz-notes.json
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not create jazz-notes'
        exit 1
    fi

    bin/dataset -quiet -nl=false keys testdata/FavoriteThings.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not keys'
        exit 1
    fi

    bin/dataset -quiet -nl=false read testdata/FavoriteThings.ds beverage jazz-notes > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not read multiple keys'
        exit 1
    fi

    # Cleanup after tests
    rm -fR testdata/FavoriteThings.ds
    rm testdata/jazz-notes.json
    echo "test_getting_started, OK"
}

function test_attachments() {
    echo 'test_attachments'
    if [[ -d "testdata/mydata.ds" ]]; then
        rm -fR testdata/mydata.ds
    fi
    bin/dataset -quiet -nl=false init testdata/mydata.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (375): could not testdata/mydata.ds init'
        exit 1
    fi

    cat << EOT > testdata/freda.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    cat << EOT > testdata/mojo.csv
Name,EMail,Office,Count
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT

    bin/dataset -quiet -nl=false import testdata/mydata.ds testdata/freda.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (389): (failed) testdata/mydata.ds import-csv testdata/freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false attach testdata/mydata.ds freda testdata/freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (395): (failed) testdata/mydata.ds attach freda testdata/freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/mydata.ds testdata/mojo.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (399): (failed) testdata/mydata.ds import-csv testdata/mojo.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false attach testdata/mydata.ds mojo testdata/mojo.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (404): (failed) testdata/mydata.ds attach testdata/mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false attachments testdata/mydata.ds mojo > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (410): (failed) attachments testdata/mydata.ds mojo'
        exit 1
    fi
    if [[ -f "mojo.csv" ]]; then
        rm mojo.csv
    fi
    bin/dataset -quiet -nl=false detach testdata/mydata.ds mojo mojo.csv > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (417): (failed) detach testdata/mydata.ds mojo mojo.csv'
        exit 1
    fi
    if [[ ! -f "mojo.csv" ]]; then
        echo 'test_attachments (417): (failed) detatch testdata/mydata.ds mojo mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false prune testdata/mydata.ds freda freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (426): (failed) prune testdata/mydata.ds freda freda.csv'
        exit 1
    fi

    # Success, cleanup our test data
    if [[ -f fred.csv ]]; then
        rm freda.csv
    fi
    if [[ -f mojo.csv ]]; then
        rm mojo.csv
    fi
    rm -fR testdata/mydata.ds
	echo "test_attachments, OK"
}

function test_check_and_repair() {
    echo 'test_check_and_repair'
    if [[ -d "testdata/myfix.ds" ]]; then
        rm -fR testdata/myfix.ds
    fi
    cat << EOT > testdata/myfix.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT
    bin/dataset -quiet -nl=false init testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) init testdata/myfix.ds'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/myfix.ds testdata/myfix.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) import testdata/myfix.ds testdata/myfix.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false check testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi

    CNT=$(bin/dataset count testdata/myfix.ds)
    echo "NOTE: expecting ${CNT} warnings detected on following line"
    echo '{}' > testdata/myfix.ds/collection.json
    bin/dataset -quiet -nl=false check testdata/myfix.ds
    if [[ "$?" == "0" ]]; then
        echo 'test_check_and_repair: (failed, expected exit code 1) testdata/myfix.ds check'
        exit 1
    fi
    echo "NOTE: Initiating a repair"
    bin/dataset -quiet -nl=false repair testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) repair testdata/myfix.ds'
        exit 1
    fi
    echo "NOTE: Expecting OK on next line if repair worked"
    bin/dataset check testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi

    echo "NOTE: Final Check, expecting OK"
    bin/dataset check testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi


    # Success, cleanup
    rm -fR testdata/myfix.ds
    rm testdata/myfix.csv
    echo 'test_check_and_repair, OK'
}

function test_count() {
    echo 'test_count'
    if [[ -d "testdata/count.ds" ]]; then
        rm -fR testdata/count.ds
    fi
    cat << EOT > testdata/count.csv
Name,EMail,Office,Count,published
freda,freda@inverness.example.edu,4th Tower,1,true
mojo,mojo.sam@sams-splace.example.org,piano,2,false
EOT

    if [[ ! -f "testdata/count.csv" ]]; then
        echo 'test_count: (failed) could not create testdata/count.csv'
        exit 1
    fi

    bin/dataset -quiet -nl=false init testdata/count.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) init testdata/count.ds'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/count.ds testdata/count.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) import testdata/count.ds testdata/count.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false count testdata/count.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) testdata/count.ds count'
        exit 1
    fi
    bin/dataset -quiet -nl=false count testdata/count.ds '(eq .published true)' > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count testdata/count.ds "(eq .published true)"'
        exit 1
    fi

    # Success, cleanup
    rm -fR testdata/count.ds
    rm testdata/count.csv
    echo 'test_count, OK'
}


function test_import_export() {
    echo 'test_import_export'
    if [[ -d "testdata/pubs.ds" ]]; then
        rm -fR "testdata/pubs.ds"
    fi
    cat << EOT > testdata/in.csv
id,title,type,date_type,date
44088,Application of a laser induced fluorescence model to the numerical simulation of detonation waves in hydrogen-oxygen-diluent mixtures,article,published,2014-04-04
46001,Leaderless Deterministic Chemical Reaction Networks,book_section,published,2013
61958,"Observation of the 14 MeV resonance in ^(12)C(p, p)^(12)C with molecular ion beams",article,published,1972-02-21
21682,Effect of scale size on a rocket engine with suddenly frozen nozzle flow,article,published,1961-03
39470,A Review of the Dynamics of Cavitating Pumps,article,published,2012-11-26
62459,An Experimental Investigation of the Flow over Blunt-Nosed Cones at a Mach Number of 5.8,monograph,completed,1956-06-15
80289,CIT-4: The first synthetic analogue of brewsterite,article,published,1997-08
80630,The Allocation of a Shared Resource Within an Organization,monograph,completed,1995-01
8488,Non-Gaussian covariance of CMB B modes of polarization and parameter degradation,article,published,2007-04-15
EOT

    bin/dataset -quiet -nl=false init testdata/pubs.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) init testdata/pubs.ds'
        exit 1
    fi

    bin/dataset -quiet -nl=false import testdata/pubs.ds testdata/in.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) testdata/pubs.ds import-csv testdata/in.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false keys testdata/pubs.ds >/dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) testdata/pubs.ds keys'
        exit 1
    fi
    #FIXME: export uses a frame to define exported content
    bin/dataset -quiet -nl=false frame -all testdata/pubs.ds outframe \
         "._Key=EPrint ID" ".title=Title" ".type=Type" \
         ".date_type=Date Type" ".date=Date" > /dev/null
    bin/dataset -quiet -nl=false export testdata/pubs.ds outframe "testdata/out.csv"
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) export testdata/pubs.ds outframe testdata/out.csv'
        exit 1
    fi

    # Success, cleanup
    rm -fR testdata/data.ds
    rm testdata/in.csv
    rm testdata/out.csv
    echo 'test_import_export, OK'
}

function test_sync() {
	echo "test_sync()"
	mkdir -p testdata
	cat << EOF > testdata/expected.csv
id,one,two,three,four,five
0,A,B,C,D,E
1,B,C,D,E,F
2,C,D,E,F,G
3,D,E,F,G,H
4,E,F,G,H,I
EOF

	cat << EOF > testdata/initial.csv
id,one,two
0,A,B
1,B,C
2,C,D
3,D,E
4,E,F
EOF

	if [[ -d testdata/merge4.ds ]]; then
		rm -fR testdata/merge4.ds
	fi
	bin/dataset -quiet -nl=false init testdata/merge4.ds
	bin/dataset -quiet -nl=false import testdata/merge4.ds testdata/initial.csv 1
	bin/dataset -quiet -nl=false frame -a testdata/merge4.ds f4 "._Key=id" ".one=one" ".two=two" ".three=three" ".four=four" ".five=five" >/dev/null

	# Now generate an updated result CSV
	cp testdata/initial.csv testdata/result.csv
	cat testdata/expected.csv |\
        bin/dataset -quiet -nl=false sync-send \
            -i - testdata/merge4.ds f4 \
            >testdata/result.csv

    #FIXME: need to check to see if our tables make sense
    T=$(diff testdata/expected.csv testdata/result.csv)
    if [[ "$?" != "0" ]]; then
        echo "Diff returned: $?"
        echo "Diff found: $T"
        #diff testdata/expected.csv testdata/result.csv
		exit 1
    fi
    if [[ "$T" != "" ]]; then
        echo "Diff found: $T"
        exit 1
    fi
	echo "test_sync, OK"
}


echo "Testing command line tools"
test_dataset
test_issue19
test_readme
test_getting_started
test_attachments
test_count
test_import_export
#NOTE: test will be skip if there is no etc/client_secret.json found
#test_gsheet credentials.json # etc/client_secret.json
test_check_and_repair
test_sync
echo 'PASS'
echo "Ok $(basename "$0")"
*/
