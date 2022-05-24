package dataset

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// textProcessor takes the a topic document and replaces all the keys
// (e.g. "{app_name}") the their value (e.g. "dataset") in the topic
// test.
func textProcessor(varMap map[string]string, topic string) string {
	src := topic[:]
	for key, val := range varMap {
		if strings.Contains(src, key) {
			src = strings.ReplaceAll(src, key, val)
		}
	}
	return src
}

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string, license string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, textProcessor(m, license))
}

// DisplayVersion returns the of the dataset application.
func DisplayVersion(out io.Writer, appName string) {
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	fmt.Fprintf(out, textProcessor(m, "{app_name} {version}\n"))
}

// DisplayUsage displays a usage message.
func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet, description string, examples string, license string) {
	// Replacable text vars
	m := map[string]string{
		"{app_name}": appName,
		"{version}":  Version,
	}
	// Convert {app_name} and {version} in description
	if description != "" {
		fmt.Fprintf(out, textProcessor(m, description))
	}
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	if examples != "" {
		fmt.Fprintf(out, textProcessor(m, examples))
	}
	if license != "" {
		DisplayLicense(out, appName, textProcessor(m, license))
	}
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
	default:
		return fmt.Errorf("verb %q not supported", verb)
	}
}
