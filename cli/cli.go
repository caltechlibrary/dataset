/**
 * cli is a package intended to encourage some standardization in the
 * command line user interface for programs developed for Caltech Library.
 *
 * @author R. S. Doiel, <rsdoiel@library.caltech.edu>
 *
 * Copyright (c) 2021, Caltech
 * All rights not granted herein are expressly reserved by Caltech.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// v0.0.17 Added go1.16 support
//
// v0.0.16 Verb.AddParams() was renamed Verb.SetParams()
//
// v0.0.15 merged pkgassets package as subpackage to cli. Removed support for "Actions" replaced by "Verbs".
//
// v0.0.14 removes depreciated pre-v0.0.5 code, adds more flexible
// Verb metaphor to replace Actions. Supports adding options to
// Verbs. Removes AddHelp() called in from AddVerb() and AddAction().
//
// v0.0.13 renames func GenerateMarkdownDocs() to GenerateMarkdown()
// to better match the GeneratgeManPage().
//
// v0.0.5 brings a more wholistic approach to building a cli
// (command line interface), not just configuring one.
//
// Below are the structs and funcs that will remain after v0.0.6-dev
//

// Open accepts a filename, fallbackFile (usually os.Stdout, os.Stdin, os.Stderr) and returns
// a file pointer and error.  It is a conviences function for wrapping stdin, stdout, stderr
// If filename is "-" or filename is "" then fallbackFile is used.
func Open(filename string, fallbackFile *os.File) (*os.File, error) {
	if len(filename) == 0 || filename == "-" {
		return fallbackFile, nil
	}
	return os.Open(filename)
}

// Create accepts a filename, fallbackFile (usually os.Stdout, os.Stdin, os.Stderr) and returns
// a file pointer and error.  It is a conviences function for wrapping stdin, stdout, stderr
// If filename is "-" or filename is "" then fallbackFile is used.
func Create(filename string, fallbackFile *os.File) (*os.File, error) {
	if len(filename) == 0 || filename == "-" {
		return fallbackFile, nil
	}
	return os.Create(filename)
}

// CloseFile accepts a filename and os.File pointer, if filename is "" or "-" it skips the close
// otherwise is does a fp.Close() on the file.
func CloseFile(filename string, fp *os.File) error {
	if len(filename) == 0 || filename == "-" {
		return nil
	}
	return fp.Close()
}

// ReadLines accepts a file pointer to read and returns an array of lines.
func ReadLines(in *os.File) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	return lines, err
}

// IsPipe accepts a file pointer and returns true if data on file pointer is
// from a pipe or false if not.
func IsPipe(in *os.File) bool {
	finfo, err := in.Stat()
	if err == nil && (finfo.Mode()&os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

//NOTE: PopArg was renamed ShiftArg since it is coming off the zero position in array

// ShiftArg takes an array of strings and if array is not empty returns a string and the rest of the args.
// If array is empty it returns an empty string, when there are no more args it returns nil for the
// arg parameter
func ShiftArg(args []string) (string, []string) {
	if args != nil && len(args) >= 1 {
		if len(args) > 1 {
			return args[0], args[1:]
		}
		return args[0], nil
	}
	return "", nil
}

// Unshift takes a string and an array of strings and returns a new array with
// the value of string as the first element and the rest are elements of the array of strings
func UnshiftArg(s string, args []string) []string {
	nArgs := []string{s}
	if args != nil {
		nArgs = append(nArgs, args...)
	}
	return nArgs
}

// PopArg takes an array pops the last element returning the element and an updated
// array. If no elements are left then a empty string and nil are turned.
func PopArg(args []string) (string, []string) {
	if args == nil {
		return "", nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	last := len(args) - 1
	return args[last], args[0:last]
}

// PushArg takes and array and adds a string to the end
func PushArg(s string, args []string) []string {
	if args != nil {
		args = append(args, s)
		return args
	}
	return []string{s}
}

// Cli models the metadata for running a common cli program
type Cli struct {
	// In is usually set to os.Stdin
	In *os.File
	// Out is usually set to os.Stdout
	Out *os.File
	// Eout is usually set to os.Stderr
	Eout *os.File
	// Documentation specific help pages, e.g. -help example1
	Documentation map[string][]byte
	// SectionNo is the numeric value section value for man page generation
	SectionNo int

	/*
		// ActionsRequired is true then the USAGE line shows ACTION rather than [ACTION]
		ActionsRequired bool
	*/

	// VerbsRequired is true then USAGE line shows VERB rather than [VERB]
	VerbsRequired bool

	// application name based on os.Args[0]
	appName string
	// application version based on string passed in New
	version string
	// expected environmental variables used by app
	env map[string]*EnvAttribute
	// description of additoinal command line parameters
	params []string
	// description of short/long options and their doc strings
	options map[string]string

	/*
		// (depreciated) non-flag options, e.g. in the command line "go test", "test" would be the action string.
		// Any additional parameters would be handed of the associated Action.
		actions map[string]*Action
	*/

	// non-flag options, e.g. in the command line "go test", "test"
	// would be the verb string.  Additional options and parameters can be
	// associated with a verb phrase
	verbs map[string]*Verb
}

// NewCli creates an Cli instance, an Cli describes the running of the command line interface
// making it easy to expose the functionality in packages as command line tools.
func NewCli(version string) *Cli {
	appName := path.Base(os.Args[0])
	env := make(map[string]*EnvAttribute)
	options := make(map[string]string)
	/*
		//NOTE: actions is depreciated
		actions := make(map[string]*Action)
	*/
	//NOTE: verbs will eventually replace actions
	verbs := make(map[string]*Verb)
	documentation := make(map[string][]byte)
	sectionNo := 0
	return &Cli{
		In:            os.Stdin,
		Out:           os.Stdout,
		Eout:          os.Stderr,
		Documentation: documentation,
		SectionNo:     sectionNo,
		appName:       appName,
		version:       fmt.Sprintf("%s %s", appName, version),
		env:           env,
		params:        []string{},
		options:       options,
		/*
			actions:       actions,
		*/
		verbs: verbs,
	}
}

// AppName returns the application name as a string
func (c *Cli) AppName() string {
	return c.appName
}

// Verb returns the verb in an arg list without poping the verb.
// If arg list is empty then an empty string is returned.
func (c *Cli) Verb(args []string) string {
	var ok bool
	for _, verb := range args {
		if _, ok = c.verbs[verb]; ok == true {
			return verb
		}
	}
	return ""
}

// AddHelp takes a string keyword and byte slice of content and
// updates the Documentation attribute
func (c *Cli) AddHelp(keyword string, usage []byte) error {
	c.Documentation[keyword] = usage
	_, ok := c.Documentation[keyword]
	if ok == false {
		return fmt.Errorf("could not add help for %q", keyword)
	}
	return nil
}

// Help returns documentation on a topic. If first looks in
// Documentation map and and if not there return a
// not documented string.
func (c *Cli) Help(keywords ...string) string {
	var sections []string

	for _, keyword := range keywords {
		if description, ok := c.Documentation[keyword]; ok == false {
			//FIXME: see if there is a documented verb...
			if verb, ok := c.verbs[keyword]; ok == true {
				description := verb.Help(keywords[1:]...)
				sections = append(sections, fmt.Sprintf("%s\n", description))
				continue
			}
			sections = append(sections, fmt.Sprintf("%q not documented", keyword))
			continue
		} else {
			sections = append(sections, fmt.Sprintf("%s\n\n%s", strings.ToUpper(keyword), description))
		}
	}
	return strings.Join(sections, "\n\n")
}

func splitOps(names string) []string {
	var ops []string

	switch {
	case strings.Contains(names, ","):
		// If we have a comma separated list of options, e.g. "h,help" for -h, -help
		ops = strings.Split(names, ",")
	case strings.Contains(names, " "):
		// If we have a space separated list of options
		ops = strings.Split(names, ",")
	default:
		ops = []string{names}
	}
	// clean up spaces and dash prefixes
	for i := 0; i < len(ops); i++ {
		op := strings.Trim(ops[i], "- ")
		ops[i] = op
	}
	return ops
}

func opsLabel(ops []string) string {
	parts := []string{}
	for _, op := range ops {
		parts = append(parts, "-"+op)
	}
	return strings.Join(parts, ", ")
}

// BoolVar updates c.options doc strings, then splits options and calls flag.BoolVar()
func (c *Cli) BoolVar(p *bool, names string, value bool, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		flag.BoolVar(p, op, value, usage)
	}
}

// IntVar updates c.options doc strings, then splits options and calls flag.IntVar()
func (c *Cli) IntVar(p *int, names string, value int, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.IntVar(p, op, value, usage)
	}
}

// Int64Var updates c.options doc strings, then splits options and calls flag.Int64Var()
func (c *Cli) Int64Var(p *int64, names string, value int64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.Int64Var(p, op, value, usage)
	}
}

// UintVar updates c.options doc strings, then splits options and calls flag.Int64Var()
func (c *Cli) UintVar(p *uint, names string, value uint, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.UintVar(p, op, value, usage)
	}
}

// Uint64Var updates c.options doc strings, then splits options and calls flag.Int64Var()
func (c *Cli) Uint64Var(p *uint64, names string, value uint64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.Uint64Var(p, op, value, usage)
	}
}

// StringVar updates c.options doc strings, then splits options and calls flag.StringVar()
func (c *Cli) StringVar(p *string, names string, value string, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.StringVar(p, op, value, usage)
	}
}

// Float64Var updates c.options doc strings, then splits options and calls flag.Float64Var()
func (c *Cli) Float64Var(p *float64, names string, value float64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.Float64Var(p, op, value, usage)
	}
}

// DurationVar updates c.options doc strings, then splits options and calls flag.DurationVar()
func (c *Cli) DurationVar(p *time.Duration, names string, value time.Duration, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	c.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		flag.DurationVar(p, op, value, usage)
	}
}

// Option returns an option's document string or unsupported string
func (c *Cli) Option(op string) string {
	op = strings.Trim(op, " ")
	for k, doc := range c.options {
		if strings.Contains(k, op) {
			return doc
		}
	}
	return fmt.Sprintf("%q is an unsupported option", op)
}

// Options returns a map of option values and doc strings
func (c *Cli) Options() map[string]string {
	return c.options
}

// ParseOptions envokes flag.Parse() updating variables set in AddOptions
func (c *Cli) ParseOptions() {
	flag.Parse()
	//FIXME: need to parse options for verbs is present...
}

// Parse process both the environment and any flags
func (c *Cli) Parse() error {
	err := c.ParseEnv()
	if err != nil {
		return err
	}
	c.ParseOptions()
	return nil
}

// Args returns flag.Args()
func (c *Cli) Args() []string {
	return flag.Args()
}

// Arg returns an argument by pos index
func (c *Cli) Arg(i int) string {
	return flag.Arg(i)
}

// NArg returns flag.NArg()
func (c *Cli) NArg() int {
	return flag.NArg()
}

// Set Params generates explicit documentation for expected parameters
// instead of using those defined by verbs and actions.
func (c *Cli) SetParams(params ...string) {
	for _, param := range params {
		c.params = append(c.params, param)
	}
}

// NewVerb associates a verb, synopsis and function with a
// command line interface. It supercedes AddVerb(),
// and AddAction(). Verbs can have their own options and
// documentation.
func (c *Cli) NewVerb(name string, usage string, fn func(io.Reader, io.Writer, io.Writer, []string, *flag.FlagSet) int) *Verb {
	verb := NewVerb(name, usage, fn)
	c.verbs[name] = verb
	return verb
}

// Verbs returns a map of verbs and their doc strings
func (c *Cli) Verbs() map[string]string {
	verbs := map[string]string{}
	for k, verb := range c.verbs {
		verbs[k] = verb.Usage
	}
	return verbs
}

// (c *Cli) Run takes a list of non-option arguments and runs them if the fist arg (i.e. arg[0]
// has a corresponding action. Returns an int suitable to passing to os.Exit()
func (c *Cli) Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(c.Eout, "Nothing to do\n")
		return 1
	}
	key, restOfArgs := ShiftArg(args)
	key = strings.TrimSpace(key)
	verb, ok := c.verbs[key]
	if ok == false {
		fmt.Fprintf(c.Eout, "do not known how to %q\n", key)
		return 1
	}
	return verb.Fn(c.In, c.Out, c.Eout, restOfArgs, verb.FlagSet)
}

func padRight(s, p string, maxWidth int) string {
	r := []string{s}
	cnt := maxWidth - len(s)
	for i := 0; i < cnt; i++ {
		r = append(r, p)
	}
	return strings.Join(r, "")
}

func (c *Cli) License() string {
	license, ok := c.Documentation["license"]
	if ok == true {
		return fmt.Sprintf("%s", license)
	}
	return ""
}

func (c *Cli) Version() string {
	return c.version
}
