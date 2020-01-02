package internal

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/api/trace"
)

type traceProvider struct {
	lock    sync.Mutex
	tracers []*tracer

	delegate trace.Provider
}

var _ trace.Provider = &traceProvider{}

func (p *traceProvider) setDelegate(provider trace.Provider) {
	if p.delegate != nil {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.delegate = provider
	for _, t := range p.tracers {
		t.setDelegate(provider)
	}

	p.tracers = nil
}

// Tracer creates a trace.Tracer with the given name. When a delegate is set,
// all previously returned trace.Tracers will be swapped to equivalent
// trace.Tracers created from the delegate.
func (p *traceProvider) Tracer(name string) trace.Tracer {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.delegate != nil {
		return p.delegate.Tracer(name)
	}

	t := &tracer{name: name}
	p.tracers = append(p.tracers, t)
	return t
}

type tracer struct {
	once sync.Once
	name string

	delegate trace.Tracer
}

var _ trace.Tracer = &tracer{}

func (t *tracer) setDelegate(provider trace.Provider) {
	t.once.Do(func() { t.delegate = provider.Tracer(t.name) })
}

// WithSpan wraps around execution of func with delegated Tracer.
func (t *tracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	if t.delegate != nil {
		return t.delegate.WithSpan(ctx, name, body)
	}
	return trace.NoopTracer{}.WithSpan(ctx, name, body)
}

// Start starts a span from the delegated tracer.
func (t *tracer) Start(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, trace.Span) {
	if t.delegate != nil {
		return t.delegate.Start(ctx, name, opts...)
	}
	return trace.NoopTracer{}.Start(ctx, name, opts...)
}
