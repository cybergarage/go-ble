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
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// macAddressLen is the length of a Bluetooth MAC address, which Linux
	// and Windows identify a device by.
	macAddressLen = 6
	// uuidAddressLen is the length of a UUID address, which macOS identifies
	// a device by.
	uuidAddressLen = 16
)

// Address represents a Bluetooth address.
//
// The bytes are stored in the little-endian order of the underlying Bluetooth
// stack. A device is identified by a MAC address on Linux and Windows, and by a
// UUID which the system assigns on macOS, so the string representation differs
// between the platforms.
type Address []byte

// NewAddressFromString creates an address from its string representation.
//
// A MAC address is accepted with or without the ':' separators
// (11:22:33:AA:BB:CC), and a UUID with or without the '-' separators.
func NewAddressFromString(s string) (Address, error) {
	hexStr := strings.NewReplacer(":", "", "-", "").Replace(strings.TrimSpace(s))

	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("%w address (%s): %w", ErrInvalid, s, err)
	}

	switch len(b) {
	case macAddressLen, uuidAddressLen:
	default:
		return nil, fmt.Errorf("%w address (%s): %d bytes", ErrInvalid, s, len(b))
	}

	// The string representation is big-endian, and the address holds the
	// bytes in the little-endian order of the Bluetooth stack.
	addr := make(Address, len(b))
	for n, v := range b {
		addr[len(b)-1-n] = v
	}

	return addr, nil
}

// Bytes returns the raw bytes of the address.
func (addr Address) Bytes() []byte {
	return addr
}

// Equal returns true if the address is the same as the specified address.
func (addr Address) Equal(other Address) bool {
	return bytes.Equal(addr, other)
}

// String returns the string representation of the Bluetooth address.
//
// A MAC address is formatted as 11:22:33:AA:BB:CC, and a UUID address as
// 0102030A-0B0C-0D0E-0F10-111213141516.
func (addr Address) String() string {
	// The bytes are stored in the little-endian order, and both
	// representations are written in the big-endian order.
	b := make([]byte, len(addr))
	for n, v := range addr {
		b[len(addr)-1-n] = v
	}

	hexStr := strings.ToUpper(hex.EncodeToString(b))

	switch len(addr) {
	case macAddressLen:
		macStr := make([]string, 0, macAddressLen)
		for n := 0; n < len(hexStr); n += 2 {
			macStr = append(macStr, hexStr[n:n+2])
		}
		return strings.Join(macStr, ":")
	case uuidAddressLen:
		return strings.Join([]string{
			hexStr[0:8],
			hexStr[8:12],
			hexStr[12:16],
			hexStr[16:20],
			hexStr[20:32],
		}, "-")
	}

	return hexStr
}
