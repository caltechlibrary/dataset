package dataset

import (
	"flag"
	"fmt"
	"io"
)

// ApiDisplayUsage displays a usage message.
func ApiDisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet) {
	// Replacable text vars
	description := `
USAGE
=====

    {app_name} SETTINGS_JSON_FILE

DESCRIPTION
--------

Configures or runs a web service for one or more dataset collections. Requires
the collections to exist (e.g. created previously with the dataset
cli). It requires a settings JSON file that decribes the web service
configuration and permissions per collection that are available via
the web service.

EXAMPLE
-------

Starting up the web service

` + "```" + `
   {app_name} settings.yaml
` + "```" + `

`
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
