// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"sync"
	"testing"
	"time"

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

		want Recording
	}{
		{
			name: "default scope",
			want: Recording{
				Scope{}: nil,
			},
		},
		{
			name:       "configured scope",
			loggerName: "test",
			loggerOptions: []log.LoggerOption{
				log.WithInstrumentationVersion("logtest v42"),
				log.WithSchemaURL("https://example.com"),
				log.WithInstrumentationAttributes(attribute.String("foo", "bar")),
			},
			want: Recording{
				Scope{
					Name:       "test",
					Version:    "logtest v42",
					SchemaURL:  "https://example.com",
					Attributes: attribute.NewSet(attribute.String("foo", "bar")),
				}: nil,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rec := NewRecorder(tt.options...)
			rec.Logger(tt.loggerName, tt.loggerOptions...)
			got := rec.Result()
			assert.Equal(t, tt.want, got)
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

func TestRecorderLoggerEmitAndReset(t *testing.T) {
	rec := NewRecorder()
	ts := time.Now()

	l := rec.Logger(t.Name())
	ctx := context.Background()
	r := log.Record{}
	r.SetTimestamp(ts)
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.StringValue("Hello there"))
	r.AddAttributes(log.Int("n", 1))
	r.AddAttributes(log.String("foo", "bar"))
	l.Emit(ctx, r)

	l2 := rec.Logger(t.Name())
	r2 := log.Record{}
	r2.SetBody(log.StringValue("Logger with the same scope"))
	l2.Emit(ctx, r2)

	want := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Context:   ctx,
				Timestamp: ts,
				Severity:  log.SeverityInfo,
				Body:      log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
			{
				Context:    ctx,
				Body:       log.StringValue("Logger with the same scope"),
				Attributes: []log.KeyValue{},
			},
		},
	}
	got := rec.Result()
	assert.Equal(t, want, got)

	rec.Reset()

	want = Recording{
		Scope{Name: t.Name()}: nil,
	}
	got = rec.Result()
	assert.Equal(t, want, got)
}

func TestRecorderConcurrentSafe(*testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := &Recorder{}

	for range goRoutineN {
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
