// cli is part of dataset
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
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/models"

	// 3rd Party packages
	"github.com/pkg/fileutils"
	"gopkg.in/yaml.v3"
)

var (
	showHelp bool
	appName  = path.Base(os.Args[0])

	helpDocs = map[string]string{
		"usage":          cliDescription,
		"examples":       cliExamples,
		"init":           cliInit,
		"model":          cliModel,
		"create":         cliCreate,
		"read":           cliRead,
		"update":         cliUpdate,
		"delete":         cliDelete,
		"query":          cliQuery,
		"keys":           cliKeys,
		"haskey":         cliHasKey,
		"has-key":        cliHasKey,
		"updated-keys":   cliUpdatedKeys,
		"count":          cliCount,
		"set-versioning": cliVersioning,
		"get-versioning": cliVersioning,
		"versions":       cliVersioning,
		"attachments":    cliAttachments,
		"attach":         cliAttach,
		"detach":         cliRetrieve,
		"retrieve":       cliRetrieve,
		"prune":          cliPrune,
		"check":          cliCheck,
		"repair":         cliRepair,
		"migrate":        cliMigrate,
		"codemeta":       cliCodemeta,
		"license":        License,
		"load":           cliLoad,
		"dump":           cliDump,
	}

	verbs = map[string]func(io.Reader, io.Writer, io.Writer, []string) error{
		"help":           CliDisplayHelp,
		"init":           doInit,
		"model":          doModel,
		"create":         doCreate,
		"read":           doRead,
		"update":         doUpdate,
		"delete":         doDelete,
		"query":          doQuery,
		"keys":           doKeys,
		"updated-keys":   doUpdatedKeys,
		"haskey":         doHasKey,
		"has-key":        doHasKey,
		"count":          doCount,
		"attachments":    doAttachments,
		"attach":         doAttach,
		"detach":         doRetrieve,
		"retrieve":       doRetrieve,
		"prune":          doPrune,
		"check":          doCheck,
		"repair":         doRepair,
		"migrate":        doMigrate,
		"codemeta":       doCodemeta,
		"get-versioning": doGetVersioning,
		"set-versioning": doSetVersioning,
		"versions":       doVersions,
		"load":           doLoad,
		"dump":           doDump,
	}
)

func prettyPrintJSON(src []byte) ([]byte, error) {
	// Force output to be pretty printed. I can't rely on a
	// standard way to implement this in SQL.
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	if err := json.Indent(w, src, "", "    "); err != nil {
		return nil, err
	}
	src = w.Bytes()
	return src, nil
}

// CliDisplayHelp writes out help on a supported topic
func CliDisplayHelp(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var topic string
	if len(args) > 0 {
		topic = args[0]
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

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprint(out, StringProcessor(m, License))
}

// DisplayVersion returns the of the dataset application.
func DisplayVersion(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprint(out, StringProcessor(m, "{app_name} {version}\n"))
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
	DisplayLicense(out, appName)
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
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME [DSN_URI], got %q", strings.Join(args, " "))
	}
	c, err := Init(cName, dsnURI)
	if err == nil {
		defer c.Close()
	}
	return err
}

func doModel(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected: COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	cName := args[0]
	modelId := "model" //strings.TrimSuffix(path.Base(cName), ".ds")
	fName := path.Join(cName, "model.yaml")
	model, err := models.NewModel(modelId)
	if err != nil {
		return err
	}
	backupFile := false
	if _, err := os.Stat(fName); err == nil {
		backupFile = true
		src, err := os.ReadFile(fName)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(src, &model); err != nil {
			return err
		}
	}
	// Decide if I'm going to create or open an existing YAML file.
	fout, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()

	models.SetDefaultTypes(model)
	model.Register("yaml", models.ModelToYAML)
	if err := models.ModelInteractively(model); err != nil {
		return err
	}
	if model.HasChanges() {
		prompt := models.NewPrompt(in, out, eout)
		fmt.Fprintf(out, "Save %s (Y/n)?", fName)
		answer := prompt.GetAnswer("y", true)
		if answer == "y" {
			// backup file if needed
			if backupFile {
				if err := fileutils.CopyFile(fName+".bak", fName); err != nil {
					return err
				}
			}
			if err := model.Render(fout, "yaml"); err != nil {
				return err
			}
		}
	}
	return nil
}

func doSetVersioning(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName      string
		versioning string
	)
	flagSet := flag.NewFlagSet("set_versioning", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"set_versioning"})
	}
	switch {
	case len(args) == 2:
		cName, versioning = args[0], strings.ToLower(strings.TrimSpace(args[1]))
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME VERSIONING_SETTING, got %q", strings.Join(append([]string{appName, "set_versioning"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetVersioning(versioning)
}

func doGetVersioning(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
	)
	flagSet := flag.NewFlagSet("get_versioning", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"get_versioning"})
	}
	switch {
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME, got %q", strings.Join(append([]string{appName, "get_versioning"}, args...), " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	fmt.Fprintf(out, "%s\n", c.Versioning)
	return nil
}

func doVersions(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		key    string
		output string
	)
	flagSet := flag.NewFlagSet("versions", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&output, "o", "-", "write output to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"versioning"})
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
	return WriteSource(output, out, []byte(strings.Join(versions, "\n")))
}

func doReadVersion(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName   string
		key     string
		version string
		src     []byte
		output  string
		pretty  bool
	)
	flagSet := flag.NewFlagSet("read-version", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.BoolVar(&pretty, "pretty", false, "pretty print output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"versioning"})
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
		// NOTE: SQL databases will store JSON in an un-pretty way.
		// I want to pretty print the JSON I output.
		src, err = c.SQLStore.ReadVersion(key, version)
	default:
		return fmt.Errorf("%q storage not supported", c.StoreType)
	}
	if pretty {
		src, err = prettyPrintJSON(src)
	}
	if err == nil {
		return WriteSource(output, out, src)
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
		overwrite bool
	)
	flagSet := flag.NewFlagSet("create", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&input, "i", "-", "read JSON from file, use '-' for stdin")
	flagSet.StringVar(&input, "input", "-", "read JSON from file, use '-' for stdin")
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite object if it previously exists")
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
	if overwrite && c.HasKey(key) {
		return c.Update(key, obj)
	}
	if err := c.Create(key, obj); err != nil {
		return err
	}
	return nil
}

func doRead(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		key    string
		output string
		pretty bool
	)
	flagSet := flag.NewFlagSet("read", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.BoolVar(&pretty, "pretty", true, "pretty print output")
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
	case PTSTORE:
		src, err = c.PTStore.Read(key)
	case SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	if pretty {
		src, err = prettyPrintJSON(src)
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
		cName      string
		output     string
		sampleSize int
		keys       []string
		err        error
	)
	flagSet := flag.NewFlagSet("keys", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.IntVar(&sampleSize, "sample", 0, "generate sample N of keys, where N is greater than zero")
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
	if sampleSize > 0 {
		keys, err = c.Sample(sampleSize)
	} else {
		keys, err = c.Keys()
	}
	if err != nil {
		return err
	}
	src := []byte(strings.Join(keys, "\n"))
	if (len(keys) > 0) {
		src = append(src, []byte("\n")...)
	}
	return WriteSource(output, out, src)
}

func doUpdatedKeys(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName  string
		start  string
		end    string
		output string
		keys   []string
		err    error
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
	case len(args) == 3:
		cName, start, end = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME START END, got %q", strings.Join(args, " "))
	}
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	keys, err = c.UpdatedKeys(start, end)
	if err != nil {
		return err
	}
	src := []byte(strings.Join(keys, "\n"))
	return WriteSource(output, out, src)
}

func doHasKey(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName string
		key   string
	)
	flagSet := flag.NewFlagSet("has-key", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"has-key"})
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
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"count"})
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
		overwrite   bool
		maxCapacity = 0
	)
	flagSet := flag.NewFlagSet("load", flag.ContinueOnError)
	flagSet.BoolVar(&overwrite, "o", false, "overwrite existing objects on load")
	flagSet.BoolVar(&overwrite, "overwrite", false, "overwrite existing objects on load")
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
	sqlFName, showHelp := "", false
	flagSet := flag.NewFlagSet("query", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&sqlFName, "sql", sqlFName, "read SQL statement from a file")
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
	if sqlFName != "" {
		if sqlFName != "-" {
			var err error
			in, err = os.Open(sqlFName)
			if err != nil {
				return err
			}
			//defer in.Close()
		}
		src, err := io.ReadAll(in)
		if err != nil {
			return err
		}
		stmt = fmt.Sprintf("%s", src)
	}
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
	if cName == "" {
		return fmt.Errorf("missing C_NAME")
	}
	if stmt == "" {
		return fmt.Errorf("missing SQL_STATEMENT")
	}
	if err := app.Run(os.Stdin, os.Stdout, os.Stderr, cName, stmt, params); err != nil {
		return err
	}
	return nil
}

// doAttachments
func doAttachments(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		key     string
		err     error
		output  string
	)
	flagSet := flag.NewFlagSet("attachments", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"attachments"})
	}
	switch {
	case len(args) == 2:
		srcName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	source, err := Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	attachments, err := source.Attachments(key)
	if err != nil {
		return fmt.Errorf("failed to get attachments %q, %q, %s", srcName, key, err)
	}
	return WriteKeys(output, out, attachments)
}

// doAttach
func doAttach(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		filenames []string
		key       string
		err       error
		output    string
	)
	flagSet := flag.NewFlagSet("attach", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"attach"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filenames = args[0], args[1], args[2:]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY FILENAME[FILENAME ...], got %q", strings.Join(args, " "))
	}
	source, err := Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	for _, filename := range filenames {
		err := source.AttachFile(key, filename)
		if err != nil {
			fmt.Fprintf(eout, "failed to attach %q to %q, %s", filename, key, err)
		}
	}
	return nil
}

// doRetrieve
func doRetrieve(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName  string
		filename string
		key      string
		err      error
		output   string
	)
	flagSet := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "", "save to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"retrieve"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filename = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTION] COLLECTION_NAME KEY FILENAME got, %q", strings.Join(args, " "))
	}
	c, err := Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer c.Close()
	// Handle writing to stdout
	if output == "-" {
		fOut := out
		err = c.RetrieveStream(key, filename, fOut)
		if err != nil {
			return fmt.Errorf("failed to retrieve %q from %q, %s", filename, key, err)
		}
		return nil
	}
	// Handle writing to designated file
	if output == "" {
		output = filename
	}
	fOut, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("could not create %q, %s", output, err)
	}
	defer fOut.Close()
	err = c.RetrieveStream(key, filename, fOut)
	if err != nil {
		return fmt.Errorf("failed to retrieve %q from %q, %s", filename, key, err)
	}
	return nil
}

// doPrune
func doPrune(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		filenames []string
		key       string
		err       error
		output    string
	)
	flagSet := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"prune"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filenames = args[0], args[1], args[2:]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY FILENAME[FILENAME ...], got %q", strings.Join(args, " "))
	}
	source, err := Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	for _, filename := range filenames {
		err := source.Prune(key, filename)
		if err != nil {
		}
		if err != nil {
			fmt.Fprintf(eout, "failed to prune %q from %q, %s", filename, key, err)
		}
	}
	return nil
}

// doCheck
func doCheck(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		verbose bool
	)
	flagSet := flag.NewFlagSet("check", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&verbose, "verbose", false, "set verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"check"})
	}
	switch {
	case len(args) == 1:
		srcName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return Analyzer(srcName, verbose)
}

// doRepair
func doRepair(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		verbose bool
	)
	flagSet := flag.NewFlagSet("check", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"repair"})
	}
	switch {
	case len(args) == 1:
		srcName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return Repair(srcName, verbose)
}

// doMigrate
func doMigrate(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		dstName string
		verbose bool
	)
	flagSet := flag.NewFlagSet("migrate", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		CliDisplayHelp(in, out, eout, []string{"migrate"})
	}
	switch {
	case len(args) == 2:
		srcName, dstName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] OLD_COLLECTION_NAME NEW_COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return Migrate(srcName, dstName, verbose)
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
