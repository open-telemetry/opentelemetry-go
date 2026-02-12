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

	bag, _ := baggage.Parse(bStr)
	if bag.Len() == 0 {
		return parent
	}
	return baggage.ContextWithBaggage(parent, bag)
}

func extractMultiBaggage(parent context.Context, carrier ValuesGetter) context.Context {
	bVals := carrier.Values(baggageHeader)
	if len(bVals) == 0 {
		return parent
	}

	// W3C Baggage spec limits: https://www.w3.org/TR/baggage/#limits
	// "If either of the above conditions is not met, a platform MAY drop
	// list-members until both conditions are met."
	// We keep the first N complete members that fit within the limits.
	var members []baggage.Member
	var totalBytes int
	limitReached := false
	for _, bStr := range bVals {
		if limitReached {
			break
		}

		currBag, _ := baggage.Parse(bStr)
		if currBag.Len() == 0 {
			// Header produced no members (invalid format or exceeded limits).
			// Skip this header and continue processing subsequent headers.
			continue
		}

		for _, m := range currBag.Members() {
			// Check member count limit.
			if len(members) >= maxMembers {
				limitReached = true
				break
			}

			// Check byte size limit.
			// Account for comma separator between members.
			memberBytes := len(m.String())
			if len(members) > 0 {
				memberBytes++ // comma separator
			}
			if totalBytes+memberBytes > maxBytesPerBaggageString {
				limitReached = true
				break
			}

			members = append(members, m)
			totalBytes += memberBytes
		}
	}

	b, err := baggage.New(members...)
	if err != nil || b.Len() == 0 {
		return parent
	}
	return baggage.ContextWithBaggage(parent, b)
}
