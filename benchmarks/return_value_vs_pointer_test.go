package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

const (
	returnBenchmarkRecords = 256
	returnTestRecords      = 32
)

type ownedResult8 = scoringRequest8
type ownedResult16 = scoringRequest16
type ownedResult24 = scoringRequest24
type ownedResult32 = scoringRequest32
type ownedResult64 = scoringRequest64
type ownedResult96 = scoringRequest96
type ownedResult128 = scoringRequest128
type ownedResult192 = scoringRequest192
type ownedResult256 = scoringRequest256
type ownedResult384 = scoringRequest384
type ownedResult512 = scoringRequest512

var returnValueVsPointerSink uint64

func TestReturnValueVsPointerOutputsMatch(t *testing.T) {
	t.Run("payload_sizes", func(t *testing.T) {
		require.Equal(t, uintptr(8), unsafe.Sizeof(ownedResult8{}))
		require.Equal(t, uintptr(16), unsafe.Sizeof(ownedResult16{}))
		require.Equal(t, uintptr(24), unsafe.Sizeof(ownedResult24{}))
		require.Equal(t, uintptr(32), unsafe.Sizeof(ownedResult32{}))
		require.Equal(t, uintptr(64), unsafe.Sizeof(ownedResult64{}))
		require.Equal(t, uintptr(96), unsafe.Sizeof(ownedResult96{}))
		require.Equal(t, uintptr(128), unsafe.Sizeof(ownedResult128{}))
		require.Equal(t, uintptr(192), unsafe.Sizeof(ownedResult192{}))
		require.Equal(t, uintptr(256), unsafe.Sizeof(ownedResult256{}))
		require.Equal(t, uintptr(384), unsafe.Sizeof(ownedResult384{}))
		require.Equal(t, uintptr(512), unsafe.Sizeof(ownedResult512{}))
	})

	t.Run("size=8", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult8Value, buildOwnedResult8Pointer)
	})
	t.Run("size=16", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult16Value, buildOwnedResult16Pointer)
	})
	t.Run("size=24", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult24Value, buildOwnedResult24Pointer)
	})
	t.Run("size=32", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult32Value, buildOwnedResult32Pointer)
	})
	t.Run("size=64", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult64Value, buildOwnedResult64Pointer)
	})
	t.Run("size=96", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult96Value, buildOwnedResult96Pointer)
	})
	t.Run("size=128", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult128Value, buildOwnedResult128Pointer)
	})
	t.Run("size=192", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult192Value, buildOwnedResult192Pointer)
	})
	t.Run("size=256", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult256Value, buildOwnedResult256Pointer)
	})
	t.Run("size=384", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult384Value, buildOwnedResult384Pointer)
	})
	t.Run("size=512", func(t *testing.T) {
		requireOwnedResultsEquivalent(t, buildOwnedResult512Value, buildOwnedResult512Pointer)
	})
}

func requireOwnedResultsEquivalent[T comparable](t *testing.T, buildValue func(int) T, buildPointer func(int) *T) {
	t.Helper()

	for i := 0; i < returnTestRecords; i++ {
		expected := buildValue(i)
		actual := buildPointer(i)
		require.NotNil(t, actual)
		require.Equal(t, expected, *actual)
	}
}

func BenchmarkReturnValueVsPointer(b *testing.B) {
	cases := []struct {
		size         int
		valueBench   func(*testing.B)
		pointerBench func(*testing.B)
	}{
		{size: int(unsafe.Sizeof(ownedResult8{})), valueBench: benchmarkOwnedResults(sumOwnedResults8Value), pointerBench: benchmarkOwnedResults(sumOwnedResults8Pointer)},
		{size: int(unsafe.Sizeof(ownedResult16{})), valueBench: benchmarkOwnedResults(sumOwnedResults16Value), pointerBench: benchmarkOwnedResults(sumOwnedResults16Pointer)},
		{size: int(unsafe.Sizeof(ownedResult24{})), valueBench: benchmarkOwnedResults(sumOwnedResults24Value), pointerBench: benchmarkOwnedResults(sumOwnedResults24Pointer)},
		{size: int(unsafe.Sizeof(ownedResult32{})), valueBench: benchmarkOwnedResults(sumOwnedResults32Value), pointerBench: benchmarkOwnedResults(sumOwnedResults32Pointer)},
		{size: int(unsafe.Sizeof(ownedResult64{})), valueBench: benchmarkOwnedResults(sumOwnedResults64Value), pointerBench: benchmarkOwnedResults(sumOwnedResults64Pointer)},
		{size: int(unsafe.Sizeof(ownedResult96{})), valueBench: benchmarkOwnedResults(sumOwnedResults96Value), pointerBench: benchmarkOwnedResults(sumOwnedResults96Pointer)},
		{size: int(unsafe.Sizeof(ownedResult128{})), valueBench: benchmarkOwnedResults(sumOwnedResults128Value), pointerBench: benchmarkOwnedResults(sumOwnedResults128Pointer)},
		{size: int(unsafe.Sizeof(ownedResult192{})), valueBench: benchmarkOwnedResults(sumOwnedResults192Value), pointerBench: benchmarkOwnedResults(sumOwnedResults192Pointer)},
		{size: int(unsafe.Sizeof(ownedResult256{})), valueBench: benchmarkOwnedResults(sumOwnedResults256Value), pointerBench: benchmarkOwnedResults(sumOwnedResults256Pointer)},
		{size: int(unsafe.Sizeof(ownedResult384{})), valueBench: benchmarkOwnedResults(sumOwnedResults384Value), pointerBench: benchmarkOwnedResults(sumOwnedResults384Pointer)},
		{size: int(unsafe.Sizeof(ownedResult512{})), valueBench: benchmarkOwnedResults(sumOwnedResults512Value), pointerBench: benchmarkOwnedResults(sumOwnedResults512Pointer)},
	}

	for _, tc := range cases {
		run := tc.valueBench
		switch variant {
		case "return_value":
		case "return_pointer":
			run = tc.pointerBench
		default:
			b.Fatalf("specify which benchmark to run: -args -variant return_value|return_pointer")
		}

		b.Run(fmt.Sprintf("size=%d", tc.size), run)
	}
}

func benchmarkOwnedResults(sum func([]int) uint64) func(*testing.B) {
	inputs := buildOwnedResultInputs(returnBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			returnValueVsPointerSink = sum(inputs)
		}
	}
}

func buildOwnedResultInputs(size int) []int {
	inputs := make([]int, size)
	for i := range inputs {
		inputs[i] = i
	}
	return inputs
}

func sumOwnedResults8Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult8Value(seed)
		acc += result.Core.AccountID
	}
	return acc
}

func sumOwnedResults8Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult8Pointer(seed)
		acc += result.Core.AccountID
	}
	return acc
}

func sumOwnedResults16Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult16Value(seed)
		acc += result.Core.AccountID
		acc += uint64(result.Core.Score)
		acc += uint64(result.Core.Attempts)
		acc += uint64(result.Core.Flags)
	}
	return acc
}

func sumOwnedResults16Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult16Pointer(seed)
		acc += result.Core.AccountID
		acc += uint64(result.Core.Score)
		acc += uint64(result.Core.Attempts)
		acc += uint64(result.Core.Flags)
	}
	return acc
}

func sumOwnedResults24Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult24Value(seed)
		acc += result.Core.AccountID ^ result.Core.RequestID
		acc += uint64(result.Core.Score)
		acc += uint64(result.Core.Attempts)
	}
	return acc
}

func sumOwnedResults24Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult24Pointer(seed)
		acc += result.Core.AccountID ^ result.Core.RequestID
		acc += uint64(result.Core.Score)
		acc += uint64(result.Core.Attempts)
	}
	return acc
}

func sumOwnedResults32Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult32Value(seed)
		acc += scoreHotFields(result.Core, requestLatency32{})
	}
	return acc
}

func sumOwnedResults32Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult32Pointer(seed)
		acc += scoreHotFields(result.Core, requestLatency32{})
	}
	return acc
}

func sumOwnedResults64Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult64Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults64Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult64Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults96Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult96Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults96Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult96Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults128Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult128Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults128Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult128Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults192Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult192Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults192Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult192Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults256Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult256Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults256Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult256Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults384Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult384Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults384Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult384Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults512Value(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult512Value(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

func sumOwnedResults512Pointer(inputs []int) uint64 {
	var acc uint64
	for _, seed := range inputs {
		result := buildOwnedResult512Pointer(seed)
		acc += scoreHotFields(result.Core, result.Latency)
	}
	return acc
}

//go:noinline
func buildOwnedResult8Value(seed int) ownedResult8 {
	return ownedResult8{
		Core: requestCore8{
			AccountID: uint64(1_000 + seed),
		},
	}
}

//go:noinline
func buildOwnedResult8Pointer(seed int) *ownedResult8 {
	result := ownedResult8{
		Core: requestCore8{
			AccountID: uint64(1_000 + seed),
		},
	}
	return &result
}

//go:noinline
func buildOwnedResult16Value(seed int) ownedResult16 {
	return ownedResult16{
		Core: requestCore16{
			AccountID: uint64(1_000 + seed),
			Score:     uint32((seed%97)+1) * 17,
			Attempts:  uint16((seed % 32) + 1),
			Flags:     uint16((seed * 5) % 64),
		},
	}
}

//go:noinline
func buildOwnedResult16Pointer(seed int) *ownedResult16 {
	result := ownedResult16{
		Core: requestCore16{
			AccountID: uint64(1_000 + seed),
			Score:     uint32((seed%97)+1) * 17,
			Attempts:  uint16((seed % 32) + 1),
			Flags:     uint16((seed * 5) % 64),
		},
	}
	return &result
}

//go:noinline
func buildOwnedResult24Value(seed int) ownedResult24 {
	return ownedResult24{
		Core: requestCore24{
			AccountID: uint64(1_000 + seed),
			RequestID: uint64(100_000 + seed*3),
			Score:     uint32((seed%97)+1) * 17,
			Attempts:  uint32((seed % 32) + 1),
		},
	}
}

//go:noinline
func buildOwnedResult24Pointer(seed int) *ownedResult24 {
	result := ownedResult24{
		Core: requestCore24{
			AccountID: uint64(1_000 + seed),
			RequestID: uint64(100_000 + seed*3),
			Score:     uint32((seed%97)+1) * 17,
			Attempts:  uint32((seed % 32) + 1),
		},
	}
	return &result
}

//go:noinline
func buildOwnedResult32Value(seed int) ownedResult32 {
	return ownedResult32{
		Core: buildRequestCore(seed),
	}
}

//go:noinline
func buildOwnedResult32Pointer(seed int) *ownedResult32 {
	result := ownedResult32{
		Core: buildRequestCore(seed),
	}
	return &result
}

//go:noinline
func buildOwnedResult64Value(seed int) ownedResult64 {
	return ownedResult64{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
	}
}

//go:noinline
func buildOwnedResult64Pointer(seed int) *ownedResult64 {
	result := ownedResult64{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
	}
	return &result
}

//go:noinline
func buildOwnedResult96Value(seed int) ownedResult96 {
	return ownedResult96{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
}

//go:noinline
func buildOwnedResult96Pointer(seed int) *ownedResult96 {
	result := ownedResult96{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	return &result
}

//go:noinline
func buildOwnedResult128Value(seed int) ownedResult128 {
	return ownedResult128{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
		Window:  buildRequestWindow(seed, 0),
	}
}

//go:noinline
func buildOwnedResult128Pointer(seed int) *ownedResult128 {
	result := ownedResult128{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
		Window:  buildRequestWindow(seed, 0),
	}
	return &result
}

//go:noinline
func buildOwnedResult192Value(seed int) ownedResult192 {
	result := ownedResult192{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return result
}

//go:noinline
func buildOwnedResult192Pointer(seed int) *ownedResult192 {
	result := ownedResult192{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return &result
}

//go:noinline
func buildOwnedResult256Value(seed int) ownedResult256 {
	result := ownedResult256{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return result
}

//go:noinline
func buildOwnedResult256Pointer(seed int) *ownedResult256 {
	result := ownedResult256{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return &result
}

//go:noinline
func buildOwnedResult384Value(seed int) ownedResult384 {
	result := ownedResult384{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return result
}

//go:noinline
func buildOwnedResult384Pointer(seed int) *ownedResult384 {
	result := ownedResult384{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return &result
}

//go:noinline
func buildOwnedResult512Value(seed int) ownedResult512 {
	result := ownedResult512{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return result
}

//go:noinline
func buildOwnedResult512Pointer(seed int) *ownedResult512 {
	result := ownedResult512{
		Core:    buildRequestCore(seed),
		Latency: buildRequestLatency(seed),
		Profile: buildRequestProfile(seed),
	}
	for i := range result.Windows {
		result.Windows[i] = buildRequestWindow(seed, i)
	}
	return &result
}
