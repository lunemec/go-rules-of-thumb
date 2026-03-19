package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"hash/fnv"
	"math/rand"
)

func init() {
	flag.StringVar(&variant, "variant", "", "specifies which variant to run")
}

var (
	variant string
	sizes   = []int{10, 50, 100, 500, 1_000, 5_000, 10_000, 50_000, 100_000, 500_000, 1_000_000}
	// sizesReduced is used for tests where it does not make sense to do that many iterations.
	sizesReduced = []int{1, 10, 100, 1_000}
	needles      = []int{10, 50, 100, 500, 1_000}
)

func copySlice[T any](in []T) []T {
	out := make([]T, len(in))
	copy(out, in)
	return out
}

// stableSeed derives deterministic fixture seeds from benchmark parameters so
// separate go test processes still compare identical inputs.
func stableSeed(parts ...any) uint64 {
	h := fnv.New64a()
	var buf [8]byte

	for _, part := range parts {
		switch value := part.(type) {
		case string:
			_, _ = h.Write([]byte("s:"))
			_, _ = h.Write([]byte(value))
		case int:
			_, _ = h.Write([]byte("i:"))
			binary.LittleEndian.PutUint64(buf[:], uint64(value))
			_, _ = h.Write(buf[:])
		default:
			_, _ = fmt.Fprintf(h, "%T:%v", value, value)
		}
		_, _ = h.Write([]byte{0})
	}

	return h.Sum64()
}

func newDeterministicRand(seed uint64) *rand.Rand {
	return rand.New(rand.NewSource(int64(seed)))
}

func testingSlice(size int, seed uint64) []int {
	return testingSliceFromRand(size, newDeterministicRand(seed))
}

func testingSliceFromRand(size int, rng *rand.Rand) []int {
	ts := make([]int, 0, size)
	for range size {
		ts = append(ts, rng.Intn(size))
	}
	return ts
}

func fillPattern(dst []byte, seed byte) {
	for i := range dst {
		dst[i] = seed + byte(i*13)
	}
}
