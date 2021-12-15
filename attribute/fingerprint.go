package attribute

import "go.opentelemetry.io/otel/internal/metric/fprint"

type Fingerprint uint64

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Hash(attrs ...KeyValue) Fingerprint {
	var x uint64
	for _, attr := range attrs {
		sp := fprint.String(string(attr.Key))
		switch attr.Value.Type() {
		case INVALID:
			x += fprint.Mix(sp, fprint.Int(int(INVALID)))
		case BOOL:
			x += fprint.Mix(sp, fprint.Int(b2i(attr.Value.AsBool())))
		case INT64:
			x += fprint.Mix(sp, fprint.Int64(attr.Value.AsInt64()))
		case FLOAT64:
			x += fprint.Mix(sp, fprint.Float64(attr.Value.AsFloat64()))
		case STRING:
			x += fprint.Mix(sp, fprint.String(attr.Value.AsString()))
		case BOOLSLICE:
			t := attr.Value.AsBoolSlice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, b := range t {
				x = fprint.Mix(x, fprint.Int(b2i(b)))
			}
		case INT64SLICE:
			t := attr.Value.AsInt64Slice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, i := range t {
				x = fprint.Mix(x, fprint.Int64(i))
			}
		case FLOAT64SLICE:
			t := attr.Value.AsFloat64Slice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, f := range t {
				x = fprint.Mix(x, fprint.Float64(f))
			}
		case STRINGSLICE:
			t := attr.Value.AsStringSlice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, s := range t {
				x = fprint.Mix(x, fprint.String(s))
			}
		}
	}
	return Fingerprint(x)
}
