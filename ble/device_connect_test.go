// Copyright (C) 2026 The go-ble Authors. All rights reserved.
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
	"errors"
	"testing"
	"time"

	"tinygo.org/x/bluetooth"
)

func TestConnectWithContextReturnsTheConnection(t *testing.T) {
	want := bluetooth.Device{}
	want.Address = bluetooth.Address{}
	dev, err := connectWithContext(context.Background(),
		func() (bluetooth.Device, error) { return want, nil },
		func() { t.Error("cancel called for a connection which completed") },
		func(bluetooth.Device) { t.Error("closeLate called for a connection which completed") },
	)
	if err != nil || dev != want {
		t.Fatalf("connectWithContext() = %v, %v", dev, err)
	}

	failure := errors.New("le-connection-abort-by-local")
	if _, err := connectWithContext(context.Background(),
		func() (bluetooth.Device, error) { return bluetooth.Device{}, failure },
		func() { t.Error("cancel called for a failed connection") },
		func(bluetooth.Device) {},
	); !errors.Is(err, failure) {
		t.Fatalf("failed connection: %v", err)
	}
}

// A connection which hangs, as tinygo's Linux implementation does while
// BlueZ does not answer, ends with the context; cancelling it makes the
// pending call return.
func TestConnectWithContextCancelsAHangingConnection(t *testing.T) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelCtx()

	unblock := make(chan struct{})
	returned := make(chan struct{})
	cancelled := make(chan struct{}, 1)
	start := time.Now()
	_, err := connectWithContext(ctx,
		func() (bluetooth.Device, error) {
			defer close(returned)
			<-unblock
			return bluetooth.Device{}, errors.New("le-connection-abort-by-local")
		},
		func() {
			cancelled <- struct{}{}
			close(unblock)
		},
		func(bluetooth.Device) { t.Error("closeLate called for a cancelled connection") },
	)
	if !errors.Is(err, ErrConnectionCanceled) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("connectWithContext() error = %v", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("returned after %v", elapsed)
	}
	select {
	case <-cancelled:
	default:
		t.Fatal("the pending connection was not cancelled")
	}
	select {
	case <-returned:
	case <-time.After(2 * time.Second):
		t.Fatal("the pending call did not return after the cancel")
	}
}

// A connection which completes after the caller gave up is closed.
func TestConnectWithContextClosesALateConnection(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	release := make(chan struct{})
	closed := make(chan bluetooth.Device, 1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancelCtx()
	}()
	_, err := connectWithContext(ctx,
		func() (bluetooth.Device, error) {
			<-release
			return bluetooth.Device{}, nil
		},
		func() { close(release) },
		func(late bluetooth.Device) { closed <- late },
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("connectWithContext() error = %v", err)
	}
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("the late connection was left open")
	}
}

func TestConnectWithAContextWhichHasEnded(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dev := &tinyDevice{baseDevice: newBaseDevice()}
	// Connect returns before it touches the adapter.
	if err := dev.Connect(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Connect() = %v", err)
	}
}

// Connect waits for a scan which still holds the adapter to stop, for up to
// scanStopWait, and not past its context.
func TestWaitForScanStop(t *testing.T) {
	saved := scanStopWait
	t.Cleanup(func() { scanStopWait = saved })
	scanStopWait = 100 * time.Millisecond

	// No scan: no wait.
	start := time.Now()
	if err := waitForScanStop(context.Background()); err != nil || time.Since(start) > 50*time.Millisecond {
		t.Fatalf("without a scan: %v after %v", err, time.Since(start))
	}

	// A scan which stops: the wait ends with it.
	scanOwner <- struct{}{}
	go func() {
		time.Sleep(20 * time.Millisecond)
		<-scanOwner
	}()
	start = time.Now()
	if err := waitForScanStop(context.Background()); err != nil || time.Since(start) < 15*time.Millisecond {
		t.Fatalf("with a stopping scan: %v after %v", err, time.Since(start))
	}

	// A scan which does not stop: the wait ends after scanStopWait, or
	// with the context.
	scanOwner <- struct{}{}
	defer func() { <-scanOwner }()
	start = time.Now()
	if err := waitForScanStop(context.Background()); err != nil || time.Since(start) < scanStopWait {
		t.Fatalf("with a scan which does not stop: %v after %v", err, time.Since(start))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := waitForScanStop(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled: %v", err)
	}
}
