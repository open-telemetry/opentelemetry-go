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

func TestNewBatchingConfig(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name    string
		envars  map[string]string
		options []BatchingOption
		want    batchingConfig
	}{
		{
			name: "Defaults",
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
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
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
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
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(100 * time.Millisecond),
				expTimeout:      newSetting(1000 * time.Millisecond),
				expMaxBatchSize: newSetting(10),
			},
		},
		{
			name: "InvalidOptions",
			options: []BatchingOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "-1",
				envarExpInterval:     "-1",
				envarExpTimeout:      "-1",
				envarExpMaxBatchSize: "-1",
			},
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
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
			want: batchingConfig{
				maxQSize:        newSetting(3),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newBatchingConfig(tc.options))
		})
	}
}
