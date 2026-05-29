//go:build wasip1

package main

import (
	"fmt"
	"os"
	"path/filepath"

	dataset "github.com/caltechlibrary/dataset/v2"
)

// entry holds an open collection and its pre-configured query map.
type entry struct {
	c       *dataset.Collection
	queryFn map[string]string
}

// registry maps collection basename to its open entry.
var registry = map[string]*entry{}

// setupFromConfig reads a settings YAML file, opens every collection listed
// in it, and registers their query maps. Existing open collections are closed
// and re-opened so the config is always authoritative.
func setupFromConfig(configPath string) error {
	settings, err := dataset.ConfigOpen(configPath)
	if err != nil {
		return err
	}
	for _, cfg := range settings.Collections {
		name := filepath.Base(cfg.CName)
		// Close any previously open handle for this name.
		if e, ok := registry[name]; ok {
			e.c.Close()
			delete(registry, name)
		}
		var c *dataset.Collection
		if _, statErr := os.Stat(filepath.Join(cfg.CName, "collection.json")); os.IsNotExist(statErr) {
			c, err = dataset.Init(cfg.CName, cfg.DsnURI)
		} else {
			c, err = dataset.Open(cfg.CName)
		}
		if err != nil {
			return fmt.Errorf("opening %q: %s", cfg.CName, err)
		}
		if cfg.Model != nil {
			c.Model = cfg.Model
		}
		registry[name] = &entry{c: c, queryFn: cfg.QueryFn}
	}
	return nil
}

// getEntry returns the registry entry for cName or an error if not open.
func getEntry(cName string) (*entry, error) {
	e, ok := registry[filepath.Base(cName)]
	if !ok {
		return nil, fmt.Errorf("collection %q is not open", cName)
	}
	return e, nil
}

// openCollection opens an existing collection and adds it to the registry.
// A no-op if already open.
func openCollection(cName string) error {
	name := filepath.Base(cName)
	if _, ok := registry[name]; ok {
		return nil
	}
	c, err := dataset.Open(cName)
	if err != nil {
		return err
	}
	registry[name] = &entry{c: c}
	return nil
}

// closeCollection closes the named collection and removes it from the registry.
func closeCollection(cName string) error {
	name := filepath.Base(cName)
	e, ok := registry[name]
	if !ok {
		return fmt.Errorf("collection %q is not open", cName)
	}
	err := e.c.Close()
	delete(registry, name)
	return err
}

// initCollection creates a new collection and registers it.
func initCollection(cName, dsnURI string) error {
	name := filepath.Base(cName)
	if _, ok := registry[name]; ok {
		return fmt.Errorf("collection %q is already open", name)
	}
	c, err := dataset.Init(cName, dsnURI)
	if err != nil {
		return err
	}
	registry[name] = &entry{c: c}
	return nil
}
