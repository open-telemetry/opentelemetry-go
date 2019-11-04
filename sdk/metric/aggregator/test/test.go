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

package test

import (
	"context"
	"math/rand"
	"sort"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/export"
)

var _ export.MetricBatcher = &metricBatcher{}
var _ export.MetricRecord = &metricRecord{}

const Magnitude = 1000

type Profile struct {
	NumberKind core.NumberKind
	Random     func(sign int) core.Number
}

func newProfiles() []Profile {
	rnd := rand.New(rand.NewSource(rand.Int63()))
	return []Profile{
		{
			NumberKind: core.Int64NumberKind,
			Random: func(sign int) core.Number {
				return core.NewInt64Number(int64(sign) * int64(rnd.Intn(Magnitude+1)))
			},
		},
		{
			NumberKind: core.Float64NumberKind,
			Random: func(sign int) core.Number {
				return core.NewFloat64Number(float64(sign) * rnd.Float64() * Magnitude)
			},
		},
	}
}

type metricBatcher struct {
}

type metricRecord struct {
	descriptor *export.Descriptor
}

func NewAggregatorTest(mkind export.MetricKind, nkind core.NumberKind, alternate bool) (export.MetricBatcher, export.MetricRecord) {
	desc := export.NewDescriptor("test.name", mkind, nil, "", "", nkind, alternate)
	return &metricBatcher{}, &metricRecord{descriptor: desc}
}

func (t *metricRecord) Descriptor() *export.Descriptor {
	return t.descriptor
}

func (t *metricRecord) Labels() []core.KeyValue {
	return nil
}

func (m *metricBatcher) AggregatorFor(rec export.MetricRecord) export.MetricAggregator {
	return nil
}

func (m *metricBatcher) Export(context.Context, export.MetricRecord, export.MetricAggregator) {
}

func RunProfiles(t *testing.T, f func(*testing.T, Profile)) {
	for _, profile := range newProfiles() {
		t.Run(profile.NumberKind.String(), func(t *testing.T) {
			f(t, profile)
		})
	}
}

type Numbers struct {
	kind    core.NumberKind
	numbers []core.Number
}

func NewNumbers(kind core.NumberKind) Numbers {
	return Numbers{
		kind: kind,
	}
}

func (n *Numbers) Append(v core.Number) {
	n.numbers = append(n.numbers, v)
}

func (n *Numbers) Sort() {
	sort.Sort(n)
}

func (n *Numbers) Less(i, j int) bool {
	return n.numbers[i].CompareNumber(n.kind, n.numbers[j]) < 0
}

func (n *Numbers) Len() int {
	return len(n.numbers)
}

func (n *Numbers) Swap(i, j int) {
	n.numbers[i], n.numbers[j] = n.numbers[j], n.numbers[i]
}

func (n *Numbers) Sum() core.Number {
	var sum core.Number
	for _, num := range n.numbers {
		sum.AddNumber(n.kind, num)
	}
	return sum
}

func (n *Numbers) Count() int64 {
	return int64(len(n.numbers))
}

func (n *Numbers) Min() core.Number {
	return n.numbers[0]
}

func (n *Numbers) Max() core.Number {
	return n.numbers[len(n.numbers)-1]
}

// Median() is an alias for Quantile(0.5).
func (n *Numbers) Median() core.Number {
	// Note that len(n.numbers) is 1 greater than the max element
	// index, so dividing by two rounds up.  This gives the
	// intended definition for Quantile() in tests, which is to
	// return the smallest element that is at or above the
	// specified quantile.
	return n.numbers[len(n.numbers)/2]
}
