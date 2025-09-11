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
	"github.com/cybergarage/go-ble/ble/types"
)

// UUID represents a Bluetooth UUID.
type UUID = types.UUID

// NewUUIDFromUUIDString creates a new UUID from the given UUID string.
func NewUUIDFromUUIDString(s string) (UUID, error) {
	return types.NewUUIDFromString(s)
}

// MustUUIDFromUUIDString creates a UUID from a string representation and return nil UUID if it fails.
func MustUUIDFromUUIDString(s string) UUID {
	return types.MustUUIDFromString(s)
}

// NewUUIDFromUUID16 creates a new UUID from the given 16-bit UUID.
func NewUUIDFromUUID16(u uint16) UUID {
	return types.NewUUIDFromUUID16(u)
}

// NewUUIDFrom creates a UUID from various types.
func NewUUIDFrom(v any) (UUID, error) {
	return types.NewUUIDFrom(v)
}

// MustUUIDFrom creates a UUID from various types and return nil UUID if it fails.
func MustUUIDFrom(v any) UUID {
	return types.MustUUIDFrom(v)
}

// NewNilUUID creates a nil UUID.
func NewNilUUID() UUID {
	return types.NewNilUUID()
}
