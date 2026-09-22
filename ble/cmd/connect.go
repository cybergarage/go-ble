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
	"context"
	"fmt"

	"github.com/cybergarage/go-ble/ble"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "connect <address> <service>",
	Short: "Connect to a device and list the characteristics of a service",
	Long: `Scan for the device of the specified address, connect to it, and print the
characteristics of the specified service with the negotiated ATT MTU.

The device must be advertising, because it is found by a scan first.`,
	Example: `  blelookup connect 11:22:33:AA:BB:CC 0xFFF6
  blelookup connect --format json 0102030A-0B0C-0D0E-0F10-111213141516 0xFFF6`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, serviceUUID := args[0], args[1]

		central := SharedCentral()

		ctx, cancel := context.WithTimeout(cmd.Context(), operationTimeout())
		defer cancel()

		// The device is found by a scan, because a connection is made to
		// an address which the adapter has seen advertising.
		scanCtx, stopScan := context.WithCancel(ctx)
		defer stopScan()

		err := central.Scan(scanCtx,
			ble.WithScanAddresses(addr),
			ble.WithScanHandler(func(dev ble.Device) {
				stopScan()
			}),
		)
		if err != nil {
			return err
		}

		dev, ok := central.LookupDevice(addr)
		if !ok {
			return fmt.Errorf("device %s is %w", addr, ble.ErrNotFound)
		}

		if err := central.Connect(ctx, dev); err != nil {
			return err
		}
		defer dev.Disconnect()

		service, ok := dev.LookupService(serviceUUID)
		if !ok {
			return fmt.Errorf("service %s of %s is %w", serviceUUID, addr, ble.ErrNotFound)
		}

		return outputCharacteristics(service)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
