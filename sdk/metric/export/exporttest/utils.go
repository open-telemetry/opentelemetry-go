package exporttest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func compareDiff[T any](extraExpected, extraActual []T) (equal bool, explination string) {
	if len(extraExpected) == 0 && len(extraActual) == 0 {
		return true, explination
	}

	formater := func(v T) string {
		return fmt.Sprintf("%#v", v)
	}

	var msg bytes.Buffer
	if len(extraExpected) > 0 {
		msg.WriteString("missing expected values:\n")
		for _, v := range extraExpected {
			msg.WriteString(formater(v) + "\n")
		}
	}

	if len(extraActual) > 0 {
		msg.WriteString("unexpected additional values:\n")
		for _, v := range extraActual {
			msg.WriteString(formater(v) + "\n")
		}
	}

	return false, msg.String()
}

func assertCompare(equal bool, explination []string) func(*testing.T) bool {
	if equal {
		return func(*testing.T) bool { return true }
	}
	return func(t *testing.T) bool {
		return assert.Fail(t, strings.Join(explination, "\n"))
	}
}
