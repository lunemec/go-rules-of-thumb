package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type concatBench func(first string, nops int) string

func TestConcat(t *testing.T) {
	for _, stringSize := range sizesReduced {
		for _, nOperations := range sizesReduced {
			teststr := testingString(stringSize)

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
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("size=%d iterations=%d", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, runF),
				)
			}
		}
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
		b.Errorf("speficy which test to run: -args -test plus|sprintf|join|builder|builder_pool")
	}
}

func benchmarkConcat(strSize, nOps int, runF concatBench) func(*testing.B) {
	teststr := testingString(strSize)

	return func(b *testing.B) {
		for b.Loop() {
			runF(teststr, nOps)
		}
	}
}

func testingString(size int) string {
	var builder strings.Builder
	for range size {
		builder.WriteString(strconv.Itoa(rand.Intn(size)))
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
	builder.Reset()
	builder.Grow(len(teststr))
	defer builderPool.Put(builder)
	for range nOps {
		builder.WriteString(teststr)
		builder.WriteString("...")
	}
	return builder.String()
}
