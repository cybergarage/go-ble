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

// Service represents a Bluetooth service.
type Service interface {
	// UUID returns the UUID of the service.
	UUID() UUID
	// Name returns the name of the service.
	Name() string
	// ID returns the ID of the service.
	ID() string
	// String returns a string representation of the service.
	String() string
}

type service struct {
	uuid UUID
	name string
	id   string
}

func newService(uuid UUID, name string, id string) Service {
	return &service{
		uuid: uuid,
		name: name,
		id:   id,
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

// ID returns the ID of the service.
func (s *service) ID() string {
	return s.id
}

// String returns a string representation of the service.
func (s *service) String() string {
	return fmt.Sprintf("Service[UUID: %s, Name: %s, ID: %s]", s.uuid, s.name, s.id)
}
