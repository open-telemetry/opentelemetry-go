// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/internal/attrdedup"
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
	attrCntLim    setting[int]
	attrValLenLim setting[int]
	allowDupKeys  setting[bool]
}

type experimentalOption interface {
	Experimental()
}

func newProviderConfig(opts []LoggerProviderOption) providerConfig {
	var c providerConfig
	for _, opt := range opts {
		if _, ok := opt.(experimentalOption); ok {
			continue
		}
		c = opt.apply(c)
	}

	if c.resource == nil {
		c.resource = resource.Default()
	}

	c.attrCntLim = c.attrCntLim.Resolve(
		getenv[int](envarAttrCntLim),
		fallback[int](defaultAttrCntLim),
	)

	c.attrValLenLim = c.attrValLenLim.Resolve(
		getenv[int](envarAttrValLenLim),
		fallback[int](defaultAttrValLenLim),
	)

	return c
}

// LoggerProvider handles the creation and coordination of Loggers. All Loggers
// created by a LoggerProvider will be associated with the same Resource.
type LoggerProvider struct {
	embedded.LoggerProvider

	resource                  *resource.Resource
	processors                []Processor
	attributeCountLimit       int
	attributeValueLengthLimit int
	allowDupKeys              bool

	loggersMu sync.Mutex
	loggers   map[instrumentation.Scope]*logger

	stopped         atomic.Bool
	processorCallMu sync.Mutex
	processorCalls  sync.WaitGroup

	noCmp [0]func() //nolint: unused  // This is indeed used.
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
		attributeCountLimit:       cfg.attrCntLim.Value,
		attributeValueLengthLimit: cfg.attrValLenLim.Value,
		allowDupKeys:              cfg.allowDupKeys.Value,
	}
}

// beginProcessorCall admits processor work that Shutdown needs to wait for.
func (p *LoggerProvider) beginProcessorCall() bool {
	if p.stopped.Load() {
		return false
	}

	p.processorCallMu.Lock()
	defer p.processorCallMu.Unlock()
	if p.stopped.Load() {
		return false
	}
	p.processorCalls.Add(1)
	return true
}

func (p *LoggerProvider) endProcessorCall() {
	p.processorCalls.Done()
}

// stop closes admission. The mutex ensures no call to Add can race with Wait.
func (p *LoggerProvider) stop() bool {
	if p.stopped.Load() {
		return false
	}

	p.processorCallMu.Lock()
	defer p.processorCallMu.Unlock()
	return !p.stopped.Swap(true)
}

// Logger returns a new [log.Logger] with the provided name and configuration.
//
// If p is shut down, a [noop.Logger] instance is returned.
//
// This method can be called concurrently.
func (p *LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	if name == "" {
		global.Warn("Invalid Logger name.", "name", name)
	}

	if p.stopped.Load() {
		return noop.NewLoggerProvider().Logger(name, opts...)
	}

	cfg := log.NewLoggerConfig(opts...)
	attrs := cfg.InstrumentationAttributes()
	if !p.allowDupKeys {
		attrs, _ = attrdedup.Set(attrs)
	}
	scope := instrumentation.Scope{
		Name:       name,
		Version:    cfg.InstrumentationVersion(),
		SchemaURL:  cfg.SchemaURL(),
		Attributes: attrs,
	}

	p.loggersMu.Lock()
	defer p.loggersMu.Unlock()

	if p.loggers == nil {
		l := newLogger(p, scope)
		p.loggers = map[instrumentation.Scope]*logger{scope: l}
		return l
	}

	l, ok := p.loggers[scope]
	if !ok {
		l = newLogger(p, scope)
		p.loggers[scope] = l
	}

	return l
}

// Shutdown shuts down the provider and all processors in the order they were
// registered.
//
// Shutdown first prevents new processor calls and waits for in-flight Enabled,
// OnEmit, and ForceFlush calls to complete. It then invokes each processor's
// Shutdown once. Processor Shutdown is therefore not called concurrently with
// any processor method, including itself.
//
// Shutdown waits for in-flight processor calls to complete even if ctx is
// canceled. The canceled context is then passed to the processors.
//
// After the first call to Shutdown, subsequent calls to the provider and its
// loggers will not invoke processors. Other concurrent and subsequent Shutdown
// calls return nil without invoking processors.
//
// Shutdown must not be called from a Processor's Enabled, OnEmit, or ForceFlush
// method.
//
// This method can be called concurrently.
func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	if !p.stop() {
		return nil
	}
	p.processorCalls.Wait()

	var err error
	for _, processor := range p.processors {
		err = errors.Join(err, processor.Shutdown(ctx))
	}
	return err
}

// ForceFlush flushes all processors.
//
// This method can be called concurrently.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	if !p.beginProcessorCall() {
		return nil
	}
	defer p.endProcessorCall()

	var err error
	for _, processor := range p.processors {
		err = errors.Join(err, processor.ForceFlush(ctx))
	}
	return err
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
		var err error
		cfg.resource, err = resource.Merge(resource.Environment(), res)
		if err != nil {
			otel.Handle(err)
		}
		return cfg
	})
}

// WithProcessor associates Processor with a LoggerProvider.
//
// By default, if this option is not used, the LoggerProvider will perform no
// operations; no data will be exported without a processor.
//
// The SDK invokes the processors sequentially in the same order as they were
// registered.
//
// The LoggerProvider assumes exclusive ownership of the Processor's lifecycle.
// Registering the same Processor with multiple LoggerProviders or multiple
// times with the same LoggerProvider is not supported and may result in
// concurrent or repeated calls to Shutdown.
//
// For production, use [NewBatchProcessor] to batch log records before they are exported.
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
		cfg.attrCntLim = newSetting(limit)
		return cfg
	})
}

// WithAttributeValueLengthLimit sets the maximum allowed attribute value length.
//
// This limit only applies to string, string slice, and byte slice attribute values.
// Strings and byte slices longer than this value will be truncated to this length.
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
		cfg.attrValLenLim = newSetting(limit)
		return cfg
	})
}

// WithAllowKeyDuplication sets whether deduplication is skipped for log record
// and instrumentation scope key-value collections.
//
// By default, the key-value collections within a log record and
// instrumentation scope are deduplicated to comply with the OpenTelemetry
// Specification.
// Deduplication means that if multiple key-value pairs with the same key are
// present, only a single pair is retained and others are discarded. Resource
// attributes are always deduplicated by go.opentelemetry.io/otel/sdk/resource.
//
// Disabling deduplication with this option can improve performance e.g. of
// adding attributes to the log record.
//
// Receivers may handle duplicate keys unpredictably. If you disable
// deduplication, you are responsible for ensuring that duplicate keys within a
// single collection are not emitted, or that the telemetry receiver can handle
// such duplicates.
func WithAllowKeyDuplication() LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.allowDupKeys = newSetting(true)
		return cfg
	})
}
