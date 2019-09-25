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

package jaeger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/api/support/bundler"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	gen "go.opentelemetry.io/exporter/trace/jaeger/internal/gen-go/jaeger"
	"go.opentelemetry.io/sdk/trace"
)

const defaultServiceName = "OpenTelemetry"

// Options are the options to be used when initializing a Jaeger exporter.
type Options struct {
	// CollectorEndpoint is the full url to the Jaeger HTTP Thrift collector.
	// For example, http://localhost:14268/api/traces
	CollectorEndpoint string

	// AgentEndpoint instructs exporter to send spans to jaeger-agent at this address.
	// For example, localhost:6831.
	AgentEndpoint string

	// OnError is the hook to be called when there is
	// an error occurred when uploading the stats data.
	// If no custom hook is set, errors are logged.
	// Optional.
	OnError func(err error)

	// Username to be used if basic auth is required.
	// Optional.
	Username string

	// Password to be used if basic auth is required.
	// Optional.
	Password string

	// Process contains the information about the exporting process.
	Process Process

	//BufferMaxCount defines the total number of traces that can be buffered in memory
	BufferMaxCount int
}

// NewExporter returns a trace.Exporter implementation that exports
// the collected spans to Jaeger.
func NewExporter(o Options) (*Exporter, error) {
	if o.CollectorEndpoint == "" && o.AgentEndpoint == "" {
		return nil, errors.New("missing endpoint for Jaeger exporter")
	}

	var endpoint string
	var client *agentClientUDP
	var err error
	if o.CollectorEndpoint != "" {
		endpoint = o.CollectorEndpoint
	} else {
		client, err = newAgentClientUDP(o.AgentEndpoint, udpPacketMaxLength)
		if err != nil {
			return nil, err
		}
	}
	onError := func(err error) {
		if o.OnError != nil {
			o.OnError(err)
			return
		}
		log.Printf("Error when uploading spans to Jaeger: %v", err)
	}
	service := o.Process.ServiceName
	if service == "" {
		service = defaultServiceName
	}
	tags := make([]*gen.Tag, len(o.Process.Tags))
	for i, tag := range o.Process.Tags {
		tags[i] = attributeToTag(tag.key, tag.value)
	}
	e := &Exporter{
		endpoint:      endpoint,
		agentEndpoint: o.AgentEndpoint,
		client:        client,
		username:      o.Username,
		password:      o.Password,
		process: &gen.Process{
			ServiceName: service,
			Tags:        tags,
		},
	}
	bundler := bundler.NewBundler((*gen.Span)(nil), func(bundle interface{}) {
		if err := e.upload(bundle.([]*gen.Span)); err != nil {
			onError(err)
		}
	})

	// Set BufferedByteLimit with the total number of spans that are permissible to be held in memory.
	// This needs to be done since the size of messages is always set to 1. Failing to set this would allow
	// 1G messages to be held in memory since that is the default value of BufferedByteLimit.
	if o.BufferMaxCount != 0 {
		bundler.BufferedByteLimit = o.BufferMaxCount
	}

	e.bundler = bundler
	return e, nil
}

// Process contains the information exported to jaeger about the source
// of the trace data.
type Process struct {
	// ServiceName is the Jaeger service name.
	ServiceName string

	// Tags are added to Jaeger Process exports
	Tags []Tag
}

// Tag defines a key-value pair
// It is limited to the possible conversions to *jaeger.Tag by attributeToTag
type Tag struct {
	key   string
	value interface{}
}

// Exporter is an implementation of trace.Exporter that uploads spans to Jaeger.
type Exporter struct {
	endpoint      string
	agentEndpoint string
	process       *gen.Process
	bundler       *bundler.Bundler
	client        *agentClientUDP

	username, password string
}

var _ trace.Exporter = (*Exporter)(nil)

// ExportSpan exports a SpanData to Jaeger.
func (e *Exporter) ExportSpan(data *trace.SpanData) {
	_ = e.bundler.Add(spanDataToThrift(data), 1)
	// TODO(jbd): Handle oversized bundlers.
}

func spanDataToThrift(data *trace.SpanData) *gen.Span {
	tags := make([]*gen.Tag, 0, len(data.Attributes))
	for _, kv := range data.Attributes {
		tag := coreAttributeToTag(kv)
		tags = append(tags, tag)
	}

	tags = append(tags, getInt64Tag("status.code", int64(data.Status)),
		getStringTag("status.message", data.Status.String()),
	)

	// Ensure that if Status.Code is not OK, that we set the "error" tag on the Jaeger span.
	// See Issue https://github.com/census-instrumentation/opencensus-go/issues/1041
	if data.Status != codes.OK {
		tags = append(tags, getBoolTag("error", true))
	}

	var logs []*gen.Log
	for _, a := range data.MessageEvents {
		fields := make([]*gen.Tag, 0, len(a.Attributes))
		for _, kv := range a.Attributes {
			tag := coreAttributeToTag(kv)
			if tag != nil {
				fields = append(fields, tag)
			}
		}
		fields = append(fields, getStringTag("message", a.Message))
		logs = append(logs, &gen.Log{
			Timestamp: a.Time.UnixNano() / 1000,
			Fields:    fields,
		})
	}
	//TODO: [rghetia] add links.
	//
	//var refs []*gen.SpanRef
	//for _, link := range data.Links {
	//	refs = append(refs, &gen.SpanRef{
	//		TraceIdHigh: bytesToInt64(link.TraceID[0:8]),
	//		TraceIdLow:  bytesToInt64(link.TraceID[8:16]),
	//		SpanId:      bytesToInt64(link.SpanID[:]),
	//	})
	//}

	return &gen.Span{
		TraceIdHigh:   int64(data.SpanContext.TraceID.High),
		TraceIdLow:    int64(data.SpanContext.TraceID.Low),
		SpanId:        int64(data.SpanContext.SpanID),
		ParentSpanId:  int64(data.ParentSpanID),
		OperationName: data.Name, // TODO: if span kind is added then add prefix "Sent"/"Recv"
		Flags:         int32(data.SpanContext.TraceFlags),
		StartTime:     data.StartTime.UnixNano() / 1000,
		Duration:      data.EndTime.Sub(data.StartTime).Nanoseconds() / 1000,
		Tags:          tags,
		Logs:          logs,
		// TODO: goes with Links.
		// References:    refs,
	}
}

func coreAttributeToTag(kv core.KeyValue) *gen.Tag {
	var tag *gen.Tag
	switch kv.Value.Type {
	case core.STRING:
		tag = &gen.Tag{
			Key:   kv.Key.Name,
			VStr:  &kv.Value.String,
			VType: gen.TagType_STRING,
		}
	case core.BOOL:
		tag = &gen.Tag{
			Key:   kv.Key.Name,
			VBool: &kv.Value.Bool,
			VType: gen.TagType_BOOL,
		}
	case core.INT32, core.INT64:
		tag = &gen.Tag{
			Key:   kv.Key.Name,
			VLong: &kv.Value.Int64,
			VType: gen.TagType_LONG,
		}
	case core.FLOAT32, core.FLOAT64:
		tag = &gen.Tag{
			Key:     kv.Key.Name,
			VDouble: &kv.Value.Float64,
			VType:   gen.TagType_DOUBLE,
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

// TODO(rghetia): remove interface{}. see https://github.com/open-telemetry/opentelemetry-go/pull/112/files#r321444786
func attributeToTag(key string, a interface{}) *gen.Tag {
	var tag *gen.Tag
	switch value := a.(type) {
	case bool:
		tag = &gen.Tag{
			Key:   key,
			VBool: &value,
			VType: gen.TagType_BOOL,
		}
	case string:
		tag = &gen.Tag{
			Key:   key,
			VStr:  &value,
			VType: gen.TagType_STRING,
		}
	case int64:
		tag = &gen.Tag{
			Key:   key,
			VLong: &value,
			VType: gen.TagType_LONG,
		}
	case int32:
		v := int64(value)
		tag = &gen.Tag{
			Key:   key,
			VLong: &v,
			VType: gen.TagType_LONG,
		}
	case float64:
		v := float64(value)
		tag = &gen.Tag{
			Key:     key,
			VDouble: &v,
			VType:   gen.TagType_DOUBLE,
		}
	}
	return tag
}

// Flush waits for exported trace spans to be uploaded.
//
// This is useful if your program is ending and you do not want to lose recent spans.
func (e *Exporter) Flush() {
	e.bundler.Flush()
}

func (e *Exporter) upload(spans []*gen.Span) error {
	batch := &gen.Batch{
		Spans:   spans,
		Process: e.process,
	}
	if e.endpoint != "" {
		return e.uploadCollector(batch)
	}
	return e.uploadAgent(batch)
}

func (e *Exporter) uploadAgent(batch *gen.Batch) error {
	return e.client.EmitBatch(batch)
}

func (e *Exporter) uploadCollector(batch *gen.Batch) error {
	body, err := serialize(batch)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", e.endpoint, body)
	if err != nil {
		return err
	}
	if e.username != "" && e.password != "" {
		req.SetBasicAuth(e.username, e.password)
	}
	req.Header.Set("Content-Type", "application/x-thrift")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to upload traces; HTTP status code: %d", resp.StatusCode)
	}
	return nil
}

func serialize(obj thrift.TStruct) (*bytes.Buffer, error) {
	buf := thrift.NewTMemoryBuffer()
	if err := obj.Write(thrift.NewTBinaryProtocolTransport(buf)); err != nil {
		return nil, err
	}
	return buf.Buffer, nil
}
