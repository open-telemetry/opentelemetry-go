// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	defaultLoggerName = "go.opentelemetry.io/otel/sdk/logger"
)

// loggerProviderConfig.
type loggerProviderConfig struct {
	// processors contains collection of LogRecordProcessors that are processing pipeline
	// for logRecords in the trace signal.
	// LogRecordProcessors registered with a LoggerProvider and are called at the start
	// and end of a Span's lifecycle, and are called in the order they are
	// registered.
	processors []LogRecordProcessor

	// logRecordLimits defines the attribute, event, and link limits for logRecords.
	logRecordLimits LogRecordLimits

	// resource contains attributes representing an entity that produces telemetry.
	resource *resource.Resource
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (cfg loggerProviderConfig) MarshalLog() interface{} {
	return struct {
		LogRecordProcessors []LogRecordProcessor
		LogRecordLimits     LogRecordLimits
		Resource            *resource.Resource
	}{
		LogRecordProcessors: cfg.processors,
		LogRecordLimits:     cfg.logRecordLimits,
		Resource:            cfg.resource,
	}
}

// LoggerProvider is an OpenTelemetry LoggerProvider. It provides Tracers to
// instrumentation so it can trace operational flow through a system.
type LoggerProvider struct {
	mu                  sync.Mutex
	namedLogger         map[instrumentation.Scope]*logger
	logRecordProcessors atomic.Value
	isShutdown          bool

	// These fields are not protected by the lock mu. They are assumed to be
	// immutable after creation of the LoggerProvider.
	logRecordLimits LogRecordLimits
	resource        *resource.Resource
}

var _ log.LoggerProvider = &LoggerProvider{}

// NewLoggerProvider returns a new and configured LoggerProvider.
//
// By default the returned LoggerProvider is configured with:
//   - a ParentBased(AlwaysSample) Sampler
//   - a random number IDGenerator
//   - the resource.Default() Resource
//   - the default LogRecordLimits.
//
// The passed opts are used to override these default values and configure the
// returned LoggerProvider appropriately.
func NewLoggerProvider(opts ...LoggerProviderOption) *LoggerProvider {
	o := loggerProviderConfig{
		logRecordLimits: NewLogRecordLimits(),
	}

	for _, opt := range opts {
		o = opt.apply(o)
	}

	o = ensureValidLoggerProviderConfig(o)

	tp := &LoggerProvider{
		namedLogger:     make(map[instrumentation.Scope]*logger),
		logRecordLimits: o.logRecordLimits,
		resource:        o.resource,
	}
	global.Info("LoggerProvider created", "config", o)

	spss := logRecordProcessorStates{}
	for _, sp := range o.processors {
		spss = append(spss, newLogRecordProcessorState(sp))
	}
	tp.logRecordProcessors.Store(spss)

	return tp
}

// Logger returns a Logger with the given name and options. If a Logger for
// the given name and options does not exist it is created, otherwise the
// existing Logger is returned.
//
// If name is empty, DefaultLoggerName is used instead.
//
// This method is safe to be called concurrently.
func (p *LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	c := log.NewLoggerConfig(opts...)

	p.mu.Lock()
	defer p.mu.Unlock()
	if name == "" {
		name = defaultLoggerName
	}
	is := instrumentation.Scope{
		Name:      name,
		Version:   c.InstrumentationVersion(),
		SchemaURL: c.SchemaURL(),
	}
	t, ok := p.namedLogger[is]
	if !ok {
		t = &logger{
			provider:             p,
			instrumentationScope: is,
		}
		p.namedLogger[is] = t
		global.Info("Logger created", "name", name, "version", c.InstrumentationVersion(), "schemaURL", c.SchemaURL())
	}
	return t
}

// RegisterLogRecordProcessor adds the given LogRecordProcessor to the list of LogRecordProcessors.
func (p *LoggerProvider) RegisterLogRecordProcessor(sp LogRecordProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isShutdown {
		return
	}
	newSPS := logRecordProcessorStates{}
	newSPS = append(newSPS, p.logRecordProcessors.Load().(logRecordProcessorStates)...)
	newSPS = append(newSPS, newLogRecordProcessorState(sp))
	p.logRecordProcessors.Store(newSPS)
}

// UnregisterLogRecordProcessor removes the given LogRecordProcessor from the list of LogRecordProcessors.
func (p *LoggerProvider) UnregisterLogRecordProcessor(sp LogRecordProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isShutdown {
		return
	}
	old := p.logRecordProcessors.Load().(logRecordProcessorStates)
	if len(old) == 0 {
		return
	}
	spss := logRecordProcessorStates{}
	spss = append(spss, old...)

	// stop the span processor if it is started and remove it from the list
	var stopOnce *logRecordProcessorState
	var idx int
	for i, sps := range spss {
		if sps.sp == sp {
			stopOnce = sps
			idx = i
		}
	}
	if stopOnce != nil {
		stopOnce.state.Do(
			func() {
				if err := sp.Shutdown(context.Background()); err != nil {
					otel.Handle(err)
				}
			},
		)
	}
	if len(spss) > 1 {
		copy(spss[idx:], spss[idx+1:])
	}
	spss[len(spss)-1] = nil
	spss = spss[:len(spss)-1]

	p.logRecordProcessors.Store(spss)
}

// ForceFlush immediately exports all logRecords that have not yet been exported for
// all the registered span processors.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	spss := p.logRecordProcessors.Load().(logRecordProcessorStates)
	if len(spss) == 0 {
		return nil
	}

	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := sps.sp.ForceFlush(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shuts down LoggerProvider. All registered span processors are shut down
// in the order they were registered and any held computational resources are released.
func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	spss := p.logRecordProcessors.Load().(logRecordProcessorStates)
	if len(spss) == 0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.isShutdown = true

	var retErr error
	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var err error
		sps.state.Do(
			func() {
				err = sps.sp.Shutdown(ctx)
			},
		)
		if err != nil {
			if retErr == nil {
				retErr = err
			} else {
				// Poor man's list of errors
				retErr = fmt.Errorf("%v; %v", retErr, err)
			}
		}
	}
	p.logRecordProcessors.Store(logRecordProcessorStates{})
	return retErr
}

// LoggerProviderOption configures a LoggerProvider.
type LoggerProviderOption interface {
	apply(loggerProviderConfig) loggerProviderConfig
}

type loggerProviderOptionFunc func(loggerProviderConfig) loggerProviderConfig

func (fn loggerProviderOptionFunc) apply(cfg loggerProviderConfig) loggerProviderConfig {
	return fn(cfg)
}

// WithSyncer registers the exporter with the LoggerProvider using a
// SimpleSpanProcessor.
//
// This is not recommended for production use. The synchronous nature of the
// SimpleSpanProcessor that will wrap the exporter make it good for testing,
// debugging, or showing examples of other feature, but it will be slow and
// have a high computation resource usage overhead. The WithBatcher option is
// recommended for production use instead.
func WithSyncer(e LogRecordExporter) LoggerProviderOption {
	return WithLogRecordProcessor(NewSimpleLogRecordProcessor(e))
}

// WithBatcher registers the exporter with the LoggerProvider using a
// BatchSpanProcessor configured with the passed opts.
func WithBatcher(e LogRecordExporter, opts ...BatchLogRecordProcessorOption) LoggerProviderOption {
	return WithLogRecordProcessor(NewBatchLogRecordProcessor(e, opts...))
}

// WithLogRecordProcessor registers the LogRecordProcessor with a LoggerProvider.
func WithLogRecordProcessor(sp LogRecordProcessor) LoggerProviderOption {
	return loggerProviderOptionFunc(
		func(cfg loggerProviderConfig) loggerProviderConfig {
			cfg.processors = append(cfg.processors, sp)
			return cfg
		},
	)
}

// WithResource returns a LoggerProviderOption that will configure the
// Resource r as a LoggerProvider's Resource. The configured Resource is
// referenced by all the Tracers the LoggerProvider creates. It represents the
// entity producing telemetry.
//
// If this option is not used, the LoggerProvider will use the
// resource.Default() Resource by default.
func WithResource(r *resource.Resource) LoggerProviderOption {
	return loggerProviderOptionFunc(
		func(cfg loggerProviderConfig) loggerProviderConfig {
			var err error
			cfg.resource, err = resource.Merge(resource.Environment(), r)
			if err != nil {
				otel.Handle(err)
			}
			return cfg
		},
	)
}

// WithLogRecordLimits returns a LoggerProviderOption that configures a
// LoggerProvider to use the LogRecordLimits sl. These LogRecordLimits bound any Span
// created by a Tracer from the LoggerProvider.
//
// If any field of sl is zero or negative it will be replaced with the default
// value for that field.
//
// If this or WithRawLogRecordLimits are not provided, the LoggerProvider will use
// the limits defined by environment variables, or the defaults if unset.
// Refer to the NewLogRecordLimits documentation for information about this
// relationship.
//
// Deprecated: Use WithRawLogRecordLimits instead which allows setting unlimited
// and zero limits. This option will be kept until the next major version
// incremented release.
func WithLogRecordLimits(sl LogRecordLimits) LoggerProviderOption {
	if sl.AttributeValueLengthLimit <= 0 {
		sl.AttributeValueLengthLimit = DefaultAttributeValueLengthLimit
	}
	if sl.AttributeCountLimit <= 0 {
		sl.AttributeCountLimit = DefaultAttributeCountLimit
	}
	return loggerProviderOptionFunc(
		func(cfg loggerProviderConfig) loggerProviderConfig {
			cfg.logRecordLimits = sl
			return cfg
		},
	)
}

// WithRawLogRecordLimits returns a LoggerProviderOption that configures a
// LoggerProvider to use these limits. These limits bound any Span created by
// a Tracer from the LoggerProvider.
//
// The limits will be used as-is. Zero or negative values will not be changed
// to the default value like WithLogRecordLimits does. Setting a limit to zero will
// effectively disable the related resource it limits and setting to a
// negative value will mean that resource is unlimited. Consequentially, this
// means that the zero-value LogRecordLimits will disable all span resources.
// Because of this, limits should be constructed using NewLogRecordLimits and
// updated accordingly.
//
// If this or WithLogRecordLimits are not provided, the LoggerProvider will use the
// limits defined by environment variables, or the defaults if unset. Refer to
// the NewLogRecordLimits documentation for information about this relationship.
func WithRawLogRecordLimits(limits LogRecordLimits) LoggerProviderOption {
	return loggerProviderOptionFunc(
		func(cfg loggerProviderConfig) loggerProviderConfig {
			cfg.logRecordLimits = limits
			return cfg
		},
	)
}

// ensureValidLoggerProviderConfig ensures that given TracerProviderConfig is valid.
func ensureValidLoggerProviderConfig(cfg loggerProviderConfig) loggerProviderConfig {
	if cfg.resource == nil {
		cfg.resource = resource.Default()
	}
	return cfg
}
