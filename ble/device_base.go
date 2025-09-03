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
	"fmt"
	"strings"
	"time"
)

type baseDevice struct {
	discoveredAt time.Time
	modifiedAt   time.Time
	lastSeenAt   time.Time
}

func newBaseDevice() *baseDevice {
	now := time.Now()
	return &baseDevice{
		discoveredAt: now,
		modifiedAt:   now,
		lastSeenAt:   now,
	}
}

// DiscoveredAt returns the time when the device was discovered.
func (baseDev *baseDevice) DiscoveredAt() time.Time {
	return baseDev.discoveredAt
}

// ModifiedAt returns the time when the device was last modified.
func (baseDev *baseDevice) ModifiedAt() time.Time {
	return baseDev.modifiedAt
}

// LastSeenAt returns the time when the device was last seen.
func (baseDev *baseDevice) LastSeenAt() time.Time {
	return baseDev.lastSeenAt
}

// LookupService looks up a Bluetooth service by its UUID.
func (baseDev *baseDevice) lookupServiceFrom(services []Service, uuid UUID) (Service, bool) {
	for _, s := range services {
		if s.UUID().Equal(uuid) {
			return s, true
		}
	}
	return nil, false
}

// String returns a string representation of the device.
func (baseDev *baseDevice) StringFrom(dev Device) string {
	services := []string{}
	for _, s := range dev.Services() {
		services = append(services, s.String())
	}
	return fmt.Sprintf("Device[Address: %s, LocalName: %s, Manufacturer: %s, RSSI: %d, Services: (%s)]",
		dev.Address(), dev.LocalName(), dev.Manufacturer().String(), dev.RSSI(), strings.Join(services, ", "))
}
