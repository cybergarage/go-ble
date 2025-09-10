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
