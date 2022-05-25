package dataset

import (
	"bytes"
	"strings"
)

// StringProcessor takes the a text and replaces all the keys
// (e.g. "{app_name}") with their value (e.g. "dataset"). It is
// used to prepare command line and daemon document for display.
func StringProcessor(varMap map[string]string, text string) string {
	src := text[:]
	for key, val := range varMap {
		if strings.Contains(src, key) {
			src = strings.ReplaceAll(src, key, val)
		}
	}
	return src
}

// BytesProcessor takes the a text and replaces all the keys
// (e.g. "{app_name}") with their value (e.g. "dataset"). It is
// used to prepare command line and daemon document for display.
func BytesProcessor(varMap map[string]string, text []byte) []byte {
	src := text[:]
	for key, val := range varMap {
		if bytes.Contains(src, []byte(key)) {
			src = bytes.ReplaceAll(src, []byte(key), []byte(val))
		}
	}
	return src
}
