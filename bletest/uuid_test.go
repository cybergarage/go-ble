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

package bletest

import (
	"testing"

	ble "github.com/cybergarage/go-ble/ble/types"
)

func TestParseUUID(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    string
			Is16Bit  bool
			Is32Bit  bool
			Is128Bit bool
		}{
			{"5F9B34FB-8000-0080-0000-10000000FD3D", true, false, false}, // 16-bit UUID
			{"5F9B34FB-8000-0080-0000-10000000FFF6", true, false, false}, // 16-bit UUID
			{"5F9B34FB-8000-0080-0000-10000001ABCD", false, true, false}, // 32-bit UUID
			{"00000000-8000-0080-0000-10000001ABCD", false, false, true}, // 128-bit UUID
		}

		for _, test := range tests {
			uuid, err := ble.NewUUIDFromString(test.input)
			if err != nil {
				t.Errorf("Failed to parse UUID %s: %v", test.input, err)
				continue
			}
			if _, ok := uuid.UUID16(); ok != test.Is16Bit {
				t.Errorf("Expected IsUUID16() to be %v for UUID %s, got %v", test.Is16Bit, test.input, uuid.IsUUID16())
			}
			if _, ok := uuid.UUID32(); ok != test.Is32Bit {
				t.Errorf("Expected IsUUID32() to be %v for UUID %s, got %v", test.Is32Bit, test.input, uuid.IsUUID32())
			}
			if uuid.IsUUID128() != test.Is128Bit {
				t.Errorf("Expected IsUUID128() to be %v for UUID %s, got %v", test.Is128Bit, test.input, uuid.IsUUID128())
			}
		}
	})
}

func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		uuid     ble.UUID
		expected ble.UUID
	}{
		{ble.NewUUIDFromUUID16(0xFFF6), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FFF6}},
		{ble.MustUUIDFrom(uint16(0xFFF6)), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FFF6}},
		{ble.MustUUIDFrom(uint32(0xFFF6)), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FFF6}},
		{ble.MustUUIDFrom(int(0xFFF6)), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FFF6}},
		{ble.NewUUIDFromUUID32(0xFFFF6), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x000FFFF6}},
		{ble.MustUUIDFrom(int(0xFFFF6)), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x000FFFF6}},
		{ble.MustUUIDFrom("5F9B34FB-8000-0080-0000-10000000FFF6"), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FFF6}},
		{ble.MustUUIDFrom([]byte{0x5F, 0x9B, 0x34, 0xFB, 0x80, 0x00, 0x00, 0x80, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0xFD, 0x3D}), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FD3D}},
		{ble.MustUUIDFrom(ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FD3D}), ble.UUID{0x5F9B34FB, 0x80000080, 0x00001000, 0x0000FD3D}},
	}

	for _, test := range tests {
		result := test.uuid
		if result != test.expected {
			t.Errorf("Expected UUID %v, got %v", test.expected, result)
		}
	}
}
