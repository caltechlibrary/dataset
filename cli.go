package dataset

import (
	"flag"
	"fmt"
	"io"
)

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string, license string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, TextProcessor(m, license))
}

// DisplayVersion returns the of the dataset application.
func DisplayVersion(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, TextProcessor(m, "{app_name} {version}\n"))
}

// DisplayUsage displays a usage message.
func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet) {
	// Replacable text vars
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	// Convert {app_name} and {version} in description
	fmt.Fprintf(out, TextProcessor(m, CLIDescription))
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	fmt.Fprintf(out, TextProcessor(m, CLIExamples))
	DisplayLicense(out, appName, TextProcessor(m, License))
}

/// RunCLI implemented the functionlity used by the cli.
func RunCLI(in io.Reader, out io.Writer, eout io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Missing parameters")
	}
	verb, args := args[0], args[1:]
	switch verb {
	case "help":
		return fmt.Errorf("verb %q not implemented", verb)
	case "init":
		return fmt.Errorf("verb %q not implemented", verb)
	case "create":
		return fmt.Errorf("verb %q not implemented", verb)
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
