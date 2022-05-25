package dataset

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
