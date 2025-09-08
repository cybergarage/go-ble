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

package bletest

import (
	"testing"

	"github.com/cybergarage/go-ble/ble/db"
	ble "github.com/cybergarage/go-ble/ble/types"
)

func TestEmbeddedDatabase(t *testing.T) {
	db := db.DefaultDatabase()

	t.Run("BluetoothSIG", func(t *testing.T) {
		t.Run("Company", func(t *testing.T) {
			// Check a few known companies in the embedded database.
			companyTests := []struct {
				ID   int
				Name string
			}{
				{ID: 0x0001, Name: "Nokia Mobile Phones"},
			}
			for _, tt := range companyTests {
				company, ok := db.LookupCompany(tt.ID)
				if !ok {
					t.Errorf("expected company %d to be found", tt.ID)
					continue
				}
				if company.Name() != tt.Name {
					t.Errorf("expected company name to be '%s', got '%s'", tt.Name, company.Name())
				}
			}

			// Check a non-existent company.
			_, ok := db.LookupCompany(0xFFFF)
			if ok {
				t.Errorf("expected company 0xFFFF to not be found")
			}
		})

		t.Run("Service", func(t *testing.T) {
			// Check a few known services in the embedded database.
			serviceTests := []struct {
				UUID ble.UUID
				Name string
			}{
				{UUID: ble.NewUUIDFromUUID16(0x1800), Name: "GAP"},
				{UUID: ble.NewUUIDFromUUID16(0xFFF6), Name: "Matter Profile ID"},
			}
			for _, tt := range serviceTests {
				service, ok := db.LookupService(tt.UUID)
				if !ok {
					t.Errorf("expected service 0x%04X to be found", tt.UUID)
					continue
				}
				if service.Name() != tt.Name {
					t.Errorf("expected service name to be '%s', got '%s'", tt.Name, service.Name())
				}
			}

			// Check a non-existent service.
			_, ok := db.LookupService(ble.NewUUIDFromUUID16(0xFFFF))
			if ok {
				t.Errorf("expected service 0xFFFF to not be found")
			}
		})

		t.Run("Characteristic", func(t *testing.T) {
			// Check a few known characteristics in the embedded database.
			charTests := []struct {
				UUID ble.UUID
				Name string
			}{
				{UUID: ble.NewUUIDFromUUID16(0x2A00), Name: "Device Name"},
				{UUID: ble.NewUUIDFromUUID16(0x2A6E), Name: "Temperature"},
			}
			for _, tt := range charTests {
				characteristic, ok := db.LookupCharacteristic(tt.UUID)
				if !ok {
					t.Errorf("expected characteristic 0x%04X to be found", tt.UUID)
					continue
				}
				if characteristic.Name() != tt.Name {
					t.Errorf("expected characteristic name to be '%s', got '%s'", tt.Name, characteristic.Name())
				}
			}

			// Check a non-existent characteristic.
			_, ok := db.LookupCharacteristic(ble.NewUUIDFromUUID16(0xFFFF))
			if ok {
				t.Errorf("expected characteristic 0xFFFF to not be found")
			}
		})
	})
	t.Run("VendorSpecific", func(t *testing.T) {
		t.Run("Matter", func(t *testing.T) {
			t.Run("Characteristic", func(t *testing.T) {
				charTests := []struct {
					UUID ble.UUID
					Name string
				}{
					{UUID: ble.MustUUIDFromUUIDString("18EE2EF5-263D-4559-959F-4F9C429F9D11"), Name: "C1 Characteristic"},
					{UUID: ble.MustUUIDFromUUIDString("18EE2EF5-263D-4559-959F-4F9C429F9D12"), Name: "C2 Characteristic"},
					{UUID: ble.MustUUIDFromUUIDString("64630238-8772-45F2-B87D-748A83218F04"), Name: "C3 Characteristic"},
				}
				for _, tt := range charTests {
					characteristic, ok := db.LookupCharacteristic(tt.UUID)
					if !ok {
						t.Errorf("expected characteristic 0x%04X to be found", tt.UUID)
						continue
					}
					if characteristic.Name() != tt.Name {
						t.Errorf("expected characteristic name to be '%s', got '%s'", tt.Name, characteristic.Name())
					}
				}
			})
		})
	})
}
