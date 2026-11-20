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
	"encoding/binary"
	"runtime"
	"testing"

	"tinygo.org/x/bluetooth"
)

type scanTestPayload struct {
	bluetooth.AdvertisementPayload
	services []bluetooth.ServiceDataElement
}

func (p scanTestPayload) LocalName() string { return "Sensor" }

func (p scanTestPayload) ManufacturerData() []bluetooth.ManufacturerDataElement { return nil }

func (p scanTestPayload) ServiceData() []bluetooth.ServiceDataElement { return p.services }

func cacheTestScanResult(t *testing.T, addressID uint32, services ...bluetooth.ServiceDataElement) *tinyDevice {
	t.Helper()
	var addressLength = macAddressLen
	if runtime.GOOS == "darwin" {
		addressLength = uuidAddressLen
	}
	address := make(Address, addressLength)
	binary.LittleEndian.PutUint32(address, addressID)
	tinyAddress, err := addressToTiny(address)
	if err != nil {
		t.Fatal(err)
	}
	return newDeviceFromScanResult(bluetooth.ScanResult{
		Address:              tinyAddress,
		AdvertisementPayload: scanTestPayload{services: services},
	})
}

func TestScannerDeviceCacheCapacity(t *testing.T) {
	scanner := NewScanner().(*tinyScanner)
	first := scanner.updateDevice(cacheTestScanResult(t, 1))
	second := scanner.updateDevice(cacheTestScanResult(t, 2))
	for id := 3; id <= maxCachedDevices; id++ {
		scanner.updateDevice(cacheTestScanResult(t, uint32(id)))
	}
	if got, ok := scanner.LookupDevice(second.Address().String()); !ok || got != second {
		t.Fatal("device at capacity was not available for lookup")
	}
	if got := scanner.updateDevice(cacheTestScanResult(t, 1)); got != first {
		t.Fatal("a repeated discovery replaced the cached device")
	}
	newest := scanner.updateDevice(cacheTestScanResult(t, maxCachedDevices+1))
	if got := len(scanner.Devices()); got != maxCachedDevices {
		t.Fatalf("cached %d devices, want %d", got, maxCachedDevices)
	}
	if _, ok := scanner.LookupDevice(second.Address().String()); ok {
		t.Fatal("least recently seen device was not evicted")
	}
	for _, retained := range []*tinyDevice{first, newest} {
		if got, ok := scanner.LookupDevice(retained.Address().String()); !ok || got != retained {
			t.Fatal("recently seen device was not available for lookup")
		}
	}
}

func TestScannerAdvertisedServiceCacheCapacity(t *testing.T) {
	scanner := NewScanner().(*tinyScanner)
	first := bluetooth.ServiceDataElement{UUID: uuidToTinygo(NewUUIDFromUUID16(1)), Data: []byte{1}}
	device := scanner.updateDevice(cacheTestScanResult(t, 1, first))
	for id := 2; id <= maxCachedServicesPerDevice+1; id++ {
		next := bluetooth.ServiceDataElement{UUID: uuidToTinygo(NewUUIDFromUUID16(uint16(id))), Data: []byte{2}}
		scanner.updateDevice(cacheTestScanResult(t, 1, next))
	}
	if got := len(device.Services()); got != maxCachedServicesPerDevice {
		t.Fatalf("cached %d services, want %d", got, maxCachedServicesPerDevice)
	}
	service, ok := device.LookupService("0002")
	if !ok || service.Device() != device || len(service.Data()) != 1 || service.Data()[0] != 2 {
		t.Fatal("merged service was not available with its original data and cached device")
	}
	if _, ok := device.LookupService(uint16(maxCachedServicesPerDevice + 1)); ok {
		t.Fatal("service beyond capacity was cached")
	}

	// The constructor must also enforce the limit before the first result is
	// stored, independently of the limit on merging subsequent results.
	services := make([]bluetooth.ServiceDataElement, maxCachedServicesPerDevice+1)
	for i := range services {
		services[i].UUID = uuidToTinygo(NewUUIDFromUUID16(uint16(i + 1)))
	}
	if got := len(cacheTestScanResult(t, 2, services...).Services()); got != maxCachedServicesPerDevice {
		t.Fatalf("first result retained %d services, want %d", got, maxCachedServicesPerDevice)
	}
}
