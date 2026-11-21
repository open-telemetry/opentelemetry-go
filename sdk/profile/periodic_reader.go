// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

const (
	defaultInterval = 60 * time.Second
	defaultTimeout  = 30 * time.Second
)

type periodicReaderConfig struct {
	interval  time.Duration
	timeout   time.Duration
	producers []Producer
}

func newPeriodicReaderConfig(options []PeriodicReaderOption) periodicReaderConfig {
	c := periodicReaderConfig{
		interval: defaultInterval,
		timeout:  defaultTimeout,
	}
	for _, o := range options {
		if o != nil {
			c = o.applyPeriodic(c)
		}
	}
	return c
}

// PeriodicReaderOption applies configuration to a [PeriodicReader].
type PeriodicReaderOption interface {
	applyPeriodic(periodicReaderConfig) periodicReaderConfig
}

type periodicReaderOptionFunc func(periodicReaderConfig) periodicReaderConfig

func (f periodicReaderOptionFunc) applyPeriodic(c periodicReaderConfig) periodicReaderConfig {
	return f(c)
}

// WithInterval configures the time between periodic profile collections.
// Non-positive durations are ignored and the default of 60 seconds is used.
func WithInterval(d time.Duration) PeriodicReaderOption {
	return periodicReaderOptionFunc(func(c periodicReaderConfig) periodicReaderConfig {
		if d > 0 {
			c.interval = d
		}
		return c
	})
}

// WithTimeout configures the maximum time a PeriodicReader waits for a single
// collect and export cycle before canceling it. Non-positive durations are
// ignored and the default of 30 seconds is used.
func WithTimeout(d time.Duration) PeriodicReaderOption {
	return periodicReaderOptionFunc(func(c periodicReaderConfig) periodicReaderConfig {
		if d > 0 {
			c.timeout = d
		}
		return c
	})
}

// PeriodicReader is a [Reader] that periodically collects and exports profiles
// to an [Exporter].
type PeriodicReader struct {
	sdkProducer atomic.Value

	mu                sync.Mutex
	isShutdown        bool
	externalProducers atomic.Value

	interval time.Duration
	timeout  time.Duration
	exporter Exporter
	flushCh  chan chan error

	done         chan struct{}
	cancel       context.CancelFunc
	shutdownOnce sync.Once
}

var _ Reader = (*PeriodicReader)(nil)

// NewPeriodicReader returns a [PeriodicReader] that collects and exports
// profiles to exporter on a fixed interval. If exporter is nil, a no-op
// exporter is used.
func NewPeriodicReader(exporter Exporter, options ...PeriodicReaderOption) *PeriodicReader {
	if exporter == nil {
		exporter = noopExporter{}
	}
	conf := newPeriodicReaderConfig(options)
	ctx, cancel := context.WithCancel( //nolint:gosec // Canceled during PeriodicReader.Shutdown.
		context.Background(),
	)
	r := &PeriodicReader{
		interval: conf.interval,
		timeout:  conf.timeout,
		exporter: exporter,
		flushCh:  make(chan chan error),
		cancel:   cancel,
		done:     make(chan struct{}),
	}
	r.externalProducers.Store(conf.producers)

	go func() {
		defer close(r.done)
		r.run(ctx)
	}()

	return r
}

func (r *PeriodicReader) run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.collectAndExport(ctx); err != nil {
				otel.Handle(err)
			}
		case errCh := <-r.flushCh:
			errCh <- r.collectAndExport(ctx)
			ticker.Reset(r.interval)
		case <-ctx.Done():
			return
		}
	}
}

func (r *PeriodicReader) register(p sdkProducer) {
	if !r.sdkProducer.CompareAndSwap(nil, produceHolder{produce: p.produce}) {
		global.Error(errDuplicateRegister, "did not register periodic profile reader")
	}
}

func (r *PeriodicReader) collectAndExport(ctx context.Context) error {
	ctx, cancel := context.WithTimeoutCause(ctx, r.timeout, errors.New("profile reader collect and export timeout"))
	defer cancel()

	var rp profiledata.ResourceProfiles
	collectErr := r.collect(ctx, r.sdkProducer.Load(), &rp)
	if len(rp.ScopeProfiles) == 0 {
		return collectErr
	}
	exportErr := r.exporter.Export(ctx, &rp)
	return joinErrors(collectErr, exportErr)
}

// Collect gathers profile data into rp without exporting it.
//
// This method is safe to call concurrently.
func (r *PeriodicReader) Collect(ctx context.Context, rp *profiledata.ResourceProfiles) error {
	if rp == nil {
		return errors.New("periodic reader: *profiledata.ResourceProfiles is nil")
	}
	return r.collect(ctx, r.sdkProducer.Load(), rp)
}

func (r *PeriodicReader) collect(ctx context.Context, p any, rp *profiledata.ResourceProfiles) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if p == nil {
		return ErrReaderNotRegistered
	}

	ph, ok := p.(produceHolder)
	if !ok {
		return fmt.Errorf("periodic reader: invalid producer: %T", p)
	}

	err := ph.produce(ctx, rp)
	if err != nil && errors.Is(err, ErrReaderShutdown) {
		return err
	}

	producers, _ := r.externalProducers.Load().([]Producer)
	for _, producer := range producers {
		sps, e := producer.Produce(ctx)
		if e != nil {
			err = joinErrors(err, e)
		}
		rp.ScopeProfiles = append(rp.ScopeProfiles, sps...)
	}
	return err
}

// ForceFlush collects and exports any pending profiles immediately and calls
// ForceFlush on the underlying exporter.
//
// This method is safe to call concurrently.
func (r *PeriodicReader) ForceFlush(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeoutCause(ctx, r.timeout, errors.New("profile reader force flush timeout"))
		defer cancel()
	}

	errCh := make(chan error, 1)
	select {
	case r.flushCh <- errCh:
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	case <-r.done:
		return ErrReaderShutdown
	case <-ctx.Done():
		return ctx.Err()
	}
	return r.exporter.ForceFlush(ctx)
}

// Shutdown stops the background collection loop, performs a final collection
// and export, shuts down any reader-attached producers, and shuts down the
// exporter. Subsequent calls return [ErrReaderShutdown].
//
// This method is safe to call concurrently.
func (r *PeriodicReader) Shutdown(ctx context.Context) error {
	err := ErrReaderShutdown
	r.shutdownOnce.Do(func() {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeoutCause(ctx, r.timeout, errors.New("profile reader shutdown timeout"))
			defer cancel()
		}

		// Stop the background loop before performing the final flush so no
		// concurrent export races with Shutdown.
		r.cancel()
		<-r.done

		ph := r.sdkProducer.Swap(produceHolder{produce: shutdownProducer{}.produce})
		err = nil
		if ph != nil {
			var rp profiledata.ResourceProfiles
			collectErr := r.collect(ctx, ph, &rp)
			err = joinErrors(err, collectErr)
			if len(rp.ScopeProfiles) > 0 {
				err = joinErrors(err, r.exporter.Export(ctx, &rp))
			}
		}

		r.mu.Lock()
		r.isShutdown = true
		producers, _ := r.externalProducers.Swap([]Producer{}).([]Producer)
		r.mu.Unlock()

		err = joinErrors(err, shutdownProducers(ctx, producers))
		err = joinErrors(err, r.exporter.Shutdown(ctx))
	})
	return err
}

// MarshalLog returns logging data about the PeriodicReader.
func (r *PeriodicReader) MarshalLog() any {
	r.mu.Lock()
	down := r.isShutdown
	r.mu.Unlock()
	return struct {
		Type       string
		Exporter   Exporter
		Registered bool
		Shutdown   bool
		Interval   time.Duration
		Timeout    time.Duration
	}{
		Type:       "PeriodicReader",
		Exporter:   r.exporter,
		Registered: r.sdkProducer.Load() != nil,
		Shutdown:   down,
		Interval:   r.interval,
		Timeout:    r.timeout,
	}
}
