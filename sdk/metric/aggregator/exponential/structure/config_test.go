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

package structure // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/structure"

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValid(t *testing.T) {
	require.True(t, Config{}.Valid())
	require.True(t, NewConfig().Valid())
	require.True(t, NewConfig(WithMaxSize(MinSize)).Valid())
	require.True(t, NewConfig(WithMaxSize(MaximumMaxSize)).Valid())
	require.True(t, NewConfig(WithMaxSize((MinSize+MaximumMaxSize)/2)).Valid())

	require.False(t, NewConfig(WithMaxSize(-1)).Valid())
	require.False(t, NewConfig(WithMaxSize(1<<20)).Valid())
	require.False(t, NewConfig(WithMaxSize(1)).Valid())
}
