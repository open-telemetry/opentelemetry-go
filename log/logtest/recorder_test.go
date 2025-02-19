// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

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
	rec := NewRecorder()

	// Emit a record.
	l := rec.Logger(t.Name())
	ctx := context.Background()
	r := log.Record{}
	r.SetSeverity(log.SeverityInfo)
	r.SetTimestamp(time.Now())
	r.SetBody(log.StringValue("Hello there"))
	r.AddAttributes(log.Int("n", 1))
	r.AddAttributes(log.String("foo", "bar"))
	l.Emit(ctx, r)

	got := rec.Result()
	// Ignore Timestamp.
	for _, recs := range got {
		for i, r := range recs {
			r.Timestamp = time.Time{}
			recs[i] = r
		}
	}
	want := RecordedRecords{
		Scope{Name: t.Name()}: []Record{
			{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.String("foo", "bar"),
					log.Int("n", 1),
				},
			},
		},
	}
	if !got.Equal(want) {
		t.Errorf("Recorded records mismatch\ngot:\n%+v\nwant:\n%+v", got, want)
	}

	rec.Reset()

	got = rec.Result()
	want = RecordedRecords{
		Scope{Name: t.Name()}: nil,
	}
	if !got.Equal(want) {
		t.Errorf("Records should be cleared\ngot:\n%+v\nwant:\n%+v", got, want)
	}
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
