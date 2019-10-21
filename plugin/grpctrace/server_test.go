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

func TestUnaryServerInterceptor(t *testing.T) {
	t.Run("calls the original handler with the expected arguments", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		var actualCtx context.Context
		var actualReq interface{}

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			actualCtx = ctx
			actualReq = req
			return nil, nil
		}
		subject := grpctrace.NewUnaryServerInterceptor()

		var ctxKey testCtxKey
		expectedCtxValue := "expected value"
		ctx := context.WithValue(context.Background(), ctxKey, expectedCtxValue)
		expectedReq := "expected request"

		_, err := subject(ctx, expectedReq, &grpc.UnaryServerInfo{}, handler)
		e.Expect(err).ToBeNil()

		e.Expect(actualCtx.Value(ctxKey)).ToEqual(expectedCtxValue)
		e.Expect(actualReq).ToEqual(expectedReq)
	})

	t.Run("returns the original response if the handler succeeds", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		expectedRes := "expected response"
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return expectedRes, nil
		}
		subject := grpctrace.NewUnaryServerInterceptor()

		res, err := subject(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		e.Expect(err).ToBeNil()
		e.Expect(res).ToEqual(expectedRes)
	})

	t.Run("returns the original error if the handler errors", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		expectedErr := errors.New("expected error")
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, expectedErr
		}
		subject := grpctrace.NewUnaryServerInterceptor()

		_, err := subject(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

		e.Expect(err).ToEqual(expectedErr)
	})
}

func TestStreamServerInterceptor(t *testing.T) {
	t.Run("calls the original handler with the expected arguments, including a functional ServerStream", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		var actualSrv interface{}
		var actualServerStream grpc.ServerStream

		handler := func(srv interface{}, ss grpc.ServerStream) error {
			actualSrv = srv
			actualServerStream = ss

			return nil
		}
		subject := grpctrace.NewStreamServerInterceptor()

		expectedSrv := 123
		var expectedSetHeaderErr error
		var expectedSendHeaderErr error
		var ctxKey testCtxKey
		expectedCtxValue := "expected value"
		expectedCtx := context.WithValue(context.Background(), ctxKey, expectedCtxValue)
		var expectedSendMsgErr error
		var expectedRecvMsgErr error

		var actualSetHeaderMetadata metadata.MD
		var actualSendHeaderMetadata metadata.MD
		var actualSetTrailerMetadata metadata.MD
		var actualSendMsgM interface{}
		var actualRecvMsgM interface{}

		serverStream := &testServerStream{
			setHeader: func(md metadata.MD) error {
				actualSetHeaderMetadata = md

				return expectedSetHeaderErr
			},
			sendHeader: func(md metadata.MD) error {
				actualSendHeaderMetadata = md

				return expectedSendHeaderErr
			},
			setTrailer: func(md metadata.MD) {
				actualSetTrailerMetadata = md
			},
			context: func() context.Context {
				return expectedCtx
			},
			sendMsg: func(m interface{}) error {
				actualSendMsgM = m

				return expectedSendMsgErr
			},
			recvMsg: func(m interface{}) error {
				actualRecvMsgM = m

				return expectedRecvMsgErr
			},
		}

		err := subject(expectedSrv, serverStream, nil, handler)
		e.Expect(err).ToBeNil()

		e.Expect(actualSrv).ToEqual(expectedSrv)

		expectedSetHeaderMetadata := map[string][]string{
			"SetHeader": {"value"},
		}
		err = actualServerStream.SetHeader(expectedSetHeaderMetadata)
		e.Expect(err).ToBeNil()
		for k, v := range expectedSetHeaderMetadata {
			e.Expect(actualSetHeaderMetadata[k]).ToEqual(v)
		}

		expectedSetHeaderErr = errors.New("test SetHeader error")
		err = actualServerStream.SetHeader(nil)
		e.Expect(err).ToEqual(expectedSetHeaderErr)

		expectedSendHeaderMetadata := map[string][]string{
			"SendHeader": {"value"},
		}
		err = actualServerStream.SendHeader(expectedSendHeaderMetadata)
		e.Expect(err).ToBeNil()
		for k, v := range expectedSendHeaderMetadata {
			e.Expect(actualSendHeaderMetadata[k]).ToEqual(v)
		}

		err = actualServerStream.SendHeader(nil)
		e.Expect(err).ToBeNil()

		expectedSendHeaderErr = errors.New("test SendHeader error")
		err = actualServerStream.SendHeader(nil)
		e.Expect(err).ToEqual(expectedSendHeaderErr)

		expectedSetTrailerMetadata := map[string][]string{
			"SetTrailer": {"value"},
		}
		actualServerStream.SetTrailer(expectedSetTrailerMetadata)
		for k, v := range expectedSetTrailerMetadata {
			e.Expect(actualSetTrailerMetadata[k]).ToEqual(v)
		}

		ctx := actualServerStream.Context()
		e.Expect(ctx.Value(ctxKey)).ToEqual(expectedCtxValue)

		expectedSendMsgM := 123
		err = actualServerStream.SendMsg(expectedSendMsgM)
		e.Expect(err).ToBeNil()
		e.Expect(actualSendMsgM).ToEqual(expectedSendMsgM)

		expectedSendMsgErr = errors.New("test SendMsg error")
		err = actualServerStream.SendMsg(nil)
		e.Expect(err).ToEqual(expectedSendMsgErr)

		expectedRecvMsgM := "abc"
		err = actualServerStream.RecvMsg(expectedRecvMsgM)
		e.Expect(err).ToBeNil()
		e.Expect(actualRecvMsgM).ToEqual(expectedRecvMsgM)

		expectedRecvMsgErr = errors.New("test RecvMsg error")
		err = actualServerStream.RecvMsg(nil)
		e.Expect(err).ToEqual(expectedRecvMsgErr)
	})

	t.Run("succeeds if the handler succeeds", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		handler := func(srv interface{}, ss grpc.ServerStream) error {
			return nil
		}
		subject := grpctrace.NewStreamServerInterceptor()

		err := subject(nil, nil, nil, handler)

		e.Expect(err).ToBeNil()
	})

	t.Run("returns the original error if the handler errors", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		expectedErr := errors.New("test error")

		handler := func(srv interface{}, ss grpc.ServerStream) error {
			return expectedErr
		}
		subject := grpctrace.NewStreamServerInterceptor()

		err := subject(nil, nil, nil, handler)

		e.Expect(err).ToEqual(expectedErr)
	})
}

var _ grpc.ServerStream = (*testServerStream)(nil)

type testServerStream struct {
	setHeader  func(md metadata.MD) error
	sendHeader func(md metadata.MD) error
	setTrailer func(md metadata.MD)
	context    func() context.Context
	sendMsg    func(m interface{}) error
	recvMsg    func(m interface{}) error
}

func (ss *testServerStream) SetHeader(md metadata.MD) error {
	return ss.setHeader(md)
}

func (ss *testServerStream) SendHeader(md metadata.MD) error {
	return ss.sendHeader(md)
}

func (ss *testServerStream) SetTrailer(md metadata.MD) {
	ss.setTrailer(md)
}

func (ss *testServerStream) Context() context.Context {
	return ss.context()
}

func (ss *testServerStream) SendMsg(m interface{}) error {
	return ss.sendMsg(m)
}

func (ss *testServerStream) RecvMsg(m interface{}) error {
	return ss.recvMsg(m)
}
