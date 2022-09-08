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

package otlpmetrichttp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"
)

func TestClient(t *testing.T) {
	factory := func() (otlpmetric.Client, otest.Collector) {
		coll, err := otest.NewHTTPCollector("", nil)
		require.NoError(t, err)

		addr := coll.Addr().String()
		client, err := newClient(WithEndpoint(addr), WithInsecure())
		require.NoError(t, err)
		return client, coll
	}

	t.Run("Integration", otest.RunClientTests(factory))
}
