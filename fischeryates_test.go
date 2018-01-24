package dataset

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestShuffleStrings(t *testing.T) {
	original := []string{"a", "b", "c", "d", "e", "f", "G", "h", "i", "j", "k", "L", "m", "n"}
	a := []string{"a", "b", "c", "d", "e", "f", "G", "h", "i", "j", "k", "L", "m", "n"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ShuffleStrings(a, r)
	s1 := strings.Join(original, "")
	s2 := strings.Join(a, "")
	if strings.Compare(s1, s2) == 0 {
		t.Errorf("Expected a shuffled string s1 %q, s2 %q", s1, s2)
	}

}
