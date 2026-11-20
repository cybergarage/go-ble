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
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/cybergarage/go-ble/ble"
)

type outputTestDevice struct {
	ble.Device
	name string
}

type outputTestManufacturer struct{ ble.Manufacturer }

func (outputTestManufacturer) Name() string { return "Manufacturer" }

func (dev outputTestDevice) LocalName() string { return dev.name }

func (outputTestDevice) Address() ble.Address { return ble.Address{1, 2, 3, 4, 5, 6} }

func (outputTestDevice) RSSI() int { return -60 }

func (outputTestDevice) Services() []ble.Service { return nil }

func (outputTestDevice) Manufacturer() ble.Manufacturer { return outputTestManufacturer{} }

func writeTestDevices(t *testing.T, format string, names ...string) string {
	t.Helper()
	t.Setenv("BLE_LOOKUP_FORMAT", format)
	var output strings.Builder
	writer, err := newDeviceWriter(&output)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if err := writer.Write(outputTestDevice{name: name}); err != nil {
			t.Fatal(err)
		}
	}
	writer.Flush()
	return output.String()
}

func TestDeviceTableEscapesControlCharacters(t *testing.T) {
	output := writeTestDevices(t, "table", "Capteur café", "Label\twith\ncontrols\x1b\u0085\u2028\u202e", "Invalid\xff")
	if strings.Count(output, "\n") != 4 || !utf8.ValidString(output) || strings.ContainsAny(output, "\t\x1b\u0085\u2028\u202e") {
		t.Fatalf("table contains unescaped characters: %q", output)
	}
	for _, want := range []string{"Capteur café", `Label\twith\ncontrols\x1b\u0085\u2028\u202e`, `Invalid\xff`} {
		if !strings.Contains(output, want) {
			t.Errorf("table did not contain %q: %q", want, output)
		}
	}
}

func TestDeviceCSVWritesFormulaLikeNamesAsText(t *testing.T) {
	names := []string{"Capteur café", `Sensor, "Front"`, "=Label", "+Label", "-Label", "@Label", " \u200b＝Label", "\tLabel", "\rLabel", "\nLabel"}
	output := writeTestDevices(t, "csv", names...)
	rows, err := csv.NewReader(strings.NewReader(output)).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(names)+1 {
		t.Fatalf("got %d rows, want %d", len(rows), len(names)+1)
	}
	for i, name := range names {
		want := name
		if i >= 2 {
			want = "'" + want
		}
		if got := rows[i+1][1]; got != want {
			t.Errorf("CSV name = %q, want %q", got, want)
		}
		if got := rows[i+1][3]; got != "-60" {
			t.Errorf("CSV RSSI = %q, want -60", got)
		}
	}
}

func TestDeviceJSONPreservesName(t *testing.T) {
	const name = "=Sensor\ncafé"
	output := writeTestDevices(t, "json", name)
	var device deviceObject
	if err := json.Unmarshal([]byte(output), &device); err != nil {
		t.Fatal(err)
	}
	if device.LocalName != name {
		t.Errorf("JSON name = %q, want %q", device.LocalName, name)
	}
}
