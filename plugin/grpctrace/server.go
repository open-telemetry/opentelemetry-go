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

func NewUnaryServerInterceptor(opts ...UnaryServerInterceptorOption) grpc.UnaryServerInterceptor {
	c := newUnaryServerInterceptorConfig(opts...)
	tracer := c.tracer

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var name string
		var service string

		if info != nil {
			name = nameFromMethod(info.FullMethod)
			service = serviceFromMethod(info.FullMethod)
		}

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

		res, err := handler(ctx, req)

		span.SetStatus(status.Convert(err).Code())

		return res, err
	}
}

type UnaryServerInterceptorOption func(*unaryServerInterceptorConfig)

func UnaryServerInterceptorWithTracer(tracer trace.Tracer) UnaryServerInterceptorOption {
	return func(c *unaryServerInterceptorConfig) {
		c.tracer = tracer
	}
}

type unaryServerInterceptorConfig struct {
	tracer trace.Tracer
}

func newUnaryServerInterceptorConfig(opts ...UnaryServerInterceptorOption) unaryServerInterceptorConfig {
	var c unaryServerInterceptorConfig
	defaultOpts := []UnaryServerInterceptorOption{
		UnaryServerInterceptorWithTracer(trace.NoopTracer{}),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	return c
}

func NewStreamServerInterceptor(opts ...StreamServerInterceptorOption) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

type StreamServerInterceptorOption func(*streamServerInterceptorConfig)

type streamServerInterceptorConfig struct{}
