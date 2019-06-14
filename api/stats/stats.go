package stats

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
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
