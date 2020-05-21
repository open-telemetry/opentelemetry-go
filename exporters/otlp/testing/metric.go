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
	"sync"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

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
