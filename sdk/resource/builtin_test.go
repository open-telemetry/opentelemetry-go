// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"context"
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
	require.ErrorIs(t, err, E)
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
				return "", fmt.Errorf("k-is-missing")
			}),
			errContains: "k-is-missing",
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
		require.Equal(t, map[string]string{"A": "B"}, m)
	}
}
