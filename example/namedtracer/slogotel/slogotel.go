package slogotel

import "golang.org/x/exp/slog"

type OtelSlogHandler struct {
}

func (o OtelSlogHandler) Enabled(level slog.Level) bool {
	//TODO implement me
	panic("implement me")
}

func (o OtelSlogHandler) Handle(r slog.Record) error {
	//TODO implement me
	panic("implement me")
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
