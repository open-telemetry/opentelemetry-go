// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: NOTE this is a temporary space, it may be moved following the
// discussion of #2813, or #2841

package export // import "go.opentelemetry.io/otel/sdk/metric/export"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Metrics is the result of a single collection.
	//
	// This struct supports re-use of the nested memory structure
	// in its Scopes slice such that repeated calls Produce will
	// not reallocate the same quantity of memory again and again.
	//
	// To re-use the memory from a previous Metrics value, pass a
	// pointer to the former result to Produce(). This is safe for
	// push-based exporters that perform sequential collection.
	Metrics struct {
		// Resource is the MeterProvider's configured Resource.
		Resource *resource.Resource

		// Scopes is a slice of metric data, one per Meter.
		Scopes []Scope
	}

	// Scope is the result of a single collection for a single Meter.
	//
	// See the comments on Metrics about re-use of slices in this struct.
	Scope struct {
		// Library describes the instrumentation scope.
		Library instrumentation.Library

		// Instruments is a slice of metric data, one per Instrument
		// in the scope.
		Instruments []Instrument
	}

	// Instrument is the result of a single collection for a single Instrument.
	//
	// See the comments on Metrics about re-use of slices in this struct.
	Instrument struct {
		// Descriptor describes an instrument created through a View,
		// including name, unit, description, instrument and number kinds.
		// TODO Define sdkinsturment package #2813
		// Descriptor sdkinstrument.Descriptor

		// Points is a slice of metric data, one per attribute.Set value.
		Points []Point
	}

	// Point is a timeseries data point resulting from a single collection.
	Point struct {
		// Attributes are the coordinates of this series.
		Attributes attribute.Set

		// Aggregation determines the kind of data point
		// recorded in this series.
		// TODO: Define aggregation types #2827, #2828
		// Aggregation aggregation.Aggregation

		// Start indicates the start of the collection
		// interval reflected in this series, which is set
		// according to the configured temporality.
		Start time.Time

		// End indicates the moment at which the collection
		// was performed.
		End time.Time
	}
)
