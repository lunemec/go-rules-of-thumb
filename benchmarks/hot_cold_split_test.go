package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var hotColdSplitSizes = []int{1, 10, 50, 100, 500, 1_000, 5_000, 10_000, 50_000, 100_000}

type hotColdSplitBench func(size int) func(*testing.B)

type inlineHotColdRecord struct {
	ID      uint64
	Score   int64
	Visits  uint32
	Flags   uint16
	Tier    uint16
	Active  bool
	Name    [24]byte
	Email   [32]byte
	Region  [16]byte
	Notes   [160]byte
	History [128]byte
	Profile [64]byte
	Audit   [64]byte
}

type splitHotRecord struct {
	ID     uint64
	Score  int64
	Visits uint32
	Flags  uint16
	Tier   uint16
	Active bool
	Cold   *splitColdRecord
}

type splitColdRecord struct {
	Name    [24]byte
	Email   [32]byte
	Region  [16]byte
	Notes   [160]byte
	History [128]byte
	Profile [64]byte
	Audit   [64]byte
}

type hotColdSnapshot struct {
	ID      uint64
	Score   int64
	Visits  uint32
	Flags   uint16
	Tier    uint16
	Active  bool
	Name    [24]byte
	Email   [32]byte
	Region  [16]byte
	Notes   [160]byte
	History [128]byte
	Profile [64]byte
	Audit   [64]byte
}

var (
	hotColdScanSink     uint64
	hotColdSnapshotSink []hotColdSnapshot
)

func TestHotColdSplitOutputsMatch(t *testing.T) {
	for _, size := range sizesReduced {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			inline, split := buildHotColdLayouts(size)

			require.Equal(t, scanInlineHot(inline), scanSplitHot(split))
			require.Equal(t, buildInlineSnapshots(inline), buildSplitSnapshots(split))
		})
	}
}

func BenchmarkHotColdSplit(b *testing.B) {
	runBenchmark := func(runF hotColdSplitBench) {
		for _, size := range hotColdSplitSizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				runF(size),
			)
		}
	}

	switch variant {
	case "inline_hot_scan":
		runBenchmark(benchmarkInlineHotScan)
	case "split_hot_scan":
		runBenchmark(benchmarkSplitHotScan)
	case "inline_snapshot":
		runBenchmark(benchmarkInlineSnapshot)
	case "split_snapshot":
		runBenchmark(benchmarkSplitSnapshot)
	default:
		b.Errorf("specify which test to run: -args -variant inline_hot_scan|split_hot_scan|inline_snapshot|split_snapshot")
	}
}

func benchmarkInlineHotScan(size int) func(*testing.B) {
	inline, _ := buildHotColdLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotColdScanSink = scanInlineHot(inline)
		}
	}
}

func benchmarkSplitHotScan(size int) func(*testing.B) {
	_, split := buildHotColdLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotColdScanSink = scanSplitHot(split)
		}
	}
}

func benchmarkInlineSnapshot(size int) func(*testing.B) {
	inline, _ := buildHotColdLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotColdSnapshotSink = buildInlineSnapshots(inline)
		}
	}
}

func benchmarkSplitSnapshot(size int) func(*testing.B) {
	_, split := buildHotColdLayouts(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotColdSnapshotSink = buildSplitSnapshots(split)
		}
	}
}

func buildHotColdLayouts(size int) ([]inlineHotColdRecord, []splitHotRecord) {
	inline := make([]inlineHotColdRecord, size)
	split := make([]splitHotRecord, size)
	// Keep cold targets contiguous so the benchmark isolates the tradeoff
	// between a narrower hot record and the extra pointer indirection.
	coldBacking := make([]splitColdRecord, size)

	for i := range inline {
		hot := splitHotRecord{
			ID:     uint64(i + 1),
			Score:  int64((i%97)+1) * 29,
			Visits: uint32((i % 64) + 1),
			Flags:  uint16(i % 13),
			Tier:   uint16((i / 5) % 4),
			Active: i%5 != 0,
		}
		cold := splitColdRecord{}

		fillPattern(cold.Name[:], byte(i+11))
		fillPattern(cold.Email[:], byte(i+37))
		fillPattern(cold.Region[:], byte(i+59))
		fillPattern(cold.Notes[:], byte(i+83))
		fillPattern(cold.History[:], byte(i+101))
		fillPattern(cold.Profile[:], byte(i+149))
		fillPattern(cold.Audit[:], byte(i+181))

		coldBacking[i] = cold
		hot.Cold = &coldBacking[i]
		split[i] = hot

		inline[i] = inlineHotColdRecord{
			ID:      hot.ID,
			Score:   hot.Score,
			Visits:  hot.Visits,
			Flags:   hot.Flags,
			Tier:    hot.Tier,
			Active:  hot.Active,
			Name:    cold.Name,
			Email:   cold.Email,
			Region:  cold.Region,
			Notes:   cold.Notes,
			History: cold.History,
			Profile: cold.Profile,
			Audit:   cold.Audit,
		}
	}

	return inline, split
}

func scanInlineHot(records []inlineHotColdRecord) uint64 {
	var acc uint64

	for i := range records {
		acc += records[i].ID
		acc += uint64(records[i].Visits)
		acc += uint64(records[i].Score)
		acc += uint64(records[i].Flags)
		acc += uint64(records[i].Tier)
		if records[i].Active {
			acc++
		}
	}

	return acc
}

func scanSplitHot(records []splitHotRecord) uint64 {
	var acc uint64

	for i := range records {
		acc += records[i].ID
		acc += uint64(records[i].Visits)
		acc += uint64(records[i].Score)
		acc += uint64(records[i].Flags)
		acc += uint64(records[i].Tier)
		if records[i].Active {
			acc++
		}
	}

	return acc
}

func buildInlineSnapshots(records []inlineHotColdRecord) []hotColdSnapshot {
	out := make([]hotColdSnapshot, len(records))

	for i := range records {
		out[i] = hotColdSnapshot{
			ID:      records[i].ID,
			Score:   records[i].Score,
			Visits:  records[i].Visits,
			Flags:   records[i].Flags,
			Tier:    records[i].Tier,
			Active:  records[i].Active,
			Name:    records[i].Name,
			Email:   records[i].Email,
			Region:  records[i].Region,
			Notes:   records[i].Notes,
			History: records[i].History,
			Profile: records[i].Profile,
			Audit:   records[i].Audit,
		}
	}

	return out
}

func buildSplitSnapshots(records []splitHotRecord) []hotColdSnapshot {
	out := make([]hotColdSnapshot, len(records))

	for i := range records {
		cold := records[i].Cold
		out[i] = hotColdSnapshot{
			ID:      records[i].ID,
			Score:   records[i].Score,
			Visits:  records[i].Visits,
			Flags:   records[i].Flags,
			Tier:    records[i].Tier,
			Active:  records[i].Active,
			Name:    cold.Name,
			Email:   cold.Email,
			Region:  cold.Region,
			Notes:   cold.Notes,
			History: cold.History,
			Profile: cold.Profile,
			Audit:   cold.Audit,
		}
	}

	return out
}
