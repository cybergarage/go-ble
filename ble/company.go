// Copyright (C) 2025 The go-ble Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ble

import "fmt"

// Company represents a Bluetooth company.
type Company interface {
	// ID returns the company ID.
	ID() int
	// Name returns the company name.
	Name() string
	// String returns a string representation of the company.
	String() string
}

type company struct {
	Value int    `yaml:"value"`
	Nam   string `yaml:"name"`
}

// nolint: tagliatelle
type companies struct {
	Companies []*company `yaml:"company_identifiers"`
}

// ID returns the company ID.
func (c *company) ID() int {
	return c.Value
}

// Name returns the company name.
func (c *company) Name() string {
	return c.Nam
}

// String returns a string representation of the company.
func (c *company) String() string {
	return fmt.Sprintf("Company(ID: %d, Name: %s)", c.Value, c.Nam)
}
