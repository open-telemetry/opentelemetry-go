package core

import (
	"strconv"
	"strings"
	"unsafe"
)

type Key string

type KeyValue struct {
	Key   Key
	Value Value
}

type ValueType int

type Value struct {
	Type    ValueType
	Bool    bool
	Int64   int64
	Uint64  uint64
	Float64 float64
	String  string
	Bytes   []byte

	// TODO See how segmentio/stats handles this type, it's much smaller.
	// TODO Lazy value type?
}

const (
	INVALID ValueType = iota
	BOOL
	INT32
	INT64
	UINT32
	UINT64
	FLOAT32
	FLOAT64
	STRING
	BYTES
)

func (k Key) Bool(v bool) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type: BOOL,
			Bool: v,
		},
	}
}

func (k Key) Int64(v int64) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:  INT64,
			Int64: v,
		},
	}
}

func (k Key) Uint64(v uint64) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:   UINT64,
			Uint64: v,
		},
	}
}

func (k Key) Float64(v float64) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:    FLOAT64,
			Float64: v,
		},
	}
}

func (k Key) Int32(v int32) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:  INT32,
			Int64: int64(v),
		},
	}
}

func (k Key) Uint32(v uint32) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:   UINT32,
			Uint64: uint64(v),
		},
	}
}

func (k Key) Float32(v float32) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:    FLOAT32,
			Float64: float64(v),
		},
	}
}

func (k Key) String(v string) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:   STRING,
			String: v,
		},
	}
}

func (k Key) Bytes(v []byte) KeyValue {
	return KeyValue{
		Key: k,
		Value: Value{
			Type:  BYTES,
			Bytes: v,
		},
	}
}

func (k Key) Int(v int) KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Int32(int32(v))
	}
	return k.Int64(int64(v))
}

func (k Key) Uint(v uint) KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Uint32(uint32(v))
	}
	return k.Uint64(uint64(v))
}

func (k Key) Defined() bool {
	return len(k) != 0
}

// Encoder supports formatting values without allocating temporary strings.
type Encoder interface {
	Write([]byte) (int, error)
	WriteRune(rune) (int, error)
	WriteString(string) (int, error)
}

// Emit constructs a new string of this value, whereas Encode() emits
// the string into an existing buffer.
func (v Value) Emit() string {
	var sb strings.Builder
	var tmp [32]byte
	_, _ = v.Encode(&sb, tmp[:])
	return sb.String()
}

// Encode writes the encoded value to `w`.  `tmp` provides a temporary
// buffer for this to work.  A correct result is ensured, even if the
// temporary buffer is too small.  `tmp` should be at least 32 bytes to
// ensure no allocations.
func (v Value) Encode(w Encoder, tmp []byte) (n int, err error) {
	tmp = tmp[:0]
	switch v.Type {
	case BOOL:
		if v.Bool {
			return w.WriteString("true")
		}
		return w.WriteString("false")
	case INT32, INT64:
		tmp = strconv.AppendInt(tmp, v.Int64, 10)
	case UINT32, UINT64:
		tmp = strconv.AppendUint(tmp, v.Uint64, 10)
	case FLOAT32, FLOAT64:
		// Note: 'x' format is much faster
		tmp = strconv.AppendFloat(tmp, v.Float64, 'g', OutputFloatingPointPrecision, 64)
	case STRING:
		return w.WriteString(v.String)
	case BYTES:
		// TODO This must be removed, it's not safe
		return w.Write(v.Bytes)
	default:
		// INVALID -> empty string
	}
	return w.Write(tmp)
}
