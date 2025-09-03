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
	"context"
)

// OnScanResult is the callback function type for scan results.
type OnScanResult func(Device)

// Scanner defines the interface for a Bluetooth scanner.
type Scanner interface {
	// Devices returns the list of discovered devices.
	Devices() []Device
	// Scan starts scanning for Bluetooth devices.
	Scan(ctx context.Context, onResult OnScanResult) error
}
