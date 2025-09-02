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

// Device represents a Bluetooth device.
type Device interface {
	// Manufacturer returns the Bluetooth manufacturer of the device.
	Manufacturer() Manufacturer
	// LocalName returns the local name of the device.
	LocalName() string
	// Address returns the Bluetooth address of the device.
	Address() Address
	// Services returns the Bluetooth services of the device.
	Services() []Service
	// RSSI returns the received signal strength indicator of the device.
	RSSI() int
	// String returns a string representation of the device.
	String() string
}
