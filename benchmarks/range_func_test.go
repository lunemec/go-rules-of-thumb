package main

import (
	"fmt"
	"iter"
	"math/rand"
	"testing"
)

func BenchmarkIterate(b *testing.B) {
	switch variant {
	case "slice_iterate":
		for _, size := range sizes {
			for _, iterations := range sizes {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", size, iterations),
					benchmarkSliceIterate(size, iterations),
				)
			}
		}
	case "range_func":
		for _, size := range sizes {
			for _, iterations := range sizes {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", size, iterations),
					benchmarkRangeFuncIterate(size, iterations),
				)
			}
		}
	case "direct":
		for _, size := range sizes {
			for _, iterations := range sizes {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", size, iterations),
					benchmarkDirect(size, iterations),
				)
			}
		}
	default:
		b.Errorf("speficy which test to run: -args -test slice_iterate|range_func|direct")
	}
}

func benchmarkSliceIterate(size, iterations int) func(*testing.B) {
	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			for i := 0; i <= iterations; i++ {
				// Here we have to allocate the iteration slice
				// as that is the main benefit of range over func
				// - no upfront allocation.
				iterSlice := testingSlice(size)
				for _, val := range iterSlice {
					acc += val
				}
			}
		}
	}
}

func benchmarkDirect(size, iterations int) func(*testing.B) {
	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			for i := 0; i <= iterations; i++ {
				for range size {
					acc += rand.Intn(size)
				}
			}
		}
	}
}

func benchmarkRangeFuncIterate(size, iterations int) func(*testing.B) {
	return func(b *testing.B) {
		var acc int

		for b.Loop() {
			for i := 0; i <= iterations; i++ {
				for val := range testingIter(size) {
					acc += val
				}
			}
		}
	}
}

func TestIter(t *testing.T) {
	size := 10
	var accum []int

	for val := range testingIter(size) {
		accum = append(accum, val)
	}

	if len(accum) != size {
		t.Errorf("expect size of accum to be: %v, got: %v", size, len(accum))
	}
}

func testingSlice(size int) []int {
	// var ts = make([]int, size)
	ts := []int{}
	for range size {
		// ts[i] = rand.Intn(size)
		ts = append(ts, rand.Intn(size))
	}
	return ts
}

func testingIter(size int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for range size {
			if !yield(rand.Intn(size)) {
				return
			}
		}
	}
}
