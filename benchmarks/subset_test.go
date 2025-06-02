package main

import (
	"fmt"
	"slices"
	"sort"
	"testing"
)

type subsetBench func(first []int, second []int) bool

func BenchmarkSubset(b *testing.B) {
	runBenchmark := func(runF subsetBench) {
		for _, subsetSize := range sizes {
			for _, setSize := range sizes {
				if subsetSize > setSize {
					continue
				}
				b.Run(
					fmt.Sprintf("size=%d subset=%d", setSize, subsetSize),
					benchmarkSubset(subsetSize, setSize, runF),
				)
			}
		}
	}

	switch variant {
	case "slice":
		runBenchmark(subsetSlice)
	case "map":
		runBenchmark(subsetMap)
	case "slice_sort_binsearch":
		runBenchmark(subsetSortBinSearch)
	default:
		b.Errorf("speficy which test to run: -args -test map|slice|slice_sort_binsearch")
	}
}

func benchmarkSubset(sizeFirst, sizeSecond int, runF subsetBench) func(*testing.B) {
	// Slice A is smaller copy of Slice B, this is to force worst case for
	// for loop approach, so that it iterates all the values.
	second := testingSlice(sizeSecond)
	first := second[:sizeFirst]

	return func(b *testing.B) {
		for b.Loop() {
			runF(first, second)
		}
	}
}

func subsetSlice(first, second []int) bool {
	if len(first) > len(second) {
		return false
	}
	for _, firstValue := range first {
		var found bool
		for _, secondValue := range second {
			if firstValue == secondValue {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func subsetSortBinSearch(first, second []int) bool {
	if len(first) > len(second) {
		return false
	}
	// Need to copy because sort modifies the slice.
	// This adds time to the execution, but it is ok because
	// the other implementations don't have to do this.
	secondCopy := make([]int, len(second))
	copy(secondCopy, second)
	sort.Ints(secondCopy)

	for _, firstValue := range first {
		_, found := slices.BinarySearch(secondCopy, firstValue)
		if !found {
			return false
		}
	}
	return true
}

func subsetMap(first, second []int) bool {
	if len(first) > len(second) {
		return false
	}
	set := make(map[int]struct{}, len(second))
	for _, value := range second {
		set[value] = struct{}{}
	}

	for _, value := range first {
		if _, found := set[value]; !found {
			return false
		}
	}

	return true
}
