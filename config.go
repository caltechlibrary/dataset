// config is a part of dataset
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

//
// Configure a web service for dataset
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/models"
)

// Settings holds the specific settings for the web service.
type Settings struct {
	// Host holds the URL to listen to for the web API
	Host string `json:"host" yaml:"host"`

	// Htdocs holds the path to static content that will be
	// provided by the web service.
	Htdocs string `json:"htdocs" yaml:"htdocs"`

	// Collections holds an array of collection configurations that
	// will be supported by the web service.
	Collections []*Config `json:"collections" yaml:"collections"`
}

// Config holds the collection specific configuration.
type Config struct {
	// Dname holds the dataset collection name/path.
	CName string `json:"dataset,omitempty" yaml:"dataset,omitempty"`

	// Dsn URI describes how to connection to a SQL storage engine
	// use by the collection(s).
	// e.g. "sqlite://my_collection.ds/collection.db".
	//
	// The Dsn URI may be past in from the environment via the
	// variable DATASET_DSN_URI. E.g. where all the collections
	// are stored in a common database.
	DsnURI string `json:"dsn_uri,omitemtpy" yaml:"dsn_uri,omitempty"`

	// QueryFn maps a query name to a SQL statement used to query the
	// dataset collection. Multiple query statments can be defaulted. They
	// Need to conform to the SQL dialect of the store. NOTE: Only collections
	// using SQL stores are supported.
	QueryFn map[string]string `json:"query,omitempty" yaml:"query,omitempty"`

	// Model describes the record structure to store. It is to validate
	// URL encoded POST and PUT tot the collection.
	Model *models.Model `json:"model,omitempty" yaml:"model,omitempty"`

	// SuccessPage is used for form submissions that are succcessful, i.e. HTTP Status OK (200)
	SuccessPage string `json:"success_page,omitempty" yaml:"success_page,omitempty"`

	// FailPage is used to for form submissions that are unsuccessful, i.g. HTTP response other than OK
	FailPage string `json:"fail_page,omitempty" yaml:"fail_page,omitempty"`

	// Permissions for accessing the collection through the web service
	// At least some of these should be set to true otherwise you
	// don't have much of a web service.

	// Keys lets you get a list of keys in a collection
	Keys bool `json:"keys,omitempty" yaml:"keys,omitempty"`

	// Create allows you to add objects to a collection
	Create bool `json:"create,omitempty" yaml:"create,omitempty"`

	// Read allows you to retrive an object from a collection
	Read bool `json:"read,omitempty" yaml:"read,omitempty"`

	// Update allows you to replace objects in a collection
	Update bool `json:"update,omitempty" yaml:"update,omitempty"`

	// Delete allows ytou to removes objects, object versions,
	// and attachments from a collection
	Delete bool `json:"delete,omitempty" yaml:"delete,omitempty"`

	// Attachments allows you to attached documents for an object in the
	// collection.
	Attachments bool `json:"attachments,omitempty" yaml:"attachments,omitempty"`

	// Attach allows you to store an attachment for an object in
	// the collection
	Attach bool `json:"attach,omitempty" yaml:"attach,omitempty"`

	// Retrieve allows you to get an attachment in the collection for
	// a given object.
	Retrieve bool `json:"retrieve,omitempty" yaml:"retreive,omitempty"`

	// Prune allows you to remove an attachment from an object in
	// a collection
	Prune bool `json:"prune,omitempty" yaml:"prune,omitempty"`

	// FrameRead allows you to see a list of frames, check for
	// a frame's existance and read the content of a frame, e.g.
	// it's definition, keys, object list.
	FrameRead bool `json:"frame_read,omitempty" yaml:"frame_read,omitempty"`

	// FrameWrite allows you to create a frame, change the frame's
	// content or remove the frame completely.
	FrameWrite bool `json:"frame_write,omitempty" yaml:"frame_write,omitempty"`

	// Versions allows you to list versions, read and delete
	// versioned objects and attachments in a collection.
	Versions bool `json:"versions,omitempty" yaml:"versions,omitempty"`
}

// String renders the configuration as a JSON string.
func (settings *Settings) String() string {
	src, _ := JSONMarshalIndent(settings, "", "    ")
	return fmt.Sprintf("%s", src)
}

// ConfigOpen reads the JSON or YAML configuration file provided, validates it
// and returns a Settings structure and error.
//
// NOTE: if the dsn string isn't specified
//
// ```
//
//	settings := "settings.yaml"
//	settings, err := ConfigOpen(settings)
//	if err != nil {
//	   ...
//	}
//
// ```
func ConfigOpen(fName string) (*Settings, error) {
	settings := new(Settings)
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	// Setup defaults from the environment.
	defaultHtdocs := os.Getenv("DATASET_HTDOCS")
	defaultHost := os.Getenv("DATASET_HOST")
	defaultDsnURI := os.Getenv("DATASET_DSN_URI")

	// Make sure we have a fallback for Host
	if defaultHost == "" {
		defaultHost = "localhost:8485"
	}

	// Since we should be OK, unmarshal in into active settings
	if strings.HasSuffix(fName, ".yaml") {
		if err = YAMLUnmarshal(src, settings); err != nil {
			return nil, fmt.Errorf("Unmarshaling %q failed, %s", fName, err)
		}
	} else {
		if err = json.Unmarshal(src, settings); err != nil {
			return nil, fmt.Errorf("Unmarshaling %q failed, %s", fName, err)
		}
	}

	// Apply defaults if needed
	if settings.Host == "" {
		settings.Host = defaultHost
	}
	if settings.Htdocs == "" {
		settings.Htdocs = defaultHtdocs
	}
	if settings.Htdocs != "" {
		info, err := os.Stat(settings.Htdocs)
		if err != nil {
			return nil, fmt.Errorf("error accesss %q, %s", settings.Htdocs, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("htdocs needs to be a directory")
		}
	}
	if defaultDsnURI != "" {
		// Propagate the default DsnURI for the collections
		for _, cfg := range settings.Collections {
			// Override the empty DsnURI with the default
			if cfg.DsnURI == "" {
				cfg.DsnURI = defaultDsnURI
			}
		}
	}
	return settings, nil
}

// Write will save a configuration to the filename provided.
//
// ```
//
//	fName := "new-settings.yaml"
//	mysql_dsn_uri := os.Getenv("DATASET_DSN_URI")
//
//	settings := new(Settings)
//	settings.Host = "localhost:8001"
//	settings.Htdocs = "/usr/local/www/htdocs"
//
//	cfg := &Config{
//		DsnURI: mysql_dsn_uri,
//	   CName: "my_collection.ds",
//	   Keys: true,
//	   Create: true,
//	   Read:  true,
//		Update: true
//		Delete: true
//		Attach: false
//		Retrieve: false
//		Prune: false
//	}}
//	settings.Collections = append(settings.Collections, cfg)
//
//	if err := api.WriteFile(fName, 0664); err != nil {
//	   ...
//	}
//
// ```
func (settings *Settings) WriteFile(name string, perm os.FileMode) error {
	var (
		src []byte
		err error
	)
	if strings.HasSuffix(name, ".yaml") {
		if src, err = YAMLMarshal(settings); err != nil {
			return err
		}
	} else {
		if src, err = JSONMarshalIndent(settings, "", "    "); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(name, src, perm)
}


// GetCfg retrieves a collection configuration from a Settings object using dataset name.
func (settings *Settings) GetCfg(cName string) (*Config, error) {
	if settings.Collections != nil && len(settings.Collections) > 0 {
		for _, cfg := range settings.Collections {
			if cfg.CName == cName {
				return cfg, nil
			}
		}
	}
	return nil, fmt.Errorf("%s not found", cName)
}
