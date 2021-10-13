// markdown.go - this is a part of the cli package. This code focuses on
// generating Markdown docs from the internal help information.
package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// GenerateMarkdown writes a Markdown page to io.Writer provided.
// Documentation is based on the application's metadata like app name,
// version, options, actions, etc.
func (c *Cli) GenerateMarkdown(w io.Writer) {
	var parts []string
	parts = append(parts, c.appName)

	if len(c.options) > 0 {
		parts = append(parts, "[OPTIONS]")
	}
	// NOTE: setup explicit parameter documentation
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
	fmt.Fprintf(w, "\nUSAGE\n=====\n\n	%s\n\n", strings.Join(parts, " "))

	if section, ok := c.Documentation["synopsis"]; ok == true {
		fmt.Fprintf(w, "SYNOPSIS\n--------\n\n%s\n\n", section)
	}

	if section, ok := c.Documentation["description"]; ok == true {
		fmt.Fprintf(w, "DESCRIPTION\n-----------\n\n%s\n\n", section)
	}

	if len(c.env) > 0 {
		fmt.Fprintf(w, "ENVIRONMENT\n-----------\n\n")
		if len(c.options) > 0 {
			fmt.Fprintf(w, "Environment variables can be overridden by corresponding options\n\n")
		}
		keys := []string{}
		padding := 0
		for k, _ := range c.env {
			keys = append(keys, k)
			if len(k) > padding {
				padding = len(k) + 1
			}
		}
		// Sort the keys alphabetically and display output
		sort.Strings(keys)
		fmt.Fprintf(w, "```\n")
		for _, k := range keys {
			fmt.Fprintf(w, "    %s  # %s\n", padRight(k, " ", padding), c.env[k].Usage)
		}
		fmt.Fprintf(w, "```\n\n")
	}

	if len(c.options) > 0 {
		fmt.Fprintf(w, "OPTIONS\n-------\n\n")
		parts := []string{}
		parts = append(parts, "Below are a set of options available.")
		if len(c.env) > 0 {
			parts = append(parts, "Options will override any corresponding environment settings.")
		}
		if len(parts) > 0 {
			fmt.Fprintf(w, "%s\n\n", strings.Join(parts, " "))
		}
		keys := []string{}
		padding := 0
		for k, _ := range c.options {
			keys = append(keys, k)
			if len(k) > padding {
				padding = len(k) + 1
			}
		}
		// Sort the keys alphabetically and display output
		sort.Strings(keys)
		fmt.Fprintf(w, "```\n")
		for _, k := range keys {
			fmt.Fprintf(w, "    %s  %s\n", padRight(k, " ", padding), c.options[k])
		}
		fmt.Fprintf(w, "```\n")
		fmt.Fprintf(w, "\n\n")
	}

	if section, ok := c.Documentation["examples"]; ok == true {
		fmt.Fprintf(w, "EXAMPLES\n--------\n\n%s\n\n", section)
	}

	if len(c.Documentation) > 0 {
		keys := []string{}
		for k, _ := range c.verbs {
			if _, hasDocs := c.Documentation[k]; hasDocs == true {
				keys = append(keys, k)
			}
		}
		if len(keys) > 0 {
			// Sort the keys alphabetically and display output
			sort.Strings(keys)
			links := []string{}
			for _, key := range keys {
				links = append(links, fmt.Sprintf("[%s](%s.html)", key, key))
			}
			fmt.Fprintf(w, "Related: %s\n\n", strings.Join(links, ", "))
		}
	}

	fmt.Fprintf(w, "%s\n", c.version)
}
