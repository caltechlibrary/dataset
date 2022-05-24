package main

//
// datasetd provide collection access/management via a simple
// HTTP/HTTPS API
//

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/dataset"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	description = dataset.WEBDescription
	examples    = dataset.WEBExamples
	license     = dataset.License
)

func main() {
	appName := path.Base(os.Args[0])
	flagSet := flag.NewFlagSet(appName, flag.ContinueOnError)
	// Standard Options
	flagSet.BoolVar(&showHelp, "help", false, "display detailed help")
	flagSet.BoolVar(&showLicense, "license", false, "display license")
	flagSet.BoolVar(&showVersion, "version", false, "display version")

	flagSet.Parse(os.Args[1:])
	args := flagSet.Args()

	if showHelp {
		dataset.DisplayUsage(os.Stdout, appName, flagSet, description, examples, license)
		os.Exit(0)
	}

	if showLicense {
		dataset.DisplayLicense(os.Stdout, appName, license)
		os.Exit(0)
	}

	if showVersion {
		dataset.DisplayVersion(os.Stdout, appName)
		os.Exit(0)
	}

	/* Looking for settings.json */
	settings := "settings.json"
	if len(args) > 0 {
		settings = args[0]
	}
	if _, err := os.Stat(settings); err != nil {
		fmt.Fprintf(os.Stderr, `Could not find %s

Try %s --help for usage details
`, settings, appName)
		os.Exit(1)
	}

	cfg, err := dataset.LoadConfig(settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cound not read configuration %q, %s", settings, err)
		os.Exit(1)
	}

	/* Open SQL database holding collections */
	if err := dataset.OpenCollections(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "OpenCollections(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}

	/* Run API */
	if err := dataset.RunAPI(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "RunWebAPI(%q) failed, %s\n", settings, err)
		os.Exit(1)
	}
}
