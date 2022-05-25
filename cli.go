package dataset

import (
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
	case "create":
		fmt.Fprint(out, StringProcessor(m, cliCreate))
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
		cName     string
		dsnURI    string
		storeType string
	)
	flagSet := flag.NewFlagSet("init", flag.ContinueOnError)
	flagSet.BoolVar(&showHelp, "h", false, "help for init")
	flagSet.Parse(args)
	if showHelp {
		DisplayHelp(out, eout, "init")
		return nil
	}
	switch {
	case len(args) == 3:
		cName, dsnURI, storeType = args[0], args[1], args[3]
	case len(args) == 2:
		cName, dsnURI = args[0], args[1]
	case len(args) == 1:
		cName = args[0]
	default:
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME [DSN_URI] [COLLECTION_TYPE], got %s", strings.Join(args, " "))
	}
	c, err := Init(cName, dsnURI, storeType)
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
		return fmt.Errorf("Expected: [OPTIONS] COLLECTION_NAME KEY [JSON_SRC]")
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
		return fmt.Errorf("verb %q not implemented", verb)
	case "update":
		return fmt.Errorf("verb %q not implemented", verb)
	case "delete":
		return fmt.Errorf("verb %q not implemented", verb)
	case "keys":
		return fmt.Errorf("verb %q not implemented", verb)
	case "has-keys":
		return fmt.Errorf("verb %q not implemented", verb)
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
