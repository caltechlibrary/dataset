//
// manpage.go - adds man page output formatting suitable for running through
// `nroff -man` and rendering a man page. It is a part of the cli page.
//
// Author: R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2021, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package cli

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	midUnderscoreRE = regexp.MustCompile(`([[:alpha:]0-9]_[[:alpha:]0-9]|\\_)`)
)

// hasMidUnderscore sees if we have an embedded, mid word underscore
// base on re in a slice of letters.
func hasMidUnderscore(s string, cur int) bool {
	start := cur - 1
	end := cur + 2
	if start < 0 || end >= len(s) {
		return false
	}
	return midUnderscoreRE.MatchString(s[start:end])
}

// inlineMarkdown2man scans a line and converts any inline formatting
// (e.g. *word*, _word_ to \fBword\fP, \fIword\fP).
func inlineMarkdown2man(s string) string {
	letters := strings.Split(s, "")
	inAster := false
	inUnderscore := false
	i := 0
	for {
		// Exit loop early
		if i >= len(letters) {
			break
		}
		// Skip any escape characters
		if letters[i] == "\\" {
			i++
			continue
		}
		if letters[i] == "*" {
			// If inAster is true then close the markup with \fP and set inAster to false
			// Else if inAster == false open the markup with \fB and set inAster to true
			if inAster {
				letters[i] = "\\fP"
				inAster = false
			} else {
				letters[i] = "\\fB"
				inAster = true
			}
		} else if letters[i] == "_" && hasMidUnderscore(s, i) == false {
			// If startUnderscore is true then close the markdup with \fP and set startUnderscore to false
			// Else if startUnderscore is false open markup with \fI and set startUnderscore to true
			if inUnderscore {
				letters[i] = "\\fP"
				inUnderscore = false
			} else {
				letters[i] = "\\fI"
				inUnderscore = true
			}
		}
		i++
	}
	return strings.Join(letters, "")
}

// md2man will try to do a crude conversion of Markdown
// to nroff man macros.
func md2man(src []byte) []byte {
	codeBlock := false
	lines := strings.Split(string(src), "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if codeBlock == false {
			// Scan line for formatting conversions
			lines[i] = inlineMarkdown2man(line)
		} else {
			if strings.HasSuffix(line, " \\") {
				lines[i] = line + "\\"
			}
		}
		// Scan line for code block handling
		if strings.HasPrefix(line, "```") {
			if codeBlock {
				lines[i] = ".EP\n"
				codeBlock = false
			} else {
				lines[i] = "\n.EX\n"
				codeBlock = true
			}
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

// GenerateManPage writes a manual page suitable for running through
// `nroff --man`. May need some human clean up depending on content and
// internal formatting (e.g markdown style, spacing, etc.)
func (c *Cli) GenerateManPage(w io.Writer) {
	var parts []string

	// .TH {appName} {section_no} {version} {date}
	fmt.Fprintf(w, ".TH %s %d %q %q\n", c.appName, c.SectionNo, time.Now().Format("2006 Jan 02"), strings.TrimSpace(c.Version()))

	parts = append(parts, fmt.Sprintf(".TP\n\\fB%s\\fP", c.appName))
	if len(c.options) > 0 {
		parts = append(parts, "[OPTIONS]")
	}

	// NOTE: if we've explicitly defined the parameters them here, otherwise
	// extrapoliate from verbs and actions.
	if len(c.params) > 0 {
		parts = append(parts, c.params...)
	}
	if len(c.verbs) > 0 && len(c.params) == 0 {
		if c.VerbsRequired {
			parts = append(parts, "VERB")
		} else {
			parts = append(parts, "[VERB]")
		}
		// Check for verb options...
		for _, verb := range c.verbs {
			if len(verb.options) > 0 {
				parts = append(parts, "[VERB OPTIONS]")
				break
			}
		}
		// Check for verb params
		for _, verb := range c.verbs {
			if len(verb.params) > 0 {
				parts = append(parts, "[VERB PARAMETERS...]")
				break
			}
		}
	}

	// .SH USAGE
	fmt.Fprintf(w, ".SH USAGE\n%s\n", strings.Join(parts, " "))

	// .SH SYNOPSIS
	if section, ok := c.Documentation["synopsis"]; ok == true {
		fmt.Fprintf(w, ".SH SYNOPSIS\n%s\n", md2man(section))
	}
	// .SH DESCRIPTION
	if section, ok := c.Documentation["description"]; ok == true {
		fmt.Fprintf(w, ".SH DESCRIPTION\n%s\n", md2man(section))
	}

	if len(c.options) > 0 {
		fmt.Fprintf(w, ".SH OPTIONS\n")
		parts := []string{}
		if len(c.env) > 0 {
			parts = append(parts, ".TP\nOptions will override any corresponding environment settings.\n")
		}
		if len(parts) > 0 {
			fmt.Fprintf(w, "%s", strings.Join(parts, "\n"))
		}
		fmt.Fprintf(w, ".TP\nThe following options are supported.\n")
		keys := []string{}
		for k, _ := range c.options {
			keys = append(keys, k)
		}
		// Sort the keys alphabetically and display output
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, ".TP\n\\fB%s\\fP\n%s\n", k, c.options[k])
		}
	}

	if len(c.env) > 0 {
		fmt.Fprintf(w, "\n.SS ENVIRONMENT\n")
		if len(c.options) > 0 {
			fmt.Fprintf(w, "Environment variables can be overridden by corresponding options\n")
		}
		keys := []string{}
		for k, _ := range c.env {
			keys = append(keys, k)
		}
		// Sort the keys alphabetically and display output
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, ".TP\n\\fB%s\\fP\n%s\n", k, c.env[k].Usage)
		}
	}

	// .SH EXAMPLES
	if section, ok := c.Documentation["examples"]; ok == true {
		//FIXME: Need to convert Markdown of examples into nroff with
		// with man macros.
		fmt.Fprintf(w, ".SH EXAMPLES\n")
		fmt.Fprintf(w, "\n%s\n", md2man(section))
	}

	/*
		// .SH SEE ALSO
		fmt.Fprintf(w, ".SH SEE ALSO\n")
			if len(c.Documentation) > 0 {
				keys := []string{}
				for k, _ := range c.verbs {
					if k != "description" && k != "examples" && k != "index" {
						keys = append(keys, k)
					}
				}
				if len(keys) > 0 {
					// Sort the keys alphabetically and display output
					sort.Strings(keys)
					links := []string{}
					for _, key := range keys {
						links = append(links, "%s", key)
					}
					fmt.Fprintf(w, ".SH SEE ALSO\n%s\n", strings.Join(links, ", "))
				}
			}
		// .BUGS
		fmt.Fprintf(w, ".SH BUGS\n")
		// AUTHORS
		fmt.Fprintf(w, ".SH AUTHORS\n")
		// COPYRIGHT
		fmt.Fprintf(w, ".SH COPYRIGHT\n")
	*/
}
