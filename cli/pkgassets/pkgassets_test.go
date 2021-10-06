package pkgassets

import (
	"io/ioutil"
	"testing"
)

func TestFile(t *testing.T) {
	txt, err := ioutil.ReadFile("testdata/helloworld.md")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if _, err := ByteArrayToDecl(txt); err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
}
