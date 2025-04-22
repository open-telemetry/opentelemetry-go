// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/dns"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	inst metric.Float64Histogram
}

// NewLookupDuration returns a new LookupDuration instrument.
func NewLookupDuration(m metric.Meter) (LookupDuration, error) {
	i, err := m.Float64Histogram(
	    "dns.lookup.duration",
	    metric.WithDescription("Measures the time taken to perform a DNS lookup."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return LookupDuration{}, err
	}
	return LookupDuration{i}, nil
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
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("dns.question.name", questionName),
			)...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes the error the DNS lookup failed with.
func (LookupDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}