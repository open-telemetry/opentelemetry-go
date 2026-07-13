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

	stopped                     atomic.Bool
	processorOperationMu        sync.Mutex
	processorOperations         atomic.Int64
	processorOperationsDone     chan struct{}
	processorOperationsDoneOnce sync.Once
	shutdownState               *shutdownState

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// shutdownState publishes one cleanup result to all Shutdown callers. The err
// field is written before done is closed and read only after done is closed.
type shutdownState struct {
	done chan struct{}
	err  error
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
// The first call to Shutdown stops admitting new operations that invoke
// processor Enabled, OnEmit, or ForceFlush methods. It then starts a background
// cleanup that waits for admitted operations to complete before invoking each
// processor's Shutdown once. Processor Shutdown is therefore not called
// concurrently with any processor method, including itself.
//
// Each call to Shutdown waits for the shared cleanup to finish or for its ctx to
// be canceled. Cancellation and deadlines bound only that call's wait, not how
// long cleanup runs. The context from the call that starts cleanup supplies
// values to processor Shutdown, but its cancellation and deadline are detached.
// Contexts from later calls only bound those calls' waits. If cleanup and ctx
// cancellation are both ready to be observed, cleanup takes priority.
//
// If a call's ctx is canceled before cleanup finishes, Shutdown returns
// ctx.Err(). A later call can wait for the same cleanup and receive its result.
//
// Once shutdown starts, no new processor operation is admitted through
// ForceFlush or Logger methods. Concurrent and subsequent Shutdown calls do not
// invoke processor Shutdown again.
//
// Shutdown must not be called directly or indirectly from any Processor method.
//
// This method can be called concurrently.
func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	state := p.initShutdown(ctx)
	return waitForShutdown(ctx, state)
}

// ForceFlush flushes all processors.
//
// This method can be called concurrently.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	if !p.beginProcessorOperation() {
		return nil
	}
	defer p.endProcessorOperation()

	var err error
	for _, processor := range p.processors {
		err = errors.Join(err, processor.ForceFlush(ctx))
	}
	return err
}

// beginProcessorOperation admits processor work that Shutdown needs to wait for.
func (p *LoggerProvider) beginProcessorOperation() bool {
	if p.stopped.Load() {
		return false
	}

	p.processorOperationMu.Lock()
	defer p.processorOperationMu.Unlock()
	if p.stopped.Load() {
		return false
	}
	p.processorOperations.Add(1)
	return true
}

func (p *LoggerProvider) endProcessorOperation() {
	if p.processorOperations.Add(-1) == 0 && p.stopped.Load() {
		p.processorOperationsDoneOnce.Do(func() { close(p.processorOperationsDone) })
	}
}

// initShutdown closes processor-operation admission, starts processor cleanup
// once, and returns the shared shutdown state to all callers.
func (p *LoggerProvider) initShutdown(ctx context.Context) *shutdownState {
	p.processorOperationMu.Lock()
	if p.stopped.Load() {
		state := p.shutdownState
		p.processorOperationMu.Unlock()
		return state
	}

	operationsDone := make(chan struct{})
	state := &shutdownState{done: make(chan struct{})}
	p.processorOperationsDone = operationsDone
	p.shutdownState = state
	p.stopped.Store(true)
	if p.processorOperations.Load() == 0 {
		p.processorOperationsDoneOnce.Do(func() { close(operationsDone) })
	}
	p.processorOperationMu.Unlock()

	go p.shutdownProcessors(context.WithoutCancel(ctx), operationsDone, state)
	return state
}

func (p *LoggerProvider) shutdownProcessors(
	ctx context.Context,
	operationsDone <-chan struct{},
	state *shutdownState,
) {
	// initShutdown serializes with admission, so this closes only after all
	// operations admitted before shutdown have ended.
	<-operationsDone
	var err error
	for _, processor := range p.processors {
		err = errors.Join(err, processor.Shutdown(ctx))
	}
	state.err = err
	close(state.done)
}

// waitForShutdown waits for cleanup to end. Completion takes priority when it
// races with context cancellation.
func waitForShutdown(ctx context.Context, state *shutdownState) error {
	select {
	case <-state.done:
		return state.err
	case <-ctx.Done():
		select {
		case <-state.done:
			return state.err
		default:
			return ctx.Err()
		}
	}
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
