package texts

import (
	"bytes"
	"testing"
)

func TestReadSource(t *testing.T) {
	src := []byte(`This is a test file.
It have three lines.
The End.
`)
	buf := bytes.NewBuffer([]byte{})
	filename := ""
	err := WriteSource(filename, buf, src)
	if err != nil {
		t.Errorf("expected to write the buffer, %s", err)
		t.FailNow()

	}
	txt, err := ReadSource(filename, buf)
	if err != nil {
		t.Errorf("expected to read the buffer, %s", err)
		t.FailNow()
	}
	if len(txt) != len(src) {
		t.Errorf("expected src %q, got txt %q", src, txt)
	}
}
