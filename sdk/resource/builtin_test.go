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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestBuiltinStringDetector(t *testing.T) {
	E := fmt.Errorf("no K")
	res, err := resource.StringDetector("", attribute.Key("K"), func() (string, error) {
		return "", E
	}).Detect(context.Background())
	require.True(t, errors.Is(err, E))
	require.NotEqual(t, E, err)
	require.Nil(t, res)
}

func TestStringDetectorErrors(t *testing.T) {
	tests := []struct {
		desc        string
		s           resource.Detector
		errContains string
	}{
		{
			desc: "explicit error from func should be returned",
			s: resource.StringDetector("", attribute.Key("K"), func() (string, error) {
				return "", fmt.Errorf("K-IS-MISSING")
			}),
			errContains: "K-IS-MISSING",
		},
		{
			desc: "empty key is an invalid",
			s: resource.StringDetector("", attribute.Key(""), func() (string, error) {
				return "not-empty", nil
			}),
			errContains: "invalid attribute: \"\" -> \"not-empty\"",
		},
	}

	for _, test := range tests {
		res, err := resource.New(
			context.Background(),
			resource.WithAttributes(attribute.String("A", "B")),
			resource.WithDetectors(test.s),
		)
		require.Error(t, err, test.desc)
		require.Contains(t, err.Error(), test.errContains)
		require.NotNil(t, res, "resource contains remaining valid entries")

		m := map[string]string{}
		for _, kv := range res.Attributes() {
			m[string(kv.Key)] = kv.Value.Emit()
		}
		require.EqualValues(t, map[string]string{"A": "B"}, m)
	}

}
