package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type appendBench func(first []int, second []int) []int

func TestAppend(t *testing.T) {
	for _, subsetSize := range sizes {
		for _, setSize := range sizes {
			set := testingSlice(setSize)
			subset := testingSlice(subsetSize)

			out1 := appendExpand(copySlice(set), copySlice(subset))
			out2 := appendFor(copySlice(set), copySlice(subset))
			out3 := appendForPrealloc(copySlice(set), copySlice(subset))
			out4 := appendForIdx(copySlice(set), copySlice(subset))

			require.Equal(t, out1, out2)
			require.Equal(t, out1, out3)
			require.Equal(t, out1, out4)
		}
	}
}

func BenchmarkAppend(b *testing.B) {
	runBenchmark := func(runF appendBench) {
		for _, subsetSize := range sizes {
			for _, setSize := range sizes {
				b.Run(
					fmt.Sprintf("size=%d subset=%d", setSize, subsetSize),
					benchmarkAppend(subsetSize, setSize, runF),
				)
			}
		}
	}
	switch variant {
	case "expand":
		runBenchmark(appendExpand)
	case "for":
		runBenchmark(appendFor)
	case "for_prealloc":
		runBenchmark(appendForPrealloc)
	case "for_index":
		runBenchmark(appendForIdx)
	default:
		b.Errorf("speficy which test to run: -args -test expand|for|for_prealloc|for_index")
	}
}

func benchmarkAppend(sizeFirst, sizeSecond int, runF appendBench) func(*testing.B) {
	first := testingSlice(sizeFirst)
	second := testingSlice(sizeSecond)

	return func(b *testing.B) {
		for b.Loop() {
			runF(first, second)
		}
	}
}

func appendExpand(first, second []int) []int {
	return append(first, second...)
}

func appendFor(first, second []int) []int {
	for _, secondVal := range second {
		first = append(first, secondVal)
	}
	return first
}

func appendForPrealloc(first, second []int) []int {
	out := make([]int, len(first), len(first)+len(second))
	copy(out, first)
	for _, secondVal := range second {
		out = append(out, secondVal)
	}
	return out
}

func appendForIdx(first, second []int) []int {
	out := make([]int, len(first)+len(second))
	copy(out, first)
	for i, secondVal := range second {
		out[(len(first))+i] = secondVal
	}
	return out
}
