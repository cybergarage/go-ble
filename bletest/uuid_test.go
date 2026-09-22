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

// The Bluetooth Base UUID is 00000000-0000-1000-8000-00805F9B34FB, and a 16-bit
// or a 32-bit UUID is the base UUID with its leading bytes replaced.
func TestParseUUID(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tests := []struct {
			input    string
			Is16Bit  bool
			Is32Bit  bool
			Is128Bit bool
		}{
			{"0000FD3D-0000-1000-8000-00805F9B34FB", true, false, false}, // 16-bit UUID
			{"0000FFF6-0000-1000-8000-00805F9B34FB", true, false, false}, // 16-bit UUID
			{"0001ABCD-0000-1000-8000-00805F9B34FB", false, true, false}, // 32-bit UUID
			{"0001ABCD-0000-1000-8000-000000000000", false, false, true}, // 128-bit UUID
			{"18EE2EF5-263D-4559-959F-4F9C429F9D11", false, false, true}, // 128-bit UUID
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
			// The string representation round trips.
			if uuid.String() != test.input {
				t.Errorf("Expected %s, got %s", test.input, uuid.String())
			}
		}
	})
}

func TestGenerateUUID(t *testing.T) {
	// 0000FFF6-0000-1000-8000-00805F9B34FB
	uuidFFF6 := ble.UUID{0x0000FFF6, 0x00001000, 0x80000080, 0x5F9B34FB}
	// 0000FD3D-0000-1000-8000-00805F9B34FB
	uuidFD3D := ble.UUID{0x0000FD3D, 0x00001000, 0x80000080, 0x5F9B34FB}
	// 000FFFF6-0000-1000-8000-00805F9B34FB
	uuidFFFF6 := ble.UUID{0x000FFFF6, 0x00001000, 0x80000080, 0x5F9B34FB}

	tests := []struct {
		uuid     ble.UUID
		expected ble.UUID
	}{
		{ble.NewUUIDFromUUID16(0xFFF6), uuidFFF6},
		{ble.MustUUIDFrom(uint16(0xFFF6)), uuidFFF6},
		{ble.MustUUIDFrom(uint32(0xFFF6)), uuidFFF6},
		{ble.MustUUIDFrom(int(0xFFF6)), uuidFFF6},
		{ble.NewUUIDFromUUID32(0xFFFF6), uuidFFFF6},
		{ble.MustUUIDFrom(int(0xFFFF6)), uuidFFFF6},
		// A 16-bit UUID and its expanded string are the same UUID.
		{ble.MustUUIDFrom("0000FFF6-0000-1000-8000-00805F9B34FB"), uuidFFF6},
		{ble.MustUUIDFrom([]byte{0x00, 0x00, 0xFD, 0x3D, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0x80, 0x5F, 0x9B, 0x34, 0xFB}), uuidFD3D},
		{ble.MustUUIDFrom(uuidFD3D), uuidFD3D},
	}

	for _, test := range tests {
		result := test.uuid
		if result != test.expected {
			t.Errorf("Expected UUID %s, got %s", test.expected.String(), result.String())
		}
	}
}

// The database is keyed by the assigned 16-bit UUIDs, so a UUID which a device
// advertises must look up its name.
func TestUUIDDatabaseLookup(t *testing.T) {
	uuid16 := ble.NewUUIDFromUUID16(0xFFF6)
	uuidStr, err := ble.NewUUIDFromString("0000FFF6-0000-1000-8000-00805F9B34FB")
	if err != nil {
		t.Fatal(err)
	}
	if !uuid16.Equal(uuidStr) {
		t.Errorf("%s != %s", uuid16.String(), uuidStr.String())
	}
}
