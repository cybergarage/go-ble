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
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	// DefaultTransportTimeout is the default timeout for transport operations.
	DefaultTransportTimeout = 5 * time.Second
)

// TransportOption represents a function type to set transport options.
type TransportOption func(*transport)

// Transport represents the BLE transport layer.
type Transport interface {
	// Open opens the transport for communication.
	Open() error
	// Close closes the transport and releases resources.
	Close() error
	// WriterCharacteristic returns the characteristic used for writing data.
	WriteCharacteristic() (Characteristic, error)
	// ReadCharacteristic returns the characteristic used for reading data.
	ReadCharacteristic() (Characteristic, error)
	// NotifyCharacteristic returns the characteristic used for notifications.
	NotifyCharacteristic() (Characteristic, error)
	// Read reads bytes from the transport.
	Read(ctx context.Context) ([]byte, error)
	// Write writes the specified bytes to the transport.
	Write(ctx context.Context, data []byte) (int, error)
	// WriteWithoutResponse writes the specified bytes to the transport without waiting for a response.
	WriteWithoutResponse(ctx context.Context, data []byte) (int, error)
}

type transport struct {
	sync.Mutex
	notifyBytes *list.List
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
		Mutex:       sync.Mutex{},
		notifyBytes: list.New(),
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
	if t.notifyCh != nil {
		notifyHandler := func(char Characteristic, buf []byte) {
			t.Lock()
			t.notifyBytes.PushBack(buf)
			t.Unlock()
		}
		if err := t.notifyCh.Notify(notifyHandler); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the transport and releases resources.
func (t *transport) Close() error {
	return nil
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
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				t.Lock()
				if 0 < t.notifyBytes.Len() {
					elem := t.notifyBytes.Front()
					t.notifyBytes.Remove(elem)
					t.Unlock()
					b, ok := elem.Value.([]byte)
					if !ok {
						return nil, fmt.Errorf("%w type: %T", ErrInvalid, elem.Value)
					}
					return b, nil
				}
				t.Unlock()
			}
			time.Sleep(100 * time.Millisecond)
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
