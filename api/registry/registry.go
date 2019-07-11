// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type Sequence uint64

type Option func(Variable) Variable

type Variable struct {
	Name        string
	Description string
	Unit        unit.Unit
	Type        Type
}

type Type interface {
	String() string
}

func Register(name string, vtype Type, opts ...Option) Variable {
	return newVar(name, vtype, opts...)
}

func newVar(name string, vtype Type, opts ...Option) Variable {
	v := Variable{
		Name: name,
	}
	for _, o := range opts {
		v = o(v)
	}
	return v
}

func (v *Variable) Defined() bool {
	return len(v.Name) != 0
}

// WithDescription applies the provided description.
func WithDescription(desc string) Option {
	return func(v Variable) Variable {
		v.Description = desc
		return v
	}
}

// WithUnit applies the provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(v Variable) Variable {
		v.Unit = unit
		return v
	}
}
