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
	"time"

	"tinygo.org/x/bluetooth"
)

// nolint: staticcheck
type tinyCharacteristic struct {
	*characteristic
	tinyChar *bluetooth.DeviceCharacteristic
	readBuf  []byte
}

func newTinyCharacteristic(service Service, uuid UUID, char *bluetooth.DeviceCharacteristic) *tinyCharacteristic {
	return &tinyCharacteristic{
		characteristic: newCharacteristic(service, uuid),
		tinyChar:       char,
		readBuf:        make([]byte, 512),
	}
}

// MTU returns the ATT maximum transmission unit of the connection.
//
// A Bluetooth stack negotiates the MTU when a connection is established, so the
// value is known only while the device is connected. Sizing the payloads by it
// avoids a write which the peer cannot accept.
func (char *tinyCharacteristic) MTU() (int, error) {
	if char.tinyChar == nil {
		return 0, fmt.Errorf("%w: %s", ErrNotConnected, char.String())
	}
	mtu, err := char.tinyChar.GetMTU()
	if err != nil {
		return 0, fmt.Errorf("%w: %s", err, char.String())
	}
	return int(mtu), nil
}

// Read reads the characteristic value.
func (char *tinyCharacteristic) Read() ([]byte, error) {
	if char.tinyChar == nil {
		return nil, fmt.Errorf("%w: %s", ErrNotConnected, char.String())
	}
	n, err := char.tinyChar.Read(char.readBuf)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, char.String())
	}
	return char.readBuf[:n], nil
}

// WriteWithoutResponse writes the characteristic value without response.
func (char *tinyCharacteristic) WriteWithoutResponse(data []byte) (int, error) {
	if char.tinyChar == nil {
		return 0, fmt.Errorf("%w: %s", ErrNotConnected, char.String())
	}
	nWrote, err := char.tinyChar.WriteWithoutResponse(data)
	if err != nil {
		return nWrote, fmt.Errorf("%w: %s", err, char.String())
	}
	time.Sleep(defaultCharacteristicWriteWithoutResponseWait)
	return nWrote, nil
}

// Notify subscribes to characteristic notifications.
func (char *tinyCharacteristic) Notify(callback OnCharacteristicNotification) error {
	if char.tinyChar == nil {
		return fmt.Errorf("%w: %s", ErrNotConnected, char.String())
	}
	tinyCallback := func(buf []byte) {
		if callback == nil {
			return
		}
		callback(char, buf)
	}
	if err := char.tinyChar.EnableNotifications(tinyCallback); err != nil {
		return fmt.Errorf("%w: %s", err, char.String())
	}
	return nil
}
