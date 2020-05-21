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

package otlp_testing

import (
	"context"
	"net"
	"sort"
	"sync"
	"time"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

// Let the system define the port
const defaultServerAddr = "127.0.0.1:0"

type ServerSuite struct {
	suite.Suite

	ServerOpts     []grpc.ServerOption
	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server

	ExporterOpts []otlp.ExporterOption
	Exporter     *otlp.Exporter
}

func (s *ServerSuite) SetupSuite() {
	s.serverAddr = defaultServerAddr

	var err error
	s.ServerListener, err = net.Listen("tcp", s.serverAddr)
	s.serverAddr = s.ServerListener.Addr().String()
	require.NoError(s.T(), err, "failed to allocate a port for server")

	s.Server = grpc.NewServer(s.ServerOpts...)
}

func (s *ServerSuite) StartServer() {
	go func() {
		s.Server.Serve(s.ServerListener)
	}()

	if s.Exporter == nil {
		s.Exporter = s.NewExporter()
	}
}

func (s *ServerSuite) ServerAddr() string {
	return s.serverAddr
}

func (s *ServerSuite) NewExporter() *otlp.Exporter {
	opts := []otlp.ExporterOption{
		otlp.WithInsecure(),
		otlp.WithAddress(s.serverAddr),
		otlp.WithReconnectionPeriod(10 * time.Millisecond),
	}
	exp, err := otlp.NewExporter(append(opts, s.ExporterOpts...)...)
	require.NoError(s.T(), err, "failed to create exporter")
	return exp
}

func (s *ServerSuite) TearDownSuite() {
	if s.ServerListener != nil {
		s.Server.GracefulStop()
		s.T().Logf("stopped grpc.Server at: %v", s.ServerAddr())
	}
	s.Exporter.Stop()
}

type TraceSuite struct {
	ServerSuite

	TraceProvider *trace.Provider

	mu              sync.RWMutex
	resourceSpanMap map[string]*tracepb.ResourceSpans
}

func (ts *TraceSuite) SetupSuite() {
	ts.ServerSuite.SetupSuite()

	ts.Reset()
	coltracepb.RegisterTraceServiceServer(ts.Server, ts)

	ts.ServerSuite.StartServer()

	if ts.TraceProvider == nil {
		ts.TraceProvider = ts.NewTraceProvider(ts.Exporter, nil)
	}
}

func (ts *TraceSuite) NewTraceProvider(exp *otlp.Exporter, opts []trace.ProviderOption) *trace.Provider {
	defaultOpts := []trace.ProviderOption{
		trace.WithConfig(
			trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			},
		),
		trace.WithSyncer(exp),
	}
	p, err := trace.NewProvider(append(defaultOpts, opts...)...)
	require.NoError(ts.T(), err, "failed to create trace provider")
	return p
}

func (ts *TraceSuite) Reset() {
	ts.resourceSpanMap = map[string]*tracepb.ResourceSpans{}
}

func (ts *TraceSuite) GetSpans() []*tracepb.Span {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	spans := []*tracepb.Span{}
	for _, rs := range ts.resourceSpanMap {
		for _, ils := range rs.InstrumentationLibrarySpans {
			spans = append(spans, ils.Spans...)
		}
	}
	return spans
}

func (ts *TraceSuite) GetResourceSpans() []*tracepb.ResourceSpans {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	rss := make([]*tracepb.ResourceSpans, 0, len(ts.resourceSpanMap))
	for _, rs := range ts.resourceSpanMap {
		rss = append(rss, rs)
	}
	return rss
}

func (ts *TraceSuite) Export(ctx context.Context, req *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	rss := req.GetResourceSpans()
	for _, rs := range rss {
		rstr := resourceString(rs.Resource)
		existingRs, ok := ts.resourceSpanMap[rstr]
		if !ok {
			ts.resourceSpanMap[rstr] = rs
			// TODO (rghetia): Add support for library Info.
			if len(rs.InstrumentationLibrarySpans) == 0 {
				rs.InstrumentationLibrarySpans = []*tracepb.InstrumentationLibrarySpans{
					{
						Spans: []*tracepb.Span{},
					},
				}
			}
		} else {
			if len(rs.InstrumentationLibrarySpans) > 0 {
				existingRs.InstrumentationLibrarySpans[0].Spans =
					append(existingRs.InstrumentationLibrarySpans[0].Spans,
						rs.InstrumentationLibrarySpans[0].GetSpans()...)
			}
		}
	}
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

func resourceString(res *resourcepb.Resource) string {
	sAttrs := sortedAttributes(res.GetAttributes())
	rstr := ""
	for _, attr := range sAttrs {
		rstr = rstr + attr.String()

	}
	return rstr
}

func sortedAttributes(attrs []*commonpb.AttributeKeyValue) []*commonpb.AttributeKeyValue {
	sort.Slice(attrs[:], func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
}

type MetricSuite struct {
	ServerSuite

	MetricProvider   metric.Provider
	metricController *push.Controller

	mu      sync.RWMutex
	metrics []*metricpb.Metric
}

func (ms *MetricSuite) SetupSuite() {
	ms.ServerSuite.SetupSuite()

	colmetricpb.RegisterMetricsServiceServer(ms.Server, ms)

	ms.ServerSuite.StartServer()

	if ms.MetricProvider == nil {
		ms.metricController = ms.NewPushController(ms.Exporter, nil)
		ms.metricController.SetErrorHandler(func(err error) {
			ms.T().Errorf("testing push controller: %w", err)
		})
		ms.metricController.Start()
		ms.MetricProvider = ms.metricController.Provider()
	}
}

func (ms *MetricSuite) NewPushController(exp *otlp.Exporter, opts []push.Option) *push.Controller {
	integrator := integrator.New(simple.NewWithExactDistribution(), true)
	pusher := push.New(integrator, exp, opts...)
	return pusher
}

func (ms *MetricSuite) GetMetrics() []*metricpb.Metric {
	// copy in order to not change.
	m := make([]*metricpb.Metric, 0, len(ms.metrics))
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return append(m, ms.metrics...)
}

func (ms *MetricSuite) Export(ctx context.Context, exp *colmetricpb.ExportMetricsServiceRequest) (*colmetricpb.ExportMetricsServiceResponse, error) {
	ms.mu.Lock()
	for _, rm := range exp.GetResourceMetrics() {
		// TODO (rghetia) handle multiple resource and library info.
		if len(rm.InstrumentationLibraryMetrics) > 0 {
			ms.metrics = append(ms.metrics, rm.InstrumentationLibraryMetrics[0].Metrics...)
		}
	}
	ms.mu.Unlock()
	return &colmetricpb.ExportMetricsServiceResponse{}, nil
}
func (ms *MetricSuite) TearDownSuite() {
	ms.metricController.Stop()
	ms.ServerSuite.TearDownSuite()
}
