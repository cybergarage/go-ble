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
	"errors"
)

var (
	// ErrNotConnectable indicates that the device is not connectable.
	ErrNotConnected = errors.New("not connected")
	// ErrNotSet indicates that the value is not set.
	ErrNotSet = errors.New("not set")
	// ErrInvalid indicates that the value is invalid.
	ErrInvalid = errors.New("invalid")
	// ErrNotFound indicates that the value was not found.
	ErrNotFound = errors.New("not found")
)
