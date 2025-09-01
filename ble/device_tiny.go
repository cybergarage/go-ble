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
	"tinygo.org/x/bluetooth"
)

type tinyDevice struct {
	company    Company
	scanResult bluetooth.ScanResult
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) Device {
	return &tinyDevice{
		company:    nil,
		scanResult: scanResult,
	}
}

// LocalName returns the local name of the device.
func (dev *tinyDevice) LocalName() string {
	return dev.scanResult.LocalName()
}

// Address returns the Bluetooth address of the device.
func (dev *tinyDevice) Address() Address {
	scanAddr := dev.scanResult.Address.Bytes()
	return mustAddressFromBytes(scanAddr[:])
}
