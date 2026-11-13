// Copyright (C) 2026 The go-ble Authors. All rights reserved.
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

//go:build linux

package ble

import (
	"strings"

	"github.com/godbus/dbus/v5"
	"tinygo.org/x/bluetooth"
)

// bluezAdapterPath is the D-Bus object of the adapter tinygo's
// DefaultAdapter uses on Linux.
const bluezAdapterPath = "/org/bluez/hci0"

// cancelPendingConnection asks BlueZ to drop a connection which is still
// being made. The Linux implementation of tinygo waits for BlueZ to report
// the connection, without a timeout. Calling Disconnect on the device
// shuts the pending LE connection down, and BlueZ answers the pending
// Connect call with an error (src/device.c, device_request_disconnect), so
// the waiting call returns as well.
func cancelPendingConnection(addr bluetooth.Address) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return
	}
	path := dbus.ObjectPath(bluezAdapterPath + "/dev_" + strings.ReplaceAll(addr.MAC.String(), ":", "_"))
	_ = conn.Object("org.bluez", path).Call("org.bluez.Device1.Disconnect", 0).Err
}
