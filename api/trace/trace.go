package trace

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/log"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type (
	span struct {
		tracer      *tracer
		spanContext core.SpanContext
		lock        sync.Mutex
		eventID     core.EventID
		finishOnce  sync.Once
	}

	tracer struct {
		resources core.EventID
	}
)

var (
	ServiceKey      = tag.New("service")
	ComponentKey    = tag.New("component")
	ErrorKey        = tag.New("error")
	SpanIDKey       = tag.New("span_id")
	TraceIDKey      = tag.New("trace_id")
	ParentSpanIDKey = tag.New("parent_span_id")
	MessageKey      = tag.New("message",
		tag.WithDescription("message text: info, error, etc"),
	)

	// The process global tracer could have process-wide resource
	// tags applied directly, or we can have a SetGlobal tracer to
	// install a default tracer w/ resources.
	global atomic.Value
	empty  = &tracer{}
)

func (t *tracer) ScopeID() core.ScopeID {
	return t.resources.Scope()
}

func (t *tracer) WithResources(attributes ...core.KeyValue) Tracer {
	s := scope.New(t.resources.Scope(), attributes...)
	return &tracer{
		resources: s.ScopeID().EventID,
	}
}

func (g *tracer) WithComponent(name string) Tracer {
	return g.WithResources(ComponentKey.String(name))
}

func (g *tracer) WithService(name string) Tracer {
	return g.WithResources(ServiceKey.String(name))
}

func (t *tracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	// TODO: use runtime/trace.WithRegion for execution tracer support
	// TODO: use runtime/pprof.Do for profile tags support
	ctx, span := t.Start(ctx, name)
	defer span.Finish()

	if err := body(ctx); err != nil {
		span.SetAttribute(ErrorKey.Bool(true))
		log.Log(ctx, "span error", MessageKey.String(err.Error()))
		return err
	}
	return nil
}

func (t *tracer) Start(ctx context.Context, name string, opts ...Option) (context.Context, Span) {
	var child core.SpanContext

	child.SpanID = rand.Uint64()

	var startTime time.Time
	var attributes []core.KeyValue
	var reference Reference

	for _, opt := range opts {
		if !opt.startTime.IsZero() {
			startTime = opt.startTime
		}
		if len(opt.attributes) != 0 {
			attributes = append(opt.attributes, attributes...)
		}
		if opt.attribute.Key != nil {
			attributes = append(attributes, opt.attribute)
		}
		if opt.reference.HasTraceID() {
			reference = opt.reference
		}
	}

	var parentScope core.ScopeID

	if reference.HasTraceID() {
		parentScope = reference.Scope()
	} else {
		parentScope = Active(ctx).ScopeID()
	}

	if parentScope.HasTraceID() {
		parent := parentScope.SpanContext
		child.TraceIDHigh = parent.TraceIDHigh
		child.TraceIDLow = parent.TraceIDLow
	} else {
		child.TraceIDHigh = rand.Uint64()
		child.TraceIDLow = rand.Uint64()
	}

	childScope := core.ScopeID{
		SpanContext: child,
		EventID:     t.resources,
	}

	span := &span{
		spanContext: child,
		tracer:      t,
		eventID: observer.Record(observer.Event{
			Time:    startTime,
			Type:    observer.START_SPAN,
			Scope:   scope.New(childScope, attributes...).ScopeID(),
			Context: ctx,
			Parent:  parentScope,
			String:  name,
		}),
	}
	return scope.SetActive(ctx, span), span
}

func (t *tracer) Inject(ctx context.Context, span Span, injector Injector) {
	injector.Inject(span.ScopeID().SpanContext, tag.FromContext(ctx))
}
