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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	showHelp bool
	appName  = path.Base(os.Args[0])
)

// DisplayHelp writes out help on a supported topic
func DisplayHelp(out io.Writer, eout io.Writer, topic string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	switch topic {
	case "usage":
		fmt.Fprintf(out, StringProcessor(m, CLIDescription))
	case "examples":
		fmt.Fprintf(out, StringProcessor(m, CLIExamples))
	case "init":
		fmt.Fprintf(out, StringProcessor(m, cliInit))
	case "create":
		fmt.Fprint(out, StringProcessor(m, cliCreate))
	case "read":
		fmt.Fprint(out, StringProcessor(m, cliRead))
	case "update":
		fmt.Fprint(out, StringProcessor(m, cliUpdate))
	case "delete":
		fmt.Fprint(out, StringProcessor(m, cliDelete))
	case "keys":
		fmt.Fprint(out, StringProcessor(m, cliKeys))
	case "has-key":
		fmt.Fprint(out, StringProcessor(m, cliHasKey))
	case "count":
		fmt.Fprint(out, StringProcessor(m, cliCount))
	case "versioning":
		fmt.Fprint(out, StringProcessor(m, cliVersioning))
	case "versions":
		fmt.Fprint(out, StringProcessor(m, cliVersioning))
	case "read-version":
		fmt.Fprint(out, StringProcessor(m, cliVersioning))
	default:
		fmt.Fprintf(eout, "Unable to find help on %q\n", topic)
	}
}

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string, license string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, StringProcessor(m, license))
}

// DisplayVersion returns the of the dataset application.
func DisplayVersion(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, StringProcessor(m, "{app_name} {version}\n"))
}

// DisplayUsage displays a usage message.
func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet, description string, examples string, license string) {
	// Replacable text vars
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	// Convert {app_name} and {version} in description
	fmt.Fprintf(out, StringProcessor(m, description))
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	fmt.Fprintf(out, StringProcessor(m, examples))
	DisplayLicense(out, appName, StringProcessor(m, license))
}

func doInit(out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		dsnURI string
	)
	flagSet := flag.NewFlagSet("init", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for init")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "init")
		return nil
	}
	switch {
	case len(args) == 2:
		cName, dsnURI = args[0], args[1]
	case len(args) == 1:
		cName, dsnURI = args[0], ""
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME [DSN_URI], got %q", strings.Join(args, " "))
	}
	c, err := Init(cName, dsnURI)
	if err == nil {
		defer c.Close()
	}
	return err
}

func doVersioning(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName      string
		versioning string
	)
	flagSet := flag.NewFlagSet("versioning", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "versioning")
	}
	switch {
	case len(args) == 2:
		cName, versioning = args[0], strings.ToLower(strings.TrimSpace(args[1]))
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME VERSIONING_SETTING, got %q", strings.Join(append([]string{appName, "versioning"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetVersioning(versioning)
}

func doVersions(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("versions", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "versioning")
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(append([]string{appName, "versions"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	versions, err := c.Versions(key)
	if err != nil {
		return fmt.Errorf("version errors for %q, %s", key, err)
	}
	for _, val := range versions {
		fmt.Fprintf(out, "%s\n", val)
	}
	return nil
}

func doReadVersion(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName   string
		key     string
		version string
		src     []byte
	)
	flagSet := flag.NewFlagSet("read-version", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "versioning")
	}
	switch {
	case len(args) == 3:
		cName, key, version = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME KEY VERSION, got %q", strings.Join(append([]string{appName, "read-version"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	switch c.StoreType {
	case PTSTORE:
		src, err = c.PTStore.ReadVersion(key, version)
	case SQLSTORE:
		src, err = c.SQLStore.ReadVersion(key, version)
	default:
		return fmt.Errorf("%q storage not supported", c.StoreType)
	}
	if err == nil {
		fmt.Fprintf(out, "%s\n", src)
	}
	return err
}

func doCreate(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
		src   []byte
		input string
		err   error
	)
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&input, "i", "-", "read JSON from file, use '-' for stdin")
	flagSet.StringVar(&input, "input", "-", "read JSON from file, use '-' for stdin")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "create")
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		if input == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(input)
		}
		if err != nil {
			return fmt.Errorf("could not read JSON file, %s", err)
		}
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY [JSON_SRC], got %q", strings.Join(append([]string{appName, "create"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	obj := map[string]interface{}{}
	if err := DecodeJSON(src, &obj); err != nil {
		return err
	}
	if err := c.Create(key, obj); err != nil {
		return err
	}
	return nil
}

func doRead(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("read", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "create")
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	src := []byte{}
	defer c.Close()
	switch c.StoreType {
	case PTSTORE:
		src, err = c.PTStore.Read(key)
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	if err == nil {
		fmt.Fprintf(out, "%s\n", src)
	}
	return err
}

func doUpdate(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
		src   []byte
		input string
		err   error
	)
	flagSet := flag.NewFlagSet("update", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&input, "i", "-", "read JSON from file, use '-' for stdin")
	flagSet.StringVar(&input, "input", "-", "read JSON from file, use '-' for stdin")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "create")
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		if input == "-" {
			src, err = ioutil.ReadAll(in)
		} else {
			src, err = ioutil.ReadFile(input)
		}
		if err != nil {
			return fmt.Errorf("could not read JSON file, %s", err)
		}
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY [JSON_SRC], got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	obj := map[string]interface{}{}
	if err := DecodeJSON(src, &obj); err != nil {
		return err
	}
	if err := c.Update(key, obj); err != nil {
		return err
	}
	return nil
}

func doDelete(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("delete", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "create")
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	switch c.StoreType {
	case PTSTORE:
		err = c.PTStore.Delete(key)
	case SQLSTORE:
		err = c.SQLStore.Delete(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	return err
}

func doKeys(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
	)
	flagSet := flag.NewFlagSet("keys", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "keys")
	}
	switch {
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	keys, err := c.Keys()
	if err != nil {
		return err
	}
	src, err := json.MarshalIndent(keys, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to encode keys, %s", err)
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}

func doHasKey(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("has-key", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "has-key")
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	if c.HasKey(key) {
		fmt.Fprintln(out, "true")
	} else {
		fmt.Fprintln(out, "false")
	}
	return nil
}

func doCount(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
	)
	flagSet := flag.NewFlagSet("count", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "count")
	}
	switch {
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	cnt := c.Length()
	fmt.Fprintf(out, "%d\n", cnt)
	return nil
}

func doClone(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		dstName   string
		dstDsnURI string
		keysName  string
		verbose   bool
		src       []byte
		err       error
	)
	flagSet := flag.NewFlagSet("clone", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for read")
	flagSet.BoolVar(&showHelp, "help", false, "help for read")
	flagSet.StringVar(&keysName, "i", "-", "filename to read keys from")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(out, eout, "clone")
	}
	switch {
	case len(args) == 2:
		srcName, dstName, dstDsnURI = args[0], args[1], ""
	case len(args) == 3:
		srcName, dstName, dstDsnURI = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTIONS] SRC_COLLECTION_NAME DEST_COLLECTION_NAME [DEST_DSN_URI], got %q", strings.Join(args, " "))
	}
	if keysName == "-" {
		src, err = ioutil.ReadAll(in)
		if err != nil {
			return fmt.Errorf("failed to read keys from stdin, %s", err)
		}
	} else {
		src, err = ioutil.ReadFile(keysName)
		if err != nil {
			return fmt.Errorf("failed to read keys from %q, %s", keysName, err)
		}
	}
	keys := strings.Split(fmt.Sprintf("%s", src), "\n")
	source, err := Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	if err := source.Clone(dstName, dstDsnURI, keys, verbose); err != nil {
		return fmt.Errorf("clone failed %s", err)
	}
	return nil
}

/// RunCLI implemented the functionlity used by the cli.
func RunCLI(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Missing parameters")
	}
	verb, args := args[0], args[1:]
	switch verb {
	case "help":
		if len(args) > 0 {
			DisplayHelp(out, eout, args[0])
			return nil
		}
		DisplayHelp(out, eout, "usage")
		return nil
	case "init":
		return doInit(out, eout, args)
	case "create":
		return doCreate(in, out, eout, args)
	case "read":
		return doRead(in, out, eout, args)
	case "update":
		return doUpdate(in, out, eout, args)
	case "delete":
		return doDelete(in, out, eout, args)
	case "keys":
		return doKeys(in, out, eout, args)
	case "has-key":
		return doHasKey(in, out, eout, args)
	case "count":
		return doCount(in, out, eout, args)
	case "frames":
		return fmt.Errorf("verb %q not implemented", verb)
	case "frame":
		return fmt.Errorf("verb %q not implemented", verb)
	case "frame-def":
		return fmt.Errorf("verb %q not implemented", verb)
	case "frame-objects":
		return fmt.Errorf("verb %q not implemented", verb)
	case "refresh":
		return fmt.Errorf("verb %q not implemented", verb)
	case "reframe":
		return fmt.Errorf("verb %q not implemented", verb)
	case "delete-frame":
		return fmt.Errorf("verb %q not implemented", verb)
	case "has-frame":
		return fmt.Errorf("verb %q not implemented", verb)
	case "attachments":
		return fmt.Errorf("verb %q not implemented", verb)
	case "attach":
		return fmt.Errorf("verb %q not implemented", verb)
	case "retrieve":
		return fmt.Errorf("verb %q not implemented", verb)
	case "prune":
		return fmt.Errorf("verb %q not implemented", verb)
	case "sample":
		return fmt.Errorf("verb %q not implemented", verb)
	case "clone":
		return doClone(in, out, eout, args)
	case "clone-sample":
		return fmt.Errorf("verb %q not implemented", verb)
	case "check":
		return fmt.Errorf("verb %q not implemented", verb)
	case "repair":
		return fmt.Errorf("verb %q not implemented", verb)
	case "codemeta":
		return fmt.Errorf("verb %q not implemented", verb)
	case "versioning":
		return doVersioning(in, out, eout, args)
	case "versions":
		return doVersions(in, out, eout, args)
	case "read-version":
		return doReadVersion(in, out, eout, args)
	default:
		return fmt.Errorf("verb %q not supported", verb)
	}
}
