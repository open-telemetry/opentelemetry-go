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
	"fmt"
	"sort"
	"sync"
	"testing"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"
)

var _ coltracepb.TraceServiceServer = (*TraceService)(nil)
var _ colmetricpb.MetricsServiceServer = (*MetricService)(nil)

var errAlreadyStopped = fmt.Errorf("already stopped")

type TraceService struct {
	T *testing.T

	mu  sync.RWMutex
	rsm map[string]*tracepb.ResourceSpans
}

func (ts *TraceService) GetSpans() []*tracepb.Span {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	spans := []*tracepb.Span{}
	for _, rs := range ts.rsm {
		spans = append(spans, rs.InstrumentationLibrarySpans[0].Spans...)
	}
	return spans
}

func (ts *TraceService) GetResourceSpans() []*tracepb.ResourceSpans {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	rss := make([]*tracepb.ResourceSpans, 0, len(ts.rsm))
	for _, rs := range ts.rsm {
		rss = append(rss, rs)
	}
	return rss
}

func (ts *TraceService) Export(ctx context.Context, exp *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	rss := exp.GetResourceSpans()
	for _, rs := range rss {
		rstr := resourceString(rs.Resource)
		existingRs, ok := ts.rsm[rstr]
		if !ok {
			ts.rsm[rstr] = rs
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

type MetricService struct {
	T *testing.T

	mu      sync.RWMutex
	metrics []*metricpb.Metric
}

func (ms *MetricService) GetMetrics() []*metricpb.Metric {
	// copy in order to not change.
	m := make([]*metricpb.Metric, 0, len(ms.metrics))
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return append(m, ms.metrics...)
}

func (ms *MetricService) Export(ctx context.Context, exp *colmetricpb.ExportMetricsServiceRequest) (*colmetricpb.ExportMetricsServiceResponse, error) {
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
