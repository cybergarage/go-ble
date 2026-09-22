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

/*
Package db holds the Bluetooth SIG assigned numbers, so that a company
identifier, a service UUID or a characteristic UUID can be named.

The tables are the YAML documents which the Bluetooth SIG publishes, and they
are embedded in the package, so a lookup needs no network and no data file.

The API of this package is not stable until v1.0.0.
*/
package db
