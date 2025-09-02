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

import (
	"encoding/hex"
	"fmt"
)

// Manufacturer represents a Bluetooth manufacturer.
type Manufacturer interface {
	// Company returns the Bluetooth company information of the manufacturer.
	Company() Company
	// Data returns the manufacturer specific data.
	Data() []byte
	// String returns a string representation of the manufacturer.
	String() string
}

type manufacturer struct {
	company Company
	data    []byte
}

func newManufacturer(id int, data []byte) Manufacturer {
	company, _ := DefaultDatabase().LookupCompany(id)
	return &manufacturer{
		company: company,
		data:    data,
	}
}

// Company returns the Bluetooth company information of the manufacturer.
func (m *manufacturer) Company() Company {
	return m.company
}

// Data returns the manufacturer specific data.
func (m *manufacturer) Data() []byte {
	return m.data
}

// String returns a string representation of the manufacturer.
func (m *manufacturer) String() string {
	return fmt.Sprintf("%s, Data: %x)", m.company.String(), hex.EncodeToString(m.data))
}
