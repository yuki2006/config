// Copyright 2016  The "config" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Context structure handles the parsing of app.conf
// It has a "preferred" section that is checked first for option queries.
// If the preferred section does not have the option, the DEFAULT section is
// checked fallback.
type Context struct {
	config  *Config
	section string // Check this section first, then fall back to DEFAULT
}

func NewContext() *Context {
	return &Context{config: NewDefault()}
}

func LoadContext(confName string, confPaths []string) (*Context, error) {
	var err error
	var conf *Config
	for _, confPath := range confPaths {
		path := filepath.Join(confPath, confName)
		conf, err = ReadDefault(path)
		if err == nil {
			return &Context{config: conf}, nil
		} else if _, isPathErr := err.(*os.PathError); !isPathErr {
			return nil, errors.New(path + " " + err.Error())
		}

	}
	return nil, err
}

func (c *Context) Raw() *Config {
	return c.config
}

func (c *Context) SetSection(section string) {
	c.section = section
}

func (c *Context) SetOption(name, value string) {
	c.config.AddOption(c.section, name, value)
}

func (c *Context) Int(option string) (result int, found bool) {
	result, err := c.config.Int(c.section, option)
	if err == nil {
		return result, true
	}
	if _, ok := err.(OptionError); ok {
		return 0, false
	}

	// If it wasn't an OptionError, it must have failed to parse.
	return 0, false
}

func (c *Context) IntDefault(option string, dfault int) int {
	if r, found := c.Int(option); found {
		return r
	}
	return dfault
}

func (c *Context) Bool(option string) (result, found bool) {
	result, err := c.config.Bool(c.section, option)
	if err == nil {
		return result, true
	}
	if _, ok := err.(OptionError); ok {
		return false, false
	}

	// If it wasn't an OptionError, it must have failed to parse.
	return false, false
}

func (c *Context) BoolDefault(option string, dfault bool) bool {
	if r, found := c.Bool(option); found {
		return r
	}
	return dfault
}

func (c *Context) String(option string) (result string, found bool) {
	if r, err := c.config.String(c.section, option); err == nil {
		return stripQuotes(r), true
	}
	return "", false
}

func (c *Context) StringDefault(option, dfault string) string {
	if r, found := c.String(option); found {
		return r
	}
	return dfault
}

func (c *Context) HasSection(section string) bool {
	return c.config.HasSection(section)
}

// Options returns all configuration option keys.
// If a prefix is provided, then that is applied as a filter.
func (c *Context) Options(prefix string) []string {
	var options []string
	keys, _ := c.config.Options(c.section)
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			options = append(options, key)
		}
	}
	return options
}

// Helpers

func stripQuotes(s string) string {
	if s == "" {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}
