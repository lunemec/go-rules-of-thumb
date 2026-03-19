package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type concatBench func(first string, nops int) string

var concatSink string

var (
	concatStringSizes      = []int{10, 50, 100, 500, 1_000}
	concatOpCounts         = []int{10, 50, 100, 500, 1_000, 5_000}
	concatLargeStringSizes = []int{500, 1_000, 5_000}
	concatLargeOpCounts    = []int{100, 500, 1_000, 5_000}
)

func TestConcat(t *testing.T) {
	for _, stringSize := range sizesReduced {
		for _, nOperations := range sizesReduced {
			teststr := testingString(stringSize, stableSeed("TestConcat", stringSize, nOperations))

			out1 := concatPlus(teststr, nOperations)
			out2 := concatSprintf(teststr, nOperations)
			out3 := concatJoin(teststr, nOperations)
			out4 := concatBuilder(teststr, nOperations)
			out5 := concatBuilderPool(teststr, nOperations)

			require.Equal(t, out1, out2)
			require.Equal(t, out1, out3)
			require.Equal(t, out1, out4)
			require.Equal(t, out1, out5)
		}
	}
}

func BenchmarkConcat(b *testing.B) {
	runBenchmark := func(runF concatBench) {
		runConcatMatrix(b, "BenchmarkConcat", concatStringSizes, concatOpCounts, runF)
	}
	switch variant {
	case "plus":
		runBenchmark(concatPlus)
	case "sprintf":
		runBenchmark(concatSprintf)
	case "join":
		runBenchmark(concatJoin)
	case "builder":
		runBenchmark(concatBuilder)
	case "builder_pool":
		runBenchmark(concatBuilderPool)
	default:
		b.Errorf("specify which benchmark to run: -args -variant plus|sprintf|join|builder|builder_pool")
	}
}

func BenchmarkConcatLarge(b *testing.B) {
	runBenchmark := func(runF concatBench) {
		runConcatMatrix(b, "BenchmarkConcatLarge", concatLargeStringSizes, concatLargeOpCounts, runF)
	}
	switch variant {
	case "plus":
		runBenchmark(concatPlus)
	case "sprintf":
		runBenchmark(concatSprintf)
	case "join":
		runBenchmark(concatJoin)
	case "builder":
		runBenchmark(concatBuilder)
	case "builder_pool":
		runBenchmark(concatBuilderPool)
	default:
		b.Errorf("specify which benchmark to run: -args -variant plus|sprintf|join|builder|builder_pool")
	}
}

func runConcatMatrix(b *testing.B, benchName string, stringSizes, opCounts []int, runF concatBench) {
	for _, stringSize := range stringSizes {
		for _, nOperations := range opCounts {
			b.Run(
				fmt.Sprintf("size=%d iterations=%d", stringSize, nOperations),
				benchmarkConcat(benchName, stringSize, nOperations, runF),
			)
		}
	}
}

func benchmarkConcat(benchName string, strSize, nOps int, runF concatBench) func(*testing.B) {
	teststr := testingString(strSize, stableSeed(benchName, strSize, nOps))

	return func(b *testing.B) {
		for b.Loop() {
			concatSink = runF(teststr, nOps)
		}
	}
}

func testingString(size int, seed uint64) string {
	rng := newDeterministicRand(seed)
	var builder strings.Builder
	for range size {
		builder.WriteString(strconv.Itoa(rng.Intn(size)))
	}
	return builder.String()
}

func concatPlus(teststr string, nOps int) string {
	var out string
	for range nOps {
		out += teststr
		out += "..."
	}
	return out
}

func concatSprintf(teststr string, nOps int) string {
	var out string
	for range nOps {
		out = fmt.Sprintf("%s%s%s", out, teststr, "...")
	}
	return out
}

func concatJoin(teststr string, nOps int) string {
	var out string
	for range nOps {
		out = strings.Join([]string{out, teststr, "..."}, "")
	}
	return out
}

func concatBuilder(teststr string, nOps int) string {
	var builder strings.Builder
	builder.Grow(len(teststr))
	for range nOps {
		builder.WriteString(teststr)
		builder.WriteString("...")
	}
	return builder.String()
}

var builderPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

func concatBuilderPool(teststr string, nOps int) string {
	builder := builderPool.Get().(*strings.Builder)
	// On current Go releases, Reset discards the backing buffer, so this pool
	// measures sync.Pool coordination overhead rather than reusable capacity.
	builder.Reset()
	builder.Grow(len(teststr))
	defer builderPool.Put(builder)
	for range nOps {
		builder.WriteString(teststr)
		builder.WriteString("...")
	}
	return builder.String()
}
