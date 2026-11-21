// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

// ManualReader is a [Reader] that collects profiles on demand when [ManualReader.Collect]
// is called.
type ManualReader struct {
	sdkProducer  atomic.Value
	shutdownOnce sync.Once

	mu                sync.Mutex
	isShutdown        bool
	externalProducers atomic.Value
}

var _ Reader = (*ManualReader)(nil)

// ManualReaderOption applies configuration to a [ManualReader].
type ManualReaderOption interface {
	applyManual(manualReaderConfig) manualReaderConfig
}

type manualReaderConfig struct {
	producers []Producer
}

// NewManualReader returns a [ManualReader] that collects profiles on demand.
func NewManualReader(opts ...ManualReaderOption) *ManualReader {
	var cfg manualReaderConfig
	for _, opt := range opts {
		if opt != nil {
			cfg = opt.applyManual(cfg)
		}
	}
	r := &ManualReader{}
	r.externalProducers.Store(cfg.producers)
	return r
}

func (mr *ManualReader) register(p sdkProducer) {
	if !mr.sdkProducer.CompareAndSwap(nil, produceHolder{produce: p.produce}) {
		global.Error(errDuplicateRegister, "did not register manual profile reader")
	}
}

// Collect gathers profile data from the registered [ProfilerProvider] and any
// reader-attached [Producer] instances and stores it in rp.
//
// This method is safe to call concurrently.
func (mr *ManualReader) Collect(ctx context.Context, rp *profiledata.ResourceProfiles) error {
	if rp == nil {
		return errors.New("manual reader: *profiledata.ResourceProfiles is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	p := mr.sdkProducer.Load()
	if p == nil {
		return ErrReaderNotRegistered
	}

	ph, ok := p.(produceHolder)
	if !ok {
		return fmt.Errorf("manual reader: invalid producer: %T", p)
	}

	err := ph.produce(ctx, rp)
	if err != nil && errors.Is(err, ErrReaderShutdown) {
		return err
	}

	producers, _ := mr.externalProducers.Load().([]Producer)
	for _, producer := range producers {
		sps, e := producer.Produce(ctx)
		if e != nil {
			err = joinErrors(err, e)
		}
		rp.ScopeProfiles = append(rp.ScopeProfiles, sps...)
	}
	return err
}

// ForceFlush performs no operation because ManualReader holds no buffered state.
func (mr *ManualReader) ForceFlush(ctx context.Context) error {
	mr.mu.Lock()
	down := mr.isShutdown
	mr.mu.Unlock()
	if down {
		return ErrReaderShutdown
	}
	return ctx.Err()
}

// Shutdown shuts down the ManualReader and any producers registered directly
// on it. Subsequent calls return [ErrReaderShutdown].
//
// This method is safe to call concurrently.
func (mr *ManualReader) Shutdown(ctx context.Context) error {
	err := ErrReaderShutdown
	mr.shutdownOnce.Do(func() {
		mr.sdkProducer.Store(produceHolder{produce: shutdownProducer{}.produce})
		mr.mu.Lock()
		mr.isShutdown = true
		producers, _ := mr.externalProducers.Swap([]Producer{}).([]Producer)
		mr.mu.Unlock()
		err = shutdownProducers(ctx, producers)
	})
	return err
}

// MarshalLog returns logging data about the ManualReader.
func (mr *ManualReader) MarshalLog() any {
	mr.mu.Lock()
	down := mr.isShutdown
	mr.mu.Unlock()
	return struct {
		Type       string
		Registered bool
		Shutdown   bool
	}{
		Type:       "ManualReader",
		Registered: mr.sdkProducer.Load() != nil,
		Shutdown:   down,
	}
}
