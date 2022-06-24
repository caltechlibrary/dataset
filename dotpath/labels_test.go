package dotpath

import (
	"testing"
)

func TestToLabels(t *testing.T) {
	data := map[string]string{
		".":             "root object",
		".Column_1":     "Column_1",
		".paths[:]":     "paths range",
		".dirs[:].size": "dirs range size",
		".titles[0]":    "titles[0]",
		".titles[0:1]":  "titles[0:1]",
	}
	for input, expected := range data {
		result := ToLabel(input)
		if result != expected {
			t.Errorf("expected %q, got %q, for %q", expected, result, input)
		}
	}
}
