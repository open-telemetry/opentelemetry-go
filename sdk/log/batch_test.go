// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
)

func TestNewBatchingProcessorConfiguration(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name                   string
		envars                 map[string]string
		options                []BatchingOption
		wantExporter           Exporter
		wantMaxQueueSize       int
		wantExportInterval     time.Duration
		wantExportTimeout      time.Duration
		wantExportMaxBatchSize int
	}{
		{
			name:                   "Defaults",
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
		},
		{
			name: "Options",
			options: []BatchingOption{
				WithMaxQueueSize(1),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       1,
			wantExportInterval:     time.Microsecond,
			wantExportTimeout:      time.Hour,
			wantExportMaxBatchSize: 2,
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       1,
			wantExportInterval:     100 * time.Millisecond,
			wantExportTimeout:      1000 * time.Millisecond,
			wantExportMaxBatchSize: 10,
		},
		{
			name: "InvalidOptions",
			options: []BatchingOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "-1",
				envarExpInterval:     "-1",
				envarExpTimeout:      "-1",
				envarExpMaxBatchSize: "-1",
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
		},
		{
			name: "Precedence",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			options: []BatchingOption{
				// These override the environment variables.
				WithMaxQueueSize(3),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       3,
			wantExportInterval:     time.Microsecond,
			wantExportTimeout:      time.Hour,
			wantExportMaxBatchSize: 2,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}

			b := NewBatchingProcessor(nil, tc.options...)
			t.Cleanup(func() {
				assert.NoError(t, b.Shutdown(context.Background()))
			})
			assert.Equal(t, tc.wantExporter, b.exporter, "exporter")
			assert.Equal(t, tc.wantExportInterval, b.exportInterval, "exportInterval")
			assert.Equal(t, tc.wantExportTimeout, b.exportTimeout, "exportTimeout")
			assert.Equal(t, tc.wantExportMaxBatchSize, b.exportMaxBatchSize, "exportMaxBatchSize")
			assert.Equal(t, tc.wantMaxQueueSize, b.maxQueueSize, "maxQueueSize")
		})
	}
}
