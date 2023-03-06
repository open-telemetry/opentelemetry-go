# Shared code

This directory contains internal code
that is shared between modules.

It also contains the application used for copying
the shared code to the different Go modules.

[`gen.go`](gen.go) contains `go:generate` directives
which copies the shared code into multiple packages.

**Warning**.
This module should not be published.

**Note**.
We might consier moving [`gocpy`](gocpy/gocpy.go)
to [`opentelemetry-go-build-tools`](https://github.com/open-telemetry/opentelemetry-go-build-tools).
