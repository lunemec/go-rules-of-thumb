package main

import (
	"fmt"
	"testing"
)

type assertBench func()

func _assert(truth bool, msg string) {
	if !truth {
		panic(msg)
	}
}

func BenchmarkAssert(b *testing.B) {
	runBenchmark := func(runF assertBench) {
		for _, size := range sizesReduced {
			b.Run(
				fmt.Sprintf("size=%d", size),
				benchmarkAssert(size, runF),
			)
		}
	}
	switch variant {
	case "no_assert":
		runBenchmark(noAssert)
	case "assert":
		runBenchmark(assert)
	case "defer_assert":
		runBenchmark(deferAssert)
	default:
		b.Errorf("speficy which test to run: -args -test no_assert|assert|defer_assert")
	}
}

var Truth = true

func benchmarkAssert(size int, runF assertBench) func(*testing.B) {
	return func(b *testing.B) {
		for b.Loop() {
			for range size {
				runF()
			}
		}
	}
}

func noAssert() {
}

func assert() {
	_assert(Truth, "n must be larger than 0")
}

func deferAssert() {
	defer _assert(Truth, "n must be larger than 0")
}
