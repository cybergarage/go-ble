# Using go-ble

go-ble is a central (client) library for Bluetooth Low Energy. This guide covers what a central does: scanning for the devices nearby, connecting to one of them, and talking to its characteristics.

The peripheral role is not implemented, so go-ble cannot advertise a service. See the [README](../README.md#status) for the current status.

## The central

A central owns the adapter of the system and the devices which it has found.

```go
central := ble.NewCentral()
```

The adapter is enabled on the first scan or connection, and it is shared by every central in the process, so a program which creates more than one works.

## Scanning

`Scan` blocks until the context is done, and it calls the handler for every device which passes the filters.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := central.Scan(ctx,
	ble.WithScanHandler(func(dev ble.Device) {
		log.Printf("%s %s (%d dBm)", dev.Address(), dev.LocalName(), dev.RSSI())
	}),
)
```

When the context has no deadline, the scan stops after `ble.DefaultScanTimeout`. Cancelling the context stops it immediately, even when nothing is advertising.

### Filtering

A link is usually busy, so a scan is narrowed to the devices which are worth looking at. The filters apply before the handler is called and before the device is stored, so `central.Devices()` holds only the matched devices.

```go
central.Scan(ctx,
	ble.WithScanServiceUUIDs(0xFFF6),          // advertises this service
	ble.WithScanAddresses("11:22:33:AA:BB:CC"), // is this device
	ble.WithScanLocalNames("My Device"),        // advertises this name
	ble.WithScanRSSIThreshold(-70),             // is close enough
	ble.WithScanHandler(handler),
)
```

### The advertisement

```go
dev.Address()       // 11:22:33:AA:BB:CC, or a UUID on macOS
dev.LocalName()     // the advertised name
dev.RSSI()          // the signal strength in dBm
dev.Services()      // the advertised services, with their service data
dev.Manufacturer()  // the first manufacturer specific data element
dev.Manufacturers() // all of them
dev.LastSeenAt()    // when the device last advertised
```

The service data of an advertisement is read without connecting, which is how a device announces what it is:

```go
service, ok := dev.LookupService(0xFFF6)
if ok {
	payload := service.Data()
}
```

## Connecting

```go
if err := central.Connect(ctx, dev); err != nil {
	return err
}
defer dev.Disconnect()
```

A service is looked up by its UUID on a connected device, which discovers it and its characteristics. There is no call which discovers every service, so the UUID has to be known.

```go
service, ok := dev.LookupService(0xFFF6)
if !ok {
	return ble.ErrNotFound
}

for _, char := range service.Characteristics() {
	log.Println(char.UUID(), char.Name())
}
```

A peripheral can drop a connection at any time, and a central which does not notice it keeps writing to a connection which is gone:

```go
central.SetConnectionHandler(func(addr ble.Address, connected bool) {
	if !connected {
		log.Printf("%s went away", addr)
	}
})
```

## Characteristics

```go
char, ok := service.LookupCharacteristic("18EE2EF5-263D-4559-959F-4F9C429F9D11")

data, err := char.Read()
n, err := char.Write(payload)                 // with an ATT response
n, err := char.WriteWithoutResponse(payload)  // without one

err := char.Notify(func(char ble.Characteristic, buf []byte) {
	// called from the Bluetooth stack: do not block here
})
```

### The MTU

A payload larger than the negotiated ATT MTU is rejected by the peer, and the MTU is known only while the device is connected.

```go
mtu, err := char.MTU()
maxPayload := mtu - 3 // the ATT write header takes three bytes
```

go-ble does not split a payload for you, so a caller which sends more than that has to segment it itself.

### With or without a response

Not every peripheral advertises the GATT "Write" property on a characteristic which it accepts writes on. On macOS, CoreBluetooth rejects a with-response write to such a characteristic before anything is transmitted, so falling back is safe:

```go
if _, err := char.Write(payload); err != nil {
	if _, err := char.WriteWithoutResponse(payload); err != nil {
		return err
	}
}
```

## The transport

A service which carries a stream over a write characteristic and a notify characteristic is opened as a transport, which reads and writes bytes.

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

if err := transport.Subscribe(); err != nil {
	return err
}

response, err := transport.Read(ctx)
```

`Subscribe` is separate from `Open` on purpose: some peripherals flush their response only once the client enables the notifications *after* the first write, and never deliver it when the notifications were already enabled beforehand.

`Read` waits for the next notification until the context is done. The notifications are buffered, and one is dropped with a warning when the buffer is full, because blocking the callback of the Bluetooth stack would stall the connection.

## UUIDs

A UUID is built from an assigned 16-bit or 32-bit number, from a string, or from bytes, and the short form is expanded with the Bluetooth Base UUID:

```go
ble.NewUUIDFromUUID16(0xFFF6)
ble.MustUUIDFromString("0000FFF6-0000-1000-8000-00805F9B34FB")
ble.MustUUIDFromString("0xFFF6")
ble.MustUUIDFrom(0xFFF6)
```

All four are the same UUID, and `String()` writes it as `0000FFF6-0000-1000-8000-00805F9B34FB`.

Anything which takes a UUID accepts these forms directly:

```go
dev.LookupService(0xFFF6)
service.LookupCharacteristic("18EE2EF5-263D-4559-959F-4F9C429F9D11")
```

## The assigned numbers

The Bluetooth SIG tables are embedded in the library, so a UUID or a company identifier can be named without a network.

```go
db := ble.DefaultDatabase()

service, _ := db.LookupService(ble.MustUUIDFrom(0xFFF6))     // Matter Profile ID
char, _ := db.LookupCharacteristic(ble.MustUUIDFrom(0x2A00)) // Device Name
company, _ := db.LookupCompany(0x004C)                       // Apple, Inc.
```

A discovered service and characteristic are named from the same tables, so `service.Name()` and `char.Name()` are filled in when the UUID is an assigned one.

## Permissions

- **macOS**: the application needs the Bluetooth permission. A command line tool inherits it from the terminal application, which is asked on the first scan.
- **Linux**: `bluetoothd` must be running, and the user needs access to the D-Bus interface of BlueZ.

## The command

The `blelookup` command does the same from a terminal, and it is the quickest way to see what is advertising nearby. See the [command reference](blelookup.md).

```
$ blelookup scan
$ blelookup scan --service 0xFFF6 --timeout 30s
$ blelookup connect 11:22:33:AA:BB:CC 0xFFF6
$ blelookup lookup 0xFFF6
```
