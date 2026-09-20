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

	"tinygo.org/x/bluetooth"
)

// TestUUIDTinygoConversion verifies that converting between this package's
// UUID and tinygo.org/x/bluetooth's UUID preserves equality with UUIDs
// derived from short (16/32-bit) forms, such as MatterServiceUUID. The two
// libraries lay out their four 32-bit words in reverse order of each other,
// so a naive standard (RFC4122 byte order) round trip silently produces a
// different value than NewUUIDFromUUID16/32 for the same short UUID.
func TestUUIDTinygoConversion(t *testing.T) {
	// The Bluetooth SIG base UUID with the 16-bit UUID 0xFFF6 substituted,
	// i.e. what a real BLE stack (macOS/Linux) expands 0xFFF6 into.
	tinyUUID, err := bluetooth.ParseUUID("0000fff6-0000-1000-8000-00805f9b34fb")
	if err != nil {
		t.Fatalf("ParseUUID() error = %v", err)
	}

	want := NewUUIDFromUUID16(0xFFF6)

	got := uuidFromTinygo(tinyUUID)
	if !got.Equal(want) {
		t.Errorf("uuidFromTinygo(%s) = %s, want %s", tinyUUID.String(), got.String(), want.String())
	}

	back := uuidToTinygo(want)
	if back.String() != tinyUUID.String() {
		t.Errorf("uuidToTinygo(%s) = %s, want %s", want.String(), back.String(), tinyUUID.String())
	}
}
