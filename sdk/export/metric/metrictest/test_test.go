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

package metrictest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnslice(t *testing.T) {
	in := make([]NoopAggregator, 2)

	a, b := Unslice2(in)

	require.Equal(t, a.(*NoopAggregator), &in[0])
	require.Equal(t, b.(*NoopAggregator), &in[1])
}
