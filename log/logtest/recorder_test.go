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

func TestRecordingEqualWithStdLib(t *testing.T) {
	got := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Timestamp: time.Now(),
				Severity:  log.SeverityInfo,
				Body:      log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.String("foo", "bar"),
					log.Int("n", 1),
				},
			},
		},
		Scope{Name: "Empty"}: nil,
	}

	want := Recording{
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
	// Ignore Timestamp.
	for _, recs := range got {
		for i, r := range recs {
			r.Timestamp = time.Time{}
			recs[i] = r
		}
	}
	if !Equal(want, got) {
		t.Errorf("Recording mismatch\na:\n%+v\nb:\n%+v", want, got)
	}
}

func TestRecordingEqualWithTestify(t *testing.T) {
	got := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Timestamp: time.Now(),
				Severity:  log.SeverityInfo,
				Body:      log.StringValue("Hello there"),
				Attributes: []log.KeyValue{

					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
		Scope{Name: "Empty"}: nil,
	}

	want := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Severity: log.SeverityInfo,
				Body:     log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					// Attributes order has to be the same.
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
		// Nil and empty slices are different for testify.
		Scope{Name: "Empty"}: nil,
	}
	// Ignore Timestamp.
	for _, recs := range got {
		for i, r := range recs {
			r.Timestamp = time.Time{}
			recs[i] = r
		}
	}
	assert.Equal(t, want, got)
}

func TestEqualRecordingWithGoCmp(t *testing.T) {
	got := Recording{
		Scope{Name: t.Name()}: []Record{
			{
				Timestamp: time.Now(),
				Severity:  log.SeverityInfo,
				Body:      log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.String("foo", "bar"),
					log.Int("n", 1),
				},
			},
		},
		Scope{Name: "Empty"}: nil,
	}

	want := Recording{
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
	seq := cmpopts.EquateComparable(context.Background())
	ordattr := cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key })
	ignstamp := cmpopts.IgnoreTypes(time.Time{}) // ignore Timestamps
	if diff := cmp.Diff(want, got, seq, ordattr, ignstamp, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("Recorded records mismatch (-want +got):\n%s", diff)
	}
}

func TestEqualRecord(t *testing.T) {
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

func TestRecorderEmitAndReset(t *testing.T) {
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
				Context:   ctx,
				Severity:  log.SeverityInfo,
				Timestamp: ts,
				Body:      log.StringValue("Hello there"),
				Attributes: []log.KeyValue{
					log.Int("n", 1),
					log.String("foo", "bar"),
				},
			},
		},
	}
	assert.Equal(t, want, got)

	rec.Reset()

	got = rec.Result()
	want = Recording{
		// For testify nil and empty slice is imporant.
		Scope{Name: t.Name()}: nil,
	}
	assert.Equal(t, want, got)
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
