package main

import (
	"fmt"
	"slices"
	"sort"
	"testing"
)

type subsetBench func(first []int, second []int) bool

func BenchmarkSubset(b *testing.B) {
	for _, sizeFirst := range sizes {
		for _, sizeSecond := range sizes {
			if sizeFirst > sizeSecond {
				continue
			}
			b.Run(
				fmt.Sprintf("slice(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkSubset(sizeFirst, sizeSecond, subsetSlice),
			)
			b.Run(
				fmt.Sprintf("slice_sort_binsearch(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkSubset(sizeFirst, sizeSecond, subsetSortBinSearch),
			)
			b.Run(
				fmt.Sprintf("map(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkSubset(sizeFirst, sizeSecond, subsetMap),
			)
		}
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
	var secondCopy = make([]int, len(second))
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
