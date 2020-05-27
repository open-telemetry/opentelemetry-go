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
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/kv/value"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testExporter struct {
	mu      sync.Mutex
	spanMap map[string]*export.SpanData
}

func (t *testExporter) ExportSpan(ctx context.Context, s *export.SpanData) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spanMap[s.Name] = s
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
	exp := &testExporter{spanMap: make(map[string]*export.SpanData)}
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
		expectedAttr map[kv.Key]value.Value
		eventsAttr   []map[kv.Key]value.Value
	}{
		{
			name: "/github.com.serviceName/bar",
			expectedAttr: map[kv.Key]value.Value{
				rpcServiceKey:  value.String("serviceName"),
				netPeerIPKey:   value.String("fake"),
				netPeerPortKey: value.String("connection"),
			},
			eventsAttr: []map[kv.Key]value.Value{
				{
					messageTypeKey:             value.String("SENT"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(req))),
				},
				{
					messageTypeKey:             value.String("RECEIVED"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			name: "/serviceName/bar",
			expectedAttr: map[kv.Key]value.Value{
				rpcServiceKey:  value.String("serviceName"),
				netPeerIPKey:   value.String("fake"),
				netPeerPortKey: value.String("connection"),
			},
			eventsAttr: []map[kv.Key]value.Value{
				{
					messageTypeKey:             value.String("SENT"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(req))),
				},
				{
					messageTypeKey:             value.String("RECEIVED"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			name: "serviceName/bar",
			expectedAttr: map[kv.Key]value.Value{
				rpcServiceKey:  value.String("serviceName"),
				netPeerIPKey:   value.String("fake"),
				netPeerPortKey: value.String("connection"),
			},
			eventsAttr: []map[kv.Key]value.Value{
				{
					messageTypeKey:             value.String("SENT"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(req))),
				},
				{
					messageTypeKey:             value.String("RECEIVED"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			name: "invalidName",
			expectedAttr: map[kv.Key]value.Value{
				rpcServiceKey:  value.String(""),
				netPeerIPKey:   value.String("fake"),
				netPeerPortKey: value.String("connection"),
			},
			eventsAttr: []map[kv.Key]value.Value{
				{
					messageTypeKey:             value.String("SENT"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(req))),
				},
				{
					messageTypeKey:             value.String("RECEIVED"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			name: "/github.com.foo.serviceName_123/method",
			expectedAttr: map[kv.Key]value.Value{
				rpcServiceKey:  value.String("serviceName_123"),
				netPeerIPKey:   value.String("fake"),
				netPeerPortKey: value.String("connection"),
			},
			eventsAttr: []map[kv.Key]value.Value{
				{
					messageTypeKey:             value.String("SENT"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(req))),
				},
				{
					messageTypeKey:             value.String("RECEIVED"),
					messageIDKey:               value.Int(1),
					messageUncompressedSizeKey: value.Int(proto.Size(proto.Message(reply))),
				},
			},
		},
	}

	for _, check := range checks {
		err = unaryInterceptor(context.Background(), check.name, req, reply, clientConn, uniInterceptorInvoker.invoker)
		if err != nil {
			t.Errorf("failed to run unary interceptor: %v", err)
			continue
		}

		spanData, ok := exp.spanMap[check.name]
		if !ok {
			t.Errorf("no span data found for name < %s >", check.name)
			continue
		}

		attrs := spanData.Attributes
		if len(check.expectedAttr) > len(attrs) {
			t.Errorf("attributes received are less than expected attributes, received %d, expected %d",
				len(attrs), len(check.expectedAttr))
		}
		for _, attr := range attrs {
			expectedAttr, ok := check.expectedAttr[attr.Key]
			if ok {
				if expectedAttr != attr.Value {
					t.Errorf("name: %s invalid %s found. expected %s, actual %s", check.name, string(attr.Key),
						expectedAttr.AsString(), attr.Value.AsString())
				}
				delete(check.expectedAttr, attr.Key)
			} else {
				t.Errorf("attribute %s not found in expected attributes map", string(attr.Key))
			}
		}

		// Check if any expected attr not seen
		if len(check.expectedAttr) > 0 {
			for attr := range check.expectedAttr {
				t.Errorf("missing attribute %s in span", string(attr))
			}
		}

		events := spanData.MessageEvents
		if len(check.eventsAttr) > len(events) {
			t.Errorf("events received are less than expected events, received %d, expected %d",
				len(events), len(check.eventsAttr))
		}
		for event := 0; event < len(check.eventsAttr); event++ {
			for _, attr := range events[event].Attributes {
				expectedAttr, ok := check.eventsAttr[event][attr.Key]
				if ok {
					if attr.Value != expectedAttr {
						t.Errorf("invalid value for attribute %s in events, expected %s actual %s",
							string(attr.Key), attr.Value.AsString(), expectedAttr.AsString())
					}
					delete(check.eventsAttr[event], attr.Key)
				} else {
					t.Errorf("attribute in event %s not found in expected attributes map", string(attr.Key))
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

type mockClientStream struct {
	Desc *grpc.StreamDesc
	Ctx  context.Context
}

func (mockClientStream) SendMsg(m interface{}) error  { return nil }
func (mockClientStream) RecvMsg(m interface{}) error  { return nil }
func (mockClientStream) CloseSend() error             { return nil }
func (c mockClientStream) Context() context.Context   { return c.Ctx }
func (mockClientStream) Header() (metadata.MD, error) { return nil, nil }
func (mockClientStream) Trailer() metadata.MD         { return nil }

func TestStreamClientInterceptor(t *testing.T) {
	exp := &testExporter{spanMap: make(map[string]*export.SpanData)}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
		},
		))
	clientConn, err := grpc.Dial("fake:connection", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to create client connection: %v", err)
	}

	// tracer
	tracer := tp.Tracer("grpctrace/Server")
	streamCI := StreamClientInterceptor(tracer)

	var mockClStr mockClientStream
	methodName := "/github.com.serviceName/bar"

	streamClient, err := streamCI(context.Background(),
		&grpc.StreamDesc{ServerStreams: true},
		clientConn,
		methodName,
		func(ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			mockClStr = mockClientStream{Desc: desc, Ctx: ctx}
			return mockClStr, nil
		})

	if err != nil {
		t.Fatalf("failed to initialize grpc stream client: %v", err)
	}

	// no span exported while stream is open
	if _, ok := exp.spanMap[methodName]; ok {
		t.Fatalf("span shouldn't end while stream is open")
	}

	req := &mockProtoMessage{}
	reply := &mockProtoMessage{}

	// send and receive fake data
	for i := 0; i < 10; i++ {
		_ = streamClient.SendMsg(req)
		_ = streamClient.RecvMsg(reply)
	}

	// close client and server stream
	_ = streamClient.CloseSend()
	mockClStr.Desc.ServerStreams = false
	_ = streamClient.RecvMsg(reply)

	// added retry because span end is called in separate go routine
	var spanData *export.SpanData
	for retry := 0; retry < 5; retry++ {
		ok := false
		exp.mu.Lock()
		spanData, ok = exp.spanMap[methodName]
		exp.mu.Unlock()
		if ok {
			break
		}
		time.Sleep(time.Second * 1)
	}
	if spanData == nil {
		t.Fatalf("no span data found for name < %s >", methodName)
	}

	attrs := spanData.Attributes
	expectedAttr := map[kv.Key]string{
		rpcServiceKey:  "serviceName",
		netPeerIPKey:   "fake",
		netPeerPortKey: "connection",
	}

	for _, attr := range attrs {
		expected, ok := expectedAttr[attr.Key]
		if ok {
			if expected != attr.Value.AsString() {
				t.Errorf("name: %s invalid %s found. expected %s, actual %s", methodName, string(attr.Key),
					expected, attr.Value.AsString())
			}
		}
	}

	events := spanData.MessageEvents
	if len(events) != 20 {
		t.Fatalf("incorrect number of events expected 20 got %d", len(events))
	}
	for i := 0; i < 20; i += 2 {
		msgID := i/2 + 1
		validate := func(eventName string, attrs []kv.KeyValue) {
			for _, attr := range attrs {
				if attr.Key == messageTypeKey && attr.Value.AsString() != eventName {
					t.Errorf("invalid event on index: %d expecting %s event, receive %s event", i, eventName, attr.Value.AsString())
				}
				if attr.Key == messageIDKey && attr.Value != value.Int(msgID) {
					t.Errorf("invalid id for message event expected %d received %d", msgID, attr.Value.AsInt32())
				}
			}
		}
		validate("SENT", events[i].Attributes)
		validate("RECEIVED", events[i+1].Attributes)
	}

	// ensure CloseSend can be subsequently called
	_ = streamClient.CloseSend()
}

func TestServerInterceptorError(t *testing.T) {
	exp := &testExporter{spanMap: make(map[string]*export.SpanData)}
	tp, err := sdktrace.NewProvider(
		sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
		}),
	)
	require.NoError(t, err)

	tracer := tp.Tracer("grpctrace/Server")
	usi := UnaryServerInterceptor(tracer)
	deniedErr := status.Error(codes.PermissionDenied, "PERMISSION_DENIED_TEXT")
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, deniedErr
	}
	_, err = usi(context.Background(), &mockProtoMessage{}, &grpc.UnaryServerInfo{}, handler)
	require.Error(t, err)
	assert.Equal(t, err, deniedErr)

	span, ok := exp.spanMap[""]
	if !ok {
		t.Fatalf("failed to export error span")
	}
	assert.Equal(t, span.StatusCode, codes.PermissionDenied)
	assert.Contains(t, deniedErr.Error(), span.StatusMessage)
	assert.Len(t, span.MessageEvents, 2)
	assert.Equal(t, []kv.KeyValue{
		kv.String("message.type", "SENT"),
		kv.Int("message.id", 1),
		kv.Int("message.uncompressed_size", 26),
	}, span.MessageEvents[1].Attributes)
}
