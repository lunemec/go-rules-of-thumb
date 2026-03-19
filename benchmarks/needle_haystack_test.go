package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type needleInHaystackBench func(checks int, needle int, haystack []int) bool

var needleInAHaystackSink bool

func TestNeedleInAHaystack(t *testing.T) {
	for _, size := range sizesReduced {
		haystack := testingSlice(size, stableSeed("TestNeedleInAHaystack", size, "haystack"))
		needles := testingNeedles(size, 8, stableSeed("TestNeedleInAHaystack", size, "needles"))

		for _, needle := range needles {
			out1 := needleInAHaystackSlice(needle, haystack)
			out2 := needleInAHaystackMap(needle, haystackToMap(haystack))

			require.Equal(t, out1, out2)
		}
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
	haystack := testingSlice(size, stableSeed("BenchmarkNeedleInAHaystack", size, checks, "haystack"))
	needles := testingNeedles(size, 256, stableSeed("BenchmarkNeedleInAHaystack", size, checks, "needles"))

	return func(b *testing.B) {
		needleIndex := 0
		for b.Loop() {
			needleInAHaystackSink = runF(checks, needles[needleIndex], haystack)
			needleIndex++
			if needleIndex == len(needles) {
				needleIndex = 0
			}
		}
	}
}

func testingNeedles(size int, count int, seed uint64) []int {
	rng := newDeterministicRand(seed)
	needles := make([]int, 0, count)
	for range count {
		needles = append(needles, rng.Intn(size*2))
	}
	return needles
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
