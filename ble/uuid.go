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
)

// UUID represents a Bluetooth UUID.
type UUID [4]uint32

var baseUUID = UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x00000000}

// Equal checks if two UUIDs are equal.
func (u UUID) Equal(other UUID) bool {
	return u == other
}

// IsUUID32 checks if the UUID is a 32-bit UUID.
func (u UUID) IsUUID32() bool {
	return u[1] == baseUUID[1] && u[2] == baseUUID[2] && u[3] == baseUUID[3]
}

// IsUUID16 checks if the UUID is a 16-bit UUID.
func (u UUID) IsUUID16() bool {
	return u.IsUUID32() && (u[3] == uint32(uint16(u[3])))
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
