package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type subsetBench func(first []int, second []int) bool

func TestSubset(t *testing.T) {
	for _, subsetSize := range sizesReduced {
		for _, setSize := range sizesReduced {
			set, subset := testingSets(setSize, subsetSize)

			out1 := subsetSlice(subset, set)
			out2 := subsetMap(subset, set)
			out3 := subsetSortBinSearch(subset, set)

			require.Equal(t, out1, out2)
			require.Equal(t, out1, out3)
		}
	}
}

func BenchmarkSubset(b *testing.B) {
	runBenchmark := func(runF subsetBench) {
		for _, subsetSize := range sizes {
			for _, setSize := range sizes {
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

func testingSets(setSize, subsetSize int) ([]int, []int) {
	subset := make([]int, subsetSize)
	set := testingSlice(setSize)

	// In 50% of cases we make completely random subset.
	if rand.Intn(2) == 0 {
		subset = testingSlice(subsetSize)
	} else {
		if setSize < subsetSize {
			// In the case where subset is larger than set, we can't
			// just take slice of it, so we will copy it multiple times,
			// We want this branch of the random flip to be when the subset is present.
			for i := 0; i < subsetSize; i += setSize {
				copy(subset[i:], set)
			}
		} else {
			// Take end of the bigger set to force
			// us to iterate further (worst case).
			subset = set[setSize-subsetSize:]
		}
	}

	return set, subset
}

func benchmarkSubset(setSize, subsetSize int, runF subsetBench) func(*testing.B) {
	set, subset := testingSets(setSize, subsetSize)
	return func(b *testing.B) {
		for b.Loop() {
			runF(subset, set)
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
