// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opentelemetry.io/otel/log"
)

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
		},
	}
	cmpCtx := cmpopts.EquateComparable(context.Background())
	cmpKVs := cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key })
	cmpEpty := cmpopts.EquateEmpty()
	got := rec.Result()
	if diff := cmp.Diff(want, got, cmpCtx, cmpKVs, cmpEpty); diff != "" {
		t.Errorf("Recorded records mismatch (-want +got):\n%s", diff)
	}

	rec.Reset()

	want = Recording{
		Scope{Name: t.Name()}: nil,
	}
	got = rec.Result()
	if diff := cmp.Diff(want, got, cmpEpty); diff != "" {
		t.Errorf("Recorded records mismatch (-want +got):\n%s", diff)
	}
}

func TestRecorderLoggerEnabled(t *testing.T) {
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
			got := NewRecorder(tt.options...).Logger("test").Enabled(tt.ctx, tt.enabledParams)
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
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
