package push

import (
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Config contains configuration for a push Controller.
type Config struct {
	// ErrorHandler is the function called when the Controller encounters an error.
	//
	// This option can be overridden after instantiation of the Controller
	// with the `SetErrorHandler` method.
	ErrorHandler sdk.ErrorHandler

	// Resource is the OpenTelemetry resource associated with all Meters
	// created by the Controller.
	Resource resource.Resource
}

// Option is the interface that applies the value to a configuration option.
type Option interface {
	// Apply sets the Option value of a Config.
	Apply(*Config)
}

// WithErrorHandler sets the ErrorHandler configuration option of a Config.
func WithErrorHandler(fn sdk.ErrorHandler) Option {
	return errorHandlerOption(fn)
}

type errorHandlerOption sdk.ErrorHandler

func (o errorHandlerOption) Apply(config *Config) {
	config.ErrorHandler = sdk.ErrorHandler(o)
}

// WithResource sets the Resource configuration option of a Config.
func WithResource(r resource.Resource) Option {
	return resourceOption(r)
}

type resourceOption resource.Resource

func (o resourceOption) Apply(config *Config) {
	config.Resource = resource.Resource(o)
}
