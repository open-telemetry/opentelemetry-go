package fprint

import (
	"math"

	// Our use of farmhash is sort of arbitrary: we want a fast,
	// never-changing fingerprint function and farmhash happens to
	// be a familiar one that meets those criteria.
	"github.com/dgryski/go-farm"
)

// Mix combines multiple fingerprints together.
func Mix(is ...uint64) uint64 {
	if len(is) == 0 {
		return 0
	}
	accumulator := is[0]
	for _, i := range is[1:] {
		accumulator = mix(accumulator, i)
	}
	return accumulator
}

// Borrowed from farmhash.
func mix(x uint64, y uint64) uint64 {
	const mul uint64 = 0x9ddfea08eb382d69
	a := (x ^ y) * mul
	a ^= a >> 47
	b := (y ^ a) * mul
	b ^= b >> 47
	b *= mul
	return b
}

func Bytes(s []byte) uint64 {
	return farm.Fingerprint64(s)
}

func Uint64(i uint64) uint64 {
	return i
}

func Int64(i int64) uint64 {
	return uint64(i)
}

func Int(i int) uint64 {
	return uint64(i)
}

func Float64(f float64) uint64 {
	return math.Float64bits(f)
}

func unsafeFingerprintString64(s string) uint64 {
	bs, err := unsafeStringToBytes(s)
	if err != nil {
		// TODO: Not needed, can do ... better.
		bs = []byte(s)
	}
	return Bytes(bs)
}

func String(s string) uint64 {
	// We know that the go-farm implementation
	// we use does not modify the []byte it is passed,
	// so we use an unsafe conversion here from string to
	// []byte to avoid a copy.
	return unsafeFingerprintString64(s)
}
