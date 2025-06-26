// cli is part of dataset
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2025, Caltech
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
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	showHelp bool
	appName  = filepath.Base(os.Args[0])

	helpDocs = map[string]string{
		"description":          cliDescription,
		"examples":       cliExamples,
		"init":           cliInit,
		"write":         cliWrite,
		"read":           cliRead,
		"delete":         cliDelete,
		"query":          cliQuery,
		"keys":           cliKeys,
		"haskey":         cliHasKey,
		"codemeta":       cliCodemeta,
		"load":           cliLoad,
		"dump":           cliDump,
		"help":           cliHelp,
	}

	verbs = map[string]func(io.Reader, io.Writer, io.Writer, []string) error{
		"help":           CliDisplayHelp,
		"init":           doInit,
		"create":         doCreate,
		"read":           doRead,
		"update":         doUpdate,
		"delete":         doDelete,
		"query":          doQuery,
		"keys":           doKeys,
		"haskey":         doHasKey,
		"codemeta":       doCodemeta,
		"load":           doLoad,
		"dump":           doDump,
	}
)

// CliDisplayHelp writes out help on a supported topic
func CliDisplayHelp(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var topic string
	if len(args) > 0 {
		topic = strings.ToLower(args[0])
	} else {
		topic = "description"
	}
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	if text, ok := helpDocs[topic]; ok {
		fmt.Fprint(out, StringProcessor(m, text))
	} else {
		fmt.Fprintf(eout, "Unable to find help on %q\n", topic)
	}
	return nil
}

// CliDisplayUsage displays a usage message.
func CliDisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet) {
	// Replacable text vars
	description, examples := cliDescription, cliExamples
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	// Convert {app_name} and {version} in description
	fmt.Fprint(out, StringProcessor(m, description))
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	fmt.Fprint(out, StringProcessor(m, examples))
}

func doInit(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		dsnURI string
	)
	flagSet := flag.NewFlagSet("init", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for init")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"init"})
		return nil
	}
	switch {
	case len(args) == 2:
		cName, dsnURI = args[0], args[1]
	case len(args) == 1:
		cName, dsnURI = args[0], ""
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	c, err := Init(cName, dsnURI)
	if err == nil {
		defer c.Close()
	}
	return err
}


func doCreate(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName     string
		key       string
		src       []byte
		input     string
		err       error
	)
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&input, "i", "-", "read JSON from file, use '-' for stdin")
	flagSet.StringVar(&input, "input", "-", "read JSON from file, use '-' for stdin")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"create"})
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		// Read the JSON object from a file or standard input
		src, err = ReadSource(input, in)
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
	if err := JSONUnmarshal(src, &obj); err != nil {
		return err
	}
	if err := c.Write(key, obj); err != nil {
		return err
	}
	return nil
}

func doRead(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		key    string
		output string
		jsonl bool
	)
	flagSet := flag.NewFlagSet("read", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&jsonl, "jsonl", false, "output as JSON line")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"create"})
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
	src := []byte{}
	switch c.StoreType {
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	if jsonl {
		src = FmtJSONL(src)
	}
	if err == nil {
		return WriteSource(output, out, src)
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
		CliDisplayHelp(in, out, eout, []string{"update"})
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		// Read JSON source
		src, err = ReadSource(input, in)
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
	if err := JSONUnmarshal(src, &obj); err != nil {
		return err
	}
	if err := c.Write(key, obj); err != nil {
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
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"delete"})
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
	case SQLSTORE:
		err = c.SQLStore.Delete(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	return err
}

func doKeys(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName      string
		output     string
		keys       []string
		err        error
	)
	flagSet := flag.NewFlagSet("keys", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"keys"})
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
	keys, err = c.Keys()
	if err != nil {
		return err
	}
	src := []byte(strings.Join(keys, "\n"))
	if len(keys) > 0 {
		// NOTE: when we have one or more keys make sure we have a trailing newline.
		src = append(src, []byte("\n")...)
	}
	return WriteSource(output, out, src)
}


func doHasKey(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("haskey", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"haskey"})
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

func doDump(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
	)
	flagSet := flag.NewFlagSet("dump", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"dump"})
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
	return c.Dump(os.Stdout)
}

func doLoad(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName       string
		maxCapacity = 0
		overwrite bool
	)
	flagSet := flag.NewFlagSet("load", flag.ContinueOnError)
	flagSet.BoolVar(&overwrite, "o", false, "overwrite existing objects on load")
	flagSet.IntVar(&maxCapacity, "m", maxCapacity, "set a maximum size for single object in megabytes")
	flagSet.IntVar(&maxCapacity, "max-capacity", maxCapacity, "set a maximum size for single object in megabytes")
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"dump"})
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
	return c.Load(os.Stdin, overwrite, maxCapacity)
}

func doQuery(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	sqlFName, showHelp, jsonl := "", false, false
	flagSet := flag.NewFlagSet("query", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&sqlFName, "sql", sqlFName, "read SQL statement from a file")
	flagSet.BoolVar(&jsonl, "jsonl", false, "output using JSON lines")
	flagSet.Parse(args)
	args = flagSet.Args()

	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"query"})
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("missing C_NAME and SQL_STATEMENT")
	}
	// Create a DSQuery object and evaluate the command line options
	app := new(DSQuery)
	cName, stmt, params := "", "", []string{}
	for _, arg := range args {
		switch {
		case cName == "":
			cName = arg
		case stmt == "":
			stmt = arg
		default:
			params = append(params, arg)
		}
	}
	if sqlFName != "" {
		if sqlFName != "-" {
			var err error
			in, err = os.Open(sqlFName)
			if err != nil {
				return err
			}
			//defer in.Close()
		}
	}
	if cName == "" {
		return fmt.Errorf("missing C_NAME")
	}
	if stmt == "" {
		src, err := io.ReadAll(in)
		if err != nil {
			return err
		}
		stmt = fmt.Sprintf("%s", src)
	}
	if err := app.RunQuery(os.Stdin, os.Stdout, os.Stderr, cName, stmt, jsonl, params); err != nil {
		return err
	}
	return nil
}

// doCodemeta
func doCodemeta(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var cPath string
	switch {
	case len(args) == 1:
		cPath = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	src, err := ioutil.ReadFile(path.Join(cPath, "codemeta.json"))
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}

// / RunCLI implemented the functionlity used by the cli.
func RunCLI(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var err error

	if len(args) == 0 {
		CliDisplayHelp(in, out, eout, []string{"usage"})
		return fmt.Errorf(` `)
	}
	verb, args := args[0], args[1:]
	if fn, ok := verbs[verb]; ok {
		err = fn(in, out, eout, args)
	} else {
		return fmt.Errorf("verb %q not supported", verb)
	}
	return err
}
