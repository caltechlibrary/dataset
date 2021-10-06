package dataset

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

func DisplayLicense(out io.Writer, appName string, license string) {
	fmt.Fprintf(out, strings.ReplaceAll(strings.ReplaceAll(license, "{app_name}", appName), "{version}", Version))
}

func DisplayVersion(out io.Writer, appName string) {
	fmt.Fprintf(out, "%s %s\n", appName, Version)
}

func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet, description string, examples string, license string) {
	// Convert {app_name} and {version} in description
	if description != "" {
		fmt.Fprintf(out, strings.ReplaceAll(description, "{app_name}", appName))
	}
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	if examples != "" {
		fmt.Fprintf(out, strings.ReplaceAll(examples, "{app_name}", appName))
	}
	if license != "" {
		DisplayLicense(out, appName, license)
	}
}
