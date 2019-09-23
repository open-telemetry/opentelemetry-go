// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"

	"go.opentelemetry.io/experimental/streaming/exporter"
	"go.opentelemetry.io/experimental/streaming/sdk/internal"
)

func TestEvents(t *testing.T) {
	obs := internal.NewTestObserver()
	_ = New(obs).WithSpan(context.Background(), "test", func(ctx context.Context) error {
		type test1Type struct{}
		type test2Type struct{}
		span := trace.CurrentSpan(ctx)
		k1v1 := key.New("k1").String("v1")
		k2v2 := key.New("k2").String("v2")
		k3v3 := key.New("k3").String("v3")
		ctx1 := context.WithValue(ctx, test1Type{}, 42)
		span.AddEvent(ctx1, "one two three", k1v1)
		ctx2 := context.WithValue(ctx1, test2Type{}, "foo")
		span.AddEvent(ctx2, "testing", k2v2, k3v3)

		got := obs.Events(exporter.ADD_EVENT)
		for idx := range got {
			if got[idx].Time.IsZero() {
				t.Errorf("Event %d has zero timestamp", idx)
			}
			got[idx].Time = time.Time{}
		}
		if len(got) != 2 {
			t.Errorf("Expected two events, got %d", len(got))
		}
		want := []exporter.Event{
			{
				Type:       exporter.ADD_EVENT,
				String:     "one two three",
				Attributes: []core.KeyValue{k1v1},
			},
			{
				Type:       exporter.ADD_EVENT,
				String:     "testing",
				Attributes: []core.KeyValue{k2v2, k3v3},
			},
		}
		if diffEvents(t, got, want) {
			checkContext(t, got[0].Context, test1Type{}, 42)
			checkContextMissing(t, got[0].Context, test2Type{})
			checkContext(t, got[1].Context, test1Type{}, 42)
			checkContext(t, got[1].Context, test2Type{}, "foo")
		}
		return nil
	})
}

func TestCustomStartEndTime(t *testing.T) {
	startTime := time.Date(2019, time.August, 27, 14, 42, 0, 0, time.UTC)
	endTime := startTime.Add(time.Second * 20)
	obs := internal.NewTestObserver()
	tracer := New(obs)
	_, span := tracer.Start(
		context.Background(),
		"testspan",
		trace.WithStartTime(startTime),
	)
	span.Finish(trace.WithFinishTime(endTime))
	want := []exporter.Event{
		{
			Type:   exporter.START_SPAN,
			Time:   startTime,
			String: "testspan",
		},
		{
			Type: exporter.FINISH_SPAN,
			Time: endTime,
		},
	}
	got := append(obs.Events(exporter.START_SPAN), obs.Events(exporter.FINISH_SPAN)...)
	diffEvents(t, got, want, "Scope")
}

func checkContextMissing(t *testing.T, ctx context.Context, key interface{}) bool {
	gotValue := ctx.Value(key)
	if gotValue != nil {
		keyType := reflect.TypeOf(key)
		t.Errorf("Expected %s to be missing in context", keyType)
		return false
	}
	return true
}

func checkContext(t *testing.T, ctx context.Context, key, wantValue interface{}) bool {
	gotValue := ctx.Value(key)
	if gotValue == nil {
		keyType := reflect.TypeOf(key)
		t.Errorf("Expected %s to exist in context", keyType)
		return false
	}
	if diff := cmp.Diff(gotValue, wantValue); diff != "" {
		keyType := reflect.TypeOf(key)
		t.Errorf("Context value for key %s: -got +want %s", keyType, diff)
		return false
	}
	return true
}

func diffEvents(t *testing.T, got, want []exporter.Event, extraIgnoredFields ...string) bool {
	ignoredPaths := map[string]struct{}{
		"Sequence": struct{}{},
		"Context":  struct{}{},
	}
	for _, field := range extraIgnoredFields {
		ignoredPaths[field] = struct{}{}
	}
	opts := []cmp.Option{
		cmp.FilterPath(func(path cmp.Path) bool {
			_, found := ignoredPaths[path.String()]
			return found
		}, cmp.Ignore()),
	}
	if diff := cmp.Diff(got, want, opts...); diff != "" {
		t.Errorf("Events: -got +want %s", diff)
		return false
	}
	return true
}
