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

package ext

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed matter/characteristic_uuids.yaml
var matterCharacteristicUUIDs []byte

// Database represents a Bluetooth database.
type Database interface {
	// LookupCharacteristic looks up a characteristic by its UUID.
	LookupCharacteristic(uuid UUID) (Characteristic, bool)
}

var sharedDatabase *database

func init() {
	// Characteristic UUIDs

	var chars characteristics
	err := yaml.Unmarshal(matterCharacteristicUUIDs, &chars)
	if err != nil {
		panic(err)
	}
	characteristicMap := make(map[UUID]*characteristic)
	for _, c := range chars.Characteristics {
		var err error
		c.uuid, err = NewUUIDFromUUIDString(c.Uuid)
		if err != nil {
			panic(err)
		}
		characteristicMap[c.uuid] = c
	}

	sharedDatabase = &database{
		chars: characteristicMap,
	}
}

// DefaultDatabase returns the default database instance.
func DefaultDatabase() Database {
	return sharedDatabase
}

type database struct {
	chars map[UUID]*characteristic
}

// LookupCharacteristic looks up a characteristic by its UUID.
func (db *database) LookupCharacteristic(uuid UUID) (Characteristic, bool) {
	dbChar, ok := db.chars[uuid]
	if ok {
		return dbChar, true
	}
	return &characteristic{
		uuid: uuid,
		Uuid: "",
		Nam:  "",
		Id:   "",
	}, false
}
