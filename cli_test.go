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
	t.Errorf("cli sample cloning not implemented")
}

func TestCLIOnFrames(t *testing.T) {
	t.Errorf("cli frames commands not implemented")
}

func TestCLIOnAttachments(t *testing.T) {
	t.Errorf("cli attachment command not implemented")
}
