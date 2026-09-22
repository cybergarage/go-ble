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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/cybergarage/go-ble/ble"
	"github.com/spf13/viper"
)

// The column headers which more than one table shares.
const (
	uuidColumnStr = "UUID"
	nameColumnStr = "NAME"
)

// outputFormat returns the configured output format.
func outputFormat() (Format, error) {
	return NewFormatFromString(viper.GetString(FormatParamStr))
}

// deviceObject represents a device for the JSON and the CSV output.
type deviceObject struct {
	Address      string   `json:"address"`
	LocalName    string   `json:"localName"`
	Manufacturer string   `json:"manufacturer"`
	RSSI         int      `json:"rssi"`
	Services     []string `json:"services"`
}

func newDeviceObject(dev ble.Device) deviceObject {
	devServices := dev.Services()
	services := make([]string, 0, len(devServices))
	for _, service := range devServices {
		name := service.Name()
		if len(name) == 0 {
			name = service.UUID().String()
		}
		services = append(services, name)
	}

	manufacturer := dev.Manufacturer().Name()
	if len(manufacturer) == 0 {
		manufacturer = strconv.Itoa(dev.Manufacturer().ID())
	}

	return deviceObject{
		Address:      dev.Address().String(),
		LocalName:    dev.LocalName(),
		Manufacturer: manufacturer,
		RSSI:         dev.RSSI(),
		Services:     services,
	}
}

func (obj deviceObject) record() []string {
	return []string{
		obj.Address,
		obj.LocalName,
		obj.Manufacturer,
		strconv.Itoa(obj.RSSI),
		strings.Join(obj.Services, " "),
	}
}

func deviceHeader() []string {
	return []string{"ADDRESS", nameColumnStr, "MANUFACTURER", "RSSI", "SERVICES"}
}

// deviceWriter writes the devices in the configured output format as they are
// discovered.
type deviceWriter struct {
	mutex       sync.Mutex
	format      Format
	wroteHeader bool
	tabWriter   *tabwriter.Writer
	csvWriter   *csv.Writer
}

func newDeviceWriter() (*deviceWriter, error) {
	format, err := outputFormat()
	if err != nil {
		return nil, err
	}

	writer := &deviceWriter{
		mutex:       sync.Mutex{},
		format:      format,
		wroteHeader: false,
		tabWriter:   nil,
		csvWriter:   nil,
	}

	switch format {
	case FormatTable:
		writer.tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	case FormatCSV:
		writer.csvWriter = csv.NewWriter(os.Stdout)
	case FormatJSON:
	}

	return writer, nil
}

// Write writes the specified device.
func (writer *deviceWriter) Write(dev ble.Device) error {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	obj := newDeviceObject(dev)

	switch writer.format {
	case FormatJSON:
		// The devices are written as JSON Lines, because they are
		// written as they are discovered.
		b, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		if !writer.wroteHeader {
			if err := writer.csvWriter.Write(deviceHeader()); err != nil {
				return err
			}
			writer.wroteHeader = true
		}
		if err := writer.csvWriter.Write(obj.record()); err != nil {
			return err
		}
		writer.csvWriter.Flush()
	case FormatTable:
		if !writer.wroteHeader {
			fmt.Fprintln(writer.tabWriter, strings.Join(deviceHeader(), "\t"))
			writer.wroteHeader = true
		}
		fmt.Fprintln(writer.tabWriter, strings.Join(obj.record(), "\t"))
		writer.tabWriter.Flush()
	}

	return nil
}

// Flush flushes the buffered output.
func (writer *deviceWriter) Flush() {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	if writer.tabWriter != nil {
		writer.tabWriter.Flush()
	}
	if writer.csvWriter != nil {
		writer.csvWriter.Flush()
	}
}

// characteristicObject represents a characteristic for the JSON and the CSV
// output.
type characteristicObject struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	ID   string `json:"id"`
	MTU  int    `json:"mtu"`
}

// outputCharacteristics writes the characteristics of the specified service.
func outputCharacteristics(service ble.Service) error {
	format, err := outputFormat()
	if err != nil {
		return err
	}

	chars := service.Characteristics()
	objs := make([]characteristicObject, 0, len(chars))
	for _, char := range chars {
		mtu, err := char.MTU()
		if err != nil {
			mtu = 0
		}
		objs = append(objs, characteristicObject{
			UUID: char.UUID().String(),
			Name: char.Name(),
			ID:   char.ID(),
			MTU:  mtu,
		})
	}

	record := func(obj characteristicObject) []string {
		return []string{obj.UUID, obj.Name, obj.ID, strconv.Itoa(obj.MTU)}
	}
	header := []string{uuidColumnStr, nameColumnStr, "ID", "MTU"}

	switch format {
	case FormatJSON:
		b, err := json.Marshal(objs)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		csvWriter := csv.NewWriter(os.Stdout)
		defer csvWriter.Flush()
		if err := csvWriter.Write(header); err != nil {
			return err
		}
		for _, obj := range objs {
			if err := csvWriter.Write(record(obj)); err != nil {
				return err
			}
		}
	case FormatTable:
		tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer tabWriter.Flush()
		fmt.Fprintln(tabWriter, strings.Join(header, "\t"))
		for _, obj := range objs {
			fmt.Fprintln(tabWriter, strings.Join(record(obj), "\t"))
		}
	}

	return nil
}

// assignedNumberObject represents a Bluetooth SIG assigned number for the JSON
// and the CSV output.
type assignedNumberObject struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

// outputAssignedNumbers writes the specified assigned numbers.
func outputAssignedNumbers(objs []assignedNumberObject) error {
	format, err := outputFormat()
	if err != nil {
		return err
	}

	record := func(obj assignedNumberObject) []string {
		return []string{obj.Kind, obj.Value, obj.Name, obj.ID}
	}
	header := []string{"KIND", "VALUE", nameColumnStr, "ID"}

	switch format {
	case FormatJSON:
		b, err := json.Marshal(objs)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		csvWriter := csv.NewWriter(os.Stdout)
		defer csvWriter.Flush()
		if err := csvWriter.Write(header); err != nil {
			return err
		}
		for _, obj := range objs {
			if err := csvWriter.Write(record(obj)); err != nil {
				return err
			}
		}
	case FormatTable:
		tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer tabWriter.Flush()
		fmt.Fprintln(tabWriter, strings.Join(header, "\t"))
		for _, obj := range objs {
			fmt.Fprintln(tabWriter, strings.Join(record(obj), "\t"))
		}
	}

	return nil
}
