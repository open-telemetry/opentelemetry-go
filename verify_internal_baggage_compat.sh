#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# bridge/opentracing depends on go.opentelemetry.io/otel/internal/baggage, an
# internal package of the root module, to synchronize OpenTracing span
# baggage with context baggage. Because internal/baggage is not part of any
# module's public API, changes to it are not guarded by the usual API
# compatibility tooling (e.g. gorelease), even though bridge/opentracing is a
# separately versioned module that depends on its exact shape.
#
# This builds the latest released bridge/opentracing against the current
# (possibly unreleased) root module source, to catch a break in that
# dependency before it ships rather than after.
#
# See https://github.com/open-telemetry/opentelemetry-go/issues/9199.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BRIDGE_VERSION=$(git -C "$ROOT_DIR" tag --list 'bridge/opentracing/v*' | sed 's#bridge/opentracing/##' | grep -v -- '-' | sort -V | tail -1)

if [[ -z "$BRIDGE_VERSION" ]]; then
    echo "Error: could not determine the latest released bridge/opentracing version."
    exit 1
fi

echo "Checking go.opentelemetry.io/otel/bridge/opentracing@${BRIDGE_VERSION} against the local root module source"

TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

cat > "$TEMP_DIR/go.mod" <<EOF
module internalbaggagecompatcheck

go 1.26.0

require (
	go.opentelemetry.io/otel ${BRIDGE_VERSION}
	go.opentelemetry.io/otel/bridge/opentracing ${BRIDGE_VERSION}
)

replace go.opentelemetry.io/otel => ${ROOT_DIR}

replace go.opentelemetry.io/otel/trace => ${ROOT_DIR}/trace

replace go.opentelemetry.io/otel/metric => ${ROOT_DIR}/metric
EOF

cat > "$TEMP_DIR/main.go" <<'EOF'
package main

import (
	"context"
	"fmt"

	ot "github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/bridge/opentracing"
)

// This exercises the same baggage synchronization path bridge/opentracing
// relies on internal/baggage for: hooking a context so that baggage stored
// via the public baggage API is mirrored onto the active OpenTracing span.
func main() {
	tracer := opentracing.NewBridgeTracer()
	span := tracer.StartSpan("span")
	ctx := tracer.NewHookedContext(ot.ContextWithSpan(context.Background(), span))

	member, err := baggage.NewMember("key", "value")
	if err != nil {
		panic(err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		panic(err)
	}
	baggage.ContextWithBaggage(ctx, bag)

	if got := span.BaggageItem("key"); got != "value" {
		panic(fmt.Sprintf("baggage not synchronized to OpenTracing span: got %q, want %q", got, "value"))
	}
}
EOF

(cd "$TEMP_DIR" && go mod tidy && go build .)

echo "go.opentelemetry.io/otel/bridge/opentracing@${BRIDGE_VERSION} is compatible with the local root module source."
