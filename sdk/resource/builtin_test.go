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

package resource_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestBuiltinStringDetector(t *testing.T) {
	E := fmt.Errorf("no K")
	res, err := resource.StringDetector(label.Key("K"), func() (string, error) {
		return "", E
	}).Detect(context.Background())
	require.True(t, errors.Is(err, E))
	require.NotEqual(t, E, err)
	require.Nil(t, res)
}

func TestBuiltinStringConfig(t *testing.T) {
	res, err := resource.New(
		context.Background(),
		resource.WithoutBuiltin(),
		resource.WithAttributes(label.String("A", "B")),
		resource.WithDetectors(resource.StringDetector(label.Key("K"), func() (string, error) {
			return "", fmt.Errorf("K-IS-MISSING")
		})),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "K-IS-MISSING")
	require.NotNil(t, res)

	m := map[string]string{}
	for _, kv := range res.Attributes() {
		m[string(kv.Key)] = kv.Value.Emit()
	}
	require.EqualValues(t, map[string]string{
		"A": "B",
	}, m)
}
