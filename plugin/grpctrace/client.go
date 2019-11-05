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
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

func NewUnaryClientInterceptor(opts ...UnaryClientInterceptorOption) grpc.UnaryClientInterceptor {
	c := newUnaryClientInterceptorConfig(opts...)
	tracer := c.tracer

	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		name := nameFromMethod(method)
		service := serviceFromMethod(method)

		startAttrs := []core.KeyValue{
			core.Key(componentKey).String(componentValue),
			core.Key(peerServiceKey).String(service),
		}

		if p, ok := peer.FromContext(ctx); ok {
			if hostname, port, err := net.SplitHostPort(p.Addr.String()); err == nil {
				startAttrs = append(startAttrs, []core.KeyValue{
					core.Key(peerHostnameKey).String(hostname),
					core.Key(peerPortKey).String(port),
				}...)
			}
		}

		startOpts := []trace.SpanOption{
			trace.WithAttributes(startAttrs...),
		}

		ctx, span := tracer.Start(ctx, name, startOpts...)
		defer span.End()

		err := invoker(ctx, method, req, reply, cc, callOpts...)

		span.SetStatus(status.Convert(err).Code())

		return err
	}
}

type UnaryClientInterceptorOption func(*unaryClientInterceptorConfig)

func UnaryClientInterceptorWithTracer(tracer trace.Tracer) UnaryClientInterceptorOption {
	return func(c *unaryClientInterceptorConfig) {
		c.tracer = tracer
	}
}

type unaryClientInterceptorConfig struct {
	tracer trace.Tracer
}

func newUnaryClientInterceptorConfig(opts ...UnaryClientInterceptorOption) unaryClientInterceptorConfig {
	var c unaryClientInterceptorConfig
	defaultOpts := []UnaryClientInterceptorOption{
		UnaryClientInterceptorWithTracer(trace.NoopTracer{}),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	return c
}

func NewStreamClientInterceptor(opts ...StreamClientInterceptorOption) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, callOpts...)
	}
}

type StreamClientInterceptorOption func(*streamClientInterceptorConfig)

type streamClientInterceptorConfig struct{}
