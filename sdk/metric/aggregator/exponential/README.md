# Base-2 Exponential Histogram

## Design

This document is a placeholder for future Aggregator, once seen in [PR
2393](https://github.com/open-telemetry/opentelemetry-go/pull/2393).

Only the mapping functions have been made available at this time.  The
equations tested here are specified in the [data model for Exponential
Histogram data points](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/datamodel.md#exponentialhistogram).

### Mapping function

There are two mapping functions used, depending on the sign of the
scale.  Negative and zero scales use the `mapping/exponent` mapping
function, which computes the bucket index directly from the bits of
the `float64` exponent.  This mapping function is used with scale `-10
<= scale <= 0`.  Scales smaller than -10 map the entire normal
`float64` number range into a single bucket, thus are not considered
useful.

The `mapping/logarithm` mapping function uses `math.Log(value)` times
the scaling factor `math.Ldexp(math.Log2E, scale)`.  This mapping
function is used with `0 < scale <= 20`.  The maximum scale is
selected because at scale 21, simply, it becomes difficult to test
correctness--at this point `math.MaxFloat64` maps to index
`math.MaxInt32` and the `math/big` logic used in testing breaks down.
