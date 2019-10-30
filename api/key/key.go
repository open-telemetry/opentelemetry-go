package key

import (
	"go.opentelemetry.io/api/core"
)

func New(name string) core.Key {
	return core.Key(name)
}

func Bool(k string, v bool) core.KeyValue {
	return New(k).Bool(v)
}

func Int64(k string, v int64) core.KeyValue {
	return New(k).Int64(v)
}

func Uint64(k string, v uint64) core.KeyValue {
	return New(k).Uint64(v)
}

func Float64(k string, v float64) core.KeyValue {
	return New(k).Float64(v)
}

func Int32(k string, v int32) core.KeyValue {
	return New(k).Int32(v)
}

func Uint32(k string, v uint32) core.KeyValue {
	return New(k).Uint32(v)
}

func Float32(k string, v float32) core.KeyValue {
	return New(k).Float32(v)
}

func String(k, v string) core.KeyValue {
	return New(k).String(v)
}

func Int(k string, v int) core.KeyValue {
	return New(k).Int(v)
}

func Uint(k string, v uint) core.KeyValue {
	return New(k).Uint(v)
}
