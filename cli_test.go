package dataset

import (
	"bytes"
	"flag"
	"os"
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
func TestRunCLI(t *testing.T) {
	var (
		input, output []byte
	)
	// Map IO for testing
	in := bytes.NewBuffer(input)
	out := bytes.NewBuffer(output)
	// Setup command line args

	// Check if version, license, help returns anything
	for _, arg := range []string{"help", "init", "create", "read", "update", "delete", "keys", "has-key", "frames", "frame", "frame-objects", "frame-def", "refresh", "reframe", "delete-frame", "attachments", "attach", "retrieve", "prune"} {
		args := []string{arg}
		if err := RunCLI(in, out, os.Stderr, args); err != nil {
			t.Errorf("unexpected error when running %q, %s", arg, err)
		}
	}
	// FIXME: Run through command sequences
	t.FailNow()
}
