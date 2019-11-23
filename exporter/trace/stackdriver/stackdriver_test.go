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

package stackdriver_test

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	tracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/exporter/trace/stackdriver"
	"go.opentelemetry.io/otel/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type mockTraceServer struct {
	tracepb.TraceServiceServer
	mu            sync.Mutex
	spansUploaded []*tracepb.Span
	delay         time.Duration
}

func (s *mockTraceServer) BatchWriteSpans(ctx context.Context, req *tracepb.BatchWriteSpansRequest) (*emptypb.Empty, error) {
	var err error
	s.mu.Lock()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-time.After(s.delay):
		s.spansUploaded = append(s.spansUploaded, req.Spans...)
	}
	s.mu.Unlock()
	return &emptypb.Empty{}, err
}

func (s *mockTraceServer) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.spansUploaded)
}

// clientOpt is the option tests should use to connect to the test server.
// It is initialized by TestMain.
var clientOpt []option.ClientOption

var (
	mockTrace mockTraceServer
)

func TestMain(m *testing.M) {
	flag.Parse()

	serv := grpc.NewServer()
	tracepb.RegisterTraceServiceServer(serv, &mockTrace)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		_ = serv.Serve(lis)
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	clientOpt = []option.ClientOption{option.WithGRPCConn(conn)}

	os.Exit(m.Run())
}

func TestExporter_ExportSpans(t *testing.T) {
	// Initial test precondition
	mockTrace.spansUploaded = nil
	mockTrace.delay = 0

	// Create StackDriver Exporter
	exp, err := stackdriver.NewExporter(
		stackdriver.WithProjectID("PROJECT_ID_NOT_REAL"),
		stackdriver.WithTraceClientOptions(clientOpt),
	)
	assert.NoError(t, err)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exp, // add following two options to ensure flush
			sdktrace.WithScheduleDelayMillis(1),
			sdktrace.WithMaxExportBatchSize(1),
		))
	assert.NoError(t, err)

	global.SetTraceProvider(tp)
	_, span := global.TraceProvider().Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()
	assert.True(t, span.SpanContext().IsValid())

	// wait exporter to flush
	time.Sleep(20 * time.Millisecond)
	assert.EqualValues(t, 1, mockTrace.len())
}

func TestExporter_Timeout(t *testing.T) {
	// Initial test precondition
	mockTrace.spansUploaded = nil
	mockTrace.delay = 20 * time.Millisecond
	var exportErrors []error

	// Create StackDriver Exporter
	exp, err := stackdriver.NewExporter(
		stackdriver.WithProjectID("PROJECT_ID_NOT_REAL"),
		stackdriver.WithTraceClientOptions(clientOpt),
		stackdriver.WithTimeout(1*time.Millisecond),
		stackdriver.WithOnError(func(err error) {
			exportErrors = append(exportErrors, err)
		}),
	)
	assert.NoError(t, err)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exp))
	assert.NoError(t, err)

	global.SetTraceProvider(tp)
	_, span := global.TraceProvider().Tracer("test-tracer").Start(context.Background(), "test-span")
	span.End()
	assert.True(t, span.SpanContext().IsValid())

	assert.EqualValues(t, 0, mockTrace.len())
	if got, want := len(exportErrors), 1; got != want {
		t.Fatalf("len(exportErrors) = %q; want %q", got, want)
	}
	if got, want := exportErrors[0].Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded"; got != want {
		t.Fatalf("err.Error() = %q; want %q", got, want)
	}
}
