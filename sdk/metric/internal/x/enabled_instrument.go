package x

import "context"

// EnabledInstrument interface is implemented by synchronous instruments.
type EnabledInstrument interface {
	// Enabled reports whether the instrument will process measurements for the given context.
	//
	// This function can be used in places where measuring an instrument
	// would result in computationally expensive operations.
	Enabled(context.Context) bool
}
