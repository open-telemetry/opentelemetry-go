package pipeline

// Register is a per-pipeline slice of data.  Although it is nothing
// more than a simple slice, the use of this generic wrapper serves as
// a notice that the associated data will be indexed by an integer,
// generally named "pipe", referring to the ordinal position of the
// reader in the list of readers.
type Register[T any] []T

// NewRegister returns a new slice of per-pipeline data.
func NewRegister[T any](size int) Register[T] {
	return Register[T](make([]T, size))
}
