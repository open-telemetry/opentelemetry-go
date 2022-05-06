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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

// Reader is the interface used between the SDK and an
// exporter.  Control flow is bi-directional through the
// Reader, since the SDK initiates ForceFlush and Shutdown
// while the initiates collection.  The Register() method here
// informs the Reader that it can begin reading, signaling the
// start of bi-directional control flow.
//
// Typically, push-based exporters that are periodic will
// implement PeroidicExporter themselves and construct a
// PeriodicReader to satisfy this interface.
//
// Pull-based exporters will typically implement Register
// themselves, since they read on demand.
type Reader interface {
	// Register is called when the SDK is fully
	// configured.  The Producer passed allows the
	// Reader to begin collecting metrics using its
	// Produce() method.
	register(producer)

	// Collect gathers all metrics from the SDK, calling any callbacks necessary.
	// TODO: How does this impact Push exporters?
	// Collect will return an error if called after shutdown.
	Collect(context.Context) (export.Metrics, error)

	// ForceFlush flushes all metric measurements held in an export pipeline.
	//
	// This deadline or cancellation of the passed context are honored. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// telemetry be flushed or all resources have been released in these
	// situations.
	ForceFlush(context.Context) error

	// Shutdown flushes all metric measurements held in an export pipeline and releases any
	// held computational resources.
	//
	// This deadline or cancellation of the passed context are honored. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// telemetry be flushed or all resources have been released in these
	// situations.
	//
	// After Shutdown is called, calls to Collect will perform no operation and instead will return
	// an error indicating the shutdown state.
	Shutdown(context.Context) error
}

//  producer produces metrics for a Reader.
type producer interface {
	// produce returns aggregated metrics from a single collection.
	//
	// This method is safe to call concurrently.
	produce(context.Context) export.Metrics
}
