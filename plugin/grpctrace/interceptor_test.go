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
package grpctrace

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testExporter struct {
	spanMap map[string][]*export.SpanData
}

func (t *testExporter) ExportSpan(ctx context.Context, s *export.SpanData) {
	t.spanMap[s.Name] = append(t.spanMap[s.Name], s)
}

type mockUICInvoker struct {
	ctx context.Context
}

func (mcuici *mockUICInvoker) invoker(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	mcuici.ctx = ctx
	return nil
}

type mockProtoMessage struct{}

func (mm *mockProtoMessage) Reset() {
}

func (mm *mockProtoMessage) String() string {
	return "mock"
}

func (mm *mockProtoMessage) ProtoMessage() {
}

func TestUnaryClientInterceptor(t *testing.T) {
	exp := &testExporter{make(map[string][]*export.SpanData)}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
		},
		))
	global.SetTraceProvider(tp)

	clientConn, err := grpc.Dial("fake:connection", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to create client connection: %v", err)
	}

	tracer := tp.Tracer("grpctrace/client")
	unaryInterceptor := UnaryClientInterceptor(tracer)

	req := &mockProtoMessage{}
	reply := &mockProtoMessage{}
	uniInterceptorInvoker := &mockUICInvoker{}

	checks := []struct {
		name         string
		expectedAttr map[core.Key]core.Value
		eventsAttr   [][]core.KeyValue
	}{
		{
			name: fmt.Sprintf("/foo.%s/bar", "serviceName"),
			expectedAttr: map[core.Key]core.Value{
				rpcServiceKey:  core.String("serviceName"),
				netPeerIPKey:   core.String("fake"),
				netPeerPortKey: core.String("connection"),
			},
			eventsAttr: [][]core.KeyValue{
				{
					core.KeyValue{Key: messageTypeKey, Value: core.String("SENT")},
					core.KeyValue{Key: messageIDKey, Value: core.Int(1)},
				},
				{
					core.KeyValue{Key: messageTypeKey, Value: core.String("RECEIVED")},
					core.KeyValue{Key: messageIDKey, Value: core.Int(1)},
				},
			},
		},
	}

	for _, check := range checks {
		err = unaryInterceptor(context.Background(), check.name, req, reply, clientConn, uniInterceptorInvoker.invoker)
		if err != nil {
			t.Fatalf("failed to run unary interceptor: %v", err)
		}

		attrs := exp.spanMap[check.name][0].Attributes
		for _, attr := range attrs {
			expectedAttr, ok := check.expectedAttr[attr.Key]
			if ok {
				if expectedAttr != attr.Value {
					t.Fatalf("invalid %s found. expected %s, actual %s", string(attr.Key),
						expectedAttr.AsString(), attr.Value.AsString())
				}
			}
		}

		events := exp.spanMap[check.name][0].MessageEvents
		for event := 0; event < len(check.eventsAttr); event++ {
			for attr := 0; attr < len(check.eventsAttr[event]); attr++ {
				if events[event].Attributes[attr] != check.eventsAttr[event][attr] {
					t.Fatalf("invalid attribute in events")
				}
			}
		}
	}
}
