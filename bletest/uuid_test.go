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

	"github.com/cybergarage/go-ble/ble"
)

func TestUUID(t *testing.T) {
	tests := []struct {
		input    string
		Is16Bit  bool
		Is32Bit  bool
		Is128Bit bool
	}{
		{"5F9B34FB-8000-0080-0000-10000000FD3D", true, false, false}, // 16-bit UUID
		{"5F9B34FB-8000-0080-0000-10000001ABCD", false, true, false}, // 32-bit UUID
		{"00000000-8000-0080-0000-10000001ABCD", false, false, true}, // 128-bit UUID
	}

	for _, test := range tests {
		uuid, err := ble.NewUUIDFromUUIDString(test.input)
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
}
