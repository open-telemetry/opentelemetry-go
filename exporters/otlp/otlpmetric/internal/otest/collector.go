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

//go:build go1.18
// +build go1.18

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	collpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

// Collector is the collection target a Client sends metric uploads to.
type Collector interface {
	Collect() *Storage
}

// Storage stores uploaded OTLP metric data in their proto form.
type Storage struct {
	dataMu sync.Mutex
	data   []*mpb.ResourceMetrics
}

// NewStorage returns a configure storage ready to store received requests.
func NewStorage() *Storage {
	return &Storage{}
}

// Add adds the request to the Storage.
func (s *Storage) Add(request *collpb.ExportMetricsServiceRequest) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	s.data = append(s.data, request.ResourceMetrics...)
}

// dump returns all added ResourceMetrics and clears the storage.
func (s *Storage) dump() []*mpb.ResourceMetrics {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var data []*mpb.ResourceMetrics
	data, s.data = s.data, []*mpb.ResourceMetrics{}
	return data
}

// GRPCCollector is an OTLP gRPC server that collects all requests it receives.
type GRPCCollector struct {
	collpb.UnimplementedMetricsServiceServer

	headersMu sync.Mutex
	headers   metadata.MD
	storage   *Storage

	errCh    <-chan error
	listener net.Listener
	srv      *grpc.Server
}

// NewGRPCCollector returns a *GRPCCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listeing on
// the localhost interface at an OS chosen port.
//
// If errCh is not nil, the collector will respond to Export calls with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func NewGRPCCollector(endpoint string, errCh <-chan error) (*GRPCCollector, error) {
	if endpoint == "" {
		endpoint = "localhost:0"
	}

	c := &GRPCCollector{
		storage: NewStorage(),
		errCh:   errCh,
	}

	var err error
	c.listener, err = net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	c.srv = grpc.NewServer()
	collpb.RegisterMetricsServiceServer(c.srv, c)
	go func() { _ = c.srv.Serve(c.listener) }()

	return c, nil
}

// Shutdown shuts down the gRPC server closing all open connections and
// listeners immediately.
func (c *GRPCCollector) Shutdown() { c.srv.Stop() }

// Addr returns the net.Addr c is listening at.
func (c *GRPCCollector) Addr() net.Addr {
	return c.listener.Addr()
}

// Collect returns the Storage holding all collected requests.
func (c *GRPCCollector) Collect() *Storage {
	return c.storage
}

// Headers returns the headers received for all requests.
func (c *GRPCCollector) Headers() map[string][]string {
	// Makes a copy.
	c.headersMu.Lock()
	defer c.headersMu.Unlock()
	return metadata.Join(c.headers)
}

// Export handles the export req.
func (c *GRPCCollector) Export(ctx context.Context, req *collpb.ExportMetricsServiceRequest) (*collpb.ExportMetricsServiceResponse, error) {
	c.storage.Add(req)

	if h, ok := metadata.FromIncomingContext(ctx); ok {
		c.headersMu.Lock()
		c.headers = metadata.Join(c.headers, h)
		c.headersMu.Unlock()
	}

	var err error
	if c.errCh != nil {
		err = <-c.errCh
	}
	return &collpb.ExportMetricsServiceResponse{}, err
}
