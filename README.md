# go-ble

![](https://img.shields.io/badge/status-Work%20In%20Progress-8A2BE2)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-ble)
[![test](https://github.com/cybergarage/go-ble/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-ble/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-ble.svg)](https://pkg.go.dev/github.com/cybergarage/go-ble)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/go-ble)
 [![codecov](https://codecov.io/gh/cybergarage/go-ble/graph/badge.svg?token=7Y64KS92VD)](https://codecov.io/gh/cybergarage/go-ble)

go-ble is a Go library for Bluetooth Low Energy (BLE). It is built on [tinygo.org/x/bluetooth](https://tinygo.org/x/bluetooth), and it adds a simpler interface and the Bluetooth SIG assigned numbers, so that scanning for a device and talking to its characteristics does not require knowing the underlying stack.

## Status

go-ble is a **central (client)** library. It scans for the devices which are advertising nearby, connects to one of them, and reads, writes and subscribes to its characteristics.

**The peripheral (server) side is not implemented.** go-ble cannot advertise a service or serve a GATT database. Use [tinygo.org/x/bluetooth](https://tinygo.org/x/bluetooth) directly for that. The peripheral role is planned for v1.0.0.

| | Status |
| --- | --- |
| Scanning, connecting, GATT client (central) | Supported |
| Advertising and serving a GATT database (peripheral) | **Not supported** |
| `blelookup` command | Supported |

### What the central supports

- Scanning with the service UUID, address, local name and RSSI filters
- Reading the advertisement: local name, RSSI, service data, and every manufacturer specific data element
- Connecting to a device, and being told when a connection is dropped
- Reading, writing (with and without a response) and subscribing to a characteristic
- The negotiated ATT MTU of a connection
- A transport which reads and writes a service through its read, write and notify characteristics
- The Bluetooth SIG assigned numbers: the company identifiers, and the service and characteristic UUIDs

### What is not supported yet

- Advertising, and serving a GATT database (the peripheral role)
- Pairing, bonding and encryption
- Descriptors, including reading and writing the CCCD directly, and telling an indication from a notification
- Selecting the adapter: the default adapter of the system is used
- Long writes, and splitting a payload larger than the MTU
- Discovering all the services of a connected device: a service is looked up by its UUID
- Expiring a device which has stopped advertising
- Stopping the adapter the moment a scan is cancelled: `Scan` returns as soon as its context is done, but the adapter stops scanning when the next advertisement arrives, because [tinygo.org/x/bluetooth](https://tinygo.org/x/bluetooth) supports `StopScan` only from inside the scan callback

## Supported platforms

| Platform | Backend | Note |
| --- | --- | --- |
| macOS | CoreBluetooth | The app needs the Bluetooth permission |
| Linux | BlueZ over D-Bus | `bluetoothd` must be running |
| Windows | WinRT | |

A device is identified by a MAC address on Linux and Windows, and by a UUID which the system assigns on macOS, so an address is not portable between the platforms.

## Install

```
go get -u github.com/cybergarage/go-ble
```

The `blelookup` command is installed with:

```
go install github.com/cybergarage/go-ble/cmd/blelookup@latest
```

## Usage

### Scanning

`Scan` blocks until the context is done, and it calls the handler for every device which passes the filters.

```go
central := ble.NewCentral()

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := central.Scan(ctx,
	ble.WithScanServiceUUIDs(0xFFF6),
	ble.WithScanHandler(func(dev ble.Device) {
		log.Printf("%s %s (%d dBm)", dev.Address(), dev.LocalName(), dev.RSSI())
	}),
)
```

The filters are optional, and `central.Devices()` returns the devices which the scan found.

### Connecting and talking to a service

```go
dev, ok := central.LookupDevice("11:22:33:AA:BB:CC")
if !ok {
	return ble.ErrNotFound
}

if err := central.Connect(ctx, dev); err != nil {
	return err
}
defer dev.Disconnect()

service, ok := dev.LookupService(0xFFF6)
if !ok {
	return ble.ErrNotFound
}

for _, char := range service.Characteristics() {
	log.Println(char.UUID(), char.Name())
}
```

A characteristic is read, written and subscribed to directly:

```go
char, ok := service.LookupCharacteristic("18EE2EF5-263D-4559-959F-4F9C429F9D11")
if !ok {
	return ble.ErrNotFound
}

mtu, _ := char.MTU() // the largest payload of a single write is MTU - 3

if _, err := char.Write(payload); err != nil {
	return err
}

err := char.Notify(func(char ble.Characteristic, buf []byte) {
	log.Printf("%s: %d bytes", char.UUID(), len(buf))
})
```

### A transport over a service

A service which carries a stream over a write characteristic and a notify characteristic is opened as a transport, which reads and writes bytes instead of characteristics.

```go
transport, err := service.Open(
	ble.WithTransportWriteUUID(writeUUID),
	ble.WithTransportNotifyUUID(notifyUUID),
)
if err != nil {
	return err
}
defer transport.Close()

if _, err := transport.Write(ctx, request); err != nil {
	return err
}

// Some peripherals flush their response only once the client enables the
// notifications after the first write, so subscribing is a separate step.
if err := transport.Subscribe(); err != nil {
	return err
}

response, err := transport.Read(ctx)
```

### Being told about a disconnection

```go
central.SetConnectionHandler(func(addr ble.Address, connected bool) {
	log.Printf("%s connected=%t", addr, connected)
})
```

### The assigned numbers

```go
db := ble.DefaultDatabase()

service, _ := db.LookupService(ble.MustUUIDFrom(0xFFF6))  // Matter Profile ID
char, _ := db.LookupCharacteristic(ble.MustUUIDFrom(0x2A00)) // Device Name
company, _ := db.LookupCompany(0x004C)                     // Apple, Inc.
```

## Command

`blelookup` scans and inspects the devices from a terminal.

```
$ blelookup scan
$ blelookup scan --service 0xFFF6 --timeout 30s
$ blelookup scan --rssi -70 --format json
$ blelookup connect 11:22:33:AA:BB:CC 0xFFF6
$ blelookup lookup 0xFFF6
```

Every command takes `--format table|json|csv` and `--timeout`. `lookup` needs no Bluetooth hardware, so a UUID which was captured elsewhere can be looked up offline.

# User Guides

- Operation
  - [blelookup](doc/blelookup.md)
- Usage
  - [Using go-ble](doc/usage.md)

# References

- [tinygo.org/x/bluetooth](https://tinygo.org/x/bluetooth)
- [Bluetooth SIG: Assigned Numbers](https://www.bluetooth.com/specifications/assigned-numbers/)
