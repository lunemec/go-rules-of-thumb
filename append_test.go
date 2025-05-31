package main

import (
	"fmt"
	"testing"
)

type appendBench func(first []int, second []int) []int

func BenchmarkAppend(b *testing.B) {
	for _, sizeFirst := range sizes {
		for _, sizeSecond := range sizes {
			b.Run(
				fmt.Sprintf("append_expand(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkAppend(sizeFirst, sizeSecond, appendExpand),
			)
			b.Run(
				fmt.Sprintf("append_for(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkAppend(sizeFirst, sizeSecond, appendFor),
			)
			b.Run(
				fmt.Sprintf("append_for_prealloc(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkAppend(sizeFirst, sizeSecond, appendForPrealloc),
			)
			b.Run(
				fmt.Sprintf("append_for_index(%d)(%d)", sizeFirst, sizeSecond),
				benchmarkAppend(sizeFirst, sizeSecond, appendForIdx),
			)
		}
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
	var out = make([]int, 0, len(first)+len(second))
	copy(out, first)
	for _, secondVal := range second {
		out = append(out, secondVal)
	}
	return out
}

func appendForIdx(first, second []int) []int {
	var out = make([]int, len(first)+len(second))
	copy(out, first)
	for i, secondVal := range second {
		out[(len(first))+i] = secondVal
	}
	return out
}
