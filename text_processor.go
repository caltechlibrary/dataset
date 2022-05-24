package dataset

import (
	"strings"
)

// TextProcessor takes the a text and replaces all the keys
// (e.g. "{app_name}") with their value (e.g. "dataset"). It is
// used to prepare command line and daemon document for display.
func TextProcessor(varMap map[string]string, text string) string {
	src := text[:]
	for key, val := range varMap {
		if strings.Contains(src, key) {
			src = strings.ReplaceAll(src, key, val)
		}
	}
	return src
}
