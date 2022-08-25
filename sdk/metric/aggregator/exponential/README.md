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
the `WithMaxSize()` option, which defaults to 160.

The implementation here maintains the best resolution possible.  Since
the scale parameter is shared by the positive and negative ranges, the
best value of the scale parameter is determined by the range with the
greater difference between minimum and maximum bucket index:

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

## Layout

### Mapping function

The `mapping` sub-package contains the equations specified in the [data
model for Exponential Histogram data
points](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/data-model.md#exponentialhistogram).

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

### Data structure

The `structure` sub-package contains a Histogram aggregator for use by
the OpenTelemetry-Go Metrics SDK as well as OpenTelemetry Collector
receivers, processors, and exporters.

## Implementation

The implementation maintains a slice of buckets and grows the array in
size only as necessary given the actual range of values, up to the
maximum size.  The structure of a single range of buckets is:

```golang
type buckets struct {
	backing    bucketsVarwidth[T]  // for T = uint8 | uint16 | uint32 | uint64
	indexBase  int32
	indexStart int32
	indexEnd   int32
}
```

The `backing` field is a generic slice of `[]uint8`, `[]uint16`,
`[]uint32`, or `[]uint64`.

The positive and negative backing arrays are independent, so the
maximum space used for `buckets` by one `Aggregator` is twice the
configured maximum size.

### Backing array

The backing array is circular.  The first observation is counted in
the 0th index of the backing array and the initial bucket number is
stored in `indexBase`.  After the initial observation, the backing
array grows in either direction (i.e., larger or smaller bucket
numbers), until rescaling is necessary.  This mechanism allows the
histogram to maintain the ideal scale without shifting values inside
the array.

The `indexStart` and `indexEnd` fields store the current minimum and
maximum bucket number.  The initial condition is `indexBase ==
indexStart == indexEnd`, representing a single bucket.

Following the first observation, new observations may fall into a
bucket up to `size-1` in either direction.  Growth is possible by
adjusting either `indexEnd` or `indexStart` as long as the constraint
`indexEnd-indexStart < size` remains true.

Bucket numbers in the range `[indexBase, indexEnd]` are stored in the
interval `[0, indexEnd-indexBase]` of the backing array.  Buckets in
the range `[indexStart, indexBase-1]` are stored in the interval
`[size+indexStart-indexBase, size-1]` of the backing array.

Considering the `aggregation.Buckets` interface, `Offset()` returns
`indexStart`, `Len()` returns `indexEnd-indexStart+1`, and `At()`
locates the correct bucket in the circular array.

### Determining change of scale

The algorithm used to determine the (best) change of scale when a new
value arrives is:

```golang
func newScale(minIndex, maxIndex, scale, maxSize int32) int32 {
    return scale - changeScale(minIndex, maxIndex, scale, maxSize)
}

func changeScale(minIndex, maxIndex, scale, maxSize int32) int32 {
    var change int32
    for maxIndex - minIndex >= maxSize {
	   maxIndex >>= 1
	   minIndex >>= 1
	   change++
    }
	return change
}
```

The `changeScale` function is also used to determine how many bits to
shift during `Merge`.

### Downscale function

The downscale function rotates the circular backing array so that
`indexStart == indexBase`, using the "3 reversals" method, before
combining the buckets in place.

### Merge function

`Merge` first calculates the correct final scale by comparing the
combined positive and negative ranges.  The destination aggregator is
then downscaled, if necessary, and the `UpdateByIncr` code path to add
the source buckets to the destination buckets.

### Scale function

The `Scale` function returns the current scale of the histogram.

If the scale is variable and there are no non-zero values in the
histogram, the scale is zero by definition; when there is only a
single value in this case, its scale is MinScale (20) by definition.

If the scale is fixed because of range limits, the fixed scale will be
returned even for any size histogram.

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

## Acknowledgements

This implementation is based on work by [Yuke
Zhuge](https://github.com/yzhuge) and [Otmar
Ertl](https://github.com/oertl).  See
[NrSketch](https://github.com/newrelic-experimental/newrelic-sketch-java/blob/1ce245713603d61ba3a4510f6df930a5479cd3f6/src/main/java/com/newrelic/nrsketch/indexer/LogIndexer.java)
and
[DynaHist](https://github.com/dynatrace-oss/dynahist/blob/9a6003fd0f661a9ef9dfcced0b428a01e303805e/src/main/java/com/dynatrace/dynahist/layout/OpenTelemetryExponentialBucketsLayout.java)
repositories for more detail.
