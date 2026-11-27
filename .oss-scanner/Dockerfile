# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.27.0-bookworm@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452

ENV DEBIAN_FRONTEND=noninteractive \
    CGO_ENABLED=1 \
    GOTOOLCHAIN=local \
    GOWORK=off

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates build-essential clang cmake ninja-build pkg-config git python3 curl \
    && rm -rf /var/lib/apt/lists/*

# The scanner expects the complete checkout at /src.
COPY . /src
WORKDIR /src
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]

# The scanner cannot follow the threat model's external policy link offline.
RUN mkdir -p /opt/oss-scanner \
    && curl --fail --silent --show-error --location \
      https://raw.githubusercontent.com/open-telemetry/.github/main/SECURITY.md \
      --output /opt/oss-scanner/SECURITY.md

# Keep dependencies for every module, including development tools, in the image.
# Cache mounts would not be available in the scanner's offline environment.
RUN find . -name go.mod -print0 | sort -z | while IFS= read -r -d '' module; do \
      go -C "$(dirname "$module")" mod download all; \
    done

# Compile production code and tests, retaining Go's default debug information.
RUN make SHELL=/bin/bash '.SHELLFLAGS=-e -o pipefail -c' build

# As in the scanner template, test failures leave the image available to inspect.
# CI separately requires fresh tests to pass with networking disabled.
RUN timeout 20m make SHELL=/bin/bash '.SHELLFLAGS=-e -o pipefail -c' test-short \
    || echo 'WARNING: short tests failed or timed out; inspect the build log'
