// validation.go provides schema-based validation for dataset collections.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2024, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, 
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE 
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR 
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE
// USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package dataset

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation error for a specific field path.
type ValidationError struct {
	// Path is the field path where the error occurred (e.g., "author[0].family")
	Path string `json:"path"`
	// Message describes the validation error
	Message string `json:"message"`
	// Type is the expected type for the field
	Type string `json:"type,omitempty"`
	// Value is the actual value that failed validation (may be truncated)
	Value string `json:"value,omitempty"`
}

// Error implements the error interface.
func (ve *ValidationError) Error() string {
	if ve.Path != "" {
		return fmt.Sprintf("%s: %s", ve.Path, ve.Message)
	}
	return ve.Message
}

// ValidationResult holds the result of validating a record against a schema.
type ValidationResult struct {
	// Valid is true if the record passed all validation checks
	Valid bool `json:"valid"`
	// Errors contains all validation errors (empty if Valid is true)
	Errors []*ValidationError `json:"errors,omitempty"`
}

// String returns a human-readable representation of the validation result.
func (vr *ValidationResult) String() string {
	if vr.Valid {
		return "Validation passed"
	}
	var sb strings.Builder
	sb.WriteString("Validation failed:")
	for _, err := range vr.Errors {
		sb.WriteString(fmt.Sprintf("\n  - %s", err.Error()))
	}
	return sb.String()
}

// ValidateRecord validates a record against a collection's schema.
// Returns a ValidationResult with Valid=true if validation passes,
// or Valid=false with a list of validation errors.
// If the collection has no schema (Model == nil) or validation is disabled
// (Validate == false in config), returns Valid=true with no errors.
func ValidateRecord(c *Collection, data interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []*ValidationError{},
	}

	// No schema or validation disabled = always valid
	if c == nil || c.Model == nil {
		return result
	}

	// Use the model's ValidateInterface method for nested validation
	if !c.Model.ValidateInterface(data) {
		result.Valid = false
		// For now, we return a generic error
		// In the future, we could enhance models package to return detailed errors
		result.Errors = append(result.Errors, &ValidationError{
			Path:    "root",
			Message: "record does not match schema",
		})
	}

	return result
}

// ValidateRecordWithConfig validates a record using the provided config.
// This is useful when the collection's config has validation settings
// that differ from the collection's Model.
func ValidateRecordWithConfig(cfg *Config, data interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []*ValidationError{},
	}

	// No schema or validation disabled = always valid
	if cfg == nil || !cfg.Validate || cfg.Model == nil {
		return result
	}

	// Use the model's ValidateInterface method
	if !cfg.Model.ValidateInterface(data) {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Path:    "root",
			Message: "record does not match schema",
		})
	}

	return result
}

// FormatValidationErrors formats validation errors as a string suitable for
// API error responses or logging.
func FormatValidationErrors(errors []*ValidationError) string {
	if len(errors) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, err := range errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}
