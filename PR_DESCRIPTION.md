# Clarify per-emission `Logger.Enabled` guidance

## Summary

- document that callers should evaluate `Logger.Enabled` for every record they intend to emit;
- explain that the enabled state may change over time; and
- bring the public Logs API documentation into alignment with the Stable specification.

## Testing

- `make precommit`

## Related issue

Fixes #5679.
