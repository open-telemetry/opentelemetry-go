package sdk

import (
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type sdk struct {
	exporter  *exporter.Exporter
	resources exporter.EventID
}

type SDK interface {
	trace.Tracer
	metric.Meter
}

var _ SDK = &sdk{}
