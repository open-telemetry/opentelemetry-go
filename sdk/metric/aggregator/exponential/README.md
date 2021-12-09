# Base-2 Exponential Histogram

## Design

This is a fixed-size data structure for aggregating the OpenTelemetry
base-2 exponential histogram introduced in [OTEP
149](https://github.com/open-telemetry/oteps/blob/main/text/0149-exponential-histogram.md)
and [described in the metrics data
model](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/datamodel.md#exponentialhistogram).
The exponential histogram data point is characterized by a `scale`
factor that determines resolution.  Positive scales correspond with
more resolution, and negatives scales correspond with less resolution.

Given a maximum size, in terms of the number of buckets, the
implementation determines the best scale possible given the set of
measurements received.  The size of the histogram is configured using
the `WithMaxSize()` option, which defaults to 320.

An optional configuration supports fixing the scale in advance, which
ensures that repeated collection periods will generate consistent
histogram bucket boundaries, across multiple processes.  This option,
set by `WithRangeLimit(min, max)`, fixes the scale parameter and is 
recommended in configurations that write through to Prometheus.

When range limits are not fixed, the implementation here maintains the
best resolution possible.  Since the scale parameter is shared by the
positive and negative ranges, the best value of the scale parameter is
determined by the range with the greater difference between minimum
and maximum bucket index:

```golang
func bucketsNeeded(minValue, maxValue float64, scale int32) int32 {
	return bucketIndex(maxValue, scale) - bucketIndex(minValue, scale) + 1
}

func bucketIndex(value float64, scale int32) int32 {
	return math.Log(value) * math.Ldexp(math.Log2E, scale)
}
```

The best scale is uniquely determined when `maxSize/2 <
bucketsNeeded(minValue, maxValue, scale) <= maxSize`.  This
implementation maintains the best scale by rescaling as needed to stay
within the maximum size.

### UpdateByIncr interface

The OpenTelemetry metrics SDK `Aggregator` type supports an `Update()`
interface which implies updating the histogram by a count of 1.  This
implementation also supports `UpdateByIncr()`, which makes it possible
to support counting multiple observations in a single API call.  This
extension is useful in applying `Histogram` aggregation to _sampled_
metric events (e.g. in the [OpenTelemetry statsd
receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/statsdreceiver)).

Another use for `UpdateByIncr` is in a Span-to-metrics pipeline
following [probability sampling in OpenTelemetry tracing
(WIP)](https://github.com/open-telemetry/opentelemetry-specification/pull/2047).

## Implementation

The implementation maintains a slice of buckets and grows the array in
size only as necessary given the actual range of values, up to the
maximum size.  The structure of a single range of buckets is:

```golang
type buckets struct {
	backing    interface{}
	indexBase  int32
	indexStart int32
	indexEnd   int32
}
```

The `backing` field is a slice of variable width unsigned integer
(i.e., `[]uint8`, `[]uint16`, `[]uint32`, or `[]uint64`.  The
`indexStart` and `indexEnd` fields store the current minimum and
maximum bucket indices for the current scale.

The backing array is circular.  When the first observation is added to
a set of (positive or negative) buckets, the initial conditition is
`indexBase == indexStart == indexEnd`.  When new observations are
added at indices lower than `indexStart` and while capacity is greater
than `indexEnd - indexBase`, new values are filled in by adjusting
`indexStart` to be less than `indexBase`. This mechanism allows the
backing array to grow in either direction without moving values, up
until rescaling is necessary.

The positive and negative backing arrays are independent, so the
maximum space used for `buckets` by one `Aggregator` is twice the
configured maximum size.

### Internal mapping function

There are two mapping functions used, depending on the sign of the
scale.  Negative and zero scales use the `internal/mapping/exponent`
mapping function, which computes the bucket index directly from the
bits of the `float64` exponent.  This mapping function is used with
scale `-10 <= scale <= 0`.  Scales smaller than -10 map the entire
normal `float64` numner range into a single bucket, thus are not
considered useful.

The `internal/mapping/logarithm` mapping function uses
`math.Log(value)` times the scaling factor `math.Ldexp(math.Log2E,
scale)`.  This mapping function is used with `0 < scale <= 20`.
Scales larger than 20 exceed the resolution achievable using this
calculation method.

### Determining change of scale

The algorithm used to determine the (best) change of scale when a new
value arrives is:

```golang
func newScale(minIndex, maxIndex, scale, maxSize int32) int32 {
    return scale - changeScale(minIndex, maxIndex, scale, maxSize)
}

func changeScale(minIndex, maxIndex, scale, maxSize int32) int32 {
    var change int32
    for maxIndex - lowIndex >= maxSize {
	   maxIndex >>= 1
	   minIndex >>= 1
	   change++
    }
	return change
}
```

The `changeScale` function is also used to determine how many bits to
shift during `Merge` and to fix the initial scale when range limits
are configured.

### Merge function

The merge function rotates the circular backing array so that
`indexStart == indexBase`, using the "3 reversals" method, which
simplifies internal logic.  `Merge` uses the `UpdateByIncr` code path
to combine one `buckets` into another.

### Scale function

The `Scale` function returns the current scale of the histogram.  If
the scale is variable and there are no non-zero values in the
histogram, the scale is zero by definition.  If the scale is fixed
because of range limits, the fixed scale will be returned even for the
empty histogram.

### Handling subnormal values

Subnormal values are those in the range [0x1p-1074, 0x1p-1022), these
being numbers that "gradually underflow" and use less than 52 bits of
precision in the significand at the smallest representable exponent
(i.e., -1022).  Subnormal numbers present special challenges for both
the exponent- and logarithm-based mapping function, and to avoid
additional complexity induced by corner cases, subnormal numbers are
rounded up to 0x1p-1022 in this implementation.

Handling subnormal numbers is difficult for the logarithm mapping
function because Golang's `math.Log()` function rounds subnormal
numbers up to 0x1p-1022.  Handling subnormal numbers is difficult for
the exponent mapping function because Golang's `math.Frexp()`, the
natural API for extracting a value's base-2 exponent, also rounds
subnormal numbers up to 0x1p-1022.

While the additional complexity needed to correctly map subnormal
numbers is small in both cases, there are few real benefits in doing
so because of the inherent loss of precision.  As secondary
motivation, clamping values to the range [0x1p-1022, math.MaxFloat64]
increases symmetry. This limit means that minimum bucket index and the
maximum bucket index have similar magnitude, which helps support
greater maximum scale.  Supporting numbers smaller than 0x1p-1022
would mean changing the valid scale interval to [-11,19] compared with
[-10,20].

## Acknowledgements

This implementation is based on work by [Yuke
Zhuge](https://github.com/yzhuge) and [Otmar
Ertl](https://github.com/oertl).  See
[NrSketch](https://github.com/newrelic-experimental/newrelic-sketch-java/blob/1ce245713603d61ba3a4510f6df930a5479cd3f6/src/main/java/com/newrelic/nrsketch/indexer/LogIndexer.java)
and
[DynaHist](https://github.com/dynatrace-oss/dynahist/blob/9a6003fd0f661a9ef9dfcced0b428a01e303805e/src/main/java/com/dynatrace/dynahist/layout/OpenTelemetryExponentialBucketsLayout.java)
repositories for more detail.
