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

package otlp

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
)

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	orig := closeStopCh
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer func() {
		cancel()
		closeStopCh = orig
	}()
	closeStopCh = func(stopCh chan bool) {
		go func() {
			<-ctx.Done()
			close(stopCh)
		}()
	}

	e := NewUnstartedExporter(NewConnections(DefaultConnectionOptions...))
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	innerCtx, innerCancel := context.WithTimeout(ctx, time.Microsecond)
	if err := e.Shutdown(innerCtx); err == nil {
		t.Error("expected context DeadlineExceeded error, got nil")
	} else if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context DeadlineExceeded error, got %v", err)
	}
	innerCancel()
}

func TestExporterShutdownHonorsCancel(t *testing.T) {
	orig := closeStopCh
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer func() {
		cancel()
		closeStopCh = orig
	}()
	closeStopCh = func(stopCh chan bool) {
		go func() {
			<-ctx.Done()
			close(stopCh)
		}()
	}

	e := NewUnstartedExporter(NewConnections(DefaultConnectionOptions...))
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	var innerCancel context.CancelFunc
	ctx, innerCancel = context.WithCancel(ctx)
	innerCancel()
	if err := e.Shutdown(ctx); err == nil {
		t.Error("expected context canceled error, got nil")
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error, got %v", err)
	}
}

func TestExporterShutdownNoError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	e := NewUnstartedExporter(NewConnections(DefaultConnectionOptions...))
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}

	if err := e.Shutdown(ctx); err != nil {
		t.Errorf("shutdown errored: expected nil, got %v", err)
	}
}

func TestExporterOnlyStartsMetricsConnection(t *testing.T) {

	collectorAddr := ":9081"
	colDeferFunc := runUnimplementedCollectorAtAddress(t, collectorAddr)
	defer func() {
		_ = colDeferFunc()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	config := NewConnections().
		SetMetricOptions(WithAddress(collectorAddr), WithInsecure())
	e := NewUnstartedExporter(config)
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}
	defer func() {
		_ = e.Shutdown(ctx)
	}()

	metricsConnErr := e.metricsConnection.lastConnectError()
	if metricsConnErr != nil {
		t.Errorf("metrics connection error out: %v", metricsConnErr)
	}

	if e.tracesConnection != nil {
		t.Errorf("expected traces connection not to start")
	}

}

func TestExporterOnlyStartsTracesConnection(t *testing.T) {

	collectorAddr := ":9081"
	colDeferFunc := runUnimplementedCollectorAtAddress(t, collectorAddr)
	defer func() {
		_ = colDeferFunc()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	config := NewConnections().
		SetTraceOptions(WithAddress(collectorAddr), WithInsecure())
	e := NewUnstartedExporter(config)
	if err := e.Start(); err != nil {
		t.Fatalf("failed to start exporter: %v", err)
	}
	defer func() {
		_ = e.Shutdown(ctx)
	}()

	tracesConnectionError := e.tracesConnection.lastConnectError()
	if tracesConnectionError != nil {
		t.Errorf("traces connection error out: %v", tracesConnectionError)
	}

	if e.metricsConnection != nil {
		t.Errorf("expected metrics connection not to start")
	}

}

func runUnimplementedCollectorAtAddress(t *testing.T, addr string) func() error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to get an address: %v", err)
	}

	srv := grpc.NewServer()
	coltracepb.RegisterTraceServiceServer(srv, new(coltracepb.UnimplementedTraceServiceServer))
	colmetricpb.RegisterMetricsServiceServer(srv, new(colmetricpb.UnimplementedMetricsServiceServer))
	go func() {
		_ = srv.Serve(ln)
	}()

	deferFunc := func() error {
		srv.Stop()
		return ln.Close()
	}

	return deferFunc
}
