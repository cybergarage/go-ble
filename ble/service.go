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
)

// Service represents a Bluetooth service.
type Service interface {
	// UUID returns the UUID of the service.
	UUID() UUID
	// Name returns the name of the service.
	Name() string
	// Data returns the data of the service.
	Data() []byte
	// MarshalObject returns an object suitable for marshaling to JSON.
	MarshalObject() any
	// String returns a string representation of the service.
	String() string
}

type service struct {
	uuid UUID
	name string
	data []byte
}

func newService(uuid UUID, name string, data []byte) *service {
	return &service{
		uuid: uuid,
		name: name,
		data: data,
	}
}

// UUID returns the UUID of the service.
func (s *service) UUID() UUID {
	return s.uuid
}

// Name returns the name of the service.
func (s *service) Name() string {
	return s.name
}

// Data returns the data of the service.
func (s *service) Data() []byte {
	return s.data
}

// MarshalObject returns an object suitable for marshaling to JSON.
func (s *service) MarshalObject() any {
	return struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
		Data string `json:"data"`
	}{
		UUID: s.uuid.String(),
		Name: s.name,
		Data: hex.EncodeToString(s.data),
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
