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

package stackdriver

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	tracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testUploader struct {
	mu            sync.Mutex
	spansUploaded []*tracepb.Span
}

// testUploadSpans assigned to uploadFn when in test.
func (c *testUploader) testUploadSpans(ctx context.Context, spans []*tracepb.Span) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spansUploaded = append(c.spansUploaded, spans...)
}

func (c *testUploader) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.spansUploaded)
}

type mockTraceServer struct {
	tracepb.TraceServiceServer
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
	go serv.Serve(lis)

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	clientOpt = []option.ClientOption{option.WithGRPCConn(conn)}

	os.Exit(m.Run())
}

func TestNewExporter(t *testing.T) {
	const projectID = "project-id"

	// Create SD Exporter
	exp, err := NewExporter(
		WithProjectID(projectID),
		WithTraceClientOptions(clientOpt),
	)

	assert.NoError(t, err)
	assert.EqualValues(t, projectID, exp.traceExporter.projectID)
}

func TestExporter_ExportSpans(t *testing.T) {
	// Create StackDriver Exporter
	exp, err := NewExporter(
		WithProjectID("PROJECT_ID_NOT_REAL"),
		WithTraceClientOptions(clientOpt),
	)
	assert.NoError(t, err)

	tu := &testUploader{}
	exp.traceExporter.uploadFn = tu.testUploadSpans

	assert.NoError(t, err)

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(exp, // add following two options to ensure flush
			sdktrace.WithScheduleDelayMillis(1),
			sdktrace.WithMaxExportBatchSize(1),
		))

	assert.NoError(t, err)

	global.SetTraceProvider(tp)
	_, span := global.TraceProvider().GetTracer("test-tracer").Start(context.Background(), "test-span")
	span.End()

	assert.True(t, span.SpanContext().IsValid())

	// wait exporter to flush
	time.Sleep(20 * time.Millisecond)
	assert.EqualValues(t, 1, tu.len())
}
