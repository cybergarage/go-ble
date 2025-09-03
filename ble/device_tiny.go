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
	*baseDevice
	scanResult   bluetooth.ScanResult
	manufacturer Manufacturer
	services     []Service
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) Device {
	return &tinyDevice{
		baseDevice:   newBaseDevice(),
		manufacturer: nil,
		scanResult:   scanResult,
		services:     nil,
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
	return Address(mustUUIDFromBytes(scanAddr[:]))
}

// RSSI returns the received signal strength indicator of the device.
func (dev *tinyDevice) RSSI() int {
	return int(dev.scanResult.RSSI)
}

// LookupService looks up a Bluetooth service by its UUID.
func (dev *tinyDevice) LookupService(uuid UUID) (Service, bool) {
	return dev.lookupServiceFrom(dev.Services(), uuid)
}

// Services returns the Bluetooth services of the device.
func (dev *tinyDevice) Services() []Service {
	if dev.services == nil {
		dev.services = []Service{}
		for _, sd := range dev.scanResult.ServiceData() {
			uuidBytes := sd.UUID.Bytes()
			service := newService(
				mustUUIDFromBytes(uuidBytes[:]),
				"",
				sd.Data,
			)
			dev.services = append(dev.services, service)
		}
	}
	return dev.services
}

// String returns a string representation of the device.
func (dev *tinyDevice) String() string {
	return dev.StringFrom(dev)
}
