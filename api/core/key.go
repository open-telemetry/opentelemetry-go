package core

import (
	"fmt"
	"unsafe"
)

type Key struct {
	Name string
}

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
	return len(k.Name) != 0
}

// TODO make this a lazy one-time conversion.
func (v Value) Emit() string {
	switch v.Type {
	case BOOL:
		return fmt.Sprint(v.Bool)
	case INT32, INT64:
		return fmt.Sprint(v.Int64)
	case UINT32, UINT64:
		return fmt.Sprint(v.Uint64)
	case FLOAT32, FLOAT64:
		return fmt.Sprint(v.Float64)
	case STRING:
		return v.String
	case BYTES:
		return string(v.Bytes)
	}
	return "unknown"
}
