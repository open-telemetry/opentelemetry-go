package stats

import (
	"context"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/scope"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

type (
	Interface interface {
		Record(ctx context.Context, m ...core.Measurement)
		RecordSingle(ctx context.Context, m core.Measurement)
	}

	Recorder struct {
		core.ScopeID
	}
)

func With(scope scope.Scope) Recorder {
	return Recorder{scope.ScopeID()}
}

func Record(ctx context.Context, m ...core.Measurement) {
	With(scope.Active(ctx)).Record(ctx, m...)
}

func RecordSingle(ctx context.Context, m core.Measurement) {
	With(scope.Active(ctx)).RecordSingle(ctx, m)
}

func (r Recorder) Record(ctx context.Context, m ...core.Measurement) {
	observer.Record(observer.Event{
		Type:    observer.RECORD_STATS,
		Scope:   r.ScopeID,
		Context: ctx,
		Stats:   m,
	})
}

func (r Recorder) RecordSingle(ctx context.Context, m core.Measurement) {
	observer.Record(observer.Event{
		Type:    observer.RECORD_STATS,
		Scope:   r.ScopeID,
		Context: ctx,
		Stat:    m,
	})
}
