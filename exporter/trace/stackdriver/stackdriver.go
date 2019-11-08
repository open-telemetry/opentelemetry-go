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

package stackdriver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	traceapi "cloud.google.com/go/trace/apiv2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

// Option is function type that is passed to the exporter initialization function.
type Option func(*options)

// options contains options for configuring the exporter.
type options struct {
	// ProjectID is the identifier of the Stackdriver
	// project the user is uploading the stats data to.
	// If not set, this will default to your "Application Default Credentials".
	// For details see: https://developers.google.com/accounts/docs/application-default-credentials.
	//
	// It will be used in the project_id label of a Stackdriver monitored
	// resource if the resource does not inherently belong to a specific
	// project, e.g. on-premise resource like k8s_container or generic_task.
	ProjectID string

	// Location is the identifier of the GCP or AWS cloud region/zone in which
	// the data for a resource is stored.
	// If not set, it will default to the location provided by the metadata server.
	//
	// It will be used in the location label of a Stackdriver monitored resource
	// if the resource does not inherently belong to a specific project, e.g.
	// on-premise resource like k8s_container or generic_task.
	Location string

	// OnError is the hook to be called when there is
	// an error uploading the stats or tracing data.
	// If no custom hook is set, errors are logged.
	// Optional.
	OnError func(err error)

	// MonitoringClientOptions are additional options to be passed
	// to the underlying Stackdriver Monitoring API client.
	// Optional.
	MonitoringClientOptions []option.ClientOption

	// TraceClientOptions are additional options to be passed
	// to the underlying Stackdriver Trace API client.
	// Optional.
	TraceClientOptions []option.ClientOption

	// TraceSpansBufferMaxBytes is the maximum size (in bytes) of spans that
	// will be buffered in memory before being dropped.
	//
	// If unset, a default of 8MB will be used.
	// TraceSpansBufferMaxBytes int

	// DefaultTraceAttributes will be appended to every span that is exported to
	// Stackdriver Trace.
	DefaultTraceAttributes map[string]interface{}

	// Context allows you to provide a custom context for API calls.
	//
	// This context will be used several times: first, to create Stackdriver
	// trace and metric clients, and then every time a new batch of traces or
	// stats needs to be uploaded.
	//
	// Do not set a timeout on this context. Instead, set the Timeout option.
	//
	// If unset, context.Background() will be used.
	Context context.Context

	// SkipCMD enforces to skip all the CreateMetricDescriptor calls.
	// These calls are important in order to configure the unit of the metrics,
	// but in some cases all the exported metrics are builtin (unit is configured)
	// or the unit is not important.
	SkipCMD bool

	// Timeout for all API calls. If not set, defaults to 5 seconds.
	Timeout time.Duration

	// ReportingInterval sets the interval between reporting metrics.
	// If it is set to zero then default value is used.
	ReportingInterval time.Duration

	// NumberOfWorkers sets the number of go rountines that send requests
	// to Stackdriver Monitoring. This is only used for Proto metrics export
	// for now. The minimum number of workers is 1.
	NumberOfWorkers int
}

// WithProjectID sets Google Cloud Platform project as projectID.
// Without using this option, it automatically detects the project ID
// from the default credential detection process.
// Please find the detailed order of the default credentail detection proecess on the doc:
// https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials
func WithProjectID(projectID string) func(o *options) {
	return func(o *options) {
		o.ProjectID = projectID
	}
}

// WithOnError sets the hook to be called when there is an error
// occurred on uploading the span data to Stackdriver.
// If no custom hook is set, errors are logged.
func WithOnError(onError func(err error)) func(o *options) {
	return func(o *options) {
		o.OnError = onError
	}
}

// WithTraceClientOptions sets additionial client options for tracing.
func WithTraceClientOptions(opts []option.ClientOption) func(o *options) {
	return func(o *options) {
		o.TraceClientOptions = opts
	}
}

// WithContext sets the context that trace exporter and metric exporter
// relies on.
func WithContext(ctx context.Context) func(o *options) {
	return func(o *options) {
		o.Context = ctx
	}
}

func (o *options) handleError(err error) {
	if o.OnError != nil {
		o.OnError(err)
		return
	}
	log.Printf("Failed to export to Stackdriver: %v", err)
}

// defaultTimeout is used as default when timeout is not set in newContextWithTimout.
const defaultTimeout = 5 * time.Second

// Exporter is a trace exporter that uploads data to Stackdriver.
//
// TODO(yoshifumi): add a metrics exporter once the spec definition
// process and the sampler implementation are done.
type Exporter struct {
	traceExporter *traceExporter
}

// NewExporter creates a new Exporter thats implements trace.Exporter.
//
// TODO(yoshifumi): add a metrics exporter one the spec definition
// process and the sampler implementation are done.
func NewExporter(opts ...Option) (*Exporter, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	if o.ProjectID == "" {
		ctx := o.Context
		if ctx == nil {
			ctx = context.Background()
		}
		creds, err := google.FindDefaultCredentials(ctx, traceapi.DefaultAuthScopes()...)
		if err != nil {
			return nil, fmt.Errorf("Stackdriver: %v", err)
		}
		if creds.ProjectID == "" {
			return nil, errors.New("Stackdriver: no project found with application default credentials")
		}
		o.ProjectID = creds.ProjectID
	}
	te, err := newTraceExporter(&o)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		traceExporter: te,
	}, nil
}

func newContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return context.WithTimeout(ctx, timeout)
}

// ExportSpan exports a SpanData to Stackdriver Trace.
func (e *Exporter) ExportSpan(ctx context.Context, sd *export.SpanData) {
	if len(e.traceExporter.o.DefaultTraceAttributes) > 0 {
		sd = e.sdWithDefaultTraceAttributes(sd)
	}
	e.traceExporter.ExportSpan(ctx, sd)
}

// ExportSpans exports a slice of SpanData to Stackdriver Trace in batch
func (e *Exporter) ExportSpans(ctx context.Context, sds []*export.SpanData) {
	e.traceExporter.ExportSpans(ctx, sds)
}

func (e *Exporter) sdWithDefaultTraceAttributes(sd *export.SpanData) *export.SpanData {
	newSD := *sd
	for k, v := range e.traceExporter.o.DefaultTraceAttributes {
		switch val := v.(type) {
		case bool:
			newSD.Attributes = append(newSD.Attributes, key.New(k).Bool(val))
		case int64:
			newSD.Attributes = append(newSD.Attributes, key.New(k).Int64(val))
		case float64:
			newSD.Attributes = append(newSD.Attributes, key.New(k).Float64(val))
		case string:
			newSD.Attributes = append(newSD.Attributes, key.New(k).String(val))
		}
	}
	newSD.Attributes = append(newSD.Attributes, sd.Attributes...)
	return &newSD
}
