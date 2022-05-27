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
	fmt.Printf("DEBUG cName %q, dsnURI %q\n", cName, dsnURI)
	c, err := Init(cName, dsnURI)
	if err == nil {
		defer c.Close()
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
		if input == "" || input == "-" {
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
		return fmt.Errorf("verb %q not implemented", verb)
	case "clone-sample":
		return fmt.Errorf("verb %q not implemented", verb)
	case "check":
		return fmt.Errorf("verb %q not implemented", verb)
	case "repair":
		return fmt.Errorf("verb %q not implemented", verb)
	case "codemeta":
		return fmt.Errorf("verb %q not implemented", verb)
	default:
		return fmt.Errorf("verb %q not supported", verb)
	}
}
