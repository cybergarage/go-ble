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
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-ble/ble"
	"github.com/spf13/cobra"
)

const (
	serviceParamStr = "service"
	addressParamStr = "address"
	nameParamStr    = "name"
	rssiParamStr    = "rssi"
)

var scanCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "scan",
	Short: "Scan for the advertising BLE devices",
	Long: `Scan for the Bluetooth Low Energy devices which are advertising nearby, and
print them as they are discovered.

The scan stops after the timeout, or when it is interrupted.`,
	Example: `  blelookup scan
  blelookup scan --service 0xFFF6 --timeout 30s
  blelookup scan --rssi -70 --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := cmd.Flags().GetStringSlice(serviceParamStr)
		if err != nil {
			return err
		}
		addrs, err := cmd.Flags().GetStringSlice(addressParamStr)
		if err != nil {
			return err
		}
		names, err := cmd.Flags().GetStringSlice(nameParamStr)
		if err != nil {
			return err
		}
		rssi, err := cmd.Flags().GetInt(rssiParamStr)
		if err != nil {
			return err
		}

		writer, err := newDeviceWriter()
		if err != nil {
			return err
		}
		defer writer.Flush()

		opts := []ble.ScannerOption{
			ble.WithScanHandler(func(dev ble.Device) {
				writer.Write(dev)
			}),
		}
		if 0 < len(services) {
			serviceUUIDs := make([]any, 0, len(services))
			for _, service := range services {
				serviceUUIDs = append(serviceUUIDs, service)
			}
			opts = append(opts, ble.WithScanServiceUUIDs(serviceUUIDs...))
		}
		if 0 < len(addrs) {
			opts = append(opts, ble.WithScanAddresses(addrs...))
		}
		if 0 < len(names) {
			opts = append(opts, ble.WithScanLocalNames(names...))
		}
		if cmd.Flags().Changed(rssiParamStr) {
			opts = append(opts, ble.WithScanRSSIThreshold(rssi))
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		ctx, cancel := context.WithTimeout(ctx, operationTimeout())
		defer cancel()

		return SharedCentral().Scan(ctx, opts...)
	},
}

func init() {
	scanCmd.Flags().StringSliceP(serviceParamStr, "s", []string{}, "scan only the devices which advertise the service UUID")
	scanCmd.Flags().StringSlice(addressParamStr, []string{}, "scan only the devices of the address")
	scanCmd.Flags().StringSlice(nameParamStr, []string{}, "scan only the devices of the local name")
	scanCmd.Flags().Int(rssiParamStr, 0, "scan only the devices whose RSSI is equal to or greater than this value")
	rootCmd.AddCommand(scanCmd)
}
