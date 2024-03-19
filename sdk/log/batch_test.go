// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
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
		name    string
		envars  map[string]string
		options []BatchingOption
		want    *BatchingProcessor
	}{
		{
			name: "Defaults",
			want: &BatchingProcessor{
				exporter:           defaultNoopExporter,
				maxQueueSize:       dfltMaxQSize,
				exportInterval:     dfltExpInterval,
				exportTimeout:      dfltExpTimeout,
				exportMaxBatchSize: dfltExpMaxBatchSize,
			},
		},
		{
			name: "Options",
			options: []BatchingOption{
				WithMaxQueueSize(1),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			want: &BatchingProcessor{
				exporter:           defaultNoopExporter,
				maxQueueSize:       1,
				exportInterval:     time.Microsecond,
				exportTimeout:      time.Hour,
				exportMaxBatchSize: 2,
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			want: &BatchingProcessor{
				exporter:           defaultNoopExporter,
				maxQueueSize:       1,
				exportInterval:     100 * time.Millisecond,
				exportTimeout:      1000 * time.Millisecond,
				exportMaxBatchSize: 10,
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "invalid envarMaxQSize",
				envarExpInterval:     "invalid envarExpInterval",
				envarExpTimeout:      "invalid envarExpTimeout",
				envarExpMaxBatchSize: "invalid envarExpMaxBatchSize",
			},
			want: &BatchingProcessor{
				exporter:           defaultNoopExporter,
				maxQueueSize:       dfltMaxQSize,
				exportInterval:     dfltExpInterval,
				exportTimeout:      dfltExpTimeout,
				exportMaxBatchSize: dfltExpMaxBatchSize,
			},
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
			want: &BatchingProcessor{
				exporter:           defaultNoopExporter,
				maxQueueSize:       3,
				exportInterval:     time.Microsecond,
				exportTimeout:      time.Hour,
				exportMaxBatchSize: 2,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, NewBatchingProcessor(nil, tc.options...))
		})
	}
}
