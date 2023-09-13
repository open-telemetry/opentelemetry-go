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

func Example() {
	// This reader is used as a stand-in for a reader that will actually export
	// data. See exporters in the go.opentelemetry.io/otel/exporters package
	// for more information.
	reader := metric.NewManualReader()

	// See the go.opentelemetry.io/otel/sdk/resource package for more
	// information about how to create and use Resources.
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("0.1.0"),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	otel.SetMeterProvider(meterProvider)
	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// The MeterProvider is configured and registered globally. You can now run
	// your code instrumented with the OpenTelemetry API that uses the global
	// MeterProvider without having to pass this MeterProvider instance. Or,
	// you can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	//
	// See the go.opentelemetry.io/otel/metric package for more information
	// about the metric API.
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
	// SDK using the WithView option. Below is an example of how the view will
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
	// SDK using the WithView option. Below is an example of how the view will
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
	// SDK using the WithView option. Below is an example of how the view will
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
	// SDK using the WithView option. Below is an example of how the view
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
