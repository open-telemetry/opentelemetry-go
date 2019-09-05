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

package tag

import (
	"go.opentelemetry.io/api/core"
)

type MutatorOp int

const (
	INSERT MutatorOp = iota
	UPDATE
	UPSERT
	DELETE
)

type Mutator struct {
	MutatorOp
	core.KeyValue
	MeasureMetadata
}

func (m Mutator) WithTTL(hops int) Mutator {
	m.TTL = hops
	return m
}

func Insert(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: INSERT,
		KeyValue:  kv,
	}
}

func Update(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: UPDATE,
		KeyValue:  kv,
	}
}

func Upsert(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: UPSERT,
		KeyValue:  kv,
	}
}

func Delete(k core.Key) Mutator {
	return Mutator{
		MutatorOp: DELETE,
		KeyValue: core.KeyValue{
			Key: k,
		},
	}
}
