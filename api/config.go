//
// api is a submodule of dataset
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
//
package api

//
// Service configuration management used by datasetd. This just gets
// us to the SQL database with the JSON columns, that database contains
// the configuration for each collection in a table called _collections
// where each collection is it's own table.
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config holds the specific settings for a collection.
type Config struct {
	// Host holds the URL to listen to for the web API
	Host string `json:"host"`

	// CName holds the Collection Name.
	// The CName may be passed in from the environment via
	// the environment variable DATASET_CNAME.
	CName string `json:"dataset,omitempty"`

	// Dsn URI describes how to connection to a SQL storage engine.
	// e.g. "sqlite://my_collection.ds/collection.db".
	// The Dsn URI may be past in from the environment via the
	// variable DATASET_DSN_URI.
	DsnURI string `json:"dsn_uri,omitemtpy"`

	// Htdocs holds a path to static content. This content will be
	// made available via the web service if htdocs is valid.
	//
	// NOTE: Htdocs maybe set through the environment via the
	// environment variable DATASET_HTDOCS
	Htdocs string `json:"htdocs,emitempty"`

	// Permissions for access the collection through the web service
	// At least some of these should be set to true otherwise you
	// don't have much of a web service.
	Keys     bool `json:"keys,omitempty"`
	Create   bool `json:"create,omitempty"`
	Read     bool `json:"read,omitempty"`
	Update   bool `json:"update,omitempty"`
	Delete   bool `json:"delete,omitempty"`
	Attach   bool `json:"attach,omitempty"`
	Retrieve bool `json:"retrieve,omitempty"`
	Prune    bool `json:"prune,omitempty"`
}

// String renders the configuration as a JSON string.
func (config *Config) String() string {
	src, _ := json.MarshalIndent(config, "", "    ")
	return fmt.Sprintf("%s", src)
}

// LoadConfig reads the JSON configuration file provided, validates it
// and either returns a Config structure or error.
//
// NOTE: if the dsn string isn't specified
//
// ```
//    settings := "settings.json"
//    cfg, err := api.Config(settings)
//    if err != nil {
//       ...
//    }
// ```
//
func LoadConfig(fName string) (*Config, error) {
	config := new(Config)
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	// Since we should be OK, unmarshal in into active config
	if err = json.Unmarshal(src, config); err != nil {
		return nil, fmt.Errorf("Unmarshaling %q failed, %s", fName, err)
	}
	if config.Host == "" {
		config.Host = "localhost:8485"
	}
	if config.Htdocs == "" {
		config.Htdocs = os.Getenv("DATASET_HTDOCS")
	}
	if config.Htdocs != "" {
		info, err := os.Stat(config.Htdocs)
		if err != nil {
			return nil, fmt.Errorf("error accesss %q, %s", config.Htdocs, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("htdocs needs to be a directory")
		}
	}
	if config.CName == "" {
		config.CName = os.Getenv("DATASET_CNAME")
		if config.CName == "" {
			return nil, fmt.Errorf("missing collection name")
		}
	}
	if config.DsnURI == "" {
		config.DsnURI = os.Getenv("DATASET_DSN_URI")
	}
	return config, nil
}

// SaveConfig will save a configuration to the filename provided.
//
// ```
//   fName := "settings.json"
//   cfg := new(Config)
//   cfg.Host = "localhost:8001"
//   cfg.Keys = true
//   cfg.Create = true
//   cfg.Read = true
//   cfg.Update = true
//   cfg.Delete = true
//   cfg.Attach = false
//   cfg.Retrieve = false
//   cfg.Prune = false
//   cfg.SaveConfig(fName)
// ```
//
func (cfg *Config) SaveConfig(name string) error {
	src, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, src, 664)
}
