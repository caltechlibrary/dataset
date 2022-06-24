package dotpath

import (
	"strings"
)

func ToLabel(s string) string {
	if s == "." {
		return "root object"
	}
	if strings.Contains(s, ".") {
		s = strings.Replace(s, ".", " ", -1)
	}
	if strings.Contains(s, "[:]") {
		s = strings.Replace(s, "[:]", " range", -1)
	}
	return strings.TrimSpace(s)
}
