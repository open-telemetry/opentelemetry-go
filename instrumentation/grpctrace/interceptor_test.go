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

	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/api/trace/testtrace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/protobuf/proto" //nolint:staticcheck

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/api/kv"
)

type SpanRecorder struct {
	mu    sync.RWMutex
	spans map[string]*testtrace.Span
}

func NewSpanRecorder() *SpanRecorder {
	return &SpanRecorder{spans: make(map[string]*testtrace.Span)}
}

func (sr *SpanRecorder) OnStart(span *testtrace.Span) {}

func (sr *SpanRecorder) OnEnd(span *testtrace.Span) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.spans[span.Name()] = span
}

func (sr *SpanRecorder) Get(name string) (*testtrace.Span, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	s, ok := sr.spans[name]
	return s, ok
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
	clientConn, err := grpc.Dial("fake:connection", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to create client connection: %v", err)
	}

	sr := NewSpanRecorder()
	tp := testtrace.NewProvider(testtrace.WithSpanRecorder(sr))
	tracer := tp.Tracer("grpctrace/client")
	unaryInterceptor := UnaryClientInterceptor(tracer)

	req := &mockProtoMessage{}
	reply := &mockProtoMessage{}
	uniInterceptorInvoker := &mockUICInvoker{}

	checks := []struct {
		method       string
		name         string
		expectedAttr map[kv.Key]kv.Value
		eventsAttr   []map[kv.Key]kv.Value
	}{
		{
			method: "/github.com.serviceName/bar",
			name:   "github.com.serviceName/bar",
			expectedAttr: map[kv.Key]kv.Value{
				standard.RPCSystemKey:   kv.StringValue("grpc"),
				standard.RPCServiceKey:  kv.StringValue("github.com.serviceName"),
				standard.RPCMethodKey:   kv.StringValue("bar"),
				standard.NetPeerIPKey:   kv.StringValue("fake"),
				standard.NetPeerPortKey: kv.StringValue("connection"),
			},
			eventsAttr: []map[kv.Key]kv.Value{
				{
					standard.RPCMessageTypeKey:             kv.StringValue("SENT"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(req))),
				},
				{
					standard.RPCMessageTypeKey:             kv.StringValue("RECEIVED"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			method: "/serviceName/bar",
			name:   "serviceName/bar",
			expectedAttr: map[kv.Key]kv.Value{
				standard.RPCSystemKey:   kv.StringValue("grpc"),
				standard.RPCServiceKey:  kv.StringValue("serviceName"),
				standard.RPCMethodKey:   kv.StringValue("bar"),
				standard.NetPeerIPKey:   kv.StringValue("fake"),
				standard.NetPeerPortKey: kv.StringValue("connection"),
			},
			eventsAttr: []map[kv.Key]kv.Value{
				{
					standard.RPCMessageTypeKey:             kv.StringValue("SENT"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(req))),
				},
				{
					standard.RPCMessageTypeKey:             kv.StringValue("RECEIVED"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			method: "serviceName/bar",
			name:   "serviceName/bar",
			expectedAttr: map[kv.Key]kv.Value{
				standard.RPCSystemKey:   kv.StringValue("grpc"),
				standard.RPCServiceKey:  kv.StringValue("serviceName"),
				standard.RPCMethodKey:   kv.StringValue("bar"),
				standard.NetPeerIPKey:   kv.StringValue("fake"),
				standard.NetPeerPortKey: kv.StringValue("connection"),
			},
			eventsAttr: []map[kv.Key]kv.Value{
				{
					standard.RPCMessageTypeKey:             kv.StringValue("SENT"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(req))),
				},
				{
					standard.RPCMessageTypeKey:             kv.StringValue("RECEIVED"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			method: "invalidName",
			name:   "invalidName",
			expectedAttr: map[kv.Key]kv.Value{
				standard.RPCSystemKey:   kv.StringValue("grpc"),
				standard.NetPeerIPKey:   kv.StringValue("fake"),
				standard.NetPeerPortKey: kv.StringValue("connection"),
			},
			eventsAttr: []map[kv.Key]kv.Value{
				{
					standard.RPCMessageTypeKey:             kv.StringValue("SENT"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(req))),
				},
				{
					standard.RPCMessageTypeKey:             kv.StringValue("RECEIVED"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(reply))),
				},
			},
		},
		{
			method: "/github.com.foo.serviceName_123/method",
			name:   "github.com.foo.serviceName_123/method",
			expectedAttr: map[kv.Key]kv.Value{
				standard.RPCSystemKey:   kv.StringValue("grpc"),
				standard.RPCServiceKey:  kv.StringValue("github.com.foo.serviceName_123"),
				standard.RPCMethodKey:   kv.StringValue("method"),
				standard.NetPeerIPKey:   kv.StringValue("fake"),
				standard.NetPeerPortKey: kv.StringValue("connection"),
			},
			eventsAttr: []map[kv.Key]kv.Value{
				{
					standard.RPCMessageTypeKey:             kv.StringValue("SENT"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(req))),
				},
				{
					standard.RPCMessageTypeKey:             kv.StringValue("RECEIVED"),
					standard.RPCMessageIDKey:               kv.IntValue(1),
					standard.RPCMessageUncompressedSizeKey: kv.IntValue(proto.Size(proto.Message(reply))),
				},
			},
		},
	}

	for _, check := range checks {
		if !assert.NoError(t, unaryInterceptor(context.Background(), check.method, req, reply, clientConn, uniInterceptorInvoker.invoker)) {
			continue
		}
		span, ok := sr.Get(check.name)
		if !assert.True(t, ok, "missing span %q", check.name) {
			continue
		}
		assert.Equal(t, check.expectedAttr, span.Attributes())
		assert.Equal(t, check.eventsAttr, eventAttrMap(span.Events()))
	}
}

func eventAttrMap(events []testtrace.Event) []map[kv.Key]kv.Value {
	maps := make([]map[kv.Key]kv.Value, len(events))
	for i, event := range events {
		maps[i] = event.Attributes
	}
	return maps
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
	clientConn, err := grpc.Dial("fake:connection", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to create client connection: %v", err)
	}

	// tracer
	sr := NewSpanRecorder()
	tp := testtrace.NewProvider(testtrace.WithSpanRecorder(sr))
	tracer := tp.Tracer("grpctrace/Server")
	streamCI := StreamClientInterceptor(tracer)

	var mockClStr mockClientStream
	method := "/github.com.serviceName/bar"
	name := "github.com.serviceName/bar"

	streamClient, err := streamCI(
		context.Background(),
		&grpc.StreamDesc{ServerStreams: true},
		clientConn,
		method,
		func(ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			mockClStr = mockClientStream{Desc: desc, Ctx: ctx}
			return mockClStr, nil
		},
	)
	require.NoError(t, err, "initialize grpc stream client")
	_, ok := sr.Get(name)
	require.False(t, ok, "span should ended while stream is open")

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
	var span *testtrace.Span
	for retry := 0; retry < 5; retry++ {
		span, ok = sr.Get(name)
		if ok {
			break
		}
		time.Sleep(time.Second * 1)
	}
	require.True(t, ok, "missing span %s", name)

	expectedAttr := map[kv.Key]kv.Value{
		standard.RPCSystemKey:   kv.StringValue("grpc"),
		standard.RPCServiceKey:  kv.StringValue("github.com.serviceName"),
		standard.RPCMethodKey:   kv.StringValue("bar"),
		standard.NetPeerIPKey:   kv.StringValue("fake"),
		standard.NetPeerPortKey: kv.StringValue("connection"),
	}
	assert.Equal(t, expectedAttr, span.Attributes())

	events := span.Events()
	require.Len(t, events, 20)
	for i := 0; i < 20; i += 2 {
		msgID := i/2 + 1
		validate := func(eventName string, attrs map[kv.Key]kv.Value) {
			for k, v := range attrs {
				if k == standard.RPCMessageTypeKey && v.AsString() != eventName {
					t.Errorf("invalid event on index: %d expecting %s event, receive %s event", i, eventName, v.AsString())
				}
				if k == standard.RPCMessageIDKey && v != kv.IntValue(msgID) {
					t.Errorf("invalid id for message event expected %d received %d", msgID, v.AsInt32())
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
	sr := NewSpanRecorder()
	tp := testtrace.NewProvider(testtrace.WithSpanRecorder(sr))
	tracer := tp.Tracer("grpctrace/Server")
	usi := UnaryServerInterceptor(tracer)
	deniedErr := status.Error(codes.PermissionDenied, "PERMISSION_DENIED_TEXT")
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, deniedErr
	}
	_, err := usi(context.Background(), &mockProtoMessage{}, &grpc.UnaryServerInfo{}, handler)
	require.Error(t, err)
	assert.Equal(t, err, deniedErr)

	span, ok := sr.Get("")
	if !ok {
		t.Fatalf("failed to export error span")
	}
	assert.Equal(t, span.StatusCode(), codes.PermissionDenied)
	assert.Contains(t, deniedErr.Error(), span.StatusMessage())
	assert.Len(t, span.Events(), 2)
	assert.Equal(t, map[kv.Key]kv.Value{
		kv.Key("message.type"):              kv.StringValue("SENT"),
		kv.Key("message.id"):                kv.IntValue(1),
		kv.Key("message.uncompressed_size"): kv.IntValue(26),
	}, span.Events()[1].Attributes)
}

func TestParseFullMethod(t *testing.T) {
	tests := []struct {
		fullMethod string
		name       string
		attr       []kv.KeyValue
	}{
		{
			fullMethod: "/grpc.test.EchoService/Echo",
			name:       "grpc.test.EchoService/Echo",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("grpc.test.EchoService"),
				standard.RPCMethodKey.String("Echo"),
			},
		}, {
			fullMethod: "/com.example.ExampleRmiService/exampleMethod",
			name:       "com.example.ExampleRmiService/exampleMethod",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("com.example.ExampleRmiService"),
				standard.RPCMethodKey.String("exampleMethod"),
			},
		}, {
			fullMethod: "/MyCalcService.Calculator/Add",
			name:       "MyCalcService.Calculator/Add",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("MyCalcService.Calculator"),
				standard.RPCMethodKey.String("Add"),
			},
		}, {
			fullMethod: "/MyServiceReference.ICalculator/Add",
			name:       "MyServiceReference.ICalculator/Add",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("MyServiceReference.ICalculator"),
				standard.RPCMethodKey.String("Add"),
			},
		}, {
			fullMethod: "/MyServiceWithNoPackage/theMethod",
			name:       "MyServiceWithNoPackage/theMethod",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("MyServiceWithNoPackage"),
				standard.RPCMethodKey.String("theMethod"),
			},
		}, {
			fullMethod: "/pkg.srv",
			name:       "pkg.srv",
			attr:       []kv.KeyValue(nil),
		}, {
			fullMethod: "/pkg.srv/",
			name:       "pkg.srv/",
			attr: []kv.KeyValue{
				standard.RPCServiceKey.String("pkg.srv"),
			},
		},
	}

	for _, test := range tests {
		n, a := parseFullMethod(test.fullMethod)
		assert.Equal(t, test.name, n)
		assert.Equal(t, test.attr, a)
	}
}
