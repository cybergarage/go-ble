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
	"strings"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

// tinyScanner scans for the Bluetooth devices with the shared adapter.
//
// The discovered devices are updated from the adapter callback goroutine and
// they are read by the caller, so the device map is guarded by the mutex.
type tinyScanner struct {
	mutex   sync.RWMutex
	devices map[string]*tinyDevice
}

// NewScanner creates a new Bluetooth scanner.
func NewScanner() Scanner {
	return &tinyScanner{
		mutex:   sync.RWMutex{},
		devices: map[string]*tinyDevice{},
	}
}

// Devices returns the list of discovered devices.
func (s *tinyScanner) Devices() []Device {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	devs := make([]Device, 0, len(s.devices))
	for _, dev := range s.devices {
		devs = append(devs, dev)
	}

	return devs
}

// LookupDevice looks up a discovered device by its address.
func (s *tinyScanner) LookupDevice(addr string) (Device, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	dev, ok := s.devices[strings.ToUpper(addr)]
	if !ok {
		return nil, false
	}

	return dev, true
}

// updateDevice adds or updates the device of the specified scan result, and
// returns the device which the scanner holds.
func (s *tinyScanner) updateDevice(scanDev *tinyDevice) *tinyDevice {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	addrKey := strings.ToUpper(scanDev.Address().String())

	discoveredDev, ok := s.devices[addrKey]
	if !ok {
		s.devices[addrKey] = scanDev
		return scanDev
	}

	discoveredDev.lastSeenAt = now
	discoveredDev.rssi = scanDev.RSSI()
	for _, scanService := range scanDev.Services() {
		if _, ok := discoveredDev.LookupService(scanService.UUID()); ok {
			continue
		}
		discoveredDev.addService(scanService)
		discoveredDev.modifiedAt = now
	}

	return discoveredDev
}

// scanOwner holds the token of the shared adapter's scan. Only the goroutine
// which holds the token may run a scan.
var scanOwner = make(chan struct{}, 1)

// Scan starts scanning for Bluetooth devices.
//
// tinygo.org/x/bluetooth supports StopScan() only from inside the scan
// callback: it reads the channel which Scan() writes without synchronizing
// them, so calling it from another goroutine is a data race, and it fails
// outright when it runs before Scan() has started. The context therefore
// requests the stop, and the callback performs it when the next advertisement
// arrives. Scan() itself returns as soon as the context is done, so a caller is
// not held past its deadline, and the scan keeps the adapter until it really
// stops, so that the next scan waits for it instead of failing.
func (s *tinyScanner) Scan(ctx context.Context, opts ...ScannerOption) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultScanTimeout)
		defer cancel()
	}

	scanOpts := newScanOptions(opts...)

	adapter, err := enableDefaultAdapter()
	if err != nil {
		return err
	}

	// The shared adapter runs one scan at a time.
	select {
	case scanOwner <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	scanDone := make(chan error, 1)

	go func() {
		defer func() { <-scanOwner }()

		scanDone <- adapter.Scan(func(adapter *bluetooth.Adapter, scanRes bluetooth.ScanResult) {
			select {
			case <-ctx.Done():
				adapter.StopScan()
				return
			default:
			}

			// The filters are checked before the device is stored, so
			// that Devices() returns only the devices which the scan
			// asked for.
			scanDev := newDeviceFromScanResult(scanRes)
			if !scanOpts.Matches(scanDev) {
				return
			}

			dev := s.updateDevice(scanDev)

			for _, scanHandler := range scanOpts.handlers {
				scanHandler(dev)
			}
		})
	}()

	select {
	case err := <-scanDone:
		return err
	case <-ctx.Done():
		return nil
	}
}
