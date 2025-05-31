package main

import (
	"fmt"
	"math/rand"
	"testing"
)

type needleInHaystackBench func(checks int, needle int, haystack []int) bool

func BenchmarkNeedleInAHaystack(b *testing.B) {
	for _, size := range sizes {
		for _, nchecks := range needles {
			b.Run(
				fmt.Sprintf("slice(%d) needles(%d)", size, nchecks),
				benchmarkNeedleInAHaystack(size, nchecks, benchNeedleInAHaystackSlice),
			)
			b.Run(
				fmt.Sprintf("map(%d) needles(%d)", size, nchecks),
				benchmarkNeedleInAHaystack(size, nchecks, benchNeedleInAHaystackMap))
		}
	}
}

func benchmarkNeedleInAHaystack(size int, checks int, runF needleInHaystackBench) func(*testing.B) {
	haystack := testingSlice(size)

	return func(b *testing.B) {
		for b.Loop() {
			needle := rand.Intn(size)
			runF(checks, needle, haystack)
		}
	}
}

func benchNeedleInAHaystackSlice(checks int, needle int, haystack []int) bool {
	var f bool
	for i := 0; i < checks; i++ {
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
	for i := 0; i < checks; i++ {
		f = needleInAHaystackMap(needle, mapHaystack)
	}
	return f
}

func haystackToMap(haystack []int) map[int]struct{} {
	var out = make(map[int]struct{}, len(haystack))
	for _, v := range haystack {
		out[v] = struct{}{}
	}
	return out
}

func needleInAHaystackMap(needle int, haystack map[int]struct{}) bool {
	_, ok := haystack[needle]
	return ok
}
