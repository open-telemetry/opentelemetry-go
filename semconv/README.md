# Semconv

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/semconv)](https://pkg.go.dev/go.opentelemetry.io/otel/semconv)

## Generation

Run `make semconv-generate TAG=vX.Y.Z` from the repository root to generate a versioned package.
The generator derives root-package migration information from the previous version.
For additional migration guidance, such as metric instrument or telemetry name changes, add a `MIGRATION_NOTES.md.inc` fragment in the versioned directory.
The generator includes that fragment in `MIGRATION.md` on each run.
