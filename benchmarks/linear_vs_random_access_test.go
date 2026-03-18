package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

var linearVsRandomAccessSizes = []int{64, 128, 256, 512, 1_024, 2_048, 4_096, 8_192, 16_384, 32_768, 65_536}

type linearVsRandomAccessBench func(records int) func(*testing.B)

// Keep each record to one cache line so the benchmark mostly reflects how
// traversal order interacts with spatial locality and hardware prefetch.
type localityRecord struct {
	Hot  uint64
	Warm uint64
	Pad  [6]uint64
}

var localityAccessSink uint64

func TestLinearVsRandomAccessOutputsMatch(t *testing.T) {
	for _, records := range sizesReduced {
		t.Run(fmt.Sprintf("records=%d", records), func(t *testing.T) {
			data, linearOrder, randomOrder := buildLinearVsRandomAccessInput(records)

			require.Equal(t, scanRecordsByOrder(data, linearOrder), scanRecordsByOrder(data, randomOrder))
			if records > 2 {
				require.NotEqual(t, linearOrder, randomOrder)
			}
		})
	}
}

func BenchmarkLinearVsRandomAccess(b *testing.B) {
	runBenchmark := func(runF linearVsRandomAccessBench) {
		for _, records := range linearVsRandomAccessSizes {
			b.Run(
				fmt.Sprintf("records=%d", records),
				runF(records),
			)
		}
	}

	switch variant {
	case "linear":
		runBenchmark(benchmarkLinearAccess)
	case "random":
		runBenchmark(benchmarkRandomAccess)
	default:
		b.Errorf("specify which test to run: -args -variant linear|random")
	}
}

func benchmarkLinearAccess(records int) func(*testing.B) {
	data, linearOrder, _ := buildLinearVsRandomAccessInput(records)
	return func(b *testing.B) {
		for b.Loop() {
			localityAccessSink = scanRecordsByOrder(data, linearOrder)
		}
	}
}

func benchmarkRandomAccess(records int) func(*testing.B) {
	data, _, randomOrder := buildLinearVsRandomAccessInput(records)
	return func(b *testing.B) {
		for b.Loop() {
			localityAccessSink = scanRecordsByOrder(data, randomOrder)
		}
	}
}

func buildLinearVsRandomAccessInput(records int) ([]localityRecord, []int, []int) {
	data := make([]localityRecord, records)
	linearOrder := make([]int, records)

	for i := range data {
		data[i] = localityRecord{
			Hot:  uint64(i+1) * 17,
			Warm: uint64((i%97)+1) * 131,
			Pad: [6]uint64{
				uint64(i*3 + 1),
				uint64(i*5 + 3),
				uint64(i*7 + 5),
				uint64(i*11 + 7),
				uint64(i*13 + 11),
				uint64(i*17 + 13),
			},
		}
		linearOrder[i] = i
	}

	randomOrder := copySlice(linearOrder)
	rng := rand.New(rand.NewSource(0x5eed))
	rng.Shuffle(len(randomOrder), func(i, j int) {
		randomOrder[i], randomOrder[j] = randomOrder[j], randomOrder[i]
	})

	return data, linearOrder, randomOrder
}

func scanRecordsByOrder(data []localityRecord, order []int) uint64 {
	var acc uint64

	for i := range order {
		record := data[order[i]]
		acc += record.Hot
		acc += record.Warm
	}

	return acc
}
