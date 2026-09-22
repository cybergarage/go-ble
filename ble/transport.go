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

package ble

import (
	"context"
	"sync"
	"time"

	"github.com/cybergarage/go-logger/log"
)

const (
	// DefaultTransportTimeout is the default timeout for transport operations.
	DefaultTransportTimeout = 5 * time.Second
	// DefaultTransportNotifyBufferSize is the number of the notifications
	// which are buffered before they are read. A notification is dropped
	// when the buffer is full, because blocking the callback of the
	// Bluetooth stack would stall the connection.
	DefaultTransportNotifyBufferSize = 256
)

// TransportOption represents a function type to set transport options.
type TransportOption func(*transport)

// Transport represents the BLE transport layer.
type Transport interface {
	// Open opens the transport for communication.
	Open() error
	// Subscribe enables notifications/indications on the notify
	// characteristic, if one was configured. It is separate from Open so
	// callers can control whether it happens before or after an initial
	// write: some peripherals only deliver a response to an initial write
	// once they observe the client enabling notifications afterwards, and
	// never flush it if notifications were already enabled beforehand.
	Subscribe() error
	// Close closes the transport and releases resources.
	Close() error
	// WriterCharacteristic returns the characteristic used for writing data.
	WriteCharacteristic() (Characteristic, error)
	// ReadCharacteristic returns the characteristic used for reading data.
	ReadCharacteristic() (Characteristic, error)
	// NotifyCharacteristic returns the characteristic used for notifications.
	NotifyCharacteristic() (Characteristic, error)
	// MTU returns the ATT maximum transmission unit of the connection. The
	// largest payload of a single write is MTU - 3.
	MTU() (int, error)
	// Read reads bytes from the transport.
	Read(ctx context.Context) ([]byte, error)
	// Write writes the specified bytes to the transport.
	Write(ctx context.Context, data []byte) (int, error)
	// WriteWithoutResponse writes the specified bytes to the transport without waiting for a response.
	WriteWithoutResponse(ctx context.Context, data []byte) (int, error)
}

type transport struct {
	mutex       sync.Mutex
	closed      bool
	notifyBytes chan []byte
	readCh      Characteristic
	writeCh     Characteristic
	notifyCh    Characteristic
}

// WithTransportReadCharacteristic sets the characteristic used for reading data.
func WithTransportReadCharacteristic(char Characteristic) TransportOption {
	return func(t *transport) {
		t.readCh = char
	}
}

// WithTransportWriteCharacteristic sets the characteristic used for writing data.
func WithTransportWriteCharacteristic(char Characteristic) TransportOption {
	return func(t *transport) {
		t.writeCh = char
	}
}

// WithTransportNotifyCharacteristic sets the characteristic used for notifications.
func WithTransportNotifyCharacteristic(char Characteristic) TransportOption {
	return func(t *transport) {
		t.notifyCh = char
	}
}

// NewTransport returns a new Transport instance.
func NewTransport(opts ...TransportOption) Transport {
	t := &transport{
		mutex:       sync.Mutex{},
		closed:      false,
		notifyBytes: make(chan []byte, DefaultTransportNotifyBufferSize),
		readCh:      nil,
		writeCh:     nil,
		notifyCh:    nil,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Open opens the transport for communication.
func (t *transport) Open() error {
	return nil
}

// Subscribe enables notifications/indications on the notify characteristic.
func (t *transport) Subscribe() error {
	if t.notifyCh == nil {
		return nil
	}
	notifyHandler := func(char Characteristic, buf []byte) {
		data := make([]byte, len(buf))
		copy(data, buf)

		t.mutex.Lock()
		closed := t.closed
		t.mutex.Unlock()
		if closed {
			return
		}

		// The notification is buffered without blocking, because this
		// handler is called from the Bluetooth stack.
		select {
		case t.notifyBytes <- data:
		default:
			log.Warnf("ble: notification dropped (%d bytes): the transport buffer of %s is full", len(data), char.UUID().String())
		}
	}
	return t.notifyCh.Notify(notifyHandler)
}

// Close closes the transport and releases resources.
func (t *transport) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.closed = true
	return nil
}

// MTU returns the ATT maximum transmission unit of the connection.
func (t *transport) MTU() (int, error) {
	for _, char := range []Characteristic{t.writeCh, t.notifyCh, t.readCh} {
		if char == nil {
			continue
		}
		return char.MTU()
	}
	return 0, ErrNotSet
}

// WriteCharacteristic returns the characteristic used for writing data.
func (t *transport) WriteCharacteristic() (Characteristic, error) {
	if t.writeCh == nil {
		return nil, ErrNotSet
	}
	return t.writeCh, nil
}

// ReadCharacteristic returns the characteristic used for reading data.
func (t *transport) ReadCharacteristic() (Characteristic, error) {
	if t.readCh == nil {
		return nil, ErrNotSet
	}
	return t.readCh, nil
}

// NotifyCharacteristic returns the characteristic used for notifications.
func (t *transport) NotifyCharacteristic() (Characteristic, error) {
	if t.notifyCh == nil {
		return nil, ErrNotSet
	}
	return t.notifyCh, nil
}

// Read reads bytes from the transport.
func (t *transport) Read(ctx context.Context) ([]byte, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTransportTimeout)
		defer cancel()
	}

	switch {
	case t.notifyCh != nil:
		// The notifications are awaited on the channel instead of being
		// polled, so a response is returned as soon as it arrives.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case data := <-t.notifyBytes:
			return data, nil
		}
	case t.readCh != nil:
		return t.readCh.Read()
	}

	return nil, ErrNotSet
}

// Write writes the specified bytes to the transport.
func (t *transport) Write(ctx context.Context, data []byte) (int, error) {
	if t.writeCh == nil {
		return 0, ErrNotSet
	}
	return t.writeCh.Write(data)
}

// WriteWithoutResponse writes the specified bytes to the transport without waiting for a response.
func (t *transport) WriteWithoutResponse(ctx context.Context, data []byte) (int, error) {
	if t.writeCh == nil {
		return 0, ErrNotSet
	}
	return t.writeCh.WriteWithoutResponse(data)
}
