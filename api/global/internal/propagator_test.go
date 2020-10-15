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

package internal_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/oteltest"
)

func TestTextMapPropagatorDelegation(t *testing.T) {
	internal.ResetForTest()
	ctx := context.Background()
	carrier := oteltest.NewTextMapCarrier(nil)

	// The default should be a noop.
	prior := internal.TextMapPropagator()
	prior.Inject(ctx, carrier)
	ctx = prior.Extract(ctx, carrier)
	if !carrier.GotN(t, 0) || !carrier.SetN(t, 0) {
		return
	}

	// Make sure the delegate woks as expected.
	delegate := oteltest.NewTextMapPropagator("test")
	delegate.Inject(ctx, carrier)
	ctx = delegate.Extract(ctx, carrier)
	if !delegate.InjectedN(t, carrier, 1) || !delegate.ExtractedN(t, ctx, 1) {
		return
	}

	// The prior propagator should use the delegate after it is set as the
	// global.
	internal.SetTextMapPropagator(delegate)
	prior.Inject(ctx, carrier)
	ctx = prior.Extract(ctx, carrier)
	delegate.InjectedN(t, carrier, 2)
	delegate.ExtractedN(t, ctx, 2)
}
