package tag

import (
	"context"
	"unsafe"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type (
	registeredKey struct {
		name    string
		desc    string
		unit    unit.Unit
		eventID core.EventID
	}

	ctxTagsType struct{}

	measure struct {
		rk *registeredKey
	}
)

var (
	ctxTagsKey = &ctxTagsType{}
)

func register(name string, opts []Option) *registeredKey {
	rk := &registeredKey{
		name: name,
	}
	for _, of := range opts {
		of(rk)
	}
	rk.eventID = observer.Record(observer.Event{
		Type:   observer.NEW_MEASURE,
		String: name,
		// TODO desc, unit
	})
	return rk
}

func (k *registeredKey) Name() string {
	if k == nil {
		return "unregistered"
	}
	return k.name
}

func (k *registeredKey) Description() string {
	if k == nil {
		return ""
	}
	return k.desc
}

func (k *registeredKey) Unit() unit.Unit {
	if k == nil {
		return unit.Dimensionless
	}
	return k.unit
}

func (k *registeredKey) DefinitionID() core.EventID {
	if k == nil {
		return 0
	}
	return k.eventID
}

func (k *registeredKey) Bool(v bool) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type: core.BOOL,
			Bool: v,
		},
	}
}

func (k *registeredKey) Int64(v int64) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:  core.INT64,
			Int64: v,
		},
	}
}

func (k *registeredKey) Uint64(v uint64) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:   core.UINT64,
			Uint64: v,
		},
	}
}

func (k *registeredKey) Float64(v float64) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:    core.FLOAT64,
			Float64: v,
		},
	}
}

func (k *registeredKey) Int32(v int32) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:  core.INT32,
			Int64: int64(v),
		},
	}
}

func (k *registeredKey) Uint32(v uint32) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:   core.UINT32,
			Uint64: uint64(v),
		},
	}
}

func (k *registeredKey) Float32(v float32) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:    core.FLOAT32,
			Float64: float64(v),
		},
	}
}

func (k *registeredKey) String(v string) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:   core.STRING,
			String: v,
		},
	}
}

func (k *registeredKey) Bytes(v []byte) core.KeyValue {
	return core.KeyValue{
		Key: k,
		Value: core.Value{
			Type:  core.BYTES,
			Bytes: v,
		},
	}
}

func (k *registeredKey) Int(v int) core.KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Int32(int32(v))
	}
	return k.Int64(int64(v))
}

func (k *registeredKey) Uint(v uint) core.KeyValue {
	if unsafe.Sizeof(v) == 4 {
		return k.Uint32(uint32(v))
	}
	return k.Uint64(uint64(v))
}

func (k *registeredKey) Value(ctx context.Context) core.KeyValue {
	v, _ := FromContext(ctx).Value(k)
	return core.KeyValue{
		Key:   k,
		Value: v,
	}
}

func (m measure) M(v float64) core.Measurement {
	return core.Measurement{
		Measure: m,
		Value:   v,
	}
}

func (m measure) V(v float64) core.KeyValue {
	return m.rk.Float64(v)
}

func (m measure) Name() string               { return m.rk.Name() }
func (m measure) Description() string        { return m.rk.Description() }
func (m measure) Unit() unit.Unit            { return m.rk.Unit() }
func (m measure) DefinitionID() core.EventID { return m.rk.DefinitionID() }
