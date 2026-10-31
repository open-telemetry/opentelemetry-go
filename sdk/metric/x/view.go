// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"slices"

	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Stream describes the stream of data an instrument produces, extending
// [sdkmetric.Stream] with the experimental Enabled configuration.
type Stream struct {
	// Name is the human-readable identifier of the stream.
	Name string
	// Description describes the purpose of the data.
	Description string
	// Unit is the unit of measurement recorded.
	Unit string
	// Aggregation the stream uses for an instrument.
	Aggregation sdkmetric.Aggregation
	// AttributeFilter is an attribute Filter applied to the attributes
	// recorded for an instrument's measurement. If the filter returns false
	// the attribute will not be recorded, otherwise, if it returns true, it
	// will record the attribute.
	//
	// If nil, the instrument's advisory attribute keys are used. If the
	// instrument has no advisory attribute keys, all attributes are recorded.
	//
	// Note that attributes filtered out by a View may still appear on Exemplars,
	// because Exemplars are recorded with the dropped measurement attributes
	// when View attribute filtering is applied.
	//
	// Use NewAllowKeysFilter from "go.opentelemetry.io/otel/attribute" to
	// provide an allow-list of attribute keys here.
	AttributeFilter attribute.Filter
	// ExemplarReservoirProviderSelector selects the
	// [go.opentelemetry.io/otel/sdk/metric/exemplar.ReservoirProvider] based
	// on the [sdkmetric.Aggregation].
	//
	// If unspecified, [sdkmetric.DefaultExemplarReservoirProviderSelector] is used.
	ExemplarReservoirProviderSelector sdkmetric.ExemplarReservoirProviderSelector
	// Enabled controls whether the stream is enabled.
	//
	// If nil, the instrument's default enabled state is preserved (enabled for
	// standard instruments, or dropped for opt-in instruments created with
	// [go.opentelemetry.io/otel/metric/x.WithOptIn]). For opt-in instruments,
	// a matching view must set Enabled to true. A view that leaves Enabled nil
	// still results in the stream being dropped.
	// If true, the stream is enabled, including for opt-in instruments.
	// If false, the stream is disabled and uses [sdkmetric.AggregationDrop].
	Enabled *bool
}

// View is an override to the default behavior of the SDK, extending
// [sdkmetric.View] with support for [Stream.Enabled]. It defines how data
// should be collected for certain instruments. It returns true and the exact
// [Stream] to use for matching [sdkmetric.Instrument]s; otherwise if the
// [sdkmetric.Instrument] does not match, false is returned.
type View func(sdkmetric.Instrument) (Stream, bool)

// NewView returns a [View] that applies the mask [Stream] to any instrument
// matching criteria. It behaves identically to [sdkmetric.NewView] with
// additional support for [Stream.Enabled].
func NewView(criteria sdkmetric.Instrument, mask Stream) View {
	var enabled *bool
	if mask.Enabled != nil {
		e := *mask.Enabled
		enabled = &e
	}
	baseView := sdkmetric.NewView(criteria, sdkmetric.Stream{
		Name:                              mask.Name,
		Description:                       mask.Description,
		Unit:                              mask.Unit,
		Aggregation:                       mask.Aggregation,
		AttributeFilter:                   mask.AttributeFilter,
		ExemplarReservoirProviderSelector: mask.ExemplarReservoirProviderSelector,
	})
	return func(inst sdkmetric.Instrument) (Stream, bool) {
		s, match := baseView(inst)
		if !match {
			return Stream{}, false
		}
		var e *bool
		if enabled != nil {
			v := *enabled
			e = &v
		}
		return Stream{
			Name:                              s.Name,
			Description:                       s.Description,
			Unit:                              s.Unit,
			Aggregation:                       s.Aggregation,
			AttributeFilter:                   s.AttributeFilter,
			ExemplarReservoirProviderSelector: s.ExemplarReservoirProviderSelector,
			Enabled:                           e,
		}, true
	}
}

type viewOption struct {
	sdkmetric.Option
	views []View
}

// Experimental prevents the SDK from applying the embedded stable option.
func (viewOption) Experimental() {}

func (o viewOption) ExperimentalViews() []func(sdkmetric.Instrument) (sdkmetric.Stream, *bool, bool) {
	out := make([]func(sdkmetric.Instrument) (sdkmetric.Stream, *bool, bool), 0, len(o.views))
	for _, v := range o.views {
		if v == nil {
			continue
		}
		out = append(out, func(inst sdkmetric.Instrument) (sdkmetric.Stream, *bool, bool) {
			s, match := v(inst)
			if !match {
				return sdkmetric.Stream{}, nil, false
			}
			var enabled *bool
			if s.Enabled != nil {
				e := *s.Enabled
				enabled = &e
			}
			return sdkmetric.Stream{
				Name:                              s.Name,
				Description:                       s.Description,
				Unit:                              s.Unit,
				Aggregation:                       s.Aggregation,
				AttributeFilter:                   s.AttributeFilter,
				ExemplarReservoirProviderSelector: s.ExemplarReservoirProviderSelector,
			}, enabled, true
		})
	}
	return out
}

// WithView associates views with a [sdkmetric.MeterProvider].
//
// Views are appended to existing views in the [sdkmetric.MeterProvider] and
// evaluated in the order options are passed to [sdkmetric.NewMeterProvider].
//
// Any nil view values are ignored.
func WithView(views ...View) sdkmetric.Option {
	return viewOption{views: slices.Clone(views)}
}
