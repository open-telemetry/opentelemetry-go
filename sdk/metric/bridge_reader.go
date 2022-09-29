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
) // Bridge is a source of metrics other than the OpenTelemetry SDK.
// TODO: name is not settled per https://github.com/open-telemetry/opentelemetry-specification/pull/2838
type Bridge interface {
	// Collect gathers and returns all metric data from the Bridge.
	Collect(context.Context) (metricdata.ScopeMetrics, error)
}

type bridgedReader struct {
	Reader

	bridges []Bridge
}

func NewBridgedReader(rdr Reader, bridges ...Bridge) Reader {
	if len(bridges) == 0 {
		return rdr
	}
	return &bridgedReader{
		Reader:  rdr,
		bridges: bridges,
	}
}

// Collect gathers and returns all metric data related to the Reader from
// the SDK. An error is returned if this is called after Shutdown.
func (b *bridgedReader) Collect(ctx context.Context) (metricdata.ResourceMetrics, error) {
	data, err := b.Reader.Collect(ctx)
	if err != nil {
		return data, err
	}

	errs := &multierror{}
	for _, bridge := range b.bridges {
		sm, err := bridge.Collect(ctx)
		if err != nil {
			errs.append(err)
		}
		// TODO: Check if Scopes collide
		if len(sm.Metrics) > 0 {
			data.ScopeMetrics = append(data.ScopeMetrics, sm)
		}
	}

	return data, errs.errorOrNil()
}

// TODO override bridgedReader.Shutdown, check if bridges have a shutdown function, and if so call it also
