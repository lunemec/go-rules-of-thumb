package main

import (
	"fmt"
	"testing"
)

type appendBench func(first []int, second []int) []int

func BenchmarkAppend(b *testing.B) {
	switch variant {
	case "expand":
		for _, sizeFirst := range sizes {
			for _, sizeSecond := range sizes {
				b.Run(
					fmt.Sprintf("size=%d subset_size=%d", sizeFirst, sizeSecond),
					benchmarkAppend(sizeFirst, sizeSecond, appendExpand),
				)
			}
		}
	case "for":
		for _, sizeFirst := range sizes {
			for _, sizeSecond := range sizes {
				b.Run(
					fmt.Sprintf("size=%d subset_size=%d", sizeFirst, sizeSecond),
					benchmarkAppend(sizeFirst, sizeSecond, appendFor),
				)
			}
		}
	case "for_prealloc":
		for _, sizeFirst := range sizes {
			for _, sizeSecond := range sizes {
				b.Run(
					fmt.Sprintf("size=%d subset_size=%d", sizeFirst, sizeSecond),
					benchmarkAppend(sizeFirst, sizeSecond, appendForPrealloc),
				)
			}
		}
	case "for_index":
		for _, sizeFirst := range sizes {
			for _, sizeSecond := range sizes {
				b.Run(
					fmt.Sprintf("size=%d subset_size=%d", sizeFirst, sizeSecond),
					benchmarkAppend(sizeFirst, sizeSecond, appendForIdx),
				)
			}
		}
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
	out := make([]int, 0, len(first)+len(second))
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
