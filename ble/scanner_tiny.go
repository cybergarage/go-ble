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
	"time"

	"tinygo.org/x/bluetooth"
)

type tinyScanner struct {
	devices map[string]*tinyDevice
}

// NewScanner creates a new Bluetooth scanner.
func NewScanner() Scanner {
	return &tinyScanner{
		devices: map[string]*tinyDevice{},
	}
}

// Devices returns the list of discovered devices.
func (s *tinyScanner) Devices() []Device {
	devs := make([]Device, 0, len(s.devices))
	for _, dev := range s.devices {
		devs = append(devs, dev)
	}
	return devs
}

// Scan starts scanning for Bluetooth devices.
func (s *tinyScanner) Scan(ctx context.Context, opts ...ScannerOption) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultScanTimeout)
		defer cancel()
	}

	scanHandlers := []ScanHandler{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case ScanHandler:
			scanHandlers = append(scanHandlers, v)
		}
	}
	err := defaultAdapter().Enable()
	if err != nil {
		return err
	}
	err = defaultAdapter().Scan(func(adapter *bluetooth.Adapter, scanRes bluetooth.ScanResult) {
		select {
		case <-ctx.Done():
			adapter.StopScan()
			return
		default:
			now := time.Now()
			addrKey := scanRes.Address.String()
			scanDev := newDeviceFromScanResult(scanRes)
			discoveredDev, ok := s.devices[addrKey]
			if ok {
				discoveredDev.lastSeenAt = now
				discoveredDev.rssi = scanDev.RSSI()
				for _, scanService := range scanDev.Services() {
					if _, ok := discoveredDev.LookupService(scanService.UUID()); !ok {
						discoveredDev.addService(scanService)
						discoveredDev.modifiedAt = now
					}
				}
			} else {
				s.devices[addrKey] = scanDev
				discoveredDev = scanDev
			}

			for _, scanHandler := range scanHandlers {
				scanHandler(discoveredDev)
			}
		}
	})
	return err
}
