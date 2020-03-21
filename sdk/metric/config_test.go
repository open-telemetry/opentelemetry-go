package metric

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestWithErrorHandler(t *testing.T) {
	errH, reg := func() (ErrorHandler, *error) {
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
	c = &Config{ErrorHandler: DefaultErrorHandler}
	WithErrorHandler(errH).Apply(c)
	err2 := fmt.Errorf("error 2")
	c.ErrorHandler(err2)
	assert.EqualError(t, *reg, err2.Error())
}

func TestWithResource(t *testing.T) {
	r := resource.New(core.Key("A").String("a"))

	c := &Config{}
	WithResource(*r).Apply(c)
	assert.Equal(t, *r, c.Resource)

	// Ensure overwriting works.
	c = &Config{Resource: resource.Resource{}}
	WithResource(*r).Apply(c)
	assert.Equal(t, *r, c.Resource)
}
