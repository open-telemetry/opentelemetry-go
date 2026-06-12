Fixes #7956

Per the OTel spec, attribute value limits must be applied recursively to map values. Previously `truncateAttr` handled `STRING`, `STRINGSLICE`, `BYTESLICE`, and `SLICE`, but not `MAP`.

Add `MAP` handling to the trace SDK attribute limiter:
- `truncateAttr`: truncates `attribute.MAP` values when needed.
- `truncateValue`: recursively truncates `MAP`, `SLICE`, `STRINGSLICE`, `STRING`, and `BYTESLICE` values.
- `needsTruncation`: pre-scan guard for `MAP` values to avoid rebuilding when no nested value changes.

```
goos: linux
goarch: amd64
pkg: go.opentelemetry.io/otel/sdk/trace
cpu: 13th Gen Intel(R) Core(TM) i7-13800H
BenchmarkSpanLimits/None-20                          223563  10077 ns/op  12448 B/op  38 allocs/op
BenchmarkSpanLimits/AttributeValueLengthLimit-20      89922  11616 ns/op  13990 B/op  56 allocs/op
BenchmarkSpanLimits/AttributeCountLimit-20           113170  10940 ns/op  11616 B/op  38 allocs/op
BenchmarkSpanLimits/EventCountLimit-20               133573   8994 ns/op  11376 B/op  35 allocs/op
BenchmarkSpanLimits/LinkCountLimit-20                235468   8065 ns/op  10976 B/op  35 allocs/op
BenchmarkSpanLimits/AttributePerEventCountLimit-20   266492  11591 ns/op  12448 B/op  38 allocs/op
BenchmarkSpanLimits/AttributePerLinkCountLimit-20    101046  10433 ns/op  12448 B/op  38 allocs/op
```

`benchstat` for `BenchmarkSpanLimits/AttributeValueLengthLimit` (`count=10`):

```
SpanLimits/AttributeValueLengthLimit-20   14.88µ ± 16%   15.56µ ± 53%       ~ (p=0.739 n=10)
SpanLimits/AttributeValueLengthLimit-20   12.71Ki ± 0%   13.66Ki ± 0%   +7.52% (p=0.000 n=10)
SpanLimits/AttributeValueLengthLimit-20    47.00 ± 0%     56.00 ± 0%   +19.15% (p=0.000 n=10)
```

The extra allocations are expected in the truncation path because truncated `MAP` values are rebuilt recursively.
