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

package db

// Service represents a Bluetooth Service.
type Service interface {
	// UUID returns the Service UUID.
	UUID() UUID
	// Name returns the Service name.
	Name() string
	// ID returns the Service ID.
	ID() string
}

// nolint: staticcheck
type service struct {
	Uuid uint16 `yaml:"uuid"`
	Nam  string `yaml:"name"`
	Id   string `yaml:"id"`
	uuid UUID   `yaml:"-"`
}

// nolint: tagliatelle
type services struct {
	Services []*service `yaml:"uuids"`
}

// UUID returns the Service UUID.
func (s *service) UUID() UUID {
	return s.uuid
}

// Name returns the Service name.
func (s *service) Name() string {
	return s.Nam
}

// ID returns the Service ID.
func (s *service) ID() string {
	return s.Id
}
