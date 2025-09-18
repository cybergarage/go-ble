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
	"fmt"
	"time"

	"github.com/cybergarage/go-ble/ble/db"
)

const (
	defaultCharacteristicWriteWithoutResponseWait = time.Duration(500 * time.Millisecond)
)

// OnCharacteristicNotification represents a callback function to be called when a notification is received.
type OnCharacteristicNotification func(char Characteristic, buf []byte)

// Characteristic represents a Bluetooth Characteristic.
type Characteristic interface {
	// CharacteristicDescriptor represents a Bluetooth Characteristic Descriptor.
	CharacteristicDescriptor
	// CharacteristicOperator represents operations that can be performed on a Bluetooth Characteristic.
	CharacteristicOperator
	// MarshalObject returns an object suitable for marshaling to JSON.
	MarshalObject() any
	// String returns a string representation of the characteristic.
	String() string
}

// CharacteristicDescriptor represents a Bluetooth Characteristic Descriptor.
type CharacteristicDescriptor interface {
	// Service returns the service that the characteristic belongs to.
	Service() Service
	// UUID returns the Characteristic UUID.
	UUID() UUID
	// Name returns the Characteristic name.
	Name() string
	// ID returns the Characteristic ID.
	ID() string
}

// CharacteristicOperator represents operations that can be performed on a Bluetooth Characteristic.
type CharacteristicOperator interface {
	// Read reads the characteristic value.
	Read() ([]byte, error)
	// Write writes the characteristic value.
	Write([]byte) (int, error)
	// WriteWithoutResponse writes the characteristic value without waiting for a response.
	WriteWithoutResponse([]byte) (int, error)
	// Notify subscribes to characteristic notifications.
	Notify(OnCharacteristicNotification) error
}

// nolint: staticcheck
type characteristic struct {
	service Service
	Uuid    UUID
	Nam     string
	Id      string
}

func newCharacteristic(service Service, uuid UUID) *characteristic {
	dbChar, _ := db.DefaultDatabase().LookupCharacteristic(uuid)
	return &characteristic{
		service: service,
		Uuid:    uuid,
		Nam:     dbChar.Name(),
		Id:      dbChar.ID(),
	}
}

// Service returns the service that the characteristic belongs to.
func (char *characteristic) Service() Service {
	return char.service
}

// UUID returns the Characteristic UUID.
func (char *characteristic) UUID() UUID {
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

// Read reads the characteristic value.
func (char *characteristic) Read() ([]byte, error) {
	return nil, fmt.Errorf("%w: %s", ErrNotConnected, char.String())
}

// Write writes the characteristic value.
func (char *characteristic) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("%w: %s", ErrNotConnected, char.String())
}

// Notify subscribes to characteristic notifications.
func (char *characteristic) Notify(callback OnCharacteristicNotification) error {
	return fmt.Errorf("%w: %s", ErrNotConnected, char.String())
}

func (char *characteristic) MarshalObject() any {
	return struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
		ID   string `json:"id"`
	}{
		UUID: char.UUID().String(),
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
