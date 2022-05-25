package dataset

import (
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	fName := path.Join("testdata", "test-settings.json")
	cfg, err := LoadConfig(fName)
	if err != nil {
		t.Errorf("LoadConfig(%q) failed, %s", fName, err)
	}
	if cfg == nil {
		t.Errorf("Configuration is nil")
	}
	if cfg.Host == "" {
		t.Errorf("Expected Host to be set %+v", cfg)
	}
	if cfg.DSN == "" {
		t.Errorf("Expected DSN to be set, %+v", cfg)
	}
}
