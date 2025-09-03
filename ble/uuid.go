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
	"strings"

	"github.com/google/uuid"
)

// UUID represents a Bluetooth UUID.
type UUID uuid.UUID

func newUUIDFromBytes(b []byte) (UUID, error) {
	addr, err := uuid.FromBytes(b)
	if err != nil {
		return UUID{}, err
	}
	return UUID(addr), nil
}

func mustUUIDFromBytes(b []byte) UUID {
	addr, err := newUUIDFromBytes(b)
	if err == nil {
		return addr
	}
	return UUID(uuid.Nil)
}

// Equal checks if two UUIDs are equal.
func (u UUID) Equal(other UUID) bool {
	return uuid.UUID(u) == uuid.UUID(other)
}

// String returns the string representation of the UUID.
func (u UUID) String() string {
	return strings.ToUpper(uuid.UUID(u).String())
}
