package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

var buildSliceLengths = []int{10, 50, 100, 500, 1_000, 5_000, 10_000, 50_000, 100_000}

type buildSliceElem8 struct {
	Word uint64
}

type buildSliceElem32 struct {
	Words [4]uint64
}

type buildSliceElem128 struct {
	Words [16]uint64
}

var buildSliceSink any

func TestBuildSliceOutputsMatch(t *testing.T) {
	require.Equal(t, uintptr(8), unsafe.Sizeof(buildSliceElem8{}))
	require.Equal(t, uintptr(32), unsafe.Sizeof(buildSliceElem32{}))
	require.Equal(t, uintptr(128), unsafe.Sizeof(buildSliceElem128{}))

	for _, length := range sizesReduced {
		t.Run(fmt.Sprintf("length=%d", length), func(t *testing.T) {
			requireBuildSliceEquivalent(t, length, makeBuildSliceElem8)
			requireBuildSliceEquivalent(t, length, makeBuildSliceElem32)
			requireBuildSliceEquivalent(t, length, makeBuildSliceElem128)
		})
	}
}

func requireBuildSliceEquivalent[T comparable](t *testing.T, length int, build func(int) T) {
	t.Helper()

	seeds := buildSliceSeeds(length)

	growth := buildSliceAppendGrowth(seeds, build)
	prealloc := buildSliceAppendPrealloc(seeds, build)
	indexed := buildSliceIndexPresized(seeds, build)

	require.Equal(t, growth, prealloc)
	require.Equal(t, growth, indexed)
}

func BenchmarkBuildSlice(b *testing.B) {
	cases := []struct {
		size         int
		appendGrowth func(int) func(*testing.B)
		appendPre    func(int) func(*testing.B)
		indexed      func(int) func(*testing.B)
	}{
		{
			size:         int(unsafe.Sizeof(buildSliceElem8{})),
			appendGrowth: benchmarkBuildSliceAppendGrowth(makeBuildSliceElem8),
			appendPre:    benchmarkBuildSliceAppendPrealloc(makeBuildSliceElem8),
			indexed:      benchmarkBuildSliceIndexPresized(makeBuildSliceElem8),
		},
		{
			size:         int(unsafe.Sizeof(buildSliceElem32{})),
			appendGrowth: benchmarkBuildSliceAppendGrowth(makeBuildSliceElem32),
			appendPre:    benchmarkBuildSliceAppendPrealloc(makeBuildSliceElem32),
			indexed:      benchmarkBuildSliceIndexPresized(makeBuildSliceElem32),
		},
		{
			size:         int(unsafe.Sizeof(buildSliceElem128{})),
			appendGrowth: benchmarkBuildSliceAppendGrowth(makeBuildSliceElem128),
			appendPre:    benchmarkBuildSliceAppendPrealloc(makeBuildSliceElem128),
			indexed:      benchmarkBuildSliceIndexPresized(makeBuildSliceElem128),
		},
	}

	for _, tc := range cases {
		run := tc.appendGrowth
		switch variant {
		case "append_growth":
		case "append_prealloc":
			run = tc.appendPre
		case "index_presized":
			run = tc.indexed
		default:
			b.Fatalf("specify which benchmark to run: -args -variant append_growth|append_prealloc|index_presized")
		}

		for _, length := range buildSliceLengths {
			b.Run(fmt.Sprintf("size=%d length=%d", tc.size, length), run(length))
		}
	}
}

func benchmarkBuildSliceAppendGrowth[T any](build func(int) T) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		seeds := buildSliceSeeds(length)
		return func(b *testing.B) {
			for b.Loop() {
				buildSliceSink = buildSliceAppendGrowth(seeds, build)
			}
		}
	}
}

func benchmarkBuildSliceAppendPrealloc[T any](build func(int) T) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		seeds := buildSliceSeeds(length)
		return func(b *testing.B) {
			for b.Loop() {
				buildSliceSink = buildSliceAppendPrealloc(seeds, build)
			}
		}
	}
}

func benchmarkBuildSliceIndexPresized[T any](build func(int) T) func(int) func(*testing.B) {
	return func(length int) func(*testing.B) {
		seeds := buildSliceSeeds(length)
		return func(b *testing.B) {
			for b.Loop() {
				buildSliceSink = buildSliceIndexPresized(seeds, build)
			}
		}
	}
}

func buildSliceSeeds(length int) []int {
	seeds := make([]int, length)
	for i := range seeds {
		seeds[i] = i
	}
	return seeds
}

func buildSliceAppendGrowth[T any](seeds []int, build func(int) T) []T {
	var out []T
	for _, seed := range seeds {
		out = append(out, build(seed))
	}
	return out
}

func buildSliceAppendPrealloc[T any](seeds []int, build func(int) T) []T {
	out := make([]T, 0, len(seeds))
	for _, seed := range seeds {
		out = append(out, build(seed))
	}
	return out
}

func buildSliceIndexPresized[T any](seeds []int, build func(int) T) []T {
	out := make([]T, len(seeds))
	for i, seed := range seeds {
		out[i] = build(seed)
	}
	return out
}

func makeBuildSliceElem8(seed int) buildSliceElem8 {
	return buildSliceElem8{
		Word: uint64(seed+1) * 17,
	}
}

func makeBuildSliceElem32(seed int) buildSliceElem32 {
	base := uint64(seed + 1)
	return buildSliceElem32{
		Words: [4]uint64{
			base * 17,
			base * 19,
			base * 23,
			base * 29,
		},
	}
}

func makeBuildSliceElem128(seed int) buildSliceElem128 {
	base := uint64(seed + 1)
	elem := buildSliceElem128{}
	for i := range elem.Words {
		elem.Words[i] = base * uint64((i+1)*13)
	}
	return elem
}
