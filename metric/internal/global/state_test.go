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

	"go.opentelemetry.io/otel/metric/nonrecording"
)

func resetGlobalMeterProvider() {
	globalMeterProvider = defaultMeterProvider()
	delegateMeterOnce = sync.Once{}
}

func TestSetMeterProvider(t *testing.T) {
	t.Cleanup(resetGlobalMeterProvider)

	t.Run("Set With default panics", func(t *testing.T) {
		resetGlobalMeterProvider()

		assert.Panics(t, func() {
			SetMeterProvider(MeterProvider())
		})

	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		resetGlobalMeterProvider()

		SetMeterProvider(nonrecording.NewNoopMeterProvider())

		_, ok := MeterProvider().(*meterProvider)
		if ok {
			t.Error("Global Meter Provider was not changed")
			return
		}
	})

	t.Run("Set() should delegate existing Meter Providers", func(t *testing.T) {
		resetGlobalMeterProvider()

		mp := MeterProvider()

		SetMeterProvider(nonrecording.NewNoopMeterProvider())

		dmp := mp.(*meterProvider)

		if dmp.delegate == nil {
			t.Error("The delegated meter providers should have a delegate")
		}
	})
}
