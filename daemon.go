package dataset

import (
	"fmt"
)

const (
	// jsonSizeLimit is the maximum size JSON object we'll accept via
	// our service. Current 1 MB (2^20)
	jsonSizeLimit = 1048576

	// attachmentSizeLimit is the maximum size of Attachments we'll
	// accept via our service. Current 250 MiB
	attachmentSizeLimit = (jsonSizeLimit * 250)
)

var (
	config *Config
)

func OpenCollections(cfg *Config) error {
	return fmt.Errorf("OpenCollections() not implemented")
}

func RunAPI(cfg *Config) error {
	return fmt.Errorf("RunAPI() not implemented")
}
