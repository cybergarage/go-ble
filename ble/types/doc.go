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
Package types holds the Bluetooth UUID type, which the rest of go-ble is built
on.

A UUID holds its four 32-bit words in the big-endian order, so the first word
holds the leading bytes of the UUID, as UUID.Bytes() and UUID.String() write
them. An assigned 16-bit or 32-bit UUID is expanded with the Bluetooth Base
UUID, so it compares equal to its full form.

The API of this package is not stable until v1.0.0.
*/
package types
