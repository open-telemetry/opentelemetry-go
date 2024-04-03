// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"sync"
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

			expectedLogger: &Recorder{
				currentScopeRecord: &ScopeRecords{},
			},
		},
		{
			name: "provides a logger with a configured scope",

			loggerName: "test",
			loggerOptions: []log.LoggerOption{
				log.WithInstrumentationVersion("logtest v42"),
				log.WithSchemaURL("https://example.com"),
			},

			expectedLogger: &Recorder{
				currentScopeRecord: &ScopeRecords{
					Name:      "test",
					Version:   "logtest v42",
					SchemaURL: "https://example.com",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRecorder(tt.options...).Logger(tt.loggerName, tt.loggerOptions...)
			// unset enabledFn to allow comparison
			l.(*Recorder).enabledFn = nil

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
			name: "the default option enables every log entry",
			ctx:  context.Background(),
			buildRecord: func() log.Record {
				return log.Record{}
			},

			isEnabled: true,
		},
		{
			name: "with everything disabled",
			options: []Option{
				WithEnabledFunc(func(context.Context, log.Record) bool {
					return false
				}),
			},
			ctx: context.Background(),
			buildRecord: func() log.Record {
				return log.Record{}
			},

			isEnabled: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e := NewRecorder(tt.options...).Enabled(tt.ctx, tt.buildRecord())
			assert.Equal(t, tt.isEnabled, e)
		})
	}
}

func TestRecorderEnabledFnUnset(t *testing.T) {
	r := &Recorder{}
	assert.True(t, r.Enabled(context.Background(), log.Record{}))
}

func TestRecorderEmitAndReset(t *testing.T) {
	r := NewRecorder()
	assert.Len(t, r.Result()[0].Records, 0)
	r.Emit(context.Background(), log.Record{})
	assert.Len(t, r.Result()[0].Records, 1)

	r.Reset()
	assert.Len(t, r.Result()[0].Records, 0)
}

func TestRecorderConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := NewRecorder()

	for i := 0; i < goRoutineN; i++ {
		go func() {
			defer wg.Done()

			nr := r.Logger("test")
			nr.Emit(context.Background(), log.Record{})

			r.Result()
			r.Reset()
		}()
	}

	wg.Wait()
}
