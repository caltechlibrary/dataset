package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	// Caltech Library Package
	"github.com/caltechlibrary/dataset/cli"
	"github.com/caltechlibrary/dataset/cli/pkgassets"
)

var (
	synopsis = `
_pkgassets_ generates Go code for mapping path to text assets
`

	description = `
_pkgassets_ generates a Go source directory whos file assets are embedded in a ` + "`" + `map[string][]byte` + "`" + ` variable. 
This is useful where you want to embed web content, template source code, help docs and other assets that 
can be used for default behavior in a Go command line program or service. 

The map content is harvested from directory holding the assets to be embedded. By default the
path key starts with a slash and does not include the hosting directory (e.g. htdocs/index.html 
would become /index.html if htdocs was used to harvest assets). The prefix and suffix on the
key can be modified based on _pkgassets_'s command line options.
`

	examples = `
` + "```" + `
    pkgassets MAP_VARAIBLE_NAME NAME_OF_DIRECTORY_HOLDING_ASSETS
` + "```" + `

This will result in a Go of type map[string][]byte holding the assets discovered by walking the directory
tree provided. The map's key will represent a path (beginning with "/") pointing at the asset ingested.

` + "```shell" + `
    pkgassets DefaultSite htdocs
` + "```" + `

Assuming that _htdocs_ held

+ index.html
+ css/site.css

In this example the htdocs directory will be crawled and all the files found harvested as a an asset. The
path in the map will not include htdocs and would result in a Go source file like

` + "```golang" + `
    package defaultsite

    var DefaultSite = map[string][]byte{
        "/index.html": []byte{}, // ... the contents of index.html would be here ...
        "/css/site.css": []byte{}, // ... the contents of css/site.css would be here ...
    }
` + "```" + `

If a package name is not provided then the package name will a lowercase name of the map variable name (e.g. 
"var DefaultSite" becomes "package defaultsite"). Likewise if a output name is not provided then the file
name will be the name of the package plus the ".go" extension.
`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	showVerbose      bool
	outputFName      string
	generateMarkdown bool
	generateManPage  bool
	quiet            bool

	// App Options
	packageName   string
	commentFName  string
	stripPrefix   string
	stripSuffix   string
	requiredExt   string
	excludeFNames string
)

func isExcluded(fName string, fList []string) bool {
	for _, s := range fList {
		if fName == s {
			return true
		}
	}
	return false
}

func main() {
	app := cli.NewCli(cli.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams("VARIABLE_NAME", "DIR_HOLDING_ASSETS", "[VARIABLE_NAME DIR_HOLDING_ASSETS ...]")

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(pkgassets.LicenseText, appName, pkgassets.Version)))
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showVerbose, "V,verbose", false, "verbose output")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

	// App Options
	app.StringVar(&packageName, "p", "", "package name, if missing defauls to lowercase of variable name")
	app.StringVar(&packageName, "package", "", "package name, if missing defauls to lowercase of variable name")
	app.StringVar(&commentFName, "c", "", "comment file to be included")
	app.StringVar(&commentFName, "comment", "", "comment file to be included")
	app.StringVar(&stripPrefix, "strip-prefix", "", "strip the prefix from the map key")
	app.StringVar(&stripSuffix, "strip-suffix", "", "strip the suffix from the map key")
	app.StringVar(&requiredExt, "ext", "", "Only include files with matching extension")
	app.StringVar(&excludeFNames, "X,exclude", "", "A colon separted list of filenames to exclude, (e.g. 'nav.md:topics.md')")

	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	/*
		app.In, err = cli.Open(inputFName, os.Stdin)
		cli.ExitOnError(app.Eout, err, quiet)
		defer cli.CloseFile(inputFName, app.In)
	*/

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintln(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	if len(args) < 2 {
		app.Usage(app.Eout)
		fmt.Fprintf(app.Eout, "\n")
		if len(args) < 1 {
			fmt.Fprintf(app.Eout, "Missing map variable name\n")
		}
		if len(args) < 2 {
			fmt.Fprintf(app.Eout, "Missing asset directory name\n")
		}
		os.Exit(1)
	}
	if (len(args) % 2) != 0 {
		fmt.Fprintf(app.Eout, "Missing variable/assets directory names must be paired\n")
		fmt.Fprintf(app.Eout, "Paramters provided, %q\n", strings.Join(args, ", "))
		os.Exit(1)
	}

	excluded := strings.Split(excludeFNames, ":")

	// For each pair of mapVName/assetDir add a map to outFName
	for i := 0; (i + 1) < len(args); i += 2 {
		mapVName, assetDir := args[i], args[i+1]
		if mapVName == "" {
			cli.ExitOnError(app.Eout, fmt.Errorf("Expected mapVName to be non-empty stirng for parameter %d", i), quiet)
		}
		if assetDir == "" {
			cli.ExitOnError(app.Eout, fmt.Errorf("Expected assetDir to be non-empty string for parameter %d", i+1), quiet)
		}

		if packageName == "" {
			packageName = strings.ToLower(mapVName)
		}
		if outputFName == "" {
			outputFName = strings.ToLower(packageName) + ".go"
		}
		commentSrc := []byte{}
		if commentFName != "" {
			commentSrc, err = ioutil.ReadFile(commentFName)
			if err != nil {
				cli.ExitOnError(app.Eout, fmt.Errorf("Can't read %s, %s\n", commentFName, err), quiet)
			}
		}

		if len(commentSrc) > 0 {
			fmt.Fprintf(app.Out, "\n//\n// %s\n//\n", bytes.Replace(commentSrc, []byte("\n"), []byte("\n// "), -1))
		}
		if i == 0 {
			fmt.Fprintf(app.Out, "package %s\n\nvar (\n\n", packageName)
		}
		fmt.Fprintf(app.Out, `    // %s is a map to asset files associated with %s package
    %s = map[string][]byte{`, mapVName, packageName, mapVName)
		// Walk the asset directory structure and for each file found at it to the map...
		if err = filepath.Walk(assetDir, func(walkingPath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if isExcluded(path.Base(walkingPath), excluded) {
				if showVerbose {
					fmt.Printf("Excluding %q\n", walkingPath)
				}
				return nil
			}
			if showVerbose {
				fmt.Printf("Processing %q\n", walkingPath)
			}

			fPath := strings.TrimPrefix(walkingPath, assetDir)
			fExt := path.Ext(fPath)
			if stripPrefix != "" {
				fPath = strings.TrimPrefix(fPath, stripPrefix)
			}
			if stripSuffix != "" {
				fPath = strings.TrimSuffix(fPath, stripSuffix)
			}
			if requiredExt == "" || requiredExt == fExt {
				bArray, err := ioutil.ReadFile(walkingPath)
				if err != nil {
					cli.OnError(app.Eout, fmt.Errorf("Can't read %q, %s", walkingPath, err), quiet)
					return nil
				}
				if bSrc, err := pkgassets.ByteArrayToDecl(bArray); err == nil {
					fmt.Fprintf(app.Out, "\n    %q: %s,\n", fPath, bSrc)
				} else {
					cli.OnError(app.Eout, fmt.Errorf("Can't convert to byte array notation %s, %s\n", walkingPath, err), quiet)
				}
			}
			return nil
		}); err != nil {
			cli.OnError(app.Eout, fmt.Errorf("Can't walk path %q, %s\n", assetDir, err), quiet)
		}
		fmt.Fprintf(app.Out, `
	}
`)
	}
	fmt.Fprintf(app.Out, `
)

`)

}
