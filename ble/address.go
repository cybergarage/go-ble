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
)

// Address represents a Bluetooth address.
type Address UUID

// String returns the string representation of the Bluetooth address.
func (a Address) String() string {
	b := [16]byte(a)
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		b[10], b[11], b[12], b[13], b[14], b[15])
}
