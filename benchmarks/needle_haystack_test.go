package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

type needleInHaystackBench func(checks int, needle int, haystack []int) bool

func TestNeedleInAHaystack(t *testing.T) {
	for _, size := range sizes {
		haystack := testingSlice(size)
		needle := rand.Intn(size * 2)

		out1 := needleInAHaystackSlice(needle, haystack)
		out2 := needleInAHaystackMap(needle, haystackToMap(haystack))

		require.Equal(t, out1, out2)
	}
}

func BenchmarkNeedleInAHaystack(b *testing.B) {
	runBenchmark := func(runF needleInHaystackBench) {
		for _, size := range sizes {
			for _, nchecks := range needles {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", size, nchecks),
					benchmarkNeedleInAHaystack(size, nchecks, runF),
				)
			}
		}
	}
	switch variant {
	case "slice":
		runBenchmark(benchNeedleInAHaystackSlice)
	case "map":
		runBenchmark(benchNeedleInAHaystackMap)
	default:
		b.Errorf("speficy which test to run: -args -test slice|map")
	}
}

func benchmarkNeedleInAHaystack(size int, checks int, runF needleInHaystackBench) func(*testing.B) {
	haystack := testingSlice(size)

	return func(b *testing.B) {
		for b.Loop() {
			// We make our needle to have 50% chance
			// to not be in the haystack.
			needle := rand.Intn(size * 2)
			runF(checks, needle, haystack)
		}
	}
}

func benchNeedleInAHaystackSlice(checks int, needle int, haystack []int) bool {
	var f bool
	for range checks {
		f = needleInAHaystackSlice(needle, haystack)
	}
	return f
}

func needleInAHaystackSlice(needle int, haystack []int) bool {
	// This is identical to `slices.Contains` implementation.
	for i := range haystack {
		if needle == haystack[i] {
			return true
		}
	}
	return false
}

func benchNeedleInAHaystackMap(checks int, needle int, haystack []int) bool {
	var f bool
	mapHaystack := haystackToMap(haystack)
	for range checks {
		f = needleInAHaystackMap(needle, mapHaystack)
	}
	return f
}

func haystackToMap(haystack []int) map[int]struct{} {
	out := make(map[int]struct{}, len(haystack))
	for _, v := range haystack {
		out[v] = struct{}{}
	}
	return out
}

func needleInAHaystackMap(needle int, haystack map[int]struct{}) bool {
	_, ok := haystack[needle]
	return ok
}
