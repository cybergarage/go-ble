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

/*
Package ble implements a Bluetooth Low Energy central.

# Status

The package is a central (client). It scans for the devices which are
advertising nearby, connects to one of them, and reads, writes and subscribes to
its characteristics. The peripheral (server) role is not implemented: go-ble
cannot advertise a service or serve a GATT database.

# Scanning

[Central.Scan] blocks until the context is done, and it calls the handler for
every device which passes the filters.

	central := ble.NewCentral()

	err := central.Scan(ctx,
		ble.WithScanServiceUUIDs(0xFFF6),
		ble.WithScanHandler(func(dev ble.Device) {
			log.Println(dev.Address(), dev.LocalName(), dev.RSSI())
		}),
	)

# Connecting

A service is looked up by its UUID on a connected device, and its
characteristics are read, written and subscribed to.

	if err := central.Connect(ctx, dev); err != nil {
		return err
	}
	defer dev.Disconnect()

	service, ok := dev.LookupService(0xFFF6)

[Service.Open] wraps a read, a write and a notify characteristic as a
[Transport], which carries bytes instead of characteristics.

# UUIDs and the assigned numbers

A [UUID] holds its words in the big-endian order, so an assigned 16-bit UUID and
its expanded form are the same value:

	ble.NewUUIDFromUUID16(0xFFF6) == ble.MustUUIDFromString("0000FFF6-0000-1000-8000-00805F9B34FB")

[DefaultDatabase] looks up the Bluetooth SIG assigned numbers, so a service, a
characteristic or a company identifier can be named without a table of its own.

# Platforms

The package is built on tinygo.org/x/bluetooth, which uses CoreBluetooth on
macOS, BlueZ on Linux and WinRT on Windows. A device is identified by a MAC
address on Linux and Windows, and by a UUID which the system assigns on macOS,
so an address is not portable between the platforms.
*/
package ble
