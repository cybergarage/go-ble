# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.0] - 2026-09-22

The first tagged release. The library is a central (client): the peripheral role is not implemented, and it is planned for v1.0.0.

### Added

- `Central.SetConnectionHandler()` reports a device being connected or disconnected.
- `Characteristic.MTU()` and `Transport.MTU()` return the negotiated ATT MTU of a connection. The largest payload of a single write is MTU - 3.
- `WithScanServiceUUIDs()`, `WithScanAddresses()`, `WithScanLocalNames()` and `WithScanRSSIThreshold()` filter a scan, and `WithScanHandler()` sets the handler.
- `Scanner.LookupDevice()` finds a discovered device by its address.
- `NewAddressFromString()` parses an address, and `Address.Bytes()` and `Address.Equal()` compare one.
- `Device.Manufacturers()` returns every manufacturer specific data element of an advertisement.
- `NewUUIDFromString()` accepts the short form of an assigned UUID, such as `FFF6` or `0xFFF6`.
- `blelookup` gained the `connect` and `lookup` commands, the `--service`, `--address`, `--name` and `--rssi` scan filters, and the `--timeout` flag. Every command supports `--format table|json|csv`.

### Changed

- `Address.String()` formats a MAC address as `11:22:33:AA:BB:CC` and a macOS device UUID as a UUID, instead of a plain hex string.
- `Device.Manufacturer()` returns the first manufacturer specific data element instead of the last one.
- `Transport.Read()` waits for a notification instead of polling every 100 milliseconds, and the notifications are buffered in a bounded channel.

### Fixed

- A UUID was built in two different word orders, so `NewUUIDFromUUID16(0xFFF6)` was not equal to `NewUUIDFromString("0000FFF6-0000-1000-8000-00805F9B34FB")` and printed as `5F9B34FB-8000-0080-0000-10000000FFF6`. A characteristic therefore never matched its entry in the assigned numbers database and its name was always empty.
- The shared adapter was enabled on every scan, and the underlying stack rejects a second `Enable()`, so a program which ran more than one scan failed after the first one.
- Cancelling the context did not stop a scan until the next advertisement arrived, so a scan on a quiet link ran until the process ended.
- The discovered devices were held in a map which the adapter callback goroutine wrote while the caller read it, and the manufacturer of a device was built on demand from those same goroutines. Both are data races which `go test -race` reports.
- Only the last manufacturer specific data element of an advertisement was kept.
- The base characteristic did not implement `WriteWithoutResponse`, so it did not satisfy the `Characteristic` interface on its own.
