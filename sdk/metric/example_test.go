// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric_test

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// To enable metrics in your application using the SDK,
// you'll need to have an initialized [MeterProvider]
// that will let you create a [go.opentelemetry.io/otel/metric.Meter].
//
// Here's how you might initialize a metrics provider.
func Example() {
	// Create resource.
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("my-service"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		log.Fatalln(err)
	}

	// This reader is used as a stand-in for a reader that will actually export
	// data. See https://pkg.go.dev/go.opentelemetry.io/otel/exporters for
	// exporters that can be used as or with readers.
	reader := metric.NewManualReader()

	// Create a meter provider.
	// You can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)

	// Handle shutdown properly so that nothing leaks.
	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used, which fails to generate data.
	otel.SetMeterProvider(meterProvider)
}

func ExampleView() {
	// The NewView function provides convenient creation of common Views
	// construction. However, it is limited in what it can create.
	//
	// When NewView is not able to provide the functionally needed, a custom
	// View can be constructed directly. Here a custom View is constructed that
	// uses Go's regular expression matching to ensure all data stream names
	// have a suffix of the units it uses.

	re := regexp.MustCompile(`[._](ms|byte)$`)
	var view metric.View = func(i metric.Instrument) (metric.Stream, bool) {
		// In a custom View function, you need to explicitly copy
		// the name, description, and unit.
		s := metric.Stream{Name: i.Name, Description: i.Description, Unit: i.Unit}
		// Any instrument that does not have a unit suffix defined, but has a
		// dimensional unit defined, update the name with a unit suffix.
		if re.MatchString(i.Name) {
			return s, false
		}
		switch i.Unit {
		case "ms":
			s.Name += ".ms"
		case "By":
			s.Name += ".byte"
		default:
			return s, false
		}
		return s, true
	}

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)

	// Below is an example of how the view will
	// function in the SDK for certain instruments.
	stream, _ := view(metric.Instrument{
		Name: "computation.time.ms",
		Unit: "ms",
	})
	fmt.Println("name:", stream.Name)

	stream, _ = view(metric.Instrument{
		Name: "heap.size",
		Unit: "By",
	})
	fmt.Println("name:", stream.Name)
	// Output:
	// name: computation.time.ms
	// name: heap.size.byte
}

func ExampleNewView() {
	// Create a view that renames the "latency" instrument from the v0.34.0
	// version of the "http" instrumentation library as "request.latency".
	view := metric.NewView(metric.Instrument{
		Name: "latency",
		Scope: instrumentation.Scope{
			Name:    "http",
			Version: "0.34.0",
		},
	}, metric.Stream{Name: "request.latency"})

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)

	// Below is an example of how the view will
	// function in the SDK for certain instruments.
	stream, _ := view(metric.Instrument{
		Name:        "latency",
		Description: "request latency",
		Unit:        "ms",
		Kind:        metric.InstrumentKindCounter,
		Scope: instrumentation.Scope{
			Name:      "http",
			Version:   "0.34.0",
			SchemaURL: "https://opentelemetry.io/schemas/1.0.0",
		},
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("description:", stream.Description)
	fmt.Println("unit:", stream.Unit)
	// Output:
	// name: request.latency
	// description: request latency
	// unit: ms
}

func ExampleNewView_wildcard() {
	// Create a view that sets unit to milliseconds for any instrument with a
	// name suffix of ".ms".
	view := metric.NewView(
		metric.Instrument{Name: "*.ms"},
		metric.Stream{Unit: "ms"},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)

	// Below is an example of how the view will
	// function in the SDK for certain instruments.
	stream, _ := view(metric.Instrument{
		Name: "computation.time.ms",
		Unit: "1",
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("unit:", stream.Unit)
	// Output:
	// name: computation.time.ms
	// unit: ms
}

func ExampleNewView_drop() {
	// Create a view that drops the "latency" instrument from the "http"
	// instrumentation library.
	view := metric.NewView(
		metric.Instrument{
			Name:  "latency",
			Scope: instrumentation.Scope{Name: "http"},
		},
		metric.Stream{Aggregation: metric.AggregationDrop{}},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)
}

func ExampleNewView_attributeFilter() {
	// Create a view that removes the "http.request.method" attribute recorded
	// by the "latency" instrument from the "http" instrumentation library.
	view := metric.NewView(
		metric.Instrument{
			Name:  "latency",
			Scope: instrumentation.Scope{Name: "http"},
		},
		metric.Stream{AttributeFilter: attribute.NewDenyKeysFilter("http.request.method")},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)
}

func ExampleNewView_exponentialHistogram() {
	// Create a view that makes the "latency" instrument from the "http"
	// instrumentation library to be reported as an exponential histogram.
	view := metric.NewView(
		metric.Instrument{
			Name:  "latency",
			Scope: instrumentation.Scope{Name: "http"},
		},
		metric.Stream{
			Aggregation: metric.AggregationBase2ExponentialHistogram{
				MaxSize:  160,
				MaxScale: 20,
			},
		},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)
}

func ExampleNewView_exemplarreservoirproviderselector() {
	// Create a view that makes all metrics use a different exemplar reservoir.
	view := metric.NewView(
		metric.Instrument{Name: "*"},
		metric.Stream{
			ExemplarReservoirProviderSelector: func(agg metric.Aggregation) exemplar.ReservoirProvider {
				// This example uses a fixed-size reservoir with a size of 10
				// for explicit bucket histograms instead of the default
				// bucket-aligned reservoir.
				if _, ok := agg.(metric.AggregationExplicitBucketHistogram); ok {
					return exemplar.FixedSizeReservoirProvider(10)
				}
				// Fall back to the default reservoir otherwise.
				return metric.DefaultExemplarReservoirProviderSelector(agg)
			},
		},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)
}

func ExampleWithExemplarFilter_disabled() {
	// Use exemplar.AlwaysOffFilter to disable exemplar collection.
	_ = metric.NewMeterProvider(
		metric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
}

func ExampleWithExemplarFilter_custom() {
	// Create a custom filter function that only offers measurements if the
	// context has an error.
	customFilter := func(ctx context.Context) bool {
		return ctx.Err() != nil
	}
	_ = metric.NewMeterProvider(
		metric.WithExemplarFilter(customFilter),
	)
}
