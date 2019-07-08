package core

import (
	"unsafe"
)

func (k *Key) Bool(v bool) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type: BOOL,
			Bool: v,
		},
	}
}

func (k *Key) Int64(v int64) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:  INT64,
			Int64: v,
		},
	}
}

func (k *Key) Uint64(v uint64) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:   UINT64,
			Uint64: v,
		},
	}
}

func (k *Key) Float64(v float64) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:    FLOAT64,
			Float64: v,
		},
	}
}

func (k *Key) Int32(v int32) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:  INT32,
			Int64: int64(v),
		},
	}
}

func (k *Key) Uint32(v uint32) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:   UINT32,
			Uint64: uint64(v),
		},
	}
}

func (k *Key) Float32(v float32) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:    FLOAT32,
			Float64: float64(v),
		},
	}
}

func (k *Key) String(v string) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:   STRING,
			String: v,
		},
	}
}

func (k *Key) Bytes(v []byte) KeyValue {
	return KeyValue{
		Key: *k,
		Value: Value{
			Type:  BYTES,
			Bytes: v,
		},
	}
}

func (k *Key) Int(v int) KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Int32(int32(v))
	}
	return k.Int64(int64(v))
}

func (k *Key) Uint(v uint) KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Uint32(uint32(v))
	}
	return k.Uint64(uint64(v))
}

func (m *Measure) M(v float64) Measurement {
	return Measurement{
		Measure: *m,
		Value:   v,
	}
}

func (m *Measure) V(v float64) KeyValue {
	k := Key{Variable: m.Variable}
	return k.Float64(v)
}

func (k *Key) Defined() bool {
	return len(k.Variable.Name) != 0
}
