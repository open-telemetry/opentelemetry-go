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

package zipkin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	zkmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	collectorURL = "http://localhost:9411/api/v2/spans"
	serviceName  = "zipkin-test"
)

func TestInstallNewPipeline(t *testing.T) {
	err := InstallNewPipeline(
		collectorURL,
		serviceName,
	)
	assert.NoError(t, err)
	assert.IsType(t, &sdktrace.TracerProvider{}, otel.GetTracerProvider())
}

func TestNewExportPipeline(t *testing.T) {
	testCases := []struct {
		name                                  string
		options                               []Option
		testSpanSampling, spanShouldBeSampled bool
	}{
		{
			name: "simple pipeline",
		},
		{
			name: "always on",
			options: []Option{
				WithSDK(&sdktrace.Config{
					DefaultSampler: sdktrace.AlwaysSample(),
				}),
			},
			testSpanSampling:    true,
			spanShouldBeSampled: true,
		},
		{
			name: "never",
			options: []Option{
				WithSDK(&sdktrace.Config{
					DefaultSampler: sdktrace.NeverSample(),
				}),
			},
			testSpanSampling:    true,
			spanShouldBeSampled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tp, err := NewExportPipeline(
				collectorURL,
				serviceName,
				tc.options...,
			)
			assert.NoError(t, err)
			assert.NotEqual(t, tp, otel.GetTracerProvider())

			if tc.testSpanSampling {
				_, span := tp.Tracer("zipkin test").Start(context.Background(), tc.name)
				spanCtx := span.SpanContext()
				assert.Equal(t, tc.spanShouldBeSampled, spanCtx.IsSampled())
				span.End()
			}
		})
	}
}

func TestNewRawExporter(t *testing.T) {
	exp, err := NewRawExporter(
		collectorURL,
		serviceName,
	)

	assert.NoError(t, err)
	assert.EqualValues(t, serviceName, exp.serviceName)
}

func TestNewRawExporterShouldFailInvalidCollectorURL(t *testing.T) {
	var (
		exp *Exporter
		err error
	)

	// cannot be empty
	exp, err = NewRawExporter(
		"",
		serviceName,
	)

	assert.Error(t, err)
	assert.EqualError(t, err, "collector URL cannot be empty")
	assert.Nil(t, exp)

	// invalid URL
	exp, err = NewRawExporter(
		"localhost",
		serviceName,
	)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid collector URL")
	assert.Nil(t, exp)
}

type mockZipkinCollector struct {
	t       *testing.T
	url     string
	closing bool
	server  *http.Server
	wg      *sync.WaitGroup

	lock   sync.RWMutex
	models []zkmodel.SpanModel
}

func startMockZipkinCollector(t *testing.T) *mockZipkinCollector {
	collector := &mockZipkinCollector{
		t:       t,
		closing: false,
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	collector.url = fmt.Sprintf("http://%s", listener.Addr().String())
	server := &http.Server{
		Handler: http.HandlerFunc(collector.handler),
	}
	collector.server = server
	wg := &sync.WaitGroup{}
	wg.Add(1)
	collector.wg = wg
	go func() {
		err := server.Serve(listener)
		require.True(t, collector.closing)
		require.Equal(t, http.ErrServerClosed, err)
		wg.Done()
	}()

	return collector
}

func (c *mockZipkinCollector) handler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	require.NoError(c.t, err)
	var models []zkmodel.SpanModel
	err = json.Unmarshal(jsonBytes, &models)
	require.NoError(c.t, err)
	// for some reason we may get the nonUTC timestamps in models,
	// fix that
	for midx := range models {
		models[midx].Timestamp = models[midx].Timestamp.UTC()
		for aidx := range models[midx].Annotations {
			models[midx].Annotations[aidx].Timestamp = models[midx].Annotations[aidx].Timestamp.UTC()
		}
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.models = append(c.models, models...)
}

func (c *mockZipkinCollector) Close() {
	if c.closing {
		return
	}
	c.closing = true
	server := c.server
	c.server = nil
	require.NoError(c.t, server.Shutdown(context.Background()))
	c.wg.Wait()
}

func (c *mockZipkinCollector) ModelsLen() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.models)
}

func (c *mockZipkinCollector) StealModels() []zkmodel.SpanModel {
	c.lock.Lock()
	defer c.lock.Unlock()
	models := c.models
	c.models = nil
	return models
}

type logStore struct {
	T        *testing.T
	Messages []string
}

func (s *logStore) Write(p []byte) (n int, err error) {
	msg := (string)(p)
	if s.T != nil {
		s.T.Logf("%s", msg)
	}
	s.Messages = append(s.Messages, msg)
	return len(p), nil
}

func logStoreLogger(s *logStore) *log.Logger {
	return log.New(s, "", 0)
}

func TestExportSpans(t *testing.T) {
	spans := []*export.SpanData{
		// parent
		{
			SpanContext: trace.SpanContext{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			},
			ParentSpanID:  trace.SpanID{},
			SpanKind:      trace.SpanKindServer,
			Name:          "foo",
			StartTime:     time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			EndTime:       time.Date(2020, time.March, 11, 19, 25, 0, 0, time.UTC),
			Attributes:    nil,
			MessageEvents: nil,
			StatusCode:    codes.Error,
			StatusMessage: "404, file not found",
		},
		// child
		{
			SpanContext: trace.SpanContext{
				TraceID: trace.TraceID{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
				SpanID:  trace.SpanID{0xDF, 0xDE, 0xDD, 0xDC, 0xDB, 0xDA, 0xD9, 0xD8},
			},
			ParentSpanID:  trace.SpanID{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8},
			SpanKind:      trace.SpanKindServer,
			Name:          "bar",
			StartTime:     time.Date(2020, time.March, 11, 19, 24, 15, 0, time.UTC),
			EndTime:       time.Date(2020, time.March, 11, 19, 24, 45, 0, time.UTC),
			Attributes:    nil,
			MessageEvents: nil,
			StatusCode:    codes.Error,
			StatusMessage: "403, forbidden",
		},
	}
	models := []zkmodel.SpanModel{
		// model of parent
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xfffefdfcfbfaf9f8),
				ParentID: nil,
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "foo",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 0, 0, time.UTC),
			Duration:  time.Minute,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "exporter-test",
			},
			RemoteEndpoint: nil,
			Annotations:    nil,
			Tags: map[string]string{
				"otel.status_code":        "Error",
				"otel.status_description": "404, file not found",
			},
		},
		// model of child
		{
			SpanContext: zkmodel.SpanContext{
				TraceID: zkmodel.TraceID{
					High: 0x001020304050607,
					Low:  0x8090a0b0c0d0e0f,
				},
				ID:       zkmodel.ID(0xdfdedddcdbdad9d8),
				ParentID: zkmodelIDPtr(0xfffefdfcfbfaf9f8),
				Debug:    false,
				Sampled:  nil,
				Err:      nil,
			},
			Name:      "bar",
			Kind:      "SERVER",
			Timestamp: time.Date(2020, time.March, 11, 19, 24, 15, 0, time.UTC),
			Duration:  30 * time.Second,
			Shared:    false,
			LocalEndpoint: &zkmodel.Endpoint{
				ServiceName: "exporter-test",
			},
			RemoteEndpoint: nil,
			Annotations:    nil,
			Tags: map[string]string{
				"otel.status_code":        "Error",
				"otel.status_description": "403, forbidden",
			},
		},
	}
	require.Len(t, models, len(spans))
	collector := startMockZipkinCollector(t)
	defer collector.Close()
	ls := &logStore{T: t}
	logger := logStoreLogger(ls)
	exporter, err := NewRawExporter(collector.url, "exporter-test", WithLogger(logger))
	require.NoError(t, err)
	ctx := context.Background()
	require.Len(t, ls.Messages, 0)
	require.NoError(t, exporter.ExportSpans(ctx, spans[0:1]))
	require.Len(t, ls.Messages, 2)
	require.Contains(t, ls.Messages[0], "send a POST request")
	require.Contains(t, ls.Messages[1], "zipkin responded")
	ls.Messages = nil
	require.NoError(t, exporter.ExportSpans(ctx, nil))
	require.Len(t, ls.Messages, 1)
	require.Contains(t, ls.Messages[0], "no spans to export")
	ls.Messages = nil
	require.NoError(t, exporter.ExportSpans(ctx, spans[1:2]))
	require.Contains(t, ls.Messages[0], "send a POST request")
	require.Contains(t, ls.Messages[1], "zipkin responded")
	checkFunc := func() bool {
		return collector.ModelsLen() == len(models)
	}
	require.Eventually(t, checkFunc, time.Second, 10*time.Millisecond)
	require.Equal(t, models, collector.StealModels())
}

func TestExporterShutdownHonorsTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	exp, err := NewRawExporter(collectorURL, serviceName)
	require.NoError(t, err)

	innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
	defer innerCancel()
	<-innerCtx.Done()
	assert.Errorf(t, exp.Shutdown(innerCtx), context.DeadlineExceeded.Error())
}

func TestExporterShutdownHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	exp, err := NewRawExporter(collectorURL, serviceName)
	require.NoError(t, err)

	innerCtx, innerCancel := context.WithCancel(ctx)
	innerCancel()
	assert.Errorf(t, exp.Shutdown(innerCtx), context.Canceled.Error())
}

func TestErrorOnExportShutdownExporter(t *testing.T) {
	exp, err := NewRawExporter(collectorURL, serviceName)
	require.NoError(t, err)
	assert.NoError(t, exp.Shutdown(context.Background()))
	assert.NoError(t, exp.ExportSpans(context.Background(), nil))
}
