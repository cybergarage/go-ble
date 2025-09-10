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
	"github.com/cybergarage/go-logger/log"
)

func ExampleScanner() {
	scanner := ble.NewScanner()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := scanner.Scan(ctx,
		ble.OnScanResult(func(dev ble.Device) {
			log.Infof("Device found: dev=%s", dev.String())
		}))
	if err != nil {
		log.Errorf("Failed to scan: %v", err)
	}

	for n, dev := range scanner.Devices() {
		log.Infof("[%d] %s", n, dev.String())
	}
}

func ExampleScanner_matter() {
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
	disc := uint16(4068)
	dev, ok := lookupMatterDevices(scanner.Devices(), disc)
	if ok {
		fmt.Printf("Matter device found: %s", dev.String())
	}
}
