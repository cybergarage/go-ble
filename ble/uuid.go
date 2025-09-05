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
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UUID represents a Bluetooth UUID.
type UUID [4]uint32

var baseUUID = UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x00000000}

// NewUUIDFromUUID creates a UUID from a uuid.UUID.
func NewUUIDFromUUID(uuid uuid.UUID) UUID {
	var b [16]byte = uuid
	return UUID{
		uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3]),
		uint32(b[4])<<24 | uint32(b[5])<<16 | uint32(b[6])<<8 | uint32(b[7]),
		uint32(b[8])<<24 | uint32(b[9])<<16 | uint32(b[10])<<8 | uint32(b[11]),
		uint32(b[12])<<24 | uint32(b[13])<<16 | uint32(b[14])<<8 | uint32(b[15]),
	}
}

// NewUUIDFromString creates a UUID from a string representation.
func NewUUIDFromUUIDString(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return NewUUIDFromUUID(u), nil
}

// NewUUIDFromUUID16 creates a UUID from a 32-bit UUID.
func NewUUIDFromUUID32(u32 uint32) UUID {
	return UUID{
		baseUUID[0],
		baseUUID[1],
		baseUUID[2],
		u32,
	}
}

// NewUUIDFromUUID16 creates a UUID from a 16-bit UUID.
func NewUUIDFromUUID16(u16 uint16) UUID {
	return UUID{
		baseUUID[0],
		baseUUID[1],
		baseUUID[2],
		uint32(u16),
	}
}

// Equal checks if two UUIDs are equal.
func (u UUID) Equal(other UUID) bool {
	return u == other
}

// IsUUID16 checks if the UUID is a 16-bit UUID.
func (u UUID) IsUUID16() bool {
	return u[0] == baseUUID[0] && u[1] == baseUUID[1] && u[2] == baseUUID[2] && (u[3] == uint32(u[3]&0xFFFF))
}

// IsUUID32 checks if the UUID is a 32-bit UUID.
func (u UUID) IsUUID32() bool {
	return u[0] == baseUUID[0] && u[1] == baseUUID[1] && u[2] == baseUUID[2] && (u[3] != uint32(u[3]&0xFFFF))
}

// IsUUID128 checks if the UUID is a 128-bit UUID.
func (u UUID) IsUUID128() bool {
	return !u.IsUUID16() && !u.IsUUID32()
}

// UUID16 returns the 16-bit representation of the UUID if it is a 16-bit UUID.
func (u UUID) UUID16() (uint16, bool) {
	if !u.IsUUID16() {
		return 0, false
	}
	return uint16(u[3]), true
}

// UUID32 returns the 32-bit representation of the UUID if it is a 32-bit UUID.
func (u UUID) UUID32() (uint32, bool) {
	if !u.IsUUID32() {
		return 0, false
	}
	return u[3], true
}

// Bytes returns the byte representation of the UUID in big-endian format.
func (u UUID) Bytes() []byte {
	bytes := make([]byte, 16)
	for i := range u {
		bytes[i*4+0] = byte(u[i] >> 24)
		bytes[i*4+1] = byte(u[i] >> 16)
		bytes[i*4+2] = byte(u[i] >> 8)
		bytes[i*4+3] = byte(u[i] >> 0)
	}
	return bytes
}

// String returns the string representation of the UUID.
func (u UUID) String() string {
	b := u.Bytes()
	return strings.ToUpper(
		fmt.Sprintf("%X-%X-%X-%X-%X",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
		),
	)
}
