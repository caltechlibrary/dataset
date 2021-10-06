/**
 * cli is a package intended to encourage some standardization in the command line user interface for programs
 * developed for Caltech Library.
 *
 * @author R. S. Doiel, <rsdoiel@caltech.edu>
 *
 * Copyright (c) 2021, Caltech
 * All rights not granted herein are expressly reserved by Caltech.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */
package cli

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
)

var (
	appName = "testcli"
)

func isSameString(expected, value string) string {
	if strings.Compare(expected, value) == 0 {
		return ""
	}
	return fmt.Sprintf("expected %q, got %q", expected, value)
}

func isSameInt(expected, value int) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("expected %d, got %d", expected, value)
}

func check(t *testing.T, msg string, failNow bool) {
	if len(msg) > 0 {
		t.Errorf(msg)
	}
}

func TestOpen(t *testing.T) {
	// Check of use of fallbackFile (os.Stdout) and cli.Open()
	fp, err := Open("", os.Stdout)
	if err != nil {
		t.Errorf("Should have fallen back to os.Stdout, got error %s", err)
		t.FailNow()
	}
	if fp != os.Stdout {
		t.Errorf("fp should be pointing at os.Stdout")
		t.FailNow()
	}
	fp, err = Open("-", os.Stdout)
	if err != nil {
		t.Errorf("Should have fallen back to os.Stdout, got error %s", err)
		t.FailNow()
	}
	if fp != os.Stdout {
		t.Errorf("fp should be pointing at os.Stdout")
		t.FailNow()
	}
	// Check if we can open an existing file, i.e. README.md
	fp, err = Open("README.md", os.Stdout)
	if err != nil {
		t.Errorf("Should have gotten pointer to README.md, got error %s", err)
		t.FailNow()
	}
	if fp == os.Stdout {
		t.Errorf("fp should be pointing at README.md NOT os.Stdout")
		t.FailNow()
	}
	if fp == nil {
		t.Errorf("fp should be pointing at README.md NOT nil")
		t.FailNow()
	}
	fp.Close()

	// Check to see if we fallack to os.Stdout for Create()
	fp, err = Create("", os.Stdout)
	if err != nil {
		t.Errorf("Should have fallen back to os.Stdout, got error %s", err)
		t.FailNow()
	}
	if fp != os.Stdout {
		t.Errorf("fp should be pointing at os.Stdout")
		t.FailNow()
	}
	fp, err = Create("-", os.Stdout)
	if err != nil {
		t.Errorf("Should have fallen back to os.Stdout, got error %s", err)
		t.FailNow()
	}
	if fp != os.Stdout {
		t.Errorf("fp should be pointing at os.Stdout")
		t.FailNow()
	}
	// Check to see if we can open test.txt with Create()
	fp, err = Create("test.txt", os.Stdout)
	if err != nil {
		t.Errorf("Should have gotten pointer to test.txt, got error %s", err)
		t.FailNow()
	}
	if fp == os.Stdout {
		t.Errorf("fp should be pointing at test.txt NOT os.Stdout")
		t.FailNow()
	}
	if fp == nil {
		t.Errorf("fp should be pointing at test.txt NOT nil")
		t.FailNow()
	}
	fp.Close()
	os.Remove("test.txt")
}

func TestShiftArg(t *testing.T) {
	args := []string{
		"one",
		"two",
		"three",
	}
	expected := args[0]
	expectedL := len(args) - 1
	r, args := ShiftArg(args)
	if r != expected {
		t.Errorf("expected %q, got %q", expected, r)
	}
	if len(args) != expectedL {
		t.Errorf("expected args length to be %d, got %d", expectedL, len(args))
	}
	expected = args[0]
	expectedL = len(args) - 1
	r, args = ShiftArg(args)
	if r != expected {
		t.Errorf("expected %q, got %q", expected, r)
	}
	if len(args) != expectedL {
		t.Errorf("expected args length to be %d, got %d", expectedL, len(args))
	}
	expected = args[0]
	expectedL = len(args) - 1
	r, args = ShiftArg(args)
	if r != expected {
		t.Errorf("expected %q, got %q", expected, r)
	}
	if args != nil {
		t.Errorf("expected args to be nil, got %+v", args)
	}
}

func TestApp(t *testing.T) {
	appName := path.Base(os.Args[0])

	app := NewCli(Version)
	if app == nil {
		t.Errorf("Expected an 'App' struct, got nil")
		t.FailNow()
	}

	version := app.Version()
	expectedS := fmt.Sprintf("%s %s", appName, Version)
	if expectedS != version {
		t.Errorf("expected %s, got %s", expectedS, version)
	}

	userName := os.Getenv("USER")
	usage := "set USER from the environment"
	err := app.EnvStringVar(&userName, "USER", userName, usage)
	if err != nil {
		t.Errorf("%s", err)
	}
	gotS := app.Env("USER")
	expectedS = usage
	if expectedS != gotS {
		t.Errorf("expected %q, got %q", expectedS, gotS)
	}
	err = app.ParseEnv()
	if err != nil {
		t.Errorf("expected ParseEnv() to return nil, got %s", err)
		t.FailNow()
	}
	// Now set a new default of "jane.doe"
	expectedUserS := "jane.doe"
	err = app.EnvStringVar(&userName, "USER", expectedUserS, usage)
	if err != nil {
		t.Errorf("EnvStringVar() returned an error, %s", err)
		t.FailNow()
	}

	e, err := app.EnvAttribute("USER")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	gotS = e.StringValue
	if expectedUserS != gotS {
		t.Errorf("expected %q, got %q", expectedUserS, gotS)
	}
	if userName != gotS {
		t.Errorf("expected %q, got %q", userName, gotS)
	}

	// We expected to get the current user when we ParseEnv
	expectedUserS = os.Getenv("USER")
	err = app.ParseEnv()
	if err != nil {
		t.Errorf("ParseEnv() returned an error, %s", err)
		t.FailNow()
	}
	// After ParseEnv(), userName should have been updated
	gotS = app.Getenv("USER")
	if expectedUserS != gotS {
		t.Errorf("expected %q, got %q", expectedUserS, gotS)
	}

	e, err = app.EnvAttribute("USER")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if expectedUserS != e.StringValue {
		t.Errorf("expected %q, got %q", expectedUserS, e.StringValue)
	}

	expectedUserS = "bessie.smith"
	expectedS = "set the username overridding the enviroment"
	app.StringVar(&userName, "u,user", expectedUserS, expectedS)
	gotS = app.Option("u")
	if expectedS != gotS {
		t.Errorf("expected %q, got %q", expectedS, gotS)
	}
	gotS = app.Option("user")
	if expectedS != gotS {
		t.Errorf("expected %q, got %q", expectedS, gotS)
	}
	if expectedUserS != userName {
		t.Errorf("expected %s, got %s", expectedUserS, userName)
	}
}

func TestMd2Man(t *testing.T) {
	s := "This is _in_line_ _italics_ in _Markdown_v1.0.0_"
	expected := false
	for i := 0; i < len(s); i++ {
		start := i
		end := i + 3
		if end >= len(s) {
			break
		}
		result := hasMidUnderscore(s, i)
		if i == 11 || i == 40 {
			expected = true
		} else {
			expected = false
		}
		if expected != result {
			t.Errorf("%d: %q <-- expected %t != %t", i, s[start:end], expected, result)
		}
	}
	src := []byte(`
This is a *test* of using _in_line_ formatting!
`)
	srcR := string(md2man(src))
	expectedS := `
This is a \fBtest\fP of using \fIin_line\fP formatting!
`
	if srcR != expectedS {
		t.Errorf("\n%q\n%q\n", expectedS, srcR)
	}

	s = `use --help EXAMPLE_NAME where EXAMPLE_NAME`
	expected = false
	for i := 0; i < len(s); i++ {
		result := hasMidUnderscore(s, i)
		if i == 18 || i == 37 {
			expected = true
		} else {
			expected = false
		}
		if expected != result {
			t.Errorf("%d: %q, expected %t, got %t", i, s, expected, result)
		}

	}

	src = []byte(`
To view a specific example use --help EXAMPLE_NAME where EXAMPLE_NAME
is one of the following: getting-started-with-dataset, cloning-and-samples,
`)

	expectedS = `
To view a specific example use --help EXAMPLE_NAME where EXAMPLE_NAME
is one of the following: getting-started-with-dataset, cloning-and-samples,
`
	srcR = string(md2man(src))
	if expectedS != srcR {
		t.Errorf("Expected %q, got %q", expectedS, srcR)
	}
}

func TestREMidUnderscoreMethod(t *testing.T) {
	re := regexp.MustCompile(`[[:alpha:]0-9]_[[:alpha:]0-9]`)
	s := "This is _in_line_ _italics_ in _Markdown_v1.0.0_"
	for i := 0; i < len(s); i++ {
		start := i
		end := i + 3
		if end >= len(s) {
			break
		}
		result := re.MatchString(s[start:end])
		expected := false
		if i == 10 || i == 39 {
			expected = true
		} else {
			expected = false
		}
		if result != expected {
			t.Errorf("%d: %q <-- %t != %t\n", i, s[start:end], result, expected)
		}
	}
}
