package view

import (
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// StandardAggregation is the specified default aggregation Kind for
// each instrument Kind.
func StandardAggregationKind(ik sdkinstrument.Kind) aggregation.Kind {
	switch ik {
	case sdkinstrument.HistogramKind:
		return aggregation.HistogramKind
	case sdkinstrument.GaugeObserverKind:
		return aggregation.GaugeKind
	case sdkinstrument.UpDownCounterKind, sdkinstrument.UpDownCounterObserverKind:
		return aggregation.NonMonotonicSumKind
	default:
		return aggregation.MonotonicSumKind
	}
}

// StandardTemporality returns the specified default Cumulative
// temporality for all instrument kinds.
func StandardTemporality(ik sdkinstrument.Kind) aggregation.Temporality {
	return aggregation.CumulativeTemporality
}

// DeltaPreferredTemporality returns the specified Delta temporality
// for all instrument kinds except UpDownCounter, which remain Cumulative.
func DeltaPreferredTemporality(ik sdkinstrument.Kind) aggregation.Temporality {
	switch ik {
	case sdkinstrument.UpDownCounterKind, sdkinstrument.UpDownCounterObserverKind:
		return aggregation.CumulativeTemporality
	default:
		return aggregation.DeltaTemporality
	}
}

// StandardConfig returns two default-initialized aggregator.Configs.
func StandardConfig(ik sdkinstrument.Kind) (ints, floats aggregator.Config) {
	return aggregator.Config{}, aggregator.Config{}
}
