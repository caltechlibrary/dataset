package dataset

//
// Service configuration management
//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config holds a configuration file structure used by EPrints Extended API
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"hostname"`

	// Collections are defined by a COLLECTION_ID (string)
	// that points at path to where the collection is saved on file system.
	Collections map[string]string `json:"collections"`
}

// DataSource can contain one or more types of datasources. E.g.
// E.g. dsn for MySQL connections and also data for REST API access.
type DataSource struct {
	// DSN is used to connect to a MySQL style DB.
	DSN string `json:"dsn,omitempty"`
	// Rest is used to connect to EPrints REST API
	// NOTE: assumes Basic Auth for authentication
	RestAPI string `json:"rest,omitempty"`
}

func LoadConfig(fname string) (*Config, error) {
	config := new(Config)
	config.Collections = map[string]string{}
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	// Since we should be OK, unmarshal in into active config
	if err = json.Unmarshal(src, config); err != nil {
		return nil, fmt.Errorf("Unmarshaling %q failed, %s", fname, err)
	}
	if config.Hostname == "" {
		config.Hostname = "localhost:8485"
	}
	return config, nil
}
