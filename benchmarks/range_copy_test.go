package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

var rangeCopyLengths = []int{10, 100, 1_000, 10_000, 100_000}

type rangeCopyCore struct {
	ID     uint64
	Score  uint32
	Flags  uint16
	Active bool
	_      [1]byte
}

type rangeCopyRecord16 struct {
	Core rangeCopyCore
}

type rangeCopyRecord32 struct {
	Core rangeCopyCore
	Pad  [16]byte
}

type rangeCopyRecord64 struct {
	Core rangeCopyCore
	Pad  [48]byte
}

type rangeCopyRecord128 struct {
	Core rangeCopyCore
	Pad  [112]byte
}

type rangeCopyRecord256 struct {
	Core rangeCopyCore
	Pad  [240]byte
}

type rangeCopyRecord512 struct {
	Core rangeCopyCore
	Pad  [496]byte
}

var rangeCopySink uint64

func TestRangeCopyOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(16), unsafe.Sizeof(rangeCopyRecord16{}))
	require.Equal(t, uintptr(32), unsafe.Sizeof(rangeCopyRecord32{}))
	require.Equal(t, uintptr(64), unsafe.Sizeof(rangeCopyRecord64{}))
	require.Equal(t, uintptr(128), unsafe.Sizeof(rangeCopyRecord128{}))
	require.Equal(t, uintptr(256), unsafe.Sizeof(rangeCopyRecord256{}))
	require.Equal(t, uintptr(512), unsafe.Sizeof(rangeCopyRecord512{}))

	for _, length := range sizesReduced {
		t.Run(fmt.Sprintf("length=%d", length), func(t *testing.T) {
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord16, scoreRangeCopyRecord16Value, scoreRangeCopyRecord16Ptr)
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord32, scoreRangeCopyRecord32Value, scoreRangeCopyRecord32Ptr)
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord64, scoreRangeCopyRecord64Value, scoreRangeCopyRecord64Ptr)
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord128, scoreRangeCopyRecord128Value, scoreRangeCopyRecord128Ptr)
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord256, scoreRangeCopyRecord256Value, scoreRangeCopyRecord256Ptr)
			requireRangeCopyEquivalent(t, length, makeRangeCopyRecord512, scoreRangeCopyRecord512Value, scoreRangeCopyRecord512Ptr)
		})
	}
}

func requireRangeCopyEquivalent[T any](t *testing.T, length int, build func(int) T, scoreValue func(T) uint64, scorePtr func(*T) uint64) {
	t.Helper()

	values := buildRangeCopyValues(length, build)
	pointers := buildRangeCopyPointers(length, build)

	require.Equal(t, scanRangeValues(values, scoreValue), scanIndexValues(values, scorePtr))
	require.Equal(t, scanRangeValues(values, scoreValue), scanRangePointers(pointers, scorePtr))
}

func BenchmarkRangeCopy(b *testing.B) {
	cases := []struct {
		size         int
		rangeBench   func(int) func(*testing.B)
		indexBench   func(int) func(*testing.B)
		pointerBench func(int) func(*testing.B)
	}{
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord16{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord16, scoreRangeCopyRecord16Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord16, scoreRangeCopyRecord16Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord16, scoreRangeCopyRecord16Ptr),
		},
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord32{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord32, scoreRangeCopyRecord32Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord32, scoreRangeCopyRecord32Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord32, scoreRangeCopyRecord32Ptr),
		},
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord64{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord64, scoreRangeCopyRecord64Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord64, scoreRangeCopyRecord64Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord64, scoreRangeCopyRecord64Ptr),
		},
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord128{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord128, scoreRangeCopyRecord128Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord128, scoreRangeCopyRecord128Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord128, scoreRangeCopyRecord128Ptr),
		},
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord256{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord256, scoreRangeCopyRecord256Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord256, scoreRangeCopyRecord256Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord256, scoreRangeCopyRecord256Ptr),
		},
		{
			size:         int(unsafe.Sizeof(rangeCopyRecord512{})),
			rangeBench:   benchmarkRangeValues(makeRangeCopyRecord512, scoreRangeCopyRecord512Value),
			indexBench:   benchmarkIndexValues(makeRangeCopyRecord512, scoreRangeCopyRecord512Ptr),
			pointerBench: benchmarkRangePointers(makeRangeCopyRecord512, scoreRangeCopyRecord512Ptr),
		},
	}

	for _, tc := range cases {
		run := tc.rangeBench
		switch variant {
		case "range_value":
		case "index_value":
			run = tc.indexBench
		case "range_pointer":
			run = tc.pointerBench
		default:
			b.Fatalf("specify which benchmark to run: -args -variant range_value|index_value|range_pointer")
		}

		for _, length := range rangeCopyLengths {
			b.Run(fmt.Sprintf("size=%d length=%d", tc.size, length), run(length))
		}
	}
}

func benchmarkRangeValues[T any](build func(int) T, score func(T) uint64) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		values := buildRangeCopyValues(length, build)
		return func(b *testing.B) {
			for b.Loop() {
				rangeCopySink = scanRangeValues(values, score)
			}
		}
	}
}

func benchmarkIndexValues[T any](build func(int) T, score func(*T) uint64) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		values := buildRangeCopyValues(length, build)
		return func(b *testing.B) {
			for b.Loop() {
				rangeCopySink = scanIndexValues(values, score)
			}
		}
	}
}

func benchmarkRangePointers[T any](build func(int) T, score func(*T) uint64) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		pointers := buildRangeCopyPointers(length, build)
		return func(b *testing.B) {
			for b.Loop() {
				rangeCopySink = scanRangePointers(pointers, score)
			}
		}
	}
}

func buildRangeCopyValues[T any](length int, build func(int) T) []T {
	values := make([]T, length)
	for i := range values {
		values[i] = build(i)
	}
	return values
}

func buildRangeCopyPointers[T any](length int, build func(int) T) []*T {
	backing := make([]T, length)
	pointers := make([]*T, length)
	for i := range backing {
		backing[i] = build(i)
		pointers[i] = &backing[i]
	}
	return pointers
}

func scanRangeValues[T any](records []T, score func(T) uint64) uint64 {
	var acc uint64
	for _, record := range records {
		acc += score(record)
	}
	return acc
}

func scanIndexValues[T any](records []T, score func(*T) uint64) uint64 {
	var acc uint64
	for i := range records {
		acc += score(&records[i])
	}
	return acc
}

func scanRangePointers[T any](records []*T, score func(*T) uint64) uint64 {
	var acc uint64
	for _, record := range records {
		acc += score(record)
	}
	return acc
}

func rangeCopyCoreScore(core rangeCopyCore) uint64 {
	if !core.Active {
		return 0
	}
	return core.ID + uint64(core.Score) + uint64(core.Flags)
}

func buildRangeCopyCore(seed int) rangeCopyCore {
	return rangeCopyCore{
		ID:     uint64(seed+1) * 17,
		Score:  uint32((seed%97)+1) * 13,
		Flags:  uint16(seed % 29),
		Active: seed%4 != 0,
	}
}

func scoreRangeCopyRecord16Value(record rangeCopyRecord16) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord32Value(record rangeCopyRecord32) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord64Value(record rangeCopyRecord64) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord128Value(record rangeCopyRecord128) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord256Value(record rangeCopyRecord256) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord512Value(record rangeCopyRecord512) uint64 {
	return rangeCopyCoreScore(record.Core)
}

func scoreRangeCopyRecord16Ptr(record *rangeCopyRecord16) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord32Ptr(record *rangeCopyRecord32) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord64Ptr(record *rangeCopyRecord64) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord128Ptr(record *rangeCopyRecord128) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord256Ptr(record *rangeCopyRecord256) uint64 {
	return rangeCopyCoreScore(record.Core)
}
func scoreRangeCopyRecord512Ptr(record *rangeCopyRecord512) uint64 {
	return rangeCopyCoreScore(record.Core)
}

func makeRangeCopyRecord16(seed int) rangeCopyRecord16 {
	return rangeCopyRecord16{Core: buildRangeCopyCore(seed)}
}

func makeRangeCopyRecord32(seed int) rangeCopyRecord32 {
	record := rangeCopyRecord32{Core: buildRangeCopyCore(seed)}
	fillPattern(record.Pad[:], byte(seed+11))
	return record
}

func makeRangeCopyRecord64(seed int) rangeCopyRecord64 {
	record := rangeCopyRecord64{Core: buildRangeCopyCore(seed)}
	fillPattern(record.Pad[:], byte(seed+29))
	return record
}

func makeRangeCopyRecord128(seed int) rangeCopyRecord128 {
	record := rangeCopyRecord128{Core: buildRangeCopyCore(seed)}
	fillPattern(record.Pad[:], byte(seed+47))
	return record
}

func makeRangeCopyRecord256(seed int) rangeCopyRecord256 {
	record := rangeCopyRecord256{Core: buildRangeCopyCore(seed)}
	fillPattern(record.Pad[:], byte(seed+73))
	return record
}

func makeRangeCopyRecord512(seed int) rangeCopyRecord512 {
	record := rangeCopyRecord512{Core: buildRangeCopyCore(seed)}
	fillPattern(record.Pad[:], byte(seed+101))
	return record
}
