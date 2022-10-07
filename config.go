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
)

// Settings holds the specific settings for the web service.
type Settings struct {
	// Host holds the URL to listen to for the web API
	Host string `json:"host"`

	// Htdocs holds the path to static content that will be
	// provided by the web service.
	Htdocs string `json:"htdocs"`

	// Collections holds an array of collection configurations that
	// will be supported by the web service.
	Collections []*Config `json:"collections"`
}

// Config holds the collection specific configuration.
type Config struct {
	// Dname holds the dataset collection name/path.
	CName string `json:"dataset,omitempty"`

	// Dsn URI describes how to connection to a SQL storage engine
	// use by the collection(s).
	// e.g. "sqlite://my_collection.ds/collection.db".
	//
	// The Dsn URI may be past in from the environment via the
	// variable DATASET_DSN_URI. E.g. where all the collections
	// are stored in a common database.
	DsnURI string `json:"dsn_uri,omitemtpy"`

	// Permissions for accessing the collection through the web service
	// At least some of these should be set to true otherwise you
	// don't have much of a web service.

	// Keys lets you get a list of keys in a collection
	Keys bool `json:"keys,omitempty"`

	// Create allows you to add objects to a collection
	Create bool `json:"create,omitempty"`

	// Read allows you to retrive an object from a collection
	Read bool `json:"read,omitempty"`

	// Update allows you to replace objects in a collection
	Update bool `json:"update,omitempty"`

	// Delete allows ytou to removes objects, object versions,
	// and attachments from a collection
	Delete bool `json:"delete,omitempty"`

	// Attachments allows you to attached documents for an object in the
	// collection.
	Attachments bool `json:"attachments,omitempty"`

	// Attach allows you to store an attachment for an object in
	// the collection
	Attach bool `json:"attach,omitempty"`

	// Retrieve allows you to get an attachment in the collection for
	// a given object.
	Retrieve bool `json:"retrieve,omitempty"`

	// Prune allows you to remove an attachment from an object in
	// a collection
	Prune bool `json:"prune,omitempty"`

	// FrameRead allows you to see a list of frames, check for
	// a frame's existance and read the content of a frame, e.g.
	// it's definition, keys, object list.
	FrameRead bool `json:"frame_read,omitempty"`

	// FrameWrite allows you to create a frame, change the frame's
	// content or remove the frame completely.
	FrameWrite bool `json:"frame_write,omitempty"`

	// Versions allows you to list versions, read and delete
	// versioned objects and attachments in a collection.
	Versions bool `json:"versions,omitempty"`
}

// String renders the configuration as a JSON string.
func (settings *Settings) String() string {
	src, _ := json.MarshalIndent(settings, "", "    ")
	return fmt.Sprintf("%s", src)
}

// ConfigOpen reads the JSON configuration file provided, validates it
// and returns a Settings structure and error.
//
// NOTE: if the dsn string isn't specified
//
// ```
//
//	settings := "settings.json"
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
	if err = json.Unmarshal(src, settings); err != nil {
		return nil, fmt.Errorf("Unmarshaling %q failed, %s", fName, err)
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
//	fName := "new-settings.json"
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
	src, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, src, perm)
}
