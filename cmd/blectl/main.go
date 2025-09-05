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

package main

import (
	"context"
	"time"

	"github.com/cybergarage/go-ble/ble"
	"github.com/cybergarage/go-logger/log"
)

func main() {
	log.EnableStdoutDebug(true)

	central, err := ble.NewCentral()
	if err != nil {
		log.Fatalf("Failed to create central: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = central.Scan(ctx, ble.OnScanResult(func(dev ble.Device) {
		log.Infof("Device found: %s", dev.String())
	}))
	if err != nil {
		log.Fatalf("Failed to scan: %v", err)
	}

	log.Infof("Discovered devices:")
	for n, dev := range central.Devices() {
		log.Infof("[%d] %s", n, dev.String())
	}
}
