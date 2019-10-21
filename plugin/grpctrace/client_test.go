// Copyright 2019, OpenTelemetry Authors
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

package grpctrace_test

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/internal/matchers"
	"go.opentelemetry.io/plugin/grpctrace"
)

func TestUnaryClientInterceptor(t *testing.T) {
	t.Run("calls the invoker with the expected arguments", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		var actualCtx context.Context
		var actualMethod string
		var actualReq interface{}
		var actualReply interface{}
		var actualClientConn *grpc.ClientConn
		var actualCallOpts []grpc.CallOption

		invoker := func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			actualCtx = ctx
			actualMethod = method
			actualReq = req
			actualReply = reply
			actualClientConn = cc
			actualCallOpts = opts

			return nil
		}
		subject := grpctrace.NewUnaryClientInterceptor()

		var ctxKey testCtxKey
		expectedCtxValue := "expected value"
		ctx := context.WithValue(context.Background(), ctxKey, expectedCtxValue)
		expectedMethod := "test method"
		expectedReq := "abc"
		expectedReply := "123"
		expectedClientConn := &grpc.ClientConn{}
		expectedCallOpts := []grpc.CallOption{
			grpc.MaxCallRecvMsgSize(3),
			grpc.MaxCallSendMsgSize(9),
		}

		err := subject(ctx, expectedMethod, expectedReq, expectedReply, expectedClientConn, invoker, expectedCallOpts...)
		e.Expect(err).ToBeNil()

		e.Expect(actualCtx.Value(ctxKey)).ToEqual(expectedCtxValue)
		e.Expect(actualCtx.Value(ctxKey)).ToEqual(expectedCtxValue)
		e.Expect(actualMethod).ToEqual(expectedMethod)
		e.Expect(actualReq).ToEqual(expectedReq)
		e.Expect(actualReply).ToEqual(expectedReply)
		e.Expect(actualClientConn).ToEqual(expectedClientConn)
		e.Expect(actualCallOpts).ToMatchInAnyOrder(expectedCallOpts)
	})

	t.Run("succeeds if the invoker succeeds", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		invoker := func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}
		subject := grpctrace.NewUnaryClientInterceptor()

		err := subject(context.Background(), "", nil, nil, nil, invoker)

		e.Expect(err).ToBeNil()
	})

	t.Run("returns the error returned by the invoker if the invoker errors", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		expectedErr := errors.New("test error")
		invoker := func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return expectedErr
		}
		subject := grpctrace.NewUnaryClientInterceptor()

		err := subject(context.Background(), "", nil, nil, nil, invoker)

		e.Expect(err).ToEqual(expectedErr)
	})
}

func TestStreamClientInterceptor(t *testing.T) {
	t.Run("calls the streamer with the expected arguments", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		var actualCtx context.Context
		var actualDesc *grpc.StreamDesc
		var actualClientConn *grpc.ClientConn
		var actualMethod string
		var actualCallOpts []grpc.CallOption

		streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			actualCtx = ctx
			actualDesc = desc
			actualClientConn = cc
			actualMethod = method
			actualCallOpts = opts

			return nil, nil
		}
		subject := grpctrace.NewStreamClientInterceptor()

		var ctxKey testCtxKey
		expectedCtxValue := "expected value"
		ctx := context.WithValue(context.Background(), ctxKey, expectedCtxValue)
		expectedDesc := &grpc.StreamDesc{
			StreamName: "test stream",
		}
		expectedClientConn := &grpc.ClientConn{}
		expectedMethod := "test method"
		expectedCallOpts := []grpc.CallOption{
			grpc.MaxCallRecvMsgSize(3),
			grpc.MaxCallSendMsgSize(9),
		}

		_, err := subject(ctx, expectedDesc, expectedClientConn, expectedMethod, streamer, expectedCallOpts...)
		e.Expect(err).ToBeNil()

		e.Expect(actualCtx.Value(ctxKey)).ToEqual(expectedCtxValue)
		e.Expect(actualDesc).ToEqual(expectedDesc)
		e.Expect(actualClientConn).ToEqual(expectedClientConn)
		e.Expect(actualMethod).ToEqual(expectedMethod)
		e.Expect(actualCallOpts).ToMatchInAnyOrder(expectedCallOpts)
	})

	t.Run("returns a stream that sends and receives messages as expected if the streamer succeeds", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		var expectedHeaderMetadata metadata.MD
		var expectedHeaderErr error
		var expectedTrailerMetadata metadata.MD
		var expectedCloseSendErr error
		var expectedSendMsgErr error
		var expectedRecvMsgErr error

		var actualSendMsgM interface{}
		var actualRecvMsgM interface{}

		streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return &testClientStream{
				header: func() (metadata.MD, error) {
					return expectedHeaderMetadata, expectedHeaderErr
				},
				trailer: func() metadata.MD {
					return expectedTrailerMetadata
				},
				closeSend: func() error {
					return expectedCloseSendErr
				},
				context: func() context.Context {
					return ctx
				},
				sendMsg: func(m interface{}) error {
					actualSendMsgM = m

					return expectedSendMsgErr
				},
				recvMsg: func(m interface{}) error {
					actualRecvMsgM = m

					return expectedRecvMsgErr
				},
			}, nil
		}
		subject := grpctrace.NewStreamClientInterceptor()

		var ctxKey testCtxKey
		expectedCtxValue := "expected value"
		ctx := context.WithValue(context.Background(), ctxKey, expectedCtxValue)

		stream, err := subject(ctx, nil, nil, "", streamer)

		e.Expect(err).ToBeNil()

		expectedHeaderMetadata = map[string][]string{
			"header": {"value"},
		}
		md, err := stream.Header()
		e.Expect(err).ToBeNil()
		for k, v := range expectedHeaderMetadata {
			e.Expect(md[k]).ToEqual(v)
		}

		expectedHeaderErr = errors.New("test Header error")
		_, err = stream.Header()
		e.Expect(err).ToEqual(expectedHeaderErr)

		expectedTrailerMetadata = map[string][]string{
			"trailer": {"value"},
		}
		md = stream.Trailer()
		for k, v := range expectedTrailerMetadata {
			e.Expect(md[k]).ToEqual(v)
		}

		err = stream.CloseSend()
		e.Expect(err).ToBeNil()

		expectedCloseSendErr = errors.New("test CloseSend error")
		err = stream.CloseSend()
		e.Expect(err).ToEqual(expectedCloseSendErr)

		expectedSendMsgM := 123
		err = stream.SendMsg(expectedSendMsgM)
		e.Expect(err).ToBeNil()
		e.Expect(actualSendMsgM).ToEqual(expectedSendMsgM)

		expectedSendMsgErr = errors.New("test SendMsg error")
		err = stream.SendMsg(nil)
		e.Expect(err).ToEqual(expectedSendMsgErr)

		expectedRecvMsgM := "abc"
		err = stream.RecvMsg(expectedRecvMsgM)
		e.Expect(err).ToBeNil()
		e.Expect(actualRecvMsgM).ToEqual(expectedRecvMsgM)

		expectedRecvMsgErr = errors.New("test RecvMsg error")
		err = stream.RecvMsg(nil)
		e.Expect(err).ToEqual(expectedRecvMsgErr)
	})

	t.Run("returns the error returned by the streamer if the streamer errors", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		expectedErr := errors.New("test error")
		streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return nil, expectedErr
		}
		subject := grpctrace.NewStreamClientInterceptor()

		_, err := subject(context.Background(), nil, nil, "", streamer)

		e.Expect(err).ToEqual(expectedErr)
	})
}

var _ grpc.ClientStream = (*testClientStream)(nil)

type testClientStream struct {
	header    func() (metadata.MD, error)
	trailer   func() metadata.MD
	closeSend func() error
	context   func() context.Context
	sendMsg   func(m interface{}) error
	recvMsg   func(m interface{}) error
}

func (cs *testClientStream) Header() (metadata.MD, error) {
	return cs.header()
}

func (cs *testClientStream) Trailer() metadata.MD {
	return cs.trailer()
}

func (cs *testClientStream) CloseSend() error {
	return cs.closeSend()
}

func (cs *testClientStream) Context() context.Context {
	return cs.context()
}

func (cs *testClientStream) SendMsg(m interface{}) error {
	return cs.sendMsg(m)
}

func (cs *testClientStream) RecvMsg(m interface{}) error {
	return cs.recvMsg(m)
}
