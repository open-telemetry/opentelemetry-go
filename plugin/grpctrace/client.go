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

package grpctrace

import (
	"context"

	"google.golang.org/grpc"
)

func NewUnaryClientInterceptor(opts ...UnaryClientInterceptorOption) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, callOpts...)
	}
}

type UnaryClientInterceptorOption func(*unaryClientInterceptorConfig)

type unaryClientInterceptorConfig struct{}

func NewStreamClientInterceptor(opts ...StreamClientInterceptorOption) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, callOpts...)
	}
}

type StreamClientInterceptorOption func(*streamClientInterceptorConfig)

type streamClientInterceptorConfig struct{}
