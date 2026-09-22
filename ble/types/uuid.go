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

package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cybergarage/go-safecast/safecast"
	"github.com/google/uuid"
)

// UUID represents a Bluetooth UUID.
type UUID [4]uint32

var (
	nilUUID = UUID{0x00000000, 0x00000000, 0x00000000, 0x00000000}
	// baseUUID is the Bluetooth Base UUID, 00000000-0000-1000-8000-00805F9B34FB.
	// The words are held in the big-endian order, so the first word holds
	// the leading bytes of the UUID, as Bytes() and String() write them.
	baseUUID = UUID{0x00000000, 0x00001000, 0x80000080, 0x5F9B34FB}
)

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

// NewUUIDFromString creates a UUID from a string.
//
// A full UUID is accepted, such as 0000FFF6-0000-1000-8000-00805F9B34FB, and so
// is the short form of an assigned 16-bit or 32-bit UUID, with or without the
// "0x" prefix, such as FFF6 or 0xFFF6. The short forms are expanded with the
// Bluetooth Base UUID.
func NewUUIDFromString(s string) (UUID, error) {
	shortStr := strings.TrimSpace(s)
	shortStr = strings.TrimPrefix(strings.TrimPrefix(shortStr, "0x"), "0X")

	switch len(shortStr) {
	case 4:
		u16, err := strconv.ParseUint(shortStr, 16, 16)
		if err == nil {
			return NewUUIDFromUUID16(uint16(u16)), nil
		}
	case 8:
		u32, err := strconv.ParseUint(shortStr, 16, 32)
		if err == nil {
			return NewUUIDFromUUID32(uint32(u32)), nil
		}
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return nilUUID, err
	}

	return NewUUIDFromUUID(u), nil
}

// MustUUIDFromString creates a UUID from a string and return nil UUID if it fails.
func MustUUIDFromString(s string) UUID {
	u, err := NewUUIDFromString(s)
	if err != nil {
		return nilUUID
	}
	return u
}

// NewUUIDFromBytes creates a UUID from a byte slice.
func NewUUIDFromBytes(b []byte) (UUID, error) {
	if len(b) != 16 {
		return nilUUID, fmt.Errorf("invalid UUID byte length: %d", len(b))
	}
	var u uuid.UUID
	copy(u[:], b)
	return NewUUIDFromUUID(u), nil
}

// NewUUIDFromUUID16 creates a UUID from a 16-bit UUID.
func NewUUIDFromUUID16(u16 uint16) UUID {
	return UUID{
		uint32(u16),
		baseUUID[1],
		baseUUID[2],
		baseUUID[3],
	}
}

// NewUUIDFromUUID32 creates a UUID from a 32-bit UUID.
func NewUUIDFromUUID32(u32 uint32) UUID {
	return UUID{
		u32,
		baseUUID[1],
		baseUUID[2],
		baseUUID[3],
	}
}

// NewUUIDFrom creates a UUID from various types.
func NewUUIDFrom(v any) (UUID, error) {
	switch v := v.(type) {
	case UUID:
		return v, nil
	case uuid.UUID:
		return NewUUIDFromUUID(v), nil
	case uint16:
		return NewUUIDFromUUID16(v), nil
	case uint32:
		return NewUUIDFromUUID32(v), nil
	case string:
		return NewUUIDFromString(v)
	case []byte:
		return NewUUIDFromBytes(v)
	}

	var uuid16 uint16
	if safecast.To(v, &uuid16) == nil {
		return NewUUIDFromUUID16(uuid16), nil
	}
	var uuid32 uint32
	if safecast.To(v, &uuid32) == nil {
		return NewUUIDFromUUID32(uuid32), nil
	}

	return nilUUID, fmt.Errorf("invalid UUID type: %T (%v)", v, v)
}

// MustUUIDFrom creates a UUID from various types and return nil UUID if it fails.
func MustUUIDFrom(v any) UUID {
	u, err := NewUUIDFrom(v)
	if err != nil {
		return nilUUID
	}
	return u
}

// NewNilUUID creates a nil UUID.
func NewNilUUID() UUID {
	return nilUUID
}

// Equal checks if two UUIDs are equal.
func (u UUID) Equal(other UUID) bool {
	return u == other
}

// IsNil checks if the UUID is a nil UUID.
func (u UUID) IsNil() bool {
	return u.Equal(nilUUID)
}

// IsUUID16 checks if the UUID is a 16-bit UUID.
func (u UUID) IsUUID16() bool {
	return u[1] == baseUUID[1] && u[2] == baseUUID[2] && u[3] == baseUUID[3] && (u[0] == uint32(u[0]&0xFFFF))
}

// IsUUID32 checks if the UUID is a 32-bit UUID.
func (u UUID) IsUUID32() bool {
	return u[1] == baseUUID[1] && u[2] == baseUUID[2] && u[3] == baseUUID[3] && (u[0] != uint32(u[0]&0xFFFF))
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
	return uint16(u[0]), true
}

// UUID32 returns the 32-bit representation of the UUID if it is a 32-bit UUID.
func (u UUID) UUID32() (uint32, bool) {
	if !u.IsUUID32() {
		return 0, false
	}
	return u[0], true
}

// Bytes returns the byte representation of the UUID in big-endian format.
func (u UUID) Bytes() [16]byte {
	var bytes [16]byte
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
