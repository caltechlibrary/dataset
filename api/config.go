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
	Host string `json:"host,omitempty"`

	// Type identifies the type of SQL storage, e.g. SQLite, MySQL
	Type string `json:"sql_type,omitempty"`

	// DSN points to a file containing a DSN string,
	// e.g. /etc/datasetd/dsn.conf if none is provided then
	// environment variables for DATASET_DSN_URI.
	DSN string `json:"dsn,omitemtpy"`
}

func (config *Config) String() string {
	src, _ := json.MarshalIndent(config, "", "    ")
	return fmt.Sprintf("%s", src)
}

// LoadConfig reads the JSON configuration file provided, validates it
// and either returns a Config structure or error.
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
	if config.Type == "" {
		config.Type = "mysq"
	}
	if config.DSN == "" {
		config.DSN = os.Getenv("DATASET_DSN_URI")
	}
	return config, nil
}
