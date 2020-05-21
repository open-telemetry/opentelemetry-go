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
	"sort"
	"sync"

	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/trace"
)

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
