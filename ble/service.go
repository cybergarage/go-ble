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
	"encoding/json"
	"strings"
	"sync"

	"github.com/cybergarage/go-ble/ble/db"
)

// Service represents a Bluetooth service.
type Service interface {
	// ServiceDescriptor represents a Bluetooth service descriptor.
	ServiceDescriptor
	// MarshalObject returns an object suitable for marshaling to JSON.
	MarshalObject() any
	// String returns a string representation of the service.
	String() string
}

// ServiceDescriptor represents a Bluetooth service descriptor.
type ServiceDescriptor interface {
	// Device returns the device that the service belongs to.
	Device() Device
	// UUID returns the UUID of the service.
	UUID() UUID
	// Name returns the name of the service.
	Name() string
	// Data returns the data of the service.
	Data() []byte
	// LookupCharacteristic looks up a characteristic by UUID.
	LookupCharacteristic(uuid UUID) (Characteristic, bool)
	// Characteristics returns the characteristics of the service.
	Characteristics() []Characteristic
}

type service struct {
	dev Device
	db.Service
	uuid    UUID
	data    []byte
	charMap sync.Map
}

func newService(dev Device, uuid UUID, data []byte, chars []Characteristic) *service {
	dbService, _ := db.DefaultDatabase().LookupService(uuid)
	s := &service{
		Service: dbService,
		dev:     dev,
		uuid:    uuid,
		data:    data,
		charMap: sync.Map{},
	}
	for _, char := range chars {
		s.charMap.Store(char.UUID(), char)
	}
	return s
}

// Device returns the device that the service belongs to.
func (s *service) Device() Device {
	return s.dev
}

// UUID returns the UUID of the service.
func (s *service) UUID() UUID {
	return s.uuid
}

// Data returns the data of the service.
func (s *service) Data() []byte {
	return s.data
}

// LookupCharacteristic looks up a characteristic by UUID.
func (s *service) LookupCharacteristic(uuid UUID) (Characteristic, bool) {
	char, ok := s.charMap.Load(uuid)
	if !ok {
		return nil, false
	}
	c, ok := char.(Characteristic)
	if !ok {
		return nil, false
	}
	return c, true
}

// Characteristics returns the characteristics of the service.
func (s *service) Characteristics() []Characteristic {
	var chars []Characteristic
	s.charMap.Range(func(key, value any) bool {
		char, ok := value.(Characteristic)
		if ok {
			chars = append(chars, char)
		}
		return true
	})
	return chars
}

// addDeviceCharacteristic adds a characteristic to the service.
func (s *service) addDeviceCharacteristic(char Characteristic) {
	s.charMap.Store(char.UUID(), char)
}

// MarshalObject returns an object suitable for marshaling to JSON.
func (s *service) MarshalObject() any {
	charObjs := make([]any, 0)
	s.charMap.Range(func(key, value any) bool {
		char, ok := value.(Characteristic)
		if ok {
			charObjs = append(charObjs, char.MarshalObject())
		}
		return true
	})
	return struct {
		UUID            string `json:"uuid"`
		Name            string `json:"name"`
		Data            string `json:"data"`
		Characteristics []any  `json:"characteristics"`
	}{
		UUID:            s.uuid.String(),
		Name:            s.Name(),
		Data:            strings.ToUpper(hex.EncodeToString(s.data)),
		Characteristics: charObjs,
	}
}

// String returns a string representation of the service.
func (s *service) String() string {
	b, err := json.Marshal(s.MarshalObject())
	if err != nil {
		return "{}"
	}
	return string(b)
}
