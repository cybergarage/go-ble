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
	"strings"
	"time"

	"github.com/cybergarage/go-ble/ble"
	"github.com/cybergarage/go-logger/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ProgramName = "blelookup"

	FormatParamStr  = "format"
	VerboseParamStr = "verbose"
	DebugParamStr   = "debug"
	TimeoutParamStr = "timeout"
)

var rootCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   ProgramName,
	Short: "Scan and inspect the Bluetooth Low Energy devices nearby",
	Long: `blelookup scans for the Bluetooth Low Energy devices which are advertising
nearby, connects to one of them, and looks up the Bluetooth SIG assigned
numbers.

  blelookup scan                             list the advertising devices
  blelookup scan --service 0xFFF6            list only the devices of a service
  blelookup connect <address> <service>      list the characteristics of a service
  blelookup lookup <uuid|company id>         look up an assigned number

This tool is a central. Advertising as a peripheral is not supported yet.`,
	Version:           ble.Version,
	DisableAutoGenTag: true,
	// The usage is not printed for a runtime error, such as a device which
	// is not found, because the command line itself is valid.
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetDefault(nil)
		verbose := viper.GetBool(VerboseParamStr)
		debug := viper.GetBool(DebugParamStr)
		if debug {
			verbose = true
		}
		if verbose {
			if debug {
				log.SetDefault(log.NewStdoutLogger(log.LevelDebug))
			} else {
				log.SetDefault(log.NewStdoutLogger(log.LevelInfo))
			}
			log.Infof("%s version %s", ProgramName, ble.Version)
			log.Infof("verbose:%t, debug:%t", verbose, debug)
		}
		return nil
	},
}

// RootCommand returns the root command.
func RootCommand() *cobra.Command {
	return rootCmd
}

var sharedCentral ble.Central

// SharedCentral returns the central which the commands run on.
func SharedCentral() ble.Central {
	return sharedCentral
}

// Execute runs the command with the specified central.
func Execute(central ble.Central) error {
	sharedCentral = central
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvPrefix("BLE_LOOKUP")

	viper.SetDefault(FormatParamStr, FormatTableStr)
	rootCmd.PersistentFlags().String(FormatParamStr, FormatTableStr, fmt.Sprintf("output format: %s", strings.Join(allSupportedFormats(), "|")))
	viper.BindPFlag(FormatParamStr, rootCmd.PersistentFlags().Lookup(FormatParamStr))
	viper.BindEnv(FormatParamStr) // BLE_LOOKUP_FORMAT

	viper.SetDefault(TimeoutParamStr, ble.DefaultScanTimeout)
	rootCmd.PersistentFlags().DurationP(TimeoutParamStr, "t", ble.DefaultScanTimeout, "operation timeout")
	viper.BindPFlag(TimeoutParamStr, rootCmd.PersistentFlags().Lookup(TimeoutParamStr))
	viper.BindEnv(TimeoutParamStr) // BLE_LOOKUP_TIMEOUT

	viper.SetDefault(VerboseParamStr, false)
	rootCmd.PersistentFlags().Bool(VerboseParamStr, false, "enable verbose output")
	viper.BindPFlag(VerboseParamStr, rootCmd.PersistentFlags().Lookup(VerboseParamStr))
	viper.BindEnv(VerboseParamStr) // BLE_LOOKUP_VERBOSE

	viper.SetDefault(DebugParamStr, false)
	rootCmd.PersistentFlags().Bool(DebugParamStr, false, "enable debug output")
	viper.BindPFlag(DebugParamStr, rootCmd.PersistentFlags().Lookup(DebugParamStr))
	viper.BindEnv(DebugParamStr) // BLE_LOOKUP_DEBUG
}

// operationTimeout returns the configured timeout.
func operationTimeout() time.Duration {
	timeout := viper.GetDuration(TimeoutParamStr)
	if timeout <= 0 {
		return ble.DefaultScanTimeout
	}
	return timeout
}
