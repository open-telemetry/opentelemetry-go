// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/trace/x"

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

// InsertOrUpdateTraceStateThKeyValue inserts or updates the threshold (th) key-value
// in the OTel tracestate "ot" field. It is exported for use by samplers that need
// to set the threshold (e.g., AlwaysOn sampler setting th:0).
func InsertOrUpdateTraceStateThKeyValue(existingOtts, thkv string) string {
	if existingOtts == "" {
		return thkv
	}

	start := -1
	var end int
	if strings.HasPrefix(existingOtts, "th:") {
		start = 0
	} else if idx := strings.Index(existingOtts, ";th:"); idx != -1 {
		start = idx + 1
	}
	if start == -1 {
		return thkv + ";" + existingOtts
	}

	for end = start; end < len(existingOtts); end++ {
		if existingOtts[end] == ';' {
			end++
			break
		}
	}

	if end == len(existingOtts) {
		return strings.TrimSuffix(thkv+";"+existingOtts[0:start], ";")
	}
	return thkv + ";" + existingOtts[0:start] + existingOtts[end:]
}

// tracestateRandomness determines whether there is a randomness "rv" sub-key in
// otts (the top-level OTel tracestate field). If present, "rv" is a 56-bit
// unsigned integer, encoded in 14 hex digits.
func tracestateRandomness(otts string) (randomness uint64, hasRandomness bool) {
	var start int
	if strings.HasPrefix(otts, "rv:") {
		start = 3
	} else if idx := strings.Index(otts, ";rv:"); idx != -1 {
		start = idx + 4
	} else {
		return 0, false
	}

	if len(otts) < start+14 || (len(otts) > start+14 && otts[start+14] != ';') {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %s", otts))
		return 0, false
	}

	rv, err := strconv.ParseUint(otts[start:start+14], 16, 56)
	if err != nil {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %s", otts))
		return 0, false
	}
	randomness = rv
	hasRandomness = true
	return randomness, hasRandomness
}

func eraseTraceStateThKeyValue(otts string) string {
	var start int
	if strings.HasPrefix(otts, "th:") {
		start = 0
	} else if idx := strings.Index(otts, ";th:"); idx != -1 {
		start = idx + 1
	} else {
		return otts
	}
	if start > 0 && otts[start-1] == ';' {
		start--
	}
	var end int
	for end = start + 1; end < len(otts); end++ {
		if otts[end] == ';' {
			if start == 0 {
				end++
			}
			break
		}
	}
	if end == len(otts) {
		return otts[0:start]
	}
	return otts[0:start] + otts[end:]
}
