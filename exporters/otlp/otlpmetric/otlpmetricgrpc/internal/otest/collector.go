// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/otest"

import (
	"context" // nolint:depguard  // This is for testing.
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

// ExportResult represents an export response.
type ExportResult struct {
	Response *collpb.ExportMetricsServiceResponse
	Err      error
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

// Dump returns all added ResourceMetrics and clears the storage.
func (s *Storage) Dump() []*mpb.ResourceMetrics {
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

	resultCh <-chan ExportResult
	listener net.Listener
	srv      *grpc.Server
}

// NewGRPCCollector returns a *GRPCCollector that is listening at the provided
// endpoint.
//
// If endpoint is an empty string, the returned collector will be listening on
// the localhost interface at an OS chosen port.
//
// If errCh is not nil, the collector will respond to Export calls with errors
// sent on that channel. This means that if errCh is not nil Export calls will
// block until an error is received.
func NewGRPCCollector(endpoint string, resultCh <-chan ExportResult) (*GRPCCollector, error) {
	if endpoint == "" {
		endpoint = "localhost:0"
	}

	c := &GRPCCollector{
		storage:  NewStorage(),
		resultCh: resultCh,
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

	if c.resultCh != nil {
		r := <-c.resultCh
		if r.Response == nil {
			return &collpb.ExportMetricsServiceResponse{}, r.Err
		}
		return r.Response, r.Err
	}
	return &collpb.ExportMetricsServiceResponse{}, nil
}
