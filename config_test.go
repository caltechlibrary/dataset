// config_test is a part of dataset
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package dataset

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestConfigOpen(t *testing.T) {
	dName := "testout"
	if _, err := os.Stat(dName); err == nil {
		os.RemoveAll(dName)
	} else {
		os.RemoveAll(dName)
		os.MkdirAll(dName, 0775)
	}
	htdocs := path.Join(dName, "htdocs")
	if _, err := os.Stat(htdocs); os.IsNotExist(err) {
		os.MkdirAll(htdocs, 0775)
	}
	fName := path.Join(dName, "settings.json")
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	// Write out sample test settings.json
	src := []byte(`{
    "host": "localhost:8000",
    "htdocs": "testout/htdocs",
    "collections": [
        { 
            "dsn_uri": "DB_USER:DB_PASSWORD@/DB_NAME",
            "dataset": "api_test.ds"
        }
    ]
}`)
	if err := ioutil.WriteFile(fName, src, 0664); err != nil {
		t.Errorf("Failed to generate %q, %s", fName, err)
		t.FailNow()
	}

	settings, err := ConfigOpen(fName)
	if err != nil {
		t.Errorf("ConfigOpen(%q) failed, %s", fName, err)
	}
	if settings == nil {
		t.Errorf("Configuration is nil")
	}
	if settings.Host == "" {
		t.Errorf("Expected Host to be set %+v", settings)
	}
	for _, cfg := range settings.Collections {
		if cfg.DsnURI == "" {
			t.Errorf("Expected DSN URI to be set, %+v", cfg)
		}
	}
}

func TestWriteFile(t *testing.T) {
	dName := "testout"
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
	}
	htdocs := path.Join(dName, "htdocs")
	if _, err := os.Stat(htdocs); os.IsNotExist(err) {
	}
	fName := path.Join(dName, "settings-saved.json")
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	settings := new(Settings)
	settings.Host = "localhost:8001"
	settings.Htdocs = htdocs

	cfg := new(Config)
	cfg.DsnURI = "mysql://$DB_USER:$DB_PASSWORD@/$DB_NAME"
	cfg.CName = "api_test.ds"
	cfg.Create = true
	cfg.Read = true
	cfg.Update = true
	cfg.Delete = true
	cfg.Keys = true
	cfg.Attach = true
	cfg.Retrieve = true
	cfg.Prune = true

	settings.Collections = append(settings.Collections, cfg)

	if err := settings.WriteFile(fName, 0664); err != nil {
		t.Errorf("WriteFile(%q, 0664) failed, %s", fName, err)
	}
	tsettings, err := ConfigOpen(fName)
	if err != nil {
		t.Errorf("expected to ConfigOpen(%q) saved config, %s", fName, err)
		t.FailNow()
	}
	if tsettings == nil {
		t.Errorf("read settings it was nil")
		t.FailNow()
	}
	if len(tsettings.Collections) != 1 {
		t.Errorf("Expected a single collection definition")
		t.FailNow()
	}
	if settings.Host != tsettings.Host {
		t.Errorf("expected Host %q, got %q", settings.Host, tsettings.Host)
	}
	if settings.Htdocs != tsettings.Htdocs {
		t.Errorf("expected Htdocs %q, got %q", settings.Htdocs, tsettings.Htdocs)
	}
	tcfg := tsettings.Collections[0]
	if cfg.DsnURI != tcfg.DsnURI {
		t.Errorf("expected DsnURI %q, got %q", cfg.DsnURI, tcfg.DsnURI)
	}
	if cfg.Create != tcfg.Create {
		t.Errorf("expected Create %t, got %t", cfg.Create, tcfg.Create)
	}
	if cfg.Read != tcfg.Read {
		t.Errorf("expected Read %t, got %t", cfg.Read, tcfg.Read)
	}
	if cfg.Update != tcfg.Update {
		t.Errorf("expected Update %t, got %t", cfg.Update, tcfg.Update)
	}
	if cfg.Delete != tcfg.Delete {
		t.Errorf("expected Delete %t, got %t", cfg.Delete, tcfg.Delete)
	}
	if cfg.Keys != tcfg.Keys {
		t.Errorf("expected Keys %t, got %t", cfg.Keys, tcfg.Keys)
	}
	if cfg.Attach != tcfg.Attach {
		t.Errorf("expected Attach %t, got %t", cfg.Attach, tcfg.Attach)
	}
	if cfg.Retrieve != tcfg.Retrieve {
		t.Errorf("expected Retrieve %t, got %t", cfg.Retrieve, tcfg.Retrieve)
	}
	if cfg.Prune != tcfg.Prune {
		t.Errorf("expected Prune %t, got %t", cfg.Prune, tcfg.Prune)
	}
}
