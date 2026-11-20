// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

var errDuplicateRegister = errors.New("duplicate reader registration")

// ErrReaderNotRegistered is returned if Collect is called before the reader is
// registered with a ProfilerProvider.
var ErrReaderNotRegistered = errors.New("reader is not registered")

// ErrReaderShutdown is returned if Collect, ForceFlush, or Shutdown is called
// after a reader has been shut down.
var ErrReaderShutdown = errors.New("reader is shutdown")

// Reader collects profile data from a [ProfilerProvider] and any attached
// [Producer] instances.
type Reader interface {
	register(sdkProducer)

	// Collect gathers all profile data related to the Reader from the SDK and
	// attached Producers and stores the result in rp.
	//
	// Collect returns an error if called after Shutdown, if rp is nil, or if
	// ctx is canceled. This method is safe to call concurrently.
	Collect(ctx context.Context, rp *profiledata.ResourceProfiles) error

	// ForceFlush flushes pending profile telemetry.
	//
	// This method is safe to call concurrently.
	ForceFlush(ctx context.Context) error

	// Shutdown flushes pending profile telemetry and releases held resources.
	//
	// This method is safe to call concurrently.
	Shutdown(ctx context.Context) error
}

type sdkProducer interface {
	produce(context.Context, *profiledata.ResourceProfiles) error
}

type produceHolder struct {
	produce func(context.Context, *profiledata.ResourceProfiles) error
}

type shutdownProducer struct{}

func (shutdownProducer) produce(context.Context, *profiledata.ResourceProfiles) error {
	return ErrReaderShutdown
}

// ProducerOption configures a [Producer] on a [ProfilerProvider],
// [ManualReader], or [PeriodicReader].
type ProducerOption interface {
	Option
	ManualReaderOption
	PeriodicReaderOption
}

type producerOption struct {
	p Producer
}

func (o producerOption) apply(c config) config {
	if o.p != nil {
		c.producers = append(c.producers, o.p)
	}
	return c
}

func (o producerOption) applyManual(c manualReaderConfig) manualReaderConfig {
	if o.p != nil {
		c.producers = append(c.producers, o.p)
	}
	return c
}

func (o producerOption) applyPeriodic(c periodicReaderConfig) periodicReaderConfig {
	if o.p != nil {
		c.producers = append(c.producers, o.p)
	}
	return c
}

// WithProducer registers a [Producer] with a [ProfilerProvider],
// [ManualReader], or [PeriodicReader]. Nil producers are ignored.
func WithProducer(p Producer) ProducerOption {
	return producerOption{p: p}
}

func joinErrors(a, b error) error {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return errors.Join(a, b)
}
