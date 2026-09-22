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
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

// tinygo.org/x/bluetooth's UUID stores its four 32-bit words in the reverse
// order of this package's own UUID (id[0] holds the *last* 4 bytes of the
// standard big-endian UUID, not the first). Before tinygo-bluetooth v0.15,
// both types were plain [4]uint32 arrays with this same reversed layout, so a
// direct array conversion transparently preserved it. v0.15 wrapped tinygo's
// array in an opaque struct, so that conversion no longer compiles; these
// helpers reconstruct the same word-for-word (not byte-standard) mapping via
// the BytesBigEndian/NewUUID accessors, so UUID equality against constants
// such as MatterServiceUUID keeps working exactly as before.
func uuidFromTinygo(u bluetooth.UUID) UUID {
	b := u.BytesBigEndian()
	return UUID{
		binary.BigEndian.Uint32(b[12:16]),
		binary.BigEndian.Uint32(b[8:12]),
		binary.BigEndian.Uint32(b[4:8]),
		binary.BigEndian.Uint32(b[0:4]),
	}
}

func uuidToTinygo(u UUID) bluetooth.UUID {
	var b [16]byte
	binary.BigEndian.PutUint32(b[0:4], u[3])
	binary.BigEndian.PutUint32(b[4:8], u[2])
	binary.BigEndian.PutUint32(b[8:12], u[1])
	binary.BigEndian.PutUint32(b[12:16], u[0])
	return bluetooth.NewUUID(b)
}

type tinyDevice struct {
	*baseDevice
	scanResult    bluetooth.ScanResult
	manufacturers []Manufacturer
	rssi          int
	adServiceMap  sync.Map
	tinyDev       *bluetooth.Device
}

func newDeviceFromScanResult(scanResult bluetooth.ScanResult) *tinyDevice {
	// The manufacturers are built here, not on demand, because a scanned
	// device is read from the caller goroutine while the adapter callback
	// goroutine is updating it.
	manufacturerData := scanResult.ManufacturerData()
	manufacturers := make([]Manufacturer, 0, len(manufacturerData))
	for _, md := range manufacturerData {
		manufacturers = append(manufacturers, newManufacturer(int(md.CompanyID), md.Data))
	}

	dev := &tinyDevice{
		baseDevice:    newBaseDevice(),
		manufacturers: manufacturers,
		scanResult:    scanResult,
		rssi:          int(scanResult.RSSI),
		adServiceMap:  sync.Map{},
		tinyDev:       nil,
	}
	for _, sd := range scanResult.ServiceData() {
		dev.addServiceDataElement(sd)
	}
	return dev
}

// Manufacturer returns the first Bluetooth manufacturer of the device.
//
// An advertisement may hold more than one manufacturer specific data element,
// and only the last one was kept before. Use Manufacturers() to read them all.
func (dev *tinyDevice) Manufacturer() Manufacturer {
	if len(dev.manufacturers) == 0 {
		return newNilManufacturer()
	}
	return dev.manufacturers[0]
}

// Manufacturers returns all the Bluetooth manufacturers of the device.
func (dev *tinyDevice) Manufacturers() []Manufacturer {
	return dev.manufacturers
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
func (dev *tinyDevice) LookupService(anyUUID any) (Service, bool) {
	lookupUUID, err := NewUUIDFrom(anyUUID)
	if err != nil {
		return nil, false
	}

	// If not connected, look up in the cached services.
	if !dev.IsConnected() {
		return dev.lookupAdvertisedService(lookupUUID)
	}

	// If connected, discover services from the device using the Bluetooth API.
	tinyServices, err := dev.tinyDev.DiscoverServices([]bluetooth.UUID{uuidToTinygo(lookupUUID)})
	if err != nil {
		return nil, false
	}
	for _, tinyService := range tinyServices {
		tinyServiceUUID := uuidFromTinygo(tinyService.UUID())
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
				uuid, err := NewUUIDFromString(tinyChar.UUID().String())
				if err != nil {
					continue
				}
				char := newTinyCharacteristic(
					service,
					uuid,
					&tinyChar,
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
		uuidFromTinygo(sd.UUID),
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
	// The adapter is enabled here as well, because a device may be built
	// from an address instead of a scan result, and connecting is then the
	// first operation of the process.
	adapter, err := enableDefaultAdapter()
	if err != nil {
		return err
	}
	tinyAddr, err := addressToTiny(dev.Address())
	if err != nil {
		return err
	}
	connParams := bluetooth.ConnectionParams{} // nolint: exhaustruct,exhaustruct_v5
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
	devServices := dev.Services()
	serviceObjs := make([]any, 0, len(devServices))
	for _, service := range devServices {
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
