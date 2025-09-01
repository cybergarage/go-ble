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

var adapter = bluetooth.DefaultAdapter

type tinyScanner struct {
}

func NewScanner() Scanner {
	return &tinyScanner{}
}

// Scan starts scanning for Bluetooth devices.
func (s *tinyScanner) Scan(ctx context.Context, onResult OnScanResult) error {
	err := adapter.Scan(func(adapter *bluetooth.Adapter, scanRes bluetooth.ScanResult) {
		onResult(newDeviceFromScanResult(scanRes))
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				adapter.StopScan()
				return
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	return nil
}
