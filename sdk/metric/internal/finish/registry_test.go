// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryShutdown(t *testing.T) {
	var stops, waits atomic.Int64
	registry := &Registry{}
	registry.Register(
		func() { stops.Add(1) },
		func(context.Context) error {
			waits.Add(1)
			return nil
		},
	)

	require.NoError(t, registry.Shutdown(t.Context()))
	require.NoError(t, registry.Shutdown(t.Context()))
	registry.Register(
		func() { stops.Add(1) },
		func(context.Context) error {
			waits.Add(1)
			return nil
		},
	)

	assert.Equal(t, int64(2), stops.Load())
	assert.Equal(t, int64(1), waits.Load())
}

func TestRegistryShutdownContext(t *testing.T) {
	registry := &Registry{}
	registry.Register(
		func() {},
		func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	assert.ErrorIs(t, registry.Shutdown(ctx), context.Canceled)
}

func TestRegistryConcurrentShutdown(t *testing.T) {
	const registrations = 100
	registry := &Registry{}
	stops := make([]atomic.Int64, registrations)
	waits := make([]atomic.Int64, registrations)
	for i := range registrations {
		registry.Register(
			func() { stops[i].Add(1) },
			func(context.Context) error {
				waits[i].Add(1)
				return nil
			},
		)
	}

	start := make(chan struct{})
	done := make(chan struct{}, 10)
	for range 10 {
		go func() {
			<-start
			_ = registry.Shutdown(t.Context())
			done <- struct{}{}
		}()
	}
	close(start)
	for range 10 {
		<-done
	}

	for i := range registrations {
		assert.Equal(t, int64(1), stops[i].Load(), "stop %d", i)
		assert.Equal(t, int64(1), waits[i].Load(), "wait %d", i)
	}
}

func TestRegistryConcurrentShutdownWaits(t *testing.T) {
	registry := &Registry{}
	waiting := make(chan struct{})
	release := make(chan struct{})
	registry.Register(
		func() {},
		func(context.Context) error {
			close(waiting)
			<-release
			return nil
		},
	)
	first := make(chan error, 1)
	go func() { first <- registry.Shutdown(t.Context()) }()
	<-waiting

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, registry.Shutdown(ctx), context.Canceled)
	second := make(chan error, 1)
	go func() { second <- registry.Shutdown(t.Context()) }()
	select {
	case <-second:
		t.Fatal("concurrent shutdown returned before the registered wait completed")
	default:
	}

	close(release)
	assert.NoError(t, <-first)
	assert.NoError(t, <-second)
}

func TestRegistryConcurrentRegisterShutdown(t *testing.T) {
	const registrations = 100
	registry := &Registry{}
	calls := make([]atomic.Int64, registrations)
	start := make(chan struct{})
	done := make(chan struct{}, registrations+10)

	for i := range registrations {
		go func() {
			<-start
			registry.Register(
				func() { calls[i].Add(1) },
				func(context.Context) error { return nil },
			)
			done <- struct{}{}
		}()
	}
	for range 10 {
		go func() {
			<-start
			_ = registry.Shutdown(t.Context())
			done <- struct{}{}
		}()
	}
	close(start)
	for range registrations + 10 {
		<-done
	}

	for i := range registrations {
		assert.Equal(t, int64(1), calls[i].Load(), "registration %d", i)
	}
}
