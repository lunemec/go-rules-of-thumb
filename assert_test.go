package main

import (
	"testing"
)

func assert(truth bool, msg string) {
	if !truth {
		panic(msg)
	}
}

func BenchmarkAssertion(b *testing.B) {
	b.Run(
		"no assert",
		benchmarkNoAssert,
	)
	b.Run(
		"assert",
		benchmarkAssert,
	)
	b.Run(
		"assert(5)",
		benchmarkAssert5,
	)
	b.Run(
		"defer assert",
		benchmarkDeferAssert,
	)
}

var (
	Truth  = true
	Truth2 = true
	Truth3 = true
	Truth4 = true
	Truth5 = true
)

func benchmarkNoAssert(b *testing.B) {
	for b.Loop() {
		func() {
		}()
	}
}

func benchmarkAssert(b *testing.B) {
	for b.Loop() {
		func() {
			assert(Truth, "n must be larger than 0")
		}()
	}
}

func benchmarkAssert5(b *testing.B) {
	for b.Loop() {
		func() {
			assert(Truth, "n must be larger than 0")
			assert(Truth2, "n must be larger than 0")
			assert(Truth3, "n must be larger than 0")
			assert(Truth4, "n must be larger than 0")
			assert(Truth5, "n must be larger than 0")
		}()
	}
}

func benchmarkDeferAssert(b *testing.B) {
	for b.Loop() {
		func() {
			defer assert(Truth, "n must be larger than 0")
		}()
	}
}
