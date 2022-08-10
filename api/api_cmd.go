package api

import (
	"flag"
	"fmt"
	"io"

	// Caltech Library packages
	ds "github.com/caltechlibrary/dataset/v2"
	"github.com/caltechlibrary/dataset/v2/texts"
)

// DisplayLicense returns the license associated with dataset application.
func DisplayLicense(out io.Writer, appName string) {
	fmt.Fprintf(out, ds.License)
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
	description := `
USAGE
=====

    {app_name} SETTINGS_JSON_FILE

DESCRIPTION
--------

Runs a web service for one or more dataset collections. Requires
the collections to exist (e.g. created previously with the dataset
cli). It requires a settings JSON file that decribes the web service
configuration and permissions per collection that are available via
the web service.

EXAMPLE
-------

Starting up the web service

` + "```" + `
   {app_name} settings.json
` + "```" + `

`
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
