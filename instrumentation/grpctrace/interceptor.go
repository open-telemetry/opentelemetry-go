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

// gRPC tracing middleware
// https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/rpc.md
import (
	"context"
	"io"
	"net"
	"regexp"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	rpcServiceKey  = kv.Key("rpc.service")
	netPeerIPKey   = kv.Key("net.peer.ip")
	netPeerPortKey = kv.Key("net.peer.port")

	messageTypeKey             = kv.Key("message.type")
	messageIDKey               = kv.Key("message.id")
	messageUncompressedSizeKey = kv.Key("message.uncompressed_size")
)

type messageType string

// Event adds an event of the messageType to the span associated with the
// passed context with id and size (if message is a proto message).
func (m messageType) Event(ctx context.Context, id int, message interface{}) {
	span := trace.SpanFromContext(ctx)
	if p, ok := message.(proto.Message); ok {
		span.AddEvent(ctx, "message",
			messageTypeKey.String(string(m)),
			messageIDKey.Int(id),
			messageUncompressedSizeKey.Int(proto.Size(p)),
		)
	} else {
		span.AddEvent(ctx, "message",
			messageTypeKey.String(string(m)),
			messageIDKey.Int(id),
		)
	}
}

const (
	messageSent     messageType = "SENT"
	messageReceived messageType = "RECEIVED"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
//
// For example:
//     tracer := global.Tracer("client-tracer")
//     s := grpc.NewServer(
//         grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor(tracer)),
//         ...,  // (existing DialOptions))
func UnaryClientInterceptor(tracer trace.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()

		var span trace.Span
		ctx, span = tracer.Start(
			ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(peerInfoFromTarget(cc.Target())...),
			trace.WithAttributes(rpcServiceKey.String(serviceFromFullMethod(method))),
		)
		defer span.End()

		Inject(ctx, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

		messageSent.Event(ctx, 1, req)

		err := invoker(ctx, method, req, reply, cc, opts...)

		messageReceived.Event(ctx, 1, reply)

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(s.Code(), s.Message())
		}

		return err
	}
}

type streamEventType int

type streamEvent struct {
	Type streamEventType
	Err  error
}

const (
	closeEvent streamEventType = iota
	receiveEndEvent
	errorEvent
)

// clientStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method call.
type clientStream struct {
	grpc.ClientStream

	desc       *grpc.StreamDesc
	events     chan streamEvent
	eventsDone chan struct{}
	finished   chan error

	receivedMessageID int
	sentMessageID     int
}

var _ = proto.Marshal

func (w *clientStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)

	if err == nil && !w.desc.ServerStreams {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err == io.EOF {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.receivedMessageID++
		messageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *clientStream) SendMsg(m interface{}) error {
	err := w.ClientStream.SendMsg(m)

	w.sentMessageID++
	messageSent.Event(w.Context(), w.sentMessageID, m)

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) Header() (metadata.MD, error) {
	md, err := w.ClientStream.Header()

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return md, err
}

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.sendStreamEvent(closeEvent, nil)
	}

	return err
}

const (
	clientClosedState byte = 1 << iota
	receiveEndedState
)

func wrapClientStream(s grpc.ClientStream, desc *grpc.StreamDesc) *clientStream {
	events := make(chan streamEvent)
	eventsDone := make(chan struct{})
	finished := make(chan error)

	go func() {
		defer close(eventsDone)

		// Both streams have to be closed
		state := byte(0)

		for event := range events {
			switch event.Type {
			case closeEvent:
				state |= clientClosedState
			case receiveEndEvent:
				state |= receiveEndedState
			case errorEvent:
				finished <- event.Err
				return
			}

			if state == clientClosedState|receiveEndedState {
				finished <- nil
				return
			}
		}
	}()

	return &clientStream{
		ClientStream: s,
		desc:         desc,
		events:       events,
		eventsDone:   eventsDone,
		finished:     finished,
	}
}

func (w *clientStream) sendStreamEvent(eventType streamEventType, err error) {
	select {
	case <-w.eventsDone:
	case w.events <- streamEvent{Type: eventType, Err: err}:
	}
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor suitable
// for use in a grpc.Dial call.
//
// For example:
//     tracer := global.Tracer("client-tracer")
//     s := grpc.Dial(
//         grpc.WithStreamInterceptor(grpctrace.StreamClientInterceptor(tracer)),
//         ...,  // (existing DialOptions))
func StreamClientInterceptor(tracer trace.Tracer) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		requestMetadata, _ := metadata.FromOutgoingContext(ctx)
		metadataCopy := requestMetadata.Copy()

		var span trace.Span
		ctx, span = tracer.Start(
			ctx, method,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(peerInfoFromTarget(cc.Target())...),
			trace.WithAttributes(rpcServiceKey.String(serviceFromFullMethod(method))),
		)

		Inject(ctx, &metadataCopy)
		ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

		s, err := streamer(ctx, desc, cc, method, opts...)
		stream := wrapClientStream(s, desc)

		go func() {
			if err == nil {
				err = <-stream.finished
			}

			if err != nil {
				s, _ := status.FromError(err)
				span.SetStatus(s.Code(), s.Message())
			}

			span.End()
		}()

		return stream, err
	}
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
//
// For example:
//     tracer := global.Tracer("client-tracer")
//     s := grpc.Dial(
//         grpc.UnaryInterceptor(grpctrace.UnaryServerInterceptor(tracer)),
//         ...,  // (existing ServerOptions))
func UnaryServerInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestMetadata, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := requestMetadata.Copy()

		entries, spanCtx := Extract(ctx, &metadataCopy)
		ctx = correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		}))

		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(ctx, spanCtx),
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(peerInfoFromContext(ctx)...),
			trace.WithAttributes(rpcServiceKey.String(serviceFromFullMethod(info.FullMethod))),
		)
		defer span.End()

		messageReceived.Event(ctx, 1, req)

		resp, err := handler(ctx, req)
		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(s.Code(), s.Message())
			messageSent.Event(ctx, 1, s.Proto())
		} else {
			messageSent.Event(ctx, 1, resp)
		}

		return resp, err
	}
}

// clientStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context

	receivedMessageID int
	sentMessageID     int
}

func (w *serverStream) Context() context.Context {
	return w.ctx
}

func (w *serverStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)

	if err == nil {
		w.receivedMessageID++
		messageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *serverStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)

	w.sentMessageID++
	messageSent.Event(w.Context(), w.sentMessageID, m)

	return err
}

func wrapServerStream(ctx context.Context, ss grpc.ServerStream) *serverStream {
	return &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor suitable
// for use in a grpc.NewServer call.
//
// For example:
//     tracer := global.Tracer("client-tracer")
//     s := grpc.Dial(
//         grpc.StreamInterceptor(grpctrace.StreamServerInterceptor(tracer)),
//         ...,  // (existing ServerOptions))
func StreamServerInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		requestMetadata, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := requestMetadata.Copy()

		entries, spanCtx := Extract(ctx, &metadataCopy)
		ctx = correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		}))

		ctx, span := tracer.Start(
			trace.ContextWithRemoteSpanContext(ctx, spanCtx),
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(peerInfoFromContext(ctx)...),
			trace.WithAttributes(rpcServiceKey.String(serviceFromFullMethod(info.FullMethod))),
		)
		defer span.End()

		err := handler(srv, wrapServerStream(ctx, ss))

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(s.Code(), s.Message())
		}

		return err
	}
}

func peerInfoFromTarget(target string) []kv.KeyValue {
	host, port, err := net.SplitHostPort(target)

	if err != nil {
		return []kv.KeyValue{}
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return []kv.KeyValue{
		netPeerIPKey.String(host),
		netPeerPortKey.String(port),
	}
}

func peerInfoFromContext(ctx context.Context) []kv.KeyValue {
	p, ok := peer.FromContext(ctx)

	if !ok {
		return []kv.KeyValue{}
	}

	return peerInfoFromTarget(p.Addr.String())
}

var fullMethodRegexp = regexp.MustCompile(`^\/?(?:\S+\.)?(\S+)\/\S+$`)

func serviceFromFullMethod(method string) string {
	match := fullMethodRegexp.FindStringSubmatch(method)
	if len(match) == 0 {
		return ""
	}

	return match[1]
}
