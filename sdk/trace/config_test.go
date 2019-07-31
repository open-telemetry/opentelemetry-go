// Copyright 2019, OpenTelemetry Authors
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

package trace

import (
	"reflect"
	"testing"
)

func TestApplyConfig(t *testing.T) {
	testCfgs := []Config{
		{},
		{
			MaxAttributesPerSpan: 1,
			MaxEventsPerSpan:     3,
			MaxLinksPerSpan:      4,
		},
		{
			MaxAttributesPerSpan: -1,
			MaxEventsPerSpan:     -3,
			MaxLinksPerSpan:      5,
		}}
	cfg := config.Load().(*Config)
	wantCfgs := []Config{
		{
			DefaultSampler:       cfg.DefaultSampler,
			IDGenerator:          cfg.IDGenerator,
			MaxAttributesPerSpan: DefaultMaxAttributesPerSpan,
			MaxEventsPerSpan:     DefaultMaxEventsPerSpan,
			MaxLinksPerSpan:      DefaultMaxLinksPerSpan,
		},
		{
			DefaultSampler:       cfg.DefaultSampler,
			IDGenerator:          cfg.IDGenerator,
			MaxAttributesPerSpan: 1,
			MaxEventsPerSpan:     3,
			MaxLinksPerSpan:      4,
		},
		{
			DefaultSampler:       cfg.DefaultSampler,
			IDGenerator:          cfg.IDGenerator,
			MaxAttributesPerSpan: 1,
			MaxEventsPerSpan:     3,
			MaxLinksPerSpan:      5,
		}}

	for i, newCfg := range testCfgs {
		ApplyConfig(newCfg)
		gotCfg := config.Load().(*Config)
		wantCfg := wantCfgs[i]

		if got, want := reflect.ValueOf(gotCfg.DefaultSampler).Pointer(), reflect.ValueOf(wantCfg.DefaultSampler).Pointer(); got != want {
			t.Fatalf("testId = %d config.DefaultSampler = %#v; want %#v", i, got, want)
		}
		if got, want := gotCfg.IDGenerator, wantCfg.IDGenerator; got != want {
			t.Fatalf("testId = %d config.IDGenerator = %#v; want %#v", i, got, want)
		}
		if got, want := gotCfg.MaxAttributesPerSpan, wantCfg.MaxAttributesPerSpan; got != want {
			t.Fatalf("testId = %d config.MaxAttributesPerSpan = %#v; want %#v", i, got, want)
		}
		if got, want := gotCfg.MaxLinksPerSpan, wantCfg.MaxLinksPerSpan; got != want {
			t.Fatalf("testId = %d config.MaxLinksPerSpan = %#v; want %#v", i, got, want)
		}
		if got, want := gotCfg.MaxEventsPerSpan, wantCfg.MaxEventsPerSpan; got != want {
			t.Fatalf("testId = %d config.MaxEventsPerSpan = %#v; want %#v", i, got, want)
		}

	}
}
