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
	"encoding/json"

	"github.com/cybergarage/go-ble/ble/db"
)

// Characteristic represents a Bluetooth Characteristic.
type Characteristic interface {
	// UUID returns the Characteristic UUID.
	UUID() uint16
	// Name returns the Characteristic name.
	Name() string
	// ID returns the Characteristic ID.
	ID() string
	// MarshalObject returns an object suitable for marshaling to JSON.
	MarshalObject() any
	// String returns a string representation of the characteristic.
	String() string
}

// nolint: staticcheck
type characteristic struct {
	Uuid uint16 `yaml:"uuid"`
	Nam  string `yaml:"name"`
	Id   string `yaml:"id"`
}

func newCharacteristic(uuid uint16) *characteristic {
	dbChar, _ := db.DefaultDatabase().LookupCharacteristic(uuid)
	return &characteristic{
		Uuid: uuid,
		Nam:  dbChar.Name(),
		Id:   dbChar.ID(),
	}
}

// nolint: tagliatelle
type characteristics struct {
	Characteristics []*characteristic `yaml:"uuids"`
}

// UUID returns the Characteristic UUID.
func (char *characteristic) UUID() uint16 {
	return char.Uuid
}

// Name returns the Characteristic name.
func (char *characteristic) Name() string {
	return char.Nam
}

// ID returns the Characteristic ID.
func (char *characteristic) ID() string {
	return char.Id
}

func (char *characteristic) MarshalObject() any {
	return struct {
		UUID uint16 `json:"uuid"`
		Name string `json:"name"`
		ID   string `json:"id"`
	}{
		UUID: char.UUID(),
		Name: char.Name(),
		ID:   char.ID(),
	}
}

// String returns a string representation of the characteristic.
func (char *characteristic) String() string {
	b, err := json.Marshal(char.MarshalObject())
	if err != nil {
		return ""
	}
	return string(b)
}
