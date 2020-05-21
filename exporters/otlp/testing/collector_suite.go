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
	"net"
	"time"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
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

type CollectorTestSuite struct {
	suite.Suite

	TraceService  *TraceService
	MetricService *MetricService

	ServerOpts     []grpc.ServerOption
	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server

	ExporterOpts []otlp.ExporterOption
	Exporter     *otlp.Exporter

	TraceProvider    *trace.Provider
	MetricProvider   metric.Provider
	metricController *push.Controller
}

func (s *CollectorTestSuite) SetupSuite() {
	s.serverAddr = defaultServerAddr

	var err error
	s.ServerListener, err = net.Listen("tcp", s.serverAddr)
	s.serverAddr = s.ServerListener.Addr().String()
	require.NoError(s.T(), err, "failed to allocate a port for server")

	s.Server = grpc.NewServer(s.ServerOpts...)

	if s.TraceService == nil {
		s.TraceService = &TraceService{
			T:   s.T(),
			rsm: map[string]*tracepb.ResourceSpans{},
		}
	}
	coltracepb.RegisterTraceServiceServer(s.Server, s.TraceService)

	if s.MetricService == nil {
		s.MetricService = &MetricService{T: s.T()}
	}
	colmetricpb.RegisterMetricsServiceServer(s.Server, s.MetricService)

	go func() {
		s.Server.Serve(s.ServerListener)
	}()

	if s.Exporter == nil {
		s.Exporter = s.NewExporter()
	}

	if s.TraceProvider == nil {
		s.TraceProvider = s.NewTraceProvider(s.Exporter, nil)
	}

	if s.MetricProvider == nil {
		s.metricController = s.NewPushController(s.Exporter, nil)
		s.metricController.SetErrorHandler(func(err error) {
			s.T().Errorf("testing push controller: %w", err)
		})
		s.metricController.Start()
		s.MetricProvider = s.metricController.Provider()
	}
}

func (s *CollectorTestSuite) ServerAddr() string {
	return s.serverAddr
}

func (s *CollectorTestSuite) NewExporter() *otlp.Exporter {
	opts := []otlp.ExporterOption{
		otlp.WithInsecure(),
		otlp.WithAddress(s.serverAddr),
		otlp.WithReconnectionPeriod(10 * time.Millisecond),
	}
	exp, err := otlp.NewExporter(append(opts, s.ExporterOpts...)...)
	require.NoError(s.T(), err, "failed to create exporter")
	return exp
}

func (s *CollectorTestSuite) NewTraceProvider(exp *otlp.Exporter, opts []trace.ProviderOption) *trace.Provider {
	defaultOpts := []trace.ProviderOption{
		trace.WithConfig(
			trace.Config{
				DefaultSampler: trace.AlwaysSample(),
			},
		),
		trace.WithSyncer(exp),
	}
	p, err := trace.NewProvider(append(defaultOpts, opts...)...)
	require.NoError(s.T(), err, "failed to create trace provider")
	return p
}

func (s *CollectorTestSuite) NewPushController(exp *otlp.Exporter, opts []push.Option) *push.Controller {
	integrator := integrator.New(simple.NewWithExactDistribution(), true)
	pusher := push.New(integrator, exp, opts...)
	return pusher
}

func (s *CollectorTestSuite) TearDownSuite() {
	s.metricController.Stop()
	if s.ServerListener != nil {
		s.Server.GracefulStop()
		s.T().Logf("stopped grpc.Server at: %v", s.ServerAddr())
	}
	s.Exporter.Stop()
}
