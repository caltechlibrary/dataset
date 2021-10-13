/**
 * cli is a package intended to encourage some standardization in the command line user interface for programs
 * developed for Caltech Library.
 *
 * @author R. S. Doiel, <rsdoiel@caltech.edu>
 *
 * Copyright (c) 2021, Caltech
 * All rights not granted herein are expressly reserved by Caltech.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvAttribute describes expected environmental attributes associated with the cli app.
// It also provides the default value of the environmental attribute if missing from the environment.
type EnvAttribute struct {
	// Name is the environment variable (e.g. DATASET, USERNAME)
	Name string
	// Type holds the type name of the attribute, e.g. int, int64, float64, string, bool, uint, uint64, time.Duration
	Type string
	// BoolValue holds the default boolean
	BoolValue bool
	// IntValue holds the default int
	IntValue int
	// Int64Value holds the default int64
	Int64Value int64
	// UintValue holds the default uint
	UintValue uint
	// Uint64Value holds the default uint64
	Uint64Value uint64
	// Float64Value holds the default float64
	Float64Value float64
	// Dura1tionValue holds the default time.Duration
	DurationValue time.Duration
	// StringValue holds the default string
	StringValue string
	// Usage describes the environment variable role and expected setting
	Usage string
}

// EnvBool adds an environment variable which is evaluate before evaluating options
// returns a pointer to the value.
func (c *Cli) EnvBool(name string, value bool, usage string) *bool {
	c.env[name] = &EnvAttribute{
		Name:      name,
		Type:      fmt.Sprintf("%T", value),
		BoolValue: value,
		Usage:     usage,
	}
	// FIXME: make sure i am creating a point to the boolean value in the map.
	var p *bool
	p = &c.env[name].BoolValue
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvBoolVar adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.BoolVar()
func (c *Cli) EnvBoolVar(p *bool, name string, value bool, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:      name,
		Type:      fmt.Sprintf("%T", value),
		BoolValue: value,
		Usage:     usage,
	}
	p = &c.env[name].BoolValue
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvInt adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.IntVar()
func (c *Cli) EnvInt(name string, value int, usage string) *int {
	c.env[name] = &EnvAttribute{
		Name:     name,
		Type:     fmt.Sprintf("%T", value),
		IntValue: value,
		Usage:    usage,
	}
	var p *int
	p = &c.env[name].IntValue
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvIntVar adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.IntVar()
func (c *Cli) EnvIntVar(p *int, name string, value int, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:     name,
		Type:     fmt.Sprintf("%T", value),
		IntValue: value,
		Usage:    usage,
	}
	p = &c.env[name].IntValue
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvInt64 adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.Int64Var()
func (c *Cli) EnvInt64(name string, value int64, usage string) *int64 {
	c.env[name] = &EnvAttribute{
		Name:       name,
		Type:       fmt.Sprintf("%T", value),
		Int64Value: value,
		Usage:      usage,
	}
	var p *int64
	p = &c.env[name].Int64Value
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvInt64Var adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.Int64Var()
func (c *Cli) EnvInt64Var(p *int64, name string, value int64, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:       name,
		Type:       fmt.Sprintf("%T", value),
		Int64Value: value,
		Usage:      usage,
	}
	p = &c.env[name].Int64Value
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvUint adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.UintVar()
func (c *Cli) EnvUint(name string, value uint, usage string) *uint {
	c.env[name] = &EnvAttribute{
		Name:      name,
		Type:      fmt.Sprintf("%T", value),
		UintValue: value,
		Usage:     usage,
	}
	var p *uint
	p = &c.env[name].UintValue
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvUintVar adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.UintVar()
func (c *Cli) EnvUintVar(p *uint, name string, value uint, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:      name,
		Type:      fmt.Sprintf("%T", value),
		UintValue: value,
		Usage:     usage,
	}
	p = &c.env[name].UintValue
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvUint64 adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.Uint64Var()
func (c *Cli) EnvUint64(name string, value uint64, usage string) *uint64 {
	c.env[name] = &EnvAttribute{
		Name:        name,
		Type:        fmt.Sprintf("%T", value),
		Uint64Value: value,
		Usage:       usage,
	}
	var p *uint64

	p = &c.env[name].Uint64Value
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvFloat64 adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.Float64Var()
func (c *Cli) EnvFloat64(name string, value float64, usage string) *float64 {
	c.env[name] = &EnvAttribute{
		Name:         name,
		Type:         fmt.Sprintf("%T", value),
		Float64Value: value,
		Usage:        usage,
	}
	var p *float64

	p = &c.env[name].Float64Value
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvUint64Var adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.Uint64Var()
func (c *Cli) EnvUint64Var(p *uint64, name string, value uint64, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:        name,
		Type:        fmt.Sprintf("%T", value),
		Uint64Value: value,
		Usage:       usage,
	}
	p = &c.env[name].Uint64Value
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvString adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.StringVar()
func (c *Cli) EnvString(name string, value string, usage string) *string {
	c.env[name] = &EnvAttribute{
		Name:        name,
		Type:        fmt.Sprintf("%T", value),
		StringValue: value,
		Usage:       usage,
	}
	var p *string

	p = &c.env[name].StringValue
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvStringVar adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.StringVar()
func (c *Cli) EnvStringVar(p *string, name string, value string, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:        name,
		Type:        fmt.Sprintf("%T", value),
		StringValue: value,
		Usage:       usage,
	}
	*p = c.env[name].StringValue
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvDuration adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.DurationVar()
func (c *Cli) EnvDuration(name string, value time.Duration, usage string) *time.Duration {
	c.env[name] = &EnvAttribute{
		Name:          name,
		Type:          fmt.Sprintf("%T", value),
		DurationValue: value,
		Usage:         usage,
	}
	var p *time.Duration

	p = &c.env[name].DurationValue
	_, ok := c.env[name]
	if ok == false {
		return nil
	}
	return p
}

// EnvDurationVar adds environment variable which is evaluate before evaluating options
// It is the environment counterpart to flag.DurationVar()
func (c *Cli) EnvDurationVar(p *time.Duration, name string, value time.Duration, usage string) error {
	c.env[name] = &EnvAttribute{
		Name:          name,
		Type:          fmt.Sprintf("%T", value),
		DurationValue: value,
		Usage:         usage,
	}
	p = &c.env[name].DurationValue
	_, ok := c.env[name]
	if ok == false {
		return fmt.Errorf("%q could not be added to environment attributes", name)
	}
	return nil
}

// EnvAttribute returns the struct corresponding to the matchine name
func (c *Cli) EnvAttribute(name string) (*EnvAttribute, error) {
	e, ok := c.env[name]
	if ok == false {
		return nil, fmt.Errorf("%q not defined for environment", name)
	}
	return e, nil
}

// Env returns an EnvAttribute documentation string for matching name
func (c *Cli) Env(name string) string {
	e, ok := c.env[name]
	if ok == false {
		return fmt.Sprintf("%q not documented for environment", name)
	}
	return e.Usage
}

// Getenv returns a given environment attribute value as a string
func (c *Cli) Getenv(name string) string {
	var s string
	e, err := c.EnvAttribute(name)
	if err != nil {
		return s
	}
	switch e.Type {
	case "bool":
		return fmt.Sprintf("%t", e.BoolValue)
	case "int":
		return fmt.Sprintf("%d", e.IntValue)
	case "int64":
		return fmt.Sprintf("%d", e.Int64Value)
	case "uint":
		return fmt.Sprintf("%d", e.UintValue)
	case "uint64":
		return fmt.Sprintf("%d", e.Uint64Value)
	case "float64":
		return fmt.Sprintf("%f", e.Float64Value)
	case "time.Duration":
		return fmt.Sprintf("%s", e.DurationValue)
	}
	return e.StringValue
}

// ParseEnv loops through the os environment using os.Getenv() and updates
// c.env EnvAttribute. Returns an error if there is a problem with environment.
func (c *Cli) ParseEnv() error {
	var (
		err error
		u64 uint64
	)
	for k, e := range c.env {
		s := strings.TrimSpace(os.Getenv(k))
		// NOTE: we only parse the environment if it is not an emprt string
		if s != "" {
			switch e.Type {
			case "bool":
				e.BoolValue, err = strconv.ParseBool(s)
			case "int":
				e.IntValue, err = strconv.Atoi(s)
			case "int64":
				e.Int64Value, err = strconv.ParseInt(s, 10, 64)
			case "uint":
				u64, err = strconv.ParseUint(s, 10, 32)
				e.UintValue = uint(u64)
			case "uint64":
				e.Uint64Value, err = strconv.ParseUint(s, 10, 64)
			case "float64":
				e.Float64Value, err = strconv.ParseFloat(s, 64)
			case "time.Duration":
				e.DurationValue, err = time.ParseDuration(s)
			default:
				e.StringValue = s
			}
			if err != nil {
				return fmt.Errorf("%q should be type %q, %s", e.Name, e.Type, err)
			}
		}
		c.env[k] = e
	}
	return err
}
