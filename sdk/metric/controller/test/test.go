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

package test

import (
	"time"

	"github.com/benbjohnson/clock"

	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
)

type MockClock struct {
	mock *clock.Mock
}

type MockTicker struct {
	ticker *clock.Ticker
}

var _ controllerTime.Clock = MockClock{}
var _ controllerTime.Ticker = MockTicker{}

func NewMockClock() MockClock {
	return MockClock{clock.NewMock()}
}

func (c MockClock) Now() time.Time {
	return c.mock.Now()
}

func (c MockClock) Ticker(period time.Duration) controllerTime.Ticker {
	return MockTicker{c.mock.Ticker(period)}
}

func (c MockClock) Add(d time.Duration) {
	c.mock.Add(d)
}

func (t MockTicker) Stop() {
	t.ticker.Stop()
}

func (t MockTicker) C() <-chan time.Time {
	return t.ticker.C
}
