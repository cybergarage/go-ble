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

	"github.com/cybergarage/go-ble/ble"
)

func TestDatabaseEmbeddedCompany(t *testing.T) {
	tests := []struct {
		ID            int
		Name          string
		expectedFound bool
	}{
		{ID: 0x0001, Name: "Nokia Mobile Phones", expectedFound: true},
	}

	db := ble.DefaultDatabase()

	for _, tt := range tests {
		company, ok := db.LookupCompany(tt.ID)
		if !ok {
			if tt.expectedFound {
				t.Errorf("expected company %d to be found", tt.ID)
			}
			continue
		}
		if company.Name() != tt.Name {
			t.Errorf("expected company name to be '%s', got '%s'", tt.Name, company.Name())
		}
	}
}
