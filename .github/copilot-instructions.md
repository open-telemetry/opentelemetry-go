# Copilot instructions for opentelemetry-go

This repository is the Go implementation of OpenTelemetry. Prefer changes that preserve specification compliance, API stability, and idiomatic Go over clever abstractions or broad refactors.

## Implementation rules

- Read the package you are editing and follow its existing patterns for naming, options, error handling, comments, and tests.
- Keep public APIs backward compatible unless the task explicitly requires a breaking change.
- Follow the OpenTelemetry specification and semantic conventions. Match span, metric, log, attribute, event, resource, and instrumentation scope behavior to the spec and existing package behavior.
- Prefer idiomatic Go and the repository's established patterns over inventing new abstractions.
- For configurable constructors, reuse the project's usual option pattern: unexported `config` types, sealed `Option` interfaces with `apply`, and exported `With...` or `Without...` helpers.
- Be conservative on hot paths. Avoid unnecessary allocations, reflection, interface churn, global state, and high-cardinality telemetry.
- Preserve concurrency safety. If code touches shared state, synchronization, or lifecycle ordering, make the invariants explicit in code and tests.
- Write comments only for intent, invariants, and non-obvious constraints. Do not add comments that restate the code.
- Add or update tests for every behavior change. Add or update benchmarks for performance-sensitive changes.

## Documentation and repository conventions

- Non-internal, non-test packages should have Go doc comments, usually in `doc.go`.
- Non-internal, non-test, non-documentation packages should also have a `README.md` with at least a title and `pkg.go.dev` badge.
- If a change is user-visible, consider whether `CHANGELOG.md`, examples, package docs, or migration notes also need updates.
- Keep license headers intact in Go and shell files.
