// Copyright 2020, OpenTelemetry Authors
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

package otlp_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"
)

func makeMockCollector(t *testing.T) *mockCol {
	return &mockCol{t: t, wg: new(sync.WaitGroup)}
}

type mockCol struct {
	t *testing.T

	spans []*tracepb.Span
	mu    sync.Mutex
	wg    *sync.WaitGroup

	address  string
	stopFunc func() error
	stopOnce sync.Once
}

var _ coltracepb.TraceServiceServer = (*mockCol)(nil)

func (mc *mockCol) Export(ctx context.Context, exp *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	resourceSpans := exp.GetResourceSpans()
	// TODO (rghetia): handle Resources
	for _, rs := range resourceSpans {
		mc.spans = append(mc.spans, rs.Spans...)
	}
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

var errAlreadyStopped = fmt.Errorf("already stopped")

func (mc *mockCol) stop() error {
	var err = errAlreadyStopped
	mc.stopOnce.Do(func() {
		if mc.stopFunc != nil {
			err = mc.stopFunc()
		}
	})
	// Give it sometime to shutdown.
	<-time.After(160 * time.Millisecond)
	mc.mu.Lock()
	mc.wg.Wait()
	mc.mu.Unlock()
	return err
}

// runMockCol is a helper function to create a mockCol
func runMockCol(t *testing.T) *mockCol {
	return runMockColAtAddr(t, "localhost:0")
}

func runMockColAtAddr(t *testing.T, addr string) *mockCol {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to get an address: %v", err)
	}

	srv := grpc.NewServer()
	mc := makeMockCollector(t)
	coltracepb.RegisterTraceServiceServer(srv, mc)
	go func() {
		_ = srv.Serve(ln)
	}()

	deferFunc := func() error {
		srv.Stop()
		return ln.Close()
	}

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	mc.address = "localhost:" + collectorPortStr
	mc.stopFunc = deferFunc

	return mc
}

func (mc *mockCol) getSpans() []*tracepb.Span {
	mc.mu.Lock()
	spans := append([]*tracepb.Span{}, mc.spans...)
	mc.mu.Unlock()

	return spans
}
