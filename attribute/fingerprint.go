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

func (aa Attributes) Equals(bb Attributes) bool {
	if aa.Fingerprint != bb.Fingerprint {
		return false
	}
	a := aa.KeyValues
	b := bb.KeyValues

	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Key != b[i].Key {
			return false
		}
		if a[i].Value.Type() != b[i].Value.Type() {
			return false
		}
		switch a[i].Value.Type() {
		case INVALID, BOOL, INT64,
			FLOAT64, STRING:
			if a[i].Value != b[i].Value {
				return false
			}
		case BOOLSLICE:
			as := a[i].Value.AsBoolSlice()
			bs := b[i].Value.AsBoolSlice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case INT64SLICE:
			as := a[i].Value.AsInt64Slice()
			bs := b[i].Value.AsInt64Slice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case FLOAT64SLICE:
			as := a[i].Value.AsFloat64Slice()
			bs := b[i].Value.AsFloat64Slice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case STRINGSLICE:
			as := a[i].Value.AsStringSlice()
			bs := b[i].Value.AsStringSlice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		}
	}
	return true
}
