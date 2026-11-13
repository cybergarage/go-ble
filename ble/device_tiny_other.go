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

//go:build !linux

package ble

import "tinygo.org/x/bluetooth"

// cancelPendingConnection does nothing on these platforms: their tinygo
// implementations end a connection attempt by themselves, after the
// ConnectionTimeout of the connection parameters.
func cancelPendingConnection(bluetooth.Address) {}
