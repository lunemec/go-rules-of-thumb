package main

import (
	"fmt"
	"iter"
	"math/rand"
	"testing"
)

type iterateBench func(int) func(*testing.B)

var iterateSink int

func BenchmarkIterate(b *testing.B) {
	runBenchmark := func(runF iterateBench) {
		for _, size := range sizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				runF(size),
			)
		}
	}
	switch variant {
	case "slice_iterate":
		runBenchmark(benchmarkSliceIterate)
	case "range_func":
		runBenchmark(benchmarkRangeFuncIterate)
	case "direct":
		runBenchmark(benchmarkDirect)
	default:
		b.Errorf("speficy which test to run: -args -test slice_iterate|range_func|direct")
	}
}

func benchmarkSliceIterate(size int) func(*testing.B) {
	seed := stableSeed("BenchmarkIterate", size)
	rng := newDeterministicRand(seed)

	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			rng.Seed(int64(seed))
			// Here we have to allocate the iteration slice
			// as that is the main benefit of range over func
			// - no upfront allocation.
			iterSlice := testingSliceFromRand(size, rng)
			for _, val := range iterSlice {
				acc += val
			}
		}

		iterateSink = acc
	}
}

func benchmarkDirect(size int) func(*testing.B) {
	seed := stableSeed("BenchmarkIterate", size)
	rng := newDeterministicRand(seed)

	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			rng.Seed(int64(seed))
			for range size {
				acc += rng.Intn(size)
			}
		}

		iterateSink = acc
	}
}

func benchmarkRangeFuncIterate(size int) func(*testing.B) {
	seed := stableSeed("BenchmarkIterate", size)
	rng := newDeterministicRand(seed)

	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			rng.Seed(int64(seed))
			for val := range testingIter(size, rng) {
				acc += val
			}
		}

		iterateSink = acc
	}
}

func TestIter(t *testing.T) {
	size := 10
	seed := stableSeed("TestIter", size)
	var accum []int

	for val := range testingIter(size, newDeterministicRand(seed)) {
		accum = append(accum, val)
	}

	if len(accum) != size {
		t.Errorf("expect size of accum to be: %v, got: %v", size, len(accum))
	}

	expected := testingSlice(size, seed)
	for i := range expected {
		if accum[i] != expected[i] {
			t.Fatalf("expect accum[%d] to be: %d, got: %d", i, expected[i], accum[i])
		}
	}
}

func testingIter(size int, rng *rand.Rand) iter.Seq[int] {
	return func(yield func(int) bool) {
		for range size {
			if !yield(rng.Intn(size)) {
				return
			}
		}
	}
}
