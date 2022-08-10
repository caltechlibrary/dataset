//
// cli is a sub module of dataset.
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
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	// dataset submodule
	ds "github.com/caltechlibrary/dataset/v2"
	"github.com/caltechlibrary/dataset/v2/texts"
)

var (
	showHelp bool
	appName  = path.Base(os.Args[0])

	helpDocs = map[string]string{
		"usage":         cliDescription,
		"examples":      cliExamples,
		"init":          cliInit,
		"create":        cliCreate,
		"read":          cliRead,
		"update":        cliUpdate,
		"delete":        cliDelete,
		"keys":          cliKeys,
		"haskey":        cliHasKey,
		"has-key":       cliHasKey,
		"updated-keys":  cliUpdatedKeys,
		"count":         cliCount,
		"versioning":    cliVersioning,
		"versions":      cliVersioning,
		"read-version":  cliVersioning,
		"sample":        cliSample,
		"clone":         cliClone,
		"clone-sample":  cliCloneSample,
		"frames":        cliFrames,
		"frame":         cliFrame,
		"frame-def":     cliFrameDef,
		"frame-keys":    cliFrameKeys,
		"frame-objects": cliFrameObjects,
		"reframe":       cliReframe,
		"refresh":       cliRefresh,
		"hasframe":      cliHasFrame,
		"has-frame":     cliHasFrame,
		"delete-frame":  cliDeleteFrame,
		"attachments":   cliAttachments,
		"attach":        cliAttach,
		"retrieve":      cliRetrieve,
		"prune":         cliPrune,
		"check":         cliCheck,
		"repair":        cliRepair,
		"migrate":       cliMigrate,
		"codemeta":      cliCodemeta,
		"license":       ds.License,
	}

	verbs = map[string]func(io.Reader, io.Writer, io.Writer, []string) error{
		"help":          DisplayHelp,
		"init":          doInit,
		"create":        doCreate,
		"read":          doRead,
		"update":        doUpdate,
		"delete":        doDelete,
		"keys":          doKeys,
		"updated-keys":  doUpdatedKeys,
		"haskey":        doHasKey,
		"has-key":       doHasKey,
		"count":         doCount,
		"frames":        doFrameNames,
		"frame":         doFrameCreate,
		"frame-def":     doFrameDef,
		"frame-keys":    doFrameKeys,
		"frame-objects": doFrameObjects,
		"refresh":       doRefresh,
		"reframe":       doReframe,
		"delete-frame":  doDeleteFrame,
		"hasframe":      doHasFrame,
		"has-frame":     doHasFrame,
		"attachments":   doAttachments,
		"attach":        doAttach,
		"retrieve":      doRetrieve,
		"prune":         doPrune,
		"sample":        doSample,
		"clone":         doClone,
		"clone-sample":  doCloneSample,
		"check":         doCheck,
		"repair":        doRepair,
		"migrate":       doMigrate,
		"codemeta":      doCodemeta,
		"versioning":    doVersioning,
		"versions":      doVersions,
		"read-version":  doReadVersion,
	}
)

// DisplayHelp writes out help on a supported topic
func DisplayHelp(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var topic string
	if len(args) > 0 {
		topic = args[0]
	}
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  ds.Version,
	}
	if text, ok := helpDocs[topic]; ok {
		fmt.Fprintf(out, texts.StringProcessor(m, text))
	} else {
		fmt.Fprintf(eout, "Unable to find help on %q\n", topic)
	}
	return nil
}

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  ds.Version,
	}
	fmt.Fprintf(out, texts.StringProcessor(m, ds.License))
}

// DisplayVersion returns the of the dataset application.
func DisplayVersion(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  ds.Version,
	}
	fmt.Fprintf(out, texts.StringProcessor(m, "{app_name} {version}\n"))
}

// DisplayUsage displays a usage message.
func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet) {
	// Replacable text vars
	description, examples := cliDescription, cliExamples
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  ds.Version,
	}
	// Convert {app_name} and {version} in description
	fmt.Fprintf(out, texts.StringProcessor(m, description))
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	fmt.Fprintf(out, texts.StringProcessor(m, examples))
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
		DisplayHelp(in, out, eout, []string{"init"})
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
	c, err := ds.Init(cName, dsnURI)
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
		DisplayHelp(in, out, eout, []string{"versioning"})
	}
	switch {
	case len(args) == 2:
		cName, versioning = args[0], strings.ToLower(strings.TrimSpace(args[1]))
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME VERSIONING_SETTING, got %q", strings.Join(append([]string{appName, "versioning"}, args...), " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	return c.SetVersioning(versioning)
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
		DisplayHelp(in, out, eout, []string{"versioning"})
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(append([]string{appName, "versions"}, args...), " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	versions, err := c.Versions(key)
	if err != nil {
		return fmt.Errorf("version errors for %q, %s", key, err)
	}
	return texts.WriteSource(output, out, []byte(strings.Join(versions, "\n")))
}

func doReadVersion(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		cName   string
		key     string
		version string
		src     []byte
		output  string
	)
	flagSet := flag.NewFlagSet("read-version", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for create")
	flagSet.BoolVar(&showHelp, "help", false, "help for create")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"versioning"})
	}
	switch {
	case len(args) == 3:
		cName, key, version = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected [OPTIONS] COLLECTION_NAME KEY VERSION, got %q", strings.Join(append([]string{appName, "read-version"}, args...), " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	switch c.StoreType {
	case ds.PTSTORE:
		src, err = c.PTStore.ReadVersion(key, version)
	case ds.SQLSTORE:
		src, err = c.SQLStore.ReadVersion(key, version)
	default:
		return fmt.Errorf("%q storage not supported", c.StoreType)
	}
	if err == nil {
		return texts.WriteSource(output, out, src)
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
		DisplayHelp(in, out, eout, []string{"create"})
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		// Read the JSON object from a file or standard input
		src, err = texts.ReadSource(input, in)
		if err != nil {
			return fmt.Errorf("could not read JSON file, %s", err)
		}
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY [JSON_SRC], got %q", strings.Join(append([]string{appName, "create"}, args...), " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	obj := map[string]interface{}{}
	if err := ds.DecodeJSON(src, &obj); err != nil {
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
	)
	flagSet := flag.NewFlagSet("read", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"create"})
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	src := []byte{}
	switch c.StoreType {
	case ds.PTSTORE:
		src, err = c.PTStore.Read(key)
	case ds.SQLSTORE:
		src, err = c.SQLStore.Read(key)
	default:
		return fmt.Errorf("%q storage not supportted", c.StoreType)
	}
	if err == nil {
		return texts.WriteSource(output, out, src)
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
		DisplayHelp(in, out, eout, []string{"update"})
	}
	switch {
	case len(args) == 3:
		cName, key, src = args[0], args[1], []byte(args[2])
	case len(args) == 2:
		cName, key = args[0], args[1]
		// Read JSON source
		src, err = texts.ReadSource(input, in)
		if err != nil {
			return fmt.Errorf("could not read JSON file, %s", err)
		}
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY [JSON_SRC], got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	obj := map[string]interface{}{}
	if err := ds.DecodeJSON(src, &obj); err != nil {
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
		DisplayHelp(in, out, eout, []string{"delete"})
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	switch c.StoreType {
	case ds.PTSTORE:
		err = c.PTStore.Delete(key)
	case ds.SQLSTORE:
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
		DisplayHelp(in, out, eout, []string{"keys"})
	}
	switch {
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
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
	return texts.WriteSource(output, out, src)
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
		DisplayHelp(in, out, eout, []string{"keys"})
	}
	switch {
	case len(args) == 3:
		cName, start, end = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME START END, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()
	keys, err = c.UpdatedKeys(start, end)
	if err != nil {
		return err
	}
	src := []byte(strings.Join(keys, "\n"))
	return texts.WriteSource(output, out, src)
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
		DisplayHelp(in, out, eout, []string{"has-key"})
	}
	switch {
	case len(args) == 2:
		cName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
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
		DisplayHelp(in, out, eout, []string{"count"})
	}
	switch {
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	c, err := ds.Open(cName)
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
		keys      []string
		err       error
	)
	flagSet := flag.NewFlagSet("clone", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&keysName, "i", "-", "filename to read keys from")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"clone"})
	}
	switch {
	case len(args) == 2:
		srcName, dstName, dstDsnURI = args[0], args[1], ""
	case len(args) == 3:
		srcName, dstName, dstDsnURI = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTIONS] SRC_COLLECTION_NAME DEST_COLLECTION_NAME [DEST_DSN_URI], got %q", strings.Join(args, " "))
	}
	keys, err = texts.ReadKeys(keysName, in)
	if err != nil {
		return err
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	if err := source.Clone(dstName, dstDsnURI, keys, verbose); err != nil {
		return fmt.Errorf("clone failed %s", err)
	}
	return nil
}

// doFrameNames
func doFrameNames(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		err     error
		output  string
	)
	flagSet := flag.NewFlagSet("frames", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"frames"})
		return nil
	}
	switch {
	case len(args) == 1:
		srcName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	frames := source.FrameNames()
	src := []byte(strings.Join(frames, "\n"))
	return texts.WriteSource(output, out, src)
}

// doFrameCreate
func doFrameCreate(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		input     string
		keys      []string
		dotPaths  []string
		labels    []string
		verbose   bool
		err       error
	)
	flagSet := flag.NewFlagSet("frame", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&input, "i", "-", "filename to read keys from")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"frame"})
	}
	switch {
	case len(args) >= 3:
		srcName, frameName, args = args[0], args[1], args[2:]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME DOT_PATH [DOT_PATH...] got %q", strings.Join(args, " "))
	}
	keys, err = texts.ReadKeys(input, in)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("missing keys for frame %q", frameName)
	}
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			dotPaths = append(dotPaths, parts[0])
			labels = append(labels, parts[1])
		} else {
			dotPaths = append(dotPaths, arg)
			labels = append(labels, arg)
		}
	}

	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	frame, err := source.FrameCreate(frameName, keys, dotPaths, labels, verbose)
	if err != nil {
		return err
	}
	if frame == nil {
		return fmt.Errorf("failed to create frame %q for %q", frameName, source.Name)
	}
	return nil
}

// doFrameDef
func doFrameDef(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
		output    string
	)
	flagSet := flag.NewFlagSet("frame-def", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"frame-def"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	m, err := source.FrameDef(frameName)
	if err != nil {
		return err
	}
	src, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}
	return texts.WriteSource(output, out, src)
}

// doFrameObjects
func doFrameObjects(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		output    string
		err       error
	)
	flagSet := flag.NewFlagSet("frame-objects", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write output to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"frame-objects"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	objects, err := source.FrameObjects(frameName)
	if err != nil {
		return err
	}
	src, err := json.MarshalIndent(objects, "", "    ")
	if err != nil {
		return err
	}
	return texts.WriteSource(output, out, src)
}

// doFrameKeys
func doFrameKeys(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
		output    string
	)
	flagSet := flag.NewFlagSet("frame-keys", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"frame-keys"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	keys := source.FrameKeys(frameName)
	return texts.WriteKeys(output, out, keys)
}

// doRefresh
func doRefresh(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
		verbose   bool
	)
	flagSet := flag.NewFlagSet("refresh", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"refresh"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	return source.FrameRefresh(frameName, verbose)
}

// doReframe
func doReframe(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
		input     string
		verbose   bool
		keys      []string
	)
	flagSet := flag.NewFlagSet("reframe", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&input, "i", "-", "read keys from a file")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"reframe"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	keys, err = texts.ReadKeys(input, in)
	if err != nil {
		return err
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	return source.FrameReframe(frameName, keys, verbose)
}

// doDeleteFrame
func doDeleteFrame(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
	)
	flagSet := flag.NewFlagSet("delete-frame", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"delete-frame"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	return source.FrameDelete(frameName)
}

// doHasFrame
func doHasFrame(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName   string
		frameName string
		err       error
	)
	flagSet := flag.NewFlagSet("has-frame", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"has-frame"})
	}
	switch {
	case len(args) == 2:
		srcName, frameName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME FRAME_NAME, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	if source.HasFrame(frameName) {
		fmt.Fprintf(out, "true\n")
	} else {
		fmt.Fprintf(out, "false\n")
		return fmt.Errorf(" ")
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
		DisplayHelp(in, out, eout, []string{"attachments"})
	}
	switch {
	case len(args) == 2:
		srcName, key = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return fmt.Errorf("failed to open %q, %s", srcName, err)
	}
	defer source.Close()
	attachments, err := source.Attachments(key)
	if err != nil {
		return fmt.Errorf("failed to get attachments %q, %q, %s", srcName, key, err)
	}
	return texts.WriteKeys(output, out, attachments)
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
		DisplayHelp(in, out, eout, []string{"attach"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filenames = args[0], args[1], args[2:]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY FILENAME[FILENAME ...], got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
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
	flagSet.StringVar(&output, "o", "-", "save to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"retrieve"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filename = args[0], args[1], args[2]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY FILENAMEgot, %q", strings.Join(args, " "))
	}
	c, err := ds.Open(srcName)
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
		DisplayHelp(in, out, eout, []string{"prune"})
	}
	switch {
	case len(args) == 3:
		srcName, key, filenames = args[0], args[1], args[2:]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY FILENAME[FILENAME ...], got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
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

// doSample
func doSample(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName string
		size    string
		err     error
		keys    []string
		output  string
	)
	flagSet := flag.NewFlagSet("clone-sample", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&output, "o", "-", "write to file")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"clone-sample"})
	}
	switch {
	case len(args) == 2:
		srcName, size = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME SIZE_OF_SAMPLE_KEYS, got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return err
	}
	defer source.Close()
	i, err := strconv.Atoi(size)
	if err != nil {
		return fmt.Errorf("size %q doesn't make sense, %s", size, err)
	}
	keys, err = source.Sample(i)
	if err != nil {
		return fmt.Errorf("sampling keys failed, %s", err)
	}
	return texts.WriteKeys(output, out, keys)
}

// doCloneSample
func doCloneSample(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var (
		srcName        string
		trainingName   string
		trainingDsnURI string
		testName       string
		testDsnURI     string
		keysName       string
		sampleSize     int
		verbose        bool
		keys           []string
		err            error
	)
	flagSet := flag.NewFlagSet("clone-sample", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "display help")
	flagSet.BoolVar(&showHelp, "help", false, "display help")
	flagSet.StringVar(&keysName, "i", "-", "filename to read keys from")
	flagSet.IntVar(&sampleSize, "size", 0, "sample size for training set")
	flagSet.BoolVar(&verbose, "verbose", false, "verbose output")
	flagSet.Parse(args)
	args = flagSet.Args()
	if showHelp {
		DisplayHelp(in, out, eout, []string{"clone-sample"})
	}
	switch {
	case len(args) == 5:
		srcName, trainingName, trainingDsnURI, testName, testDsnURI = args[0], args[1], args[2], args[3], args[4]
	case len(args) == 4:
		srcName, trainingName, trainingDsnURI, testName, testDsnURI = args[0], args[1], args[2], args[3], ""
	case len(args) == 3:
		srcName, trainingName, trainingDsnURI, testName, testDsnURI = args[0], args[1], args[2], "", ""
	default:
		return fmt.Errorf("Expected: [OPTIONS] SRC_COLLECTION_NAME TRAINING_COLLECTION TRAINING_DSN_URI [DEST_COLLECTION_NAME [TEST_DSN_URI]], got %q", strings.Join(args, " "))
	}
	source, err := ds.Open(srcName)
	if err != nil {
		return err
	}
	keys, err = texts.ReadKeys(keysName, in)
	if err != nil {
		return err
	}
	if err := source.CloneSample(trainingName, trainingDsnURI, testName, testDsnURI, keys, sampleSize, verbose); err != nil {
		return fmt.Errorf("clone-sample failed %s", err)
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
		DisplayHelp(in, out, eout, []string{"check"})
	}
	switch {
	case len(args) == 1:
		srcName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return ds.Analyzer(srcName, verbose)
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
		DisplayHelp(in, out, eout, []string{"repair"})
	}
	switch {
	case len(args) == 1:
		srcName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return ds.Repair(srcName, verbose)
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
		DisplayHelp(in, out, eout, []string{"migrate"})
	}
	switch {
	case len(args) == 2:
		srcName, dstName = args[0], args[1]
	default:
		return fmt.Errorf("Expected: [OPTIONS] OLD_COLLECTION_NAME NEW_COLLECTION_NAME, got %q", strings.Join(args, " "))
	}
	return ds.Migrate(srcName, dstName, verbose)
}

// doCodemeta
func doCodemeta(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	return fmt.Errorf("doCodemeta() not implemented")
}

/// RunCLI implemented the functionlity used by the cli.
func RunCLI(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	var err error

	if len(args) == 0 {
		DisplayHelp(in, out, eout, []string{"usage"})
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
