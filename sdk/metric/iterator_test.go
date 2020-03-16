// Copyright 2020, OpenTelemetry Authors
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

package metric_test

import (
	"testing"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/test"
	"go.opentelemetry.io/otel/sdk/metric"
)

type sdkIteratorProvider struct{}

var _ test.IteratorProvider = sdkIteratorProvider{}

func (sdkIteratorProvider) Iterators(labels []core.KeyValue) []export.LabelIterator {
	return metric.GetIteratorsForTesting(labels)
}

func (sdkIteratorProvider) EmptyIterators() []export.LabelIterator {
	return metric.GetEmptyIteratorsForTesting()
}

func TestSDKIterators(t *testing.T) {
	test.RunLabelIteratorTests(t, sdkIteratorProvider{})
}
