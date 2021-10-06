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
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Verb holds the metadata for an keyword action associated with a cli.
// e.g. `go test` would have the verb "test" and the Verb struct would
// define options related metadata associated with that verb.
type Verb struct {
	// Name of the verb
	Name string

	// Usage, a one line description
	Usage string

	// Documentation is a map to help documentation (usually detailed)
	Documentation map[string][]byte

	// SectionNo corresponds to the manual section number (used in generating man pages)
	SectionNo int

	// options holds documentation strings for flags associated with verb
	options map[string]string

	// Fn holds the main function associated with the verb, often is passed
	// stdin, stdout and stnerror returns a value suitable for passing to
	// os.Exit()
	Fn func(io.Reader, io.Writer, io.Writer, []string, *flag.FlagSet) int

	// FlagSet holds the parsable options associated with the verb.
	FlagSet *flag.FlagSet

	// params holds description of non-option command line parameters
	// e.g. for parameters `FILENAME [URL]` the options
	// array would hold "FILENAME", "[URL]".
	params []string
}

// NewVerb creates an Verb instance, and describes the running of the
// command line interface making it easy to expose the functionality
// in packages as command line tools.
func NewVerb(name, usage string, fn func(io.Reader, io.Writer, io.Writer, []string, *flag.FlagSet) int) *Verb {
	options := make(map[string]string)
	documentation := make(map[string][]byte)
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	return &Verb{
		Name:          name,
		Usage:         usage,
		Fn:            fn,
		Documentation: documentation,
		SectionNo:     0,
		FlagSet:       flagSet,
		params:        []string{},
		options:       options,
	}
}

// AddHelp takes a string keyword and byte slice of content and
// updates the Documentation attribute
func (v *Verb) AddHelp(keyword string, usage []byte) error {
	v.Documentation[keyword] = usage
	_, ok := v.Documentation[keyword]
	if ok == false {
		return fmt.Errorf("could not add help for %q", keyword)
	}
	return nil
}

// Help returns documentation on a topic. If first looks in
// Documentation map and if nothing found looks in Synopsis map
// and if not there return an empty string not documented string.
func (v *Verb) Help(keywords ...string) string {
	var sections []string

	if len(keywords) == 0 {
		if len(v.params) > 0 {
			sections = append(sections, fmt.Sprintf("VERB\n\n%s", v.Name))
			sections = append(sections, fmt.Sprintf("    %s %s", v.Name, strings.Join(v.params, " ")))
			if len(v.Usage) != 0 {
				sections = append(sections, v.Usage)
			}
			if len(v.options) > 0 {
				keys := []string{}
				for key, _ := range v.options {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				block := []string{"OPTIONS\n"}
				for _, key := range keys {
					block = append(block, fmt.Sprintf("    %s  %s", key, v.options[key]))
				}
				sections = append(sections, strings.Join(block, "\n"))
			}
			if len(v.Documentation) > 0 {
				block := []string{"DESCRIPTION\n"}
				for keyword, text := range v.Documentation {
					block = append(block, fmt.Sprintf("%s\n   %s", keyword, text))
				}
				sections = append(sections, strings.Join(block, "\n"))
			}
			sections = append(sections, "")
		}
	} else {
		for _, keyword := range keywords {
			if description, ok := v.Documentation[keyword]; ok == false {
				sections = append(sections, fmt.Sprintf("%q not documented", keyword))
				continue
			} else {
				sections = append(sections, fmt.Sprintf("%s\n\n%s", strings.ToUpper(keyword), description))
			}
		}
	}
	return strings.Join(sections, "\n\n")
}

// BoolVar updates v.options doc strings, then splits options and calls v.FlagSet.BoolVar()
func (v *Verb) BoolVar(p *bool, names string, value bool, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		v.FlagSet.BoolVar(p, op, value, usage)
	}
}

// IntVar updates v.options doc strings, then splits options and calls v.FlagSet.IntVar()
func (v *Verb) IntVar(p *int, names string, value int, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.IntVar(p, op, value, usage)
	}
}

// Int64Var updates v.options doc strings, then splits options and calls v.FlagSet.Int64Var()
func (v *Verb) Int64Var(p *int64, names string, value int64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.Int64Var(p, op, value, usage)
	}
}

// UintVar updates v.options doc strings, then splits options and calls v.FlagSet.Int64Var()
func (v *Verb) UintVar(p *uint, names string, value uint, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.UintVar(p, op, value, usage)
	}
}

// Uint64Var updates v.options doc strings, then splits options and calls v.FlagSet.Int64Var()
func (v *Verb) Uint64Var(p *uint64, names string, value uint64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.Uint64Var(p, op, value, usage)
	}
}

// StringVar updates v.options doc strings, then splits options and calls v.FlagSet.StringVar()
func (v *Verb) StringVar(p *string, names string, value string, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.StringVar(p, op, value, usage)
	}
}

// Float64Var updates v.options doc strings, then splits options and calls v.FlagSet.Float64Var()
func (v *Verb) Float64Var(p *float64, names string, value float64, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.Float64Var(p, op, value, usage)
	}
}

// DurationVar updates v.options doc strings, then splits options and calls v.FlagSet.DurationVar()
func (v *Verb) DurationVar(p *time.Duration, names string, value time.Duration, usage string) {
	// Prep to hand off to the flag package
	ops := splitOps(names)
	// Save for our internal option documentation
	label := opsLabel(ops)
	v.options[label] = usage
	// process with flag package
	for _, op := range ops {
		op = strings.TrimSpace(op)
		v.FlagSet.DurationVar(p, op, value, usage)
	}
}

// HasOptions returns true if len(v.options) > 0, false otherwise
func (v *Verb) HasOptions() bool {
	if v.options == nil || len(v.options) == 0 {
		return false
	}
	return true
}

// Option returns an option's document string or unsupported string
func (v *Verb) Option(op string) string {
	op = strings.Trim(op, " ")
	for k, doc := range v.options {
		if strings.Contains(k, op) {
			return doc
		}
	}
	return fmt.Sprintf("%q is an unsupported option", op)
}

// Options returns a map of option values and doc strings
func (v *Verb) Options() map[string]string {
	return v.options
}

// Parse envokes v.FlagSet.Parse() updating variables set in AddOptions
func (v *Verb) Parse(args []string) error {
	return v.FlagSet.Parse(args)
}

// Args returns v.FlagSet.Args()
func (v *Verb) Args() []string {
	return v.FlagSet.Args()
}

// Arg returns an argument by pos index
func (v *Verb) Arg(i int) string {
	return v.FlagSet.Arg(i)
}

// NArg returns v.FlagSet.NArg()
func (v *Verb) NArg() int {
	return v.FlagSet.NArg()
}

// SetParams documents any parameters not defined as Options,
// it is an orders list of strings
func (v *Verb) SetParams(params ...string) {
	for _, param := range params {
		v.params = append(v.params, param)
	}
}

// String prints an actions' verb and description
func (v *Verb) String() string {
	return fmt.Sprintf("%s - %s", v.Name, v.Usage)
}

// AddFlagSet adds/replaces a flag set associated with verbs
func (v *Verb) AddFlagSet(flagSet *flag.FlagSet) {
	v.FlagSet = flagSet
}
