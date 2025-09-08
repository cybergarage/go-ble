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
	"encoding/json"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

type tinyDevice struct {
	*baseDevice
	scanResult   bluetooth.ScanResult
	manufacturer Manufacturer
	rssi         int
	adServiceMap sync.Map
	tinyDev      *bluetooth.Device
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) *tinyDevice {
	dev := &tinyDevice{
		baseDevice:   newBaseDevice(),
		manufacturer: nil,
		scanResult:   scanResult,
		rssi:         int(scanResult.RSSI),
		adServiceMap: sync.Map{},
		tinyDev:      nil,
	}
	for _, sd := range scanResult.ServiceData() {
		dev.addServiceDataElement(sd)
	}
	return dev
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
	return dev.manufacturer
}

// LocalName returns the local name of the device.
func (dev *tinyDevice) LocalName() string {
	return dev.scanResult.LocalName()
}

// Address returns the Bluetooth address of the device.
func (dev *tinyDevice) Address() Address {
	addr, _ := newAddressFromTiny(dev.scanResult.Address)
	return addr
}

// RSSI returns the received signal strength indicator of the device.
func (dev *tinyDevice) RSSI() int {
	return dev.rssi
}

func (dev *tinyDevice) lookupAdvertisedService(lookupUUID UUID) (Service, bool) {
	for _, service := range dev.Services() {
		if lookupUUID.Equal(service.UUID()) {
			return service, true
		}
	}
	return nil, false
}

// LookupService looks up a Bluetooth service by its UUID.
func (dev *tinyDevice) LookupService(lookupUUID UUID) (Service, bool) {
	// If not connected, look up in the cached services.
	if !dev.IsConnected() {
		return dev.lookupAdvertisedService(lookupUUID)
	}

	// If connected, discover services from the device using the Bluetooth API.
	tinyServices, err := dev.tinyDev.DiscoverServices([]bluetooth.UUID{bluetooth.UUID(lookupUUID)})
	if err != nil {
		return nil, false
	}
	for _, tinyService := range tinyServices {
		tinyServiceUUID := UUID(tinyService.UUID())
		if lookupUUID.Equal(tinyServiceUUID) {
			tinyChars, err := tinyService.DiscoverCharacteristics(nil)
			if err != nil {
				return nil, false
			}
			adData := []byte{}
			adService, ok := dev.lookupAdvertisedService(lookupUUID)
			if ok {
				adData = adService.Data()
			}
			service := newTinyService(
				dev,
				&tinyService,
				tinyServiceUUID,
				adData,
				[]Characteristic{},
			)
			for _, tinyChar := range tinyChars {
				char := newCharacteristic(
					service,
					UUID(tinyChar.UUID()),
				)
				service.addDeviceCharacteristic(char)
			}
			return service, true
		}
	}
	return nil, false
}

func (dev *tinyDevice) addServiceDataElement(sd bluetooth.ServiceDataElement) {
	service := newService(
		dev,
		UUID(sd.UUID),
		sd.Data,
		[]Characteristic{}, // No characteristics in scan result
	)
	dev.addService(service)
}

func (dev *tinyDevice) addService(service Service) {
	dev.adServiceMap.Store(service.UUID(), service)
}

// Services returns the Bluetooth services of the device.
func (dev *tinyDevice) Services() []Service {
	services := make([]Service, 0)
	dev.adServiceMap.Range(func(key, value any) bool {
		service, ok := value.(Service)
		if ok {
			services = append(services, service)
		}
		return true
	})
	return services
}

// Connect connects to the device.
func (dev *tinyDevice) Connect(ctx context.Context) error {
	adapter := defaultAdapter()
	tinyAddr, err := addressToTiny(dev.Address())
	if err != nil {
		return err
	}
	connParams := bluetooth.ConnectionParams{} // nolint: exhaustruct
	tinyDev, err := adapter.Connect(tinyAddr, connParams)
	if err != nil {
		return err
	}
	dev.tinyDev = &tinyDev
	return nil
}

// Disconnect disconnects from the device.
func (dev *tinyDevice) Disconnect() error {
	if dev.tinyDev == nil {
		return nil
	}
	err := dev.tinyDev.Disconnect()
	if err != nil {
		return err
	}
	dev.tinyDev = nil
	return nil
}

// IsConnected returns whether the device is connected.
func (dev *tinyDevice) IsConnected() bool {
	return dev.tinyDev != nil
}

// MarshalObject returns an object suitable for marshaling to JSON.
func (dev *tinyDevice) MarshalObject() any {
	serviceObjs := []any{}
	for _, service := range dev.Services() {
		serviceObjs = append(serviceObjs, service.MarshalObject())
	}
	return struct {
		Address      string `json:"address"`
		LocalName    string `json:"localName"`
		Manufacturer any    `json:"manufacturer"`
		RSSI         int    `json:"rssi"`
		Services     []any  `json:"services"`
		DiscoveredAt string `json:"discoveredAt"`
		ModifiedAt   string `json:"modifiedAt"`
		LastSeenAt   string `json:"lastSeenAt"`
	}{
		Address:      dev.Address().String(),
		LocalName:    dev.LocalName(),
		Manufacturer: dev.Manufacturer().MarshalObject(),
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
