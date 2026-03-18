package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

const (
	paramBenchmarkRecords = 256
	paramTestRecords      = 32
)

type requestCore8 struct {
	AccountID uint64
}

type requestCore16 struct {
	AccountID uint64
	Score     uint32
	Attempts  uint16
	Flags     uint16
}

type requestCore24 struct {
	AccountID uint64
	RequestID uint64
	Score     uint32
	Attempts  uint32
}

type requestCore32 struct {
	AccountID uint64
	RequestID uint64
	Score     uint32
	Attempts  uint32
	RegionID  uint16
	Flags     uint16
	Priority  uint8
	Active    bool
	Shard     uint8
	Segment   uint8
}

type requestLatency32 struct {
	GeoHash     uint64
	LastSeenAt  uint64
	LatencyP50  float32
	LatencyP95  float32
	Risk        float32
	ErrorBudget float32
}

type requestProfile32 struct {
	SpendCents    uint64
	LifetimeValue uint64
	CountryCodes  [4]uint16
	FeatureMix    [4]uint8
	ErrorRate     float32
}

type requestWindow32 struct {
	Quantiles [4]float32
	Counts    [4]uint16
	Ratios    [4]uint8
	Version   uint16
	Status    uint8
	Tier      uint8
}

type scoringRequest8 struct {
	Core requestCore8
}

type scoringRequest16 struct {
	Core requestCore16
}

type scoringRequest24 struct {
	Core requestCore24
}

type scoringRequest32 struct {
	Core requestCore32
}

type scoringRequest64 struct {
	Core    requestCore32
	Latency requestLatency32
}

type scoringRequest96 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
}

type scoringRequest128 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
	Window  requestWindow32
}

type scoringRequest192 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
	Windows [3]requestWindow32
}

type scoringRequest256 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
	Windows [5]requestWindow32
}

type scoringRequest384 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
	Windows [9]requestWindow32
}

type scoringRequest512 struct {
	Core    requestCore32
	Latency requestLatency32
	Profile requestProfile32
	Windows [13]requestWindow32
}

var paramValueVsPointerSink uint64

func TestParamValueVsPointerOutputsMatch(t *testing.T) {
	t.Run("payload_sizes", func(t *testing.T) {
		require.Equal(t, uintptr(8), unsafe.Sizeof(scoringRequest8{}))
		require.Equal(t, uintptr(16), unsafe.Sizeof(scoringRequest16{}))
		require.Equal(t, uintptr(24), unsafe.Sizeof(scoringRequest24{}))
		require.Equal(t, uintptr(32), unsafe.Sizeof(scoringRequest32{}))
		require.Equal(t, uintptr(64), unsafe.Sizeof(scoringRequest64{}))
		require.Equal(t, uintptr(96), unsafe.Sizeof(scoringRequest96{}))
		require.Equal(t, uintptr(128), unsafe.Sizeof(scoringRequest128{}))
		require.Equal(t, uintptr(192), unsafe.Sizeof(scoringRequest192{}))
		require.Equal(t, uintptr(256), unsafe.Sizeof(scoringRequest256{}))
		require.Equal(t, uintptr(384), unsafe.Sizeof(scoringRequest384{}))
		require.Equal(t, uintptr(512), unsafe.Sizeof(scoringRequest512{}))
	})

	t.Run("size=8", func(t *testing.T) {
		requests := buildScoringRequests8(paramTestRecords)
		require.Equal(t, sumScoringRequests8Value(requests), sumScoringRequests8Pointer(requests))
	})
	t.Run("size=16", func(t *testing.T) {
		requests := buildScoringRequests16(paramTestRecords)
		require.Equal(t, sumScoringRequests16Value(requests), sumScoringRequests16Pointer(requests))
	})
	t.Run("size=24", func(t *testing.T) {
		requests := buildScoringRequests24(paramTestRecords)
		require.Equal(t, sumScoringRequests24Value(requests), sumScoringRequests24Pointer(requests))
	})
	t.Run("size=32", func(t *testing.T) {
		requests := buildScoringRequests32(paramTestRecords)
		require.Equal(t, sumScoringRequests32Value(requests), sumScoringRequests32Pointer(requests))
	})
	t.Run("size=64", func(t *testing.T) {
		requests := buildScoringRequests64(paramTestRecords)
		require.Equal(t, sumScoringRequests64Value(requests), sumScoringRequests64Pointer(requests))
	})
	t.Run("size=96", func(t *testing.T) {
		requests := buildScoringRequests96(paramTestRecords)
		require.Equal(t, sumScoringRequests96Value(requests), sumScoringRequests96Pointer(requests))
	})
	t.Run("size=128", func(t *testing.T) {
		requests := buildScoringRequests128(paramTestRecords)
		require.Equal(t, sumScoringRequests128Value(requests), sumScoringRequests128Pointer(requests))
	})
	t.Run("size=192", func(t *testing.T) {
		requests := buildScoringRequests192(paramTestRecords)
		require.Equal(t, sumScoringRequests192Value(requests), sumScoringRequests192Pointer(requests))
	})
	t.Run("size=256", func(t *testing.T) {
		requests := buildScoringRequests256(paramTestRecords)
		require.Equal(t, sumScoringRequests256Value(requests), sumScoringRequests256Pointer(requests))
	})
	t.Run("size=384", func(t *testing.T) {
		requests := buildScoringRequests384(paramTestRecords)
		require.Equal(t, sumScoringRequests384Value(requests), sumScoringRequests384Pointer(requests))
	})
	t.Run("size=512", func(t *testing.T) {
		requests := buildScoringRequests512(paramTestRecords)
		require.Equal(t, sumScoringRequests512Value(requests), sumScoringRequests512Pointer(requests))
	})
}

func BenchmarkParamValueVsPointer(b *testing.B) {
	cases := []struct {
		size         int
		valueBench   func(*testing.B)
		pointerBench func(*testing.B)
	}{
		{size: int(unsafe.Sizeof(scoringRequest8{})), valueBench: benchmarkScoringRequests8Value(), pointerBench: benchmarkScoringRequests8Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest16{})), valueBench: benchmarkScoringRequests16Value(), pointerBench: benchmarkScoringRequests16Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest24{})), valueBench: benchmarkScoringRequests24Value(), pointerBench: benchmarkScoringRequests24Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest32{})), valueBench: benchmarkScoringRequests32Value(), pointerBench: benchmarkScoringRequests32Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest64{})), valueBench: benchmarkScoringRequests64Value(), pointerBench: benchmarkScoringRequests64Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest96{})), valueBench: benchmarkScoringRequests96Value(), pointerBench: benchmarkScoringRequests96Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest128{})), valueBench: benchmarkScoringRequests128Value(), pointerBench: benchmarkScoringRequests128Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest192{})), valueBench: benchmarkScoringRequests192Value(), pointerBench: benchmarkScoringRequests192Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest256{})), valueBench: benchmarkScoringRequests256Value(), pointerBench: benchmarkScoringRequests256Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest384{})), valueBench: benchmarkScoringRequests384Value(), pointerBench: benchmarkScoringRequests384Pointer()},
		{size: int(unsafe.Sizeof(scoringRequest512{})), valueBench: benchmarkScoringRequests512Value(), pointerBench: benchmarkScoringRequests512Pointer()},
	}

	for _, tc := range cases {
		run := tc.valueBench
		switch variant {
		case "value_param":
		case "pointer_param":
			run = tc.pointerBench
		default:
			b.Fatalf("specify which benchmark to run: -args -variant value_param|pointer_param")
		}

		b.Run(fmt.Sprintf("size=%d", tc.size), run)
	}
}

func benchmarkScoringRequests8Value() func(*testing.B) {
	requests := buildScoringRequests8(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests8Value(requests)
		}
	}
}

func benchmarkScoringRequests8Pointer() func(*testing.B) {
	requests := buildScoringRequests8(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests8Pointer(requests)
		}
	}
}

func benchmarkScoringRequests16Value() func(*testing.B) {
	requests := buildScoringRequests16(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests16Value(requests)
		}
	}
}

func benchmarkScoringRequests16Pointer() func(*testing.B) {
	requests := buildScoringRequests16(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests16Pointer(requests)
		}
	}
}

func benchmarkScoringRequests24Value() func(*testing.B) {
	requests := buildScoringRequests24(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests24Value(requests)
		}
	}
}

func benchmarkScoringRequests24Pointer() func(*testing.B) {
	requests := buildScoringRequests24(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests24Pointer(requests)
		}
	}
}

func benchmarkScoringRequests32Value() func(*testing.B) {
	requests := buildScoringRequests32(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests32Value(requests)
		}
	}
}

func benchmarkScoringRequests32Pointer() func(*testing.B) {
	requests := buildScoringRequests32(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests32Pointer(requests)
		}
	}
}

func benchmarkScoringRequests64Value() func(*testing.B) {
	requests := buildScoringRequests64(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests64Value(requests)
		}
	}
}

func benchmarkScoringRequests64Pointer() func(*testing.B) {
	requests := buildScoringRequests64(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests64Pointer(requests)
		}
	}
}

func benchmarkScoringRequests96Value() func(*testing.B) {
	requests := buildScoringRequests96(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests96Value(requests)
		}
	}
}

func benchmarkScoringRequests96Pointer() func(*testing.B) {
	requests := buildScoringRequests96(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests96Pointer(requests)
		}
	}
}

func benchmarkScoringRequests128Value() func(*testing.B) {
	requests := buildScoringRequests128(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests128Value(requests)
		}
	}
}

func benchmarkScoringRequests128Pointer() func(*testing.B) {
	requests := buildScoringRequests128(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests128Pointer(requests)
		}
	}
}

func benchmarkScoringRequests192Value() func(*testing.B) {
	requests := buildScoringRequests192(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests192Value(requests)
		}
	}
}

func benchmarkScoringRequests192Pointer() func(*testing.B) {
	requests := buildScoringRequests192(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests192Pointer(requests)
		}
	}
}

func benchmarkScoringRequests256Value() func(*testing.B) {
	requests := buildScoringRequests256(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests256Value(requests)
		}
	}
}

func benchmarkScoringRequests256Pointer() func(*testing.B) {
	requests := buildScoringRequests256(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests256Pointer(requests)
		}
	}
}

func benchmarkScoringRequests384Value() func(*testing.B) {
	requests := buildScoringRequests384(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests384Value(requests)
		}
	}
}

func benchmarkScoringRequests384Pointer() func(*testing.B) {
	requests := buildScoringRequests384(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests384Pointer(requests)
		}
	}
}

func benchmarkScoringRequests512Value() func(*testing.B) {
	requests := buildScoringRequests512(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests512Value(requests)
		}
	}
}

func benchmarkScoringRequests512Pointer() func(*testing.B) {
	requests := buildScoringRequests512(paramBenchmarkRecords)
	return func(b *testing.B) {
		for b.Loop() {
			paramValueVsPointerSink = sumScoringRequests512Pointer(requests)
		}
	}
}

func buildScoringRequests8(size int) []scoringRequest8 {
	requests := make([]scoringRequest8, size)
	for i := range requests {
		requests[i] = scoringRequest8{
			Core: requestCore8{
				AccountID: uint64(1_000 + i),
			},
		}
	}
	return requests
}

func buildScoringRequests16(size int) []scoringRequest16 {
	requests := make([]scoringRequest16, size)
	for i := range requests {
		requests[i] = scoringRequest16{
			Core: requestCore16{
				AccountID: uint64(1_000 + i),
				Score:     uint32((i%97)+1) * 17,
				Attempts:  uint16((i % 32) + 1),
				Flags:     uint16((i * 5) % 64),
			},
		}
	}
	return requests
}

func buildScoringRequests24(size int) []scoringRequest24 {
	requests := make([]scoringRequest24, size)
	for i := range requests {
		requests[i] = scoringRequest24{
			Core: requestCore24{
				AccountID: uint64(1_000 + i),
				RequestID: uint64(100_000 + i*3),
				Score:     uint32((i%97)+1) * 17,
				Attempts:  uint32((i % 32) + 1),
			},
		}
	}
	return requests
}

func buildScoringRequests32(size int) []scoringRequest32 {
	requests := make([]scoringRequest32, size)
	for i := range requests {
		requests[i] = scoringRequest32{
			Core: buildRequestCore(i),
		}
	}
	return requests
}

func buildScoringRequests64(size int) []scoringRequest64 {
	requests := make([]scoringRequest64, size)
	for i := range requests {
		requests[i] = scoringRequest64{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
		}
	}
	return requests
}

func buildScoringRequests96(size int) []scoringRequest96 {
	requests := make([]scoringRequest96, size)
	for i := range requests {
		requests[i] = scoringRequest96{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
		}
	}
	return requests
}

func buildScoringRequests128(size int) []scoringRequest128 {
	requests := make([]scoringRequest128, size)
	for i := range requests {
		requests[i] = scoringRequest128{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
			Window:  buildRequestWindow(i, 0),
		}
	}
	return requests
}

func buildScoringRequests192(size int) []scoringRequest192 {
	requests := make([]scoringRequest192, size)
	for i := range requests {
		requests[i] = scoringRequest192{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
		}
		for j := range requests[i].Windows {
			requests[i].Windows[j] = buildRequestWindow(i, j)
		}
	}
	return requests
}

func buildScoringRequests256(size int) []scoringRequest256 {
	requests := make([]scoringRequest256, size)
	for i := range requests {
		requests[i] = scoringRequest256{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
		}
		for j := range requests[i].Windows {
			requests[i].Windows[j] = buildRequestWindow(i, j)
		}
	}
	return requests
}

func buildScoringRequests384(size int) []scoringRequest384 {
	requests := make([]scoringRequest384, size)
	for i := range requests {
		requests[i] = scoringRequest384{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
		}
		for j := range requests[i].Windows {
			requests[i].Windows[j] = buildRequestWindow(i, j)
		}
	}
	return requests
}

func buildScoringRequests512(size int) []scoringRequest512 {
	requests := make([]scoringRequest512, size)
	for i := range requests {
		requests[i] = scoringRequest512{
			Core:    buildRequestCore(i),
			Latency: buildRequestLatency(i),
			Profile: buildRequestProfile(i),
		}
		for j := range requests[i].Windows {
			requests[i].Windows[j] = buildRequestWindow(i, j)
		}
	}
	return requests
}

func buildRequestCore(i int) requestCore32 {
	return requestCore32{
		AccountID: uint64(1_000 + i),
		RequestID: uint64(100_000 + i*3),
		Score:     uint32((i%97)+1) * 17,
		Attempts:  uint32((i % 32) + 1),
		RegionID:  uint16(i % 256),
		Flags:     uint16((i * 5) % 64),
		Priority:  uint8(i % 8),
		Active:    i%5 != 0,
		Shard:     uint8((i / 7) % 16),
		Segment:   uint8((i / 3) % 16),
	}
}

func buildRequestLatency(i int) requestLatency32 {
	return requestLatency32{
		GeoHash:     uint64(5_000_000 + i*19),
		LastSeenAt:  uint64(1_700_000_000 + i*60),
		LatencyP50:  2.5 + float32(i%23)*0.25,
		LatencyP95:  8.0 + float32(i%37)*0.5,
		Risk:        0.25 + float32(i%11)*0.125,
		ErrorBudget: 1.0 + float32(i%13)*0.25,
	}
}

func buildRequestProfile(i int) requestProfile32 {
	return requestProfile32{
		SpendCents:    uint64(10_000 + i*13),
		LifetimeValue: uint64(100_000 + i*29),
		CountryCodes: [4]uint16{
			uint16((i + 1) % 250),
			uint16((i + 7) % 250),
			uint16((i + 13) % 250),
			uint16((i + 29) % 250),
		},
		FeatureMix: [4]uint8{
			uint8(i % 11),
			uint8((i + 1) % 11),
			uint8((i + 2) % 11),
			uint8((i + 3) % 11),
		},
		ErrorRate: 0.5 + float32(i%9)*0.25,
	}
}

func buildRequestWindow(i, j int) requestWindow32 {
	base := i*31 + j*17
	return requestWindow32{
		Quantiles: [4]float32{
			10.0 + float32(base%17)*0.5,
			20.0 + float32(base%19)*0.75,
			30.0 + float32(base%23)*1.0,
			40.0 + float32(base%29)*1.25,
		},
		Counts: [4]uint16{
			uint16(10 + base%100),
			uint16(20 + (base*3)%100),
			uint16(30 + (base*5)%100),
			uint16(40 + (base*7)%100),
		},
		Ratios: [4]uint8{
			uint8(base % 7),
			uint8((base + 1) % 7),
			uint8((base + 2) % 7),
			uint8((base + 3) % 7),
		},
		Version: uint16((base % 32) + 1),
		Status:  uint8(base % 5),
		Tier:    uint8(base % 4),
	}
}

func sumScoringRequests8Value(requests []scoringRequest8) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest8Value(requests[i])
	}
	return acc
}

func sumScoringRequests8Pointer(requests []scoringRequest8) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest8Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests16Value(requests []scoringRequest16) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest16Value(requests[i])
	}
	return acc
}

func sumScoringRequests16Pointer(requests []scoringRequest16) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest16Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests24Value(requests []scoringRequest24) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest24Value(requests[i])
	}
	return acc
}

func sumScoringRequests24Pointer(requests []scoringRequest24) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest24Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests32Value(requests []scoringRequest32) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest32Value(requests[i])
	}
	return acc
}

func sumScoringRequests32Pointer(requests []scoringRequest32) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest32Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests64Value(requests []scoringRequest64) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest64Value(requests[i])
	}
	return acc
}

func sumScoringRequests64Pointer(requests []scoringRequest64) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest64Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests96Value(requests []scoringRequest96) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest96Value(requests[i])
	}
	return acc
}

func sumScoringRequests96Pointer(requests []scoringRequest96) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest96Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests128Value(requests []scoringRequest128) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest128Value(requests[i])
	}
	return acc
}

func sumScoringRequests128Pointer(requests []scoringRequest128) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest128Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests192Value(requests []scoringRequest192) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest192Value(requests[i])
	}
	return acc
}

func sumScoringRequests192Pointer(requests []scoringRequest192) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest192Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests256Value(requests []scoringRequest256) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest256Value(requests[i])
	}
	return acc
}

func sumScoringRequests256Pointer(requests []scoringRequest256) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest256Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests384Value(requests []scoringRequest384) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest384Value(requests[i])
	}
	return acc
}

func sumScoringRequests384Pointer(requests []scoringRequest384) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest384Pointer(&requests[i])
	}
	return acc
}

func sumScoringRequests512Value(requests []scoringRequest512) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest512Value(requests[i])
	}
	return acc
}

func sumScoringRequests512Pointer(requests []scoringRequest512) uint64 {
	var acc uint64
	for i := range requests {
		acc += scoreRequest512Pointer(&requests[i])
	}
	return acc
}

//go:noinline
func scoreRequest8Value(request scoringRequest8) uint64 {
	return request.Core.AccountID
}

//go:noinline
func scoreRequest8Pointer(request *scoringRequest8) uint64 {
	return request.Core.AccountID
}

//go:noinline
func scoreRequest16Value(request scoringRequest16) uint64 {
	acc := request.Core.AccountID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	acc += uint64(request.Core.Flags)
	return acc
}

//go:noinline
func scoreRequest16Pointer(request *scoringRequest16) uint64 {
	acc := request.Core.AccountID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	acc += uint64(request.Core.Flags)
	return acc
}

//go:noinline
func scoreRequest24Value(request scoringRequest24) uint64 {
	acc := request.Core.AccountID ^ request.Core.RequestID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	return acc
}

//go:noinline
func scoreRequest24Pointer(request *scoringRequest24) uint64 {
	acc := request.Core.AccountID ^ request.Core.RequestID
	acc += uint64(request.Core.Score)
	acc += uint64(request.Core.Attempts)
	return acc
}

//go:noinline
func scoreRequest32Value(request scoringRequest32) uint64 {
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func scoreRequest32Pointer(request *scoringRequest32) uint64 {
	return scoreHotFields(request.Core, requestLatency32{})
}

//go:noinline
func scoreRequest64Value(request scoringRequest64) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest64Pointer(request *scoringRequest64) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest96Value(request scoringRequest96) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest96Pointer(request *scoringRequest96) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest128Value(request scoringRequest128) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest128Pointer(request *scoringRequest128) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest192Value(request scoringRequest192) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest192Pointer(request *scoringRequest192) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest256Value(request scoringRequest256) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest256Pointer(request *scoringRequest256) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest384Value(request scoringRequest384) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest384Pointer(request *scoringRequest384) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest512Value(request scoringRequest512) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

//go:noinline
func scoreRequest512Pointer(request *scoringRequest512) uint64 {
	return scoreHotFields(request.Core, request.Latency)
}

func scoreHotFields(core requestCore32, latency requestLatency32) uint64 {
	acc := core.AccountID ^ core.RequestID
	acc += uint64(core.Score)
	acc += uint64(core.Attempts)
	acc += uint64(core.RegionID)
	acc += uint64(core.Flags)
	acc += uint64(core.Priority)
	acc += uint64(core.Shard)
	acc += uint64(core.Segment)
	acc += latency.GeoHash
	acc += uint64(latency.LastSeenAt & 0xffff)
	acc += uint64(latency.LatencyP50 * 100)
	acc += uint64(latency.LatencyP95 * 100)
	acc += uint64(latency.Risk * 1_000)
	acc += uint64(latency.ErrorBudget * 1_000)
	if !core.Active {
		acc ^= 0x9e3779b97f4a7c15
	}
	return acc
}
