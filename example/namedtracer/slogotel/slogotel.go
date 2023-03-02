package slogotel

import (
	"context"

	"golang.org/x/exp/slog"

	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
)

type OtelSlogHandler struct {
	logger otellog.Logger
}

func NewOtelSlogHandler(logger otellog.Logger) *OtelSlogHandler {
	return &OtelSlogHandler{logger: logger}
}

func (o OtelSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (o OtelSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make([]attribute.KeyValue, 0, record.NumAttrs())
	record.Attrs(
		func(attr slog.Attr) {
			attrs = append(attrs, slogToOtelAttr(attr))
		},
	)

	o.logger.Emit(ctx, otellog.WithAttributes(attrs...))
	return nil
}

func slogToOtelAttr(attr slog.Attr) (r attribute.KeyValue) {
	r.Key = attribute.Key(attr.Key)
	switch attr.Value.Kind() {
	case slog.KindString:
		r.Value = attribute.StringValue(attr.Value.String())
	default:
		panic("implement other cases")
	}
	return r
}

func (o OtelSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func (o OtelSlogHandler) WithGroup(name string) slog.Handler {
	//TODO implement me
	panic("implement me")
}

var _ slog.Handler = (*OtelSlogHandler)(nil)
