// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	defaultAttrCntLim    = 128
	defaultAttrValLenLim = -1

	envarAttrCntLim    = "OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT"
	envarAttrValLenLim = "OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT"
)

type providerConfig struct {
	resource      *resource.Resource
	processors    []Processor
	attrCntLim    limit
	attrValLenLim limit
}

func newProviderConfig(opts []LoggerProviderOption) providerConfig {
	var c providerConfig
	for _, opt := range opts {
		c = opt.apply(c)
	}

	if c.resource == nil {
		c.resource = resource.Default()
	}

	c.attrCntLim = c.attrCntLim.Resolve(
		envarAttrCntLim,
		defaultAttrCntLim,
	)

	c.attrValLenLim = c.attrValLenLim.Resolve(
		envarAttrValLenLim,
		defaultAttrValLenLim,
	)

	return c
}

type limit struct {
	value int
	set   bool
}

func newLimit(value int) limit {
	return limit{value: value, set: true}
}

// Resolve returns the resolved form of the limit l. If l's value is set, it
// will return l. If the l's value is not set, a new limit based on the
// environment variable envar will be returned if that environment variable is
// set. Otherwise, fallback is used to construct a new limit that is returned.
func (l limit) Resolve(envar string, fallback int) limit {
	if l.set {
		return l
	}

	if v := os.Getenv(envar); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return newLimit(n)
		}
		otel.Handle(fmt.Errorf("invalid %s value %s: %w", envar, v, err))
	}

	return newLimit(fallback)
}

// Value returns the limit value if set. Otherwise, it returns -1.
func (l limit) Value() int {
	if l.set {
		return l.value
	}
	// Fail open, not closed (-1 == unlimited).
	return -1
}

// LoggerProvider handles the creation and coordination of Loggers. All Loggers
// created by a LoggerProvider will be associated with the same Resource.
type LoggerProvider struct {
	embedded.LoggerProvider

	resource                  *resource.Resource
	processors                []Processor
	attributeCountLimit       int
	attributeValueLengthLimit int

	loggersMu sync.Mutex
	loggers   map[instrumentation.Scope]*logger

	stopped atomic.Bool
}

// Compile-time check LoggerProvider implements log.LoggerProvider.
var _ log.LoggerProvider = (*LoggerProvider)(nil)

// NewLoggerProvider returns a new and configured LoggerProvider.
//
// By default, the returned LoggerProvider is configured with the default
// Resource and no Processors. Processors cannot be added after a LoggerProvider is
// created. This means the returned LoggerProvider, one created with no
// Processors, will perform no operations.
func NewLoggerProvider(opts ...LoggerProviderOption) *LoggerProvider {
	cfg := newProviderConfig(opts)
	return &LoggerProvider{
		resource:                  cfg.resource,
		processors:                cfg.processors,
		attributeCountLimit:       cfg.attrCntLim.Value(),
		attributeValueLengthLimit: cfg.attrValLenLim.Value(),
	}
}

// Logger returns a new [log.Logger] with the provided name and configuration.
//
// If p is shut down, a [noop.Logger] instace is returned.
//
// This method can be called concurrently.
func (p *LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	if p.stopped.Load() {
		return noop.NewLoggerProvider().Logger(name, opts...)
	}

	cfg := log.NewLoggerConfig(opts...)
	scope := instrumentation.Scope{
		Name:      name,
		Version:   cfg.InstrumentationVersion(),
		SchemaURL: cfg.SchemaURL(),
	}

	p.loggersMu.Lock()
	defer p.loggersMu.Unlock()

	l, ok := p.loggers[scope]
	if !ok {
		l = newLogger(p, scope)
		p.loggers[scope] = l
	}

	return l
}

// Shutdown flushes queued log records and shuts down the decorated expoter.
//
// This method can be called concurrently.
func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	stopped := p.stopped.Swap(true)
	if stopped {
		return nil
	}

	var err error
	for _, p := range p.processors {
		err = errors.Join(err, p.Shutdown(ctx))
	}
	return err
}

// ForceFlush flushes all exporters.
//
//	This method can be called concurrently.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	// TODO (#5060): Implement.
	return nil
}

// LoggerProviderOption applies a configuration option value to a LoggerProvider.
type LoggerProviderOption interface {
	apply(providerConfig) providerConfig
}

type loggerProviderOptionFunc func(providerConfig) providerConfig

func (fn loggerProviderOptionFunc) apply(c providerConfig) providerConfig {
	return fn(c)
}

// WithResource associates a Resource with a LoggerProvider. This Resource
// represents the entity producing telemetry and is associated with all Loggers
// the LoggerProvider will create.
//
// By default, if this Option is not used, the default Resource from the
// go.opentelemetry.io/otel/sdk/resource package will be used.
func WithResource(res *resource.Resource) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.resource = res
		return cfg
	})
}

// WithProcessor associates Processor with a LoggerProvider.
//
// By default, if this option is not used, the LoggerProvider will perform no
// operations; no data will be exported without a processor.
//
// Each WithProcessor creates a separate pipeline. Use custom decorators
// for advanced scenarios such as enriching with attributes.
//
// For production, use [NewBatchingProcessor] to batch log records before they are exported.
// For testing and debugging, use [NewSimpleProcessor] to synchronously export log records.
func WithProcessor(processor Processor) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.processors = append(cfg.processors, processor)
		return cfg
	})
}

// WithAttributeCountLimit sets the maximum allowed log record attribute count.
// Any attribute added to a log record once this limit is reached will be dropped.
//
// Setting this to zero means no attributes will be recorded.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 128 will be used.
func WithAttributeCountLimit(limit int) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.attrCntLim = newLimit(limit)
		return cfg
	})
}

// AttributeValueLengthLimit sets the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value will be truncated to this length.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit (-1) will be used.
func WithAttributeValueLengthLimit(limit int) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.attrValLenLim = newLimit(limit)
		return cfg
	})
}
