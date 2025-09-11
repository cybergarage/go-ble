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

func ExampleService_Open_matter() {
	scanner := ble.NewScanner()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	scanner.Scan(ctx)
	lookupMatterDevices := func(devs []ble.Device, targetDisc uint16) (ble.Device, bool) {
		for _, dev := range devs {
			if service, ok := dev.LookupService(0xFFF6); ok {
				if 3 <= len(service.Data()) {
					serviceDisc := uint16(service.Data()[1]) | (uint16(service.Data()[2]) << 8)
					if serviceDisc == targetDisc {
						return dev, true
					}
				}
			}
		}
		return nil, false
	}
	targetDisc := uint16(4068)
	dev, ok := lookupMatterDevices(scanner.Devices(), targetDisc)
	if ok {
		fmt.Printf("Matter device found: %s", dev.String())
	}
	service, _ := dev.LookupService(0xFFF6)
	c1 := ble.MustUUIDFromUUIDString("18EE2EF5-263D-4559-959F-4F9C429F9D11")
	c2 := ble.MustUUIDFromUUIDString("18EE2EF5-263D-4559-959F-4F9C429F9D12")
	transport, err := service.Open(
		ble.WithTransportReadUUID(c1),
		ble.WithTransportNotifyUUID(c2),
	)
	if err != nil {
		return
	}
	defer transport.Close()
}
