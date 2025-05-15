// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "dns" namespace.
package dnsconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes the error the DNS lookup failed with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// LookupDuration is an instrument used to record metric values conforming to the
// "dns.lookup.duration" semantic conventions. It represents the measures the
// time taken to perform a DNS lookup.
type LookupDuration struct {
	metric.Float64Histogram
}

// NewLookupDuration returns a new LookupDuration instrument.
func NewLookupDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (LookupDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return LookupDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"dns.lookup.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Measures the time taken to perform a DNS lookup."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return LookupDuration{noop.Float64Histogram{}}, err
	}
	return LookupDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LookupDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (LookupDuration) Name() string {
	return "dns.lookup.duration"
}

// Unit returns the semantic convention unit of the instrument
func (LookupDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (LookupDuration) Description() string {
	return "Measures the time taken to perform a DNS lookup."
}

// Record records val to the current distribution.
//
// The questionName is the the name being queried.
//
// All additional attrs passed are included in the recorded value.
func (m LookupDuration) Record(
	ctx context.Context,
	val float64,
	questionName string,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("dns.question.name", questionName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes the error the DNS lookup failed with.
func (LookupDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}