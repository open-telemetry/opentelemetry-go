// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func TestRecorderLogger(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []Option

		loggerName    string
		loggerOptions []log.LoggerOption

		expectedLogger log.Logger
	}{
		{
			name: "provides a default logger",

			expectedLogger: &Recorder{},
		},
		{
			name: "provides a logger with a configured scope",

			loggerName: "test",
			loggerOptions: []log.LoggerOption{
				log.WithInstrumentationVersion("logtest v42"),
				log.WithSchemaURL("https://example.com"),
			},

			expectedLogger: &Recorder{
				Scope: Scope{
					Name:      "test",
					Version:   "logtest v42",
					SchemaURL: "https://example.com",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRecorder(tt.options...).Logger(tt.loggerName, tt.loggerOptions...)
			assert.Equal(t, tt.expectedLogger, l)
		})
	}
}

func TestRecorderLoggerCreatesNewStruct(t *testing.T) {
	r := NewRecorder()
	assert.NotEqual(t, r, r.Logger("test"))
}

func TestRecorderEnabled(t *testing.T) {
	for _, tt := range []struct {
		name        string
		options     []Option
		ctx         context.Context
		buildRecord func() log.Record

		isEnabled bool
	}{
		{
			name: "the default option enables unset levels",
			ctx:  context.Background(),
			buildRecord: func() log.Record {
				return log.Record{}
			},

			isEnabled: true,
		},
		{
			name: "with a minimum severity set disables",
			options: []Option{
				WithMinSeverity(log.SeverityWarn1),
			},
			ctx: context.Background(),
			buildRecord: func() log.Record {
				return log.Record{}
			},

			isEnabled: false,
		},
		{
			name: "with a context that forces an enabled recorder",
			options: []Option{
				WithMinSeverity(log.SeverityWarn1),
			},
			ctx: ContextWithEnabledRecorder(context.Background()),
			buildRecord: func() log.Record {
				return log.Record{}
			},

			isEnabled: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e := NewRecorder(tt.options...).Enabled(tt.ctx, tt.buildRecord())
			assert.Equal(t, tt.isEnabled, e)
		})
	}
}

func TestRecorderEmitAndReset(t *testing.T) {
	r := NewRecorder()
	assert.Len(t, r.Result(), 0)
	r.Emit(context.Background(), log.Record{})
	assert.Len(t, r.Result(), 1)

	r.Reset()
	assert.Len(t, r.Result(), 0)
}
