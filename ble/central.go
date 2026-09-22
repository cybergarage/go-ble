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

// ConnectionHandler is a handler function which is called when a device is
// connected or disconnected.
type ConnectionHandler func(addr Address, connected bool)

// Central represents a Bluetooth central device.
type Central interface {
	Scanner
	// Connect connects to the specified device.
	Connect(ctx context.Context, dev Device) error
	// SetConnectionHandler sets a handler which is called when a device is
	// connected or disconnected.
	//
	// A peripheral can drop a connection at any time, and a central which
	// does not notice it keeps writing to a connection which is gone. The
	// handler is set on the shared adapter, so it reports every device of
	// the process.
	SetConnectionHandler(handler ConnectionHandler)
}
