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
	"slices"
	"strings"
	"time"
)

const (
	// DefaultScanTimeout is the default duration for scanning.
	DefaultScanTimeout = time.Duration(5 * time.Second)
)

// ScannerOption represents an option for the scanner.
//
// A ScanHandler and a ScannerOptionFunc are accepted. The type is left open so
// that the handlers which were passed directly keep working; it will be
// narrowed to ScannerOptionFunc in v1.0.0.
type ScannerOption any

// ScanHandler defines a handler function for scan results.
type ScanHandler func(Device)

// ScannerOptionFunc represents a scanner option which is set by a With
// function.
type ScannerOptionFunc func(*scanOptions)

// scanOptions holds the handlers and the filters of a scan.
type scanOptions struct {
	handlers         []ScanHandler
	serviceUUIDs     []UUID
	addresses        []string
	localNames       []string
	rssiThreshold    int
	hasRSSIThreshold bool
}

// newScanOptions returns the scan options of the specified options.
func newScanOptions(opts ...ScannerOption) *scanOptions {
	scanOpts := &scanOptions{
		handlers:         []ScanHandler{},
		serviceUUIDs:     []UUID{},
		addresses:        []string{},
		localNames:       []string{},
		rssiThreshold:    0,
		hasRSSIThreshold: false,
	}

	for _, opt := range opts {
		switch v := opt.(type) {
		case ScanHandler:
			scanOpts.handlers = append(scanOpts.handlers, v)
		case func(Device):
			scanOpts.handlers = append(scanOpts.handlers, v)
		case ScannerOptionFunc:
			v(scanOpts)
		case func(*scanOptions):
			v(scanOpts)
		}
	}

	return scanOpts
}

// WithScanHandler sets a handler which is called for every matched device.
func WithScanHandler(handler ScanHandler) ScannerOptionFunc {
	return func(opts *scanOptions) {
		opts.handlers = append(opts.handlers, handler)
	}
}

// WithScanServiceUUIDs scans only the devices which advertise one of the
// specified service UUIDs. The UUID can be of any type accepted by NewUUIDFrom,
// such as a string, a uint16 or a UUID.
func WithScanServiceUUIDs(uuids ...any) ScannerOptionFunc {
	return func(opts *scanOptions) {
		for _, anyUUID := range uuids {
			uuid, err := NewUUIDFrom(anyUUID)
			if err != nil {
				continue
			}
			opts.serviceUUIDs = append(opts.serviceUUIDs, uuid)
		}
	}
}

// WithScanAddresses scans only the devices of the specified addresses.
func WithScanAddresses(addrs ...string) ScannerOptionFunc {
	return func(opts *scanOptions) {
		for _, addr := range addrs {
			opts.addresses = append(opts.addresses, strings.ToUpper(addr))
		}
	}
}

// WithScanLocalNames scans only the devices which advertise one of the
// specified local names.
func WithScanLocalNames(names ...string) ScannerOptionFunc {
	return func(opts *scanOptions) {
		opts.localNames = append(opts.localNames, names...)
	}
}

// WithScanRSSIThreshold scans only the devices whose RSSI is equal to or
// greater than the specified threshold, such as -70 to skip the distant
// devices.
func WithScanRSSIThreshold(rssi int) ScannerOptionFunc {
	return func(opts *scanOptions) {
		opts.rssiThreshold = rssi
		opts.hasRSSIThreshold = true
	}
}

// Matches returns true if the specified device passes all the configured
// filters, otherwise false.
func (opts *scanOptions) Matches(dev Device) bool {
	if 0 < len(opts.serviceUUIDs) {
		matched := false
		for _, uuid := range opts.serviceUUIDs {
			if _, ok := dev.LookupService(uuid); ok {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	if 0 < len(opts.addresses) {
		devAddr := strings.ToUpper(dev.Address().String())
		matched := slices.Contains(opts.addresses, devAddr)
		if !matched {
			return false
		}
	}

	if 0 < len(opts.localNames) {
		devName := dev.LocalName()
		matched := slices.Contains(opts.localNames, devName)
		if !matched {
			return false
		}
	}

	if opts.hasRSSIThreshold && dev.RSSI() < opts.rssiThreshold {
		return false
	}

	return true
}

// Scanner defines the interface for a Bluetooth scanner.
type Scanner interface {
	// Devices returns the list of discovered devices.
	Devices() []Device
	// LookupDevice looks up a discovered device by its address.
	LookupDevice(addr string) (Device, bool)
	// Scan starts scanning for Bluetooth devices, and it blocks until the
	// context is done. The scan stops after DefaultScanTimeout when the
	// context has no deadline.
	Scan(ctx context.Context, opts ...ScannerOption) error
}
