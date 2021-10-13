package cli

import (
	"flag"
	"io"
	"testing"
)

func TestVerb(t *testing.T) {
	vName := "TestVerb"
	vUsage := "test verb"
	fn := func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		out.Write([]byte("Hello World!"))
		return 0
	}
	verb := NewVerb(vName, vUsage, fn)
	if verb == nil {
		t.Errorf("Expected an 'Verb' struct, got nil")
		t.FailNow()
	}

	var (
		userName string
		count    int
	)

	expectedUserS := "bessie.smith"
	expectedS := "set the username overridding the enviroment"
	verb.StringVar(&userName, "u,user", expectedUserS, expectedS)
	gotS := verb.Option("u")
	if expectedS != gotS {
		t.Errorf("expected %q, got %q", expectedS, gotS)
	}
	gotS = verb.Option("user")
	if expectedS != gotS {
		t.Errorf("expected %q, got %q", expectedS, gotS)
	}
	if expectedUserS != userName {
		t.Errorf("expected %s, got %s", expectedUserS, userName)
	}

	expectedI := 3
	expectedS = "count is an integer"
	verb.IntVar(&count, "c,count", expectedI, expectedS)
	gotS = verb.Option("c")
	if expectedS != gotS {
		t.Errorf("expected (return string) %q, got %q", expectedS, gotS)
	}
	if expectedI != count {
		t.Errorf("expected (count) %d, got %d", expectedI, count)
	}
}
