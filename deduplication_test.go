package main

import (
	"fmt"
	"slices"
	"sort"
	"testing"
)

type deduplicateBench func(haystack []int) []int

func BenchmarkDeduplication(b *testing.B) {
	switch variant {
	case "slice":
		for _, size := range sizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				benchmarkDeduplicate(size, benchDeduplicateSlice),
			)
		}
	case "slice_sort_inplace":
		for _, size := range sizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				benchmarkDeduplicate(size, benchDeduplicateSliceSortInplace),
			)
		}
	case "map":
		for _, size := range sizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				benchmarkDeduplicate(size, benchDeduplicateMap),
			)
		}
	default:
		b.Errorf("speficy which test to run: -args -test slice|map|unique")
	}
}

func benchmarkDeduplicate(size int, runF deduplicateBench) func(*testing.B) {
	haystack := testingSlice(size)
	return func(b *testing.B) {
		for b.Loop() {
			// Note: we need to copy to prevent runs pre-sorting
			// arrays for each other.
			h := make([]int, len(haystack))
			copy(h, haystack)
			runF(h)
		}
	}
}

func benchDeduplicateSlice(haystack []int) []int {
	var result []int
	for _, val := range haystack {
		if !slices.Contains(result, val) {
			result = append(result, val)
		}
	}
	return result
}

func benchDeduplicateSliceSortInplace(haystack []int) []int {
	// "borrowed" from https://go.dev/wiki/SliceTricks#in-place-deduplicate-comparable, thanks!
	// Note sort + slices.Compact is the same thing.
	sort.Ints(haystack)
	j := 0
	for i := 1; i < len(haystack); i++ {
		if haystack[j] == haystack[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		haystack[j] = haystack[i]
	}
	return haystack[:j+1]
}

func benchDeduplicateMap(haystack []int) []int {
	result := make([]int, 0, len(haystack))
	seen := make(map[int]struct{}, len(haystack))

	for _, item := range haystack {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}
