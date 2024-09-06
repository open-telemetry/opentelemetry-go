# Experimental Features

The Trace SDK contains features that have not yet stabilized.
These features are added to the OpenTelemetry Go Trace SDK prior to
stabilization so that users can start experimenting with them and provide
feedback.

These feature may change in backwards incompatible ways as feedback is applied.
See the [Compatibility and Stability](#compatibility-and-stability) section for
more information.

## Features

- [OnEnding Processor](#onending-processor)

### OnEnding Processor

Processor implementations sometimes want to be able to modify a span after it
ended, but before it becomes immutable.
A processor that implements the `OnEnding` method can use that callback to
perform such modifications.

It can be used to implement tail-based sampling for example.

## Compatibility and Stability

Experimental features do not fall within the scope of the OpenTelemetry Go
versioning and stability [policy](../../../../VERSIONING.md).
These features may be removed or modified in successive version releases,
including patch versions.

When an experimental feature is promoted to a stable feature, a migration path
will be included in the changelog entry of the release.
