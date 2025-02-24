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

	// Emit a record.
	l := rec.Logger(t.Name())
	ctx := context.Background()
	r := log.Record{}
	r.SetSeverity(log.SeverityInfo)
	r.SetTimestamp(ts)
	r.SetBody(log.StringValue("Hello there"))
	r.AddAttributes(log.Int("n", 1))
	r.AddAttributes(log.String("foo", "bar"))
	l.Emit(ctx, r)

	got := rec.Result()
	want := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
	}
	cmpCtx := cmpopts.EquateComparable(context.Background())
	cmpKVs := cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key })
	cmpStmps := cmpopts.IgnoreTypes(time.Time{})
	cmpEpty := cmpopts.EquateEmpty()
	if diff := cmp.Diff(want, got, cmpCtx, cmpKVs, cmpStmps, cmpEpty); diff != "" {
		t.Errorf("Recorded records mismatch (-want +got):\n%s", diff)
	}

	rec.Reset()

	got = rec.Result()
	want = Recording{
		Scope{Name: t.Name()}: nil,
	}
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

func TestRecordingEqual(t *testing.T) {
	a := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.String("foo", "bar"),
					log.Int("n", 1),
				},
			},
		},
		Scope{Name: "Empty"}: nil,
	}

	b := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
		Scope{Name: "Empty"}: []Record{},
	}

	if !Equal(b, a) {
		t.Errorf("Recording mismatch\na:\n%+v\nb:\n%+v", b, a)
	}
}

func TestRecordEqual(t *testing.T) {
	a := Record{
		Severity: log.SeverityInfo,
		Body:     log.StringValue("Hello there"),
		Attributes: []log.KeyValue{
			log.Int("n", 1),
			log.String("foo", "bar"),
		},
	}
	b := Record{
		Severity: log.SeverityInfo,
		Body:     log.StringValue("Hello there"),
		Attributes: []log.KeyValue{
			// Order of attributes is not important.
			log.String("foo", "bar"),
			log.Int("n", 1),
		},
	}
	if !Equal(a, b) {
		t.Errorf("Record mismatch\na:\n%+v\nb:\n%+v", a, b)
	}
}
