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

package ble_test

import (
	"context"
	"fmt"
	"time"

	"github.com/cybergarage/go-ble/ble"
)

// Scan blocks until the context is done, and it calls the handler for every
// device which passes the filters.
func ExampleCentral_Scan() {
	central := ble.NewCentral()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := central.Scan(ctx,
		ble.WithScanHandler(func(dev ble.Device) {
			fmt.Printf("%s %s (%d dBm)\n", dev.Address(), dev.LocalName(), dev.RSSI())
		}),
	)
	if err != nil {
		return
	}

	fmt.Printf("%d devices\n", len(central.Devices()))
}

// The filters narrow a scan to the devices which are worth looking at.
func ExampleCentral_Scan_filter() {
	central := ble.NewCentral()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	central.Scan(ctx,
		// The Matter service UUID, as an assigned 16-bit UUID.
		ble.WithScanServiceUUIDs(0xFFF6),
		ble.WithScanRSSIThreshold(-70),
		ble.WithScanHandler(func(dev ble.Device) {
			fmt.Println(dev.Address())
		}),
	)
}

// A service is looked up by its UUID on a connected device.
func ExampleCentral_Connect() {
	central := ble.NewCentral()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dev, ok := central.LookupDevice("11:22:33:AA:BB:CC")
	if !ok {
		return
	}

	if err := central.Connect(ctx, dev); err != nil {
		return
	}
	defer dev.Disconnect()

	service, ok := dev.LookupService(0xFFF6)
	if !ok {
		return
	}

	for _, char := range service.Characteristics() {
		mtu, _ := char.MTU()
		fmt.Printf("%s %s (mtu %d)\n", char.UUID(), char.Name(), mtu)
	}
}

// A service which carries a stream over a write and a notify characteristic is
// opened as a transport.
func ExampleService_Open() {
	var service ble.Service

	transport, err := service.Open(
		ble.WithTransportWriteUUID(ble.MustUUIDFromString("18EE2EF5-263D-4559-959F-4F9C429F9D11")),
		ble.WithTransportNotifyUUID(ble.MustUUIDFromString("18EE2EF5-263D-4559-959F-4F9C429F9D12")),
	)
	if err != nil {
		return
	}
	defer transport.Close()

	ctx := context.Background()

	if _, err := transport.Write(ctx, []byte("request")); err != nil {
		return
	}

	// Some peripherals flush their response only once the client enables
	// the notifications after the first write, so subscribing is a separate
	// step.
	if err := transport.Subscribe(); err != nil {
		return
	}

	response, err := transport.Read(ctx)
	if err != nil {
		return
	}

	fmt.Printf("%d bytes\n", len(response))
}

// A connection can be dropped by the peripheral at any time.
func ExampleCentral_SetConnectionHandler() {
	central := ble.NewCentral()

	central.SetConnectionHandler(func(addr ble.Address, connected bool) {
		fmt.Printf("%s connected=%t\n", addr, connected)
	})
}

// The Bluetooth SIG assigned numbers name a service, a characteristic or a
// company.
func ExampleDefaultDatabase() {
	db := ble.DefaultDatabase()

	if service, ok := db.LookupService(ble.MustUUIDFrom(0xFFF6)); ok {
		fmt.Println(service.Name())
	}
	if char, ok := db.LookupCharacteristic(ble.MustUUIDFrom(0x2A00)); ok {
		fmt.Println(char.Name())
	}
	if company, ok := db.LookupCompany(0x004C); ok {
		fmt.Println(company.Name())
	}

	// Output:
	// Matter Profile ID
	// Device Name
	// Apple, Inc.
}
