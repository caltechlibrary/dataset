package dataset

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	fName := "test-settings.json"
	cfg, err := LoadConfig(fName)
	if err != nil {
		t.Errorf("LoadConfig(%q) failed, %s", fName, err)
	}
	if cfg == nil {
		t.Errorf("Configuration is nil")
	}
}
