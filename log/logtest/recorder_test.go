// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestRecorderLogger(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options []Option

		loggerName    string
		loggerOptions []log.LoggerOption

		wantLogger log.Logger
	}{
		{
			name: "provides a default logger",

			wantLogger: &logger{
				scopeRecord: &ScopeRecords{},
			},
		},
		{
			name: "provides a logger with a configured scope",

			loggerName: "test",
			loggerOptions: []log.LoggerOption{
				log.WithInstrumentationVersion("logtest v42"),
				log.WithSchemaURL("https://example.com"),
				log.WithInstrumentationAttributes(attribute.String("foo", "bar")),
			},

			wantLogger: &logger{
				scopeRecord: &ScopeRecords{
					Name:       "test",
					Version:    "logtest v42",
					SchemaURL:  "https://example.com",
					Attributes: attribute.NewSet(attribute.String("foo", "bar")),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRecorder(tt.options...).Logger(tt.loggerName, tt.loggerOptions...)
			// unset enabledFn to allow comparison
			l.(*logger).enabledFn = nil

			assert.Equal(t, tt.wantLogger, l)
		})
	}
}

func TestRecorderLoggerCreatesNewStruct(t *testing.T) {
	r := &Recorder{}
	assert.NotEqual(t, r, r.Logger("test"))
}

func TestLoggerEnabled(t *testing.T) {
	for _, tt := range []struct {
		name          string
		options       []Option
		ctx           context.Context
		enabledParams log.EnabledParameters
		want          bool
	}{
		{
			name: "the default option enables every log entry",
			ctx:  context.Background(),
			want: true,
		},
		{
			name: "with everything disabled",
			options: []Option{
				WithEnabledFunc(func(context.Context, log.EnabledParameters) bool {
					return false
				}),
			},
			ctx:  context.Background(),
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			e := NewRecorder(tt.options...).Logger("test").Enabled(tt.ctx, tt.enabledParams)
			assert.Equal(t, tt.want, e)
		})
	}
}

func TestLoggerEnabledFnUnset(t *testing.T) {
	r := &logger{}
	assert.True(t, r.Enabled(context.Background(), log.EnabledParameters{}))
}

func TestRecorderEmitAndReset(t *testing.T) {
	r := NewRecorder()
	l := r.Logger("test")
	assert.Empty(t, r.Result()[0].Records)

	r1 := log.Record{}
	r1.SetSeverity(log.SeverityInfo)
	ctx := context.Background()

	l.Emit(ctx, r1)
	assert.Equal(t, []EmittedRecord{
		{r1, ctx},
	}, r.Result()[0].Records)

	nl := r.Logger("test")
	assert.Empty(t, r.Result()[1].Records)

	r2 := log.Record{}
	r2.SetSeverity(log.SeverityError)
	// We want a non-background context here so it's different from `ctx`.
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	nl.Emit(ctx2, r2)
	assert.Len(t, r.Result()[0].Records, 1)
	AssertRecordEqual(t, r.Result()[0].Records[0].Record, r1)
	assert.Equal(t, r.Result()[0].Records[0].Context(), ctx)

	assert.Len(t, r.Result()[1].Records, 1)
	AssertRecordEqual(t, r.Result()[1].Records[0].Record, r2)
	assert.Equal(t, r.Result()[1].Records[0].Context(), ctx2)

	r.Reset()
	assert.Empty(t, r.Result()[0].Records)
	assert.Empty(t, r.Result()[1].Records)
}

func TestRecorderConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := &Recorder{}

	for i := 0; i < goRoutineN; i++ {
		go func() {
			defer wg.Done()

			nr := r.Logger("test")
			nr.Enabled(context.Background(), log.EnabledParameters{})
			nr.Emit(context.Background(), log.Record{})

			r.Result()
			r.Reset()
		}()
	}

	wg.Wait()
}
