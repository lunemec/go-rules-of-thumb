package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var valuesVsPointersSizes = []int{10, 100, 1_000, 10_000, 100_000}

type valuesVsPointersBench func(size int) func(*testing.B)

type userRecord struct {
	ID     uint64
	Score  int64
	Visits uint32
	Active bool
	Name   [24]byte
	Email  [32]byte
	Region [16]byte
	Notes  [48]byte
}

type userSnapshot struct {
	ID     uint64
	Score  int64
	Visits uint32
	Name   [24]byte
	Email  [32]byte
	Region [16]byte
	Notes  [48]byte
}

var (
	hotUserScanSink  uint64
	userSnapshotSink []userSnapshot
)

func TestValuesVsPointersOutputsMatch(t *testing.T) {
	for _, size := range sizesReduced {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			values, pointers := buildUserCollections(size)

			require.Equal(t, scanHotValues(values), scanHotPointers(pointers))
			require.Equal(t, buildSnapshotsValues(values), buildSnapshotsPointers(pointers))
		})
	}
}

func BenchmarkValuesVsPointers(b *testing.B) {
	runBenchmark := func(runF valuesVsPointersBench) {
		for _, size := range valuesVsPointersSizes {
			b.Run(
				fmt.Sprintf("size=%d", size),
				runF(size),
			)
		}
	}

	switch variant {
	case "values_hot_scan":
		runBenchmark(benchmarkValuesHotScan)
	case "pointers_hot_scan":
		runBenchmark(benchmarkPointersHotScan)
	case "values_snapshot":
		runBenchmark(benchmarkValuesSnapshot)
	case "pointers_snapshot":
		runBenchmark(benchmarkPointersSnapshot)
	default:
		b.Errorf("specify which test to run: -args -variant values_hot_scan|pointers_hot_scan|values_snapshot|pointers_snapshot")
	}
}

func benchmarkValuesHotScan(size int) func(*testing.B) {
	values, _ := buildUserCollections(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotUserScanSink = scanHotValues(values)
		}
	}
}

func benchmarkPointersHotScan(size int) func(*testing.B) {
	_, pointers := buildUserCollections(size)
	return func(b *testing.B) {
		for b.Loop() {
			hotUserScanSink = scanHotPointers(pointers)
		}
	}
}

func benchmarkValuesSnapshot(size int) func(*testing.B) {
	values, _ := buildUserCollections(size)
	return func(b *testing.B) {
		for b.Loop() {
			userSnapshotSink = buildSnapshotsValues(values)
		}
	}
}

func benchmarkPointersSnapshot(size int) func(*testing.B) {
	_, pointers := buildUserCollections(size)
	return func(b *testing.B) {
		for b.Loop() {
			userSnapshotSink = buildSnapshotsPointers(pointers)
		}
	}
}

func buildUserCollections(size int) ([]userRecord, []*userRecord) {
	values := make([]userRecord, size)
	// Keep the pointer targets contiguous so this benchmark isolates the cost
	// of the extra pointer slice and indirection rather than heap fragmentation.
	pointerBacking := make([]userRecord, size)
	pointers := make([]*userRecord, size)

	for i := range values {
		record := userRecord{
			ID:     uint64(i + 1),
			Score:  int64((i%97)+1) * 17,
			Visits: uint32((i % 64) + 1),
			Active: i%4 != 0,
		}

		fillPattern(record.Name[:], byte(i+11))
		fillPattern(record.Email[:], byte(i+37))
		fillPattern(record.Region[:], byte(i+59))
		fillPattern(record.Notes[:], byte(i+83))

		values[i] = record
		pointerBacking[i] = record
		pointers[i] = &pointerBacking[i]
	}

	return values, pointers
}

func scanHotValues(values []userRecord) uint64 {
	var acc uint64

	for i := range values {
		if !values[i].Active {
			continue
		}

		acc += values[i].ID
		acc += uint64(values[i].Visits)
		acc += uint64(values[i].Score)
	}

	return acc
}

func scanHotPointers(pointers []*userRecord) uint64 {
	var acc uint64

	for i := range pointers {
		if !pointers[i].Active {
			continue
		}

		acc += pointers[i].ID
		acc += uint64(pointers[i].Visits)
		acc += uint64(pointers[i].Score)
	}

	return acc
}

func buildSnapshotsValues(values []userRecord) []userSnapshot {
	out := make([]userSnapshot, 0, len(values))

	for i := range values {
		if !values[i].Active {
			continue
		}

		out = append(out, userSnapshot{
			ID:     values[i].ID,
			Score:  values[i].Score,
			Visits: values[i].Visits,
			Name:   values[i].Name,
			Email:  values[i].Email,
			Region: values[i].Region,
			Notes:  values[i].Notes,
		})
	}

	return out
}

func buildSnapshotsPointers(pointers []*userRecord) []userSnapshot {
	out := make([]userSnapshot, 0, len(pointers))

	for i := range pointers {
		if !pointers[i].Active {
			continue
		}

		out = append(out, userSnapshot{
			ID:     pointers[i].ID,
			Score:  pointers[i].Score,
			Visits: pointers[i].Visits,
			Name:   pointers[i].Name,
			Email:  pointers[i].Email,
			Region: pointers[i].Region,
			Notes:  pointers[i].Notes,
		})
	}

	return out
}
