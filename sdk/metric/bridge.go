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

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Bridge is a source of metrics other than the OpenTelemetry SDK.
// Pull exporters can accept this as configuration, if they want,
// to support additional sources of metrics.
// This is EXPERIMENTAL, and will likely be replaced prior to the stable
// release of the OpenTelemetry metrics sdk.
type Bridge interface {
	// Collect gathers and returns all metric data from the Bridge.
	Collect(context.Context) (metricdata.ScopeMetrics, error)
}
