package registry

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/instrument"
)

type hasDescriptor interface {
	Descriptor() sdkapi.Descriptor
}

type State struct {
	names map[string]hasDescriptor
}

var ErrIncompatibleInstruments = fmt.Errorf("incompatible instrument registration")

func New() *State {
	return &State{
		names: map[string]hasDescriptor{},
	}
}

func Lookup[T hasDescriptor](reg *State, name string, opts []instrument.Option, nk number.Kind, ik sdkapi.InstrumentKind,
	f func(desc sdkapi.Descriptor) T) (T, error) {
	cfg := instrument.NewConfig(opts...)
	desc := sdkapi.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())

	lookup := reg.names[name]

	if lookup != nil {
		hasD, ok := lookup.(T)
		if ok {
			exist := hasD.Descriptor()
			
			if exist.NumberKind() == nk && exist.InstrumentKind() == ik && exist.Unit() == cfg.Unit() {
				return hasD, nil
			}
		}

		
		var t T
		return t, ErrIncompatibleInstruments
	}
	value := f(desc)
	reg.names[name] = value
	return value, nil
}
