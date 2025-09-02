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
	scanResult   bluetooth.ScanResult
	manufacturer Manufacturer
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) Device {
	return &tinyDevice{
		manufacturer: nil,
		scanResult:   scanResult,
	}
}

// Manufacturer returns the Bluetooth manufacturer of the device.
func (dev *tinyDevice) Manufacturer() Manufacturer {
	if dev.manufacturer == nil {
		manufacturers := dev.scanResult.ManufacturerData()
		switch len(manufacturers) {
		case 0:
			dev.manufacturer = newNilManufacturer()
		case 1:
			manufacturer := manufacturers[0]
			dev.manufacturer = newManufacturer(int(manufacturer.CompanyID), manufacturer.Data)
		default:
			for _, v := range manufacturers {
				dev.manufacturer = newManufacturer(int(v.CompanyID), v.Data)
			}
		}
	}
	if dev.manufacturer.Company().ID() == 2409 {
		return dev.manufacturer
	}
	return dev.manufacturer
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

// RSSI returns the received signal strength indicator of the device.
func (dev *tinyDevice) RSSI() int {
	return int(dev.scanResult.RSSI)
}

func (dev *tinyDevice) ServiceUUIDs() []UUID {
	newUUIDFromBytes := func(b bluetooth.UUID) UUID {
		var out [16]byte
		for i, v := range b {
			out[i*4+0] = byte(v)
			out[i*4+1] = byte(v >> 8)
			out[i*4+2] = byte(v >> 16)
			out[i*4+3] = byte(v >> 24)
		}
		return UUID(out)
	}
	uuids := []UUID{}
	for _, u := range dev.scanResult.ServiceUUIDs() {
		uuids = append(uuids, newUUIDFromBytes(u))
	}
	return uuids
}
