// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package propagation // import "go.opentelemetry.io/otel/propagation"

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

const (
	baggageHeader = "baggage"

	// W3C Baggage specification limits.
	// https://www.w3.org/TR/baggage/#limits
	maxMembers               = 64
	maxBytesPerBaggageString = 8192
)

// Baggage is a propagator that supports the W3C Baggage format.
//
// This propagates user-defined baggage associated with a trace. The complete
// specification is defined at https://www.w3.org/TR/baggage/.
type Baggage struct{}

var _ TextMapPropagator = Baggage{}

// Inject sets baggage key-values from ctx into the carrier.
func (Baggage) Inject(ctx context.Context, carrier TextMapCarrier) {
	bStr := baggage.FromContext(ctx).String()
	if bStr != "" {
		carrier.Set(baggageHeader, bStr)
	}
}

// Extract returns a copy of parent with the baggage from the carrier added.
// If carrier implements [ValuesGetter] (e.g. [HeaderCarrier]), Values is invoked
// for multiple values extraction. Otherwise, Get is called.
func (Baggage) Extract(parent context.Context, carrier TextMapCarrier) context.Context {
	if multiCarrier, ok := carrier.(ValuesGetter); ok {
		return extractMultiBaggage(parent, multiCarrier)
	}
	return extractSingleBaggage(parent, carrier)
}

// Fields returns the keys who's values are set with Inject.
func (Baggage) Fields() []string {
	return []string{baggageHeader}
}

func extractSingleBaggage(parent context.Context, carrier TextMapCarrier) context.Context {
	bStr := carrier.Get(baggageHeader)
	if bStr == "" {
		return parent
	}

	bag, err := baggage.Parse(bStr)
	if err != nil {
		return parent
	}
	return baggage.ContextWithBaggage(parent, bag)
}

func extractMultiBaggage(parent context.Context, carrier ValuesGetter) context.Context {
	bVals := carrier.Values(baggageHeader)
	if len(bVals) == 0 {
		return parent
	}

	// Check combined size of all baggage headers.
	// W3C spec: "If there are multiple baggage headers, all limits apply to
	// the combination of all baggage headers and not each header individually."
	var totalBytes int
	for _, bStr := range bVals {
		totalBytes += len(bStr)
	}
	if totalBytes > maxBytesPerBaggageString {
		return parent
	}

	var members []baggage.Member
	for _, bStr := range bVals {
		currBag, err := baggage.Parse(bStr)
		if err != nil {
			continue
		}

		currMembers := currBag.Members()
		if len(members)+len(currMembers) > maxMembers {
			// W3C Baggage limit exceeded, drop all baggage.
			// If a platform cannot propagate all baggage, it MUST NOT propagate any partial list-members
			return parent
		}
		members = append(members, currMembers...)
	}

	b, err := baggage.New(members...)
	if err != nil || b.Len() == 0 {
		return parent
	}
	return baggage.ContextWithBaggage(parent, b)
}
