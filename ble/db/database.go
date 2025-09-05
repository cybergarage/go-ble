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

package db

import (
	_ "embed"

	"gopkg.in/yaml.v2"
)

// Assigned Numbers | Bluetooth® Technology Website
// https://www.bluetooth.com/specifications/assigned-numbers/
//
//go:embed data/company_identifiers.yaml
var companyIdentifiers []byte

//go:embed data/service_uuids.yaml
var serviceUUIDs []byte

//go:embed data/sdo_uuids.yaml
var sdoUUIDs []byte

type Database interface {
	// LookupCompany looks up a company by its ID.
	LookupCompany(id int) (Company, bool)
	// LookupService looks up a service by its UUID.
	LookupService(uuid uint16) (Service, bool)
}

var sharedDatabase *database

func init() {
	var companies companies
	err := yaml.Unmarshal(companyIdentifiers, &companies)
	if err != nil {
		panic(err)
	}
	companyMap := make(map[int]*company)
	for _, c := range companies.Companies {
		companyMap[c.Value] = c
	}

	var svcs services
	err = yaml.Unmarshal(serviceUUIDs, &svcs)
	if err != nil {
		panic(err)
	}
	serviceMap := make(map[uint16]*service)
	for _, s := range svcs.Services {
		serviceMap[s.Uuid] = s
	}

	var sdos services
	err = yaml.Unmarshal(sdoUUIDs, &sdos)
	if err != nil {
		panic(err)
	}
	for _, s := range sdos.Services {
		serviceMap[s.Uuid] = s
	}

	sharedDatabase = &database{
		companies: companyMap,
		services:  serviceMap,
	}
}

// DefaultDatabase returns the default database instance.
func DefaultDatabase() Database {
	return sharedDatabase
}

type database struct {
	companies map[int]*company
	services  map[uint16]*service
}

// LookupCompany looks up a company by its ID.
func (db *database) LookupCompany(id int) (Company, bool) {
	dbCompany, ok := db.companies[id]
	if ok {
		return dbCompany, true
	}
	return &company{
		Value: id,
		Nam:   "",
	}, false
}

// LookupService looks up a service by its UUID.
func (db *database) LookupService(uuid uint16) (Service, bool) {
	dbService, ok := db.services[uuid]
	if ok {
		return dbService, true
	}
	return &service{
		Uuid: uuid,
		Nam:  "",
		Id:   "",
	}, false
}
