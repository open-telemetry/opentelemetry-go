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

//go:build go1.18
// +build go1.18

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeterConcurrentSafe(t *testing.T) {
	const name = "TestMeterConcurrentSafe meter"
	mp := NewMeterProvider()

	go func() {
		_ = mp.Meter(name)
	}()

	_ = mp.Meter(name)
}

func TestForceFlushConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	go func() {
		_ = mp.ForceFlush(context.Background())
	}()

	_ = mp.ForceFlush(context.Background())
}

func TestShutdownConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	go func() {
		_ = mp.Shutdown(context.Background())
	}()

	_ = mp.Shutdown(context.Background())
}

func TestMeterDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Meter("") })
}

func TestForceFlushDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.ForceFlush(context.Background()) })
}

func TestShutdownDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Shutdown(context.Background()) })
}
