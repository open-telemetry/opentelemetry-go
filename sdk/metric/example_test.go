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

package metric_test

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// To enable metrics in your application using the SDK,
// you'll need to have an initialized [MeterProvider]
// that will let you create a [go.opentelemetry.io/otel/metric.Meter].
//
// Here's how you might initialize a metrics provider.
func Example() {
	// See [go.opentelemetry.io/otel/sdk/resource] for more
	// information about how to create and use resources.
	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("my-service"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		log.Fatalln(err)
	}

	// This reader is used as a stand-in for a reader that will actually export
	// data. See [go.opentelemetry.io/otel/exporters] for exporters
	// that can be used as or with readers.
	reader := metric.NewManualReader()

	// Create a meter provider.
	// You can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)

	// Register as global meter provider so that it can
	// that can used via [go.opentelemetry.io/otel.Meter]
	// and accessed using [go.opentelemetry.io/otel.GetMeterProvider].
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used and which fails to generate data.
	otel.SetMeterProvider(meterProvider)

	// Handle shutdown properly so that nothing leaks.
	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()
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

func ExampleNewView_drop() {
	// Create a view that sets the drop aggregator for all instrumentation from
	// the "db" library, effectively turning-off all instrumentation from that
	// library.
	view := metric.NewView(
		metric.Instrument{Scope: instrumentation.Scope{Name: "db"}},
		metric.Stream{Aggregation: metric.AggregationDrop{}},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option.
	_ = metric.NewMeterProvider(
		metric.WithView(view),
	)

	// Below is an example of how the view will
	// function in the SDK for certain instruments.
	stream, _ := view(metric.Instrument{
		Name:  "queries",
		Kind:  metric.InstrumentKindCounter,
		Scope: instrumentation.Scope{Name: "db", Version: "v0.4.0"},
	})
	fmt.Println("name:", stream.Name)
	fmt.Printf("aggregation: %#v", stream.Aggregation)
	// Output:
	// name: queries
	// aggregation: metric.AggregationDrop{}
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
