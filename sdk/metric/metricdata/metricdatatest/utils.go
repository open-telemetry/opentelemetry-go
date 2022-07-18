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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func notEqualStr(prefix string, expected, actual interface{}) string {
	return fmt.Sprintf("%s not equal:\nexpected: %v\nactual: %v", prefix, expected, actual)
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func equalPtrValues[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}

func diffSlices[T any](a, b []T, equal func(T, T) bool) (extraA, extraB []T) {
	visited := make([]bool, len(b))
	for i := 0; i < len(a); i++ {
		found := false
		for j := 0; j < len(b); j++ {
			if visited[j] {
				continue
			}
			if equal(a[i], b[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, a[i])
		}
	}

	for j := 0; j < len(b); j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, b[j])
	}

	return extraA, extraB
}

func compareDiff[T any](extraExpected, extraActual []T) (equal bool, explanation string) {
	if len(extraExpected) == 0 && len(extraActual) == 0 {
		return true, explanation
	}

	formater := func(v T) string {
		return fmt.Sprintf("%#v", v)
	}

	var msg bytes.Buffer
	if len(extraExpected) > 0 {
		_, _ = msg.WriteString("missing expected values:\n")
		for _, v := range extraExpected {
			_, _ = msg.WriteString(formater(v) + "\n")
		}
	}

	if len(extraActual) > 0 {
		_, _ = msg.WriteString("unexpected additional values:\n")
		for _, v := range extraActual {
			_, _ = msg.WriteString(formater(v) + "\n")
		}
	}

	return false, msg.String()
}

func assertCompare(equal bool, explanation []string) func(*testing.T) bool { // nolint: revive  // equal is not a control flag.
	return func(t *testing.T) bool {
		t.Helper()
		if !equal {
			if len(explanation) > 0 {
				t.Error(strings.Join(explanation, "\n"))
			} else {
				t.Fail()
			}
		}
		return equal
	}
}
