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
	_ "embed"

	"gopkg.in/yaml.v2"
)

// Assigned Numbers | Bluetooth® Technology Website
// https://www.bluetooth.com/specifications/assigned-numbers/
//
//go:embed data/company_identifiers.yaml
var companyIdentifiers []byte

type Database interface {
	// LookupCompany looks up a company by its ID.
	LookupCompany(id int) (Company, bool)
}

var sharedDatabase *database

func init() {
	var companies companies
	err := yaml.Unmarshal(companyIdentifiers, &companies)
	if err != nil {
		panic(err)
	}
	sharedDatabase = &database{
		companies: companies.Companies,
	}
}

// DefaultDatabase returns the default database instance.
func DefaultDatabase() Database {
	return sharedDatabase
}

type database struct {
	companies []*company
}

// LookupCompany looks up a company by its ID.
func (db *database) LookupCompany(id int) (Company, bool) {
	for _, company := range db.companies {
		if company.ID() == id {
			return company, true
		}
	}
	return &company{
		Value: id,
		Nam:   "",
	}, false
}
