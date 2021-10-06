//
// usage.go - handles rendering the "usage" message for cli package
//
package cli

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Usage writes a help page to io.Writer provided. Documentation is based on
// the application's metadata like app name, version, options, actions, etc.
func (c *Cli) Usage(w io.Writer) {
	var parts []string
	parts = append(parts, c.appName)
	if len(c.options) > 0 {
		parts = append(parts, "[OPTIONS]")
	}
	// Add parts defined by params, verbs, or actions
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
	fmt.Fprintf(w, "\nUSAGE: %s\n\n", strings.Join(parts, " "))

	if section, ok := c.Documentation["synopsis"]; ok == true {
		fmt.Fprintf(w, "SYNOPSIS\n\n%s\n\n", bytes.TrimSpace(section))
	}
	if section, ok := c.Documentation["description"]; ok == true {
		fmt.Fprintf(w, "DESCRIPTION\n\n%s\n\n", bytes.TrimSpace(section))
	}

	if len(c.env) > 0 {
		fmt.Fprintf(w, "ENVIRONMENT\n\n")
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
		for _, k := range keys {
			fmt.Fprintf(w, "    %s  %s\n", padRight(k, " ", padding), c.env[k].Usage)
		}
		fmt.Fprintf(w, "\n\n")
	}

	if len(c.options) > 0 {
		fmt.Fprintf(w, "OPTIONS\n\n")
		if len(c.env) > 0 {
			fmt.Fprintf(w, "Options will override any corresponding environment settings\n\n")
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
		for _, k := range keys {
			fmt.Fprintf(w, "    %s  %s\n", padRight(k, " ", padding), c.options[k])
		}
		fmt.Fprintf(w, "\n\n")
	}

	if len(c.verbs) > 0 {
		fmt.Fprintf(w, "VERBS\n\n")
		keys := []string{}
		padding := 0
		for k, _ := range c.verbs {
			keys = append(keys, k)
			if len(k) > padding {
				padding = len(k) + 1
			}
		}
		// Sort the keys alphabetically and display output
		sort.Strings(keys)
		for _, k := range keys {
			usage := c.verbs[k].Usage
			fmt.Fprintf(w, "    %s  %s\n", padRight(k, " ", padding), usage)
			if len(c.verbs[k].params) > 0 {
				if len(c.verbs[k].options) == 0 {
					fmt.Fprintf(w, "    %s   `%s %s %s`\n", padRight("", " ", padding), c.appName, k, strings.Join(c.verbs[k].params, " "))
				} else {
					fmt.Fprintf(w, "    %s   `%s %s [VERB OPTIONS] %s`\n", padRight("", " ", padding), c.appName, k, strings.Join(c.verbs[k].params, " "))
				}
			}
			if len(c.verbs[k].options) > 0 {
				fmt.Fprintf(w, "    %s  verb options:\n", padRight("", " ", padding))
				for op, desc := range c.verbs[k].options {
					fmt.Fprintf(w, "    %s  %s     %s\n", padRight("", " ", padding), padRight(op, " ", padding), desc)
				}
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "\n\n")
	}

	if section, ok := c.Documentation["examples"]; ok == true {
		fmt.Fprintf(w, "EXAMPLES\n\n%s\n\n", bytes.TrimSpace(section))
	}

	if section, ok := c.Documentation["bugs"]; ok == true {
		fmt.Fprintf(w, "BUGS\n\n%s\n\n", bytes.TrimSpace(section))
	}

	if len(c.Documentation) > 0 {
		keys := []string{}
		for k, _ := range c.verbs {
			if _, hasDoc := c.Documentation[k]; hasDoc == true {
				keys = append(keys, k)
			}
		}
		if len(keys) > 0 {
			// Sort the keys alphabetically and display output
			sort.Strings(keys)
			fmt.Fprintf(w, "See %s -help TOPIC for topics - %s\n\n", c.appName, strings.Join(keys, ", "))
		}
	}

	fmt.Fprintf(w, "%s\n", c.version)
}
