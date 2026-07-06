# Derive missing exception attributes independently

## Summary

- preserve every exception attribute explicitly supplied by the caller;
- derive `exception.message` when it is missing;
- derive `exception.type` when it is missing; and
- keep a caller-provided stacktrace from suppressing message and type derivation.

This changes exception precedence from all-or-nothing suppression to per-key precedence, as required by the Stable Logs SDK specification and exception semantic conventions.

## Testing

- `go test ./...` in `sdk/log`
- `make precommit`
- `benchstat` comparison of `BenchmarkLoggerSetErrAndEmit` (6 runs): no statistically significant time or byte change; allocations remain 1/op

## Related issue

Fixes #5683.

## Integration note

The data-model audit in #8543 uses this same exception fix. The dropped-count fix for #8547 also edits the exception helper; preserve both changes when integrating the branches.
