#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# Some internal packages are deliberately shared across module boundaries
# (see "Internal packages" in CONTRIBUTING.md), currently
# go.opentelemetry.io/otel/internal/baggage and
# go.opentelemetry.io/otel/internal/global. Because an internal package is
# not part of any module's public API, changes to it are not guarded by the
# usual API compatibility tooling (e.g. gorelease), even though other,
# separately versioned modules depend on its exact shape.
#
# For each module listed in TARGETS below, this builds the latest released
# version of that module against the local (possibly unreleased) source of
# every module in this repository, using a Go workspace, to catch a break in
# a cross-module internal dependency before it ships rather than after.
#
# This already happened once: the root module's v1.47.0 release broke
# compilation of bridge/opentracing v1.46.0 and earlier. See
# https://github.com/open-telemetry/opentelemetry-go/issues/9199.
#
# To add a module to this check, add its import path to TARGETS. Its
# released go.mod does not need any special handling: this script strips any
# `replace` directives it carries before adding it to the workspace, so the
# workspace's resolution of local modules always takes precedence.

set -euo pipefail

TARGETS=(
    "go.opentelemetry.io/otel/bridge/opentracing"
    "go.opentelemetry.io/otel/exporters/prometheus"
)

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

WORK_DIR=$(mktemp -d)
trap 'rm -rf "$WORK_DIR"' EXIT

cd "$WORK_DIR"
go work init

# Add every local module as a workspace member, so targets resolve against
# local source instead of the versions recorded in their go.mod. Skip a
# module that is itself a target: that one is added later from its released
# source instead, which is the whole point of the check.
while IFS= read -r mod_dir; do
    mod_path=$(go mod edit -json "$ROOT_DIR/$mod_dir/go.mod" | sed -n 's/.*"Path": "\([^"]*\)".*/\1/p' | head -1)
    is_target=false
    for target in "${TARGETS[@]}"; do
        [[ "$mod_path" == "$target" ]] && is_target=true && break
    done
    "$is_target" || go work use "$ROOT_DIR/$mod_dir"
done < <(find "$ROOT_DIR" -name go.mod -not -path "$ROOT_DIR/internal/tools/*" -exec dirname {} \; | sed "s#^$ROOT_DIR/\\?#./#")

status=0
for target in "${TARGETS[@]}"; do
    module_dir="${target#go.opentelemetry.io/otel/}"
    version=$(git -C "$ROOT_DIR" tag --list "${module_dir}/v*" | sed "s#${module_dir}/##" | grep -v -- '-' | sort -V | tail -1)

    if [[ -z "$version" ]]; then
        echo "Error: could not determine the latest released version of ${target}."
        status=1
        continue
    fi

    echo "Checking ${target}@${version} against local module source"

    # Fetch outside the workspace: workspace mode forces -mod=readonly, and
    # the target is not required by any workspace member.
    (cd "$WORK_DIR" && GOWORK=off go mod download "${target}@${version}")
    cached_dir=$(cd "$WORK_DIR" && GOWORK=off go list -m -f '{{.Dir}}' "${target}@${version}")

    # Copy out of the (read-only) module cache and strip any replace
    # directives the release carries, so the workspace's resolution of
    # local modules is the only one in effect.
    local_copy="$WORK_DIR/targets/${module_dir//\//_}"
    mkdir -p "$(dirname "$local_copy")"
    cp -r "$cached_dir" "$local_copy"
    chmod -R u+w "$local_copy"
    while IFS= read -r replaced_mod; do
        (cd "$local_copy" && go mod edit "-dropreplace=${replaced_mod}")
    done < <(grep '^replace ' "$local_copy/go.mod" | awk '{print $2}')

    go work use "$local_copy"

    if ! go build "${target}/..."; then
        echo "Error: ${target}@${version} is incompatible with local module source."
        status=1
    fi
done

if [[ "$status" -eq 0 ]]; then
    echo "All targets are compatible with local module source."
fi

exit "$status"
