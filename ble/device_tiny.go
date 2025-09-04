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
	"encoding/json"
	"time"

	"tinygo.org/x/bluetooth"
)

type tinyDevice struct {
	*baseDevice
	scanResult   bluetooth.ScanResult
	manufacturer Manufacturer
	rssi         int
	serviceMap   map[UUID]Service
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) *tinyDevice {
	return &tinyDevice{
		baseDevice:   newBaseDevice(),
		manufacturer: nil,
		scanResult:   scanResult,
		rssi:         int(scanResult.RSSI),
		serviceMap:   nil,
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
	return Address(dev.scanResult.Address.String())
}

// RSSI returns the received signal strength indicator of the device.
func (dev *tinyDevice) RSSI() int {
	return dev.rssi
}

// LookupService looks up a Bluetooth service by its UUID.
func (dev *tinyDevice) LookupService(uuid UUID) (Service, bool) {
	service, ok := dev.serviceMap[uuid]
	return service, ok
}

// Services returns the Bluetooth services of the device.
func (dev *tinyDevice) Services() []Service {
	if dev.serviceMap == nil {
		dev.serviceMap = make(map[UUID]Service)
		for _, sd := range dev.scanResult.ServiceData() {
			uuidBytes := sd.UUID.Bytes()
			service := newService(
				mustUUIDFromBytes(uuidBytes[:]),
				"",
				sd.Data,
			)
			dev.serviceMap[service.UUID()] = service
		}
	}
	services := make([]Service, 0, len(dev.serviceMap))
	for _, service := range dev.serviceMap {
		services = append(services, service)
	}
	return services
}

func (dev *tinyDevice) MarshalObject() any {
	serviceObjs := []any{}
	for _, service := range dev.Services() {
		serviceObjs = append(serviceObjs, service.MarshalObject())
	}
	return struct {
		Address      Address      `json:"address"`
		LocalName    string       `json:"localName"`
		Manufacturer Manufacturer `json:"manufacturer"`
		RSSI         int          `json:"rssi"`
		Services     []any        `json:"services"`
		DiscoveredAt string       `json:"discoveredAt"`
		ModifiedAt   string       `json:"modifiedAt"`
		LastSeenAt   string       `json:"lastSeenAt"`
	}{
		Address:      dev.Address(),
		LocalName:    dev.LocalName(),
		Manufacturer: dev.Manufacturer(),
		RSSI:         dev.RSSI(),
		Services:     serviceObjs,
		DiscoveredAt: dev.discoveredAt.Format(time.RFC3339),
		ModifiedAt:   dev.modifiedAt.Format(time.RFC3339),
		LastSeenAt:   dev.lastSeenAt.Format(time.RFC3339),
	}
}

// String returns a string representation of the device.
func (dev *tinyDevice) String() string {
	b, err := json.Marshal(dev.MarshalObject())
	if err != nil {
		return "{}"
	}
	return string(b)
}
