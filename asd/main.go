package main

import (
	"sync/atomic"
)

type A struct {
	_ int64 // b is accessed via a sync/atomic function
	_ int32
	d int64 // d should NOT be accessed via a sync/atomic function as it's not 64bits aligned.
}

func main() {
	a := A{}
	// Atomically increment a.b. This is valid, b is the first field of the struct, so it's
	// guaranteed to be 64-bits aligned.
	atomic.AddInt64(&a.d, 1)
}
