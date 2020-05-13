// Copyright The OpenTelemetry Authors
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

package push

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/kv"

	"github.com/stretchr/testify/assert"

	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestWithErrorHandler(t *testing.T) {
	errH, reg := func() (sdk.ErrorHandler, *error) {
		e := fmt.Errorf("default invalid")
		reg := &e
		return func(err error) {
			*reg = err
		}, reg
	}()

	c := &Config{}
	WithErrorHandler(errH).Apply(c)
	err1 := fmt.Errorf("error 1")
	c.ErrorHandler(err1)
	assert.EqualError(t, *reg, err1.Error())

	// Ensure overwriting works.
	c = &Config{ErrorHandler: sdk.DefaultErrorHandler}
	WithErrorHandler(errH).Apply(c)
	err2 := fmt.Errorf("error 2")
	c.ErrorHandler(err2)
	assert.EqualError(t, *reg, err2.Error())
}

func TestWithResource(t *testing.T) {
	r := resource.New(kv.String("A", "a"))

	c := &Config{}
	WithResource(r).Apply(c)
	assert.Equal(t, r.Equivalent(), c.Resource.Equivalent())

	// Ensure overwriting works.
	c = &Config{Resource: &resource.Resource{}}
	WithResource(r).Apply(c)
	assert.Equal(t, r.Equivalent(), c.Resource.Equivalent())
}
