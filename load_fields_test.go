package dataset

import (
	"strings"
	"testing"
)

// TestLoadAcceptsBothFieldNames is tier A4: the installed Load reads "key" and
// "object", while `dataset help load` says "__key__" and "__object__". Both
// are accepted, mixed in one stream.
func TestLoadAcceptsBothFieldNames(t *testing.T) {
	c, err := Init(t.TempDir()+"/load.ds", "sqlite://collection.db")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	jsonl := strings.Join([]string{
		`{"key":"one","object":{"id":"one"}}`,
		`{"__key__":"two","__object__":{"id":"two"}}`,
		`{"key":"three","__object__":{"id":"three"}}`,
	}, "\n") + "\n"
	if err := c.Load(strings.NewReader(jsonl), false, 0); err != nil {
		t.Fatalf("Load, %s", err)
	}
	for _, k := range []string{"one", "two", "three"} {
		if !c.HasKey(k) {
			t.Errorf("key %q should have loaded", k)
		}
	}
	// A line with neither spelling is still an error.
	if err := c.Load(strings.NewReader(`{"id":"x"}`+"\n"), false, 0); err == nil {
		t.Errorf("a line without a key or object must still fail")
	}
}
