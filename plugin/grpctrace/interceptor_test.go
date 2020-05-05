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
		eventsAttr   []map[core.Key]core.Value
	}{
		{
			name: "/github.com.serviceName/bar",
			expectedAttr: map[core.Key]core.Value{
				rpcServiceKey:  core.String("serviceName"),
				netPeerIPKey:   core.String("fake"),
				netPeerPortKey: core.String("connection"),
			},
			eventsAttr: []map[core.Key]core.Value{
				{
					messageTypeKey: core.String("SENT"),
					messageIDKey:   core.Int(1),
				},
				{
					messageTypeKey: core.String("RECEIVED"),
					messageIDKey:   core.Int(1),
				},
			},
		},
		{
			name: "/serviceName/bar",
			expectedAttr: map[core.Key]core.Value{
				rpcServiceKey: core.String("serviceName"),
			},
			eventsAttr: []map[core.Key]core.Value{
				{
					messageTypeKey: core.String("SENT"),
					messageIDKey:   core.Int(1),
				},
				{
					messageTypeKey: core.String("RECEIVED"),
					messageIDKey:   core.Int(1),
				},
			},
		},
		{
			name:         "serviceName/bar",
			expectedAttr: map[core.Key]core.Value{rpcServiceKey: core.String("serviceName")},
		},
		{
			name:         "invalidName",
			expectedAttr: map[core.Key]core.Value{rpcServiceKey: core.String("")},
		},
		{
			name:         "/github.com.foo.serviceName_123/method",
			expectedAttr: map[core.Key]core.Value{rpcServiceKey: core.String("serviceName_123")},
		},
	}

	for idx, check := range checks {
		fmt.Println("================", idx, "==================")
		err = unaryInterceptor(context.Background(), check.name, req, reply, clientConn, uniInterceptorInvoker.invoker)
		if err != nil {
			t.Fatalf("failed to run unary interceptor: %v", err)
		}

		spanData, ok := exp.spanMap[check.name]
		if !ok || len(spanData) == 0 {
			t.Fatalf("no span data found for name < %s >", check.name)
		}

		attrs := spanData[0].Attributes
		for _, attr := range attrs {
			expectedAttr, ok := check.expectedAttr[attr.Key]
			if ok {
				if expectedAttr != attr.Value {
					t.Errorf("name: %s invalid %s found. expected %s, actual %s", check.name, string(attr.Key),
						expectedAttr.AsString(), attr.Value.AsString())
				}
				delete(check.expectedAttr, attr.Key)
			}
		}

		// Check if any expected attr not seen
		if len(check.expectedAttr) > 0 {
			for attr := range check.expectedAttr {
				t.Errorf("missing attribute %s in span", string(attr))
			}
		}

		events := spanData[0].MessageEvents
		for event := 0; event < len(check.eventsAttr); event++ {
			for _, attr := range events[event].Attributes {
				expectedAttr, ok := check.eventsAttr[event][attr.Key]
				if ok {
					if attr.Value != expectedAttr {
						t.Errorf("invalid value for attribute %s in events, expected %s actual %s",
							string(attr.Key), attr.Value.AsString(), expectedAttr.AsString())
					}
					delete(check.eventsAttr[event], attr.Key)
				}
			}
			if len(check.eventsAttr[event]) > 0 {
				for attr := range check.eventsAttr[event] {
					t.Errorf("missing attribute %s in span event", string(attr))
				}
			}
		}
	}
}
