// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type config struct {
	res       *resource.Resource
	readers   []Reader
	producers []Producer
}

func newConfig(options []Option) config {
	var c config
	for _, o := range options {
		if o != nil {
			c = o.apply(c)
		}
	}
	if c.res == nil {
		c.res = resource.Default()
	}
	return c
}

// Option applies configuration to a [ProfilerProvider].
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (f optionFunc) apply(c config) config { return f(c) }

// WithResource associates a [resource.Resource] with all profiles produced by
// the [ProfilerProvider]. By default, [resource.Default] is used.
func WithResource(res *resource.Resource) Option {
	return optionFunc(func(c config) config {
		c.res = res
		return c
	})
}

// WithReader registers a [Reader] with the [ProfilerProvider]. Nil readers are
// ignored.
func WithReader(r Reader) Option {
	return optionFunc(func(c config) config {
		if r != nil {
			c.readers = append(c.readers, r)
		}
		return c
	})
}

// ProfilerProvider coordinates profile [Producer] instances and [Reader]
// pipelines sharing a single [resource.Resource].
type ProfilerProvider struct {
	res       *resource.Resource
	readers   []Reader
	producers []Producer

	stopped      atomic.Bool
	shutdownOnce sync.Once
}

// NewProfilerProvider returns a new configured [ProfilerProvider].
func NewProfilerProvider(options ...Option) *ProfilerProvider {
	cfg := newConfig(options)
	pp := &ProfilerProvider{
		res:       cfg.res,
		readers:   cfg.readers,
		producers: cfg.producers,
	}

	sp := providerProducer{res: cfg.res, producers: cfg.producers}
	for _, r := range pp.readers {
		r.register(sp)
	}
	return pp
}

type providerProducer struct {
	res       *resource.Resource
	producers []Producer
}

func (p providerProducer) produce(ctx context.Context, rp *profiledata.ResourceProfiles) error {
	rp.Resource = p.res
	rp.ScopeProfiles = rp.ScopeProfiles[:0]

	var err error
	for _, prod := range p.producers {
		sps, e := prod.Produce(ctx)
		if e != nil {
			err = joinErrors(err, e)
		}
		rp.ScopeProfiles = append(rp.ScopeProfiles, sps...)
	}
	return err
}

// ForceFlush flushes all registered [Reader] instances.
//
// This method is safe to call concurrently.
func (pp *ProfilerProvider) ForceFlush(ctx context.Context) error {
	if pp.stopped.Load() {
		return nil
	}
	var err error
	for _, r := range pp.readers {
		err = joinErrors(err, r.ForceFlush(ctx))
	}
	return err
}

// Shutdown flushes and shuts down all registered [Reader] instances and then
// shuts down provider-level [Producer] instances. Subsequent calls perform no
// operation and return nil.
//
// This method is safe to call concurrently.
func (pp *ProfilerProvider) Shutdown(ctx context.Context) error {
	var err error
	pp.shutdownOnce.Do(func() {
		pp.stopped.Store(true)
		// Shut down readers first so their final flush can collect from
		// provider-level producers before those producers stop.
		for _, r := range pp.readers {
			err = joinErrors(err, r.Shutdown(ctx))
		}
		err = joinErrors(err, shutdownProducers(ctx, pp.producers))
	})
	return err
}
