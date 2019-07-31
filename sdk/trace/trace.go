// Copyright 2019, OpenTelemetry Authors
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

package trace

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/event"
	apitag "go.opentelemetry.io/api/tag"
	apitrace "go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/sdk/internal"
	"google.golang.org/grpc/codes"
)

type tracer struct {
	name      string
	component string
	resources []core.KeyValue
}

var _ apitrace.Tracer = &tracer{}

func (tr *tracer) Start(ctx context.Context, name string, o ...apitrace.SpanOption) (context.Context, apitrace.Span) {
	var opts apitrace.SpanOptions
	var parent core.SpanContext
	var remoteParent bool

	//TODO [rghetia] : Add new option for parent. If parent is configured then use that parent.
	for _, op := range o {
		op(&opts)
	}

	// TODO: [rghetia] ChildOfRelationship is used to indicate that the parent is remote
	// and its context is received as part of a request. There are two possibilities
	// 1. Remote is trusted. So continue using same trace.
	//      tracer.Start(ctx, "some name", ChildOf(remote_span_context))
	// 2. Remote is not trusted. In this case create a root span and then add the remote as link
	//      span := tracer.Start(ctx, "some name")
	//      span.Link(remote_span_context, ChildOfRelationship)
	if opts.Reference.SpanContext != core.INVALID_SPAN_CONTEXT &&
		opts.Reference.RelationshipType == apitrace.ChildOfRelationship {
		parent = opts.Reference.SpanContext
		remoteParent = true
	} else {
		if p := FromContext(ctx); p != nil {
			p.addChild()
			parent = p.spanContext
		}
	}

	span := startSpanInternal(name, parent, remoteParent, opts)
	span.tracer = tr

	ctx, end := startExecutionTracerTask(ctx, name)
	span.executionTracerTaskEnd = end
	return NewContext(ctx, span), span
}

func (tr *tracer) WithSpan(ctx context.Context, name string, body func(ctx context.Context) error) error {
	ctx, span := tr.Start(ctx, name)
	defer span.Finish()

	if err := body(ctx); err != nil {
		// TODO: set event with boolean attribute for error.
		return err
	}
	return nil
}

// WithSampler sets the sampler to make sampling decision.
// If not set probability sampler with probability of 0.00001 is selected.
func WithSampler(s Sampler) apitrace.SpanOption {
	return func(o *apitrace.SpanOptions) {
		//TODO [rghetia] ; add sampler option
	}
}

// WithParent sets the SpanContext as parent SpanContext.
// It is typically used when a request is received from remote entity such as
// http client.
func WithParent(sc core.SpanContext) apitrace.SpanOption {
	return func(o *apitrace.SpanOptions) {
		//TODO [rghetia] ; add parent option
	}
}

func (tr *tracer) WithService(name string) apitrace.Tracer {
	tr.name = name
	return tr
}

// WithResources does nothing and returns noop implementation of apitrace.Tracer.
func (tr *tracer) WithResources(res ...core.KeyValue) apitrace.Tracer {
	tr.resources = res
	return tr
}

// WithComponent does nothing and returns noop implementation of apitrace.Tracer.
func (tr *tracer) WithComponent(component string) apitrace.Tracer {
	tr.component = component
	return tr
}

func (tr *tracer) Inject(ctx context.Context, span apitrace.Span, injector apitrace.Injector) {
	injector.Inject(span.SpanContext(), nil)
}

var tr *tracer
var registerOnce sync.Once

func Register() apitrace.Tracer {
	registerOnce.Do(func() {
		tr = &tracer{}
		apitrace.SetGlobalTracer(tr)
	})
	return tr
}

// sdkSpan represents a span of a trace.  It has an associated SpanContext, and
// stores data accumulated while the span is active.
//
// Ideally users should interact with Spans by calling the functions in this
// package that take a Context parameter.
type sdkSpan struct {
	// data contains information recorded about the span.
	//
	// It will be non-nil if we are exporting the span or recording events for it.
	// Otherwise, data is nil, and the sdkSpan is simply a carrier for the
	// SpanContext, so that the trace ID is propagated.
	data        *SpanData
	mu          sync.Mutex // protects the contents of *data (but not the pointer value.)
	spanContext core.SpanContext

	// lruAttributes are capped at configured limit. When the capacity is reached an oldest entry
	// is removed to create room for a new entry.
	lruAttributes *lruMap

	// messageEvents are stored in FIFO queue capped by configured limit.
	messageEvents *evictedQueue

	// links are stored in FIFO queue capped by configured limit.
	links *evictedQueue

	// spanStore is the spanStore this span belongs to, if any, otherwise it is nil.
	//*spanStore
	endOnce sync.Once

	executionTracerTaskEnd func()          // ends the execution tracer span
	tracer                 apitrace.Tracer // tracer used to create span.
}

var _ apitrace.Span = &sdkSpan{}

type contextKey struct{}

// FromContext returns the sdkSpan stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *sdkSpan {
	s, _ := ctx.Value(contextKey{}).(*sdkSpan)
	return s
}

// NewContext returns a new context with the given sdkSpan attached.
func NewContext(parent context.Context, s *sdkSpan) context.Context {
	return context.WithValue(parent, contextKey{}, s)
}

// All available span kinds. sdkSpan kind must be either one of these values.
const (
	SpanKindUnspecified = iota
	SpanKindServer
	SpanKindClient
)

// SpancContext returns an invalid span context.
func (s *sdkSpan) SpanContext() core.SpanContext {
	if s == nil {
		return core.SpanContext{}
	}
	return s.spanContext
}

// IsRecordingEvents returns true if events are being recorded for this span.
// Use this check to avoid computing expensive annotations when they will never
// be used.
func (s *sdkSpan) IsRecordingEvents() bool {
	if s == nil {
		return false
	}
	return s.data != nil
}

// SetStatus sets the status of the span, if it is recording events.
func (s *sdkSpan) SetStatus(status codes.Code) {
	if s == nil {
		return
	}
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.data.Status = status
	s.mu.Unlock()
}

// SetError does nothing.
func (s *sdkSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (s *sdkSpan) SetAttribute(attribute core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.copyToCappedAttributes(attribute)
	s.mu.Unlock()
}

// SetAttributes does nothing.
func (s *sdkSpan) SetAttributes(attributes ...core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.copyToCappedAttributes(attributes...)
	s.mu.Unlock()
}

// ModifyAttribute does nothing.
func (s *sdkSpan) ModifyAttribute(mutator apitag.Mutator) {
}

// ModifyAttributes does nothing.
func (s *sdkSpan) ModifyAttributes(mutators ...apitag.Mutator) {
}

// Finish does nothing.
func (s *sdkSpan) Finish() {
	if s == nil {
		return
	}

	if s.executionTracerTaskEnd != nil {
		s.executionTracerTaskEnd()
	}
	if !s.IsRecordingEvents() {
		return
	}
	s.endOnce.Do(func() {
		exp, _ := exporters.Load().(exportersMap)
		mustExport := s.spanContext.IsSampled() && len(exp) > 0
		//if s.spanStore != nil || mustExport {
		if mustExport {
			sd := s.makeSpanData()
			sd.EndTime = internal.MonotonicEndTime(sd.StartTime)
			//if s.spanStore != nil {
			//	s.spanStore.finished(s, sd)
			//}
			if mustExport {
				for e := range exp {
					e.ExportSpan(sd)
				}
			}
		}
	})
}

// Tracer returns noop implementation of Tracer.
func (s *sdkSpan) Tracer() apitrace.Tracer {
	return s.tracer
}

func (s *sdkSpan) AddEvent(ctx context.Context, event event.Event) {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageEvents.add(event)
}

func (s *sdkSpan) Event(ctx context.Context, msg string, attrs ...core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	now := time.Now()
	s.mu.Lock()
	s.messageEvents.add(MessageEvent{
		msg:        msg,
		attributes: attrs,
		Time:       now,
	})
	s.mu.Unlock()
}

func startSpanInternal(name string, parent core.SpanContext, remoteParent bool, o apitrace.SpanOptions) *sdkSpan {
	var noParent bool
	span := &sdkSpan{}
	span.spanContext = parent

	cfg := config.Load().(*Config)

	if parent == core.INVALID_SPAN_CONTEXT {
		span.spanContext.TraceID = cfg.IDGenerator.NewTraceID()
		noParent = true
	}
	span.spanContext.SpanID = cfg.IDGenerator.NewSpanID()
	sampler := cfg.DefaultSampler

	// TODO: [rghetia] fix sampler
	//if !hasParent || remoteParent || o.Sampler != nil {
	if noParent || remoteParent {
		// If this span is the child of a local span and no Sampler is set in the
		// options, keep the parent's TraceOptions.
		//
		// Otherwise, consult the Sampler in the options if it is non-nil, otherwise
		// the default sampler.
		//if o.Sampler != nil {
		//	sampler = o.Sampler
		//}
		sampled := sampler(SamplingParameters{
			ParentContext:   parent,
			TraceID:         span.spanContext.TraceID,
			SpanID:          span.spanContext.SpanID,
			Name:            name,
			HasRemoteParent: remoteParent}).Sample
		if sampled {
			span.spanContext.TraceOptions = core.TraceOptionSampled
		}
	}

	// TODO: [rghetia] restore when spanstore is added.
	// if !internal.LocalSpanStoreEnabled && !span.spanContext.IsSampled() && !o.RecordEvent {
	if !span.spanContext.IsSampled() && !o.RecordEvent {
		return span
	}

	span.data = &SpanData{
		SpanContext: span.spanContext,
		StartTime:   time.Now(),
		// TODO;[rghetia] : fix spanKind
		//SpanKind:        o.SpanKind,
		Name:            name,
		HasRemoteParent: remoteParent,
	}
	span.lruAttributes = newLruMap(cfg.MaxAttributesPerSpan)
	span.messageEvents = newEvictedQueue(cfg.MaxMessageEventsPerSpan)
	span.links = newEvictedQueue(cfg.MaxLinksPerSpan)

	if !noParent {
		span.data.ParentSpanID = parent.SpanID
	}
	// TODO: [rghetia] restore when spanstore is added.
	//if internal.LocalSpanStoreEnabled {
	//	ss := spanStoreForNameCreateIfNew(name)
	//	if ss != nil {
	//		span.spanStore = ss
	//		ss.add(span)
	//	}
	//}

	return span
}

// makeSpanData produces a SpanData representing the current state of the sdkSpan.
// It requires that s.data is non-nil.
func (s *sdkSpan) makeSpanData() *SpanData {
	var sd SpanData
	s.mu.Lock()
	sd = *s.data
	if s.lruAttributes.simpleLruMap.Len() > 0 {
		sd.Attributes = s.lruAttributesToAttributeMap()
		sd.DroppedAttributeCount = s.lruAttributes.droppedCount
	}
	if len(s.messageEvents.queue) > 0 {
		sd.MessageEvents = s.interfaceArrayToMessageEventArray()
		sd.DroppedMessageEventCount = s.messageEvents.droppedCount
	}
	if len(s.links.queue) > 0 {
		sd.Links = s.interfaceArrayToLinksArray()
		sd.DroppedLinkCount = s.links.droppedCount
	}
	s.mu.Unlock()
	return &sd
}

func (s *sdkSpan) interfaceArrayToLinksArray() []Link {
	linksArr := make([]Link, 0)
	for _, value := range s.links.queue {
		linksArr = append(linksArr, value.(Link))
	}
	return linksArr
}

func (s *sdkSpan) interfaceArrayToMessageEventArray() []MessageEvent {
	messageEventArr := make([]MessageEvent, 0)
	for _, value := range s.messageEvents.queue {
		messageEventArr = append(messageEventArr, value.(MessageEvent))
	}
	return messageEventArr
}

func (s *sdkSpan) lruAttributesToAttributeMap() map[string]interface{} {
	attributes := make(map[string]interface{})
	for _, key := range s.lruAttributes.simpleLruMap.Keys() {
		value, ok := s.lruAttributes.simpleLruMap.Get(key)
		if ok {
			key := key.(core.Key)
			attributes[key.Variable.Name] = value
		}
	}
	return attributes
}

func (s *sdkSpan) copyToCappedAttributes(attributes ...core.KeyValue) {
	for _, a := range attributes {
		s.lruAttributes.add(a.Key, a.Value)
	}
}

func (s *sdkSpan) addChild() {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.data.ChildSpanCount++
	s.mu.Unlock()
}

// AddLink adds a link to the span.
func (s *sdkSpan) AddLink(l Link) {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.links.add(l)
	s.mu.Unlock()
}

func (s *sdkSpan) String() string {
	if s == nil {
		return "<nil>"
	}
	if s.data == nil {
		return fmt.Sprintf("span %s", s.spanContext.SpanIDString())
	}
	s.mu.Lock()
	str := fmt.Sprintf("span %s %q", s.spanContext.SpanIDString(), s.data.Name)
	s.mu.Unlock()
	return str
}

var config atomic.Value // access atomically

func init() {
	gen := &defaultIDGenerator{}
	// initialize traceID and spanID generators.
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &gen.traceIDAdd, &gen.nextSpanID, &gen.spanIDInc,
	} {
		_ = binary.Read(crand.Reader, binary.LittleEndian, p)
	}
	gen.traceIDRand = rand.New(rand.NewSource(rngSeed))
	gen.spanIDInc |= 1

	config.Store(&Config{
		DefaultSampler:          ProbabilitySampler(defaultSamplingProbability),
		IDGenerator:             gen,
		MaxAttributesPerSpan:    DefaultMaxAttributesPerSpan,
		MaxMessageEventsPerSpan: DefaultMaxMessageEventsPerSpan,
		MaxLinksPerSpan:         DefaultMaxLinksPerSpan,
	})
}

type defaultIDGenerator struct {
	sync.Mutex

	// Please keep these as the first fields
	// so that these 8 byte fields will be aligned on addresses
	// divisible by 8, on both 32-bit and 64-bit machines when
	// performing atomic increments and accesses.
	// See:
	// * https://github.com/census-instrumentation/opencensus-go/issues/587
	// * https://github.com/census-instrumentation/opencensus-go/issues/865
	// * https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	nextSpanID uint64
	spanIDInc  uint64

	traceIDAdd  [2]uint64
	traceIDRand *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *defaultIDGenerator) NewSpanID() uint64 {
	var id uint64
	for id == 0 {
		id = atomic.AddUint64(&gen.nextSpanID, gen.spanIDInc)
	}
	return id
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// mu should be held while this function is called.
func (gen *defaultIDGenerator) NewTraceID() core.TraceID {
	gen.Lock()
	// Construct the trace ID from two outputs of traceIDRand, with a constant
	// added to each half for additional entropy.
	tid := core.TraceID{
		High: gen.traceIDRand.Uint64() + gen.traceIDAdd[0],
		Low:  gen.traceIDRand.Uint64() + gen.traceIDAdd[1],
	}
	gen.Unlock()
	return tid
}
