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

	"tinygo.org/x/bluetooth"
)

type tinyCentral struct {
	Scanner
}

// NewCentral creates a new Bluetooth central device.
func NewCentral() Central {
	return &tinyCentral{
		Scanner: NewScanner(),
	}
}

// Connect connects to the specified device.
func (c *tinyCentral) Connect(ctx context.Context, dev Device) error {
	return dev.Connect(ctx)
}

// SetConnectionHandler sets a handler which is called when a device is
// connected or disconnected.
func (c *tinyCentral) SetConnectionHandler(handler ConnectionHandler) {
	if handler == nil {
		defaultAdapter().SetConnectHandler(nil)
		return
	}

	// The handler is set on the adapter itself, so it does not require the
	// adapter to be enabled, and a machine without Bluetooth hardware does
	// not fail here.
	defaultAdapter().SetConnectHandler(func(tinyDev bluetooth.Device, connected bool) {
		addr, err := newAddressFromTiny(tinyDev.Address)
		if err != nil {
			return
		}
		handler(addr, connected)
	})
}
