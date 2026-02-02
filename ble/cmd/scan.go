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
	"time"

	"github.com/cybergarage/go-ble/ble"
	"github.com/cybergarage/go-logger/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "scan",
	Short: "Scan for Matter devices.",
	Long:  "Scan for Matter devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// format, err := NewFormatFromString(viper.GetString(FormatParamStr))
		// if err != nil {
		// 	return err
		// }

		central := SharedCentral()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := central.Scan(ctx, ble.ScanHandler(func(dev ble.Device) {
			log.Infof("Device responded: %s", dev.String())
		}))
		if err != nil {
			log.Fatalf("Failed to scan: %v", err)
		}

		log.Infof("Discovered devices:")
		for n, dev := range central.Devices() {
			log.Infof("[%d] %s", n, dev.String())
		}
		return nil
	},
}
