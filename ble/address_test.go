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
	"testing"
)

func TestAddressString(t *testing.T) {
	tests := []struct {
		name     string
		addr     Address
		expected string
	}{
		{
			// The bytes are stored in the little-endian order of the
			// Bluetooth stack, and the string is big-endian.
			name:     "mac",
			addr:     Address{0xCC, 0xBB, 0xAA, 0x33, 0x22, 0x11},
			expected: "11:22:33:AA:BB:CC",
		},
		{
			name:     "uuid",
			addr:     Address{0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x03, 0x02, 0x01},
			expected: "0102030A-0B0C-0D0E-0F10-111213141516",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.addr.String() != test.expected {
				t.Errorf("%s != %s", test.addr.String(), test.expected)
			}

			// The string representation round trips.
			addr, err := NewAddressFromString(test.expected)
			if err != nil {
				t.Fatal(err)
			}
			if !addr.Equal(test.addr) {
				t.Errorf("%v != %v", addr, test.addr)
			}
		})
	}
}

func TestNewAddressFromString(t *testing.T) {
	tests := []struct {
		str     string
		isValid bool
	}{
		{"11:22:33:AA:BB:CC", true},
		{"112233AABBCC", true},
		{"0102030a-0b0c-0d0e-0f10-111213141516", true},
		{"0102030A0B0C0D0E0F10111213141516", true},
		{"11:22:33", false},
		{"", false},
		{"zz:22:33:AA:BB:CC", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			_, err := NewAddressFromString(test.str)
			if test.isValid && err != nil {
				t.Error(err)
			}
			if !test.isValid && err == nil {
				t.Errorf("%s should be invalid", test.str)
			}
		})
	}
}
