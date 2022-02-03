---
title: Go
description: >-
  <img width="35" src="https://raw.github.com/open-telemetry/opentelemetry.io/main/iconography/32x32/Golang_SDK.svg"></img>
  A language-specific implementation of OpenTelemetry in Go.
aliases: [/golang, /golang/metrics, /golang/tracing]
cascade:
  github_repo: &repo https://github.com/open-telemetry/opentelemetry-go
  github_subdir: website_docs
  path_base_for_github_subdir: content/en/docs/instrumentation/go/
  github_project_repo: *repo
spelling: cSpell:ignore godoc
weight: 16
---

This is the OpenTelemetry for Go documentation. OpenTelemetry is an observability framework -- an API, SDK, and tools that are designed to aid in the generation and collection of application telemetry data such as metrics, logs, and traces. This documentation is designed to help you understand how to get started using OpenTelemetry for Go.

## Status and Releases

The current status of the major functional components for OpenTelemetry Go is as follows:

| Tracing | Metrics | Logging |
| ------- | ------- | ------- |
| Stable  | Alpha   | Not Yet Implemented |

{{% latest_release "go" /%}}

## Further Reading

- [godoc](https://pkg.go.dev/go.opentelemetry.io/otel)
- [Examples](https://github.com/open-telemetry/opentelemetry-go/tree/main/example)
- [Contrib Repository](https://github.com/open-telemetry/opentelemetry-go-contrib)
