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

package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cybergarage/go-ble/ble"
	"github.com/spf13/cobra"
)

var lookupCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "lookup <uuid|company id>",
	Short: "Look up a Bluetooth SIG assigned number",
	Long: `Look up a Bluetooth SIG assigned number in the bundled database, and print the
service, the characteristic or the company which it names.

The command needs no Bluetooth hardware, so it can be used to read a UUID which
was captured elsewhere.`,
	Example: `  blelookup lookup 0xFFF6
  blelookup lookup 0000fff6-0000-1000-8000-00805f9b34fb
  blelookup lookup --format json 0x004C`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]

		objs := []assignedNumberObject{}

		db := ble.DefaultDatabase()

		if uuid, err := ble.NewUUIDFrom(arg); err == nil {
			if service, ok := db.LookupService(uuid); ok {
				objs = append(objs, assignedNumberObject{
					Kind:  "service",
					Value: uuid.String(),
					Name:  service.Name(),
					ID:    service.ID(),
				})
			}
			if char, ok := db.LookupCharacteristic(uuid); ok {
				objs = append(objs, assignedNumberObject{
					Kind:  "characteristic",
					Value: uuid.String(),
					Name:  char.Name(),
					ID:    char.ID(),
				})
			}
		}

		if id, err := parseCompanyID(arg); err == nil {
			if company, ok := db.LookupCompany(id); ok {
				objs = append(objs, assignedNumberObject{
					Kind:  "company",
					Value: fmt.Sprintf("0x%04X", id),
					Name:  company.Name(),
					ID:    strconv.Itoa(company.ID()),
				})
			}
		}

		if len(objs) == 0 {
			return fmt.Errorf("%s is %w", arg, ble.ErrNotFound)
		}

		return outputAssignedNumbers(objs)
	},
}

// parseCompanyID parses a company identifier, which is a 16-bit number.
func parseCompanyID(s string) (int, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		id, err := strconv.ParseUint(s[2:], 16, 16)
		return int(id), err
	}
	id, err := strconv.ParseUint(s, 10, 16)
	return int(id), err
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
