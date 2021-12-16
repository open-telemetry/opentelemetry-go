package attribute

import "go.opentelemetry.io/otel/internal/fprint"

type (
	Attributes struct {
		Fingerprint uint64
		KeyValues   []KeyValue
	}
)

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Fingerprint(kvs ...KeyValue) Attributes {
	var x uint64
	for _, kv := range kvs {
		sp := fprint.String(string(kv.Key))
		switch kv.Value.Type() {
		case INVALID:
			x += fprint.Mix(sp, fprint.Int(int(INVALID)))
		case BOOL:
			x += fprint.Mix(sp, fprint.Int(b2i(kv.Value.AsBool())))
		case INT64:
			x += fprint.Mix(sp, fprint.Int64(kv.Value.AsInt64()))
		case FLOAT64:
			x += fprint.Mix(sp, fprint.Float64(kv.Value.AsFloat64()))
		case STRING:
			x += fprint.Mix(sp, fprint.String(kv.Value.AsString()))
		case BOOLSLICE:
			t := kv.Value.AsBoolSlice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, b := range t {
				x = fprint.Mix(x, fprint.Int(b2i(b)))
			}
		case INT64SLICE:
			t := kv.Value.AsInt64Slice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, i := range t {
				x = fprint.Mix(x, fprint.Int64(i))
			}
		case FLOAT64SLICE:
			t := kv.Value.AsFloat64Slice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, f := range t {
				x = fprint.Mix(x, fprint.Float64(f))
			}
		case STRINGSLICE:
			t := kv.Value.AsStringSlice()
			x += fprint.Mix(sp, fprint.Int(len(t)))
			for _, s := range t {
				x = fprint.Mix(x, fprint.String(s))
			}
		}
	}
	return Attributes{
		Fingerprint: x,
		KeyValues:   kvs,
	}
}
