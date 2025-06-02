package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type concatBench func(first string, nops int) string

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
	for range nOps {
		teststr = teststr + "..."
	}
	return teststr
}

func concatSprintf(teststr string, nOps int) string {
	for range nOps {
		teststr = fmt.Sprintf("%s%s", teststr, "...")
	}
	return teststr
}

func concatJoin(teststr string, nOps int) string {
	for range nOps {
		teststr = strings.Join([]string{teststr, "..."}, "")
	}
	return teststr
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
