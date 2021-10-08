package dataset

//
// Service configuration management
//

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Config holds a configuration file structure used by EPrints Extended API
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"host" default:"localhost:8485"`

	// Collections are defined by a COLLECTION_ID (string)
	// that points at path to where the collection is saved on file system.
	Collections map[string]*Settings `json:"collections,required"`

	// Routes are mappings of collections to supported routes.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error) `json:"-"`
}

// Settings holds the specific settings for a collection.
type Settings struct {
	CName    string      `json:"dataset,required"`
	Keys     bool        `json:"keys" default:"false"`
	Create   bool        `json:"create" default:"false"`
	Read     bool        `json:"read" default:"false"`
	Update   bool        `json:"update" default:"false"`
	Delete   bool        `json:"delete" default:"false"`
	Attach   bool        `json:"attach" default:"false"`
	Retrieve bool        `json:"retrieve" default:"false"`
	Prune    bool        `json:"prune" default:"false"`
	DS       *Collection `json:"-"`
}

func (config *Config) String() string {
	src, _ := json.MarshalIndent(config, "", "    ")
	return fmt.Sprintf("%s", src)
}

// LoadConfig reads the JSON configuration file provided, validates it
// and either returns a Config structure or error.
func LoadConfig(fname string) (*Config, error) {
	config := new(Config)
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
	if len(config.Collections) == 0 {
		return nil, fmt.Errorf("No collections defined in %s", fname)
	}
	for collectionID, settings := range config.Collections {
		if settings.CName == "" {
			return nil, fmt.Errorf("Settings for %q missing dataset path", collectionID)
		}
		if _, err := os.Stat(settings.CName); err != nil {
			return nil, err
		}
	}
	return config, nil
}
