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
	switch variant {
	case "plus":
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("plus_sign(%d) ops:(%d)", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, concatPlus),
				)
			}
		}
	case "sprintf":
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("sprintf(%d) ops:(%d)", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, concatSprintf),
				)
			}
		}
	case "join":
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("strings_join(%d) ops:(%d)", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, concatJoin),
				)
			}
		}
	case "builder":
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("strings_builder(%d) ops:(%d)", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, concatBuilder),
				)
			}
		}
	case "builder_pool":
		for _, stringSize := range sizes {
			for _, nOperations := range sizes {
				b.Run(
					fmt.Sprintf("strings_builder_pool(%d) ops:(%d)", stringSize, nOperations),
					benchmarkConcat(stringSize, nOperations, concatBuilderPool),
				)
			}
		}
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
