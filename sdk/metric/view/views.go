package view

// Views is a configured set of view clauses with an associated Name
// that is used for debugging.
type Views struct {
	// Name of these views, used in error reporting.
	Name string

	// Config is the configuration for these views.
	Config
}

// New configures the clauses and default settings of a Views.
func New(name string, opts ...Option) *Views {
	return &Views{
		Name:   name,
		Config: NewConfig(opts...),
	}
}

// TODO: call views.Validate() to check for:
// - empty (?)
// - duplicate name
// - invalid inst/number/aggregation kind
// - both instrument name and regexp
// - schemaURL or Version without library name
// - empty attribute keys
// - Name w/o SingleInst
