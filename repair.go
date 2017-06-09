package dataset

import (
	"io"
	"fmt"
)

// Check takes a collection name and reports back
// problems that have been identitified in the 
// collections
//
// + checks if collection.json exists and is valid
// + checks if keys.json exits and is valid
// + checks version of collection and version of dataset tool running
// + compares keys.json with k/v pairs in collectio.keymap
// + checks if all collection.buckets exist
// + checks for unaccounted for buckets
// + checks if all keys in collection.keymap exist
// + checks for unaccounted for keys in buckets
// + checks for keys in multiple buckets and reports duplicate record modified times
// 
func Analyze(out io.Writer, collectionName string) error {
	return fmt.Errorf("Analyze not implemented.")
}

// Repair will take a collection name and attempt to recreate
// valid collection.json and keys.json files from content
// in discovered buckets and json documents
func Repair(out io.Writer, collectionName string) error {
	return fmt.Errorf("Repair not implemented")
}

