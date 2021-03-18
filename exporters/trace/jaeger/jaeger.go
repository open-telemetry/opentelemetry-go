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

package jaeger // import "go.opentelemetry.io/otel/exporters/trace/jaeger"

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"google.golang.org/api/support/bundler"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	gen "go.opentelemetry.io/otel/exporters/trace/jaeger/internal/gen-go/jaeger"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	keyInstrumentationLibraryName    = "otel.library.name"
	keyInstrumentationLibraryVersion = "otel.library.version"
)

type Option func(*options)

// options are the options to be used when initializing a Jaeger export.
type options struct {
	// Process contains the information about the exporting process.
	Process Process

	// BufferMaxCount defines the total number of traces that can be buffered in memory
	BufferMaxCount int

	// BatchMaxCount defines the maximum number of spans sent in one batch
	BatchMaxCount int

	// TracerProviderOptions defines the options for tracer provider of sdk.
	TracerProviderOptions []sdktrace.TracerProviderOption

	Disabled bool
}

// WithBufferMaxCount defines the total number of traces that can be buffered in memory
func WithBufferMaxCount(bufferMaxCount int) Option {
	return func(o *options) {
		o.BufferMaxCount = bufferMaxCount
	}
}

// WithBatchMaxCount defines the maximum number of spans in one batch
func WithBatchMaxCount(batchMaxCount int) Option {
	return func(o *options) {
		o.BatchMaxCount = batchMaxCount
	}
}

// WithSDKOptions configures options for tracer provider of sdk.
func WithSDKOptions(opts ...sdktrace.TracerProviderOption) Option {
	return func(o *options) {
		o.TracerProviderOptions = opts
	}
}

// WithDisabled option will cause pipeline methods to use
// a no-op provider
func WithDisabled(disabled bool) Option {
	return func(o *options) {
		o.Disabled = disabled
	}
}

// NewRawExporter returns an OTel Exporter implementation that exports the
// collected spans to Jaeger.
//
// It will IGNORE Disabled option.
func NewRawExporter(endpointOption EndpointOption, opts ...Option) (*Exporter, error) {
	uploader, err := endpointOption()
	if err != nil {
		return nil, err
	}

	o := options{}
	opts = append(opts, WithProcessFromEnv())
	for _, opt := range opts {
		opt(&o)
	}

	// Fetch default service.name from default resource for backup
	var defaultServiceName string
	defaultResource := resource.Default()
	if value, exists := defaultResource.Set().Value(semconv.ServiceNameKey); exists {
		defaultServiceName = value.AsString()
	}
	if defaultServiceName == "" {
		return nil, fmt.Errorf("failed to get service name from default resource")
	}

	e := &Exporter{
		uploader:            uploader,
		o:                   o,
		defaultServiceName:  defaultServiceName,
		resourceFromProcess: processToResource(o.Process),
	}
	bundler := bundler.NewBundler((*export.SpanSnapshot)(nil), func(bundle interface{}) {
		if err := e.upload(bundle.([]*export.SpanSnapshot)); err != nil {
			otel.Handle(err)
		}
	})

	// Set BufferedByteLimit with the total number of spans that are permissible to be held in memory.
	// This needs to be done since the size of messages is always set to 1. Failing to set this would allow
	// 1G messages to be held in memory since that is the default value of BufferedByteLimit.
	if o.BufferMaxCount != 0 {
		bundler.BufferedByteLimit = o.BufferMaxCount
	}

	// The default value bundler uses is 10, increase to send larger batches
	if o.BatchMaxCount != 0 {
		bundler.BundleCountThreshold = o.BatchMaxCount
	}

	e.bundler = bundler
	return e, nil
}

// NewExportPipeline sets up a complete export pipeline
// with the recommended setup for trace provider
func NewExportPipeline(endpointOption EndpointOption, opts ...Option) (trace.TracerProvider, func(), error) {
	o := options{}
	opts = append(opts, WithDisabledFromEnv())
	for _, opt := range opts {
		opt(&o)
	}
	if o.Disabled {
		return trace.NewNoopTracerProvider(), func() {}, nil
	}

	exporter, err := NewRawExporter(endpointOption, opts...)
	if err != nil {
		return nil, nil, err
	}

	pOpts := append(exporter.o.TracerProviderOptions, sdktrace.WithSyncer(exporter))
	tp := sdktrace.NewTracerProvider(pOpts...)
	return tp, exporter.Flush, nil
}

// InstallNewPipeline instantiates a NewExportPipeline with the
// recommended configuration and registers it globally.
func InstallNewPipeline(endpointOption EndpointOption, opts ...Option) (func(), error) {
	tp, flushFn, err := NewExportPipeline(endpointOption, opts...)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	return flushFn, nil
}

// Process contains the information exported to jaeger about the source
// of the trace data.
type Process struct {
	// ServiceName is the Jaeger service name.
	ServiceName string

	// Tags are added to Jaeger Process exports
	Tags []attribute.KeyValue
}

// Exporter is an implementation of an OTel SpanSyncer that uploads spans to
// Jaeger.
type Exporter struct {
	bundler  *bundler.Bundler
	uploader batchUploader
	o        options

	stoppedMu sync.RWMutex
	stopped   bool

	defaultServiceName  string
	resourceFromProcess *resource.Resource
}

var _ export.SpanExporter = (*Exporter)(nil)

// ExportSpans exports SpanSnapshots to Jaeger.
func (e *Exporter) ExportSpans(ctx context.Context, ss []*export.SpanSnapshot) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	for _, span := range ss {
		// TODO(jbd): Handle oversized bundlers.
		err := e.bundler.Add(span, 1)
		if err != nil {
			return fmt.Errorf("failed to bundle %q: %w", span.Name, err)
		}
	}
	return nil
}

// flush is used to wrap the bundler's Flush method for testing.
var flush = func(e *Exporter) {
	e.bundler.Flush()
}

// Shutdown stops the exporter flushing any pending exports.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	done := make(chan struct{}, 1)
	// Shadow so if the goroutine is leaked in testing it doesn't cause a race
	// condition when the file level var is reset.
	go func(FlushFunc func(*Exporter)) {
		// The OpenTelemetry specification is explicit in not having this
		// method block so the preference here is to orphan this goroutine if
		// the context is canceled or times out while this flushing is
		// occurring. This is a consequence of the bundler Flush method not
		// supporting a context.
		FlushFunc(e)
		done <- struct{}{}
	}(flush)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	return nil
}

func spanSnapshotToThrift(ss *export.SpanSnapshot) *gen.Span {
	tags := make([]*gen.Tag, 0, len(ss.Attributes))
	for _, kv := range ss.Attributes {
		tag := keyValueToTag(kv)
		if tag != nil {
			tags = append(tags, tag)
		}
	}

	if il := ss.InstrumentationLibrary; il.Name != "" {
		tags = append(tags, getStringTag(keyInstrumentationLibraryName, il.Name))
		if il.Version != "" {
			tags = append(tags, getStringTag(keyInstrumentationLibraryVersion, il.Version))
		}
	}

	if ss.SpanKind != trace.SpanKindInternal {
		tags = append(tags,
			getStringTag("span.kind", ss.SpanKind.String()),
		)
	}

	if ss.StatusCode != codes.Unset {
		tags = append(tags,
			getInt64Tag("status.code", int64(ss.StatusCode)),
			getStringTag("status.message", ss.StatusMessage),
		)

		if ss.StatusCode == codes.Error {
			tags = append(tags, getBoolTag("error", true))
		}
	}

	var logs []*gen.Log
	for _, a := range ss.MessageEvents {
		fields := make([]*gen.Tag, 0, len(a.Attributes))
		for _, kv := range a.Attributes {
			tag := keyValueToTag(kv)
			if tag != nil {
				fields = append(fields, tag)
			}
		}
		fields = append(fields, getStringTag("name", a.Name))
		logs = append(logs, &gen.Log{
			Timestamp: a.Time.UnixNano() / 1000,
			Fields:    fields,
		})
	}

	var refs []*gen.SpanRef
	for _, link := range ss.Links {
		tid := link.TraceID()
		sid := link.SpanID()
		refs = append(refs, &gen.SpanRef{
			TraceIdHigh: int64(binary.BigEndian.Uint64(tid[0:8])),
			TraceIdLow:  int64(binary.BigEndian.Uint64(tid[8:16])),
			SpanId:      int64(binary.BigEndian.Uint64(sid[:])),
			RefType:     gen.SpanRefType_FOLLOWS_FROM,
		})
	}

	tid := ss.SpanContext.TraceID()
	sid := ss.SpanContext.SpanID()
	return &gen.Span{
		TraceIdHigh:   int64(binary.BigEndian.Uint64(tid[0:8])),
		TraceIdLow:    int64(binary.BigEndian.Uint64(tid[8:16])),
		SpanId:        int64(binary.BigEndian.Uint64(sid[:])),
		ParentSpanId:  int64(binary.BigEndian.Uint64(ss.ParentSpanID[:])),
		OperationName: ss.Name, // TODO: if span kind is added then add prefix "Sent"/"Recv"
		Flags:         int32(ss.SpanContext.TraceFlags()),
		StartTime:     ss.StartTime.UnixNano() / 1000,
		Duration:      ss.EndTime.Sub(ss.StartTime).Nanoseconds() / 1000,
		Tags:          tags,
		Logs:          logs,
		References:    refs,
	}
}

func keyValueToTag(keyValue attribute.KeyValue) *gen.Tag {
	var tag *gen.Tag
	switch keyValue.Value.Type() {
	case attribute.STRING:
		s := keyValue.Value.AsString()
		tag = &gen.Tag{
			Key:   string(keyValue.Key),
			VStr:  &s,
			VType: gen.TagType_STRING,
		}
	case attribute.BOOL:
		b := keyValue.Value.AsBool()
		tag = &gen.Tag{
			Key:   string(keyValue.Key),
			VBool: &b,
			VType: gen.TagType_BOOL,
		}
	case attribute.INT64:
		i := keyValue.Value.AsInt64()
		tag = &gen.Tag{
			Key:   string(keyValue.Key),
			VLong: &i,
			VType: gen.TagType_LONG,
		}
	case attribute.FLOAT64:
		f := keyValue.Value.AsFloat64()
		tag = &gen.Tag{
			Key:     string(keyValue.Key),
			VDouble: &f,
			VType:   gen.TagType_DOUBLE,
		}
	case attribute.ARRAY:
		json, _ := json.Marshal(keyValue.Value.AsArray())
		a := (string)(json)
		tag = &gen.Tag{
			Key:   string(keyValue.Key),
			VStr:  &a,
			VType: gen.TagType_STRING,
		}
	}
	return tag
}

func getInt64Tag(k string, i int64) *gen.Tag {
	return &gen.Tag{
		Key:   k,
		VLong: &i,
		VType: gen.TagType_LONG,
	}
}

func getStringTag(k, s string) *gen.Tag {
	return &gen.Tag{
		Key:   k,
		VStr:  &s,
		VType: gen.TagType_STRING,
	}
}

func getBoolTag(k string, b bool) *gen.Tag {
	return &gen.Tag{
		Key:   k,
		VBool: &b,
		VType: gen.TagType_BOOL,
	}
}

// Flush waits for exported trace spans to be uploaded.
//
// This is useful if your program is ending and you do not want to lose recent spans.
func (e *Exporter) Flush() {
	flush(e)
}

func (e *Exporter) upload(spans []*export.SpanSnapshot) error {
	batchList := jaegerBatchList(spans, e.defaultServiceName, e.resourceFromProcess)
	for _, batch := range batchList {
		err := e.uploader.upload(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

// jaegerBatchList transforms a slice of SpanSnapshot into a slice of jaeger
// Batch.
func jaegerBatchList(ssl []*export.SpanSnapshot, defaultServiceName string, resourceFromProcess *resource.Resource) []*gen.Batch {
	if len(ssl) == 0 {
		return nil
	}

	batchDict := make(map[attribute.Distinct]*gen.Batch)

	for _, ss := range ssl {
		if ss == nil {
			continue
		}

		newResource := ss.Resource
		if resourceFromProcess != nil {
			// The value from process will overwrite the value from span's resources
			newResource = resource.Merge(ss.Resource, resourceFromProcess)
		}
		resourceKey := newResource.Equivalent()
		batch, bOK := batchDict[resourceKey]
		if !bOK {
			batch = &gen.Batch{
				Process: process(newResource, defaultServiceName),
				Spans:   []*gen.Span{},
			}
		}
		batch.Spans = append(batch.Spans, spanSnapshotToThrift(ss))
		batchDict[resourceKey] = batch
	}

	// Transform the categorized map into a slice
	batchList := make([]*gen.Batch, 0, len(batchDict))
	for _, batch := range batchDict {
		batchList = append(batchList, batch)
	}
	return batchList
}

// process transforms an OTel Resource into a jaeger Process.
func process(res *resource.Resource, defaultServiceName string) *gen.Process {
	var process gen.Process

	var serviceName attribute.KeyValue
	if res != nil {
		for iter := res.Iter(); iter.Next(); {
			if iter.Attribute().Key == semconv.ServiceNameKey {
				serviceName = iter.Attribute()
				// Don't convert service.name into tag.
				continue
			}
			if tag := keyValueToTag(iter.Attribute()); tag != nil {
				process.Tags = append(process.Tags, tag)
			}
		}
	}

	// If no service.name is contained in a Span's Resource,
	// that field MUST be populated from the default Resource.
	if serviceName.Value.AsString() == "" {
		serviceName = semconv.ServiceVersionKey.String(defaultServiceName)
	}
	process.ServiceName = serviceName.Value.AsString()

	return &process
}

func processToResource(process Process) *resource.Resource {
	var attrs []attribute.KeyValue
	if process.ServiceName != "" {
		attrs = append(attrs, semconv.ServiceNameKey.String(process.ServiceName))
	}
	attrs = append(attrs, process.Tags...)

	if len(attrs) == 0 {
		return nil
	}
	return resource.NewWithAttributes(attrs...)
}
