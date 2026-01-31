// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package internal

import (
	"os"
	"strings"
)

func GetStabilityMode() (emitOld, emitNew bool) {
	optIn := os.Getenv("OTEL_SEMCONV_STABILITY_OPT_IN")

	// Default: Only old conventions
	if optIn == "" {
		return true, false
	}

	modes := strings.Split(optIn, ",")
	hasRPC := false
	hasRPCDup := false
	for _, m := range modes {
		switch strings.TrimSpace(m) {
		case "rpc/dup":
			hasRPCDup = true
		case "rpc":
			hasRPC = true
		}
	}
	switch {
	case hasRPCDup:
		return true, true // Emit BOTH (transition phase)
	case hasRPC:
		return false, true // Emit ONLY new
	default:
		return true, false // Only old conventions
	}
}
