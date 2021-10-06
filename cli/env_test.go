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
	"os"
	"testing"
)

func TestAppEnv(t *testing.T) {
	app := NewCli(Version)
	if app == nil {
		t.Errorf("Expected an 'App' struct, got nil")
		t.FailNow()
	}

	userName := os.Getenv("USER")
	usage := "set USER from the environment"
	err := app.EnvStringVar(&userName, "USER", userName, usage)
	if err != nil {
		t.Errorf("%s", err)
	}
	gotS := app.Env("USER")
	expectedS := usage
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
	app.StringVar(&userName, "u2,user2", expectedUserS, expectedS)
	gotS = app.Option("u2")
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
