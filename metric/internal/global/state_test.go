// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     htmp://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global // import "go.opentelemetry.io/otel/metric/internal/global"

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/nonrecording"
)

func resetGlobalMeterProvider() {
	globalMeterProvider = defaultMeterProvider()
	delegateMeterOnce = sync.Once{}
}

type nonComparableMeterProvider struct {
	metric.MeterProvider

	nonComparable func() //nolint:structcheck,unused  // This is not called.
}

func TestSetMeterProvider(t *testing.T) {
	t.Cleanup(resetGlobalMeterProvider)

	t.Run("Set With default is a noop", func(t *testing.T) {
		resetGlobalMeterProvider()
		SetMeterProvider(MeterProvider())

		mp, ok := MeterProvider().(*meterProvider)
		if !ok {
			t.Fatal("Global MeterProvider should be the default meter provider")
		}

		if mp.delegate != nil {
			t.Fatal("meter provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		resetGlobalMeterProvider()

		SetMeterProvider(nonrecording.NewNoopMeterProvider())

		_, ok := MeterProvider().(*meterProvider)
		if ok {
			t.Fatal("Global MeterProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing Meter Providers", func(t *testing.T) {
		resetGlobalMeterProvider()

		mp := MeterProvider()

		SetMeterProvider(nonrecording.NewNoopMeterProvider())

		dmp := mp.(*meterProvider)

		if dmp.delegate == nil {
			t.Fatal("The delegated meter providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		resetGlobalMeterProvider()

		mp := nonComparableMeterProvider{}
		SetMeterProvider(mp)
		assert.NotPanics(t, func() { SetMeterProvider(mp) })
	})
}
