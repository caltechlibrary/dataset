package dataset

import ()

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = "2006-01-02 15:04:05"
	datestamp = "2006-01-02"

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
