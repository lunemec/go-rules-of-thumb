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
	for _, stringSize := range sizes {
		for _, nOperations := range sizes {
			b.Run(
				fmt.Sprintf("plus_sign(%d) ops:(%d)", stringSize, nOperations),
				benchmarkConcat(stringSize, nOperations, concatPlus),
			)
			b.Run(
				fmt.Sprintf("sprintf(%d) ops:(%d)", stringSize, nOperations),
				benchmarkConcat(stringSize, nOperations, concatSprintf),
			)
			b.Run(
				fmt.Sprintf("strings_join(%d) ops:(%d)", stringSize, nOperations),
				benchmarkConcat(stringSize, nOperations, concatJoin),
			)
			b.Run(
				fmt.Sprintf("strings_builder(%d) ops:(%d)", stringSize, nOperations),
				benchmarkConcat(stringSize, nOperations, concatBuilder),
			)
			b.Run(
				fmt.Sprintf("strings_builder_pool(%d) ops:(%d)", stringSize, nOperations),
				benchmarkConcat(stringSize, nOperations, concatBuilderPool),
			)
		}
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
	for i := 0; i < size; i++ {
		builder.WriteString(strconv.Itoa(rand.Intn(size)))
	}
	return builder.String()
}

func concatPlus(teststr string, nOps int) string {
	for i := 0; i < nOps; i++ {
		teststr = teststr + "..."
	}
	return teststr
}

func concatSprintf(teststr string, nOps int) string {
	for i := 0; i < nOps; i++ {
		teststr = fmt.Sprintf("%s%s", teststr, "...")
	}
	return teststr
}

func concatJoin(teststr string, nOps int) string {
	for i := 0; i < nOps; i++ {
		teststr = strings.Join([]string{teststr, "..."}, "")
	}
	return teststr
}

func concatBuilder(teststr string, nOps int) string {
	var builder strings.Builder
	builder.Grow(len(teststr))
	for i := 0; i < nOps; i++ {
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
	for i := 0; i < nOps; i++ {
		builder.WriteString(teststr)
		builder.WriteString("...")
	}
	return builder.String()
}
