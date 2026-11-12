#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# Check each module independently so a workspace cannot hide dependencies.
export GOWORK=off

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK_DIR=$(mktemp -d)
trap 'rm -rf "$WORK_DIR"' EXIT

HTTP_MODULES=(
    exporters/otlp/otlpmetric/otlpmetrichttp
    exporters/otlp/otlplog/otlploghttp
)
FORBIDDEN='(^|[[:space:]])(google[.]golang[.]org/grpc|github[.]com/grpc-ecosystem/grpc-gateway|go[.]opentelemetry[.]io/proto/otlp)(@|/|[[:space:]]|$)'

check_dependencies() {
    local directory=$1
    echo "Checking OTLP/HTTP dependencies in $directory"
    (
        cd "$directory"
        go mod graph > "$WORK_DIR/graph"
        go list -deps -test ./... > "$WORK_DIR/packages"
    )
    if grep -E "$FORBIDDEN" "$WORK_DIR/graph" "$WORK_DIR/packages"; then
        echo "Error: OTLP/HTTP depends on google.golang.org/grpc, github.com/grpc-ecosystem/grpc-gateway/v2, or go.opentelemetry.io/proto/otlp." >&2
        return 1
    fi
}

for module in "${HTTP_MODULES[@]}"; do
    check_dependencies "$ROOT_DIR/$module"
done

# Dependency-module replacements are not inherited by consumers. Replace every
# repository module explicitly so these consumers exercise the local source.
git -C "$ROOT_DIR" ls-files 'go.mod' '*/go.mod' > "$WORK_DIR/modules"
new_consumer() {
    local name=$1
    local fixture=$2
    local exporters=("${HTTP_MODULES[@]}")
    if [[ "$name" == mixed ]]; then
        exporters+=(
            exporters/otlp/otlplog/otlploggrpc
            exporters/otlp/otlpmetric/otlpmetricgrpc
            exporters/otlp/otlptrace/otlptracegrpc
            exporters/otlp/otlptrace/otlptracehttp
        )
    fi
    mkdir "$WORK_DIR/$name"
    cp "$ROOT_DIR/internal/shared/otlp/testdata/httpdependencies/$fixture" "$WORK_DIR/$name/consumer_test.go"
    (
        cd "$WORK_DIR/$name"
        go mod init "example.com/$name"
        while IFS= read -r mod_file; do
            mod_path=$(awk '$1 == "module" {print $2; exit}' "$ROOT_DIR/$mod_file")
            go mod edit "-replace=$mod_path=$ROOT_DIR/$(dirname "$mod_file")"
        done < "$WORK_DIR/modules"
        # Load the exporters' pinned dependencies before resolving fixture
        # imports, rather than selecting the latest external protobuf version.
        for exporter in "${exporters[@]}"; do
            mod_path=$(awk '$1 == "module" {print $2; exit}' "$ROOT_DIR/$exporter/go.mod")
            go mod edit "-require=$mod_path@v0.0.0"
        done
        go mod tidy
    )
}

new_consumer http-only http_only_test.go
check_dependencies "$WORK_DIR/http-only"
(cd "$WORK_DIR/http-only" && go test -count=1 -timeout=2m ./...)

# Keep canonical imports outside the HTTP modules. Importing all six exporters
# also catches descriptor-registration collisions during process initialization.
new_consumer mixed mixed_test.go
(cd "$WORK_DIR/mixed" && go test -race -count=1 -timeout=2m ./...)

echo "OTLP/HTTP dependency and protobuf coexistence checks passed."
