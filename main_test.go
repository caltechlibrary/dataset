package dataset

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setupBlobStoreTests(m)
	setupSearchTests(m)
	os.Exit(m.Run())
}
