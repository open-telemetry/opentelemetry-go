package fprint

import (
	"fmt"
	"reflect"
	"unsafe"
)

// unsafeStringToBytes returns a []byte slice
// that has a zero-copy view into the underlying
// array backing s.
//
// Note that strings in Go are immutable, so
// the returned byte slice must not be mutated.
//
// If it's mutated, this could lead to extremely
// strange issues: for example, strings may
// be inserted in maps. If the backing array under
// a string is mutated, that string will now
// be in the map at the wrong location.
//
// This function is based off code from
// https://groups.google.com/g/golang-nuts/c/Zsfk-VMd_fU/m/O1ru4fO-BgAJ.
func unsafeStringToBytes(s string) ([]byte, error) {
	const max = 0x7fff0000 // ~2 GiB
	if len(s) > max {
		return nil, fmt.Errorf("string with length %v exceeds max %v", len(s), max)
	}
	if len(s) == 0 {
		// We need to special case a zero-length string here, because
		// below we take an exclusive slice (e.g. arr[0:len(s)]), which
		// doesn't make a whole lot of sense (the length is zero; there's no
		// data). Practically speaking, it panics. So that's why we play this game.
		return nil, nil
	}
	// type "string" is actually {uinptr, int}.
	// Cast it to that here so we can access the fields.
	//
	// See unsafe.Pointer rule #1:
	// https://golang.org/pkg/unsafe/#Pointer
	// We can perform this cast as we know
	// the structure of a string is the same
	// as reflect.StringHeader.
	//
	// We are also, additionally, using rule #6
	// as we are casting the reflect.StringHeader.Data [uintptr]
	// field to an unsafe.Pointer. We must do this so that
	// the garbage collector has a reference to the backing array.
	//
	// As well, what we're doing here is using the three-index
	// slice syntax, described here: https://golang.org/doc/go1.2#three_index
	//
	// So we declare a pointer to an array (not slice) with length `max`.
	// This is important, because an array does not have a concept of
	// capacity or length: it is just bytes. We then take the data (bytes)
	// of the string and interpret those bytes as a pointer to an array.
	//
	// To clarify, let's imagine our string is at 0xbeef. Let's imagine
	// its backing array of that string is at 0xcafe. The string is "hello".
	//
	// This is how memory looks:
	//
	// 0xbeef: 0xcafe 0x5 // {data, length} of "hello"
	// ......:
	// 0xcafe: 0x68 0x65 0x6C 0x6C 0x6F
	//
	// We then take 0xcafe, and reinterpret 0xcafe
	// as a pointer to a byte slice of length max.
	//
	// This works because there is no "array header" or similar; an array
	// is just a bunch of bytes. So we have a pointer to a bunch of bytes now,
	// which just-so-happens to match the representation of the string's
	// backing data.
	//
	// Once we have our hands on a pointer to a byte array, Go gives us something
	// sort of special, which is that the typical slice operations work on
	// a pointer to an array; they don't require an actual array (in other words,
	// there is an operator overload for pointer to byte array). This is magic.
	//
	// So, now that we have our array with a fixed size, but we've cheated and given
	// it a backing buffer that's the wrong size, we need to turn that into a slice.
	//
	// type "[]byte" (really, any slice) is actually {uintptr, int, int}.
	// (which is data, length, capacity).
	// So it needs to be {uinptr, len(s), len(s)}.
	//
	// Refer back to the three-index notation above, which is what
	// we can use to synthesize the slice with those three fields.
	// The fields in the three-index notation mean
	// Field 1: Start index for data.
	// Field 2: End index (exclusive) for data.
	// Field 3: End index (exclusive) for capacity.
	//
	// So, we want this to be 0, len(s), len(s). We want
	// our capacity and length to match exactly, so
	// three-index fields 2 and 3 match. Then we don't want to skip
	// any data at the front, so the front is 0. We omit it here
	// as it typically stylistic for Go.
	return (*[max]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
	))[:len(s):len(s)], nil
}
