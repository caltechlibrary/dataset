//
// generator.go - provides a quick and dirty skeleton for cli based apps.
//
package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	commentBlock = `//
// %s %s
//
// %s
`
	importsBlock = `

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/cli"
)

`

	varBlock = `

var (
	synopsis = %s

	description = %s

	examples = %s 

	bugs = %s

	license = %s

	// Standard Options
	showHelp bool
	showLicense bool
	showVersion bool
	showExamples bool
	inputFName string
	outputFName string
	newLine bool
	quiet bool
	prettyPrint bool
	generateMarkdown bool
	generateManPage bool

	// Application Options
)

`

	mainBlock = `

func main() {
	//FIXME: Replace with your base package .Version attribute
	app := cli.NewCli("v0.0.0")
	//FIXME: if you need the app name then...
	//appName := app.AppName()

	// Set explicit parameter list for USAGE
	//app.SetParams(ARRAY_OF_OPTIONAL_AND_REQUIRE_PARAMETERS...)

	// Add Help Docs
	app.SectionNo = 1 // The manual page section number
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))
	app.AddHelp("bugs", []byte(bugs))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input file name")
	app.StringVar(&outputFName, "o,output", "", "output file name")
	app.BoolVar(&newLine, "nl,newline", false, "if true add a trailing newline")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&prettyPrint, "p,pretty", false, "pretty print output")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "output manpage markup")

	// Application Options
	//FIXME: Add any application specific options

	// Application Verbs
	//FIXME: If the application is verb based add your verbs here e.g. 
	//app.VerbsRequired = true
	//verbObj := app.NewVerb(STRING_VERB, STRING_DESCRIPTION, FUNC_POINTER)
	//verbObj.StringVar(&inputFName, "i,input", "", "input filename")
	// ...

	// We're ready to process args
	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
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
			fmt.Fprintf(app.Out, app.Help(args...))
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

	// Application Logic
	//FIXME: running code, e.g. os.Exit(app.Run(args))

	if newLine {
		fmt.Fprintln(app.Out, "")
	}
}

`
)

func backQuote(s string) string {
	return "`" + s + "\n`"
}

// Generate creates main.go source code
func Generate(appName, synopsis, author, descriptionFName, examplesFName, bugsFName, licenseFName string) []byte {
	var (
		description []byte
		examples    []byte
		license     []byte
		bugs        []byte
		err         error
	)
	if len(descriptionFName) > 0 {
		description, err = ioutil.ReadFile(descriptionFName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping description %q, %s", descriptionFName, err)
		}
	}
	if len(examplesFName) > 0 {
		examples, err = ioutil.ReadFile(examplesFName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping examples %q, %s", examplesFName, err)
		}
	}
	if len(bugsFName) > 0 {
		bugs, err = ioutil.ReadFile(bugsFName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping bugs %q, %s", bugsFName, err)
		}
	}
	if len(licenseFName) > 0 {
		license, err = ioutil.ReadFile(licenseFName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping license %q, %s", licenseFName, err)
		}
	}
	blocks := []string{}
	// Setup the initial comment block
	blocks = append(blocks, fmt.Sprintf(commentBlock, appName, string(synopsis), author))
	// Add the license
	blocks = append(blocks, fmt.Sprintf("//%s\n", strings.Replace(string(license), "\n", "\n// ", -1)))
	// Add Package name
	blocks = append(blocks, "package main")
	// Add Standard Imports
	blocks = append(blocks, importsBlock)
	// Add Global vars
	blocks = append(blocks, fmt.Sprintf(varBlock,
		backQuote(string(synopsis)),
		backQuote(string(description)),
		backQuote(string(examples)),
		backQuote(string(bugs)),
		backQuote(string(license)),
	))
	// Add Main
	blocks = append(blocks, mainBlock)

	// Convert to byte array and return
	return []byte(strings.Join(blocks, ""))
}
